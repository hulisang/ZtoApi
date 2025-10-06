// Additional page templates: Docs and Deploy pages

import type { ProxyConfig, Language } from "../lib/types.ts";
import { getHtmlHead, getNavLinks, getFooter, getLanguageSwitcher } from "../lib/pages.ts";
import { getTranslations } from "../lib/i18n.ts";
import { getSeoMeta } from "../lib/seo.ts";

/**
 * Generate API documentation page
 */
export function getDocsPage(config: ProxyConfig, lang: Language = "zh-CN", currentUrl?: string, extraSections?: string): string {
  const t = getTranslations(lang);
  return `${getHtmlHead(t.docsTitle, config, lang, config.seoDescription, currentUrl)}
<body class="min-h-screen bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 gradient-animated">
    ${getLanguageSwitcher(lang)}
    <div class="container mx-auto px-4 py-8 max-w-4xl">
        <!-- Header -->
        <div class="mb-8 flex justify-between items-center">
            <div>
                <h1 class="text-4xl font-bold text-white mb-2">📖 ${t.docsTitle}</h1>
                <p class="text-blue-100">${t.docsSubtitle}</p>
            </div>
            <div class="text-right text-sm">
                ${getNavLinks("/docs", t)}
            </div>
        </div>

        <div class="space-y-6">
            <!-- Authentication -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">🔐 认证</h2>
                <p class="text-blue-100 mb-3">所有请求需要在 Header 中携带 Bearer Token：</p>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-sm overflow-x-auto">
                    <pre class="text-green-300">Authorization: Bearer ${config.defaultKey}</pre>
                </div>
            </div>

            <!-- Get Models -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">📋 获取模型列表</h2>
                <p class="text-blue-100 mb-3">GET <code class="bg-white/20 px-2 py-1 rounded text-white">/v1/models</code></p>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-sm overflow-x-auto">
                    <pre class="text-green-300">curl http://localhost:${config.port}/v1/models \\
  -H "Authorization: Bearer ${config.defaultKey}"</pre>
                </div>
            </div>

            <!-- Chat Completions (Non-Streaming) -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">💬 聊天完成（非流式）</h2>
                <p class="text-blue-100 mb-3">POST <code class="bg-white/20 px-2 py-1 rounded text-white">/v1/chat/completions</code></p>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-xs overflow-x-auto mb-4">
                    <pre class="text-green-300">curl -X POST http://localhost:${config.port}/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer ${config.defaultKey}" \\
  -d '{
    "model": "${config.modelName}",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": false
  }'</pre>
                </div>
            </div>

            <!-- Chat Completions (Streaming) -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">🌊 聊天完成（流式）</h2>
                <p class="text-blue-100 mb-3">POST <code class="bg-white/20 px-2 py-1 rounded text-white">/v1/chat/completions</code> with <code class="bg-white/20 px-2 py-1 rounded text-white">stream: true</code></p>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-xs overflow-x-auto">
                    <pre class="text-green-300">curl -N -X POST http://localhost:${config.port}/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer ${config.defaultKey}" \\
  -d '{
    "model": "${config.modelName}",
    "messages": [
      {"role": "user", "content": "Tell me a story"}
    ],
    "stream": true
  }'</pre>
                </div>
            </div>

            <!-- SDK Examples -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">🔧 使用 OpenAI SDK</h2>

                <h3 class="text-xl font-semibold text-white mb-2 mt-4">Python</h3>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-xs overflow-x-auto mb-4">
                    <pre class="text-blue-300">from openai import OpenAI

client = OpenAI(
    api_key="${config.defaultKey}",
    base_url="http://localhost:${config.port}/v1"
)

response = client.chat.completions.create(
    model="${config.modelName}",
    messages=[{"role": "user", "content": "Hello!"}]
)

print(response.choices[0].message.content)</pre>
                </div>

                <h3 class="text-xl font-semibold text-white mb-2 mt-4">JavaScript/TypeScript</h3>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-xs overflow-x-auto">
                    <pre class="text-yellow-300">import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: '${config.defaultKey}',
  baseURL: 'http://localhost:${config.port}/v1'
});

const response = await client.chat.completions.create({
  model: '${config.modelName}',
  messages: [{ role: 'user', content: 'Hello!' }]
});

console.log(response.choices[0].message.content);</pre>
                </div>
            </div>

            ${extraSections || ""}

            <!-- Footer -->
            ${getFooter(config)}
        </div>
    </div>
</body>
</html>`;
}

/**
 * Generate deployment guide page
 */
