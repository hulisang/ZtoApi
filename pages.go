package main

import (
	"fmt"
)

// 生成首页 HTML
func getHomeHTML() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZtoApi - OpenAI兼容API代理</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="min-h-screen bg-gradient-to-br from-purple-600 via-purple-700 to-indigo-800">
    <div class="container mx-auto px-4 py-12 flex items-center justify-center min-h-screen">
        <div class="max-w-4xl w-full">
            <!-- Header -->
            <div class="text-center mb-12 animate-fade-in">
                <h1 class="text-6xl font-bold text-white mb-4">
                    <span class="inline-block hover:scale-110 transition-transform">🦕</span> ZtoApi
                </h1>
                <p class="text-xl text-purple-100">OpenAI 兼容 API 代理 for Z.ai GLM-4.6</p>
            </div>

            <!-- Status Card -->
            <div class="bg-white/10 backdrop-blur-lg rounded-2xl p-8 mb-8 border border-white/20 shadow-2xl">
                <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
                    <div class="text-center">
                        <div class="text-3xl mb-2">🟢</div>
                        <div class="text-white/60 text-sm mb-1">状态</div>
                        <div class="text-white font-semibold">运行中</div>
                    </div>
                    <div class="text-center">
                        <div class="text-3xl mb-2">🤖</div>
                        <div class="text-white/60 text-sm mb-1">模型</div>
                        <div class="text-white font-semibold font-mono">%s</div>
                    </div>
                    <div class="text-center">
                        <div class="text-3xl mb-2">🔌</div>
                        <div class="text-white/60 text-sm mb-1">端口</div>
                        <div class="text-white font-semibold font-mono">%s</div>
                    </div>
                    <div class="text-center">
                        <div class="text-3xl mb-2">⚡</div>
                        <div class="text-white/60 text-sm mb-1">运行时</div>
                        <div class="text-white font-semibold">Go</div>
                    </div>
                </div>
            </div>

            <!-- Navigation Cards -->
            <div class="grid md:grid-cols-4 gap-6 mb-8">
                <a href="/docs" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl">
                    <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">📖</div>
                    <h3 class="text-white text-xl font-bold mb-2">API 文档</h3>
                    <p class="text-purple-100">查看完整的 API 使用文档和示例</p>
                </a>

                <a href="/playground" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl">
                    <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">🎮</div>
                    <h3 class="text-white text-xl font-bold mb-2">Playground</h3>
                    <p class="text-purple-100">在线测试 API 请求和响应</p>
                </a>

                <a href="/deploy" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl">
                    <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">🚀</div>
                    <h3 class="text-white text-xl font-bold mb-2">部署指南</h3>
                    <p class="text-purple-100">快速部署指南和配置说明</p>
                </a>

                <a href="/dashboard" class="group bg-white/10 backdrop-blur-lg rounded-xl p-6 border border-white/20 hover:bg-white/20 hover:border-white/40 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl">
                    <div class="text-5xl mb-4 group-hover:scale-110 transition-transform">📊</div>
                    <h3 class="text-white text-xl font-bold mb-2">Dashboard</h3>
                    <p class="text-purple-100">实时监控请求和性能统计</p>
                </a>
            </div>

            <!-- Footer -->
            <div class="text-center text-white/60 text-sm space-y-3">
                <p>Powered by <span class="font-semibold text-white">Go</span> | OpenAI Compatible API</p>
                <div class="flex justify-center items-center gap-6 text-xs">
                    <a href="https://github.com/hulisang/ZtoApi" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">
                        📦 源码地址
                    </a>
                    <span class="text-white/40">|</span>
                    <a href="https://linux.do/t/topic/1000335" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">
                        💬 交流讨论
                    </a>
                </div>
                <p class="text-white/50 text-xs italic pt-2">欲买桂花同载酒 终不似 少年游</p>
            </div>
        </div>
    </div>
