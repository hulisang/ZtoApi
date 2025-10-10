package register

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db              *sql.DB
	currentTask     *RegisterTask
	currentConfig   RegisterConfig
	authUsername    string
	authPassword    string
)

// 初始化注册系统
func InitRegisterSystem(dbPath string) error {
	// 确保数据目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 打开数据库
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	// 设置连接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	// 加载认证配置
	authUsername = getEnv("ZAI_USERNAME", "admin")
	authPassword = getEnv("ZAI_PASSWORD", "123456")

	// 加载注册配置
	if err := LoadConfig(); err != nil {
		log.Printf("⚠️ 加载配置失败，使用默认配置: %v", err)
		currentConfig = DefaultConfig
	}

	log.Printf("✅ 注册系统初始化成功")
	log.Printf("   - 数据库: %s", dbPath)
	log.Printf("   - 管理账号: %s", authUsername)

	return nil
}

// LoadConfig 加载配置
func LoadConfig() error {
	var configJSON string
	err := db.QueryRow("SELECT value FROM config WHERE key = 'register'").Scan(&configJSON)
	if err == sql.ErrNoRows {
		// 配置不存在，使用默认配置
		currentConfig = DefaultConfig
		return SaveConfig(currentConfig)
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(configJSON), &currentConfig); err != nil {
		return err
	}

	return nil
}

// SaveConfig 保存配置
func SaveConfig(config RegisterConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR REPLACE INTO config (key, value, updated_at)
		VALUES ('register', ?, CURRENT_TIMESTAMP)
	`, string(configJSON))
	
	if err == nil {
		currentConfig = config
	}
	
	return err
}

// GetConfig 获取当前配置
func GetConfig() RegisterConfig {
	return currentConfig
}

// 创建数据表
func createTables() error {
	// 账号表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			token TEXT NOT NULL,
			apikey TEXT,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 会话表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 配置表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 创建索引
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_accounts_apikey ON accounts(apikey)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at)`)

	return nil
}

// 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 关闭数据库
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

