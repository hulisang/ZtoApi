package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hulisang/ZtoApi/register"
)

// 配置变量（从环境变量读取）
var (
	UPSTREAM_URL      string
	DEFAULT_KEY       string
	ZAI_TOKEN         string
	MODEL_NAME        string
	PORT              string
	DEBUG_MODE        bool
	DEFAULT_STREAM    bool
	DASHBOARD_ENABLED bool
	ENABLE_THINKING   bool
)

// 请求统计信息
type RequestStats struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	LastRequestTime     time.Time
	AverageResponseTime time.Duration
}

// 实时请求信息
type LiveRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Duration  int64     `json:"duration"`
	UserAgent string    `json:"user_agent"`
}

// 全局变量
var (
	stats         RequestStats
	liveRequests  = []LiveRequest{} // 初始化为空数组，而不是 nil
	statsMutex    sync.Mutex
	requestsMutex sync.Mutex
)

// 思考内容处理策略
const (
	THINK_TAGS_MODE = "strip" // strip: 去除<details>标签；think: 转为<think>标签；raw: 保留原样
)

// 系统配置常量
const (
	MAX_LIVE_REQUESTS      = 100        // 最多保留的实时请求记录数
	AUTH_TOKEN_TIMEOUT     = 10         // 获取匿名token的超时时间（秒）
	UPSTREAM_TIMEOUT       = 60         // 上游API调用超时时间（秒）
	TOKEN_DISPLAY_LENGTH   = 10         // token显示时的截取长度
	NANOSECONDS_TO_SECONDS = 1000000000 // 纳秒转秒的倍数
)

// 伪装前端头部（2025-09-30 更新：修复426错误）
const (
	X_FE_VERSION   = "prod-fe-1.0.94"                                                                                                  // 更新：1.0.70 → 1.0.94
	BROWSER_UA     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36" // 更新：Chrome 139 → 140
	SEC_CH_UA      = "\"Chromium\";v=\"140\", \"Not=A?Brand\";v=\"24\", \"Google Chrome\";v=\"140\""                                   // 更新：Chrome 140
	SEC_CH_UA_MOB  = "?0"
	SEC_CH_UA_PLAT = "\"Windows\""
	ORIGIN_BASE    = "https://chat.z.ai"
)

// 匿名token开关
const ANON_TOKEN_ENABLED = true

// 从环境变量初始化配置
func initConfig() {
	// 加载 .env.local 文件（如果存在）
	loadEnvFile(".env.local")
	// 也尝试加载标准的 .env 文件
	loadEnvFile(".env")

	UPSTREAM_URL = getEnv("UPSTREAM_URL", "https://chat.z.ai/api/chat/completions")
	DEFAULT_KEY = getEnv("DEFAULT_KEY", "sk-your-key")
	ZAI_TOKEN = getEnv("ZAI_TOKEN", "")
	MODEL_NAME = getEnv("MODEL_NAME", "GLM-4.6")
	PORT = getEnv("PORT", "9090")

	// 处理PORT格式，确保有冒号前缀
	if !strings.HasPrefix(PORT, ":") {
		PORT = ":" + PORT
	}

	DEBUG_MODE = getEnv("DEBUG_MODE", "true") == "true"
	DEFAULT_STREAM = getEnv("DEFAULT_STREAM", "true") == "true"
	DASHBOARD_ENABLED = getEnv("DASHBOARD_ENABLED", "true") == "true"
	ENABLE_THINKING = getEnv("ENABLE_THINKING", "false") == "true"
}

// 记录请求统计信息
func recordRequestStats(startTime time.Time, path string, status int) {
	duration := time.Since(startTime)

	statsMutex.Lock()
	defer statsMutex.Unlock()

	stats.TotalRequests++
	stats.LastRequestTime = time.Now()

	if status >= 200 && status < 300 {
		stats.SuccessfulRequests++
	} else {
		stats.FailedRequests++
	}

	// 更新平均响应时间
	if stats.TotalRequests > 0 {
		totalDuration := stats.AverageResponseTime*time.Duration(stats.TotalRequests-1) + duration
		stats.AverageResponseTime = totalDuration / time.Duration(stats.TotalRequests)
	} else {
		stats.AverageResponseTime = duration
	}
}

// 添加实时请求信息
func addLiveRequest(method, path string, status int, duration time.Duration, _, userAgent string) {
	requestsMutex.Lock()
	defer requestsMutex.Unlock()

	request := LiveRequest{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Method:    method,
		Path:      path,
		Status:    status,
		Duration:  duration.Milliseconds(),
		UserAgent: userAgent,
	}

	liveRequests = append(liveRequests, request)

	// 只保留最近的请求记录
	if len(liveRequests) > MAX_LIVE_REQUESTS {
		liveRequests = liveRequests[1:]
	}
}

// 获取实时请求数据（用于SSE）
func getLiveRequestsData() []byte {
	requestsMutex.Lock()
	defer requestsMutex.Unlock()

	// 确保 liveRequests 不为 nil
	if liveRequests == nil {
		liveRequests = []LiveRequest{}
	}

	data, err := json.Marshal(liveRequests)
	if err != nil {
		// 如果序列化失败，返回空数组
		emptyArray := []LiveRequest{}
		data, _ = json.Marshal(emptyArray)
	}
	return data
}

