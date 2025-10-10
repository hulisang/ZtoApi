package register

import (
	"encoding/json"
	"sync"
	"time"
)

var (
	// SSE全局通道
	sseClients      = make(map[chan string]bool)
	sseClientMutex  sync.RWMutex
	logHistory      = make([]map[string]interface{}, 0, 100)
	logHistoryMutex sync.RWMutex
)

// BroadcastLog 广播日志到所有SSE客户端
func BroadcastLog(level, message string) {
	logEntry := map[string]interface{}{
		"type":    "log",
		"level":   level,
		"message": message,
		"time":    time.Now().Format("15:04:05"),
	}

	// 添加到历史记录
	logHistoryMutex.Lock()
	logHistory = append(logHistory, logEntry)
	// 保留最近100条
	if len(logHistory) > 100 {
		logHistory = logHistory[len(logHistory)-100:]
	}
	logHistoryMutex.Unlock()

	// 广播到所有客户端
	data, _ := json.Marshal(logEntry)
	messageStr := string(data)

	sseClientMutex.RLock()
	defer sseClientMutex.RUnlock()

	for client := range sseClients {
		select {
		case client <- messageStr:
		default:
			// 客户端通道满，跳过
		}
	}
}

// BroadcastLogWithLink 广播带链接的日志到所有SSE客户端
func BroadcastLogWithLink(level, message, linkText, linkURL string) {
	logEntry := map[string]interface{}{
		"type":    "log",
		"level":   level,
		"message": message,
		"time":    time.Now().Format("15:04:05"),
		"link": map[string]string{
			"text": linkText,
			"url":  linkURL,
		},
	}

	// 添加到历史记录
	logHistoryMutex.Lock()
	logHistory = append(logHistory, logEntry)
	// 保留最近100条
	if len(logHistory) > 100 {
		logHistory = logHistory[len(logHistory)-100:]
	}
	logHistoryMutex.Unlock()

	// 广播到所有客户端
	data, _ := json.Marshal(logEntry)
	messageStr := string(data)

	sseClientMutex.RLock()
	defer sseClientMutex.RUnlock()

	for client := range sseClients {
		select {
		case client <- messageStr:
		default:
		}
	}
}

// BroadcastProgress 广播进度到所有SSE客户端
func BroadcastProgress(total, success, failed int) {
	progressEntry := map[string]interface{}{
		"type":    "progress",
		"total":   total,
		"success": success,
		"failed":  failed,
	}

	data, _ := json.Marshal(progressEntry)
	message := string(data)

	sseClientMutex.RLock()
	defer sseClientMutex.RUnlock()

	for client := range sseClients {
		select {
		case client <- message:
		default:
		}
	}
}

// AddSSEClient 添加SSE客户端
func AddSSEClient(client chan string) {
	sseClientMutex.Lock()
	sseClients[client] = true
	sseClientMutex.Unlock()
}

// RemoveSSEClient 移除SSE客户端
func RemoveSSEClient(client chan string) {
	sseClientMutex.Lock()
	delete(sseClients, client)
	sseClientMutex.Unlock()
	close(client)
}

// GetLogHistory 获取历史日志
func GetLogHistory() []map[string]interface{} {
	logHistoryMutex.RLock()
	defer logHistoryMutex.RUnlock()

	// 返回副本
	history := make([]map[string]interface{}, len(logHistory))
	copy(history, logHistory)
	return history
}