</body>
</html>`, MODEL_NAME, PORT)
}

// 生成 Playground HTML
func getPlaygroundHTML() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Playground - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11/build/styles/github.min.css">
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11/build/highlight.min.js"></script>
    <style>
        #response {
            line-height: 1.6;
        }
        #response h1, #response h2, #response h3 {
            font-weight: bold;
            margin-top: 1em;
            margin-bottom: 0.5em;
        }
        #response h1 { font-size: 1.5em; }
        #response h2 { font-size: 1.3em; }
        #response h3 { font-size: 1.1em; }
        #response p { margin-bottom: 0.8em; }
        #response ul, #response ol { margin-left: 1.5em; margin-bottom: 0.8em; }
        #response li { margin-bottom: 0.3em; }
        #response pre {
            background: #f6f8fa;
            padding: 1em;
            border-radius: 0.375rem;
            overflow-x: auto;
            margin-bottom: 1em;
        }
        #response code {
            background: #f6f8fa;
            padding: 0.2em 0.4em;
            border-radius: 0.25rem;
            font-size: 0.9em;
        }
        #response pre code {
            background: transparent;
            padding: 0;
        }
    </style>
</head>
<body class="bg-gray-50">
    <nav class="bg-white shadow-sm border-b">
        <div class="container mx-auto px-4 py-4">
            <div class="flex items-center justify-between">
                <a href="/" class="flex items-center space-x-2 text-purple-600 hover:text-purple-700 transition">
                    <span class="text-2xl">🦕</span>
                    <span class="text-xl font-bold">ZtoApi</span>
                </a>
                <div class="flex items-center space-x-6">
                    <a href="/" class="text-gray-600 hover:text-purple-600 transition">首页</a>
                    <a href="/docs" class="text-gray-600 hover:text-purple-600 transition">文档</a>
                    <a href="/playground" class="text-purple-600 font-semibold">Playground</a>
                    <a href="/deploy" class="text-gray-600 hover:text-purple-600 transition">部署</a>
                    <a href="/dashboard" class="text-gray-600 hover:text-purple-600 transition">Dashboard</a>
                </div>
            </div>
        </div>
    </nav>

    <div class="container mx-auto px-4 py-8 max-w-7xl">
        <div class="text-center mb-8">
            <h1 class="text-4xl font-bold text-gray-900 mb-3">🎮 Playground</h1>
            <p class="text-gray-600">在线测试 Z.ai GLM-4.6 API 请求和响应</p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Request Panel -->
            <div class="bg-white rounded-xl shadow-sm border p-6">
                <h2 class="text-xl font-bold text-gray-900 mb-4 flex items-center">
                    <span class="text-2xl mr-2">📤</span> 请求配置
                </h2>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">API Key</label>
                    <input type="text" id="apiKey" value=""
                           class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                           placeholder="请输入你的 API Key">
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">模型</label>
                    <select id="model" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500">
                        <option value="">加载中...</option>
                    </select>
                </div>

                <div class="mb-4">
                    <label class="flex items-center">
                        <input type="checkbox" id="stream" checked class="mr-2">
                        <span class="text-sm font-medium text-gray-700">启用流式响应</span>
                    </label>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Message</label>
                    <textarea id="message" rows="6"
                              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 text-sm"
                              placeholder="输入你的问题...">你好，请介绍一下你自己</textarea>
                </div>

                <button id="sendBtn"
                        class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-4 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed">
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

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">响应内容</label>
                    <div id="response"
                         class="w-full h-96 px-3 py-2 border border-gray-300 rounded-lg bg-white text-sm overflow-auto">
                        <div id="emptyState" class="flex flex-col items-center justify-center h-full text-gray-400">
                            <div class="text-6xl mb-4">💬</div>
                            <p class="text-lg font-medium mb-2">等待请求</p>
                            <p class="text-sm">配置参数后点击"发送请求"开始测试</p>
                        </div>
                        <div id="loadingState" class="hidden flex-col items-center justify-center h-full">
                            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mb-4"></div>
                            <p class="text-gray-600 font-medium">正在请求中...</p>
                        </div>
                        <div id="contentArea" class="hidden"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        const sendBtn = document.getElementById('sendBtn');
        const clearBtn = document.getElementById('clearBtn');
        const responseDiv = document.getElementById('response');
        const emptyState = document.getElementById('emptyState');
        const loadingState = document.getElementById('loadingState');
        const contentArea = document.getElementById('contentArea');

        function showState(state) {
            emptyState.classList.add('hidden');
            loadingState.classList.add('hidden');
            contentArea.classList.add('hidden');
            emptyState.classList.remove('flex');
            loadingState.classList.remove('flex');

            if (state === 'empty') {
                emptyState.classList.remove('hidden');
                emptyState.classList.add('flex');
            } else if (state === 'loading') {
                loadingState.classList.remove('hidden');
                loadingState.classList.add('flex');
            } else if (state === 'content') {
                contentArea.classList.remove('hidden');
            }
        }

        async function loadModels() {
            try {
                const response = await fetch('/v1/models');
                const data = await response.json();
                const modelSelect = document.getElementById('model');
                modelSelect.innerHTML = '';

                if (data.data && Array.isArray(data.data) && data.data.length > 0) {
                    data.data.forEach(model => {
                        const option = document.createElement('option');
                        option.value = model.id;
                        option.textContent = model.id;
                        if (model.id === '%s') {
                            option.selected = true;
                        }
                        modelSelect.appendChild(option);
                    });
                } else {
                    const option = document.createElement('option');
                    option.value = '%s';
                    option.textContent = '%s (默认)';
                    option.selected = true;
                    modelSelect.appendChild(option);
                }
            } catch (error) {
                console.error('Failed to load models:', error);
            }
        }

        loadModels();

        clearBtn.addEventListener('click', () => {
            showState('empty');
            contentArea.innerHTML = '';
        });

        sendBtn.addEventListener('click', async () => {
            const apiKey = document.getElementById('apiKey').value;
            const model = document.getElementById('model').value;
            const stream = document.getElementById('stream').checked;
            const messageText = document.getElementById('message').value.trim();

            if (!messageText) {
                alert('请输入消息内容');
                return;
            }

            const messages = [{ role: 'user', content: messageText }];
            sendBtn.disabled = true;
            sendBtn.textContent = '⏳ 请求中...';
            showState('loading');

            try {
                const response = await fetch('/v1/chat/completions', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + apiKey
                    },
                    body: JSON.stringify({ model, messages, stream })
                });

                if (!response.ok) {
                    throw new Error('HTTP error! status: ' + response.status);
                }

                showState('content');
                contentArea.innerHTML = '';

                if (stream) {
                    const reader = response.body.getReader();
                    const decoder = new TextDecoder();

                    while (true) {
                        const {value, done} = await reader.read();
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
                                    if (content) {
                                        contentArea.innerHTML += content;
                                        responseDiv.scrollTop = responseDiv.scrollHeight;
                                    }
                                } catch (e) {}
                            }
                        }
                    }
                } else {
                    const data = await response.json();
                    const content = data.choices[0]?.message?.content || '无响应';
                    contentArea.textContent = content;
                }

                sendBtn.textContent = '🚀 发送请求';
                sendBtn.disabled = false;
            } catch (error) {
                showState('content');
                contentArea.innerHTML = '<div class="text-red-600">请求失败: ' + error.message + '</div>';
                sendBtn.textContent = '🚀 发送请求';
                sendBtn.disabled = false;
            }
        });
    </script>
</body>
</html>`, MODEL_NAME, MODEL_NAME, MODEL_NAME)
}

