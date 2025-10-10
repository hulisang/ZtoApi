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

// 生成 Playground HTML (完整版，包含所有高级功能)
func getPlaygroundHTML() string {
	enableThinkingChecked := ""
	if ENABLE_THINKING {
		enableThinkingChecked = "checked"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Playground - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <!-- Marked.js for Markdown parsing -->
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <!-- Highlight.js for code syntax highlighting -->
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
        #response blockquote {
            border-left: 4px solid #ddd;
            padding-left: 1em;
            margin: 1em 0;
            color: #666;
        }
        #response a {
            color: #3b82f6;
            text-decoration: underline;
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

                <!-- API Key -->
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">API Key</label>
                    <input type="text" id="apiKey" value=""
                           class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                           placeholder="请输入你的 API Key (例如: sk-your-key)">
                </div>

                <!-- ZAI Token (Optional) -->
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">ZAI Token (可选)</label>
                    <input type="text" id="zaiToken" value=""
                           class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent font-mono text-sm"
                           placeholder="留空则使用匿名 token">
                    <p class="text-xs text-gray-500 mt-1">自定义 Z.ai 上游 token（高级选项）</p>
                </div>

                <!-- Model Selection -->
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">模型</label>
                    <select id="model" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500">
                        <option value="">加载中...</option>
                    </select>
                </div>

                <!-- Parameters Row -->
                <div class="grid grid-cols-2 gap-3 mb-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Temperature</label>
                        <input type="number" id="temperature" min="0" max="2" step="0.1" value="0.7"
                               class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
                               placeholder="0.7">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Max Tokens</label>
                        <input type="number" id="maxTokens" min="1" max="8192" step="1" value="2048"
                               class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
                               placeholder="2048">
                    </div>
                </div>

                <!-- Stream -->
                <div class="mb-4">
                    <label class="flex items-center">
                        <input type="checkbox" id="stream" checked class="mr-2">
                        <span class="text-sm font-medium text-gray-700">启用流式响应</span>
                    </label>
                </div>

                <!-- Enable Thinking -->
                <div class="mb-4">
                    <label class="flex items-center">
                        <input type="checkbox" id="enableThinking" %s class="mr-2">
                        <span class="text-sm font-medium text-gray-700">启用思维链</span>
                    </label>
                    <p class="text-xs text-gray-500 mt-1">启用后将显示 AI 的思考过程</p>
                </div>

                <!-- System Message -->
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">System (可选)</label>
                    <textarea id="system" rows="3"
                              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 text-sm"
                              placeholder="你是一个乐于助人的AI助手..."></textarea>
                    <p class="text-xs text-gray-500 mt-1">系统提示词，用于设定角色和行为</p>
                </div>

                <!-- User Message -->
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Message</label>
                    <textarea id="message" rows="6"
                              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 text-sm"
                              placeholder="输入你的问题...">你好，请介绍一下你自己</textarea>
                    <p class="text-xs text-gray-500 mt-1">用户消息内容</p>
                </div>

                <!-- Send Button -->
                <button id="sendBtn"
                        class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-4 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed">
                    🚀 发送请求
                </button>

                <!-- Clear Button -->
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

                <!-- Status -->
                <div id="status" class="mb-4 p-3 bg-gray-100 rounded-lg hidden">
                    <span class="font-mono text-sm"></span>
                </div>

                <!-- Response -->
                <div class="mb-4">
                    <div class="flex items-center justify-between mb-2">
                        <label class="block text-sm font-medium text-gray-700">响应内容</label>
                        <button id="copyBtn" class="text-xs text-purple-600 hover:text-purple-700 hidden">📋 复制</button>
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
                            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mb-4"></div>
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

                <!-- Stats -->
                <div id="stats" class="grid grid-cols-2 gap-3 hidden">
                    <div class="bg-purple-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-600">耗时</p>
                        <p id="duration" class="text-lg font-bold text-purple-600">-</p>
                    </div>
                    <div class="bg-green-50 p-3 rounded-lg">
                        <p class="text-xs text-gray-600">状态</p>
                        <p id="statusCode" class="text-lg font-bold text-green-600">-</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Request/Response Examples -->
        <div class="mt-8 bg-white rounded-xl shadow-sm border p-6">
            <h2 class="text-xl font-bold text-gray-900 mb-4 flex items-center">
                <span class="text-2xl mr-2">💡</span> 示例
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <button class="text-left p-4 border border-gray-200 rounded-lg hover:border-purple-500 hover:bg-purple-50 transition example-btn"
                        data-example="simple">
                    <p class="font-semibold text-gray-900">简单对话</p>
                    <p class="text-sm text-gray-600">单轮对话示例</p>
                </button>
                <button class="text-left p-4 border border-gray-200 rounded-lg hover:border-purple-500 hover:bg-purple-50 transition example-btn"
                        data-example="thinking">
                    <p class="font-semibold text-gray-900">思维链示例</p>
                    <p class="text-sm text-gray-600">展示 AI 思考过程</p>
                </button>
                <button class="text-left p-4 border border-gray-200 rounded-lg hover:border-purple-500 hover:bg-purple-50 transition example-btn"
                        data-example="code">
                    <p class="font-semibold text-gray-900">代码生成</p>
                    <p class="text-sm text-gray-600">生成代码示例</p>
                </button>
                <button class="text-left p-4 border border-gray-200 rounded-lg hover:border-purple-500 hover:bg-purple-50 transition example-btn"
                        data-example="creative">
                    <p class="font-semibold text-gray-900">创意写作</p>
                    <p class="text-sm text-gray-600">高温度创意输出</p>
                </button>
            </div>
        </div>
    </div>

    <footer class="bg-white border-t mt-12 py-6">
        <div class="container mx-auto px-4 text-center text-gray-500 text-sm">
            <p>Powered by <span class="font-semibold">Go</span> | <a href="/" class="text-purple-600 hover:underline">返回首页</a> | <a href="https://github.com/hulisang/ZtoApi" target="_blank" rel="noopener noreferrer" class="text-purple-600 hover:underline">⭐ GitHub</a></p>
        </div>
    </footer>

    <script>
        const examples = {
            simple: {
                model: '%s',
                system: '',
                message: '你好，请介绍一下你自己',
                enableThinking: false,
                temperature: 0.7
            },
            thinking: {
                model: '%s',
                system: '你是一个专业的数学老师，擅长用清晰的思路解决问题。',
                message: '一个正方形的边长是5厘米，求它的面积和周长。',
                enableThinking: true,
                temperature: 0.7
            },
            code: {
                model: '%s',
                system: '你是一个专业的编程助手，提供清晰、高效的代码示例。',
                message: '用 JavaScript 写一个函数，判断一个字符串是否为回文',
                enableThinking: false,
                temperature: 0.7
            },
            creative: {
                model: '%s',
                system: '你是一个富有创意的作家，擅长创作引人入胜的故事。',
                message: '写一个关于未来城市的短故事（100字以内）',
                enableThinking: false,
                temperature: 1.2
            }
        };

        const sendBtn = document.getElementById('sendBtn');
        const clearBtn = document.getElementById('clearBtn');
        const copyBtn = document.getElementById('copyBtn');
        const responseDiv = document.getElementById('response');
        const emptyState = document.getElementById('emptyState');
        const loadingState = document.getElementById('loadingState');
        const errorState = document.getElementById('errorState');
        const contentArea = document.getElementById('contentArea');
        const statsDiv = document.getElementById('stats');

        let responseContent = '';
        let startTime = 0;

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
                copyBtn.classList.add('hidden');
                statsDiv.classList.add('hidden');
            } else if (state === 'loading') {
                loadingState.classList.remove('hidden');
                loadingState.classList.add('flex');
                copyBtn.classList.add('hidden');
                statsDiv.classList.add('hidden');
            } else if (state === 'error') {
                errorState.classList.remove('hidden');
                errorState.classList.add('flex');
                copyBtn.classList.add('hidden');
                statsDiv.classList.add('hidden');
            } else if (state === 'content') {
                contentArea.classList.remove('hidden');
                copyBtn.classList.remove('hidden');
                statsDiv.classList.remove('hidden');
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
            responseContent = '';
        });

        copyBtn.addEventListener('click', () => {
            navigator.clipboard.writeText(responseContent).then(() => {
                const originalText = copyBtn.textContent;
                copyBtn.textContent = '✅ 已复制';
                setTimeout(() => {
                    copyBtn.textContent = originalText;
                }, 2000);
            });
        });

        // Example buttons
        document.querySelectorAll('.example-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                const exampleName = btn.dataset.example;
                const example = examples[exampleName];
                if (example) {
                    document.getElementById('model').value = example.model;
                    document.getElementById('system').value = example.system || '';
                    document.getElementById('message').value = example.message;
                    document.getElementById('enableThinking').checked = example.enableThinking;
                    document.getElementById('temperature').value = example.temperature || 0.7;
                }
            });
        });

        sendBtn.addEventListener('click', async () => {
            const apiKey = document.getElementById('apiKey').value;
            const zaiToken = document.getElementById('zaiToken').value.trim();
            const model = document.getElementById('model').value;
            const stream = document.getElementById('stream').checked;
            const enableThinking = document.getElementById('enableThinking').checked;
            const system = document.getElementById('system').value.trim();
            const messageText = document.getElementById('message').value.trim();
            const temperature = parseFloat(document.getElementById('temperature').value) || 0.7;
            const maxTokens = parseInt(document.getElementById('maxTokens').value) || 2048;

            if (!messageText) {
                alert('请输入消息内容');
                return;
            }

            const messages = [];
            if (system) {
                messages.push({ role: 'system', content: system });
            }
            messages.push({ role: 'user', content: messageText });

            const requestBody = {
                model,
                messages,
                stream,
                temperature,
                max_tokens: maxTokens
            };

            // 添加ZAI Token到headers或body (取决于API设计)
            const headers = {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + apiKey
            };
            if (zaiToken) {
                headers['X-ZAI-Token'] = zaiToken;
            }

            sendBtn.disabled = true;
            sendBtn.textContent = '⏳ 请求中...';
            showState('loading');
            responseContent = '';
            startTime = Date.now();

            try {
                const response = await fetch('/v1/chat/completions', {
                    method: 'POST',
                    headers: headers,
                    body: JSON.stringify(requestBody)
                });

                const duration = Date.now() - startTime;
                document.getElementById('duration').textContent = duration + 'ms';
                document.getElementById('statusCode').textContent = response.status;

                if (!response.ok) {
                    throw new Error('HTTP ' + response.status + ': ' + response.statusText);
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
                                        responseContent += content;
                                        // 使用marked渲染Markdown
                                        contentArea.innerHTML = marked.parse(responseContent);
                                        // 高亮代码块
                                        contentArea.querySelectorAll('pre code').forEach((block) => {
                                            hljs.highlightElement(block);
                                        });
                                        responseDiv.scrollTop = responseDiv.scrollHeight;
                                    }
                                } catch (e) {}
                            }
                        }
                    }
                } else {
                    const data = await response.json();
                    const content = data.choices[0]?.message?.content || '无响应';
                    responseContent = content;
                    contentArea.innerHTML = marked.parse(content);
                    contentArea.querySelectorAll('pre code').forEach((block) => {
                        hljs.highlightElement(block);
                    });
                }

                sendBtn.textContent = '🚀 发送请求';
                sendBtn.disabled = false;
            } catch (error) {
                showState('error');
                document.getElementById('errorMessage').textContent = error.message;
                sendBtn.textContent = '🚀 发送请求';
                sendBtn.disabled = false;
            }
        });
    </script>
