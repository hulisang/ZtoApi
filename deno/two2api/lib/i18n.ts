// i18n translations for OpenAI-compatible API proxy

import type { Language, I18nTranslations } from "./types.ts";

// Chinese (Simplified) translations
const zhCN: I18nTranslations = {
  // Common
  home: "首页",
  docs: "文档",
  playground: "Playground",
  deploy: "部署",
  dashboard: "Dashboard",

  // Homepage
  homeTitle: "首页",
  homeSubtitle: "OpenAI 兼容的 API 代理服务",
  homeStatusTitle: "服务状态",
  homeApiEndpoint: "API 端点",
  homeModel: "模型",
  homeApiKey: "API 密钥",
  homePort: "端口",
  homeRuntime: "运行时",
  homeQuickStart: "快速开始",
  homeViewDocs: "查看完整 API 文档",
  homeOnlineTest: "在线测试 API",
  homeDeployNow: "一键部署",
  homeViewStats: "查看统计数据",
  homeDocsDesc: "完整的 API 使用文档和示例",
  homePlaygroundDesc: "在线测试 API 请求和响应",
  homeDeployDesc: "快速部署到 Deno Deploy",
  homeDashboardDesc: "实时监控和统计信息",
  homeFooterPoweredBy: "OpenAI 兼容 API",
  homeApiDocsTitle: "API 文档",
  homeApiDocsDesc: "完整的使用文档和代码示例",
  homePlaygroundTitle: "Playground",
  homeDeployTitle: "部署指南",
  homeDashboardTitle: "Dashboard",
  homeModelsListTitle: "模型列表",
  homeModelsListDesc: "查看所有可用的 AI 模型",

  // Docs page
  docsTitle: "API 文档",
  docsSubtitle: "完整的使用指南和代码示例",
  docsAuth: "认证",
  docsAuthDesc: "所有请求需要在 Header 中携带 Bearer Token：",
  docsGetModels: "获取模型列表",
  docsChatNonStreaming: "聊天完成（非流式）",
  docsChatStreaming: "聊天完成（流式）",
  docsUseSdk: "使用 OpenAI SDK",

  // Playground page
  playgroundTitle: "Playground",
  playgroundSubtitle: "在线测试 API",
  playgroundRequestConfig: "请求配置",
  playgroundApiKey: "API Key",
  playgroundModel: "模型",
  playgroundModelDesc: "支持的模型名称",
  playgroundTemperature: "Temperature",
  playgroundMaxTokens: "Max Tokens",
  playgroundStream: "启用流式响应",
  playgroundStreamDesc: "实时流式输出",
  playgroundSystem: "System (可选)",
  playgroundSystemDesc: "系统提示词，用于设定角色和行为",
  playgroundSystemPlaceholder: "你是一个乐于助人的AI助手...",
  playgroundMessage: "Message",
  playgroundMessageDesc: "用户消息内容",
  playgroundMessagePlaceholder: "输入你的问题...",
  playgroundSendBtn: "🚀 发送请求",
  playgroundSending: "⏳ 请求中...",
  playgroundClearBtn: "🗑️ 清空响应",
  playgroundResponseTitle: "响应结果",
  playgroundResponseContent: "响应内容",
  playgroundCopyBtn: "📋 复制",
  playgroundCopied: "✅ 已复制",
  playgroundEmptyTitle: "等待请求",
  playgroundEmptyDesc: "配置参数后点击\"发送请求\"开始测试",
  playgroundLoadingTitle: "正在请求中...",
  playgroundLoadingDesc: "请稍候",
  playgroundErrorTitle: "请求失败",
  playgroundErrorEmpty: "请输入消息内容",
  playgroundDuration: "耗时",
  playgroundStatus: "状态",

  // Deploy page
  deployTitle: "部署指南",
  deploySubtitle: "一键部署到 Deno Deploy",
  deployQuickDeploy: "快速部署",
  deployGithubTitle: "方式一：GitHub 自动部署",
  deployGithubStep1: "Fork 本项目到你的 GitHub",
  deployGithubStep2: "访问 Deno Deploy",
  deployGithubStep3: "创建新项目，选择你的 GitHub 仓库",
  deployGithubStep4: "设置入口文件（参考项目 README）",
  deployGithubStep5: "配置环境变量",
  deployGithubStep6: "部署完成！",
  deployCliTitle: "命令行部署",
  deployCliDesc: "使用 deployctl 进行部署：",
  deployEnvTitle: "环境变量配置",
  deployEnvName: "变量名",
  deployEnvDesc: "说明",
  deployEnvDefault: "默认值",
  deployLocalTitle: "本地运行",

  // Dashboard page
  dashboardTitle: "Dashboard",
  dashboardSubtitle: "实时监控和统计",
  dashboardOverview: "概览",
  dashboardTotalRequests: "总请求数",
  dashboardSuccessRate: "成功率",
  dashboardAvgResponseTime: "平均响应时间",
  dashboardUptime: "运行时间",
  dashboardTopModels: "热门模型",
  dashboardRecentRequests: "最近请求",
  dashboardTime: "时间",
  dashboardMethod: "方法",
  dashboardPath: "路径",
  dashboardStatusCode: "状态码",
  dashboardDuration: "耗时",
  dashboardModel: "模型",
  dashboardNoRequests: "暂无请求记录",
  dashboardSuccessRequests: "成功请求",
  dashboardFailedRequests: "失败请求",
  dashboardAvgResponseTimeMs: "平均响应时间",
  dashboardApiCalls: "API 调用",
  dashboardPerformance: "性能指标",
  dashboardFastest: "最快响应",
  dashboardSlowest: "最慢响应",
  dashboardSystemInfo: "系统信息",
  dashboardStreamingReq: "流式请求",
  dashboardNonStreamingReq: "非流式请求",
  dashboardNoData: "暂无数据",

  // Footer
  footerPoweredBy: "Powered by",
  footerDiscussion: "讨论交流",
  footerGithub: "GitHub",
  footerBackHome: "返回首页",
};

