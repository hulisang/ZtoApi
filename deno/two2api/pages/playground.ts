// Playground page HTML generator
// This is a separate file to keep pages.ts clean

import type { ProxyConfig, Language } from "../lib/types.ts";
import { getTranslations } from "../lib/i18n.ts";
import { getSeoMeta } from "../lib/seo.ts";

/**
 * Generate Playground page HTML
 * A simple interactive testing page for the API with Markdown rendering
 */
export function getPlaygroundPage(config: ProxyConfig, lang: Language = "zh-CN", currentUrl?: string): string {
  const t = getTranslations(lang);

  return `<!DOCTYPE html>
<html lang="${lang}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${t.playgroundTitle} - ${config.serviceName}</title>
    ${getSeoMeta(config, t.playgroundTitle, t.playgroundSubtitle, currentUrl)}
    <script src="https://cdn.tailwindcss.com"></script>
    <!-- Marked.js for Markdown parsing -->
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <!-- Highlight.js for code syntax highlighting -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11/build/styles/github.min.css">
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11/build/highlight.min.js"></script>
    <style>
        #response { line-height: 1.6; }
        #response h1, #response h2, #response h3 { font-weight: bold; margin-top: 1em; margin-bottom: 0.5em; }
        #response h1 { font-size: 1.5em; }
        #response h2 { font-size: 1.3em; }
        #response h3 { font-size: 1.1em; }
        #response p { margin-bottom: 0.8em; }
        #response ul, #response ol { margin-left: 1.5em; margin-bottom: 0.8em; }
        #response li { margin-bottom: 0.3em; }
        #response pre { background: #f6f8fa; padding: 1em; border-radius: 0.375rem; overflow-x: auto; margin-bottom: 1em; }
        #response code { background: #f6f8fa; padding: 0.2em 0.4em; border-radius: 0.25rem; font-size: 0.9em; }
        #response pre code { background: transparent; padding: 0; }
        #response blockquote { border-left: 4px solid #ddd; padding-left: 1em; margin: 1em 0; color: #666; }
        #response a { color: #3b82f6; text-decoration: underline; }
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
</head>
<body class="bg-gray-50">
    <!-- Language Switcher -->
    <div class="lang-switcher">
        <select id="langSelect" onchange="window.location.href='?lang='+this.value">
            <option value="zh-CN" ${lang === "zh-CN" ? "selected" : ""}>🇨🇳 中文</option>
            <option value="en-US" ${lang === "en-US" ? "selected" : ""}>🇺🇸 English</option>
            <option value="ja-JP" ${lang === "ja-JP" ? "selected" : ""}>🇯🇵 日本語</option>
        </select>
    </div>
    <!-- Navigation -->
    <nav class="bg-white shadow-sm border-b">
        <div class="container mx-auto px-4 py-4">
            <div class="flex items-center justify-between">
                <a href="/" class="flex items-center space-x-2 text-blue-600 hover:text-blue-700 transition">
                    <span class="text-2xl">${config.serviceEmoji}</span>
                    <span class="text-xl font-bold">${config.serviceName}</span>
                </a>
                <div class="flex space-x-6">
                    <a href="/?lang=${lang}" class="text-gray-600 hover:text-blue-600 transition">${t.home}</a>
                    <a href="/docs?lang=${lang}" class="text-gray-600 hover:text-blue-600 transition">${t.docs}</a>
                    <a href="/playground?lang=${lang}" class="text-blue-600 font-semibold">${t.playground}</a>
                    <a href="/deploy?lang=${lang}" class="text-gray-600 hover:text-blue-600 transition">${t.deploy}</a>
                    <a href="/dashboard?lang=${lang}" class="text-gray-600 hover:text-blue-600 transition">${t.dashboard}</a>
                </div>
            </div>
        </div>
    </nav>

    <!-- Main Content -->
    <div class="container mx-auto px-4 py-8 max-w-7xl">
        <div class="text-center mb-8">
            <h1 class="text-4xl font-bold text-gray-900 mb-3">🎮 Playground</h1>
            <p class="text-gray-600">在线测试 ${config.serviceName} API</p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Request Panel -->
            <div class="bg-white rounded-xl shadow-sm border p-6">
                <h2 class="text-xl font-bold text-gray-900 mb-4 flex items-center">
                    <span class="text-2xl mr-2">📤</span> 请求配置
                </h2>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">API Key</label>
                    <input type="text" id="apiKey" value="${config.defaultKey}"
                           class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">模型</label>
                    <input type="text" id="model" value="${config.modelName}"
                           class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                    <p class="text-xs text-gray-500 mt-1">支持的模型名称</p>
                </div>

                <div class="grid grid-cols-2 gap-3 mb-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Temperature</label>
                        <input type="number" id="temperature" min="0" max="2" step="0.1" value="0.7"
                               class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Max Tokens</label>
                        <input type="number" id="maxTokens" min="1" max="4096" step="1" value="2048"
                               class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                    </div>
                </div>

                <div class="mb-4">
                    <label class="flex items-center">
                        <input type="checkbox" id="stream" checked class="mr-2">
                        <span class="text-sm font-medium text-gray-700">启用流式响应</span>
                    </label>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">System (可选)</label>
                    <textarea id="system" rows="3"
                              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 text-sm"
                              placeholder="你是一个乐于助人的AI助手..."></textarea>
                    <p class="text-xs text-gray-500 mt-1">系统提示词，用于设定角色和行为</p>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Message</label>
                    <textarea id="message" rows="6"
                              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 text-sm"
                              placeholder="输入你的问题...">你好，请介绍一下你自己</textarea>
                    <p class="text-xs text-gray-500 mt-1">用户消息内容</p>
                </div>

                <button id="sendBtn"
                        class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed">
                    🚀 发送请求
                </button>

                <button id="clearBtn"
                        class="w-full mt-2 bg-gray-200 hover:bg-gray-300 text-gray-700 font-bold py-2 px-4 rounded-lg transition">
                    🗑️ 清空响应
                </button>
            </div>

            <!-- Response Panel -->
            <div class="bg-white rounded-xl shadow-sm border p-6">
                <h2 class="text-xl font-bold text-gray-900 mb-4 flex items-center">
                    <span class="text-2xl mr-2">📥</span> 响应结果
                </h2>

                <div id="status" class="mb-4 p-3 bg-gray-100 rounded-lg hidden">
                    <span class="font-mono text-sm"></span>
                </div>

                <div class="mb-4">
                    <div class="flex items-center justify-between mb-2">
                        <label class="block text-sm font-medium text-gray-700">响应内容</label>
                        <button id="copyBtn" class="text-xs text-blue-600 hover:text-blue-700 hidden">📋 复制</button>
                    </div>
                    <div id="response"
                         class="w-full h-96 px-3 py-2 border border-gray-300 rounded-lg bg-white text-sm overflow-auto">
                        <!-- Empty state -->
                        <div id="emptyState" class="flex flex-col items-center justify-center h-full text-gray-400">
                            <div class="text-6xl mb-4">💬</div>
                            <p class="text-lg font-medium mb-2">等待请求</p>
                            <p class="text-sm">配置参数后点击"发送请求"开始测试</p>
                        </div>
                        <!-- Loading state -->
                        <div id="loadingState" class="hidden flex-col items-center justify-center h-full">
                            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
                            <p class="text-gray-600 font-medium">正在请求中...</p>
                            <p class="text-gray-400 text-sm mt-1">请稍候</p>
                        </div>
                        <!-- Error state -->
                        <div id="errorState" class="hidden flex-col items-center justify-center h-full text-red-600">
                            <div class="text-6xl mb-4">❌</div>
                            <p class="text-lg font-medium mb-2">请求失败</p>
                            <p id="errorMessage" class="text-sm text-gray-600 text-center px-4"></p>
                        </div>
                        <!-- Content area -->
                        <div id="contentArea" class="hidden"></div>
                    </div>
                </div>

                <div id="stats" class="grid grid-cols-2 gap-3 hidden">
                    <div class="bg-blue-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-600">耗时</p>
                        <p id="duration" class="text-lg font-bold text-blue-600">-</p>
                    </div>
                    <div class="bg-green-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-600">状态</p>
                        <p id="statusCode" class="text-lg font-bold text-green-600">-</p>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <footer class="bg-white border-t mt-12 py-6">
        <div class="container mx-auto px-4 text-center text-gray-500 text-sm">
            <p>Powered by <span class="font-semibold">Deno 🦕</span> | <a href="/" class="text-blue-600 hover:underline">返回首页</a> | <a href="${config.githubRepo}" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:underline">⭐ GitHub</a></p>
        </div>
    </footer>

    <script>
        const sendBtn = document.getElementById('sendBtn');
        const clearBtn = document.getElementById('clearBtn');
        const copyBtn = document.getElementById('copyBtn');
        const responseDiv = document.getElementById('response');
        const statusDiv = document.getElementById('status');
        const statsDiv = document.getElementById('stats');
        const emptyState = document.getElementById('emptyState');
        const loadingState = document.getElementById('loadingState');
        const errorState = document.getElementById('errorState');
        const contentArea = document.getElementById('contentArea');
        const errorMessage = document.getElementById('errorMessage');

        // Show specific state
        function showState(state) {
            emptyState.classList.add('hidden');
            loadingState.classList.add('hidden');
            errorState.classList.add('hidden');
            contentArea.classList.add('hidden');

            emptyState.classList.remove('flex');
            loadingState.classList.remove('flex');
            errorState.classList.remove('flex');

            if (state === 'empty') {
                emptyState.classList.remove('hidden');
                emptyState.classList.add('flex');
            } else if (state === 'loading') {
                loadingState.classList.remove('hidden');
                loadingState.classList.add('flex');
            } else if (state === 'error') {
                errorState.classList.remove('hidden');
                errorState.classList.add('flex');
            } else if (state === 'content') {
                contentArea.classList.remove('hidden');
            }
        }

        // Send request
        sendBtn.addEventListener('click', async () => {
            const apiKey = document.getElementById('apiKey').value;
            const model = document.getElementById('model').value;
            const stream = document.getElementById('stream').checked;
            const temperature = parseFloat(document.getElementById('temperature').value);
            const maxTokens = parseInt(document.getElementById('maxTokens').value);
            const systemText = document.getElementById('system').value.trim();
            const messageText = document.getElementById('message').value.trim();

            if (!messageText) {
                showState('error');
                errorMessage.textContent = '请输入消息内容';
                return;
            }

            // Build messages array
            const messages = [];
            if (systemText) {
                messages.push({ role: 'system', content: systemText });
            }
            messages.push({ role: 'user', content: messageText });

            sendBtn.disabled = true;
            sendBtn.textContent = '⏳ 请求中...';
            showState('loading');
            statusDiv.classList.add('hidden');
            statsDiv.classList.add('hidden');
            copyBtn.classList.add('hidden');

            const startTime = Date.now();

            try {
                const requestBody = { model, messages, stream };

                // Add optional parameters
                if (temperature !== 0.7) requestBody.temperature = temperature;
                if (maxTokens !== 2048) requestBody.max_tokens = maxTokens;

                const response = await fetch('/v1/chat/completions', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': \`Bearer \${apiKey}\`
                    },
                    body: JSON.stringify(requestBody)
                });

                const duration = Date.now() - startTime;

                // Show status
                statusDiv.classList.remove('hidden');
                statusDiv.querySelector('span').textContent = \`HTTP \${response.status} \${response.statusText}\`;
                statusDiv.className = \`mb-4 p-3 rounded-lg \${response.ok ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}\`;

                // Show stats
                statsDiv.classList.remove('hidden');
                document.getElementById('duration').textContent = duration + 'ms';
                document.getElementById('statusCode').textContent = response.status;

                if (stream && response.ok) {
                    // Handle streaming response with Markdown rendering
                    const reader = response.body.getReader();
                    const decoder = new TextDecoder();
                    let fullText = '';

                    showState('content');

                    // Configure marked for better rendering
                    marked.setOptions({
                        breaks: true,
                        gfm: true,
                        highlight: function(code, lang) {
                            if (lang && hljs.getLanguage(lang)) {
                                try {
                                    return hljs.highlight(code, { language: lang }).value;
                                } catch (e) {}
                            }
                            return hljs.highlightAuto(code).value;
                        }
                    });

                    while (true) {
                        const { done, value } = await reader.read();
                        if (done) break;

                        const chunk = decoder.decode(value);
                        const lines = chunk.split('\\n');

                        for (const line of lines) {
                            if (line.startsWith('data: ')) {
                                const data = line.slice(6);
                                if (data === '[DONE]') continue;

                                try {
                                    const json = JSON.parse(data);
                                    const content = json.choices[0]?.delta?.content || '';
                                    fullText += content;
                                    // Render as Markdown
                                    contentArea.innerHTML = marked.parse(fullText);
                                    // Apply syntax highlighting
                                    contentArea.querySelectorAll('pre code').forEach(block => {
                                        hljs.highlightElement(block);
                                    });
                                } catch (e) {
                                    // Ignore parse errors
                                }
                            }
                        }
                    }

                    copyBtn.classList.remove('hidden');
                } else {
                    // Handle non-streaming response with Markdown rendering
                    const data = await response.json();
                    if (data.choices && data.choices[0]?.message?.content) {
                        showState('content');
                        // Render response content as Markdown
                        const content = data.choices[0].message.content;
                        marked.setOptions({
                            breaks: true,
                            gfm: true,
                            highlight: function(code, lang) {
                                if (lang && hljs.getLanguage(lang)) {
                                    try {
                                        return hljs.highlight(code, { language: lang }).value;
                                    } catch (e) {}
                                }
                                return hljs.highlightAuto(code).value;
                            }
                        });
                        contentArea.innerHTML = marked.parse(content);
                        contentArea.querySelectorAll('pre code').forEach(block => {
                            hljs.highlightElement(block);
                        });
                    } else {
                        // Show raw JSON for errors or unexpected responses
                        showState('content');
                        contentArea.innerHTML = '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
                    }
                    copyBtn.classList.remove('hidden');
                }
            } catch (error) {
                showState('error');
                errorMessage.textContent = error.message;
                statusDiv.classList.remove('hidden');
                statusDiv.className = 'mb-4 p-3 rounded-lg bg-red-100 text-red-700';
                statusDiv.querySelector('span').textContent = 'Error: ' + error.message;
            } finally {
                sendBtn.disabled = false;
                sendBtn.textContent = '🚀 发送请求';
            }
        });

        // Clear response
        clearBtn.addEventListener('click', () => {
            showState('empty');
            statusDiv.classList.add('hidden');
            statsDiv.classList.add('hidden');
            copyBtn.classList.add('hidden');
        });

        // Copy response (copy plain text, not HTML)
        copyBtn.addEventListener('click', () => {
            navigator.clipboard.writeText(contentArea.innerText);
            const originalText = copyBtn.textContent;
            copyBtn.textContent = '✅ 已复制';
            setTimeout(() => {
                copyBtn.textContent = originalText;
            }, 2000);
        });

        // Support Enter key to send (Ctrl+Enter in textarea)
        document.getElementById('message').addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.key === 'Enter') {
                sendBtn.click();
            }
        });
    </script>
</body>
</html>`;
}