</body>
</html>`, enableThinkingChecked, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME, MODEL_NAME)
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
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
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

            <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-600 text-sm mb-1">首页访问</p>
                        <p class="text-3xl font-bold text-indigo-600" id="homeviews">0</p>
                    </div>
                    <div class="bg-indigo-100 p-3 rounded-lg">
                        <span class="text-3xl">🏠</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Detailed Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
            <!-- API Stats -->
            <div class="bg-white rounded-xl shadow-sm border p-6">
                <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center">
                    <span class="text-2xl mr-2">🎯</span> API 统计
                </h3>
                <div class="space-y-3">
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">Chat Completions</span>
                        <span class="font-bold text-purple-600" id="api-calls">0</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">Models 查询</span>
                        <span class="font-bold text-purple-600" id="models-calls">0</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">流式请求</span>
                        <span class="font-bold text-blue-600" id="streaming">0</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">非流式请求</span>
                        <span class="font-bold text-blue-600" id="non-streaming">0</span>
                    </div>
                </div>
            </div>

            <!-- Performance Stats -->
            <div class="bg-white rounded-xl shadow-sm border p-6">
                <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center">
                    <span class="text-2xl mr-2">⚡</span> 性能指标
                </h3>
                <div class="space-y-3">
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">平均响应</span>
                        <span class="font-bold text-blue-600" id="avg-time-detail">0ms</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">最快响应</span>
                        <span class="font-bold text-green-600" id="fastest">-</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">最慢响应</span>
                        <span class="font-bold text-orange-600" id="slowest">0ms</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">成功率</span>
                        <span class="font-bold text-green-600" id="success-rate">0%</span>
                    </div>
                </div>
            </div>

            <!-- System Info -->
            <div class="bg-white rounded-xl shadow-sm border p-6">
                <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center">
                    <span class="text-2xl mr-2">📊</span> 系统信息
                </h3>
                <div class="space-y-3">
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">运行时长</span>
                        <span class="font-bold text-indigo-600" id="uptime">0</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">Token 使用</span>
                        <span class="font-bold text-indigo-600" id="tokens">0</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">最后请求</span>
                        <span class="font-bold text-gray-600 text-xs" id="last-request">-</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600 text-sm">首页访问</span>
                        <span class="font-bold text-indigo-600" id="home-visits">0</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Top Models Card -->
        <div class="bg-white rounded-xl shadow-sm border p-6 mb-8">
            <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center">
                <span class="text-2xl mr-2">🏆</span> 热门模型 Top 3
            </h3>
            <div id="top-models" class="space-y-3">
                <p class="text-gray-500 text-sm">暂无数据</p>
            </div>
        </div>

        <!-- Chart -->
        <div class="bg-white rounded-xl shadow-sm border p-6 mb-8">
            <div class="flex items-center justify-between mb-4">
                <h2 class="text-xl font-bold text-gray-900">📉 请求趋势分析</h2>
                <div class="flex gap-2">
                    <button id="view-hourly" class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition text-sm font-semibold">按小时</button>
                    <button id="view-daily" class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition text-sm font-semibold">按天</button>
                </div>
            </div>

            <!-- Info banner -->
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-3">
                <p class="text-sm text-blue-800">
                    💡 <strong>提示：</strong>此图表显示基于 SQLite 持久化存储的历史数据。数据会在每次 API 请求后自动保存，并在服务器上永久保留。
                </p>
            </div>

            <div class="mb-3 flex items-center gap-4">
                <div class="flex items-center gap-2">
                    <span class="text-sm text-gray-600">时间范围:</span>
                    <select id="time-range" class="px-3 py-1 border rounded-lg text-sm">
                        <option value="12">最近12个</option>
                        <option value="24" selected>最近24个</option>
                        <option value="48">最近48个</option>
                        <option value="72">最近72个</option>
                    </select>
                </div>
                <div class="text-sm text-gray-500" id="chart-subtitle">显示最近24小时的数据</div>
            </div>
            <canvas id="chart" height="80"></canvas>
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
                            <th class="text-left py-3 px-4 text-gray-700 font-semibold">模型</th>
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
            <!-- Pagination -->
            <div id="pagination" class="mt-4 flex items-center justify-between">
                <div class="flex items-center gap-4">
                    <div class="text-sm text-gray-600">
                        共 <span id="total-requests">0</span> 条记录，第 <span id="current-page">1</span> / <span id="total-pages">1</span> 页
                    </div>
                    <div class="flex items-center gap-2">
                        <span class="text-sm text-gray-600">每页:</span>
                        <select id="page-size" class="px-2 py-1 border rounded text-sm">
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="20" selected>20</option>
                            <option value="50">50</option>
                            <option value="100">100</option>
                        </select>
                    </div>
                </div>
                <div class="flex gap-2">
                    <button id="prev-page" class="px-3 py-1 bg-gray-200 hover:bg-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed">上一页</button>
                    <button id="next-page" class="px-3 py-1 bg-gray-200 hover:bg-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed">下一页</button>
                </div>
            </div>
        </div>
    </div>

    <script>
        let chart = null;
        const chartData = { labels: [], data: [] };
        let currentPage = 1;
        let pageSize = 20;
        let chartViewMode = 'hourly'; // 'hourly' or 'daily'
        let chartTimeRange = 24; // hours or days

        async function update() {
            try {
                const statsRes = await fetch('/dashboard/stats');
                const stats = await statsRes.json();

                // Top cards
                document.getElementById('total').textContent = stats.totalRequests;
                document.getElementById('success').textContent = stats.successfulRequests;
                document.getElementById('failed').textContent = stats.failedRequests;
                document.getElementById('avgtime').textContent = Math.round(stats.averageResponseTime) + 'ms';
                document.getElementById('homeviews').textContent = stats.homePageViews;

                // API Stats
                document.getElementById('api-calls').textContent = stats.apiCallsCount || 0;
                document.getElementById('models-calls').textContent = stats.modelsCallsCount || 0;
                document.getElementById('streaming').textContent = stats.streamingRequests || 0;
                document.getElementById('non-streaming').textContent = stats.nonStreamingRequests || 0;

                // Performance Stats
                document.getElementById('avg-time-detail').textContent = Math.round(stats.averageResponseTime) + 'ms';
                document.getElementById('fastest').textContent = stats.fastestResponse === -1 ? '-' : Math.round(stats.fastestResponse) + 'ms';
                document.getElementById('slowest').textContent = stats.slowestResponse === 0 ? '-' : Math.round(stats.slowestResponse) + 'ms';
                const successRate = stats.totalRequests > 0 ? ((stats.successfulRequests / stats.totalRequests) * 100).toFixed(1) : '0';
                document.getElementById('success-rate').textContent = successRate + '%';

                // System Info
                const uptime = Date.now() - new Date(stats.startTime).getTime();
                const hours = Math.floor(uptime / 3600000);
                const minutes = Math.floor((uptime % 3600000) / 60000);
                document.getElementById('uptime').textContent = hours + 'h ' + minutes + 'm';
                document.getElementById('tokens').textContent = (stats.totalTokensUsed || 0).toLocaleString();
                document.getElementById('last-request').textContent = stats.lastRequestTime ? new Date(stats.lastRequestTime).toLocaleTimeString() : '-';
                document.getElementById('home-visits').textContent = stats.homePageViews;

                // Top Models
                const topModelsDiv = document.getElementById('top-models');
                if (stats.topModels && stats.topModels.length > 0) {
                    topModelsDiv.innerHTML = stats.topModels.map((m, i) => ` + "`" + `
                        <div class="flex items-center justify-between">
                            <div class="flex items-center gap-2">
                                <span class="text-lg">${i === 0 ? '🥇' : i === 1 ? '🥈' : '🥉'}</span>
                                <span class="font-mono text-sm text-gray-700">${m.model}</span>
                            </div>
                            <span class="font-bold text-purple-600">${m.count}</span>
                        </div>
                    ` + "`" + `).join('');
                } else {
                    topModelsDiv.innerHTML = '<p class="text-gray-500 text-sm">暂无数据</p>';
                }

                // Fetch paginated requests
                const reqsRes = await fetch(` + "`" + `/dashboard/requests?page=${currentPage}&pageSize=${pageSize}` + "`" + `);
                const data = await reqsRes.json();
                const tbody = document.getElementById('requests');
                const empty = document.getElementById('empty');

                tbody.innerHTML = '';

                if (data.requests.length === 0) {
                    empty.classList.remove('hidden');
                } else {
                    empty.classList.add('hidden');
                    data.requests.forEach(r => {
                        const row = tbody.insertRow();
                        const time = new Date(r.timestamp).toLocaleTimeString();
                        const statusClass = r.status >= 200 && r.status < 300 ? 'text-green-600 bg-green-50' : 'text-red-600 bg-red-50';
                        const modelDisplay = r.model ? r.model : '-';

                        row.innerHTML = ` + "`" + `
                            <td class="py-3 px-4 text-gray-700">${time}</td>
                            <td class="py-3 px-4"><span class="bg-blue-100 text-blue-700 px-2 py-1 rounded text-sm font-mono">${r.method}</span></td>
                            <td class="py-3 px-4 font-mono text-sm text-gray-600">${r.path}</td>
                            <td class="py-3 px-4 font-mono text-xs text-gray-600">${modelDisplay}</td>
                            <td class="py-3 px-4"><span class="${statusClass} px-2 py-1 rounded font-semibold text-sm">${r.status}</span></td>
                            <td class="py-3 px-4 text-gray-700">${r.duration}ms</td>
                        ` + "`" + `;
                    });

                    // Update pagination info
                    document.getElementById('total-requests').textContent = data.total;
                    document.getElementById('current-page').textContent = data.page;
                    document.getElementById('total-pages').textContent = data.totalPages;

                    // Enable/disable pagination buttons
                    document.getElementById('prev-page').disabled = data.page <= 1;
                    document.getElementById('next-page').disabled = data.page >= data.totalPages;
                }
            } catch (e) {
                console.error('Update error:', e);
            }
        }

        async function updateChartData() {
            try {
                let endpoint, labelKey, subtitle;

                if (chartViewMode === 'hourly') {
                    endpoint = ` + "`" + `/dashboard/hourly?hours=${chartTimeRange}` + "`" + `;
                    labelKey = 'hour';
                    subtitle = ` + "`显示最近${chartTimeRange}小时的数据`" + `;
                } else {
                    endpoint = ` + "`" + `/dashboard/daily?days=${chartTimeRange}` + "`" + `;
                    labelKey = 'date';
                    subtitle = ` + "`显示最近${chartTimeRange}天的数据`" + `;
                }

                const res = await fetch(endpoint);
                const data = await res.json();

                if (data && data.length > 0) {
                    chartData.labels = data.map(d => {
                        if (chartViewMode === 'hourly') {
                            // Format: 2025-09-30-14 -> 09-30 14:00
                            const parts = d[labelKey].split('-');
                            return ` + "`${parts[1]}-${parts[2]} ${parts[3]}:00`" + `;
                        } else {
                            // Format: 2025-09-30 -> 09-30
                            const parts = d[labelKey].split('-');
                            return ` + "`${parts[1]}-${parts[2]}`" + `;
                        }
                    });
                    chartData.data = data.map(d => Math.round(d.avgResponseTime));
                    subtitle += ` + "` (共${data.length}条记录)`" + `;
                } else {
                    chartData.labels = [];
                    chartData.data = [];
                    subtitle += ' - ⚠️ 暂无持久化数据，请发送API请求后稍等片刻';
                }

                document.getElementById('chart-subtitle').textContent = subtitle;
                updateChart();
            } catch (e) {
                console.error('Chart update error:', e);
                document.getElementById('chart-subtitle').textContent = '⚠️ 加载数据失败: ' + e.message;
            }
        }

        function updateChart() {
            const ctx = document.getElementById('chart').getContext('2d');

            if (chart) {
                chart.data.labels = chartData.labels;
                chart.data.datasets[0].data = chartData.data;
                chart.update();
            } else {
                chart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: chartData.labels,
                        datasets: [{
                            label: '响应时间 (ms)',
                            data: chartData.data,
                            borderColor: 'rgb(147, 51, 234)',
                            backgroundColor: 'rgba(147, 51, 234, 0.1)',
                            tension: 0.4,
                            fill: true
                        }]
                    },
                    options: {
                        responsive: true,
                        plugins: {
                            legend: { display: false },
                            tooltip: {
                                callbacks: {
                                    label: (ctx) => ` + "`响应时间: ${ctx.parsed.y}ms`" + `
                                }
                            }
                        },
                        scales: {
                            y: {
                                beginAtZero: true,
                                ticks: { callback: (val) => val + 'ms' }
                            }
                        }
                    }
                });
            }
        }

        // Pagination handlers
        document.getElementById('prev-page').addEventListener('click', () => {
            if (currentPage > 1) {
                currentPage--;
                update();
            }
        });

        document.getElementById('next-page').addEventListener('click', () => {
            currentPage++;
            update();
        });

        // Chart view mode handlers
        document.getElementById('view-hourly').addEventListener('click', () => {
            chartViewMode = 'hourly';
            document.getElementById('view-hourly').className = 'px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition text-sm font-semibold';
            document.getElementById('view-daily').className = 'px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition text-sm font-semibold';
            updateChartData();
        });

        document.getElementById('view-daily').addEventListener('click', () => {
            chartViewMode = 'daily';
            document.getElementById('view-daily').className = 'px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition text-sm font-semibold';
            document.getElementById('view-hourly').className = 'px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition text-sm font-semibold';
            updateChartData();
        });

        // Time range handler
        document.getElementById('time-range').addEventListener('change', (e) => {
            chartTimeRange = parseInt(e.target.value);
            updateChartData();
        });

        // Page size handler
        document.getElementById('page-size').addEventListener('change', (e) => {
            pageSize = parseInt(e.target.value);
            currentPage = 1; // Reset to first page
            update();
        });

        update();
        updateChartData();
        setInterval(update, 5000);
        setInterval(updateChartData, 60000); // Update chart every minute
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
            <p class="text-gray-600">ZtoApi 账号管理系统</p>
        </div>

        <form id="loginForm" class="space-y-6">
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
                <input type="text" id="username" required
                    class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">密码</label>
                <input type="password" id="password" required
                    class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
            </div>

            <div id="errorMsg" class="hidden text-red-500 text-sm text-center"></div>

            <button type="submit"
                class="w-full px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all">
                登录
            </button>
        </form>

        <div class="mt-6 text-center text-sm text-gray-500">
            <p>默认账号: admin / 123456</p>
        </div>
    </div>

    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const errorMsg = document.getElementById('errorMsg');

            errorMsg.classList.add('hidden');

            try {
                const response = await fetch('/admin/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                const result = await response.json();

                if (result.success) {
                    document.cookie = 'adminSessionId=' + result.sessionId + '; path=/; max-age=86400';
                    window.location.href = '/admin';
                } else {
                    errorMsg.textContent = result.error || '登录失败';
                    errorMsg.classList.remove('hidden');
                }
            } catch (error) {
                errorMsg.textContent = '网络错误，请重试';
                errorMsg.classList.remove('hidden');
            }
        });
    </script>
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

        <!-- Docker Compose 部署 -->
        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
                <span class="mr-3">🐳</span> Docker Compose 部署
            </h2>
            <div class="space-y-4">
                <p class="text-gray-700 mb-3">使用 Docker Compose 一键部署（推荐）：</p>
                <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
                    <pre class="text-green-400 font-mono text-sm">version: '3.8'
services:
  ztoapi:
    image: hulisang/ztoapi:latest
    container_name: ztoapi
    ports:
      - "9090:9090"
    environment:
      - DEFAULT_KEY=your-api-key
      - ZAI_TOKEN=your-zai-token
      - MODEL_NAME=GLM-4.6
      - DEBUG_MODE=true
      - ENABLE_THINKING=false
    volumes:
      - ./data:/app/data
    restart: unless-stopped</pre>
                </div>
                <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mt-3">
                    <p class="text-sm text-gray-700">💡 启动命令：<code class="bg-white px-2 py-1 rounded">docker-compose up -d</code></p>
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
            
            <!-- 基础配置 -->
            <div class="mb-6">
                <h3 class="text-lg font-bold text-gray-800 mb-3">基础配置</h3>
                <div class="space-y-3 text-sm">
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">DEFAULT_KEY</code>
                        <span class="text-gray-600 ml-2">- API 密钥（默认：sk-your-key）</span>
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
                        <code class="text-purple-600 font-mono">UPSTREAM_URL</code>
                        <span class="text-gray-600 ml-2">- 上游API地址（默认：https://chat.z.ai/api/chat/completions）</span>
                    </div>
                </div>
            </div>

            <!-- 功能开关 -->
            <div class="mb-6">
                <h3 class="text-lg font-bold text-gray-800 mb-3">功能开关</h3>
                <div class="space-y-3 text-sm">
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">DEBUG_MODE</code>
                        <span class="text-gray-600 ml-2">- 调试模式（默认：true）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">DEFAULT_STREAM</code>
                        <span class="text-gray-600 ml-2">- 默认流式响应（默认：true）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">ENABLE_THINKING</code>
                        <span class="text-gray-600 ml-2">- 启用思考功能（默认：false）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">DASHBOARD_ENABLED</code>
                        <span class="text-gray-600 ml-2">- 启用Dashboard统计面板（默认：true）</span>
                    </div>
                </div>
            </div>

            <!-- 管理系统配置 -->
            <div class="mb-6">
                <h3 class="text-lg font-bold text-gray-800 mb-3">管理系统配置</h3>
                <div class="space-y-3 text-sm">
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">REGISTER_ENABLED</code>
                        <span class="text-gray-600 ml-2">- 启用注册管理系统（默认：true）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">REGISTER_DB_PATH</code>
                        <span class="text-gray-600 ml-2">- 注册数据库路径（默认：./data/zai2api.db）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">ADMIN_ENABLED</code>
                        <span class="text-gray-600 ml-2">- 启用Admin面板（默认：true）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">ADMIN_USERNAME</code>
                        <span class="text-gray-600 ml-2">- Admin用户名（默认：admin）</span>
                    </div>
                    <div class="bg-gray-50 rounded p-3">
                        <code class="text-purple-600 font-mono">ADMIN_PASSWORD</code>
                        <span class="text-gray-600 ml-2">- Admin密码（默认：123456）</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- 管理系统说明 -->
        <div class="bg-white rounded-xl shadow-sm border p-8 mb-6">
            <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
                <span class="mr-3">🔧</span> 管理系统
            </h2>
            
            <div class="space-y-6">
                <!-- 注册管理系统 -->
                <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <h3 class="text-lg font-bold text-gray-900 mb-2">📝 注册管理系统</h3>
                    <p class="text-gray-700 mb-2">批量注册 Z.ai 账号，支持导入导出、批量获取 APIKEY</p>
                    <p class="text-sm text-gray-600">访问地址：<code class="bg-white px-2 py-1 rounded">http://localhost:9090/register/login</code></p>
                    <p class="text-sm text-gray-600 mt-1">默认账号：<code class="bg-white px-2 py-1 rounded">admin / 123456</code></p>
                </div>

                <!-- Admin 面板 -->
                <div class="bg-green-50 border border-green-200 rounded-lg p-4">
                    <h3 class="text-lg font-bold text-gray-900 mb-2">🔐 Admin 面板</h3>
                    <p class="text-gray-700 mb-2">账号管理、导入导出功能</p>
                    <p class="text-sm text-gray-600">访问地址：<code class="bg-white px-2 py-1 rounded">http://localhost:9090/admin</code></p>
                    <p class="text-sm text-gray-600 mt-1">默认账号：<code class="bg-white px-2 py-1 rounded">admin / 123456</code></p>
                </div>

                <!-- Dashboard -->
                <div class="bg-purple-50 border border-purple-200 rounded-lg p-4">
                    <h3 class="text-lg font-bold text-gray-900 mb-2">📊 Dashboard</h3>
                    <p class="text-gray-700 mb-2">实时监控 API 请求和性能统计</p>
                    <p class="text-sm text-gray-600">访问地址：<code class="bg-white px-2 py-1 rounded">http://localhost:9090/dashboard</code></p>
                </div>
            </div>
        </div>

        <div class="flex justify-center space-x-4">
            <a href="/" class="inline-block bg-purple-600 hover:bg-purple-700 text-white font-semibold px-8 py-3 rounded-lg transition">
                返回首页
            </a>
            <a href="/docs" class="inline-block bg-gray-600 hover:bg-gray-700 text-white font-semibold px-8 py-3 rounded-lg transition">
                API 文档
            </a>
        </div>
    </div>
</body>
</html>`
}