// English translations
const enUS: I18nTranslations = {
  // Common
  home: "Home",
  docs: "Docs",
  playground: "Playground",
  deploy: "Deploy",
  dashboard: "Dashboard",

  // Homepage
  homeTitle: "Home",
  homeSubtitle: "OpenAI-Compatible API Proxy Service",
  homeStatusTitle: "Service Status",
  homeApiEndpoint: "API Endpoint",
  homeModel: "Model",
  homeApiKey: "API Key",
  homePort: "Port",
  homeRuntime: "Runtime",
  homeQuickStart: "Quick Start",
  homeViewDocs: "View Full API Documentation",
  homeOnlineTest: "Test API Online",
  homeDeployNow: "Deploy Now",
  homeViewStats: "View Statistics",
  homeDocsDesc: "Complete API documentation and examples",
  homePlaygroundDesc: "Test API requests and responses online",
  homeDeployDesc: "Quick deploy to Deno Deploy",
  homeDashboardDesc: "Real-time monitoring and statistics",
  homeFooterPoweredBy: "OpenAI Compatible API",
  homeApiDocsTitle: "API Documentation",
  homeApiDocsDesc: "Complete documentation and code examples",
  homePlaygroundTitle: "Playground",
  homeDeployTitle: "Deployment Guide",
  homeDashboardTitle: "Dashboard",
  homeModelsListTitle: "Models List",
  homeModelsListDesc: "View all available AI models",

  // Docs page
  docsTitle: "API Documentation",
  docsSubtitle: "Complete usage guide and code examples",
  docsAuth: "Authentication",
  docsAuthDesc: "All requests require a Bearer Token in the header:",
  docsGetModels: "Get Models List",
  docsChatNonStreaming: "Chat Completion (Non-Streaming)",
  docsChatStreaming: "Chat Completion (Streaming)",
  docsUseSdk: "Using OpenAI SDK",

  // Playground page
  playgroundTitle: "Playground",
  playgroundSubtitle: "Test API Online",
  playgroundRequestConfig: "Request Configuration",
  playgroundApiKey: "API Key",
  playgroundModel: "Model",
  playgroundModelDesc: "Supported model name",
  playgroundTemperature: "Temperature",
  playgroundMaxTokens: "Max Tokens",
  playgroundStream: "Enable Streaming",
  playgroundStreamDesc: "Real-time streaming output",
  playgroundSystem: "System (Optional)",
  playgroundSystemDesc: "System prompt for setting role and behavior",
  playgroundSystemPlaceholder: "You are a helpful AI assistant...",
  playgroundMessage: "Message",
  playgroundMessageDesc: "User message content",
  playgroundMessagePlaceholder: "Enter your question...",
  playgroundSendBtn: "🚀 Send Request",
  playgroundSending: "⏳ Sending...",
  playgroundClearBtn: "🗑️ Clear Response",
  playgroundResponseTitle: "Response Result",
  playgroundResponseContent: "Response Content",
  playgroundCopyBtn: "📋 Copy",
  playgroundCopied: "✅ Copied",
  playgroundEmptyTitle: "Waiting for Request",
  playgroundEmptyDesc: "Configure parameters and click 'Send Request' to test",
  playgroundLoadingTitle: "Requesting...",
  playgroundLoadingDesc: "Please wait",
  playgroundErrorTitle: "Request Failed",
  playgroundErrorEmpty: "Please enter message content",
  playgroundDuration: "Duration",
  playgroundStatus: "Status",

  // Deploy page
  deployTitle: "Deployment Guide",
  deploySubtitle: "Deploy to Deno Deploy with one click",
  deployQuickDeploy: "Quick Deploy",
  deployGithubTitle: "Option 1: GitHub Auto Deploy",
  deployGithubStep1: "Fork this project to your GitHub",
  deployGithubStep2: "Visit Deno Deploy",
  deployGithubStep3: "Create new project, select your GitHub repository",
  deployGithubStep4: "Set entry file (refer to project README)",
  deployGithubStep5: "Configure environment variables",
  deployGithubStep6: "Deploy complete!",
  deployCliTitle: "Command Line Deploy",
  deployCliDesc: "Deploy using deployctl:",
  deployEnvTitle: "Environment Variables",
  deployEnvName: "Variable",
  deployEnvDesc: "Description",
  deployEnvDefault: "Default",
  deployLocalTitle: "Local Run",

  // Dashboard page
  dashboardTitle: "Dashboard",
  dashboardSubtitle: "Real-time monitoring and statistics",
  dashboardOverview: "Overview",
  dashboardTotalRequests: "Total Requests",
  dashboardSuccessRate: "Success Rate",
  dashboardAvgResponseTime: "Avg Response Time",
  dashboardUptime: "Uptime",
  dashboardTopModels: "Top Models",
  dashboardRecentRequests: "Recent Requests",
  dashboardTime: "Time",
  dashboardMethod: "Method",
  dashboardPath: "Path",
  dashboardStatusCode: "Status",
  dashboardDuration: "Duration",
  dashboardModel: "Model",
  dashboardNoRequests: "No requests yet",
  dashboardSuccessRequests: "Successful Requests",
  dashboardFailedRequests: "Failed Requests",
  dashboardAvgResponseTimeMs: "Avg Response Time",
  dashboardApiCalls: "API Calls",
  dashboardPerformance: "Performance Metrics",
  dashboardFastest: "Fastest Response",
  dashboardSlowest: "Slowest Response",
  dashboardSystemInfo: "System Info",
  dashboardStreamingReq: "Streaming Requests",
  dashboardNonStreamingReq: "Non-Streaming Requests",
  dashboardNoData: "No data available",

  // Footer
  footerPoweredBy: "Powered by",
  footerDiscussion: "Discussion",
  footerGithub: "GitHub",
  footerBackHome: "Back to Home",
};

