// OpenAI-Compatible API Proxy Template
// This is a configurable template that can be customized for different AI services

export {};

import type { ProxyConfig, RequestStats, LiveRequest, OpenAIRequest, Message, Language } from "./lib/types.ts";
import {
  generateBrowserHeaders,
  generateRequestId,
  generateChatId,
  recordRequest,
  debugLog,
  createErrorResponse,
  verifyAuth,
  createSSEData,
  parseSSELine,
} from "./lib/utils.ts";
import { getHomePage, getDashboardPage } from "./lib/pages.ts";
import { getDocsPage, getDeployPage } from "./pages/docs-deploy.ts";
import { getPlaygroundPage } from "./pages/playground.ts";
import { detectLanguage, getLanguageFromUrl } from "./lib/i18n.ts";

// ============================================================================
// CONFIGURATION - Customize these values for your service
// ============================================================================

const defaultExtraBody = (() => {
  const raw = Deno.env.get("DEFAULT_EXTRA_BODY");
  if (!raw) return undefined;
  try {
    return JSON.parse(raw);
  } catch (error) {
    console.warn("[WARN] 无法解析 DEFAULT_EXTRA_BODY，忽略该配置", error);
    return undefined;
  }
})();

const CONFIG: ProxyConfig = {
  // Server settings
  port: parseInt(Deno.env.get("PORT") || "9090"),
  debugMode: Deno.env.get("DEBUG_MODE") === "true",
  defaultStream: Deno.env.get("DEFAULT_STREAM") !== "false",
  dashboardEnabled: Deno.env.get("DASHBOARD_ENABLED") !== "false",

  // API settings
  upstreamUrl: Deno.env.get("UPSTREAM_URL") || "https://chatsutra-server.account-2b0.workers.dev/v2/chat/completions",
  sessionToken: Deno.env.get("X_SESSION_TOKEN") || "",
  sessionCookie: Deno.env.get("X_SESSION_COOKIE") || "",
  defaultKey: Deno.env.get("DEFAULT_KEY") || "sk-two-demo",
  modelName: Deno.env.get("MODEL_NAME") || "sutra-v2",
  defaultTemperature: Deno.env.has("DEFAULT_TEMPERATURE") ? Number(Deno.env.get("DEFAULT_TEMPERATURE")) : undefined,
  defaultMaxTokens: Deno.env.has("DEFAULT_MAX_TOKENS") ? Number(Deno.env.get("DEFAULT_MAX_TOKENS")) : undefined,
  defaultExtraBody,

  // Branding - Customize these for your service
  serviceName: Deno.env.get("SERVICE_NAME") || "Two2Api",
  serviceEmoji: Deno.env.get("SERVICE_EMOJI") || "🌀",
  footerText: Deno.env.get("FOOTER_TEXT") || "Two 体验转接 OpenAI 接口",
  discussionUrl: Deno.env.get("DISCUSSION_URL") || "https://github.com/your-repo/discussions",
  githubRepo: Deno.env.get("GITHUB_REPO") || "https://github.com/your-repo",

  // SEO settings - Customize for better search engine visibility
  seoTitle: Deno.env.get("SEO_TITLE") || "Two2Api - Two Chat OpenAI 兼容代理",
  seoDescription:
    Deno.env.get("SEO_DESCRIPTION") ||
    "将 Two Chat 的 Sutra 接口代理为 OpenAI 兼容格式，提供监控与即开即用体验",
  seoKeywords: Deno.env.get("SEO_KEYWORDS") || "Two,Sutra,OpenAI,API,Proxy,Deno",
  seoAuthor: Deno.env.get("SEO_AUTHOR") || "Two2Api",
  seoOgImage: Deno.env.get("SEO_OG_IMAGE") || "",
};

// ============================================================================
// STATISTICS
// ============================================================================

