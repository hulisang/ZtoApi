// Common utility functions for OpenAI-compatible API proxy

import type { RequestStats, LiveRequest } from "./types.ts";

/**
 * Generate random browser headers to avoid detection
 */
export function generateBrowserHeaders(origin: string, referer?: string): Record<string, string> {
  const chromeVersion = Math.floor(Math.random() * 13) + 128; // 128-140
  const edgeVersion = chromeVersion;

  const userAgents = [
    `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/${chromeVersion}.0.0.0 Safari/537.36 Edg/${edgeVersion}.0.0.0`,
    `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/${chromeVersion}.0.0.0 Safari/537.36`,
    `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/${chromeVersion}.0.0.0 Safari/537.36`,
    `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/${chromeVersion}.0.0.0 Safari/537.36`,
  ];

  const platforms = ['"Windows"', '"macOS"', '"Linux"'];
  const randomUA = userAgents[Math.floor(Math.random() * userAgents.length)];
  const randomPlatform = platforms[Math.floor(Math.random() * platforms.length)];

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    "Accept": "*/*",
    "User-Agent": randomUA,
    "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
    "Accept-Encoding": "gzip, deflate, br, zstd",
    "sec-ch-ua": `"Chromium";v="${chromeVersion}", "Not(A:Brand";v="24", "Microsoft Edge";v="${edgeVersion}"`,
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": randomPlatform,
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-origin",
    "Origin": origin,
    "Priority": "u=1, i",
  };

  if (referer) {
    headers["Referer"] = referer;
  }

  return headers;
}

/**
 * Generate unique request ID
 */
export function generateRequestId(): string {
  return `chatcmpl-${Date.now()}-${Math.random().toString(36).substring(7)}`;
}

/**
 * Generate chat/conversation ID
 */
export function generateChatId(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 15)}`;
}

/**
 * Record request statistics
 */
export function recordRequest(
  stats: RequestStats,
  liveRequests: LiveRequest[],
  req: {
    method: string;
    path: string;
    userAgent: string;
    model?: string;
  },
  status: number,
  duration: number
): void {
  stats.totalRequests++;
  stats.lastRequestTime = new Date();

  if (status >= 200 && status < 300) {
    stats.successfulRequests++;
  } else {
    stats.failedRequests++;
  }

  // Update average response time
  const totalTime = stats.averageResponseTime * (stats.totalRequests - 1) + duration;
  stats.averageResponseTime = totalTime / stats.totalRequests;

  // Update fastest/slowest
  if (duration < stats.fastestResponse) {
    stats.fastestResponse = duration;
  }
  if (duration > stats.slowestResponse) {
    stats.slowestResponse = duration;
  }

  // Track endpoint calls
  if (req.path.includes("/chat/completions")) {
    stats.apiCallsCount++;
  } else if (req.path.includes("/models")) {
    stats.modelsCallsCount++;
  }

  // Track model usage
  if (req.model) {
    const currentCount = stats.modelUsage.get(req.model) || 0;
    stats.modelUsage.set(req.model, currentCount + 1);
  }

  // Record live request
  const liveRequest: LiveRequest = {
    id: generateRequestId(),
    timestamp: new Date(),
    method: req.method,
    path: req.path,
    status,
    duration,
    userAgent: req.userAgent,
    model: req.model,
  };

  liveRequests.unshift(liveRequest);
  if (liveRequests.length > 100) {
    liveRequests.pop();
  }
}

/**
 * Format uptime duration
 */
export function formatUptime(startTime: Date): string {
  const now = new Date();
  const diff = now.getTime() - startTime.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) {
    return `${days}天 ${hours % 24}小时`;
  } else if (hours > 0) {
    return `${hours}小时 ${minutes % 60}分钟`;
  } else if (minutes > 0) {
    return `${minutes}分钟`;
  } else {
    return `${seconds}秒`;
  }
}

/**
 * Get top N models by usage
 */
export function getTopModels(modelUsage: Map<string, number>, n = 3): Array<[string, number]> {
  return Array.from(modelUsage.entries())
    .sort((a, b) => b[1] - a[1])
    .slice(0, n);
}

/**
 * Create OpenAI error response
 */
export function createErrorResponse(message: string, type = "invalid_request_error"): Response {
  return new Response(
    JSON.stringify({
      error: {
        message,
        type,
        code: null,
      },
    }),
    {
      status: 400,
      headers: { "Content-Type": "application/json" },
    }
  );
}

/**
 * Verify authentication
 */
export function verifyAuth(req: Request, validKey: string): boolean {
  const authHeader = req.headers.get("Authorization");
  if (!authHeader) return false;

  const token = authHeader.replace("Bearer ", "");
  return token === validKey;
}

/**
 * Log debug message
 */
export function debugLog(debugMode: boolean, ...args: unknown[]): void {
  if (debugMode) {
    console.log("[DEBUG]", new Date().toISOString(), ...args);
  }
}

/**
 * Create SSE data string
 */
export function createSSEData(data: unknown): string {
  return `data: ${JSON.stringify(data)}\n\n`;
}

/**
 * Parse SSE line
 */
export function parseSSELine(line: string): unknown | null {
  if (line.startsWith("data: ")) {
    const data = line.substring(6).trim();
    if (data === "[DONE]") return null;
    try {
      return JSON.parse(data);
    } catch {
      return null;
    }
  }
  return null;
}
