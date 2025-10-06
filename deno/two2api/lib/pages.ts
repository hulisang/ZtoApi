// Common HTML page templates for OpenAI-compatible API proxy

import type { ProxyConfig, RequestStats, LiveRequest, Language, I18nTranslations } from "./types.ts";
import { formatUptime, getTopModels } from "./utils.ts";
import { getSeoMeta, getStructuredData } from "./seo.ts";
import { getTranslations } from "./i18n.ts";

/**
 * Generate common HTML head section with SEO and i18n support
 */
function getHtmlHead(
  title: string,
  config: ProxyConfig,
  lang: Language = "zh-CN",
  pageDescription?: string,
  currentUrl?: string
): string {
  const t = getTranslations(lang);
  return `<!DOCTYPE html>
<html lang="${lang}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title} - ${config.serviceName}</title>
    ${getSeoMeta(config, title, pageDescription, currentUrl)}
    ${getStructuredData(config, currentUrl)}
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        @keyframes gradient-shift {
            0%, 100% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
        }
        @keyframes fade-in {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .gradient-animated {
            background-size: 200% 200%;
            animation: gradient-shift 10s ease infinite;
        }
        .animate-fade-in {
            animation: fade-in 0.6s ease-out forwards;
        }
        .animate-delay-1 { animation-delay: 0.1s; }
        .animate-delay-2 { animation-delay: 0.2s; }
        .animate-delay-3 { animation-delay: 0.3s; }
        .animate-delay-4 { animation-delay: 0.4s; }
        /* Language Switcher */
        .lang-switcher { position: fixed; top: 20px; right: 20px; z-index: 1000; }
        .lang-switcher select {
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.2);
            border-radius: 8px;
            padding: 8px 12px;
            font-size: 14px;
            cursor: pointer;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
    </style>
</head>`;
}

/**
 * Generate language switcher
 */
export function getLanguageSwitcher(currentLang: Language = "zh-CN"): string {
  return `
    <div class="lang-switcher">
        <select id="langSelect" onchange="window.location.href='?lang='+this.value">
            <option value="zh-CN" ${currentLang === "zh-CN" ? "selected" : ""}>🇨🇳 中文</option>
            <option value="en-US" ${currentLang === "en-US" ? "selected" : ""}>🇺🇸 English</option>
            <option value="ja-JP" ${currentLang === "ja-JP" ? "selected" : ""}>🇯🇵 日本語</option>
        </select>
    </div>
  `;
}

/**
 * Generate common navigation links with i18n
 */
function getNavLinks(currentPath: string, t: I18nTranslations): string {
  const links = [
    { href: "/", label: t.home },
    { href: "/docs", label: t.docs },
    { href: "/playground", label: t.playground },
    { href: "/deploy", label: t.deploy },
    { href: "/dashboard", label: t.dashboard },
  ];

  return links
    .map((link) => {
      const isActive = currentPath === link.href;
      const className = isActive
        ? "text-blue-600 font-semibold"
        : "text-gray-600 hover:text-blue-600 transition";
      return `<a href="${link.href}" class="${className}">${link.label}</a>`;
    })
    .join("\n                    ");
}

/**
 * Generate common footer
 */
function getFooter(config: ProxyConfig): string {
  return `<div class="text-center text-white/60 text-sm space-y-3 animate-fade-in animate-delay-4">
                <p>Powered by <span class="font-semibold text-white">Deno 🦕</span> | OpenAI Compatible API</p>
                <div class="flex justify-center items-center gap-6 text-xs">
                    <a href="${config.discussionUrl}" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors flex items-center gap-1">
                        <span>💬</span> 讨论交流
                    </a>
                    <span class="text-white/40">|</span>
                    <a href="${config.githubRepo}" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors flex items-center gap-1">
                        <span>⭐</span> GitHub
                    </a>
                </div>
                <p class="text-xs italic">${config.footerText}</p>
            </div>`;
}

/**
 * Generate homepage HTML
 */