const stats: RequestStats = {
  totalRequests: 0,
  successfulRequests: 0,
  failedRequests: 0,
  lastRequestTime: new Date(),
  averageResponseTime: 0,
  apiCallsCount: 0,
  modelsCallsCount: 0,
  streamingRequests: 0,
  nonStreamingRequests: 0,
  startTime: new Date(),
  fastestResponse: Infinity,
  slowestResponse: 0,
  modelUsage: new Map<string, number>(),
};

const liveRequests: LiveRequest[] = [];

// ============================================================================
// UPSTREAM TRANSFORMERS - Customize these for your upstream API
// ============================================================================

/**
 * Transform OpenAI request to upstream format
 * Override this function to match your upstream API's expected format
 */
function transformToUpstream(openAIReq: OpenAIRequest): Record<string, unknown> {
  const payload: Record<string, unknown> = {
    model: openAIReq.model || CONFIG.modelName,
    messages: openAIReq.messages,
  };

  // Two Chat 额外参数透传
  const passthroughKeys = [
    "temperature",
    "max_tokens",
    "top_p",
    "frequency_penalty",
    "presence_penalty",
    "stop",
    "extra_body",
  ];

  for (const key of passthroughKeys) {
    const value = (openAIReq as Record<string, unknown>)[key];
    if (value !== undefined) {
      payload[key] = value;
    }
  }

  if (openAIReq.temperature === undefined && CONFIG.defaultTemperature !== undefined) {
    payload.temperature = CONFIG.defaultTemperature;
  }

  if (openAIReq.max_tokens === undefined && CONFIG.defaultMaxTokens !== undefined) {
    payload.max_tokens = CONFIG.defaultMaxTokens;
  }

  if ((openAIReq as Record<string, unknown>).extra_body === undefined && CONFIG.defaultExtraBody) {
    payload.extra_body = CONFIG.defaultExtraBody;
  }

  // 透传其余未覆盖字段，避免遗漏新能力
  for (const [key, value] of Object.entries(openAIReq)) {
    if (payload[key] !== undefined || key === "messages") continue;
    payload[key] = value;
  }

  return payload;
}

/**
 * Transform upstream response to OpenAI format
 * Override this function to match your upstream API's response format
 */
function transformFromUpstream(upstreamData: any, requestId: string): any {
  const upstreamChoices = Array.isArray(upstreamData?.choices) ? upstreamData.choices : [];
  const usage = upstreamData?.usage || {};

  const choices = upstreamChoices.map((choice: any, index: number) => ({
    index,
    message: choice.message || {
      role: "assistant",
      content: choice.delta?.content ?? choice.text ?? choice.output ?? "",
    },
    finish_reason: choice.finish_reason ?? "stop",
  }));

  if (choices.length === 0) {
    choices.push({
      index: 0,
      message: {
        role: "assistant",
        content: upstreamData?.output || upstreamData?.content || upstreamData?.text || "",
      },
      finish_reason: upstreamData?.finish_reason ?? "stop",
    });
  }

  return {
    id: upstreamData?.id || requestId,
    object: upstreamData?.object || "chat.completion",
    created: upstreamData?.created ?? Math.floor(Date.now() / 1000),
    model: upstreamData?.model || CONFIG.modelName,
    choices,
    usage: {
      prompt_tokens: usage.prompt_tokens ?? usage.prompt ?? 0,
      completion_tokens: usage.completion_tokens ?? usage.completion ?? 0,
      total_tokens:
        usage.total_tokens ?? usage.total ?? (usage.prompt_tokens ?? 0) + (usage.completion_tokens ?? 0),
    },
  };
}

/**
 * Transform message before sending to upstream
 * Override to handle special cases (e.g., system message conversion)
 */
function transformMessages(messages: Message[]): Message[] {
  // Default: pass through unchanged
  return messages;
}

/**
 * Get additional headers for upstream request
 * Override to add service-specific headers
 */