// Japanese translations
const jaJP: I18nTranslations = {
  // Common
  home: "ホーム",
  docs: "ドキュメント",
  playground: "Playground",
  deploy: "デプロイ",
  dashboard: "Dashboard",

  // Homepage
  homeTitle: "ホーム",
  homeSubtitle: "OpenAI互換APIプロキシサービス",
  homeStatusTitle: "サービスステータス",
  homeApiEndpoint: "APIエンドポイント",
  homeModel: "モデル",
  homeApiKey: "APIキー",
  homePort: "ポート",
  homeRuntime: "ランタイム",
  homeQuickStart: "クイックスタート",
  homeViewDocs: "完全なAPIドキュメントを表示",
  homeOnlineTest: "APIをオンラインでテスト",
  homeDeployNow: "今すぐデプロイ",
  homeViewStats: "統計を表示",
  homeDocsDesc: "完全なAPIドキュメントとサンプル",
  homePlaygroundDesc: "APIリクエストとレスポンスをオンラインでテスト",
  homeDeployDesc: "Deno Deployに素早くデプロイ",
  homeDashboardDesc: "リアルタイム監視と統計情報",
  homeFooterPoweredBy: "OpenAI互換API",
  homeApiDocsTitle: "APIドキュメント",
  homeApiDocsDesc: "完全なドキュメントとコード例",
  homePlaygroundTitle: "Playground",
  homeDeployTitle: "デプロイガイド",
  homeDashboardTitle: "Dashboard",
  homeModelsListTitle: "モデルリスト",
  homeModelsListDesc: "利用可能なすべてのAIモデルを表示",

  // Docs page
  docsTitle: "APIドキュメント",
  docsSubtitle: "完全な使用ガイドとコード例",
  docsAuth: "認証",
  docsAuthDesc: "すべてのリクエストにはヘッダーにBearerトークンが必要です：",
  docsGetModels: "モデルリストを取得",
  docsChatNonStreaming: "チャット補完（非ストリーミング）",
  docsChatStreaming: "チャット補完（ストリーミング）",
  docsUseSdk: "OpenAI SDKを使用",

  // Playground page
  playgroundTitle: "Playground",
  playgroundSubtitle: "APIをオンラインでテスト",
  playgroundRequestConfig: "リクエスト設定",
  playgroundApiKey: "API Key",
  playgroundModel: "モデル",
  playgroundModelDesc: "サポートされているモデル名",
  playgroundTemperature: "Temperature",
  playgroundMaxTokens: "Max Tokens",
  playgroundStream: "ストリーミングを有効化",
  playgroundStreamDesc: "リアルタイムストリーミング出力",
  playgroundSystem: "System（オプション）",
  playgroundSystemDesc: "役割と動作を設定するシステムプロンプト",
  playgroundSystemPlaceholder: "あなたは親切なAIアシスタントです...",
  playgroundMessage: "Message",
  playgroundMessageDesc: "ユーザーメッセージの内容",
  playgroundMessagePlaceholder: "質問を入力してください...",
  playgroundSendBtn: "🚀 リクエスト送信",
  playgroundSending: "⏳ 送信中...",
  playgroundClearBtn: "🗑️ レスポンスをクリア",
  playgroundResponseTitle: "レスポンス結果",
  playgroundResponseContent: "レスポンス内容",
  playgroundCopyBtn: "📋 コピー",
  playgroundCopied: "✅ コピーしました",
  playgroundEmptyTitle: "リクエスト待機中",
  playgroundEmptyDesc: "パラメータを設定して「リクエスト送信」をクリックしてテスト開始",
  playgroundLoadingTitle: "リクエスト中...",
  playgroundLoadingDesc: "お待ちください",
  playgroundErrorTitle: "リクエスト失敗",
  playgroundErrorEmpty: "メッセージ内容を入力してください",
  playgroundDuration: "所要時間",
  playgroundStatus: "ステータス",

  // Deploy page
  deployTitle: "デプロイガイド",
  deploySubtitle: "ワンクリックでDeno Deployにデプロイ",
  deployQuickDeploy: "クイックデプロイ",
  deployGithubTitle: "方法1：GitHub自動デプロイ",
  deployGithubStep1: "このプロジェクトをあなたのGitHubにフォーク",
  deployGithubStep2: "Deno Deployにアクセス",
  deployGithubStep3: "新しいプロジェクトを作成し、GitHubリポジトリを選択",
  deployGithubStep4: "エントリーファイルを設定（プロジェクトのREADMEを参照）",
  deployGithubStep5: "環境変数を設定",
  deployGithubStep6: "デプロイ完了！",
  deployCliTitle: "コマンドラインでデプロイ",
  deployCliDesc: "deployctlを使用してデプロイ：",
  deployEnvTitle: "環境変数設定",
  deployEnvName: "変数名",
  deployEnvDesc: "説明",
  deployEnvDefault: "デフォルト値",
  deployLocalTitle: "ローカル実行",

  // Dashboard page
  dashboardTitle: "Dashboard",
  dashboardSubtitle: "リアルタイム監視と統計",
  dashboardOverview: "概要",
  dashboardTotalRequests: "総リクエスト数",
  dashboardSuccessRate: "成功率",
  dashboardAvgResponseTime: "平均レスポンス時間",
  dashboardUptime: "稼働時間",
  dashboardTopModels: "人気モデル",
  dashboardRecentRequests: "最近のリクエスト",
  dashboardTime: "時間",
  dashboardMethod: "メソッド",
  dashboardPath: "パス",
  dashboardStatusCode: "ステータス",
  dashboardDuration: "所要時間",
  dashboardModel: "モデル",
  dashboardNoRequests: "まだリクエストがありません",
  dashboardSuccessRequests: "成功リクエスト",
  dashboardFailedRequests: "失敗リクエスト",
  dashboardAvgResponseTimeMs: "平均レスポンス時間",
  dashboardApiCalls: "API呼び出し",
  dashboardPerformance: "パフォーマンス指標",
  dashboardFastest: "最速レスポンス",
  dashboardSlowest: "最遅レスポンス",
  dashboardSystemInfo: "システム情報",
  dashboardStreamingReq: "ストリーミングリクエスト",
  dashboardNonStreamingReq: "非ストリーミングリクエスト",
  dashboardNoData: "データがありません",

  // Footer
  footerPoweredBy: "Powered by",
  footerDiscussion: "ディスカッション",
  footerGithub: "GitHub",
  footerBackHome: "ホームに戻る",
};

// Translation map
const translations: Record<Language, I18nTranslations> = {
  "zh-CN": zhCN,
  "en-US": enUS,
  "ja-JP": jaJP,
};

/**
 * Get translations for a specific language
 */
export function getTranslations(lang: Language = "zh-CN"): I18nTranslations {
  return translations[lang] || translations["zh-CN"];
}

/**
 * Detect language from request headers
 */
export function detectLanguage(req: Request): Language {
  const acceptLanguage = req.headers.get("Accept-Language") || "";

  if (acceptLanguage.includes("en")) return "en-US";
  if (acceptLanguage.includes("ja")) return "ja-JP";
  if (acceptLanguage.includes("zh")) return "zh-CN";

  return "zh-CN"; // default
}

/**
 * Get language from URL query parameter or cookie
 */
export function getLanguageFromUrl(url: URL): Language | null {
  const lang = url.searchParams.get("lang");
  if (lang === "zh-CN" || lang === "en-US" || lang === "ja-JP") {
    return lang as Language;
  }
  return null;
}