export function getDeployPage(config: ProxyConfig, lang: Language = "zh-CN", currentUrl?: string, extraContent?: string): string {
  const t = getTranslations(lang);
  return `${getHtmlHead(t.deployTitle, config, lang, config.seoDescription, currentUrl)}
<body class="min-h-screen bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 gradient-animated">
    ${getLanguageSwitcher(lang)}
    <div class="container mx-auto px-4 py-8 max-w-4xl">
        <!-- Header -->
        <div class="mb-8 flex justify-between items-center">
            <div>
                <h1 class="text-4xl font-bold text-white mb-2">🚀 ${t.deployTitle}</h1>
                <p class="text-blue-100">${t.deploySubtitle}</p>
            </div>
            <div class="text-right text-sm">
                ${getNavLinks("/deploy", t)}
            </div>
        </div>

        <div class="space-y-6">
            <!-- Quick Deploy -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">⚡ 快速部署</h2>
                <div class="space-y-4">
                    <div>
                        <h3 class="text-lg font-semibold text-white mb-2">方式一：GitHub 自动部署</h3>
                        <ol class="text-blue-100 space-y-2 text-sm list-decimal list-inside">
                            <li>Fork 本项目到你的 GitHub</li>
                            <li>访问 <a href="https://dash.deno.com/" target="_blank" class="text-yellow-300 hover:underline">Deno Deploy</a></li>
                            <li>创建新项目，选择你的 GitHub 仓库</li>
                            <li>设置入口文件（参考项目 README）</li>
                            <li>配置环境变量</li>
                            <li>部署完成！</li>
                        </ol>
                    </div>
                </div>
            </div>

            <!-- Command Line Deploy -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">💻 命令行部署</h2>
                <div class="space-y-3">
                    <p class="text-blue-100 text-sm">使用 deployctl 进行部署：</p>
                    <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-sm overflow-x-auto">
                        <pre class="text-green-300"># 安装 deployctl
deno install -Arf jsr:@deno/deployctl

# 部署项目
deployctl deploy --project=your-project-name main.ts</pre>
                    </div>
                </div>
            </div>

            <!-- Environment Variables -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">🔧 环境变量配置</h2>
                <div class="overflow-x-auto">
                    <table class="w-full text-sm text-white">
                        <thead>
                            <tr class="border-b border-white/20">
                                <th class="text-left py-2 px-2">变量名</th>
                                <th class="text-left py-2 px-2">说明</th>
                                <th class="text-left py-2 px-2">默认值</th>
                            </tr>
                        </thead>
                        <tbody class="text-blue-100">
                            <tr class="border-b border-white/10">
                                <td class="py-2 px-2 font-mono text-xs">PORT</td>
                                <td class="py-2 px-2">服务端口</td>
                                <td class="py-2 px-2 font-mono text-xs">${config.port}</td>
                            </tr>
                            <tr class="border-b border-white/10">
                                <td class="py-2 px-2 font-mono text-xs">DEFAULT_KEY</td>
                                <td class="py-2 px-2">API 密钥</td>
                                <td class="py-2 px-2 font-mono text-xs">${config.defaultKey}</td>
                            </tr>
                            <tr class="border-b border-white/10">
                                <td class="py-2 px-2 font-mono text-xs">MODEL_NAME</td>
                                <td class="py-2 px-2">模型名称</td>
                                <td class="py-2 px-2 font-mono text-xs">${config.modelName}</td>
                            </tr>
                            <tr class="border-b border-white/10">
                                <td class="py-2 px-2 font-mono text-xs">DEBUG_MODE</td>
                                <td class="py-2 px-2">调试模式</td>
                                <td class="py-2 px-2 font-mono text-xs">${config.debugMode}</td>
                            </tr>
                            <tr class="border-b border-white/10">
                                <td class="py-2 px-2 font-mono text-xs">DASHBOARD_ENABLED</td>
                                <td class="py-2 px-2">启用监控面板</td>
                                <td class="py-2 px-2 font-mono text-xs">${config.dashboardEnabled}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Local Run -->
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20">
                <h2 class="text-2xl font-bold text-white mb-4">🏠 本地运行</h2>
                <div class="bg-gray-900/50 rounded-lg p-4 font-mono text-sm overflow-x-auto">
                    <pre class="text-green-300"># 使用 deno task
deno task start

# 开发模式（自动重载）
deno task dev

# 直接运行
deno run --allow-net --allow-env main.ts</pre>
                </div>
            </div>

            ${extraContent || ""}

            <!-- Footer -->
            ${getFooter(config)}
        </div>
    </div>
</body>
</html>`;
}