function getUpstreamHeaders(): Record<string, string> {
  const origin = new URL(CONFIG.upstreamUrl).origin;
  const headers = generateBrowserHeaders(origin, "https://chat.two.ai/");

  headers["Accept"] = "application/json";
  headers["Content-Type"] = "application/json";
  headers["Origin"] = "https://chat.two.ai";
  headers["Referer"] = "https://chat.two.ai/";

  if (CONFIG.sessionToken) {
    headers["x-session-token"] = CONFIG.sessionToken;
  }

  return headers;
}

// ============================================================================
// HANDLERS
// ============================================================================

async function handleModels(req: Request): Promise<Response> {
  const startTime = Date.now();
  const userAgent = req.headers.get("user-agent") || "unknown";

  if (!verifyAuth(req, CONFIG.defaultKey)) {
    recordRequest(stats, liveRequests, { method: "GET", path: "/v1/models", userAgent }, 401, Date.now() - startTime);
    return createErrorResponse("Unauthorized");
  }

  // Return default model list - customize as needed
  const models = {
    object: "list",
    data: [
      {
        id: CONFIG.modelName,
        object: "model",
        created: Math.floor(Date.now() / 1000),
        owned_by: "system",
      },
    ],
  };

  recordRequest(stats, liveRequests, { method: "GET", path: "/v1/models", userAgent }, 200, Date.now() - startTime);

  return new Response(JSON.stringify(models), {
    headers: { "Content-Type": "application/json" },
  });
}

async function handleChatCompletions(req: Request): Promise<Response> {
  const startTime = Date.now();
  const userAgent = req.headers.get("user-agent") || "unknown";

  if (!verifyAuth(req, CONFIG.defaultKey)) {
    recordRequest(stats, liveRequests, { method: "POST", path: "/v1/chat/completions", userAgent }, 401, Date.now() - startTime);
    return createErrorResponse("Unauthorized");
  }

  const openAIReq: OpenAIRequest = await req.json();
  const isStreaming = openAIReq.stream ?? CONFIG.defaultStream;
  const requestId = generateRequestId();

  debugLog(CONFIG.debugMode, "Request:", openAIReq);

  // Update stats
  if (isStreaming) {
    stats.streamingRequests++;
  } else {
    stats.nonStreamingRequests++;
  }

  // Transform messages
  const transformedMessages = transformMessages(openAIReq.messages);

  // Build upstream request
  const upstreamReq = transformToUpstream({
    ...openAIReq,
    messages: transformedMessages,
  });

  debugLog(CONFIG.debugMode, "Upstream request:", upstreamReq);

  try {
    const upstreamResponse = await fetch(CONFIG.upstreamUrl, {
      method: "POST",
      headers: getUpstreamHeaders(),
      body: JSON.stringify(upstreamReq),
    });

    if (!upstreamResponse.ok) {
      throw new Error(`Upstream error: ${upstreamResponse.status}`);
    }

    if (isStreaming) {
      return handleStreamingResponse(upstreamResponse, requestId, startTime, userAgent, openAIReq.model);
    } else {
      return handleNonStreamingResponse(upstreamResponse, requestId, startTime, userAgent, openAIReq.model);
    }
  } catch (error) {
    debugLog(CONFIG.debugMode, "Error:", error);
    recordRequest(
      stats,
      liveRequests,
      { method: "POST", path: "/v1/chat/completions", userAgent, model: openAIReq.model },
      500,
      Date.now() - startTime
    );
    return createErrorResponse(`Request failed: ${error}`);
  }
}

