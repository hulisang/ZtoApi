package register

import (
	"time"
)

// 账号信息
type Account struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
	APIKEY    string    `json:"apikey"`
	Status    string    `json:"status"` // active, inactive
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 会话信息
type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// 注册配置
type RegisterConfig struct {
	EmailTimeout         int    `json:"emailTimeout"`         // 邮件等待超时(秒)
	EmailCheckInterval   int    `json:"emailCheckInterval"`   // 邮件检查间隔(秒)
	RegisterDelay        int    `json:"registerDelay"`        // 注册间隔(毫秒)
	RetryTimes           int    `json:"retryTimes"`           // 重试次数
	Concurrency          int    `json:"concurrency"`          // 并发数
	HTTPTimeout          int    `json:"httpTimeout"`          // HTTP超时(秒)
	BatchSaveSize        int    `json:"batchSaveSize"`        // 批量保存大小
	ConnectionPoolSize   int    `json:"connectionPoolSize"`   // 连接池大小
	SkipAPIKey           bool   `json:"skipApikeyOnRegister"` // 快速模式
	EnableNotification   bool   `json:"enableNotification"`   // 启用通知
	PushPlusToken        string `json:"pushplusToken"`        // PushPlus Token
}

// 注册任务状态
type RegisterTask struct {
	ID          string    `json:"id"`
	Total       int       `json:"total"`
	Success     int       `json:"success"`
	Failed      int       `json:"failed"`
	IsRunning   bool      `json:"isRunning"`
	ShouldStop  bool      `json:"shouldStop"`
	StartTime   time.Time `json:"startTime"`
	CurrentMsg  string    `json:"currentMsg"`
	Config      RegisterConfig `json:"config"`
}

// 统计信息
type Stats struct {
	TotalAccounts   int64 `json:"totalAccounts"`
	WithAPIKEY      int64 `json:"withAPIKEY"`
	WithoutAPIKEY   int64 `json:"withoutAPIKEY"`
	ActiveAccounts  int64 `json:"activeAccounts"`
	InactiveAccounts int64 `json:"inactiveAccounts"`
}

// 临时邮箱域名
var EmailDomains = []string{
	"chatgptuk.pp.ua", "freemails.pp.ua", "email.gravityengine.cc", "gravityengine.cc",
	"3littlemiracles.com", "almiswelfare.org", "gyan-netra.com", "iraniandsa.org",
	"14club.org.uk", "aard.org.uk", "allumhall.co.uk", "cade.org.uk",
	"caye.org.uk", "cketrust.org", "club106.org.uk", "cok.org.uk",
	"cwetg.co.uk", "goleudy.org.uk", "hhe.org.uk", "hottchurch.org.uk",
}

// 默认配置
var DefaultConfig = RegisterConfig{
	EmailTimeout:       300,
	EmailCheckInterval: 5,
	RegisterDelay:      2000,
	RetryTimes:         3,
	Concurrency:        15,
	HTTPTimeout:        30,
	BatchSaveSize:      10,
	ConnectionPoolSize: 100,
	SkipAPIKey:         false,
	EnableNotification: false,
	PushPlusToken:      "",
}