// 生成 API 文档 HTML (从 deno 版本移植)
func getAPIDocsHTML() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50">
    <nav class="bg-white shadow-sm border-b">
        <div class="container mx-auto px-4 py-4">
            <div class="flex items-center justify-between">
                <a href="/" class="flex items-center space-x-2 text-purple-600 hover:text-purple-700 transition">
                    <span class="text-2xl">🦕</span>
                    <span class="text-xl font-bold">ZtoApi</span>
                </a>
                <div class="flex space-x-4">
                    <a href="/" class="text-gray-600 hover:text-purple-600 transition">首页</a>
                    <a href="/docs" class="text-purple-600 font-semibold">文档</a>
                    <a href="/playground" class="text-gray-600 hover:text-purple-600 transition">Playground</a>
                    <a href="/deploy" class="text-gray-600 hover:text-purple-600 transition">部署</a>
                    <a href="/dashboard" class="text-gray-600 hover:text-purple-600 transition">Dashboard</a>
                </div>
            </div>
        </div>
    </nav>

    <div class="container mx-auto px-4 py-8 max-w-5xl">
        <div class="text-center mb-12">
            <h1 class="text-4xl font-bold text-gray-900 mb-3">📖 API Documentation</h1>
            <p class="text-gray-600">OpenAI 兼容的 API 接口文档</p>
        </div>

        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-4">概述</h2>
            <p class="text-gray-700 mb-4">ZtoApi 是一个为 Z.ai GLM-4.6 模型提供 OpenAI 兼容 API 接口的代理服务器。</p>
            <div class="bg-purple-50 border border-purple-200 rounded-lg p-4">
                <p class="text-sm text-gray-600 mb-2">基础 URL</p>
                <code class="text-purple-700 font-mono text-lg">http://localhost%s/v1</code>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-4">🔐 身份验证</h2>
            <p class="text-gray-700 mb-4">所有 API 请求都需要在请求头中包含 Bearer Token：</p>
            <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                <code class="text-green-400 font-mono text-sm">Authorization: Bearer $你设置的 DEFAULT_KEY</code>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-6">🔌 API 端点</h2>

            <div class="mb-8">
                <div class="flex items-center space-x-3 mb-3">
                    <span class="bg-green-100 text-green-700 px-3 py-1 rounded-lg font-semibold text-sm">GET</span>
                    <code class="text-lg font-mono text-gray-800">/v1/models</code>
                </div>
                <p class="text-gray-700 mb-3">获取可用模型列表</p>
                <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                    <pre class="text-green-400 font-mono text-sm">curl http://localhost%s/v1/models \\
  -H "Authorization: Bearer $你设置的 DEFAULT_KEY"</pre>
                </div>
            </div>

            <div>
                <div class="flex items-center space-x-3 mb-3">
                    <span class="bg-blue-100 text-blue-700 px-3 py-1 rounded-lg font-semibold text-sm">POST</span>
                    <code class="text-lg font-mono text-gray-800">/v1/chat/completions</code>
                </div>
                <p class="text-gray-700 mb-4">创建聊天完成（支持流式和非流式）</p>

                <div class="bg-gray-50 rounded-lg p-4 mb-4">
                    <h4 class="font-semibold text-gray-900 mb-3">请求参数</h4>
                    <div class="space-y-2 text-sm">
                        <div class="flex items-start">
                            <code class="bg-white px-2 py-1 rounded mr-3 text-purple-600 font-mono">model</code>
                            <span class="text-gray-600">string, 必需 - 模型名称 (如 "%s")</span>
                        </div>
                        <div class="flex items-start">
                            <code class="bg-white px-2 py-1 rounded mr-3 text-purple-600 font-mono">messages</code>
                            <span class="text-gray-600">array, 必需 - 消息列表</span>
                        </div>
                        <div class="flex items-start">
                            <code class="bg-white px-2 py-1 rounded mr-3 text-purple-600 font-mono">stream</code>
                            <span class="text-gray-600">boolean, 可选 - 是否流式响应（默认: true）</span>
                        </div>
                    </div>
                </div>

                <h4 class="font-semibold text-gray-900 mb-3">请求示例</h4>
                <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                    <pre class="text-green-400 font-mono text-sm">curl -X POST http://localhost%s/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer $你设置的 DEFAULT_KEY" \\
  -d '{
    "model": "%s",
    "messages": [
      {"role": "user", "content": "你好"}
    ],
    "stream": false
  }'</pre>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-4">🐍 Python 示例</h2>
            <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                <pre class="text-green-400 font-mono text-sm">import openai