export function getHomePage(config: ProxyConfig, lang: Language = "zh-CN", currentUrl?: string): string {
  const t = getTranslations(lang);
  return `${getHtmlHead(t.homeTitle, config, lang, config.seoDescription, currentUrl)}
<body class="min-h-screen gradient-animated bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500">
    ${getLanguageSwitcher(lang)}
    <div class="container mx-auto px-4 py-12 max-w-5xl">
        <!-- Header -->
        <div class="text-center mb-12 animate-fade-in">
            <div class="text-7xl mb-4 animate-pulse">${config.serviceEmoji}</div>
            <h1 class="text-5xl font-bold text-white mb-3 drop-shadow-lg">${config.serviceName}</h1>
            <p class="text-xl text-blue-100 max-w-2xl mx-auto">${t.homeSubtitle}</p>
        </div>

        <!-- Status Cards -->
        <div class="bg-white/10 backdrop-blur-lg rounded-2xl p-8 mb-8 border border-white/20 shadow-2xl animate-fade-in animate-delay-1">
            <div class="grid grid-cols-3 gap-6 text-center">
                <div class="text-center group cursor-default">
                    <div class="text-3xl mb-2 group-hover:scale-125 transition-transform">🎯</div>
                    <div class="text-white/60 text-sm mb-1">模型</div>
                    <div class="text-white font-semibold">${config.modelName}</div>
                </div>
                <div class="text-center group cursor-default">
                    <div class="text-3xl mb-2 group-hover:scale-125 transition-transform">🔌</div>
                    <div class="text-white/60 text-sm mb-1">端口</div>
                    <div class="text-white font-semibold font-mono">${config.port}</div>
                </div>
                <div class="text-center group cursor-default">
                    <div class="text-3xl mb-2 group-hover:scale-125 transition-transform">⚡</div>
                    <div class="text-white/60 text-sm mb-1">运行时</div>
                    <div class="text-white font-semibold">Deno</div>
                </div>
            </div>
        </div>

        <!-- Navigation Cards -->
        <div class="grid md:grid-cols-5 gap-6 mb-8">
            <a href="/docs" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl animate-fade-in animate-delay-2">
                <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">📖</div>
                <h3 class="text-white text-xl font-bold mb-2">API 文档</h3>
                <p class="text-blue-100 text-sm">完整的使用文档和代码示例</p>
            </a>

            <a href="/playground" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl animate-fade-in animate-delay-3">
                <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">🎮</div>
                <h3 class="text-white text-xl font-bold mb-2">Playground</h3>
                <p class="text-blue-100 text-sm">在线测试 API 请求和响应</p>
            </a>

            <a href="/deploy" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl animate-fade-in animate-delay-3">
                <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">🚀</div>
                <h3 class="text-white text-xl font-bold mb-2">部署指南</h3>
                <p class="text-blue-100 text-sm">快速部署到 Deno Deploy</p>
            </a>

            <a href="/dashboard" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl animate-fade-in animate-delay-4">
                <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">📊</div>
                <h3 class="text-white text-xl font-bold mb-2">Dashboard</h3>
                <p class="text-blue-100 text-sm">实时监控和性能统计分析</p>
            </a>

            <a href="/v1/models" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl animate-fade-in animate-delay-4">
                <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">🤖</div>
                <h3 class="text-white text-xl font-bold mb-2">模型列表</h3>
                <p class="text-blue-100 text-sm">查看所有可用的 AI 模型</p>
            </a>
        </div>

        <!-- Quick Start -->
        <div class="bg-gradient-to-r from-white/15 to-white/10 backdrop-blur-lg rounded-xl p-6 mb-8 border border-white/20 shadow-xl animate-fade-in animate-delay-4">
            <h3 class="text-white text-lg font-bold mb-3 flex items-center">
                <span class="text-2xl mr-2">🚀</span> 快速开始
            </h3>
            <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-sm overflow-x-auto">
                <pre class="text-green-300">curl -X POST http://localhost:${config.port}/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer ${config.defaultKey}" \\
  -d '{"model":"${config.modelName}","messages":[{"role":"user","content":"Hello!"}]}'</pre>
            </div>
        </div>

        <!-- Footer -->
        ${getFooter(config)}
    </div>
</body>
</html>`;
}

/**
 * Generate dashboard HTML
 */
