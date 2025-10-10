package register

import (
	"net/http"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(mux *http.ServeMux) {
	// 页面路由
	mux.HandleFunc("/register/login", HandleLoginPage)
	mux.HandleFunc("/register/", HandleMainPage)

	// API路由
	mux.HandleFunc("/register/api/login", HandleLogin)
	mux.HandleFunc("/register/api/logout", HandleLogout)
	mux.HandleFunc("/register/api/accounts", HandleGetAccounts)
	mux.HandleFunc("/register/api/stats", HandleGetStats)
	mux.HandleFunc("/register/api/accounts/delete", HandleDeleteAccount)
	mux.HandleFunc("/register/api/accounts/batch-delete", HandleBatchDeleteAccounts)
	mux.HandleFunc("/register/api/accounts/export", HandleExportAccounts)
	mux.HandleFunc("/register/api/accounts/import", HandleImportAccounts)
	mux.HandleFunc("/register/api/register/start", HandleStartRegister)
	mux.HandleFunc("/register/api/register/stop", HandleStopRegister)
	mux.HandleFunc("/register/api/register/stream", HandleRegisterStream)
	mux.HandleFunc("/register/api/config", HandleGetConfig)
	mux.HandleFunc("/register/api/config/save", HandleSaveConfig)
	mux.HandleFunc("/register/api/batch-refetch-apikey", HandleBatchRefetchAPIKEY)
	mux.HandleFunc("/register/api/batch-check-accounts", HandleBatchCheckAccounts)
	mux.HandleFunc("/register/api/delete-inactive-accounts", HandleDeleteInactiveAccounts)
	mux.HandleFunc("/register/api/refetch-apikey", HandleRefetchSingleAPIKEY)
}