// 生成 Admin 面板 HTML
func getAdminPanelHTML() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>账号管理 - ZtoApi</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://code.jquery.com/jquery-3.7.1.min.js"></script>
</head>
<body class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 min-h-screen p-4 md:p-8">
    <div class="max-w-7xl mx-auto">
        <div class="text-center text-white mb-8">
            <div class="flex items-center justify-between">
                <div class="flex-1"></div>
                <div class="flex-1 text-center">
                    <h1 class="text-4xl md:text-5xl font-bold mb-3">🔐 ZtoApi 账号管理</h1>
                    <p class="text-lg md:text-xl opacity-90">导入导出 · 数据管理</p>
                </div>
                <div class="flex-1 flex justify-end gap-2">
                    <a href="/dashboard" class="px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg text-white font-semibold transition">
                        统计面板
                    </a>
                    <button id="logoutBtn" class="px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg text-white font-semibold transition">
                        退出登录
                    </button>
                </div>
            </div>
        </div>

        <!-- 统计面板 -->
        <div class="bg-white rounded-2xl shadow-2xl p-6 mb-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">统计信息</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-gradient-to-br from-green-400 to-emerald-500 rounded-xl p-4 text-center text-white">
                    <div class="text-sm opacity-90 mb-1">总账号数</div>
                    <div class="text-3xl font-bold" id="totalAccounts">0</div>
                </div>
                <div class="bg-gradient-to-br from-blue-400 to-indigo-500 rounded-xl p-4 text-center text-white">
                    <div class="text-sm opacity-90 mb-1">最近导入</div>
                    <div class="text-3xl font-bold" id="recentImport">0</div>
                </div>
            </div>
        </div>

        <!-- 账号列表 -->
        <div class="bg-white rounded-2xl shadow-2xl p-6 mb-6">
            <div class="flex items-center justify-between mb-4">
                <h2 class="text-2xl font-bold text-gray-800">账号列表</h2>
                <div class="flex gap-2">
                    <input type="text" id="searchInput" placeholder="搜索邮箱..."
                        class="px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                    <input type="file" id="importFileInput" accept=".txt" style="display: none;">
                    <button id="importBtn"
                        class="px-6 py-2 bg-gradient-to-r from-purple-500 to-violet-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition">
                        导入 TXT
                    </button>
                    <button id="exportBtn"
                        class="px-6 py-2 bg-gradient-to-r from-green-500 to-emerald-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition">
                        导出 TXT
                    </button>
                    <button id="refreshBtn"
                        class="px-6 py-2 bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition">
                        刷新
                    </button>
                </div>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="bg-gray-50 text-left">
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">序号</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">邮箱</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">密码</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">Token</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">APIKEY</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">创建时间</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">操作</th>
                        </tr>
                    </thead>
                    <tbody id="accountTableBody" class="divide-y divide-gray-200">
                        <tr>
                            <td colspan="7" class="px-4 py-8 text-center text-gray-400">加载中...</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- 分页控件 -->
            <div class="flex items-center justify-between mt-4 px-4 border-t pt-4">
                <div class="text-sm text-gray-600">
                    共 <span id="totalItems" class="font-semibold text-indigo-600">0</span> 条数据，
                    每页显示 <span id="currentPageSize" class="font-semibold text-indigo-600">20</span> 条
                </div>
                <div class="flex items-center gap-2">
                    <!-- 每页显示条数 -->
                    <select id="pageSizeSelect" class="px-3 py-2 border border-gray-300 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition text-sm">
                        <option value="10">10 条/页</option>
                        <option value="20" selected>20 条/页</option>
                        <option value="50">50 条/页</option>
                        <option value="100">100 条/页</option>
                    </select>

                    <!-- 页码按钮 -->
                    <div class="flex items-center gap-1">
                        <button id="firstPageBtn" class="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition text-sm font-medium" title="首页">
                            首页
                        </button>
                        <button id="prevPageBtn" class="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition text-sm font-medium" title="上一页">
                            上一页
                        </button>

                        <div class="flex items-center gap-1" id="pageNumbers"></div>

                        <button id="nextPageBtn" class="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition text-sm font-medium" title="下一页">
                            下一页
                        </button>
                        <button id="lastPageBtn" class="px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition text-sm font-medium" title="尾页">
                            尾页
                        </button>
                    </div>

                    <!-- 跳转页码 -->
                    <div class="flex items-center gap-2 ml-2">
                        <span class="text-sm text-gray-600">前往</span>
                        <input type="number" id="jumpPageInput" min="1" class="w-16 px-2 py-2 border border-gray-300 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition text-sm text-center">
                        <button id="jumpPageBtn" class="px-3 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition text-sm font-medium">
                            跳转
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        let accounts = [];
        let filteredAccounts = [];
        let currentPage = 1;
        let pageSize = 20;

        const $totalAccounts = $('#totalAccounts');
        const $recentImport = $('#recentImport');
        const $accountTableBody = $('#accountTableBody');
        const $searchInput = $('#searchInput');
        const $totalItems = $('#totalItems');
        const $currentPageSize = $('#currentPageSize');
        const $pageNumbers = $('#pageNumbers');
        const $jumpPageInput = $('#jumpPageInput');

        // 渲染表格（带分页）
        function renderTable(data = filteredAccounts) {
            const totalPages = Math.ceil(data.length / pageSize);

            // 边界检查
            if (currentPage < 1) currentPage = 1;
            if (currentPage > totalPages && totalPages > 0) currentPage = totalPages;
            if (totalPages === 0) currentPage = 1;

            const startIndex = (currentPage - 1) * pageSize;
            const endIndex = startIndex + pageSize;
            const pageData = data.slice(startIndex, endIndex);

            if (data.length === 0) {
                $accountTableBody.html('<tr><td colspan="7" class="px-4 py-8 text-center text-gray-400">暂无数据</td></tr>');
            } else if (pageData.length === 0) {
                $accountTableBody.html('<tr><td colspan="7" class="px-4 py-8 text-center text-gray-400">当前页无数据</td></tr>');
            } else {
                const rows = pageData.map((acc, idx) => {
                    const globalIndex = startIndex + idx + 1;
                    const tokenPreview = acc.token.length > 20 ? acc.token.substring(0, 20) + '...' : acc.token;

                    // 处理 APIKEY 显示
                    const apikeyDisplay = acc.apikey ?
                        '<code class="bg-gray-100 px-2 py-1 rounded text-xs">' + (acc.apikey.length > 20 ? acc.apikey.substring(0, 20) + '...' : acc.apikey) + '</code>' :
                        '<span class="text-gray-400 text-xs">未生成</span>';

                    return '<tr class="hover:bg-gray-50">' +
                        '<td class="px-4 py-3 text-sm text-gray-700">' + globalIndex + '</td>' +
                        '<td class="px-4 py-3 text-sm text-gray-700">' + acc.email + '</td>' +
                        '<td class="px-4 py-3 text-sm text-gray-700"><code class="bg-gray-100 px-2 py-1 rounded">' + acc.password + '</code></td>' +
                        '<td class="px-4 py-3 text-sm text-gray-700"><code class="bg-gray-100 px-2 py-1 rounded text-xs">' + tokenPreview + '</code></td>' +
                        '<td class="px-4 py-3 text-sm text-gray-700">' + apikeyDisplay + '</td>' +
                        '<td class="px-4 py-3 text-sm text-gray-700">' + new Date(acc.createdAt).toLocaleString('zh-CN') + '</td>' +
                        '<td class="px-4 py-3 flex gap-2">' +
                            '<button class="copy-account-btn text-blue-600 hover:text-blue-800 text-sm font-medium" data-email="' + acc.email + '" data-password="' + acc.password + '">复制账号</button>' +
                            '<button class="copy-token-btn text-green-600 hover:text-green-800 text-sm font-medium" data-token="' + acc.token + '">复制Token</button>' +
                            (acc.apikey ? '<button class="copy-apikey-btn text-purple-600 hover:text-purple-800 text-sm font-medium" data-apikey="' + acc.apikey + '">复制APIKEY</button>' : '') +
                        '</td>' +
                    '</tr>';
                });
                $accountTableBody.html(rows.join(''));

                // 绑定复制事件
                $('.copy-account-btn').on('click', function() {
                    const email = $(this).data('email');
                    const password = $(this).data('password');
                    navigator.clipboard.writeText(email + '----' + password);
                    alert('已复制账号: ' + email);
                });

                $('.copy-token-btn').on('click', function() {
                    const token = $(this).data('token');
                    navigator.clipboard.writeText(token);
                    alert('已复制 Token');
                });

                $('.copy-apikey-btn').on('click', function() {
                    const apikey = $(this).data('apikey');
                    navigator.clipboard.writeText(apikey);
                    alert('已复制 APIKEY');
                });
            }

            // 更新分页信息
            updatePagination(data.length, totalPages);
        }

        // 更新分页控件
        function updatePagination(totalItems, totalPages) {
            $totalItems.text(totalItems);
            $currentPageSize.text(pageSize);

            // 更新按钮状态
            $('#firstPageBtn, #prevPageBtn').prop('disabled', currentPage === 1);
            $('#nextPageBtn, #lastPageBtn').prop('disabled', currentPage === totalPages || totalPages === 0);

            // 渲染页码
            $pageNumbers.empty();
            if (totalPages === 0) return;

            const pagerCount = 7;
            let showPrevMore = false;
            let showNextMore = false;

            if (totalPages > pagerCount) {
                if (currentPage > pagerCount - 3) showPrevMore = true;
                if (currentPage < totalPages - 3) showNextMore = true;
            }

            const array = [];

            if (showPrevMore && !showNextMore) {
                const startPage = totalPages - (pagerCount - 2);
                for (let i = startPage; i < totalPages; i++) array.push(i);
            } else if (!showPrevMore && showNextMore) {
                for (let i = 2; i < pagerCount; i++) array.push(i);
            } else if (showPrevMore && showNextMore) {
                const offset = Math.floor(pagerCount / 2) - 1;
                for (let i = currentPage - offset; i <= currentPage + offset; i++) array.push(i);
            } else {
                for (let i = 2; i < totalPages; i++) array.push(i);
            }

            // 第一页
            addPageButton(1, $pageNumbers);

            // 前省略号
            if (showPrevMore) {
                $pageNumbers.append('<button class="px-3 py-2 text-gray-600 hover:text-indigo-600 transition more-btn" data-action="prev-more">...</button>');
            }

            // 中间页码
            array.forEach(page => addPageButton(page, $pageNumbers));

            // 后省略号
            if (showNextMore) {
                $pageNumbers.append('<button class="px-3 py-2 text-gray-600 hover:text-indigo-600 transition more-btn" data-action="next-more">...</button>');
            }

            // 最后一页
            if (totalPages > 1) addPageButton(totalPages, $pageNumbers);

            // 绑定省略号点击事件
            $('.more-btn').on('click', function() {
                const action = $(this).data('action');
                if (action === 'prev-more') {
                    currentPage = Math.max(1, currentPage - 5);
                } else if (action === 'next-more') {
                    currentPage = Math.min(totalPages, currentPage + 5);
                }
                renderTable();
            });
        }

        // 添加页码按钮
        function addPageButton(page, container) {
            const isActive = page === currentPage;
            const $btn = $('<button>', {
                text: page,
                class: 'px-3 py-2 rounded-lg transition text-sm font-medium ' +
                       (isActive ? 'bg-indigo-600 text-white' : 'border border-gray-300 hover:bg-gray-50'),
                click: () => {
                    currentPage = page;
                    renderTable();
                }
            });
            container.append($btn);
        }

        // 加载账号列表
        async function loadAccounts() {
            try {
                const response = await fetch('/admin/api/accounts');
                accounts = await response.json();
                filteredAccounts = accounts;
                $totalAccounts.text(accounts.length);
                currentPage = 1;
                renderTable();
            } catch (error) {
                alert('加载账号失败: ' + error.message);
            }
        }

        // 搜索功能
        $searchInput.on('input', function() {
            const keyword = $(this).val().toLowerCase();
            filteredAccounts = accounts.filter(acc => acc.email.toLowerCase().includes(keyword));
            currentPage = 1;
            renderTable();
        });

        // 分页按钮事件
        $('#firstPageBtn').on('click', () => { currentPage = 1; renderTable(); });
        $('#prevPageBtn').on('click', () => { if (currentPage > 1) { currentPage--; renderTable(); } });
        $('#nextPageBtn').on('click', () => {
            const totalPages = Math.ceil(filteredAccounts.length / pageSize);
            if (currentPage < totalPages) { currentPage++; renderTable(); }
        });
        $('#lastPageBtn').on('click', () => {
            const totalPages = Math.ceil(filteredAccounts.length / pageSize);
            currentPage = totalPages;
            renderTable();
        });

        // 每页显示条数变更
        $('#pageSizeSelect').on('change', function() {
            pageSize = parseInt($(this).val());
            currentPage = 1;
            renderTable();
        });

        // 跳转到指定页
        $('#jumpPageBtn').on('click', () => {
            const targetPage = parseInt($jumpPageInput.val());
            const totalPages = Math.ceil(filteredAccounts.length / pageSize);

            if (!targetPage || targetPage < 1 || targetPage > totalPages) {
                alert('请输入有效的页码（1-' + totalPages + '）');
                return;
            }

            currentPage = targetPage;
            $jumpPageInput.val('');
            renderTable();
        });

        // 回车快速跳转
        $jumpPageInput.on('keypress', function(e) {
            if (e.which === 13) $('#jumpPageBtn').click();
        });

        $('#refreshBtn').on('click', loadAccounts);

        $('#exportBtn').on('click', async function() {
            try {
                const response = await fetch('/admin/api/export');
                const blob = await response.blob();
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'zai_accounts_' + Date.now() + '.txt';
                a.click();
                alert('导出成功！');
            } catch (error) {
                alert('导出失败: ' + error.message);
            }
        });

        $('#importBtn').on('click', function() {
            $('#importFileInput').click();
        });

        $('#importFileInput').on('change', async function(e) {
            const file = e.target.files[0];
            if (!file) return;

            try {
                const text = await file.text();
                const lines = text.split('\n').filter(line => line.trim());

                const importData = [];
                const emailSet = new Set();

                for (const line of lines) {
                    const parts = line.split('----');
                    let email, password, token, apikey;

                    if (parts.length >= 5) {
                        email = parts[0].trim();
                        password = parts[1].trim();
                        token = parts[2].trim() + '----' + parts[3].trim();
                        apikey = parts[4].trim() || null;
                    } else if (parts.length === 4) {
                        email = parts[0].trim();
                        password = parts[1].trim();
                        token = parts[2].trim();
                        apikey = parts[3].trim() || null;
                    } else if (parts.length === 3) {
                        email = parts[0].trim();
                        password = parts[1].trim();
                        token = parts[2].trim();
                        apikey = null;
                    } else {
                        continue;
                    }

                    if (!emailSet.has(email)) {
                        emailSet.add(email);
                        importData.push({ email, password, token, apikey });
                    }
                }

                const response = await fetch('/admin/api/import-batch', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ accounts: importData })
                });

                const result = await response.json();
                if (result.success) {
                    alert('导入完成！成功: ' + result.imported + ', 跳过重复: ' + result.skipped);
                    $recentImport.text(result.imported);
                    await loadAccounts();
                } else {
                    alert('导入失败: ' + result.error);
                }

                $(this).val('');
            } catch (error) {
                alert('导入失败: ' + error.message);
            }
        });

        $('#logoutBtn').on('click', async function() {
            if (confirm('确定要退出登录吗？')) {
                await fetch('/admin/api/logout', { method: 'POST' });
                document.cookie = 'adminSessionId=; path=/; max-age=0';
                window.location.href = '/admin/login';
            }
        });

        $(document).ready(function() {
            loadAccounts();
        });
    </script>
</body>
</html>`
}
