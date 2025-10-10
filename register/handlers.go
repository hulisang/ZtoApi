package register

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 登录页面处理器
func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, LoginPage)
}

// 主页面处理器
func HandleMainPage(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Redirect(w, r, "/register/login", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, MainPage)
}

// 登录API处理器
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)

	// 检查IP是否被锁定
	if locked, remaining := IsIPLocked(clientIP); locked {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":       false,
			"locked":        true,
			"remainingTime": remaining.Seconds(),
			"error":         "账号已锁定",
		})
		return
	}

	// 验证凭证
	if !ValidateLogin(req.Username, req.Password) {
		RecordLoginFailure(clientIP)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "用户名或密码错误",
		})
		return
	}

	// 创建会话
	session, err := CreateSession()
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// 清除登录失败记录
	ClearLoginFailure(clientIP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"sessionId": session.ID,
	})
}

// 登出API处理器
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionId")
	if err == nil {
		DeleteSession(cookie.Value)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 获取账号列表API处理器
func HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	filter := r.URL.Query().Get("filter")
	search := r.URL.Query().Get("search")

	accounts, total, err := GetAccounts(page, pageSize, filter, search)
	if err != nil {
		http.Error(w, "Failed to get accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"accounts": accounts,
		"pagination": map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}

// 获取统计信息API处理器
func HandleGetStats(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	stats, err := GetStats()
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// 删除账号API处理器
func HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := DeleteAccount(req.Email); err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 批量删除账号API处理器
func HandleBatchDeleteAccounts(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Emails []string `json:"emails"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := BatchDeleteAccounts(req.Emails); err != nil {
		http.Error(w, "Failed to delete accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 导出账号API处理器
func HandleExportAccounts(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accounts, _, err := GetAccounts(1, 100000, "", "")
	if err != nil {
		http.Error(w, "Failed to get accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=zai_accounts.txt")

	for _, acc := range accounts {
		fmt.Fprintf(w, "%s----%s----%s----%s\n", acc.Email, acc.Password, acc.Token, acc.APIKEY)
	}
}

// 开始注册API处理器
func HandleStartRegister(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Count  int            `json:"count"`
		Config RegisterConfig `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	task := GetCurrentTask()
	if task != nil && task.IsRunning {
		http.Error(w, "Task already running", http.StatusConflict)
		return
	}

	// 启动后台注册任务（使用SSE广播，不需要通道）
	go BatchRegisterAccounts(req.Count, req.Config, nil, nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 停止注册API处理器
func HandleStopRegister(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	StopCurrentTask()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// SSE流处理器（实时日志）
func HandleRegisterStream(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 创建客户端通道
	client := make(chan string, 100)
	AddSSEClient(client)
	defer RemoveSSEClient(client)

	// 发送连接状态
	task := GetCurrentTask()
	isRunning := task != nil && task.IsRunning
	connectData := map[string]interface{}{
		"type":      "connected",
		"isRunning": isRunning,
	}
	jsonData, _ := json.Marshal(connectData)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()

	// 发送历史日志（最近50条）
	history := GetLogHistory()
	startIdx := 0
	if len(history) > 50 {
		startIdx = len(history) - 50
	}
	for i := startIdx; i < len(history); i++ {
		logData, _ := json.Marshal(history[i])
		fmt.Fprintf(w, "data: %s\n\n", logData)
	}
	flusher.Flush()

	// 保持连接并发送消息
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-client:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ticker.C:
			// 发送保活消息
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

// 导入账号API处理器
func HandleImportAccounts(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取上传的文件
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// 解析账号数据（格式: email----password----token----apikey）
	lines := strings.Split(string(content), "\n")
	imported := 0
	failed := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "----")
		if len(parts) < 3 {
			failed++
			continue
		}

		account := &Account{
			Email:     strings.TrimSpace(parts[0]),
			Password:  strings.TrimSpace(parts[1]),
			Token:     strings.TrimSpace(parts[2]),
			Status:    "unknown",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if len(parts) > 3 && strings.TrimSpace(parts[3]) != "" {
			account.APIKEY = strings.TrimSpace(parts[3])
		}

		if err := SaveAccount(account); err != nil {
			failed++
		} else {
			imported++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"imported": imported,
		"failed":   failed,
	})
}

// 获取配置API处理器
func HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	config := GetConfig()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// 保存配置API处理器
func HandleSaveConfig(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config RegisterConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := SaveConfig(config); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 批量补充APIKEY处理器
func HandleBatchRefetchAPIKEY(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Emails []string `json:"emails"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	logChan := make(chan string, 100)
	go func() {
		config := GetConfig()
		BatchRefetchAPIKEY(req.Emails, config, logChan)
		close(logChan)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 批量检测存活处理器
func HandleBatchCheckAccounts(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Emails []string `json:"emails"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	logChan := make(chan string, 100)
	go func() {
		BatchCheckAccounts(req.Emails, logChan)
		close(logChan)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// 删除失效账号处理器
func HandleDeleteInactiveAccounts(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count, err := DeleteInactiveAccounts()
	if err != nil {
		http.Error(w, "Failed to delete inactive accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   count,
	})
}

// 单独获取APIKEY
func HandleRefetchSingleAPIKEY(w http.ResponseWriter, r *http.Request) {
	if !CheckAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求参数",
		})
		return
	}

	if req.Email == "" || req.Token == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "缺少必需参数: email 或 token",
		})
		return
	}

	// 尝试使用Token获取APIKEY
	config := GetConfig()
	apikey, err := getAPIKEY(req.Token, config)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 更新数据库中的账号APIKEY
	if err := UpdateAccountAPIKEY(req.Email, apikey); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "更新APIKEY失败: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"apikey":  apikey,
	})
}

