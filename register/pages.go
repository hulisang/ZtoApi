package register

// 登录页面
const LoginPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>登录 - Z.AI 注册管理系统</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 min-h-screen flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        <div class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-2">🤖 Z.AI 注册管理系统</h1>
            <p class="text-gray-600">请登录以继续</p>
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
        <div class="mt-2 text-center text-sm text-gray-500">
            <p>📦 <a href="https://github.com/hulisang/ZtoApi" target="_blank" class="text-cyan-600 underline">源码地址 (GitHub)</a></p>
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
                const response = await fetch('/register/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                const result = await response.json();

                if (result.success) {
                    document.cookie = 'sessionId=' + result.sessionId + '; path=/; max-age=86400';
                    window.location.href = '/register/';
                } else {
                    let errorText = result.error || '登录失败';
                    if (result.locked) {
                        errorText += ' (账号已锁定' + Math.ceil(result.remainingTime) + '秒)';
                    }
                    errorMsg.textContent = errorText;
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

// 主管理页面（包含注册、账号管理、实时日志等所有功能）
const MainPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Z.AI 注册管理系统</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://code.jquery.com/jquery-3.7.1.min.js"></script>
    <style>
        @keyframes slideIn { from { transform: translateX(100%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
        .toast-enter { animation: slideIn 0.3s ease-out; }
        ::-webkit-scrollbar { width: 8px; height: 8px; }
        ::-webkit-scrollbar-track { background: #f1f5f9; border-radius: 4px; }
        ::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 4px; }
        
        /* PC端表格悬停效果 */
        @media (min-width: 769px) {
            tbody tr {
                transition: all 0.2s ease;
            }
            
            tbody tr:hover {
                background-color: #f8fafc;
                transform: translateX(4px);
                box-shadow: -4px 0 0 0 #6366f1;
            }
            
            /* 可复制单元格样式 */
            .clickable-cell {
                cursor: pointer;
                transition: all 0.15s ease;
                position: relative;
            }
            
            .clickable-cell:hover {
                opacity: 0.7;
            }
            
            /* 复制图标 */
            .clickable-cell::before {
                content: '📋';
                position: absolute;
                left: 0;
                top: 50%;
                transform: translateY(-50%);
                opacity: 0;
                transition: opacity 0.2s ease;
                font-size: 0.75rem;
            }
            
            .clickable-cell:hover::before {
                opacity: 0.5;
            }
        }
    </style>
</head>
<body class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 min-h-screen p-4">
    <div id="toastContainer" class="fixed top-4 right-4 z-50 space-y-2"></div>

    <div class="max-w-7xl mx-auto">
        <div class="text-center text-white mb-8">
            <h1 class="text-4xl font-bold mb-2">🤖 Z.AI 注册管理系统</h1>
            <p class="text-xl opacity-90">批量注册 · 数据管理 · 实时监控</p>
            <button id="logoutBtn" class="mt-4 px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg text-white font-semibold transition">
                退出登录
            </button>
        </div>

        <!-- 注册控制面板 -->
        <div class="bg-white rounded-2xl shadow-2xl p-6 mb-6">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 sm:gap-0 mb-4 sm:mb-6">
                <h2 class="text-xl sm:text-2xl font-bold text-gray-800">注册控制</h2>
                <div class="flex gap-2 w-full sm:w-auto">
                    <button id="settingsBtn" class="flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg font-semibold transition text-sm sm:text-base">
                        ⚙️ 高级设置
                    </button>
                    <span id="statusBadge" class="flex-1 sm:flex-none px-3 sm:px-4 py-2 rounded-full text-xs sm:text-sm font-semibold bg-gray-400 text-white text-center">闲置中</span>
                </div>
            </div>
            
            <!-- 高级设置面板 -->
            <div id="settingsPanel" class="mb-6 p-4 bg-gray-50 rounded-lg hidden">
                <h3 class="font-semibold text-gray-700 mb-4">⚙️ 高级设置</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">邮件等待超时 (秒)</label>
                        <input type="number" id="emailTimeout" value="300" min="60" max="600"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">建议：300秒（5分钟），最多10分钟</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">账号间隔 (毫秒)</label>
                        <input type="number" id="registerDelay" value="2000" min="500" max="10000" step="500"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">建议：2000ms（2秒），更稳定</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">邮件轮询间隔（秒）</label>
                        <input type="number" id="emailCheckInterval" value="5" min="1" max="30" step="1"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">建议：3-10秒，过小可能触发限流</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">并发数</label>
                        <input type="number" id="concurrency" value="15" min="1" max="100"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">建议：10-30，过高可能被封</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">API 重试次数</label>
                        <input type="number" id="retryTimes" value="3" min="1" max="10"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                    </div>
                    <div class="flex items-center">
                        <input type="checkbox" id="skipApikey" class="w-5 h-5 text-indigo-600 rounded">
                        <label class="ml-3 text-sm font-medium text-gray-700">🚀 快速模式（注册后稍后批量获取APIKEY）</label>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">PushPlus Token</label>
                        <input type="text" id="pushplusToken" value="" placeholder="留空则不发送通知"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">HTTP超时 (秒)</label>
                        <input type="number" id="httpTimeout" value="30" min="5" max="120"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">默认30秒</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">批量保存大小</label>
                        <input type="number" id="batchSaveSize" value="10" min="1" max="100"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">默认10条</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">连接池大小</label>
                        <input type="number" id="connectionPoolSize" value="100" min="10" max="500"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                        <p class="text-xs text-gray-500 mt-1">默认100</p>
                    </div>
                    <div class="flex items-center md:col-span-2">
                        <input type="checkbox" id="enableNotification" class="w-5 h-5 text-indigo-600 rounded">
                        <label class="ml-3 text-sm font-medium text-gray-700">启用 PushPlus 通知</label>
                    </div>
                </div>
                <div class="mt-4 flex gap-2">
                    <button id="saveSettingsBtn" class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition">
                        保存设置
                    </button>
                    <button id="cancelSettingsBtn" class="px-6 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition">
                        取消
                    </button>
                </div>
            </div>
            
            <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 mb-4">
                <div class="flex-1">
                    <label for="registerCount" class="block text-sm font-medium text-gray-700 mb-2">注册数量</label>
                    <input type="number" id="registerCount" value="5" min="1" max="1000"
                        placeholder="输入要注册的账号数量"
                        class="w-full px-4 py-3 text-base border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                </div>
                <button id="startRegisterBtn"
                    class="w-full sm:w-auto px-6 sm:px-8 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all disabled:opacity-60 disabled:cursor-not-allowed text-base self-end">
                    开始注册
                </button>
                <button id="stopRegisterBtn" style="display: none;"
                    class="w-full sm:w-auto px-6 sm:px-8 py-3 bg-gradient-to-r from-red-500 to-pink-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all text-base">
                    停止注册
                </button>
            </div>

            <div id="progressContainer" style="display: none;" class="mt-4">
                <div class="flex justify-between text-sm text-gray-600 mb-2">
                    <span>注册进度</span>
                    <span id="progressText">0/0 (0%)</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-4">
                    <div id="progressBar" class="h-full bg-gradient-to-r from-indigo-500 to-purple-600 rounded-full transition-all"></div>
                </div>
            </div>
        </div>

        <!-- 统计信息 -->
        <div class="bg-white rounded-2xl shadow-2xl p-6 mb-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">统计信息</h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="bg-gradient-to-br from-green-400 to-emerald-500 rounded-xl p-4 text-center text-white">
                    <div class="text-sm opacity-90 mb-1">总账号</div>
                    <div class="text-3xl font-bold" id="totalAccounts">0</div>
                </div>
                <div class="bg-gradient-to-br from-purple-400 to-violet-500 rounded-xl p-4 text-center text-white">
                    <div class="text-sm opacity-90 mb-1">有APIKEY</div>
                    <div class="text-3xl font-bold" id="withApikey">0</div>
                </div>
                <div class="bg-gradient-to-br from-orange-400 to-red-500 rounded-xl p-4 text-center text-white">
                    <div class="text-sm opacity-90 mb-1">无APIKEY</div>
                    <div class="text-3xl font-bold" id="withoutApikey">0</div>
                </div>
                <div class="bg-gradient-to-br from-blue-400 to-indigo-500 rounded-xl p-4 text-center text-white">
                    <div class="text-sm opacity-90 mb-1">活跃账号</div>
                    <div class="text-3xl font-bold" id="activeAccounts">0</div>
                </div>
            </div>
        </div>

        <!-- 账号列表 -->
        <div class="bg-white rounded-2xl shadow-2xl p-6 mb-6">
            <div class="flex items-center justify-between mb-4">
                <div class="flex items-center gap-3">
                    <h2 class="text-2xl font-bold text-gray-800">账号列表</h2>
                    <span id="selectedCount" class="hidden px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-sm font-semibold">已选 0 项</span>
                </div>
                <div class="flex flex-wrap gap-2">
                    <input type="text" id="searchInput" placeholder="搜索邮箱/密码/Token/APIKEY..."
                        class="px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500">
                    <button id="batchRefetchApikeyBtn" class="px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-600 text-white rounded-lg hover:shadow-lg transition text-sm">
                        🔑 批量补充APIKEY
                    </button>
                    <button id="batchCheckAccountsBtn" class="px-4 py-2 bg-gradient-to-r from-yellow-500 to-orange-600 text-white rounded-lg hover:shadow-lg transition text-sm">
                        🔍 批量检测存活
                    </button>
                    <button id="deleteInactiveBtn" class="px-4 py-2 bg-gradient-to-r from-red-500 to-pink-600 text-white rounded-lg hover:shadow-lg transition text-sm">
                        🗑️ 删除失效账号
                    </button>
                    <button id="refreshBtn" class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition">
                        🔃 刷新
                    </button>
                    <button id="exportBtn" class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition">
                        📤 导出
                    </button>
                    <button id="importBtn" class="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition">
                        📥 导入
                    </button>
                    <input type="file" id="importFileInput" accept=".txt" style="display: none;">
                </div>
            </div>
            
            <!-- 快速筛选标签 -->
            <div class="flex flex-wrap gap-2 mb-4">
                <span class="text-sm text-gray-600 font-medium self-center">快速筛选:</span>
                <button class="quick-filter-btn px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-full transition" data-filter="today">
                    📅 今日注册
                </button>
                <button class="quick-filter-btn px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-full transition" data-filter="week">
                    📆 本周注册
                </button>
                <button class="quick-filter-btn px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-full transition" data-filter="inactive">
                    ⚠️ 失效账号
                </button>
                <button class="quick-filter-btn px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-full transition" data-filter="no-apikey">
                    🔒 无APIKEY
                </button>
                <button class="quick-filter-btn px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-full transition" data-filter="has-apikey">
                    🔑 有APIKEY
                </button>
                <button id="clearFilterBtn" class="px-3 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 rounded-full transition hidden">
                    ✖ 清除筛选
                </button>
            </div>
            
            <!-- 批量操作按钮区域 -->
            <div id="batchActionsBar" class="hidden mb-4 p-3 bg-indigo-50 rounded-lg border-2 border-indigo-200">
                <div class="flex flex-wrap gap-2">
                    <button id="batchDeleteBtn" class="px-4 py-2 bg-red-500 hover:bg-red-600 text-white font-semibold rounded-lg transition text-sm">
                        🗑️ 批量删除
                    </button>
                    <button id="batchCopyEmailsBtn" class="px-4 py-2 bg-purple-500 hover:bg-purple-600 text-white font-semibold rounded-lg transition text-sm">
                        📋 复制邮箱
                    </button>
                    <button id="batchCopyTokensBtn" class="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white font-semibold rounded-lg transition text-sm">
                        🔑 复制Token
                    </button>
                    <button id="cancelSelectionBtn" class="ml-auto px-4 py-2 bg-gray-400 hover:bg-gray-500 text-white font-semibold rounded-lg transition text-sm">
                        ✖️ 取消选择
                    </button>
                </div>
            </div>
            
            <div class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="bg-gray-50 text-left">
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">
                                <input type="checkbox" id="selectAllCheckbox" class="w-4 h-4 text-indigo-600 rounded cursor-pointer">
                            </th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">序号</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">邮箱</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">密码</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">Token</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">APIKEY</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">状态</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">创建时间</th>
                            <th class="px-4 py-3 text-sm font-semibold text-gray-700">操作</th>
                        </tr>
                    </thead>
                    <tbody id="accountTableBody" class="divide-y divide-gray-200">
                        <tr><td colspan="9" class="px-4 py-8 text-center text-gray-400">暂无数据</td></tr>
                    </tbody>
                </table>
            </div>
            
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mt-4">
                <div class="text-sm text-gray-600">
                    共 <span id="totalItems">0</span> 条数据
                </div>
                <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 w-full sm:w-auto">
                    <div class="flex items-center gap-1 sm:gap-2 overflow-x-auto">
                        <button id="firstPageBtn" class="px-2 sm:px-3 py-1 text-xs sm:text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap">首页</button>
                        <button id="prevPageBtn" class="px-2 sm:px-3 py-1 text-xs sm:text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap">上一页</button>
                        <div class="flex items-center gap-1" id="pageNumbers"></div>
                        <button id="nextPageBtn" class="px-2 sm:px-3 py-1 text-xs sm:text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap">下一页</button>
                        <button id="lastPageBtn" class="px-2 sm:px-3 py-1 text-xs sm:text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap">尾页</button>
                    </div>
                    <select id="pageSizeSelect" class="px-2 py-1 text-xs sm:text-sm border border-gray-300 rounded w-full sm:w-auto">
                        <option value="10">10条/页</option>
                        <option value="20" selected>20条/页</option>
                        <option value="50">50条/页</option>
                        <option value="100">100条/页</option>
                    </select>
                </div>
            </div>
        </div>

        <!-- 实时日志 -->
        <div class="bg-white rounded-2xl shadow-2xl p-6">
            <div class="flex items-center justify-between mb-4">
                <h2 class="text-2xl font-bold text-gray-800">实时日志</h2>
                <button id="clearLogBtn" class="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition">
                    清空日志
                </button>
            </div>
            <div id="logContainer" class="bg-gray-900 rounded-lg p-4 h-64 overflow-y-auto font-mono text-sm text-blue-400">
                <div>等待任务启动...</div>
            </div>
        </div>
    </div>

    <script>
        let accounts = [];
        let currentPage = 1;
        let pageSize = 20;
        let totalPages = 1;
        let totalItems = 0;
        let isRunning = false;
        let eventSource = null;
        let currentConfig = null;
        let selectedAccounts = new Set(); // 选中的账号ID
        let currentFilter = ''; // 当前筛选条件
        let searchKeyword = ''; // 搜索关键词

        function showToast(message, type = 'info') {
            const colors = { success: 'bg-green-500', error: 'bg-red-500', info: 'bg-blue-500' };
            const $toast = $('<div>', {
                class: 'toast-enter ' + colors[type] + ' text-white px-6 py-3 rounded-lg shadow-lg',
                text: message
            });
            $('#toastContainer').append($toast);
            setTimeout(() => $toast.remove(), 3000);
        }

        function addLog(message, level = 'info', link = null) {
            const colors = { success: 'text-green-400', error: 'text-red-400', info: 'text-blue-400', warning: 'text-yellow-400' };
            const time = new Date().toLocaleTimeString();
            
            let logHTML = '<span class="text-gray-500">[' + time + ']</span> ' + message;
            
            // 如果有链接，添加链接按钮
            if (link && link.url && link.text) {
                logHTML += ' <a href="' + link.url + '" target="_blank" class="ml-2 px-2 py-0.5 bg-blue-600 hover:bg-blue-700 text-white text-xs rounded transition">' + link.text + '</a>';
            }
            
            const $log = $('<div>', {
                class: colors[level] + ' mb-1',
                html: logHTML
            });
            $('#logContainer').append($log);
            $('#logContainer')[0].scrollTop = $('#logContainer')[0].scrollHeight;
        }

        async function loadAccounts() {
            try {
                let url = '/register/api/accounts?page=' + currentPage + '&pageSize=' + pageSize;
                if (currentFilter) {
                    url += '&filter=' + currentFilter;
                }
                if (searchKeyword) {
                    url += '&search=' + encodeURIComponent(searchKeyword);
                }
                const response = await fetch(url);
                const data = await response.json();
                accounts = data.accounts || [];
                
                // 保存分页信息
                if (data.pagination) {
                    totalItems = data.pagination.total || 0;
                    totalPages = Math.ceil(totalItems / pageSize);
                }
                
                renderTable();
                updateStats();
            } catch (error) {
                showToast('加载账号列表失败', 'error');
            }
        }

        async function updateStats() {
            try {
                const response = await fetch('/register/api/stats');
                const stats = await response.json();
                $('#totalAccounts').text(stats.totalAccounts);
                $('#withApikey').text(stats.withAPIKEY);
                $('#withoutApikey').text(stats.withoutAPIKEY);
                $('#activeAccounts').text(stats.activeAccounts);
            } catch (error) {
                console.error('Failed to load stats:', error);
            }
        }

        function renderTable() {
            if (accounts.length === 0) {
                $('#accountTableBody').html('<tr><td colspan="9" class="px-4 py-8 text-center text-gray-400">暂无数据</td></tr>');
                return;
            }

            const rows = accounts.map((acc, idx) => {
                const apikeyDisplay = acc.apikey ?
                    '<code class="bg-indigo-50 text-indigo-700 px-2 py-1 rounded text-xs">' + acc.apikey.substring(0, 20) + '...</code>' :
                    '<span class="text-gray-400 text-xs">未生成</span>';
                
                let statusDisplay;
                if (acc.status === 'active') {
                    statusDisplay = '<span class="px-2 py-1 bg-green-100 text-green-700 rounded-full text-xs">✓ 正常</span>';
                } else if (acc.status === 'inactive') {
                    statusDisplay = '<span class="px-2 py-1 bg-red-100 text-red-700 rounded-full text-xs">✗ 失效</span>';
                } else {
                    statusDisplay = '<span class="px-2 py-1 bg-gray-100 text-gray-700 rounded-full text-xs">? 未知</span>';
                }

                const isChecked = selectedAccounts.has(acc.email);

                return '<tr>' +
                    '<td class="px-4 py-3"><input type="checkbox" class="account-checkbox w-4 h-4 text-indigo-600 rounded cursor-pointer" data-email="' + acc.email + '" ' + (isChecked ? 'checked' : '') + '></td>' +
                    '<td class="px-4 py-3 text-sm">' + (idx + 1) + '</td>' +
                    '<td class="px-4 py-3 text-sm clickable-cell" data-copy="' + acc.email + '" title="点击复制 邮箱">' + acc.email + '</td>' +
                    '<td class="px-4 py-3 text-sm clickable-cell" data-copy="' + acc.password + '" title="点击复制 密码"><code class="bg-blue-50 text-blue-700 px-2 py-1 rounded text-xs">' + acc.password + '</code></td>' +
                    '<td class="px-4 py-3 text-sm clickable-cell" data-copy="' + acc.token + '" title="点击复制 Token"><code class="bg-green-50 text-green-700 px-2 py-1 rounded text-xs">' + acc.token.substring(0, 20) + '...</code></td>' +
                    '<td class="px-4 py-3 text-sm' + (acc.apikey ? ' clickable-cell' : '') + '"' + (acc.apikey ? ' data-copy="' + acc.apikey + '" title="点击复制 APIKEY"' : '') + '>' + apikeyDisplay + '</td>' +
                    '<td class="px-4 py-3 text-center">' + statusDisplay + '</td>' +
                    '<td class="px-4 py-3 text-sm">' + new Date(acc.createdAt).toLocaleString() + '</td>' +
                    '<td class="px-4 py-3"><div class="flex gap-2">' +
                        (!acc.apikey ? '<button class="refetch-apikey-btn text-green-600 hover:text-green-800 text-sm font-medium whitespace-nowrap" data-email="' + acc.email + '" data-token="' + acc.token + '">🔑 获取KEY</button>' : '') +
                        '<button class="delete-btn text-red-600 hover:text-red-800 text-sm" data-email="' + acc.email + '">删除</button>' +
                    '</div></td>' +
                    '</tr>';
            });
            $('#accountTableBody').html(rows.join(''));

            // 更新选中计数
            updateSelectionCount();

            // 绑定复选框事件
            $('.account-checkbox').on('change', function() {
                const email = $(this).data('email');
                if ($(this).is(':checked')) {
                    selectedAccounts.add(email);
                } else {
                    selectedAccounts.delete(email);
                }
                updateSelectionCount();
            });

            // 绑定点击复制事件
            $('.clickable-cell').on('click', function() {
                const text = $(this).data('copy');
                navigator.clipboard.writeText(text);
                showToast('已复制', 'success');
            });

            // 绑定删除事件
            $('.delete-btn').on('click', async function() {
                const email = $(this).data('email');
                if (!confirm('确定要删除账号 ' + email + ' 吗？')) return;
                try {
                    await fetch('/register/api/accounts/delete', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ email: email })
                    });
                    showToast('删除成功', 'success');
                    loadAccounts();
                } catch (error) {
                    showToast('删除失败', 'error');
                }
            });

            // 绑定"获取APIKEY"按钮事件
            $('.refetch-apikey-btn').on('click', async function() {
                const email = $(this).data('email');
                const token = $(this).data('token');
                const $btn = $(this);
                const originalText = $btn.text();
                
                $btn.prop('disabled', true).text('获取中...');
                
                try {
                    const response = await fetch('/register/api/refetch-apikey', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ email, token })
                    });

                    const result = await response.json();

                    if (result.success) {
                        showToast('✓ ' + email + ' APIKEY获取成功', 'success');
                        loadAccounts();
                    } else {
                        showToast('✗ ' + email + ' ' + result.error, 'error');
                        $btn.prop('disabled', false).text(originalText);
                    }
                } catch (error) {
                    showToast('✗ ' + email + ' 获取失败: ' + error, 'error');
                    $btn.prop('disabled', false).text(originalText);
                }
            });
            
            // 更新分页控件
            updatePagination();
        }

        // 更新分页控件
        function updatePagination() {
            $('#totalItems').text(totalItems);

            // 更新按钮状态
            $('#firstPageBtn, #prevPageBtn').prop('disabled', currentPage === 1);
            $('#nextPageBtn, #lastPageBtn').prop('disabled', currentPage === totalPages || totalPages === 0);

            // 渲染页码
            const $pageNumbers = $('#pageNumbers');
            $pageNumbers.empty();

            if (totalPages <= 7) {
                // 总页数 <= 7，显示所有页码
                for (let i = 1; i <= totalPages; i++) {
                    addPageButton(i, $pageNumbers);
                }
            } else {
                // 总页数 > 7，智能显示
                addPageButton(1, $pageNumbers);
                if (currentPage > 3) {
                    $pageNumbers.append('<span class="px-2 text-gray-400">...</span>');
                }

                let start = Math.max(2, currentPage - 1);
                let end = Math.min(totalPages - 1, currentPage + 1);

                for (let i = start; i <= end; i++) {
                    addPageButton(i, $pageNumbers);
                }

                if (currentPage < totalPages - 2) {
                    $pageNumbers.append('<span class="px-2 text-gray-400">...</span>');
                }
                addPageButton(totalPages, $pageNumbers);
            }
        }

        // 添加页码按钮
        function addPageButton(page, container) {
            const isActive = page === currentPage;
            const $btn = $('<button>', {
                text: page,
                class: 'px-2 sm:px-3 py-1 text-xs sm:text-sm border rounded ' + (isActive ? 'bg-indigo-600 text-white border-indigo-600' : 'border-gray-300 hover:bg-gray-100'),
                click: () => {
                    currentPage = page;
                    loadAccounts();
                }
            });
            container.append($btn);
        }

        // 更新选中计数
        function updateSelectionCount() {
            if (selectedAccounts.size > 0) {
                $('#selectedCount').text('已选 ' + selectedAccounts.size + ' 项').removeClass('hidden');
                $('#batchActionsBar').removeClass('hidden');
            } else {
                $('#selectedCount').addClass('hidden');
                $('#batchActionsBar').addClass('hidden');
            }
        }

        $('#stopRegisterBtn').on('click', async function() {
            try {
                await fetch('/register/api/register/stop', { method: 'POST' });
                $('#startRegisterBtn').show();
                $('#stopRegisterBtn').hide();
                $('#statusBadge').text('闲置中').removeClass('bg-green-500').addClass('bg-gray-400');
                showToast('正在停止...', 'info');
            } catch (error) {
                showToast('停止失败', 'error');
            }
        });

        $('#exportBtn').on('click', function() {
            window.open('/register/api/accounts/export', '_blank');
        });

        $('#importBtn').on('click', function() {
            $('#importFileInput').click();
        });

        $('#importFileInput').on('change', async function(e) {
            const file = e.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            try {
                const response = await fetch('/register/api/accounts/import', {
                    method: 'POST',
                    body: formData
                });

                const result = await response.json();
                if (result.success) {
                    showToast('导入成功！成功: ' + result.imported + ', 失败: ' + result.failed, 'success');
                    loadAccounts();
                } else {
                    showToast('导入失败', 'error');
                }
            } catch (error) {
                showToast('导入失败: ' + error, 'error');
            }

            $(this).val('');
        });

        $('#refreshBtn').on('click', function() {
            loadAccounts();
        });

        $('#clearLogBtn').on('click', function() {
            $('#logContainer').html('<div class="text-blue-400">日志已清空</div>');
        });

        $('#logoutBtn').on('click', function() {
            document.cookie = 'sessionId=; path=/; max-age=0';
            window.location.href = '/register/login';
        });

        // 加载配置
        async function loadConfig() {
            try {
                const response = await fetch('/register/api/config');
                const config = await response.json();
                currentConfig = config;
                
                // 填充高级设置表单
                $('#emailTimeout').val(config.emailTimeout);
                $('#registerDelay').val(config.registerDelay);
                $('#emailCheckInterval').val(config.emailCheckInterval);
                $('#concurrency').val(config.concurrency);
                $('#retryTimes').val(config.retryTimes);
                $('#skipApikey').prop('checked', config.skipApikeyOnRegister);
                $('#pushplusToken').val(config.pushplusToken);
                $('#httpTimeout').val(config.httpTimeout);
                $('#batchSaveSize').val(config.batchSaveSize);
                $('#connectionPoolSize').val(config.connectionPoolSize);
                $('#enableNotification').prop('checked', config.enableNotification);
            } catch (error) {
                console.error('Failed to load config:', error);
            }
        }

        // 保存配置
        async function saveConfig() {
            const config = {
                emailTimeout: parseInt($('#emailTimeout').val()),
                emailCheckInterval: parseInt($('#emailCheckInterval').val()),
                registerDelay: parseInt($('#registerDelay').val()),
                retryTimes: parseInt($('#retryTimes').val()),
                concurrency: parseInt($('#concurrency').val()),
                httpTimeout: parseInt($('#httpTimeout').val()),
                batchSaveSize: parseInt($('#batchSaveSize').val()),
                connectionPoolSize: parseInt($('#connectionPoolSize').val()),
                skipApikeyOnRegister: $('#skipApikey').is(':checked'),
                enableNotification: $('#enableNotification').is(':checked'),
                pushplusToken: $('#pushplusToken').val()
            };

            try {
                const response = await fetch('/register/api/config/save', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(config)
                });

                if (response.ok) {
                    currentConfig = config;
                    showToast('配置保存成功', 'success');
                    $('#settingsPanel').slideUp();
                } else {
                    showToast('配置保存失败', 'error');
                }
            } catch (error) {
                showToast('配置保存失败: ' + error, 'error');
            }
        }

        // 高级设置按钮
        $('#settingsBtn').on('click', function() {
            $('#settingsPanel').slideToggle();
        });

        $('#saveSettingsBtn').on('click', function() {
            saveConfig();
        });

        $('#cancelSettingsBtn').on('click', function() {
            $('#settingsPanel').slideUp();
            loadConfig(); // 恢复配置
        });

        $('#startRegisterBtn').on('click', async function() {
            const count = parseInt($('#registerCount').val());

            if (count < 1 || count > 1000) {
                showToast('注册数量必须在1-1000之间', 'error');
                return;
            }

            // 使用当前配置
            if (!currentConfig) {
                await loadConfig();
            }

            // 启动注册任务
            try {
                const response = await fetch('/register/api/register/start', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        count: count,
                        config: currentConfig
                    })
                });

                if (response.ok) {
                    $('#startRegisterBtn').hide();
                    $('#stopRegisterBtn').show();
                    $('#progressContainer').show();
                    $('#statusBadge').text('运行中').removeClass('bg-gray-400').addClass('bg-green-500');
                    showToast('开始注册', 'success');
                } else {
                    showToast('启动失败', 'error');
                }
            } catch (error) {
                showToast('启动失败: ' + error, 'error');
            }
        });

        // 全选复选框
        $('#selectAllCheckbox').on('change', function() {
            if ($(this).is(':checked')) {
                accounts.forEach(acc => selectedAccounts.add(acc.email));
            } else {
                selectedAccounts.clear();
            }
            renderTable();
        });

        // 搜索功能
        $('#searchInput').on('input', function() {
            searchKeyword = $(this).val().toLowerCase();
            loadAccounts();
        });

        // 快速筛选
        $('.quick-filter-btn').on('click', function() {
            currentFilter = $(this).data('filter');
            $('#clearFilterBtn').removeClass('hidden');
            loadAccounts();
        });

        $('#clearFilterBtn').on('click', function() {
            currentFilter = '';
            $(this).addClass('hidden');
            loadAccounts();
        });

        // 取消选择
        $('#cancelSelectionBtn').on('click', function() {
            selectedAccounts.clear();
            renderTable();
        });

        // 批量删除
        $('#batchDeleteBtn').on('click', async function() {
            if (selectedAccounts.size === 0) {
                showToast('请先选择账号', 'error');
                return;
            }

            if (!confirm('确定要删除选中的 ' + selectedAccounts.size + ' 个账号吗？')) return;

            try {
                await fetch('/register/api/accounts/batch-delete', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ emails: Array.from(selectedAccounts) })
                });
                showToast('批量删除成功', 'success');
                selectedAccounts.clear();
                loadAccounts();
            } catch (error) {
                showToast('批量删除失败', 'error');
            }
        });

        // 复制邮箱
        $('#batchCopyEmailsBtn').on('click', function() {
            const emails = Array.from(selectedAccounts).join('\\n');
            navigator.clipboard.writeText(emails);
            showToast('已复制 ' + selectedAccounts.size + ' 个邮箱', 'success');
        });

        // 复制Token
        $('#batchCopyTokensBtn').on('click', function() {
            const tokens = accounts.filter(acc => selectedAccounts.has(acc.email)).map(acc => acc.token).join('\\n');
            navigator.clipboard.writeText(tokens);
            showToast('已复制 ' + selectedAccounts.size + ' 个Token', 'success');
        });

        // 批量补充APIKEY
        $('#batchRefetchApikeyBtn').on('click', async function() {
            const accountsWithoutKey = accounts.filter(acc => !acc.apikey);
            if (accountsWithoutKey.length === 0) {
                showToast('所有账号均已有APIKEY', 'info');
                return;
            }

            if (!confirm('发现 ' + accountsWithoutKey.length + ' 个账号缺少APIKEY，确定要批量获取吗？')) return;

            const emails = accountsWithoutKey.map(acc => acc.email);

            try {
                await fetch('/register/api/batch-refetch-apikey', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ emails: emails })
                });
                showToast('开始批量获取APIKEY，请查看日志...', 'success');
                setTimeout(() => loadAccounts(), 5000);
            } catch (error) {
                showToast('启动失败: ' + error, 'error');
            }
        });

        // 批量检测存活
        $('#batchCheckAccountsBtn').on('click', async function() {
            if (accounts.length === 0) {
                showToast('没有账号可以检测', 'info');
                return;
            }

            const scope = selectedAccounts.size > 0 ? '选中' : '所有';
            const emails = selectedAccounts.size > 0 ? Array.from(selectedAccounts) : accounts.map(acc => acc.email);

            if (!confirm('开始批量检测' + scope + ' ' + emails.length + ' 个账号的存活性？')) return;

            try {
                await fetch('/register/api/batch-check-accounts', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ emails: emails })
                });
                showToast('开始批量检测，请查看实时日志...', 'success');
            } catch (error) {
                showToast('启动失败: ' + error, 'error');
            }
        });

        // 删除失效账号
        $('#deleteInactiveBtn').on('click', async function() {
            if (!confirm('确定要删除所有失效账号吗？此操作不可恢复！')) return;

            try {
                const response = await fetch('/register/api/delete-inactive-accounts', {
                    method: 'POST'
                });
                const data = await response.json();
                showToast('成功删除 ' + data.count + ' 个失效账号', 'success');
                loadAccounts();
            } catch (error) {
                showToast('删除失败: ' + error, 'error');
            }
        });

        // 分页按钮事件
        $('#firstPageBtn').on('click', () => { 
            currentPage = 1; 
            loadAccounts(); 
        });
        
        $('#prevPageBtn').on('click', () => { 
            if (currentPage > 1) { 
                currentPage--; 
                loadAccounts(); 
            } 
        });
        
        $('#nextPageBtn').on('click', () => { 
            if (currentPage < totalPages) { 
                currentPage++; 
                loadAccounts(); 
            } 
        });
        
        $('#lastPageBtn').on('click', () => { 
            currentPage = totalPages; 
            loadAccounts(); 
        });
        
        $('#pageSizeSelect').on('change', function() {
            pageSize = parseInt($(this).val());
            currentPage = 1;
            loadAccounts();
        });

        // 建立SSE连接（页面加载时立即连接）
        function connectSSE() {
            if (eventSource) {
                eventSource.close();
            }

            eventSource = new EventSource('/register/api/register/stream');
            
            eventSource.onopen = function() {
                console.log('SSE连接已建立');
            };

            eventSource.onmessage = function(e) {
                const data = JSON.parse(e.data);
                
                if (data.type === 'connected') {
                    console.log('SSE已连接, 运行状态:', data.isRunning);
                    addLog('✓ 已连接到服务器', 'success');
                    // 根据运行状态更新UI（包括按钮和状态标签）
                    if (data.isRunning) {
                        $('#statusBadge').text('运行中').removeClass('bg-gray-400').addClass('bg-green-500');
                        $('#startRegisterBtn').hide();
                        $('#stopRegisterBtn').show();
                    } else {
                        $('#statusBadge').text('闲置中').removeClass('bg-green-500').addClass('bg-gray-400');
                        $('#startRegisterBtn').show();
                        $('#stopRegisterBtn').hide();
                        $('#progressContainer').hide();
                    }
                } else if (data.type === 'log') {
                    addLog(data.message, data.level, data.link || null);
                } else if (data.type === 'progress') {
                    const percent = Math.round((data.success + data.failed) / data.total * 100);
                    $('#progressBar').css('width', percent + '%');
                    $('#progressText').text((data.success + data.failed) + '/' + data.total + ' (' + percent + '%)');
                    $('#progressContainer').show();
                } else if (data.type === 'complete') {
                    $('#startRegisterBtn').show();
                    $('#stopRegisterBtn').hide();
                    $('#progressContainer').hide();
                    $('#statusBadge').text('闲置中').removeClass('bg-green-500').addClass('bg-gray-400');
                    loadAccounts();
                    showToast('注册任务完成！', 'success');
                } else if (data.type === 'check_complete') {
                    // 批量检测完成，自动刷新列表
                    loadAccounts();
                    showToast('检测完成！正常: ' + data.active + ', 失效: ' + data.inactive, 'success');
                }
            };

            eventSource.onerror = function(e) {
                console.error('SSE连接错误:', e);
                addLog('⚠️ 连接断开，5秒后重连...', 'warning');
                // 5秒后重连
                setTimeout(connectSSE, 5000);
            };
        }

        // 初始加载
        connectSSE();
        loadAccounts();
        loadConfig();
        setInterval(updateStats, 30000);
    </script>
</body>
</html>`