client = openai.OpenAI(
    api_key="$你设置的 DEFAULT_KEY",
    base_url="http://localhost%s/v1"
)

response = client.chat.completions.create(
    model="%s",
    messages=[{"role": "user", "content": "你好"}]
)

print(response.choices[0].message.content)</pre>
            </div>
        </div>

        <div class="text-center">
            <a href="/" class="inline-block bg-purple-600 hover:bg-purple-700 text-white font-semibold px-6 py-3 rounded-lg transition">
                返回首页
            </a>
        </div>
    </div>
</body>
</html>`, PORT, PORT, MODEL_NAME, PORT, MODEL_NAME, PORT, MODEL_NAME)
}

// 生成 Dashboard HTML (从 deno 版本完整移植)
func getDashboardHTMLNew() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body class="bg-gray-50">
    <nav class="bg-white shadow-sm border-b">
        <div class="container mx-auto px-4 py-4">
            <div class="flex items-center justify-between">
                <a href="/" class="flex items-center space-x-2 text-purple-600 hover:text-purple-700 transition">
                    <span class="text-2xl">🦕</span>
                    <span class="text-xl font-bold">ZtoApi</span>
                </a>
                <div class="flex space-x-4">
                    <a href="/" class="text-gray-600 hover:text-purple-600 transition">首页</a>
                    <a href="/docs" class="text-gray-600 hover:text-purple-600 transition">文档</a>
                    <a href="/playground" class="text-gray-600 hover:text-purple-600 transition">Playground</a>
                    <a href="/deploy" class="text-gray-600 hover:text-purple-600 transition">部署</a>
                    <a href="/dashboard" class="text-purple-600 font-semibold">Dashboard</a>
                </div>
            </div>
        </div>
    </nav>

    <div class="container mx-auto px-4 py-8 max-w-7xl">
        <div class="text-center mb-8">
            <h1 class="text-4xl font-bold text-gray-900 mb-3">📊 Dashboard</h1>
            <p class="text-gray-600">实时监控 API 请求和性能统计</p>
        </div>

        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-600 text-sm mb-1">总请求数</p>
                        <p class="text-3xl font-bold text-gray-900" id="total">0</p>
                    </div>
                    <div class="bg-purple-100 p-3 rounded-lg">
                        <span class="text-3xl">📈</span>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-600 text-sm mb-1">成功请求</p>
                        <p class="text-3xl font-bold text-green-600" id="success">0</p>
                    </div>
                    <div class="bg-green-100 p-3 rounded-lg">
                        <span class="text-3xl">✅</span>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-600 text-sm mb-1">失败请求</p>
                        <p class="text-3xl font-bold text-red-600" id="failed">0</p>
                    </div>
                    <div class="bg-red-100 p-3 rounded-lg">
                        <span class="text-3xl">❌</span>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-600 text-sm mb-1">平均响应时间</p>
                        <p class="text-3xl font-bold text-blue-600" id="avgtime">0ms</p>
                    </div>
                    <div class="bg-blue-100 p-3 rounded-lg">
                        <span class="text-3xl">⚡</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Requests Table -->
        <div class="bg-white rounded-xl shadow-sm border p-6">
            <div class="flex items-center justify-between mb-4">
                <h2 class="text-xl font-bold text-gray-900">🔔 实时请求</h2>
                <span class="text-sm text-gray-500">自动刷新（每5秒）</span>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b">
                            <th class="text-left py-3 px-4 text-gray-700 font-semibold">时间</th>
                            <th class="text-left py-3 px-4 text-gray-700 font-semibold">方法</th>
                            <th class="text-left py-3 px-4 text-gray-700 font-semibold">路径</th>
                            <th class="text-left py-3 px-4 text-gray-700 font-semibold">状态</th>
                            <th class="text-left py-3 px-4 text-gray-700 font-semibold">耗时</th>
                        </tr>
                    </thead>
                    <tbody id="requests" class="divide-y"></tbody>
                </table>
            </div>
            <div id="empty" class="text-center py-8 text-gray-500 hidden">
                暂无请求记录
            </div>
        </div>
    </div>

    <script>
        async function update() {
            try {
                const statsRes = await fetch('/dashboard/stats');
                const stats = await statsRes.json();

                document.getElementById('total').textContent = stats.totalRequests || 0;
                document.getElementById('success').textContent = stats.successfulRequests || 0;
                document.getElementById('failed').textContent = stats.failedRequests || 0;
                document.getElementById('avgtime').textContent = Math.round(stats.averageResponseTime) + 'ms';

                const reqsRes = await fetch('/dashboard/requests');
                const requests = await reqsRes.json();
                const tbody = document.getElementById('requests');
                const empty = document.getElementById('empty');

                tbody.innerHTML = '';

                if (!requests || requests.length === 0) {
                    empty.classList.remove('hidden');
                } else {
                    empty.classList.add('hidden');
                    requests.forEach(r => {
                        const row = tbody.insertRow();
                        const time = new Date(r.timestamp).toLocaleTimeString();
                        const statusClass = r.status >= 200 && r.status < 300 ? 'text-green-600 bg-green-50' : 'text-red-600 bg-red-50';

                        row.innerHTML = '<td class="py-3 px-4 text-gray-700">' + time + '</td>' +
                            '<td class="py-3 px-4"><span class="bg-blue-100 text-blue-700 px-2 py-1 rounded text-sm font-mono">' + r.method + '</span></td>' +
                            '<td class="py-3 px-4 font-mono text-sm text-gray-600">' + r.path + '</td>' +
                            '<td class="py-3 px-4"><span class="' + statusClass + ' px-2 py-1 rounded font-semibold text-sm">' + r.status + '</span></td>' +
                            '<td class="py-3 px-4 text-gray-700">' + r.duration + 'ms</td>';
                    });
                }
            } catch (e) {
                console.error('Update error:', e);
            }
        }

        update();
        setInterval(update, 5000);
    </script>
</body>
</html>`
}