async function handleStreamingResponse(
  upstreamResponse: Response,
  requestId: string,
  startTime: number,
  userAgent: string,
  model?: string
): Promise<Response> {
  const stream = new ReadableStream({
    async start(controller) {
      const reader = upstreamResponse.body?.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      try {
        while (true) {
          const { done, value } = await reader!.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split("\n");
          buffer = lines.pop() || "";

          for (const line of lines) {
            if (!line.trim()) continue;

            const parsed = parseSSELine(line);
            if (!parsed) continue;

            const parsedChoices = (parsed as any).choices;
            const firstChoice = Array.isArray(parsedChoices) ? parsedChoices[0] : undefined;
            const deltaContent =
              firstChoice?.delta?.content ??
              firstChoice?.message?.content ??
              (parsed as any).content ??
              (parsed as any).text ??
              "";
            const finishReason = firstChoice?.finish_reason ?? (parsed as any).finish_reason ?? null;

            const chunk = {
              id: requestId,
              object: "chat.completion.chunk",
              created: Math.floor(Date.now() / 1000),
              model: (parsed as any).model || CONFIG.modelName,
              choices: [
                {
                  index: 0,
                  delta: {
                    content: deltaContent,
                  },
                  finish_reason: finishReason,
                },
              ],
            };

            controller.enqueue(new TextEncoder().encode(createSSEData(chunk)));

            if (finishReason) {
              controller.enqueue(new TextEncoder().encode("data: [DONE]\n\n"));
              break;
            }
          }
        }
      } catch (error) {
        debugLog(CONFIG.debugMode, "Streaming error:", error);
      } finally {
        controller.close();
        recordRequest(
          stats,
          liveRequests,
          { method: "POST", path: "/v1/chat/completions", userAgent, model },
          200,
          Date.now() - startTime
        );
      }
    },
  });

  return new Response(stream, {
    headers: {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
    },
  });
}

async function handleNonStreamingResponse(
  upstreamResponse: Response,
  requestId: string,
  startTime: number,
  userAgent: string,
  model?: string
): Promise<Response> {
  const upstreamData = await upstreamResponse.json();
  debugLog(CONFIG.debugMode, "Upstream response:", upstreamData);

  const openAIResponse = transformFromUpstream(upstreamData, requestId);

  recordRequest(
    stats,
    liveRequests,
    { method: "POST", path: "/v1/chat/completions", userAgent, model },
    200,
    Date.now() - startTime
  );

  return new Response(JSON.stringify(openAIResponse), {
    headers: { "Content-Type": "application/json" },
  });
}

// ============================================================================
// SERVER
// ============================================================================

async function handler(req: Request): Promise<Response> {
  const url = new URL(req.url);
  const path = url.pathname;

  // Detect language from URL parameter or browser Accept-Language header
  const urlLang = getLanguageFromUrl(url);
  const browserLang = detectLanguage(req);
  const lang: Language = urlLang || browserLang;
  const currentUrl = url.toString();

  // API endpoints
  if (path === "/v1/models") {
    return handleModels(req);
  }

  if (path === "/v1/chat/completions") {
    return handleChatCompletions(req);
  }

  // Web pages with i18n and SEO support
  if (path === "/" || path === "/index.html") {
    return new Response(getHomePage(CONFIG, lang, currentUrl), {
      headers: { "Content-Type": "text/html" },
    });
  }

  if (path === "/docs") {
    return new Response(getDocsPage(CONFIG, lang, currentUrl), {
      headers: { "Content-Type": "text/html" },
    });
  }

  if (path === "/deploy") {
    return new Response(getDeployPage(CONFIG, lang, currentUrl), {
      headers: { "Content-Type": "text/html" },
    });
  }

  if (path === "/playground") {
    return new Response(getPlaygroundPage(CONFIG, lang, currentUrl), {
      headers: { "Content-Type": "text/html" },
    });
  }

  if (path === "/dashboard" && CONFIG.dashboardEnabled) {
    return new Response(getDashboardPage(CONFIG, stats, liveRequests, lang, currentUrl), {
      headers: { "Content-Type": "text/html" },
    });
  }

  // 404
  return new Response("Not Found", { status: 404 });
}

console.log(`🚀 ${CONFIG.serviceName} ${CONFIG.serviceEmoji}`);
console.log(`🔗 Server starting on http://localhost:${CONFIG.port}`);
console.log(`📊 Dashboard: ${CONFIG.dashboardEnabled ? "enabled" : "disabled"}`);
console.log(`🤖 Model: ${CONFIG.modelName}`);
console.log(`🔑 API Key: ${CONFIG.defaultKey}`);

Deno.serve({ port: CONFIG.port }, handler);
