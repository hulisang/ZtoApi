package register

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// PushPlus通知
func sendNotification(title, content string, config RegisterConfig) {
	if !config.EnableNotification || config.PushPlusToken == "" {
		return
	}

	payload := map[string]interface{}{
		"token":    config.PushPlusToken,
		"title":    title,
		"content":  content,
		"template": "markdown",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://www.pushplus.plus/send", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	client.Do(req)
	// 忽略错误
}

// 生成随机邮箱
func generateEmail() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	username := hex.EncodeToString(bytes)
	domain := EmailDomains[time.Now().UnixNano()%int64(len(EmailDomains))]
	return fmt.Sprintf("%s@%s", username, domain)
}

// 生成随机密码
func generatePassword() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	b := make([]byte, 14)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// 注册Z.AI账号（完整流程与Deno版本一致）
func registerZAIAccount(email, password string, config RegisterConfig) (*Account, error) {
	client := &http.Client{Timeout: time.Duration(config.HTTPTimeout) * time.Second}
	name := strings.Split(email, "@")[0]

	// 1. 调用signup API
	BroadcastLog("info", "  → 注册...")
	signupPayload := map[string]interface{}{
		"name":              name,
		"email":             email,
		"password":          password,
		"profile_image_url": "data:image/png;base64,",
		"sso_redirect":      nil,
	}
	
	signupBody, _ := json.Marshal(signupPayload)
	req, err := http.NewRequest("POST", "https://chat.z.ai/api/v1/auths/signup", bytes.NewBuffer(signupBody))
	if err != nil {
		return nil, fmt.Errorf("  ✗ 创建请求失败:%v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Origin", "https://chat.z.ai")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("  ✗ 注册失败:%v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("  ✗ 注册失败:HTTP%d:%s", resp.StatusCode, string(body))
	}

	var signupResp struct {
		Success bool `json:"success"`
	}
	json.NewDecoder(resp.Body).Decode(&signupResp)
	if !signupResp.Success {
		return nil, fmt.Errorf("  ✗ 被拒绝")
	}

	BroadcastLog("success", "  ✓ 注册成功")

	// 2. 等待验证邮件
	emailCheckURL := fmt.Sprintf("https://mail.chatgpt.org.uk/api/get-emails?email=%s", email)
	BroadcastLogWithLink("info", fmt.Sprintf("  → 等待邮件:%s", email), "打开邮箱", emailCheckURL)
	emailContent, err := waitForVerificationEmail(email, config)
	if err != nil {
		return nil, fmt.Errorf("  ✗ %v", err)
	}

	// 3. 提取验证链接并解析参数
	BroadcastLog("info", "  → 提取链接...")
	verifyURL := extractVerificationURL(emailContent)
	if verifyURL == "" {
		preview := emailContent
		if len(preview) > 500 {
			preview = preview[:500]
		}
		return nil, fmt.Errorf("  ✗ 未找到链接:%s...", strings.ReplaceAll(preview, "\n", " "))
	}

	token, emailFromURL, username := parseVerificationURL(verifyURL)
	if token == "" || emailFromURL == "" || username == "" {
		return nil, fmt.Errorf("  ✗ 链接格式错")
	}
	
	BroadcastLog("success", "  ✓ 链接已提取")

	// 4. 完成注册（finish_signup）
	BroadcastLog("info", "  → 验证...")
	finishPayload := map[string]interface{}{
		"email":             emailFromURL,
		"password":          password,
		"profile_image_url": "data:image/png;base64,",
		"sso_redirect":      nil,
		"token":             token,
		"username":          username,
	}
	
	finishBody, _ := json.Marshal(finishPayload)
	req, _ = http.NewRequest("POST", "https://chat.z.ai/api/v1/auths/finish_signup", bytes.NewBuffer(finishBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Origin", "https://chat.z.ai")
	
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("  ✗ 验证失败:%v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("  ✗ 验证失败:HTTP%d", resp.StatusCode)
	}

	var finishResp struct {
		Success bool `json:"success"`
		User    struct {
			Token string `json:"token"`
		} `json:"user"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&finishResp); err != nil {
		return nil, fmt.Errorf("  ✗ 解析响应失败:%v", err)
	}
	
	if !finishResp.Success || finishResp.User.Token == "" {
		return nil, fmt.Errorf("  ✗ 验证拒绝或无Token")
	}
	
	userToken := finishResp.User.Token
	BroadcastLog("success", "  ✓ 获得Token")

	account := &Account{
		Email:     email,
		Password:  password,
		Token:     userToken,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 5. 快速模式：跳过APIKEY获取
	if config.SkipAPIKey {
		return account, nil
	}

	// 6. 正常模式：获取APIKEY
	BroadcastLog("info", "  → 登录API...")
	apikey, err := getAPIKEY(userToken, config)
	if err != nil {
		BroadcastLog("warning", fmt.Sprintf("  ⚠️ API登录失败:%v(仅Token)", err))
		return account, nil // 返回账号，但没有APIKEY
	}
	
	account.APIKEY = apikey
	return account, nil
}

// 等待验证邮件
func waitForVerificationEmail(email string, config RegisterConfig) (string, error) {
	// 使用相同的邮箱API
	apiURL := fmt.Sprintf("https://mail.chatgpt.org.uk/api/get-emails?email=%s", email)
	client := &http.Client{Timeout: 10 * time.Second}
	
	startTime := time.Now()
	attempts := 0
	maxAttempts := config.EmailTimeout / config.EmailCheckInterval
	lastReportTime := 0
	
	for i := 0; i < maxAttempts; i++ {
		attempts++
		elapsed := int(time.Since(startTime).Seconds())
		
		resp, err := client.Get(apiURL)
		if err != nil {
			time.Sleep(time.Duration(config.EmailCheckInterval) * time.Second)
			continue
		}
		
		var data struct {
			Emails []struct {
				From    string `json:"from"`
				Subject string `json:"subject"`
				Content string `json:"content"`
			} `json:"emails"`
		}
		
		json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()
		
		// 每10秒报告进度
		if elapsed-lastReportTime >= 10 && elapsed > 0 {
			progress := int(float64(elapsed) / float64(config.EmailTimeout) * 100)
			if progress > 99 {
				progress = 99
			}
			remaining := config.EmailTimeout - elapsed
			BroadcastLog("info", fmt.Sprintf("  等待邮件[%d%%] 已用:%ds/剩余:%ds(尝试%d次)", progress, elapsed, remaining, attempts))
			lastReportTime = elapsed
		}
		
		// 查找Z.AI的验证邮件
		if data.Emails != nil {
			for _, emailData := range data.Emails {
				if strings.Contains(strings.ToLower(emailData.From), "z.ai") {
					BroadcastLog("success", fmt.Sprintf("  ✓ 收到邮件(%ds)", elapsed))
					return emailData.Content, nil
				}
			}
		}
		
		time.Sleep(time.Duration(config.EmailCheckInterval) * time.Second)
	}
	
	return "", fmt.Errorf("邮件超时(%ds)", config.EmailTimeout)
}

// 提取验证URL（多种匹配方式，与Deno版本一致）
func extractVerificationURL(emailContent string) string {
	// 方式1: /auth/verify_email
	re := regexp.MustCompile(`https://chat\.z\.ai/auth/verify_email\?[^\s<>"']+`)
	if match := re.FindString(emailContent); match != "" {
		return strings.ReplaceAll(strings.ReplaceAll(match, "&amp;", "&"), "&#39;", "'")
	}
	
	// 方式2: /verify_email
	re = regexp.MustCompile(`https://chat\.z\.ai/verify_email\?[^\s<>"']+`)
	if match := re.FindString(emailContent); match != "" {
		return strings.ReplaceAll(strings.ReplaceAll(match, "&amp;", "&"), "&#39;", "'")
	}
	
	// 方式3: HTML编码
	re = regexp.MustCompile(`https?://chat\.z\.ai/(?:auth/)?verify_email[^"'\s]*`)
	if match := re.FindString(emailContent); match != "" {
		return strings.ReplaceAll(strings.ReplaceAll(match, "&amp;", "&"), "&#39;", "'")
	}
	
	// 方式4: JSON格式
	re = regexp.MustCompile(`"(https?://[^"]*verify_email[^"]*)"`)
	if matches := re.FindStringSubmatch(emailContent); len(matches) > 1 {
		match := matches[1]
		match = strings.ReplaceAll(match, "\\u0026", "&")
		match = strings.ReplaceAll(match, "&amp;", "&")
		match = strings.ReplaceAll(match, "&#39;", "'")
		return match
	}
	
	return ""
}

// 解析验证URL参数
func parseVerificationURL(verifyURL string) (token, email, username string) {
	u, err := url.Parse(verifyURL)
	if err != nil {
		return "", "", ""
	}
	
	query := u.Query()
	return query.Get("token"), query.Get("email"), query.Get("username")
}

// 获取APIKEY（完整流程）
func getAPIKEY(token string, config RegisterConfig) (string, error) {
	client := &http.Client{Timeout: time.Duration(config.HTTPTimeout) * time.Second}
	
	// 1. 登录API获取accessToken
	accessToken, err := loginToAPI(token, client)
	if err != nil {
		return "", err
	}
	BroadcastLog("success", "  ✓ API登录成功")
	
	// 2. 获取客户信息（组织和项目ID）
	BroadcastLog("info", "  → 组织...")
	orgID, projectID, err := getCustomerInfo(accessToken, client)
	if err != nil {
		BroadcastLog("error", fmt.Sprintf("  ✗ 组织失败:%v", err))
		return "", err
	}
	BroadcastLog("success", "  ✓ 获取组织成功")
	
	// 3. 创建APIKEY
	BroadcastLog("info", "  → APIKEY...")
	apikey, err := createAPIKey(accessToken, orgID, projectID, client)
	if err != nil {
		BroadcastLog("error", fmt.Sprintf("  ✗ APIKEY创建失败:%v", err))
		return "", err
	}
	BroadcastLog("success", "  ✓ APIKEY创建成功")
	
	return apikey, nil
}

// 登录到Z.AI API
func loginToAPI(token string, client *http.Client) (string, error) {
	payload := map[string]interface{}{"token": token}
	body, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", "https://api.z.ai/api/auth/z/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Origin", "https://z.ai")
	req.Header.Set("Referer", "https://z.ai/")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	if !result.Success || result.Code != 200 || result.Data.AccessToken == "" {
		return "", fmt.Errorf("API登录失败")
	}
	
	return result.Data.AccessToken, nil
}

// 获取客户信息
func getCustomerInfo(accessToken string, client *http.Client) (string, string, error) {
	req, _ := http.NewRequest("GET", "https://api.z.ai/api/biz/customer/getCustomerInfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Origin", "https://z.ai")
	req.Header.Set("Referer", "https://z.ai/")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    struct {
			Organizations []struct {
				OrganizationID string `json:"organizationId"`
				Projects       []struct {
					ProjectID string `json:"projectId"`
				} `json:"projects"`
			} `json:"organizations"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	
	if !result.Success || result.Code != 200 || len(result.Data.Organizations) == 0 {
		return "", "", fmt.Errorf("获取组织失败")
	}
	
	org := result.Data.Organizations[0]
	if len(org.Projects) == 0 {
		return "", "", fmt.Errorf("无可用项目")
	}
	
	return org.OrganizationID, org.Projects[0].ProjectID, nil
}

// 创建APIKEY
func createAPIKey(accessToken, orgID, projectID string, client *http.Client) (string, error) {
	url := fmt.Sprintf("https://api.z.ai/api/biz/v1/organization/%s/projects/%s/api_keys", orgID, projectID)
	
	randomName := fmt.Sprintf("key_%d", time.Now().UnixNano())
	payload := map[string]interface{}{"name": randomName}
	body, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Origin", "https://z.ai")
	req.Header.Set("Referer", "https://z.ai/")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    struct {
			APIKey    string `json:"apiKey"`
			SecretKey string `json:"secretKey"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	if !result.Success || result.Code != 200 {
		return "", fmt.Errorf("创建APIKEY失败")
	}
	
	// 拼接APIKEY（格式: apiKey.secretKey）
	finalKey := fmt.Sprintf("%s.%s", result.Data.APIKey, result.Data.SecretKey)
	if finalKey == "." || result.Data.APIKey == "" || result.Data.SecretKey == "" {
		return "", fmt.Errorf("APIKEY无效")
	}
	
	return finalKey, nil
}

// 保存账号到数据库
func SaveAccount(account *Account) error {
	_, err := db.Exec(`
		INSERT INTO accounts (email, password, token, apikey, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, account.Email, account.Password, account.Token, account.APIKEY, account.Status, account.CreatedAt, account.UpdatedAt)
	return err
}

// 批量注册账号
func BatchRegisterAccounts(count int, config RegisterConfig, logChan chan<- string, progressChan chan<- map[string]interface{}) {
	if currentTask != nil && currentTask.IsRunning {
		BroadcastLog("error", "❌ 已有任务正在运行")
		if logChan != nil {
			logChan <- "❌ 已有任务正在运行"
		}
		return
	}

	currentTask = &RegisterTask{
		ID:         fmt.Sprintf("task-%d", time.Now().Unix()),
		Total:      count,
		Success:    0,
		Failed:     0,
		IsRunning:  true,
		ShouldStop: false,
		StartTime:  time.Now(),
		Config:     config,
	}

	BroadcastLog("info", fmt.Sprintf("🚀 开始批量注册 %d 个账号", count))
	skipMode := "否"
	if config.SkipAPIKey {
		skipMode = "是(稍后批量获取)"
	}
	BroadcastLog("info", fmt.Sprintf("⚙️ 配置: 并发=%d 间隔=%dms 快速=%s 超时=%ds", 
		config.Concurrency, config.RegisterDelay, skipMode, config.EmailTimeout))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Concurrency)
	
	for i := 0; i < count; i++ {
		if currentTask.ShouldStop {
			BroadcastLog("warning", "⏹️ 用户停止注册")
			break
		}

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			email := generateEmail()
			password := generatePassword()
			emailCheckURL := fmt.Sprintf("https://mail.chatgpt.org.uk/api/get-emails?email=%s", email)
			
			BroadcastLogWithLink("info", fmt.Sprintf("▶ 开始:%s", email), "邮箱", emailCheckURL)
			
			for retry := 0; retry < config.RetryTimes; retry++ {
				account, err := registerZAIAccount(email, password, config)
				if err != nil {
					if retry < config.RetryTimes-1 {
						BroadcastLog("warning", fmt.Sprintf("⚠️ 重试%d/%d:%s", retry+1, config.RetryTimes, err.Error()))
						time.Sleep(2 * time.Second)
						continue
					}
					BroadcastLog("error", fmt.Sprintf("❌ 失败:%s", err.Error()))
					currentTask.Failed++
					break
				}
				
				if err := SaveAccount(account); err != nil {
					BroadcastLog("warning", fmt.Sprintf("⚠️ 保存失败:%v", err))
				}
				
				// 根据模式和APIKEY情况输出不同消息
				if config.SkipAPIKey {
					BroadcastLog("success", fmt.Sprintf("✅ 快速完成:%s(稍后获取KEY)", email))
				} else if account.APIKEY != "" {
					BroadcastLog("success", fmt.Sprintf("✅ 完成:%s(含KEY)", email))
				} else {
					BroadcastLog("warning", fmt.Sprintf("⚠️ 成功但KEY失败:%s(仅Token)", email))
				}
				currentTask.Success++
				break
			}
			
			// 广播进度
			BroadcastProgress(currentTask.Total, currentTask.Success, currentTask.Failed)
			
			// 间隔延迟
			time.Sleep(time.Duration(config.RegisterDelay) * time.Millisecond)
		}(i)
	}

	wg.Wait()
	currentTask.IsRunning = false
	
	elapsed := time.Since(currentTask.StartTime)
	BroadcastLog("success", fmt.Sprintf("🎉 注册完成! 成功: %d, 失败: %d, 耗时: %.1fs", 
		currentTask.Success, currentTask.Failed, elapsed.Seconds()))
	
	// 发送完成事件
	completeData := map[string]interface{}{
		"type":    "complete",
		"total":   currentTask.Total,
		"success": currentTask.Success,
		"failed":  currentTask.Failed,
		"elapsed": elapsed.Seconds(),
	}
	data, _ := json.Marshal(completeData)
	
	sseClientMutex.RLock()
	for client := range sseClients {
		select {
		case client <- string(data):
		default:
		}
	}
	sseClientMutex.RUnlock()
	
	// 发送通知
	notifyContent := fmt.Sprintf("## 注册任务完成\n\n- 总数: %d\n- 成功: %d\n- 失败: %d\n- 耗时: %.1fs", 
		currentTask.Total, currentTask.Success, currentTask.Failed, elapsed.Seconds())
	sendNotification("Z.AI 注册完成", notifyContent, config)
}

// 获取账号列表（支持搜索和筛选）
func GetAccounts(page, pageSize int, filter, search string) ([]Account, int64, error) {
	// 构建查询条件
	where := "1=1"
	args := []interface{}{}
	
	// 搜索功能
	if search != "" {
		where += " AND (email LIKE ? OR password LIKE ? OR token LIKE ? OR apikey LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	// 快速筛选
	if filter == "has-apikey" {
		where += " AND apikey IS NOT NULL AND apikey != ''"
	} else if filter == "no-apikey" {
		where += " AND (apikey IS NULL OR apikey = '')"
	} else if filter == "inactive" {
		where += " AND status = 'inactive'"
	} else if filter == "today" {
		where += " AND DATE(created_at) = DATE('now')"
	} else if filter == "week" {
		where += " AND created_at >= DATE('now', '-7 days')"
	}

	// 获取总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM accounts WHERE " + where
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, email, password, token, COALESCE(apikey, ''), status, created_at, updated_at
		FROM accounts
		WHERE %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, where)
	
	queryArgs := append(args, pageSize, offset)
	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		var acc Account
		var apikey sql.NullString
		err := rows.Scan(&acc.ID, &acc.Email, &acc.Password, &acc.Token, &apikey, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt)
		if err != nil {
			continue
		}
		if apikey.Valid {
			acc.APIKEY = apikey.String
		}
		accounts = append(accounts, acc)
	}

	return accounts, total, nil
}

// 获取统计信息
func GetStats() (*Stats, error) {
	stats := &Stats{}
	
	// 总账号数
	db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&stats.TotalAccounts)
	
	// 有APIKEY的账号数
	db.QueryRow("SELECT COUNT(*) FROM accounts WHERE apikey != ''").Scan(&stats.WithAPIKEY)
	
	// 无APIKEY的账号数
	db.QueryRow("SELECT COUNT(*) FROM accounts WHERE apikey IS NULL OR apikey = ''").Scan(&stats.WithoutAPIKEY)
	
	// 活跃账号数
	db.QueryRow("SELECT COUNT(*) FROM accounts WHERE status = 'active'").Scan(&stats.ActiveAccounts)
	
	// 失效账号数
	db.QueryRow("SELECT COUNT(*) FROM accounts WHERE status = 'inactive'").Scan(&stats.InactiveAccounts)
	
	return stats, nil
}

// 删除账号
func DeleteAccount(email string) error {
	_, err := db.Exec("DELETE FROM accounts WHERE email = ?", email)
	return err
}

// 批量删除账号
func BatchDeleteAccounts(emails []string) error {
	if len(emails) == 0 {
		return nil
	}
	
	placeholders := strings.Repeat("?,", len(emails)-1) + "?"
	query := fmt.Sprintf("DELETE FROM accounts WHERE email IN (%s)", placeholders)
	
	args := make([]interface{}, len(emails))
	for i, email := range emails {
		args[i] = email
	}
	
	_, err := db.Exec(query, args...)
	return err
}

// 停止当前任务
func StopCurrentTask() {
	if currentTask != nil {
		currentTask.ShouldStop = true
	}
}

// 获取当前任务状态
func GetCurrentTask() *RegisterTask {
	return currentTask
}

// 批量补充APIKEY（为无APIKEY的账号获取）
func BatchRefetchAPIKEY(emails []string, config RegisterConfig, logChan chan<- string) (int, int) {
	success := 0
	failed := 0

	logChan <- fmt.Sprintf("🔑 开始批量获取APIKEY，共 %d 个账号", len(emails))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Concurrency)
	resultMutex := sync.Mutex{}

	for i, email := range emails {
		wg.Add(1)
		go func(idx int, email string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 获取账号Token
			var token string
			err := db.QueryRow("SELECT token FROM accounts WHERE email = ?", email).Scan(&token)
			if err != nil {
				logChan <- fmt.Sprintf("[%d/%d] ❌ %s: 获取Token失败", idx+1, len(emails), email)
				resultMutex.Lock()
				failed++
				resultMutex.Unlock()
				return
			}

			// 获取APIKEY
			apikey, err := getAPIKEY(token, config)
			if err != nil {
				logChan <- fmt.Sprintf("[%d/%d] ❌ %s: 获取APIKEY失败 - %v", idx+1, len(emails), email, err)
				resultMutex.Lock()
				failed++
				resultMutex.Unlock()
				return
			}

			// 更新数据库
			_, err = db.Exec("UPDATE accounts SET apikey = ?, updated_at = CURRENT_TIMESTAMP WHERE email = ?", apikey, email)
			if err != nil {
				logChan <- fmt.Sprintf("[%d/%d] ⚠️ %s: 更新失败", idx+1, len(emails), email)
				resultMutex.Lock()
				failed++
				resultMutex.Unlock()
				return
			}

			logChan <- fmt.Sprintf("[%d/%d] ✅ %s: 获取成功", idx+1, len(emails), email)
			resultMutex.Lock()
			success++
			resultMutex.Unlock()

			time.Sleep(time.Duration(config.RegisterDelay) * time.Millisecond)
		}(i, email)
	}

	wg.Wait()
	logChan <- fmt.Sprintf("🎉 批量获取完成！成功: %d, 失败: %d", success, failed)
	return success, failed
}

// 批量检测账号存活性
func BatchCheckAccounts(emails []string, logChan chan<- string) (int, int) {
	active := 0
	inactive := 0

	BroadcastLog("info", fmt.Sprintf("🔍 开始批量检测，共 %d 个账号", len(emails)))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 检测并发数固定为10
	resultMutex := sync.Mutex{}

	for i, email := range emails {
		wg.Add(1)
		go func(idx int, email string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 获取账号Token
			var token string
			err := db.QueryRow("SELECT token FROM accounts WHERE email = ?", email).Scan(&token)
			if err != nil {
				BroadcastLog("error", fmt.Sprintf("[%d/%d] %s ⚠️ API登录失败", idx+1, len(emails), email))
				resultMutex.Lock()
				inactive++
				resultMutex.Unlock()
				return
			}

			// 检测账号有效性（使用API登录测试）
			isActive := checkAccountStatus(token)

			// 更新状态
			status := "active"
			if !isActive {
				status = "inactive"
			}

			db.Exec("UPDATE accounts SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE email = ?", status, email)

			if isActive {
				BroadcastLog("success", fmt.Sprintf("[%d/%d] %s ✓ API登录成功", idx+1, len(emails), email))
				resultMutex.Lock()
				active++
				resultMutex.Unlock()
			} else {
				BroadcastLog("warning", fmt.Sprintf("[%d/%d] %s ⚠️ API登录失败", idx+1, len(emails), email))
				resultMutex.Lock()
				inactive++
				resultMutex.Unlock()
			}

			time.Sleep(500 * time.Millisecond)
		}(i, email)
	}

	wg.Wait()
	BroadcastLog("success", fmt.Sprintf("🎉 批量检测完成！正常: %d, 失效: %d", active, inactive))
	
	// 发送检测完成事件
	checkCompleteData := map[string]interface{}{
		"type":     "check_complete",
		"active":   active,
		"inactive": inactive,
		"total":    len(emails),
	}
	data, _ := json.Marshal(checkCompleteData)
	
	sseClientMutex.RLock()
	for client := range sseClients {
		select {
		case client <- string(data):
		default:
		}
	}
	sseClientMutex.RUnlock()
	
	return active, inactive
}

// 检测账号状态（使用API登录测试）
func checkAccountStatus(token string) bool {
	config := GetConfig()
	client := &http.Client{Timeout: time.Duration(config.HTTPTimeout) * time.Second}
	
	// 尝试API登录
	_, err := loginToAPI(token, client)
	return err == nil
}

// 删除失效账号
func DeleteInactiveAccounts() (int, error) {
	result, err := db.Exec("DELETE FROM accounts WHERE status = 'inactive'")
	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()
	return int(count), nil
}

// UpdateAccountAPIKEY 更新账号的APIKEY
func UpdateAccountAPIKEY(email, apikey string) error {
	_, err := db.Exec(`
		UPDATE accounts 
		SET apikey = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE email = ?
	`, apikey, email)
	return err
}