// 生成 Admin 登录页面 HTML
func getAdminLoginHTML() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>管理员登录 - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 min-h-screen flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        <div class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-2">🔐 管理员登录</h1>
            <p class="text-gray-600">ZtoApi 管理系统</p>
        </div>

        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
            <p class="text-sm text-blue-800 text-center">
                📌 该功能需要配置 register 模块<br>
                请访问 <a href="/register" class="underline font-semibold">/register</a> 进行账号管理
            </p>
        </div>

        <div class="flex justify-center">
            <a href="/" class="inline-block bg-purple-600 hover:bg-purple-700 text-white font-semibold px-6 py-3 rounded-lg transition">
                返回首页
            </a>
        </div>
    </div>
</body>
</html>`
}

// 生成部署指南 HTML
func getDeployHTML() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>部署指南 - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50">
    <nav class="bg-white shadow-sm border-b">
        <div class="container mx-auto px-4 py-4">
            <div class="flex items-center justify-between">
                <a href="/" class="flex items-center space-x-2 text-purple-600 hover:text-purple-700 transition">
                    <span class="text-2xl">🦕</span>
                    <span class="text-xl font-bold">ZtoApi</span>
                </a>
                <div class="flex space-x-4">
                    <a href="/" class="text-gray-600 hover:text-purple-600 transition">首页</a>
                    <a href="/docs" class="text-gray-600 hover:text-purple-600 transition">文档</a>
                    <a href="/playground" class="text-gray-600 hover:text-purple-600 transition">Playground</a>
                    <a href="/deploy" class="text-purple-600 font-semibold">部署</a>
                    <a href="/dashboard" class="text-gray-600 hover:text-purple-600 transition">Dashboard</a>
                </div>
            </div>
        </div>
    </nav>

    <div class="container mx-auto px-4 py-8 max-w-5xl">
        <div class="text-center mb-12">
            <h1 class="text-4xl font-bold text-gray-900 mb-3">🚀 部署指南</h1>
            <p class="text-gray-600">快速部署 ZtoApi</p>
        </div>

        <!-- Docker 部署 -->
        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
                <span class="mr-3">🐋</span> Docker 部署
            </h2>
            <div class="space-y-4">
                <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                    <pre class="text-green-400 font-mono text-sm">docker run -d -p 9090:9090 \\
  -e DEFAULT_KEY=your-api-key \\
  -e ZAI_TOKEN=your-zai-token \\
  -e MODEL_NAME=GLM-4.6 \\
  hulisang/ztoapi:latest</pre>
                </div>
            </div>
        </div>

        <!-- 源码部署 -->
        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
                <span class="mr-3">📦</span> 源码部署
            </h2>
            <div class="space-y-4">
                <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                    <pre class="text-green-400 font-mono text-sm">git clone https://github.com/hulisang/ZtoApi.git
cd ZtoApi

# 配置环境变量
export DEFAULT_KEY=your-api-key
export ZAI_TOKEN=your-zai-token
export MODEL_NAME=GLM-4.6

# 编译运行
go build -o ztoapi main.go
./ztoapi</pre>
                </div>
            </div>
        </div>

        <!-- 环境变量 -->
        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
                <span class="mr-3">🔐</span> 环境变量配置
            </h2>
            <div class="space-y-3 text-sm">
                <div class="bg-gray-50 rounded p-3">
                    <code class="text-purple-600 font-mono">DEFAULT_KEY</code>
                    <span class="text-gray-600 ml-2">- API 密钥（必需）</span>
                </div>
                <div class="bg-gray-50 rounded p-3">
                    <code class="text-purple-600 font-mono">ZAI_TOKEN</code>
                    <span class="text-gray-600 ml-2">- Z.ai Token（可选，不设置将使用匿名token）</span>
                </div>
                <div class="bg-gray-50 rounded p-3">
                    <code class="text-purple-600 font-mono">MODEL_NAME</code>
                    <span class="text-gray-600 ml-2">- 模型名称（默认：GLM-4.6）</span>
                </div>
                <div class="bg-gray-50 rounded p-3">
                    <code class="text-purple-600 font-mono">PORT</code>
                    <span class="text-gray-600 ml-2">- 服务端口（默认：9090）</span>
                </div>
                <div class="bg-gray-50 rounded p-3">
                    <code class="text-purple-600 font-mono">DEBUG_MODE</code>
                    <span class="text-gray-600 ml-2">- 调试模式（默认：true）</span>
                </div>
            </div>
        </div>

        <div class="flex justify-center space-x-4">
            <a href="/" class="inline-block bg-purple-600 hover:bg-purple-700 text-white font-semibold px-8 py-3 rounded-lg transition">
                返回首页
            </a>
        </div>
    </div>
</body>
</html>`
}