export function getDashboardPage(
  config: ProxyConfig,
  stats: RequestStats,
  liveRequests: LiveRequest[],
  lang: Language = "zh-CN",
  currentUrl?: string
): string {
  const t = getTranslations(lang);
  const topModels = getTopModels(stats.modelUsage, 3);
  const successRate = stats.totalRequests > 0
    ? ((stats.successfulRequests / stats.totalRequests) * 100).toFixed(1)
    : "0.0";
  const uptime = formatUptime(stats.startTime);

  const recentRequests = liveRequests.slice(0, 20);

  return `${getHtmlHead("Dashboard", config, lang, config.seoDescription, currentUrl)}
<body class="min-h-screen bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 gradient-animated">
    ${getLanguageSwitcher(lang)}
    <div class="container mx-auto px-4 py-8 max-w-7xl">
        <!-- Header -->
        <div class="mb-8 flex justify-between items-center">
            <div>
                <h1 class="text-4xl font-bold text-white mb-2">📊 Dashboard</h1>
                <p class="text-blue-100">实时监控和性能统计</p>
            </div>
            <div class="text-right">
                ${getNavLinks("/dashboard", t)}
            </div>
        </div>

        <!-- Top Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-5 gap-4 mb-8">
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <div class="text-blue-200 text-sm mb-1">📈 总请求数</div>
                <div class="text-3xl font-bold text-white">${stats.totalRequests}</div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <div class="text-green-200 text-sm mb-1">✅ 成功请求</div>
                <div class="text-3xl font-bold text-white">${stats.successfulRequests}</div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <div class="text-red-200 text-sm mb-1">❌ 失败请求</div>
                <div class="text-3xl font-bold text-white">${stats.failedRequests}</div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <div class="text-yellow-200 text-sm mb-1">⚡ 平均响应时间</div>
                <div class="text-3xl font-bold text-white">${stats.averageResponseTime.toFixed(0)}<span class="text-lg">ms</span></div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <div class="text-purple-200 text-sm mb-1">🔌 API 调用</div>
                <div class="text-3xl font-bold text-white">${stats.apiCallsCount}</div>
            </div>
        </div>

        <!-- System Info -->
        <div class="grid md:grid-cols-3 gap-6 mb-8">
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h3 class="text-white font-bold mb-4 flex items-center"><span class="mr-2">⚡</span> 性能指标</h3>
                <div class="space-y-3 text-sm">
                    <div class="flex justify-between"><span class="text-blue-200">最快响应:</span><span class="text-white font-mono">${stats.fastestResponse === Infinity ? "-" : stats.fastestResponse.toFixed(0)}ms</span></div>
                    <div class="flex justify-between"><span class="text-blue-200">最慢响应:</span><span class="text-white font-mono">${stats.slowestResponse.toFixed(0)}ms</span></div>
                    <div class="flex justify-between"><span class="text-blue-200">成功率:</span><span class="text-white font-mono">${successRate}%</span></div>
                </div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h3 class="text-white font-bold mb-4 flex items-center"><span class="mr-2">📊</span> 系统信息</h3>
                <div class="space-y-3 text-sm">
                    <div class="flex justify-between"><span class="text-blue-200">运行时长:</span><span class="text-white">${uptime}</span></div>
                    <div class="flex justify-between"><span class="text-blue-200">流式请求:</span><span class="text-white">${stats.streamingRequests}</span></div>
                    <div class="flex justify-between"><span class="text-blue-200">非流式请求:</span><span class="text-white">${stats.nonStreamingRequests}</span></div>
                </div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h3 class="text-white font-bold mb-4 flex items-center"><span class="mr-2">🏆</span> 热门模型 Top 3</h3>
                <div class="space-y-2">
                    ${topModels.map(([model, count], idx) => `
                    <div class="flex items-center justify-between text-sm">
                        <span class="text-white">${["🥇", "🥈", "🥉"][idx]} ${model}</span>
                        <span class="bg-white/20 px-2 py-1 rounded text-white text-xs">${count}</span>
                    </div>`).join("")}
                    ${topModels.length === 0 ? '<div class="text-blue-200 text-sm text-center">暂无数据</div>' : ""}
                </div>
            </div>
        </div>

        <!-- Recent Requests -->
        <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
            <h3 class="text-white font-bold mb-4">📋 最近请求</h3>
            <div class="overflow-x-auto">
                <table class="w-full text-sm">
                    <thead>
                        <tr class="text-blue-200 border-b border-white/20">
                            <th class="text-left py-2 px-2">时间</th>
                            <th class="text-left py-2 px-2">方法</th>
                            <th class="text-left py-2 px-2">路径</th>
                            <th class="text-left py-2 px-2">状态</th>
                            <th class="text-left py-2 px-2">耗时</th>
                            <th class="text-left py-2 px-2">模型</th>
                        </tr>
                    </thead>
                    <tbody class="text-white">
                        ${recentRequests.map((req) => `
                        <tr class="border-b border-white/10 hover:bg-white/5">
                            <td class="py-2 px-2 font-mono text-xs">${req.timestamp.toLocaleTimeString()}</td>
                            <td class="py-2 px-2"><span class="bg-blue-500/30 px-2 py-1 rounded text-xs">${req.method}</span></td>
                            <td class="py-2 px-2 font-mono text-xs">${req.path}</td>
                            <td class="py-2 px-2"><span class="bg-${req.status >= 200 && req.status < 300 ? "green" : "red"}-500/30 px-2 py-1 rounded text-xs">${req.status}</span></td>
                            <td class="py-2 px-2 font-mono text-xs">${req.duration.toFixed(0)}ms</td>
                            <td class="py-2 px-2 text-xs">${req.model || "-"}</td>
                        </tr>`).join("")}
                        ${recentRequests.length === 0 ? '<tr><td colspan="6" class="text-center py-4 text-blue-200">暂无请求数据</td></tr>' : ""}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</body>
</html>`;
}

// Export simpler page generators
export { getHtmlHead, getNavLinks, getFooter };
