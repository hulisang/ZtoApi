package register

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

// 生成会话ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// 创建会话
func CreateSession() (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时过期
	}

	_, err = db.Exec(
		"INSERT INTO sessions (id, created_at, expires_at) VALUES (?, ?, ?)",
		session.ID, session.CreatedAt, session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// 验证会话
func ValidateSession(sessionID string) bool {
	var expiresAt time.Time
	err := db.QueryRow(
		"SELECT expires_at FROM sessions WHERE id = ?",
		sessionID,
	).Scan(&expiresAt)

	if err != nil {
		return false
	}

	return time.Now().Before(expiresAt)
}

// 删除会话
func DeleteSession(sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

// 清理过期会话
func CleanExpiredSessions() error {
	_, err := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}

// 检查认证
func CheckAuth(r *http.Request) bool {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return false
	}
	return ValidateSession(cookie.Value)
}

// 获取客户端IP
func getClientIP(r *http.Request) string {
	// X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

// 验证登录凭证
func ValidateLogin(username, password string) bool {
	return username == authUsername && password == authPassword
}

// 登录失败记录（简单的内存版本，可以后续扩展到数据库）
var loginAttempts = make(map[string]*struct {
	attempts   int
	lockedUntil time.Time
})

// 检查IP是否被锁定
func IsIPLocked(ip string) (bool, time.Duration) {
	if record, exists := loginAttempts[ip]; exists {
		if time.Now().Before(record.lockedUntil) {
			remaining := time.Until(record.lockedUntil)
			return true, remaining
		}
		// 过期则删除
		delete(loginAttempts, ip)
	}
	return false, 0
}

// 记录登录失败
func RecordLoginFailure(ip string) {
	if record, exists := loginAttempts[ip]; exists {
		record.attempts++
		if record.attempts >= 5 {
			record.lockedUntil = time.Now().Add(15 * time.Minute)
		}
	} else {
		loginAttempts[ip] = &struct {
			attempts   int
			lockedUntil time.Time
		}{
			attempts: 1,
		}
	}
}

// 清除登录失败记录
func ClearLoginFailure(ip string) {
	delete(loginAttempts, ip)
}

