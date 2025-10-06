// Common types for OpenAI-compatible API proxy

// OpenAI request/response types
export interface Message {
  role: string;
  content: string;
}

export interface OpenAIRequest {
  model: string;
  messages: Message[];
  stream?: boolean;
  temperature?: number;
  max_tokens?: number;
  top_p?: number;
  [key: string]: unknown; // Allow additional properties
}

export interface OpenAIResponse {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: Array<{
    index: number;
    message: Message;
    finish_reason: string;
  }>;
  usage: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

// Statistics types
export interface RequestStats {
  totalRequests: number;
  successfulRequests: number;
  failedRequests: number;
  lastRequestTime: Date;
  averageResponseTime: number; // in milliseconds
  apiCallsCount: number; // /v1/chat/completions calls
  modelsCallsCount: number; // /v1/models calls
  streamingRequests: number; // streaming mode requests
  nonStreamingRequests: number; // non-streaming mode requests
  startTime: Date; // server start time
  fastestResponse: number; // fastest response time in ms
  slowestResponse: number; // slowest response time in ms
  modelUsage: Map<string, number>; // model name -> count
}

export interface LiveRequest {
  id: string;
  timestamp: Date;
  method: string;
  path: string;
  status: number;
  duration: number;
  userAgent: string;
  model?: string;
}

// Configuration interface
export interface ProxyConfig {
  // Server settings
  port: number;
  debugMode: boolean;
  defaultStream: boolean;
  dashboardEnabled: boolean;

  // API settings
  upstreamUrl: string;
  sessionToken?: string;
  sessionCookie?: string;
  defaultKey: string;
  modelName: string;
  defaultTemperature?: number;
  defaultMaxTokens?: number;
  defaultExtraBody?: unknown;

  // Branding
  serviceName: string;
  serviceEmoji: string;
  footerText: string;
  discussionUrl: string;
  githubRepo: string;

  // Optional features
  enableThinking?: boolean;
  thinkTagsMode?: "strip" | "think" | "raw";

  // SEO settings
  seoTitle?: string;
  seoDescription?: string;
  seoKeywords?: string;
  seoAuthor?: string;
  seoOgImage?: string;
}

// i18n Language type
export type Language = "zh-CN" | "en-US" | "ja-JP";

// i18n translations interface
export interface I18nTranslations {
  // Common
  home: string;
  docs: string;
  playground: string;
  deploy: string;
  dashboard: string;

  // Homepage
  homeTitle: string;
  homeSubtitle: string;
  homeStatusTitle: string;
  homeApiEndpoint: string;
  homeModel: string;
  homeApiKey: string;
  homePort: string;
  homeRuntime: string;
  homeQuickStart: string;
  homeViewDocs: string;
  homeOnlineTest: string;
  homeDeployNow: string;
  homeViewStats: string;
  homeDocsDesc: string;
  homePlaygroundDesc: string;
  homeDeployDesc: string;
  homeDashboardDesc: string;
  homeFooterPoweredBy: string;

  // Docs page
  docsTitle: string;
  docsSubtitle: string;
  docsAuth: string;
  docsAuthDesc: string;
  docsGetModels: string;
  docsChatNonStreaming: string;
  docsChatStreaming: string;
  docsUseSdk: string;

  // Playground page
  playgroundTitle: string;
  playgroundSubtitle: string;
  playgroundRequestConfig: string;
  playgroundApiKey: string;
  playgroundModel: string;
  playgroundModelDesc: string;
  playgroundTemperature: string;
  playgroundMaxTokens: string;
  playgroundStream: string;
  playgroundStreamDesc: string;
  playgroundSystem: string;
  playgroundSystemDesc: string;
  playgroundSystemPlaceholder: string;
  playgroundMessage: string;
  playgroundMessageDesc: string;
  playgroundMessagePlaceholder: string;
  playgroundSendBtn: string;
  playgroundSending: string;
  playgroundClearBtn: string;
  playgroundResponseTitle: string;
  playgroundResponseContent: string;
  playgroundCopyBtn: string;
  playgroundCopied: string;
  playgroundEmptyTitle: string;
  playgroundEmptyDesc: string;
  playgroundLoadingTitle: string;
  playgroundLoadingDesc: string;
  playgroundErrorTitle: string;
  playgroundErrorEmpty: string;
  playgroundDuration: string;
  playgroundStatus: string;

  // Deploy page
  deployTitle: string;
  deploySubtitle: string;
  deployQuickDeploy: string;
  deployGithubTitle: string;
  deployGithubStep1: string;
  deployGithubStep2: string;
  deployGithubStep3: string;
  deployGithubStep4: string;
  deployGithubStep5: string;
  deployGithubStep6: string;
  deployCliTitle: string;
  deployCliDesc: string;
  deployEnvTitle: string;
  deployEnvName: string;
  deployEnvDesc: string;
  deployEnvDefault: string;
  deployLocalTitle: string;

  // Dashboard page
  dashboardTitle: string;
  dashboardSubtitle: string;
  dashboardOverview: string;
  dashboardTotalRequests: string;
  dashboardSuccessRate: string;
  dashboardAvgResponseTime: string;
  dashboardUptime: string;
  dashboardTopModels: string;
  dashboardRecentRequests: string;
  dashboardTime: string;
  dashboardMethod: string;
  dashboardPath: string;
  dashboardStatusCode: string;
  dashboardDuration: string;
  dashboardModel: string;
  dashboardNoRequests: string;
  dashboardSuccessRequests: string;
  dashboardFailedRequests: string;
  dashboardAvgResponseTimeMs: string;
  dashboardApiCalls: string;
  dashboardPerformance: string;
  dashboardFastest: string;
  dashboardSlowest: string;
  dashboardSystemInfo: string;
  dashboardStreamingReq: string;
  dashboardNonStreamingReq: string;
  dashboardNoData: string;
  homeApiDocsTitle: string;
  homeApiDocsDesc: string;
  homePlaygroundTitle: string;
  homeDeployTitle: string;
  homeDashboardTitle: string;
  homeModelsListTitle: string;
  homeModelsListDesc: string;

  // Footer
  footerPoweredBy: string;
  footerDiscussion: string;
  footerGithub: string;
  footerBackHome: string;
}