// 获取统计数据（用于SSE）
func getStatsData() []byte {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	data, _ := json.Marshal(stats)
	return data
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 加载 .env 文件
func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		// 文件不存在时不报错，这样 .env.local 是可选的
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 只有当环境变量未设置时才从文件加载
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

// 获取客户端IP地址
func getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 检查X-Real-IP头
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 使用RemoteAddr
	ip := r.RemoteAddr
	// 移除端口号
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

// OpenAI 请求结构
type OpenAIRequest struct {
	Model          string    `json:"model"`
	Messages       []Message `json:"messages"`
	Stream         bool      `json:"stream,omitempty"`
	Temperature    float64   `json:"temperature,omitempty"`
	MaxTokens      int       `json:"max_tokens,omitempty"`
	EnableThinking *bool     `json:"enable_thinking,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 上游请求结构
type UpstreamRequest struct {
	Stream          bool                   `json:"stream"`
	Model           string                 `json:"model"`
	Messages        []Message              `json:"messages"`
	Params          map[string]interface{} `json:"params"`
	Features        map[string]interface{} `json:"features"`
	BackgroundTasks map[string]bool        `json:"background_tasks,omitempty"`
	ChatID          string                 `json:"chat_id,omitempty"`
	ID              string                 `json:"id,omitempty"`
	MCPServers      []string               `json:"mcp_servers,omitempty"`
	ModelItem       struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		OwnedBy string `json:"owned_by"`
	} `json:"model_item,omitempty"`
	ToolServers []string          `json:"tool_servers,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// OpenAI 响应结构
type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message,omitempty"`
	Delta        Delta   `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// 上游SSE响应结构
type UpstreamData struct {
	Type string `json:"type"`
	Data struct {
		DeltaContent string         `json:"delta_content"`
		Phase        string         `json:"phase"`
		Done         bool           `json:"done"`
		Usage        Usage          `json:"usage,omitempty"`
		Error        *UpstreamError `json:"error,omitempty"`
		Inner        *struct {
			Error *UpstreamError `json:"error,omitempty"`
		} `json:"data,omitempty"`
	} `json:"data"`
	Error *UpstreamError `json:"error,omitempty"`
}

type UpstreamError struct {
	Detail string `json:"detail"`
	Code   int    `json:"code"`
}

// 模型列表响应
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// debug日志函数
func debugLog(format string, args ...interface{}) {
	if DEBUG_MODE {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// 转换思考内容的通用函数
func transformThinkingContent(s string) string {
	// 去除 <summary>…</summary>
	s = regexp.MustCompile(`(?s)<summary>.*?</summary>`).ReplaceAllString(s, "")
	// 清理残留自定义标签，如 </thinking>、<Full> 等
	s = strings.ReplaceAll(s, "</thinking>", "")
	s = strings.ReplaceAll(s, "<Full>", "")
	s = strings.ReplaceAll(s, "</Full>", "")
	s = strings.TrimSpace(s)

	switch THINK_TAGS_MODE {
	case "think":
		s = regexp.MustCompile(`<details[^>]*>`).ReplaceAllString(s, "<think>")
		s = strings.ReplaceAll(s, "</details>", "</think>")
	case "strip":
		s = regexp.MustCompile(`<details[^>]*>`).ReplaceAllString(s, "")
		s = strings.ReplaceAll(s, "</details>", "")
	}

	// 处理每行前缀 "> "（包括起始位置）
	s = strings.TrimPrefix(s, "> ")
	s = strings.ReplaceAll(s, "\n> ", "\n")
	return strings.TrimSpace(s)
}

// 根据模型名称获取上游实际模型ID
func getUpstreamModelID(modelName string) string {
	switch modelName {
	case "GLM-4.6":
		return "GLM-4-6-API-V1" // 使用官方API的真实模型名称
	default:
		debugLog("未知模型名称: %s，使用GLM-4.6作为默认", modelName)
		return "GLM-4-6-API-V1" // 默认使用GLM-4.6
	}
}

// 获取匿名token（每次对话使用不同token，避免共享记忆）
func getAnonymousToken() (string, error) {
	client := &http.Client{Timeout: AUTH_TOKEN_TIMEOUT * time.Second}
	req, err := http.NewRequest("GET", ORIGIN_BASE+"/api/v1/auths/", nil)
	if err != nil {
		return "", err
	}
	// 伪装浏览器头
	req.Header.Set("User-Agent", BROWSER_UA)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("X-FE-Version", X_FE_VERSION)
	req.Header.Set("sec-ch-ua", SEC_CH_UA)
	req.Header.Set("sec-ch-ua-mobile", SEC_CH_UA_MOB)
	req.Header.Set("sec-ch-ua-platform", SEC_CH_UA_PLAT)
	req.Header.Set("Origin", ORIGIN_BASE)
	req.Header.Set("Referer", ORIGIN_BASE+"/")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anon token status=%d", resp.StatusCode)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Token == "" {
		return "", fmt.Errorf("anon token empty")
	}
	return body.Token, nil
}

// urlsafeB64Decode 解码URL安全的base64字符串（自动添加padding）
func urlsafeB64Decode(data string) ([]byte, error) {
	// 添加必要的padding
	padding := len(data) % 4
	if padding > 0 {
		data += strings.Repeat("=", 4-padding)
	}
	return base64.URLEncoding.DecodeString(data)
}

// decodeJWTPayload 解码JWT的payload部分（不验证签名）
func decodeJWTPayload(token string) map[string]interface{} {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return map[string]interface{}{}
	}

	payloadBytes, err := urlsafeB64Decode(parts[1])
	if err != nil {
		return map[string]interface{}{}
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return map[string]interface{}{}
	}

	return payload
}

// extractUserIDFromToken 从JWT token中提取user_id
func extractUserIDFromToken(token string) string {
	if token == "" {
		return "guest"
	}

	payload := decodeJWTPayload(token)

	// 尝试多个可能的字段名
	for _, key := range []string{"id", "user_id", "uid", "sub"} {
		if val, ok := payload[key]; ok {
			if strVal, ok := val.(string); ok && strVal != "" {
				return strVal
			}
		}
	}

	return "guest"
}

// generateSignature 生成双层HMAC-SHA256签名
// Layer1: derived_key = HMAC(secret, window_index)
// Layer2: signature = HMAC(derived_key, canonical_string)
// canonical_string = "requestId,<id>,timestamp,<ts>,user_id,<uid>|<msg>|<ts>"
func generateSignature(messageText, requestID string, timestampMs int64, userID, secret string) string {
	if secret == "" {
		secret = "junjie"
	}

	// 构建规范字符串
	r := fmt.Sprintf("%d", timestampMs)
	e := fmt.Sprintf("requestId,%s,timestamp,%d,user_id,%s", requestID, timestampMs, userID)
	canonicalString := fmt.Sprintf("%s|%s|%s", e, messageText, r)

	// Layer1: 基于5分钟时间窗口生成派生密钥
	windowIndex := timestampMs / (5 * 60 * 1000)
	rootKey := []byte(secret)

	mac1 := hmac.New(sha256.New, rootKey)
	mac1.Write([]byte(fmt.Sprintf("%d", windowIndex)))
	derivedHex := fmt.Sprintf("%x", mac1.Sum(nil))

	// Layer2: 使用派生密钥对规范字符串签名
	mac2 := hmac.New(sha256.New, []byte(derivedHex))
	mac2.Write([]byte(canonicalString))
	signature := fmt.Sprintf("%x", mac2.Sum(nil))

	return signature
}

// extractLastUserMessage 提取最后一条用户消息的文本内容
func extractLastUserMessage(messages []Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return messages[i].Content
		}
	}
	return ""
}

func main() {
	// 初始化配置
	initConfig()

	// 注册路由
	http.HandleFunc("/v1/models", handleModels)
	http.HandleFunc("/v1/chat/completions", handleChatCompletions)
	http.HandleFunc("/docs", handleAPIDocs)
	http.HandleFunc("/", handleOptions)

	// Dashboard路由
	if DASHBOARD_ENABLED {
		http.HandleFunc("/dashboard", handleDashboard)
		http.HandleFunc("/dashboard/stats", handleDashboardStats)
		http.HandleFunc("/dashboard/requests", handleDashboardRequests)
		log.Printf("Dashboard已启用，访问地址: http://localhost%s/dashboard", PORT)
	}

	// 初始化注册管理系统
	registerEnabled := getEnv("REGISTER_ENABLED", "true")
	if registerEnabled == "true" || registerEnabled == "1" {
		dbPath := getEnv("REGISTER_DB_PATH", "./data/zai2api.db")
		if err := register.InitRegisterSystem(dbPath); err != nil {
			log.Printf("❌ 注册系统初始化失败: %v", err)
		} else {
			// 注册路由
			register.RegisterRoutes(http.DefaultServeMux)
			log.Printf("🔐 注册管理: http://localhost%s/register/login", PORT)
		}
	}

	log.Printf("OpenAI兼容API服务器启动在端口%s", PORT)
	log.Printf("模型: %s", MODEL_NAME)
	log.Printf("上游: %s", UPSTREAM_URL)
	log.Printf("API密钥: %s", func() string {
		if len(DEFAULT_KEY) > TOKEN_DISPLAY_LENGTH {
			return DEFAULT_KEY[:TOKEN_DISPLAY_LENGTH] + "..."
		}
		return DEFAULT_KEY
	}())
	log.Printf("Debug模式: %v", DEBUG_MODE)
	log.Printf("默认流式响应: %v", DEFAULT_STREAM)
	log.Printf("Dashboard启用: %v", DASHBOARD_ENABLED)
	log.Printf("思考功能: %v", ENABLE_THINKING)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

// Dashboard页面处理器
func handleDashboard(w http.ResponseWriter, r *http.Request) {
	// 只允许GET请求
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 动态HTML模板，使用当前配置的模型名称
	tmpl := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API调用看板</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 20px;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .stats-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background-color: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            text-align: center;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .stat-value {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
        }
        .stat-label {
            font-size: 14px;
            color: #6c757d;
            margin-top: 5px;
        }
        .requests-container {
            margin-top: 30px;
        }
        .requests-table {
            width: 100%%;
            border-collapse: collapse;
        }
        .requests-table th, .requests-table td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .requests-table th {
            background-color: #f8f9fa;
        }
        .status-success {
            color: #28a745;
        }
        .status-error {
            color: #dc3545;
        }
        .refresh-info {
            text-align: center;
            margin-top: 20px;
            color: #6c757d;
            font-size: 14px;
        }
        .pagination-container {
            display: flex;
            justify-content: center;
            align-items: center;
            margin-top: 20px;
            gap: 10px;
        }
        .pagination-container button {
            padding: 5px 10px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        .pagination-container button:disabled {
            background-color: #cccccc;
            cursor: not-allowed;
        }
        .pagination-container button:hover:not(:disabled) {
            background-color: #0056b3;
        }
        .chart-container {
            margin-top: 30px;
            height: 300px;
            background-color: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>API调用看板</h1>

        <div class="stats-container">
            <div class="stat-card">
                <div class="stat-value" id="total-requests">0</div>
                <div class="stat-label">总请求数</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="successful-requests">0</div>
                <div class="stat-label">成功请求</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="failed-requests">0</div>
                <div class="stat-label">失败请求</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="avg-response-time">0s</div>
                <div class="stat-label">平均响应时间</div>
            </div>
        </div>

        <div class="chart-container">
            <h2>请求统计图表</h2>
            <canvas id="requestsChart"></canvas>
        </div>

        <div class="requests-container">
            <h2>实时请求</h2>
            <table class="requests-table">
                <thead>
                    <tr>
                        <th>时间</th>
                        <th>模型</th>
                        <th>方法</th>
                        <th>状态</th>
                        <th>耗时</th>
                        <th>User Agent</th>
                    </tr>
                </thead>
                <tbody id="requests-tbody">
                    <!-- 请求记录将通过JavaScript动态添加 -->
                </tbody>
            </table>
            <div class="pagination-container">
                <button id="prev-page" disabled>上一页</button>
                <span id="page-info">第 1 页，共 1 页</span>
                <button id="next-page" disabled>下一页</button>
            </div>
        </div>

        <div class="refresh-info">
            数据每5秒自动刷新一次
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script>
        // 全局变量
        let allRequests = [];
        let currentPage = 1;
        const itemsPerPage = 10;
        let requestsChart = null;

        // 更新统计数据
        function updateStats() {
            fetch('/dashboard/stats')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('total-requests').textContent = data.TotalRequests;
                    document.getElementById('successful-requests').textContent = data.SuccessfulRequests;
                    document.getElementById('failed-requests').textContent = data.FailedRequests;
                    document.getElementById('avg-response-time').textContent = (data.AverageResponseTime / 1000000000).toFixed(2) + 's';
                })
                .catch(error => console.error('Error fetching stats:', error));
        }

        // 更新请求列表
        function updateRequests() {
            fetch('/dashboard/requests')
                .then(response => response.json())
                .then(data => {
                    // 检查数据是否为数组
                    if (!Array.isArray(data)) {
                        console.error('返回的数据不是数组:', data);
                        return;
                    }

                    // 保存所有请求数据
                    allRequests = data;

                    // 按时间倒序排列
                    allRequests.sort((a, b) => {
                        const timeA = new Date(a.timestamp);
                        const timeB = new Date(b.timestamp);
                        return timeB - timeA;
                    });

                    // 更新表格
                    updateTable();

                    // 更新图表
                    updateChart();

                    // 更新分页信息
                    updatePagination();
                })
                .catch(error => console.error('Error fetching requests:', error));
        }

        // 更新表格显示
        function updateTable() {
            const tbody = document.getElementById('requests-tbody');
            tbody.innerHTML = '';

            // 计算当前页的数据范围
            const startIndex = (currentPage - 1) * itemsPerPage;
            const endIndex = startIndex + itemsPerPage;
            const currentRequests = allRequests.slice(startIndex, endIndex);

            currentRequests.forEach(request => {
                const row = document.createElement('tr');

                // 格式化时间 - 检查时间戳是否有效
                let timeStr = "Invalid Date";
                if (request.timestamp) {
                    try {
                        const time = new Date(request.timestamp);
                        if (!isNaN(time.getTime())) {
                            timeStr = time.toLocaleTimeString();
                        }
                    } catch (e) {
                        console.error("时间格式化错误:", e);
                    }
                }

                // 状态样式
                const statusClass = request.status >= 200 && request.status < 300 ? 'status-success' : 'status-error';

                // 截断 User Agent，避免过长
                let userAgent = request.user_agent || "undefined";
                if (userAgent.length > 30) {
                    userAgent = userAgent.substring(0, 30) + "...";
                }

                row.innerHTML =
                   "<td>" + timeStr + "</td>" +
                   "<td>%s</td>" +
                   "<td>" + (request.method || "undefined") + "</td>" +
                   "<td class=\"" + statusClass + "\">" + (request.status || "undefined") + "</td>" +
                   "<td>" + ((request.duration / 1000).toFixed(2) || "undefined") + "s</td>" +
                   "<td title=\"" + (request.user_agent || "") + "\">" + userAgent + "</td>";

                tbody.appendChild(row);
            });
        }

        // 更新分页信息
        function updatePagination() {
            const totalPages = Math.ceil(allRequests.length / itemsPerPage);
            document.getElementById('page-info').textContent = "第 " + currentPage + " 页，共 " + totalPages + " 页";

            document.getElementById('prev-page').disabled = currentPage <= 1;
            document.getElementById('next-page').disabled = currentPage >= totalPages;
        }

        // 更新图表
        function updateChart() {
            const ctx = document.getElementById('requestsChart').getContext('2d');

            // 准备图表数据 - 最近20条请求的响应时间
            const chartData = allRequests.slice(0, 20).reverse();
            const labels = chartData.map(req => {
                const time = new Date(req.timestamp);
                return time.toLocaleTimeString();
            });
            const responseTimes = chartData.map(req => req.duration);

            // 如果图表已存在，先销毁
            if (requestsChart) {
                requestsChart.destroy();
            }

            // 创建新图表
            requestsChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: '响应时间 (s)',
                        data: responseTimes.map(time => time / 1000),
                        borderColor: '#007bff',
                        backgroundColor: 'rgba(0, 123, 255, 0.1)',
                        tension: 0.1,
                        fill: true
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: '响应时间 (s)'
                            }
                        },
                        x: {
                            title: {
                                display: true,
                                text: '时间'
                            }
                        }
                    },
                    plugins: {
                        title: {
                            display: true,
                            text: '最近20条请求的响应时间趋势 (s)'
                        }
                    }
                }
            });
        }

        // 分页按钮事件
        document.getElementById('prev-page').addEventListener('click', function() {
            if (currentPage > 1) {
                currentPage--;
                updateTable();
                updatePagination();
            }
        });

        document.getElementById('next-page').addEventListener('click', function() {
            const totalPages = Math.ceil(allRequests.length / itemsPerPage);
            if (currentPage < totalPages) {
                currentPage++;
                updateTable();
                updatePagination();
            }
        });

        // 初始加载
        updateStats();
        updateRequests();

        // 定时刷新
        setInterval(updateStats, 5000);
        setInterval(updateRequests, 5000);
    </script>
</body>
</html>`, MODEL_NAME)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, tmpl)
}

// Dashboard统计数据处理器
func handleDashboardStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(getStatsData())
}

// Dashboard请求数据处理器
func handleDashboardRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(getLiveRequestsData())
}

// API文档页面处理器
func handleAPIDocs(w http.ResponseWriter, r *http.Request) {
	// 只允许GET请求
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 动态API文档HTML模板，使用当前配置的模型名称
	tmpl := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZtoApi 文档</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
            line-height: 1.6;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 30px;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 2px solid #007bff;
            padding-bottom: 10px;
        }
        h2 {
            color: #007bff;
            margin-top: 30px;
            margin-bottom: 15px;
        }
        h3 {
            color: #333;
            margin-top: 25px;
            margin-bottom: 10px;
        }
        .endpoint {
            background-color: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            margin-bottom: 20px;
            border-left: 4px solid #007bff;
        }
        .method {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            color: white;
            font-weight: bold;
            margin-right: 10px;
            font-size: 14px;
        }
        .get { background-color: #28a745; }
        .post { background-color: #007bff; }
        .path {
            font-family: monospace;
            background-color: #e9ecef;
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 16px;
        }
        .description {
            margin: 15px 0;
        }
        .parameters {
            margin: 15px 0;
        }
        table {
            width: 100%%;
            border-collapse: collapse;
            margin: 15px 0;
        }
        th, td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        .example {
            background-color: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            margin: 15px 0;
            font-family: monospace;
            white-space: pre-wrap;
            overflow-x: auto;
        }
        .note {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 10px 15px;
            margin: 15px 0;
            border-radius: 0 4px 4px 0;
        }
        .response {
            background-color: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            margin: 15px 0;
            font-family: monospace;
            white-space: pre-wrap;
            overflow-x: auto;
        }
        .tab {
            overflow: hidden;
            border: 1px solid #ccc;
            background-color: #f1f1f1;
            border-radius: 4px 4px 0 0;
        }
        .tab button {
            background-color: inherit;
            float: left;
            border: none;
            outline: none;
            cursor: pointer;
            padding: 14px 16px;
            transition: 0.3s;
            font-size: 16px;
        }
        .tab button:hover {
            background-color: #ddd;
        }
        .tab button.active {
            background-color: #ccc;
        }
        .tabcontent {
            display: none;
            padding: 6px 12px;
            border: 1px solid #ccc;
            border-top: none;
            border-radius: 0 0 4px 4px;
        }
        .toc {
            background-color: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            margin-bottom: 20px;
        }
        .toc ul {
            padding-left: 20px;
        }
        .toc li {
            margin: 5px 0;
        }
        .toc a {
            color: #007bff;
            text-decoration: none;
        }
        .toc a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>ZtoApi 文档</h1>

        <div class="toc">
            <h2>目录</h2>
            <ul>
                <li><a href="#overview">概述</a></li>
                <li><a href="#authentication">身份验证</a></li>
                <li><a href="#endpoints">API端点</a>
                    <ul>
                        <li><a href="#models">获取模型列表</a></li>
                        <li><a href="#chat-completions">聊天完成</a></li>
                    </ul>
                </li>
                <li><a href="#examples">使用示例</a></li>
                <li><a href="#error-handling">错误处理</a></li>
            </ul>
        </div>

        <section id="overview">
            <h2>概述</h2>
            <p>这是一个为Z.ai %s模型提供OpenAI兼容API接口的代理服务器。它允许你使用标准的OpenAI API格式与Z.ai的%s模型进行交互，支持流式和非流式响应。</p>
            <p><strong>基础URL:</strong> <code>http://localhost:9090/v1</code></p>
            <div class="note">
                <strong>注意:</strong> 默认端口为9090，可以通过环境变量PORT进行修改。
            </div>
        </section>

        <section id="authentication">
            <h2>身份验证</h2>
            <p>所有API请求都需要在请求头中包含有效的API密钥进行身份验证：</p>
            <div class="example">
Authorization: Bearer your-api-key</div>
            <p>默认的API密钥为 <code>sk-your-key</code>，可以通过环境变量 <code>DEFAULT_KEY</code> 进行修改。</p>
        </section>

        <section id="endpoints">
            <h2>API端点</h2>

            <div class="endpoint" id="models">
                <h3>获取模型列表</h3>
                <div>
                    <span class="method get">GET</span>
                    <span class="path">/v1/models</span>
                </div>
                <div class="description">
                    <p>获取可用模型列表。</p>
                </div>
                <div class="parameters">
                    <h4>请求参数</h4>
                    <p>无</p>
                </div>
                <div class="response">
{
  "object": "list",
  "data": [
    {
      "id": "%s",
      "object": "model",
      "created": 1756788845,
      "owned_by": "z.ai"
    }
  ]
}</div>
            </div>

            <div class="endpoint" id="chat-completions">
                <h3>聊天完成</h3>
                <div>
                    <span class="method post">POST</span>
                    <span class="path">/v1/chat/completions</span>
                </div>
                <div class="description">
                    <p>基于消息列表生成模型响应。支持流式和非流式两种模式。</p>
                </div>
                <div class="parameters">
                    <h4>请求参数</h4>
                    <table>
                        <thead>
                            <tr>
                                <th>参数名</th>
                                <th>类型</th>
                                <th>必需</th>
                                <th>说明</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>model</td>
                                <td>string</td>
                                <td>是</td>
                                <td>要使用的模型ID，例如 "%s"</td>
                            </tr>
                            <tr>
                                <td>messages</td>
                                <td>array</td>
                                <td>是</td>
                                <td>消息列表，包含角色和内容</td>
                            </tr>
                            <tr>
                                <td>stream</td>
                                <td>boolean</td>
                                <td>否</td>
                                <td>是否使用流式响应，默认为true</td>
                            </tr>
                            <tr>
                                <td>temperature</td>
                                <td>number</td>
                                <td>否</td>
                                <td>采样温度，控制随机性</td>
                            </tr>
                            <tr>
                               <td>max_tokens</td>
                               <td>integer</td>
                               <td>否</td>
                               <td>生成的最大令牌数</td>
                           </tr>
                           <tr>
                               <td>enable_thinking</td>
                               <td>boolean</td>
                               <td>否</td>
                               <td>是否启用思考功能，默认使用环境变量 ENABLE_THINKING 的值</td>
                           </tr>
                        </tbody>
                    </table>
                </div>
                <div class="parameters">
                    <h4>消息格式</h4>
                    <table>
                        <thead>
                            <tr>
                                <th>字段</th>
                                <th>类型</th>
                                <th>说明</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>role</td>
                                <td>string</td>
                                <td>消息角色，可选值：system、user、assistant</td>
                            </tr>
                            <tr>
                                <td>content</td>
                                <td>string</td>
                                <td>消息内容</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </section>

        <section id="examples">
            <h2>使用示例</h2>

            <div class="tab">
                <button class="tablinks active" onclick="openTab(event, 'python-tab')">Python</button>
                <button class="tablinks" onclick="openTab(event, 'curl-tab')">cURL</button>
                <button class="tablinks" onclick="openTab(event, 'javascript-tab')">JavaScript</button>
            </div>

            <div id="python-tab" class="tabcontent" style="display: block;">
                <h3>Python示例</h3>
                <div class="example">
import openai

# 配置客户端
client = openai.OpenAI(
    api_key="your-api-key",  # 对应 DEFAULT_KEY
    base_url="http://localhost:9090/v1"
)

# 非流式请求
response = client.chat.completions.create(
    model="%s",
    messages=[{"role": "user", "content": "你好，请介绍一下自己"}]
)

print(response.choices[0].message.content)

# 流式请求
response = client.chat.completions.create(
    model="%s",
    messages=[{"role": "user", "content": "请写一首关于春天的诗"}],
    stream=True
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")</div>
            </div>

            <div id="curl-tab" class="tabcontent">
                <h3>cURL示例</h3>
                <div class="example">
# 非流式请求
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "model": "GLM-4.6",
    "messages": [{"role": "user", "content": "你好"}],
    "stream": false
  }'

# 流式请求
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "model": "GLM-4.6",
    "messages": [{"role": "user", "content": "你好"}],
    "stream": true
  }'</div>

# 启用思考功能的请求
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "model": "GLM-4.6",
    "messages": [{"role": "user", "content": "请分析一下这个问题"}],
    "enable_thinking": true
  }'
            </div>

            <div id="javascript-tab" class="tabcontent">
                <h3>JavaScript示例</h3>
                <div class="example">
const fetch = require('node-fetch');

async function chatWithGLM(message, stream = false) {
  const response = await fetch('http://localhost:9090/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer your-api-key'
    },
    body: JSON.stringify({
      model: '%s',
      messages: [{ role: 'user', content: message }],
      stream: stream
    })
  });

  if (stream) {
    // 处理流式响应
    const reader = response.body.getReader();
    const decoder = new TextDecoder();

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunk = decoder.decode(value);
      const lines = chunk.split('\n');

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6);
          if (data === '[DONE]') {
            console.log('\n流式响应完成');
            return;
          }

          try {
            const parsed = JSON.parse(data);
            const content = parsed.choices[0]?.delta?.content;
            if (content) {
              process.stdout.write(content);
            }
          } catch (e) {
            // 忽略解析错误
          }
        }
      }
    }
  } else {
    // 处理非流式响应
    const data = await response.json();
    console.log(data.choices[0].message.content);
  }
}

// 使用示例
chatWithGLM('你好，请介绍一下JavaScript', false);</div>
            </div>
        </section>

        <section id="error-handling">
            <h2>错误处理</h2>
            <p>API使用标准HTTP状态码来表示请求的成功或失败：</p>
            <table>
                <thead>
                    <tr>
                        <th>状态码</th>
                        <th>说明</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>200 OK</td>
                        <td>请求成功</td>
                    </tr>
                    <tr>
                        <td>400 Bad Request</td>
                        <td>请求格式错误或参数无效</td>
                    </tr>
                    <tr>
                        <td>401 Unauthorized</td>
                        <td>API密钥无效或缺失</td>
                    </tr>
                    <tr>
                        <td>502 Bad Gateway</td>
                        <td>上游服务错误</td>
                    </tr>
                </tbody>
            </table>
            <div class="note">
                <strong>注意:</strong> 在调试模式下，服务器会输出详细的日志信息，可以通过设置环境变量 DEBUG_MODE=true 来启用。
            </div>
        </section>
    </div>

    <script>
        function openTab(evt, tabName) {
            var i, tabcontent, tablinks;
            tabcontent = document.getElementsByClassName("tabcontent");
            for (i = 0; i < tabcontent.length; i++) {
                tabcontent[i].style.display = "none";
            }
            tablinks = document.getElementsByClassName("tablinks");
            for (i = 0; i < tablinks.length; i++) {
                tablinks[i].className = tablinks[i].className.replace(" active", "");
            }
            document.getElementById(tabName).style.display = "block";
            evt.currentTarget.className += " active";
        }
    </script>
</body>
</html>`, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, tmpl)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func handleModels(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := ModelsResponse{
		Object: "list",
		Data: []Model{
			{
				ID:      "GLM-4.6",
				Object:  "model",
				Created: time.Now().Unix(),
				OwnedBy: "z.ai",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	path := r.URL.Path
	clientIP := getClientIP(r)
	userAgent := r.UserAgent()

	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	debugLog("收到chat completions请求")

	// 验证API Key
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		debugLog("缺少或无效的Authorization头")
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusUnauthorized)
		addLiveRequest(r.Method, path, http.StatusUnauthorized, duration, "", userAgent)
		return
	}

	apiKey := strings.TrimPrefix(authHeader, "Bearer ")
	if apiKey != DEFAULT_KEY {
		debugLog("无效的API key: %s", apiKey)
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusUnauthorized)
		addLiveRequest(r.Method, path, http.StatusUnauthorized, duration, "", userAgent)
		return
	}

	debugLog("API key验证通过")

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		debugLog("读取请求体失败: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusBadRequest)
		addLiveRequest(r.Method, path, http.StatusBadRequest, duration, "", userAgent)
		return
	}

	// 解析请求
	var req OpenAIRequest
	if err := json.Unmarshal(body, &req); err != nil {
		debugLog("JSON解析失败: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusBadRequest)
		addLiveRequest(r.Method, path, http.StatusBadRequest, duration, "", userAgent)
		return
	}

	// 如果客户端没有明确指定stream参数，使用默认值
	if !bytes.Contains(body, []byte(`"stream"`)) {
		req.Stream = DEFAULT_STREAM
		debugLog("客户端未指定stream参数，使用默认值: %v", DEFAULT_STREAM)
	}

	debugLog("请求解析成功 - 模型: %s, 流式: %v, 消息数: %d", req.Model, req.Stream, len(req.Messages))

	// 生成会话相关ID
	chatID := fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
	msgID := fmt.Sprintf("%d", time.Now().UnixNano())

	// 决定是否启用思考功能：优先使用请求参数，其次使用环境变量
	enableThinking := ENABLE_THINKING // 默认使用环境变量值
	if req.EnableThinking != nil {
		enableThinking = *req.EnableThinking
		debugLog("使用请求参数中的思考功能设置: %v", enableThinking)
	} else {
		debugLog("使用环境变量中的思考功能设置: %v", enableThinking)
	}

	// 构造上游请求
	upstreamReq := UpstreamRequest{
		Stream:   true, // 总是使用流式从上游获取
		ChatID:   chatID,
		ID:       msgID,
		Model:    getUpstreamModelID(MODEL_NAME), // 根据模型名称获取上游实际模型ID
		Messages: req.Messages,
		Params:   map[string]interface{}{},
		Features: map[string]interface{}{
			"enable_thinking": enableThinking,
		},
		BackgroundTasks: map[string]bool{
			"title_generation": false,
			"tags_generation":  false,
		},
		MCPServers: []string{},
		ModelItem: struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			OwnedBy string `json:"owned_by"`
		}{ID: getUpstreamModelID(MODEL_NAME), Name: MODEL_NAME, OwnedBy: "openai"},
		ToolServers: []string{},
		Variables: map[string]string{
			"{{USER_NAME}}":        "User",
			"{{USER_LOCATION}}":    "Unknown",
			"{{CURRENT_DATETIME}}": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	// 选择本次对话使用的token：优先使用配置的ZAI_TOKEN，否则获取匿名token
	authToken := ZAI_TOKEN
	if authToken == "" && ANON_TOKEN_ENABLED {
		if t, err := getAnonymousToken(); err == nil {
			authToken = t
			debugLog("使用匿名token: %s...", func() string {
				if len(t) > TOKEN_DISPLAY_LENGTH {
					return t[:TOKEN_DISPLAY_LENGTH]
				}
				return t
			}())
		} else {
			debugLog("匿名token获取失败: %v", err)
		}
	} else if authToken != "" {
		debugLog("使用配置的ZAI_TOKEN: %s...", func() string {
			if len(authToken) > TOKEN_DISPLAY_LENGTH {
				return authToken[:TOKEN_DISPLAY_LENGTH]
			}
			return authToken
		}())
	}

	// 调用上游API
	if req.Stream {
		handleStreamResponseWithIDs(w, upstreamReq, chatID, authToken, startTime, path, clientIP, userAgent)
	} else {
		handleNonStreamResponseWithIDs(w, upstreamReq, chatID, authToken, startTime, path, clientIP, userAgent)
	}
}

func callUpstreamWithHeaders(upstreamReq UpstreamRequest, refererChatID string, authToken string) (*http.Response, error) {
	reqBody, err := json.Marshal(upstreamReq)
	if err != nil {
		debugLog("上游请求序列化失败: %v", err)
		return nil, err
	}

	// 构建带URL参数的完整URL
	baseURL := UPSTREAM_URL
	timestampMs := time.Now().UnixMilli()
	timestamp := fmt.Sprintf("%d", timestampMs)

	// 生成UUID (简化版，使用crypto/rand会更好)
	requestID := fmt.Sprintf("%x-%x-%x-%x-%x",
		time.Now().UnixNano(), time.Now().Unix(),
		time.Now().Nanosecond(), time.Now().Second(), time.Now().Minute())

	// 从token中提取user_id（而不是随机生成）
	userID := extractUserIDFromToken(authToken)

	// 提取最后一条用户消息用于签名
	lastUserMessage := extractLastUserMessage(upstreamReq.Messages)

	// 获取签名密钥（从环境变量或使用默认值）
	secret := getEnv("ZAI_SIGNING_SECRET", "junjie")

	// 生成双层HMAC-SHA256签名
	signature := generateSignature(lastUserMessage, requestID, timestampMs, userID, secret)

	debugLog("签名参数 - user_id: %s, message: %s..., timestamp: %d",
		userID,
		func() string {
			if len(lastUserMessage) > 20 {
				return lastUserMessage[:20]
			}
			return lastUserMessage
		}(),
		timestampMs)
	debugLog("生成签名: %s (双层HMAC-SHA256)", signature)

	// 构建URL参数 - 添加所有必要的指纹参数
	fullURL := fmt.Sprintf("%s?timestamp=%s&requestId=%s&user_id=%s&version=0.0.1&platform=web&token=%s"+
		"&user_agent=%s&language=zh-CN&languages=zh-CN,zh&timezone=Asia/Shanghai"+
		"&cookie_enabled=true&screen_width=1680&screen_height=1050&screen_resolution=1680x1050"+
		"&viewport_height=812&viewport_width=1087&viewport_size=1087x812"+
		"&color_depth=30&pixel_ratio=2"+
		"&current_url=%s&pathname=/c/%s&search=&hash="+
		"&host=chat.z.ai&hostname=chat.z.ai&protocol=https:&referrer="+
		"&title=%s"+
		"&timezone_offset=-480&local_time=%s&utc_time=%s"+
		"&is_mobile=false&is_touch=false&max_touch_points=0"+
		"&browser_name=Chrome&os_name=Mac+OS&signature_timestamp=%s",
		baseURL, timestamp, requestID, userID, authToken,
		url.QueryEscape(BROWSER_UA),
		url.QueryEscape(ORIGIN_BASE+"/c/"+refererChatID), refererChatID,
		url.QueryEscape("Z.ai Chat - Free AI powered by GLM-4.6"),
		url.QueryEscape(time.Now().Format("2006-01-02T15:04:05.000Z")),
		url.QueryEscape(time.Now().UTC().Format(time.RFC1123)),
		timestamp,
	)

	debugLog("调用上游API: %s", fullURL)
	debugLog("上游请求体: %s", string(reqBody))

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		debugLog("创建HTTP请求失败: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("User-Agent", BROWSER_UA)
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("sec-ch-ua", SEC_CH_UA)
	req.Header.Set("sec-ch-ua-mobile", SEC_CH_UA_MOB)
	req.Header.Set("sec-ch-ua-platform", SEC_CH_UA_PLAT)
	req.Header.Set("X-FE-Version", X_FE_VERSION)
	req.Header.Set("X-Signature", signature)
	req.Header.Set("Origin", ORIGIN_BASE)
	req.Header.Set("Referer", ORIGIN_BASE+"/c/"+refererChatID)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	// 添加Cookie
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", authToken))

	// 创建HTTP客户端 - 流式请求专用配置
	// 不设置总超时(Timeout)，只设置连接和响应头超时
	// 这样只要数据持续到达，连接就会保持，支持长时间思考
	client := &http.Client{
		Transport: &http.Transport{
			// 连接超时：30秒
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			// 响应头超时：60秒（等待服务器开始响应）
			ResponseHeaderTimeout: 60 * time.Second,
			// TLS握手超时
			TLSHandshakeTimeout: 10 * time.Second,
			// 最大空闲连接
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			// 不设置整体超时，让流式响应可以持续任意长时间
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		debugLog("上游请求失败: %v", err)
		return nil, err
	}

	debugLog("上游响应状态: %d %s", resp.StatusCode, resp.Status)
	return resp, nil
}

func handleStreamResponseWithIDs(w http.ResponseWriter, upstreamReq UpstreamRequest, chatID string, authToken string, startTime time.Time, path string, clientIP, userAgent string) {
	debugLog("开始处理流式响应 (chat_id=%s)", chatID)

	resp, err := callUpstreamWithHeaders(upstreamReq, chatID, authToken)
	if err != nil {
		debugLog("调用上游失败: %v", err)
		http.Error(w, "Failed to call upstream", http.StatusBadGateway)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusBadGateway)
		addLiveRequest("POST", path, http.StatusBadGateway, duration, "", userAgent)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		debugLog("上游返回错误状态: %d", resp.StatusCode)
		// 读取错误响应体
		if DEBUG_MODE {
			body, _ := io.ReadAll(resp.Body)
			debugLog("上游错误响应: %s", string(body))
		}
		http.Error(w, "Upstream error", http.StatusBadGateway)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusBadGateway)
		addLiveRequest("POST", path, http.StatusBadGateway, duration, "", userAgent)
		return
	}

	// 策略2：总是展示thinking + answer

	// 设置SSE头部
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 发送第一个chunk（role）
	firstChunk := OpenAIResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   MODEL_NAME,
		Choices: []Choice{
			{
				Index: 0,
				Delta: Delta{Role: "assistant"},
			},
		},
	}
	writeSSEChunk(w, firstChunk)
	flusher.Flush()

	// 读取上游SSE流
	debugLog("开始读取上游SSE流")
	scanner := bufio.NewScanner(resp.Body)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		dataStr := strings.TrimPrefix(line, "data: ")
		if dataStr == "" {
			continue
		}

		debugLog("收到SSE数据 (第%d行): %s", lineCount, dataStr)

		var upstreamData UpstreamData
		if err := json.Unmarshal([]byte(dataStr), &upstreamData); err != nil {
			debugLog("SSE数据解析失败: %v", err)
			continue
		}

		// 错误检测（data.error 或 data.data.error 或 顶层error）
		if (upstreamData.Error != nil) || (upstreamData.Data.Error != nil) || (upstreamData.Data.Inner != nil && upstreamData.Data.Inner.Error != nil) {
			errObj := upstreamData.Error
			if errObj == nil {
				errObj = upstreamData.Data.Error
			}
			if errObj == nil && upstreamData.Data.Inner != nil {
				errObj = upstreamData.Data.Inner.Error
			}
			debugLog("上游错误: code=%d, detail=%s", errObj.Code, errObj.Detail)
			// 结束下游流
			endChunk := OpenAIResponse{
				ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   MODEL_NAME,
				Choices: []Choice{{Index: 0, Delta: Delta{}, FinishReason: "stop"}},
			}
			writeSSEChunk(w, endChunk)
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			break
		}

		debugLog("解析成功 - 类型: %s, 阶段: %s, 内容长度: %d, 完成: %v",
			upstreamData.Type, upstreamData.Data.Phase, len(upstreamData.Data.DeltaContent), upstreamData.Data.Done)

		// 策略2：总是展示thinking + answer
		if upstreamData.Data.DeltaContent != "" {
			var out = upstreamData.Data.DeltaContent
			if upstreamData.Data.Phase == "thinking" {
				out = transformThinkingContent(out)
			}
			if out != "" {
				debugLog("发送内容(%s): %s", upstreamData.Data.Phase, out)
				chunk := OpenAIResponse{
					ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
					Object:  "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model:   MODEL_NAME,
					Choices: []Choice{
						{
							Index: 0,
							Delta: Delta{Content: out},
						},
					},
				}
				writeSSEChunk(w, chunk)
				flusher.Flush()
			}
		}

		// 检查是否结束
		if upstreamData.Data.Done || upstreamData.Data.Phase == "done" {
			debugLog("检测到流结束信号")
			// 发送结束chunk
			endChunk := OpenAIResponse{
				ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   MODEL_NAME,
				Choices: []Choice{
					{
						Index:        0,
						Delta:        Delta{},
						FinishReason: "stop",
					},
				},
			}
			writeSSEChunk(w, endChunk)
			flusher.Flush()

			// 发送[DONE]
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			debugLog("流式响应完成，共处理%d行", lineCount)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		debugLog("扫描器错误: %v", err)
	}

	// 记录成功请求统计
	duration := time.Since(startTime)
	recordRequestStats(startTime, path, http.StatusOK)
	addLiveRequest("POST", path, http.StatusOK, duration, "", userAgent)
}

func writeSSEChunk(w http.ResponseWriter, chunk OpenAIResponse) {
	data, _ := json.Marshal(chunk)
	fmt.Fprintf(w, "data: %s\n\n", data)
}

func handleNonStreamResponseWithIDs(w http.ResponseWriter, upstreamReq UpstreamRequest, chatID string, authToken string, startTime time.Time, path string, clientIP, userAgent string) {
	debugLog("开始处理非流式响应 (chat_id=%s)", chatID)

	resp, err := callUpstreamWithHeaders(upstreamReq, chatID, authToken)
	if err != nil {
		debugLog("调用上游失败: %v", err)
		http.Error(w, "Failed to call upstream", http.StatusBadGateway)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusBadGateway)
		addLiveRequest("POST", path, http.StatusBadGateway, duration, "", userAgent)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		debugLog("上游返回错误状态: %d", resp.StatusCode)
		// 读取错误响应体
		if DEBUG_MODE {
			body, _ := io.ReadAll(resp.Body)
			debugLog("上游错误响应: %s", string(body))
		}
		http.Error(w, "Upstream error", http.StatusBadGateway)
		// 记录请求统计
		duration := time.Since(startTime)
		recordRequestStats(startTime, path, http.StatusBadGateway)
		addLiveRequest("POST", path, http.StatusBadGateway, duration, "", userAgent)
		return
	}

	// 收集完整响应（策略2：thinking与answer都纳入，thinking转换）
	var fullContent strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	debugLog("开始收集完整响应内容")
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		debugLog("收到原始行[%d]: %s", lineCount, line)

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		dataStr := strings.TrimPrefix(line, "data: ")
		if dataStr == "" {
			continue
		}

		debugLog("解析SSE数据: %s", dataStr)

		var upstreamData UpstreamData
		if err := json.Unmarshal([]byte(dataStr), &upstreamData); err != nil {
			debugLog("JSON解析失败: %v", err)
			continue
		}

		debugLog("解析成功 - type:%s phase:%s content_len:%d done:%v",
			upstreamData.Type, upstreamData.Data.Phase,
			len(upstreamData.Data.DeltaContent), upstreamData.Data.Done)

		if upstreamData.Data.DeltaContent != "" {
			out := upstreamData.Data.DeltaContent
			if upstreamData.Data.Phase == "thinking" {
				out = transformThinkingContent(out)
			}
			if out != "" {
				debugLog("添加内容: %s", out)
				fullContent.WriteString(out)
			}
		}

		if upstreamData.Data.Done || upstreamData.Data.Phase == "done" {
			debugLog("检测到完成信号，停止收集")
			break
		}
	}

	debugLog("扫描器共处理%d行", lineCount)

	finalContent := fullContent.String()
	debugLog("内容收集完成，最终长度: %d", len(finalContent))

	// 构造完整响应
	response := OpenAIResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   MODEL_NAME,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: finalContent,
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	debugLog("非流式响应发送完成")

	// 记录成功请求统计
	duration := time.Since(startTime)
	recordRequestStats(startTime, path, http.StatusOK)
	addLiveRequest("POST", path, http.StatusOK, duration, "", userAgent)
}
