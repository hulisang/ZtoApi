/**
 * Z.AI账号注册管理V2
 * 登录鉴权/批量注册/实时监控/账号管理/高级配置
 * 存储: Deno KV
 * @author dext7r
 */

import { serve } from "https://deno.land/std@0.208.0/http/server.ts";

// ==================== 配置区域 ====================

const PORT = 8001;  // 端口
const NOTIFY_INTERVAL = 3600;  // 通知间隔秒
const MAX_LOGIN_ATTEMPTS = 5;  // 最大登录失败
const LOGIN_LOCK_DURATION = 900000;  // 锁定15分钟

// 鉴权配置
const AUTH_USERNAME = Deno.env.get("ZAI_USERNAME") || "admin";
const AUTH_PASSWORD = Deno.env.get("ZAI_PASSWORD") || "123456";

// 邮箱域名
const DOMAINS = [
  "chatgptuk.pp.ua", "freemails.pp.ua", "email.gravityengine.cc", "gravityengine.cc",
  "3littlemiracles.com", "almiswelfare.org", "gyan-netra.com", "iraniandsa.org",
  "14club.org.uk", "aard.org.uk", "allumhall.co.uk", "cade.org.uk",
  "caye.org.uk", "cketrust.org", "club106.org.uk", "cok.org.uk",
  "cwetg.co.uk", "goleudy.org.uk", "hhe.org.uk", "hottchurch.org.uk"
];

// ==================== 数据存储 ====================

// KV数据库
let kv: Deno.Kv;

// 初始化KV
async function initKV() {
  try {
    kv = await Deno.openKv();
    console.log("[DEBUG] KV已初始化");
  } catch (error) {
    console.error("❌ KV初始化失败:", error);
    console.error("⚠️ 需要--unstable-kv标志");
    console.error("   运行: deno run --allow-net --allow-env --allow-read --unstable-kv zai_register.ts");
    throw new Error("KV初始化失败");
  }
}

// ==================== 全局状态 ====================

let isRunning = false;  // 运行中
let shouldStop = false;  // 停止标志
const sseClients = new Set<ReadableStreamDefaultController>();  // SSE连接
let stats = { success: 0, failed: 0, startTime: 0, lastNotifyTime: 0 };  // 统计
const logHistory: any[] = [];  // 日志缓存
const MAX_LOG_HISTORY = 500;  // 最大日志数
let logSaveTimer: number | null = null;  // 日志定时器
const LOG_SAVE_INTERVAL = 30000;  // 保存间隔30秒

// 登录失败跟踪
const loginAttempts = new Map<string, { attempts: number; lockedUntil: number }>();

// 批量保存日志(节流)
async function saveLogs(): Promise<void> {
  if (logHistory.length === 0) return;

  try {
    const logKey = ["logs", "recent"];
    const now = Date.now();

    // 保存1小时内日志
    const oneHourAgo = now - 3600000;
    const recentLogs = logHistory
      .filter(log => log.timestamp > oneHourAgo)
      .slice(-200);

    if (recentLogs.length > 0) {
      await kv.set(logKey, recentLogs, { expireIn: 3600000 });  // 1小时过期
    } else {
      await kv.delete(logKey);
    }
  } catch (error) {
    console.error("保存日志失败:", error);
  }
}

// 调度日志保存(防抖)
function scheduleSaveLogs() {
  if (logSaveTimer) {
    clearTimeout(logSaveTimer);
  }

  logSaveTimer = setTimeout(() => {
    saveLogs();
    logSaveTimer = null;
  }, LOG_SAVE_INTERVAL);
}

// 广播消息
function broadcast(data: any) {
  const message = `data: ${JSON.stringify(data)}\n\n`;
  console.log(`📤 broadcast: type=${data.type}, sseClients=${sseClients.size}, message=${message.substring(0, 100)}...`);

  for (const controller of sseClients) {
    try {
      controller.enqueue(new TextEncoder().encode(message));
    } catch (err) {
      console.log(`⚠️ SSE发送失败:`, err);
      sseClients.delete(controller);
    }
  }

  // 保存到内存
  if (data.type === 'log' || data.type === 'start' || data.type === 'complete') {
    logHistory.push({ ...data, timestamp: Date.now() });

    // 清理1小时外日志
    const oneHourAgo = Date.now() - 3600000;
    while (logHistory.length > 0 && logHistory[0].timestamp < oneHourAgo) {
      logHistory.shift();
    }

    // 限制最大数量
    if (logHistory.length > MAX_LOG_HISTORY) {
      logHistory.shift();
    }

    // 调度批量保存
    scheduleSaveLogs();

    // 完成或错误时立即保存
    if (data.type === 'complete' || (data.type === 'log' && data.level === 'error')) {
      saveLogs().catch(() => {});
    }
  }
}

// 生成SessionID
function generateSessionId(): string {
  return crypto.randomUUID();
}

// 获取客户端IP
function getClientIP(req: Request): string {
  // X-Forwarded-For
  const forwarded = req.headers.get("X-Forwarded-For");
  if (forwarded) {
    return forwarded.split(',')[0].trim();
  }

  // X-Real-IP
  const realIP = req.headers.get("X-Real-IP");
  if (realIP) {
    return realIP;
  }

  return "unknown";
}

// 检查IP锁定
function checkIPLocked(ip: string): { locked: boolean; remainingTime?: number } {
  const record = loginAttempts.get(ip);
  if (!record) {
    return { locked: false };
  }

  const now = Date.now();
  if (record.lockedUntil > now) {
    return {
      locked: true,
      remainingTime: Math.ceil((record.lockedUntil - now) / 1000)
    };
  }

  // 过期清除
  loginAttempts.delete(ip);
  return { locked: false };
}

// 记录登录失败
function recordLoginFailure(ip: string): void {
  const record = loginAttempts.get(ip) || { attempts: 0, lockedUntil: 0 };
  record.attempts++;

  if (record.attempts >= MAX_LOGIN_ATTEMPTS) {
    record.lockedUntil = Date.now() + LOGIN_LOCK_DURATION;
    console.log(`🔒 IP ${ip} 已锁定 ${LOGIN_LOCK_DURATION / 60000} 分钟(失败${record.attempts}次)`);
  }

  loginAttempts.set(ip, record);
}

// 清除登录失败
function clearLoginFailure(ip: string): void {
  loginAttempts.delete(ip);
}

// 注册配置
let registerConfig = {
  emailTimeout: 120,  // 邮件超时秒
  emailCheckInterval: 1,  // 轮询间隔秒
  registerDelay: 2000,  // 间隔毫秒
  retryTimes: 3,  // 重试次数
  concurrency: 10,  // 并发1-10
  enableNotification: false,  // 通知默认关
  pushplusToken: "",  // PushPlus Token
};

// ==================== 鉴权相关 ====================

// 检查请求认证
async function checkAuth(req: Request): Promise<{ authenticated: boolean; sessionId?: string }> {
  const cookies = req.headers.get("Cookie") || "";
  const sessionMatch = cookies.match(/sessionId=([^;]+)/);

  if (sessionMatch) {
    const sessionId = sessionMatch[1];
    // KV检查session
    const sessionKey = ["sessions", sessionId];
    const session = await kv.get(sessionKey);

    if (session.value) {
      return { authenticated: true, sessionId };
    }
  }

  return { authenticated: false };
}

// ==================== 工具函数 ====================

// 生成随机邮箱
function createEmail(): string {
  const randomHex = Array.from({ length: 12 }, () =>
    Math.floor(Math.random() * 16).toString(16)
  ).join('');
  const domain = DOMAINS[Math.floor(Math.random() * DOMAINS.length)];
  return `${randomHex}@${domain}`;
}

// 生成随机密码
function createPassword(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';
  return Array.from({ length: 14 }, () =>
    chars[Math.floor(Math.random() * chars.length)]
  ).join('');
}

// PushPlus通知
async function sendNotification(title: string, content: string): Promise<void> {
  // 检查配置
  if (!registerConfig.enableNotification || !registerConfig.pushplusToken) return;

  try {
    await fetch("https://www.pushplus.plus/send", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        token: registerConfig.pushplusToken,
        title,
        content,
        template: "markdown"
      })
    });
  } catch {
    // 忽略错误
  }
}

// 获取验证邮件
async function fetchVerificationEmail(email: string): Promise<string | null> {
  const actualTimeout = registerConfig.emailTimeout;
  const checkInterval = registerConfig.emailCheckInterval;
  const startTime = Date.now();
  const apiUrl = `https://mail.chatgpt.org.uk/api/get-emails?email=${encodeURIComponent(email)}`;

  let attempts = 0;
  let lastReportTime = 0;
  const reportInterval = 10;

  // 格式化时间
  const formatTime = (seconds: number): string => {
    if (seconds < 60) return `${seconds}s`;
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}m${secs}s`;
  };

  while (Date.now() - startTime < actualTimeout * 1000) {
    attempts++;
    try {
      const response = await fetch(apiUrl, { signal: AbortSignal.timeout(10000) });
      const data = await response.json();

      // 每10秒报告
      const elapsed = Math.floor((Date.now() - startTime) / 1000);
      if (elapsed - lastReportTime >= reportInterval && elapsed > 0) {
        const progress = Math.min(Math.floor((elapsed / actualTimeout) * 100), 99);
        const remaining = actualTimeout - elapsed;
        broadcast({
          type: 'log',
          level: 'info',
          message: `  等待邮件[${progress}%] 已用:${formatTime(elapsed)}/剩余:${formatTime(remaining)}(尝试${attempts}次)`
        });
        lastReportTime = elapsed;
      }

      if (data?.emails) {
        for (const emailData of data.emails) {
          if (emailData.from?.toLowerCase().includes("z.ai")) {
            broadcast({ type: 'log', level: 'success', message: `  ✓ 收到邮件(${Math.floor((Date.now() - startTime) / 1000)}s)` });
            return emailData.content || null;
          }
        }
      }
    } catch {
      // 重试
    }
    await new Promise(resolve => setTimeout(resolve, checkInterval * 1000));
  }

  broadcast({ type: 'log', level: 'error', message: `  ✗ 邮件超时(${actualTimeout}s)` });
  return null;
}

function parseVerificationUrl(url: string): { token: string | null; email: string | null; username: string | null } {
  try {
    const urlObj = new URL(url);
    return {
      token: urlObj.searchParams.get('token'),
      email: urlObj.searchParams.get('email'),
      username: urlObj.searchParams.get('username')
    };
  } catch {
    return { token: null, email: null, username: null };
  }
}

// API登录
async function loginToApi(token: string): Promise<string | null> {
  const url = 'https://api.z.ai/api/auth/z/login';
  const headers = {
    'User-Agent': 'Mozilla/5.0',
    'Origin': 'https://z.ai',
    'Referer': 'https://z.ai/',
    'Content-Type': 'application/json'
  };

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify({ token }),
      signal: AbortSignal.timeout(15000)
    });

    const result = await response.json();
    if (result.success && result.code === 200) {
      const accessToken = result.data?.access_token;
      if (accessToken) {
        broadcast({ type: 'log', level: 'success', message: `  ✓ API登录成功` });
        return accessToken;
      }
    }
    broadcast({ type: 'log', level: 'error', message: `  ✗ API登录失败:${JSON.stringify(result)}` });
    return null;
  } catch (error) {
    broadcast({ type: 'log', level: 'error', message: `  ✗ API登录异常:${error}` });
    return null;
  }
}

// 获取客户信息
async function getCustomerInfo(accessToken: string): Promise<{ orgId: string | null; projectId: string | null }> {
  const url = 'https://api.z.ai/api/biz/customer/getCustomerInfo';
  const headers = {
    'Authorization': `Bearer ${accessToken}`,
    'User-Agent': 'Mozilla/5.0',
    'Origin': 'https://z.ai',
    'Referer': 'https://z.ai/'
  };

  try {
    const response = await fetch(url, {
      method: 'GET',
      headers,
      signal: AbortSignal.timeout(20000)
    });

    const result = await response.json();
    if (result.success && result.code === 200) {
      const orgs = result.data?.organizations || [];
      if (orgs.length > 0) {
        const orgId = orgs[0].organizationId;
        const projects = orgs[0].projects || [];
        const projectId = projects.length > 0 ? projects[0].projectId : null;

        if (orgId && projectId) {
          broadcast({ type: 'log', level: 'success', message: `  ✓ 获取组织成功` });
          return { orgId, projectId };
        }
      }
    }
    broadcast({ type: 'log', level: 'error', message: `  ✗ 获取组织失败:${JSON.stringify(result)}` });
    return { orgId: null, projectId: null };
  } catch (error) {
    broadcast({ type: 'log', level: 'error', message: `  ✗ 获取组织异常:${error}` });
    return { orgId: null, projectId: null };
  }
}

// 创建APIKEY
async function createApiKey(accessToken: string, orgId: string, projectId: string): Promise<string | null> {
  const url = `https://api.z.ai/api/biz/v1/organization/${orgId}/projects/${projectId}/api_keys`;
  const headers = {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json',
    'User-Agent': 'Mozilla/5.0',
    'Origin': 'https://z.ai',
    'Referer': 'https://z.ai/'
  };

  try {
    const randomName = 'key_' + Math.random().toString(36).slice(2, 10) + Date.now().toString(36);
    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify({ name: randomName }),
      signal: AbortSignal.timeout(30000)
    });

    const result = await response.json();
    if (result.success && result.code === 200) {
      const apiKeyData = result.data || {};
      const finalKey = `${apiKeyData.apiKey}.${apiKeyData.secretKey}`;
      if (finalKey && finalKey !== 'undefined.undefined') {
        broadcast({ type: 'log', level: 'success', message: `  ✓ APIKEY创建成功` });
        return finalKey;
      }
    }
    broadcast({ type: 'log', level: 'error', message: `  ✗ APIKEY创建失败:${JSON.stringify(result)}` });
    return null;
  } catch (error) {
    broadcast({ type: 'log', level: 'error', message: `  ✗ APIKEY创建异常:${error}` });
    return null;
  }
}

// 清理Token
function cleanToken(token: string): string {
  return token.includes('----') ? token.split('----')[0].trim() : token.trim();
}

// 检查账号有效性
async function checkAccountStatus(token: string): Promise<boolean> {
  try {
    const accessToken = await loginToApi(cleanToken(token));
    return accessToken !== null;
  } catch (error) {
    return false;
  }
}

async function saveAccount(email: string, password: string, token: string, apikey?: string, status: string = 'active'): Promise<boolean> {
  try {
    const timestamp = Date.now();
    const key = ["zai_accounts", timestamp, email];
    await kv.set(key, {
      email,
      password,
      token,
      apikey: apikey || null,
      status: status,
      createdAt: new Date().toISOString()
    });
    return true;
  } catch (error) {
    console.error("❌ KV保存失败:", error);

    const errorMessage = error instanceof Error ? error.message : String(error);
    if (errorMessage.includes("quota is exhausted")) {
      broadcast({
        type: 'log',
        level: 'error',
        message: `❌ KV配额耗尽,保存本地:${email}`
      });
      return false;
    }

    throw error;
  }
}

interface RegisterResult {
  success: boolean;
  account?: { email: string; password: string; token: string; apikey: string | null };
}

async function registerAccount(): Promise<RegisterResult> {
  try {
    const email = createEmail();
    const password = createPassword();
    const name = email.split("@")[0];
    const emailCheckUrl = `https://mail.chatgpt.org.uk/api/get-emails?email=${encodeURIComponent(email)}`;

    broadcast({
      type: 'log',
      level: 'info',
      message: `▶ 开始:${email}`,
      link: { text: '邮箱', url: emailCheckUrl }
    });

    // 1. 注册
    broadcast({ type: 'log', level: 'info', message: `  → 注册...` });
    const signupResponse = await fetch("https://chat.z.ai/api/v1/auths/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, email, password, profile_image_url: "data:image/png;base64,", sso_redirect: null }),
      signal: AbortSignal.timeout(30000)
    });

    if (signupResponse.status !== 200) {
      broadcast({ type: 'log', level: 'error', message: `  ✗ 注册失败:HTTP${signupResponse.status}` });
      stats.failed++;
      return { success: false };
    }

    const signupResult = await signupResponse.json();
    if (!signupResult.success) {
      broadcast({ type: 'log', level: 'error', message: `  ✗ 被拒绝:${JSON.stringify(signupResult)}` });
      stats.failed++;
      return { success: false };
    }

    broadcast({ type: 'log', level: 'success', message: `  ✓ 注册成功` });

    // 2. 获取验证邮件
    broadcast({
      type: 'log',
      level: 'info',
      message: `  → 等待邮件:${email}`,
      link: { text: '打开邮箱', url: emailCheckUrl }
    });
    const emailContent = await fetchVerificationEmail(email);
    if (!emailContent) {
      stats.failed++;
      return { success: false };
    }

    // 3. 提取验证链接
    broadcast({ type: 'log', level: 'info', message: `  → 提取链接...` });

    // 多种匹配
    let verificationUrl = null;

    // 方式1: /auth/verify_email
    let match = emailContent.match(/https:\/\/chat\.z\.ai\/auth\/verify_email\?[^\s<>"']+/);
    if (match) {
      verificationUrl = match[0].replace(/&amp;/g, '&').replace(/&#39;/g, "'");
    }

    // 方式2: /verify_email
    if (!verificationUrl) {
      match = emailContent.match(/https:\/\/chat\.z\.ai\/verify_email\?[^\s<>"']+/);
      if (match) {
        verificationUrl = match[0].replace(/&amp;/g, '&').replace(/&#39;/g, "'");
        broadcast({ type: 'log', level: 'success', message: `  ✓ 旧版路径` });
      }
    }

    // 方式3: HTML编码
    if (!verificationUrl) {
      match = emailContent.match(/https?:\/\/chat\.z\.ai\/(?:auth\/)?verify_email[^"'\s]*/);
      if (match) {
        verificationUrl = match[0].replace(/&amp;/g, '&').replace(/&#39;/g, "'");
        broadcast({ type: 'log', level: 'success', message: `  ✓ HTML解码` });
      }
    }

    // 方式4: JSON格式
    if (!verificationUrl) {
      try {
        const urlMatch = emailContent.match(/"(https?:\/\/[^"]*verify_email[^"]*)"/);
        if (urlMatch) {
          verificationUrl = urlMatch[1].replace(/\\u0026/g, '&').replace(/&amp;/g, '&').replace(/&#39;/g, "'");
          broadcast({ type: 'log', level: 'success', message: `  ✓ JSON格式` });
        }
      } catch (e) {
        // 忽略
      }
    }

    if (!verificationUrl) {
      const preview = emailContent.substring(0, 500).replace(/\n/g, ' ');
      broadcast({ type: 'log', level: 'error', message: `  ✗ 未找到链接:${preview}...` });
      stats.failed++;
      return { success: false };
    }


    const { token, email: emailFromUrl, username } = parseVerificationUrl(verificationUrl);
    if (!token || !emailFromUrl || !username) {
      broadcast({ type: 'log', level: 'error', message: `  ✗ 链接格式错` });
      stats.failed++;
      return { success: false };
    }

    broadcast({ type: 'log', level: 'success', message: `  ✓ 链接已提取` });

    // 4. 完成注册
    broadcast({ type: 'log', level: 'info', message: `  → 验证...` });
    const finishResponse = await fetch("https://chat.z.ai/api/v1/auths/finish_signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: emailFromUrl, password, profile_image_url: "data:image/png;base64,", sso_redirect: null, token, username }),
      signal: AbortSignal.timeout(30000)
    });

    if (finishResponse.status !== 200) {
      broadcast({ type: 'log', level: 'error', message: `  ✗ 验证失败:HTTP${finishResponse.status}` });
      stats.failed++;
      return { success: false };
    }

    const finishResult = await finishResponse.json();
    if (!finishResult.success) {
      broadcast({ type: 'log', level: 'error', message: `  ✗ 验证拒绝:${JSON.stringify(finishResult)}` });
      stats.failed++;
      return { success: false };
    }

    // 5. 获取Token
    const userToken = finishResult.user?.token;
    if (!userToken) {
      broadcast({ type: 'log', level: 'error', message: `  ✗ 无Token` });
      stats.failed++;
      return { success: false };
    }

    broadcast({ type: 'log', level: 'success', message: `  ✓ 获得Token` });

    // 6. API登录
    broadcast({ type: 'log', level: 'info', message: `  → 登录API...` });
    const accessToken = await loginToApi(userToken);
    if (!accessToken) {
      const account = { email, password, token: userToken, apikey: null, createdAt: new Date().toISOString() };
      const saved = await saveAccount(email, password, userToken);

      if (saved) {
        stats.success++;
        broadcast({
          type: 'log',
          level: 'warning',
          message: `⚠️ 成功但API登录失败:${email}(仅Token)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'account_added', account });
      } else {
        stats.success++;
        broadcast({
          type: 'log',
          level: 'warning',
          message: `⚠️ 成功但API登录失败:${email}(仅Token,本地)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'local_account_added', account });
      }

      return { success: true, account };
    }

    // 7. 获取组织
    broadcast({ type: 'log', level: 'info', message: `  → 组织...` });
    const { orgId, projectId } = await getCustomerInfo(accessToken);
    if (!orgId || !projectId) {
      const account = { email, password, token: userToken, apikey: null, createdAt: new Date().toISOString() };
      const saved = await saveAccount(email, password, userToken);

      if (saved) {
        stats.success++;
        broadcast({
          type: 'log',
          level: 'warning',
          message: `⚠️ 成功但组织失败:${email}(仅Token)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'account_added', account });
      } else {
        stats.success++;
        broadcast({
          type: 'log',
          level: 'warning',
          message: `⚠️ 成功但组织失败:${email}(仅Token,本地)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'local_account_added', account });
      }

      return { success: true, account };
    }

    // 8. 创建APIKEY
    broadcast({ type: 'log', level: 'info', message: `  → APIKEY...` });
    const apiKey = await createApiKey(accessToken, orgId, projectId);

    // 9. 保存
    const account = { email, password, token: userToken, apikey: apiKey || null, createdAt: new Date().toISOString() };
    const saved = await saveAccount(email, password, userToken, apiKey || undefined);

    stats.success++;

    if (saved) {
      if (apiKey) {
        broadcast({
          type: 'log',
          level: 'success',
          message: `✅ 完成:${email}(含KEY)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'account_added', account });
      } else {
        broadcast({
          type: 'log',
          level: 'warning',
          message: `⚠️ 成功但KEY失败:${email}(仅Token)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'account_added', account });
      }
    } else {
      if (apiKey) {
        broadcast({
          type: 'log',
          level: 'success',
          message: `✅ 完成:${email}(含KEY,本地)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'local_account_added', account });
      } else {
        broadcast({
          type: 'log',
          level: 'warning',
          message: `⚠️ 成功但KEY失败:${email}(仅Token,本地)`,
          stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed },
          link: { text: '邮箱', url: emailCheckUrl }
        });
        broadcast({ type: 'local_account_added', account });
      }
    }

    return { success: true, account };
  } catch (error: any) {
    const msg = error instanceof Error ? error.message : String(error);
    broadcast({ type: 'log', level: 'error', message: `  ✗ 异常:${msg}` });
    stats.failed++;
    return { success: false };
  }
}

async function batchRegister(count: number): Promise<void> {
  console.log(`🚀 批量注册,count=${count},clients=${sseClients.size}`);

  isRunning = true;
  shouldStop = false;
  stats = { success: 0, failed: 0, startTime: Date.now(), lastNotifyTime: Date.now() };

  console.log(`📡 广播start`);
  broadcast({ type: 'start', config: { count } });
  console.log(`✓ start已广播`);

  const concurrency = registerConfig.concurrency || 1;
  let completed = 0;
  const successAccounts: Array<{ email: string; password: string; token: string; apikey: string | null }> = [];

  // 并发注册
  while (completed < count && !shouldStop) {
    const batchSize = Math.min(concurrency, count - completed);
    const batchPromises: Promise<RegisterResult>[] = [];

    // 创建任务
    for (let i = 0; i < batchSize; i++) {
      const taskIndex = completed + i + 1;
      const progress = Math.floor((taskIndex / count) * 100);
      const elapsed = Math.floor((Date.now() - stats.startTime) / 1000);
      const avgTimePerAccount = completed > 0 ? elapsed / completed : 0;
      const remaining = count - taskIndex;
      const eta = avgTimePerAccount > 0 ? Math.ceil(remaining * avgTimePerAccount) : 0;

      // 格式化时间
      const formatTime = (seconds: number): string => {
        if (seconds < 60) return `${seconds}s`;
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}m${secs}s`;
      };

      broadcast({
        type: 'log',
        level: 'info',
        message: `\n[${taskIndex}/${count}] ━━━━━━━━━━━━━━━━━━━━ [${progress}%] 已用:${formatTime(elapsed)}/预计:${formatTime(eta)}`
      });
      batchPromises.push(registerAccount());
    }

    // 等待完成
    const results = await Promise.allSettled(batchPromises);

    // 收集成功账号
    for (const result of results) {
      if (result.status === 'fulfilled' && result.value.success && result.value.account) {
        successAccounts.push(result.value.account);
      }
    }

    completed += batchSize;

    // 批次延迟
    if (completed < count && !shouldStop) {
      await new Promise(resolve => setTimeout(resolve, registerConfig.registerDelay));
    }
  }

  if (shouldStop) {
    broadcast({ type: 'log', level: 'warning', message: `⚠️ 手动停止,已完成${completed}/${count}` });
  }

  const elapsedTime = (Date.now() - stats.startTime) / 1000;

  broadcast({
    type: 'complete',
    stats: { success: stats.success, failed: stats.failed, total: stats.success + stats.failed, elapsedTime: elapsedTime.toFixed(1) }
  });

  // 总账号数
  let totalAccounts = 0;
  try {
    const entries = kv.list({ prefix: ["zai_accounts"] });
    for await (const _ of entries) {
      totalAccounts++;
    }
  } catch {
    // 忽略
  }

  // 详情(最多10个)
  let accountsDetail = '';
  if (successAccounts.length > 0) {
    accountsDetail += '\n\n### 📋 详情\n';
    const displayCount = Math.min(successAccounts.length, 10);
    for (let i = 0; i < displayCount; i++) {
      const acc = successAccounts[i];
      accountsDetail += `${i + 1}. **${acc.email}**\n`;
      accountsDetail += `   - 密码:\`${acc.password}\`\n`;
      accountsDetail += `   - Token:\`${acc.token.substring(0, 20)}...\`\n`;
      if (acc.apikey) {
        accountsDetail += `   - KEY:\`${acc.apikey.substring(0, 20)}...\`\n`;
      }
    }
    if (successAccounts.length > displayCount) {
      accountsDetail += `\n*还有${successAccounts.length - displayCount}个未显示*\n`;
    }
  }

  // 发送通知
  await sendNotification(
    "✅ Z.AI注册完成",
    `## ✅ Z.AI注册完成

### 📊 结果
- 成功:${stats.success}
- 失败:${stats.failed}
- 本次:${stats.success + stats.failed}
- 总计:${totalAccounts}

### ⏱️ 耗时
- 总:${elapsedTime.toFixed(1)}s (${(elapsedTime / 60).toFixed(1)}min)
- 速度:${((stats.success + stats.failed) / (elapsedTime / 60)).toFixed(1)}/min
- 单个:${stats.success + stats.failed > 0 ? (elapsedTime / (stats.success + stats.failed)).toFixed(1) : 0}s

### 📈 成功率
- 成功:${stats.success + stats.failed > 0 ? ((stats.success / (stats.success + stats.failed)) * 100).toFixed(1) : 0}%
- 失败:${stats.success + stats.failed > 0 ? ((stats.failed / (stats.success + stats.failed)) * 100).toFixed(1) : 0}%${accountsDetail}`
  );

  isRunning = false;
  shouldStop = false;
}

// 登录页面
const LOGIN_PAGE = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>登录 - Z.AI 管理系统</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 min-h-screen flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        <div class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-2">🤖 Z.AI 管理系统</h1>
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
      <p>📦 <a href="https://github.com/dext7r/ZtoApi/tree/main/deno/zai/zai_register.ts" target="_blank" class="text-cyan-600 underline">源码地址 (GitHub)</a> |
      💬 <a href="https://linux.do/t/topic/1009939" target="_blank" class="text-cyan-600 underline">交流讨论</a></p>
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
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                const result = await response.json();

                if (result.success) {
                    document.cookie = 'sessionId=' + result.sessionId + '; path=/; max-age=86400';
                    window.location.href = '/';
                } else {
                    // 显示错误信息
                    let errorText = result.error || '登录失败';

                    // 如果账号被锁定，显示剩余时间
                    if (result.code === 'ACCOUNT_LOCKED' && result.remainingTime) {
                        const minutes = Math.floor(result.remainingTime / 60);
                        const seconds = result.remainingTime % 60;
                        errorText += ' (' + minutes + '分' + seconds + '秒后可重试)';
                    }
                    // 如果有剩余尝试次数，显示提示
                    else if (result.attemptsRemaining !== undefined) {
                        errorText += ' (剩余 ' + result.attemptsRemaining + ' 次尝试机会)';
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
</html>`;

// 主页面
const HTML_PAGE = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>Z.AI 账号管理系统</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://code.jquery.com/jquery-3.7.1.min.js"></script>
    <style>
        @keyframes slideIn {
            from { transform: translateX(100%); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
        @keyframes slideOut {
            from { transform: translateX(0); opacity: 1; }
            to { transform: translateX(100%); opacity: 0; }
        }
        .toast-enter { animation: slideIn 0.3s ease-out; }
        .toast-exit { animation: slideOut 0.3s ease-in; }

        /* 移动端优化 */
        @media (max-width: 768px) {
            .mobile-scroll {
                overflow-x: auto;
                -webkit-overflow-scrolling: touch;
            }

            table {
                font-size: 0.75rem;
            }

            /* 移动端固定Toast位置到底部 */
            #toastContainer {
                left: 0.5rem;
                right: 0.5rem;
                top: auto;
                bottom: 0.5rem;
            }

            /* 优化日志容器高度 */
            #logContainer {
                height: 10rem !important;
            }

            /* 隐藏部分列 */
            .hide-mobile {
                display: none;
            }

            /* 移动端按钮组优化 */
            .btn-group-mobile {
                flex-wrap: wrap;
            }

            /* 移动端可点击单元格 */
            .clickable-cell {
                cursor: pointer;
            }

            .clickable-cell:active {
                opacity: 0.5;
            }
        }

        /* 触摸优化 */
        button, a, input[type="checkbox"] {
            -webkit-tap-highlight-color: rgba(0, 0, 0, 0.1);
        }

        /* 防止双击缩放 */
        * {
            touch-action: manipulation;
        }

        /* 统计卡片选中状态 */
        .stat-card.active {
            ring: 4px;
            ring-color: white;
            box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.5), 0 10px 25px -5px rgba(0, 0, 0, 0.3);
        }

        /* PC端优化 */
        @media (min-width: 769px) {
            /* 表格悬停效果 */
            tbody tr {
                transition: all 0.2s ease;
            }

            tbody tr:hover {
                background-color: #f8fafc;
                transform: translateX(4px);
                box-shadow: -4px 0 0 0 #6366f1;
            }

            /* 操作按钮悬停效果 */
            .action-btn {
                transition: all 0.15s ease;
                position: relative;
            }

            .action-btn:hover {
                transform: translateY(-1px);
            }

            .action-btn::after {
                content: '';
                position: absolute;
                bottom: -2px;
                left: 0;
                right: 0;
                height: 2px;
                background: currentColor;
                transform: scaleX(0);
                transition: transform 0.2s ease;
            }

            .action-btn:hover::after {
                transform: scaleX(1);
            }

            /* 表格单元格内边距优化 */
            td, th {
                padding: 1rem !important;
            }

            /* 代码块样式优化 */
            code {
                font-family: 'Courier New', Consolas, Monaco, monospace;
                letter-spacing: -0.5px;
            }

            /* 可点击单元格样式 */
            .clickable-cell {
                cursor: pointer;
                transition: all 0.15s ease;
                position: relative;
            }

            .clickable-cell:hover {
                opacity: 0.7;
            }

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

        /* 滚动条美化 */
        ::-webkit-scrollbar {
            width: 8px;
            height: 8px;
        }

        ::-webkit-scrollbar-track {
            background: #f1f5f9;
            border-radius: 4px;
        }

        ::-webkit-scrollbar-thumb {
            background: #cbd5e1;
            border-radius: 4px;
        }

        ::-webkit-scrollbar-thumb:hover {
            background: #94a3b8;
        }

        /* 日志链接样式优化 */
        #logContainer a {
            text-decoration: none;
            transition: all 0.2s ease;
        }

        #logContainer a:hover {
            opacity: 0.8;
            transform: translateX(2px);
        }
    </style>
</head>
<body class="bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 min-h-screen p-2 sm:p-4 md:p-8">
    <!-- Toast 容器 -->
    <div id="toastContainer" class="fixed top-4 right-4 z-50 space-y-2"></div>

    <div class="max-w-7xl mx-auto">
        <div class="text-center text-white mb-4 sm:mb-8">
            <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
                <div class="hidden sm:block flex-1"></div>
                <div class="flex-1 text-center">
                    <h1 class="text-2xl sm:text-3xl md:text-5xl font-bold mb-2">🤖 Z.AI 管理系统 V2</h1>
                    <p class="text-sm sm:text-base md:text-xl opacity-90">批量注册 · 数据管理 · 实时监控</p>
          <p class="text-xs sm:text-sm mt-2 opacity-80">📦 <a href="https://github.com/dext7r/ZtoApi/tree/main/deno/zai/zai_register.ts" target="_blank" class="text-cyan-200 underline">源码</a> |
          💬 <a href="https://linux.do/t/topic/1009939" target="_blank" class="text-cyan-200 underline">讨论</a></p>
                </div>
                <div class="w-full sm:w-auto sm:flex-1 sm:flex sm:justify-end">
                    <button id="logoutBtn" class="w-full sm:w-auto px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg text-white font-semibold transition">
                        退出登录
                    </button>
                </div>
            </div>
        </div>

        <!-- 控制面板 + 高级设置 -->
        <div class="bg-white rounded-2xl shadow-2xl p-3 sm:p-6 mb-4 sm:mb-6">
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
                        <input type="number" id="emailTimeout" value="120" min="30" max="300"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">账号间隔 (毫秒)</label>
                        <input type="number" id="registerDelay" value="2000" min="500" max="10000" step="500"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">邮件轮询间隔（秒）</label>
                        <input type="number" id="emailCheckInterval" value="1" min="0.5" max="10" step="0.5"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                        <p class="text-xs text-gray-500 mt-1">建议：0.5-2秒，过小可能触发限流</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">并发数</label>
                        <input type="number" id="concurrency" value="1" min="1" max="10"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                        <p class="text-xs text-gray-500 mt-1">同时注册的账号数量，建议3-5</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">API 重试次数</label>
                        <input type="number" id="retryTimes" value="3" min="1" max="10"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">PushPlus Token</label>
                        <input type="text" id="pushplusToken" value="" placeholder="留空则不发送通知"
                            class="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                    </div>
                    <div class="flex items-center md:col-span-2">
                        <input type="checkbox" id="enableNotification" checked class="w-5 h-5 text-indigo-600 rounded">
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
                <input type="number" id="registerCount" value="5" min="1" max="100"
                    class="flex-1 px-4 py-3 text-base border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">
                <button id="startRegisterBtn"
                    class="w-full sm:w-auto px-6 sm:px-8 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all disabled:opacity-60 disabled:cursor-not-allowed text-base">
                    开始注册
                </button>
                <button id="stopRegisterBtn" style="display: none;"
                    class="w-full sm:w-auto px-6 sm:px-8 py-3 bg-gradient-to-r from-red-500 to-pink-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all text-base">
                    停止注册
                </button>
            </div>

            <!-- 进度条 -->
            <div id="progressContainer" style="display: none;" class="mb-4">
                <div class="flex justify-between text-sm text-gray-600 mb-2">
                    <span>注册进度</span>
                    <span id="progressText">0/0 (0%)</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-4 overflow-hidden">
                    <div id="progressBar" class="h-full bg-gradient-to-r from-indigo-500 to-purple-600 rounded-full transition-all duration-300 flex items-center justify-center">
                        <span id="progressPercent" class="text-xs text-white font-semibold"></span>
                    </div>
                </div>
                <div class="flex justify-between text-xs text-gray-500 mt-1">
                    <span id="progressSpeed">速度: 0/分钟</span>
                    <span id="progressETA">预计剩余: --</span>
                </div>
            </div>
        </div>

        <!-- 统计面板 -->
        <div class="bg-white rounded-2xl shadow-2xl p-3 sm:p-6 mb-4 sm:mb-6">
            <h2 class="text-xl sm:text-2xl font-bold text-gray-800 mb-3 sm:mb-4">统计信息 <span class="text-sm text-gray-500 font-normal">(点击切换显示)</span></h2>
            <div class="grid grid-cols-2 md:grid-cols-5 gap-2 sm:gap-4">
                <div id="totalAccountsCard" class="stat-card bg-gradient-to-br from-green-400 to-emerald-500 rounded-xl p-3 sm:p-4 text-center text-white cursor-pointer transform transition-all hover:scale-105 active:scale-95">
                    <div class="text-xs sm:text-sm opacity-90 mb-1">总账号</div>
                    <div class="text-2xl sm:text-3xl font-bold" id="totalAccounts">0</div>
                </div>
                <div id="localAccountsCard" class="stat-card bg-gradient-to-br from-cyan-400 to-teal-500 rounded-xl p-3 sm:p-4 text-center text-white cursor-pointer transform transition-all hover:scale-105 active:scale-95">
                    <div class="text-xs sm:text-sm opacity-90 mb-1">本地账号</div>
                    <div class="text-2xl sm:text-3xl font-bold" id="localAccountsCount">0</div>
                </div>
                <div class="bg-gradient-to-br from-blue-400 to-indigo-500 rounded-xl p-3 sm:p-4 text-center text-white">
                    <div class="text-xs sm:text-sm opacity-90 mb-1">本次成功</div>
                    <div class="text-2xl sm:text-3xl font-bold" id="successCount">0</div>
                </div>
                <div class="bg-gradient-to-br from-red-400 to-pink-500 rounded-xl p-3 sm:p-4 text-center text-white">
                    <div class="text-xs sm:text-sm opacity-90 mb-1">本次失败</div>
                    <div class="text-2xl sm:text-3xl font-bold" id="failedCount">0</div>
                </div>
                <div class="bg-gradient-to-br from-purple-400 to-fuchsia-500 rounded-xl p-3 sm:p-4 text-center text-white">
                    <div class="text-xs sm:text-sm opacity-90 mb-1">耗时</div>
                    <div class="text-2xl sm:text-3xl font-bold" id="timeValue">0s</div>
                </div>
            </div>
        </div>

        <!-- 账号列表 -->
        <div class="bg-white rounded-2xl shadow-2xl p-3 sm:p-6 mb-4 sm:mb-6">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4">
                <h2 class="text-xl sm:text-2xl font-bold text-gray-800">账号列表</h2>
                <div class="flex flex-wrap gap-2 w-full sm:w-auto">
                    <input type="text" id="searchInput" placeholder="搜索邮箱..."
                        class="flex-1 sm:flex-none px-3 sm:px-4 py-2 text-sm sm:text-base border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:ring focus:ring-indigo-200 transition">

                    <!-- 服务端操作 -->
                    <input type="file" id="importFileInput" accept=".txt" style="display: none;">
                    <button id="importBtn"
                        class="local-operation-btn flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-purple-500 to-violet-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        📥 导入到服务器
                    </button>
                    <button id="exportBtn"
                        class="local-operation-btn flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        📤 导出服务器
                    </button>

                    <!-- 本地操作 -->
                    <input type="file" id="importLocalFileInput" accept=".txt" style="display: none;">
                    <button id="importLocalBtn"
                        class="local-operation-btn flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-cyan-500 to-teal-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        💾 导入本地
                    </button>
                    <button id="exportLocalBtn"
                        class="local-operation-btn flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-orange-500 to-amber-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        📦 导出本地
                    </button>
                    <button id="syncToServerBtn"
                        class="local-operation-btn flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        🔄 同步到服务器
                    </button>

                    <!-- APIKEY批量操作 -->
                    <button id="batchRefetchApikeyBtn"
                        class="flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-pink-500 to-rose-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        🔑 批量补充APIKEY
                    </button>

                    <!-- 存活性检测 -->
                    <button id="batchCheckAccountsBtn"
                        class="flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-yellow-500 to-orange-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        🔍 批量检测存活
                    </button>

                    <button id="deleteInactiveBtn"
                        class="flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-red-500 to-pink-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        🗑️ 删除失效账号
                    </button>

                    <button id="refreshBtn"
                        class="flex-1 sm:flex-none px-3 sm:px-4 py-2 bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-xs sm:text-sm whitespace-nowrap">
                        🔃 刷新
                    </button>
                </div>
            </div>
            <div class="overflow-x-auto mobile-scroll">
                <table class="w-full min-w-[640px]">
                    <thead>
                        <tr class="bg-gray-50 text-left">
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700">序号</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700">邮箱</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700 hide-mobile">密码</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700 hide-mobile">Token</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700 hide-mobile">APIKEY</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700 hide-mobile">创建时间</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700">状态</th>
                            <th class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm font-semibold text-gray-700">操作</th>
                        </tr>
                    </thead>
                    <tbody id="accountTableBody" class="divide-y divide-gray-200">
                        <tr>
                            <td colspan="7" class="px-4 py-8 text-center text-gray-400">暂无数据</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <!-- 分页控件 -->
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mt-4 px-2 sm:px-4">
                <div class="text-xs sm:text-sm text-gray-600">
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
        <div class="bg-white rounded-2xl shadow-2xl p-3 sm:p-6">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4">
                <h2 class="text-xl sm:text-2xl font-bold text-gray-800">实时日志</h2>
                <button id="clearLogBtn"
                    class="w-full sm:w-auto px-3 sm:px-4 py-2 bg-gradient-to-r from-gray-500 to-gray-600 text-white font-semibold rounded-lg shadow hover:shadow-lg transition text-sm sm:text-base">
                    清空日志
                </button>
            </div>
            <div id="logContainer" class="bg-gray-900 rounded-lg p-3 sm:p-4 h-40 sm:h-64 overflow-y-auto font-mono text-xs sm:text-sm">
                <div class="text-blue-400">等待任务启动...</div>
            </div>
        </div>
    </div>

    <script>
        let accounts = [];
        let filteredAccounts = [];
        let isRunning = false;
        let currentPage = 1;
        let pageSize = 20;
        let taskStartTime = 0;
        let totalTaskCount = 0;
        let filterMode = 'all'; // 'all' 或 'local'

        const $statusBadge = $('#statusBadge');
        const $startRegisterBtn = $('#startRegisterBtn');
        const $stopRegisterBtn = $('#stopRegisterBtn');
        const $logContainer = $('#logContainer');
        const $totalAccounts = $('#totalAccounts');
        const $successCount = $('#successCount');
        const $failedCount = $('#failedCount');
        const $timeValue = $('#timeValue');
        const $accountTableBody = $('#accountTableBody');
        const $searchInput = $('#searchInput');
        const $progressContainer = $('#progressContainer');
        const $progressBar = $('#progressBar');
        const $progressText = $('#progressText');
        const $progressPercent = $('#progressPercent');
        const $progressSpeed = $('#progressSpeed');
        const $progressETA = $('#progressETA');

        // 更新进度条
        function updateProgress(current, total, success, failed) {
            const completed = success + failed;
            const percent = total > 0 ? Math.round((completed / total) * 100) : 0;

            $progressBar.css('width', percent + '%');
            $progressPercent.text(percent + '%');
            $progressText.text(completed + '/' + total + ' (' + percent + '%)');

            // 计算速度和预计剩余时间
            if (taskStartTime > 0 && completed > 0) {
                const elapsed = (Date.now() - taskStartTime) / 1000 / 60; // 分钟
                const speed = completed / elapsed;
                const remaining = total - completed;
                const eta = remaining / speed;

                $progressSpeed.text('速度: ' + speed.toFixed(1) + '/分钟');

                if (eta < 1) {
                    $progressETA.text('预计剩余: <1分钟');
                } else if (eta < 60) {
                    $progressETA.text('预计剩余: ' + Math.ceil(eta) + '分钟');
                } else {
                    const hours = Math.floor(eta / 60);
                    const mins = Math.ceil(eta % 60);
                    $progressETA.text('预计剩余: ' + hours + '小时' + mins + '分钟');
                }
            }
        }

        // Toast 消息提示
        function showToast(message, type = 'info') {
            const colors = {
                success: 'bg-green-500',
                error: 'bg-red-500',
                warning: 'bg-yellow-500',
                info: 'bg-blue-500'
            };
            const icons = {
                success: '✓',
                error: '✗',
                warning: '⚠',
                info: 'ℹ'
            };

            const $toast = $('<div>', {
                class: 'toast-enter ' + colors[type] + ' text-white px-6 py-3 rounded-lg shadow-lg flex items-center gap-2 min-w-[300px]',
                html: '<span class="text-xl">' + icons[type] + '</span><span>' + message + '</span>'
            });

            $('#toastContainer').append($toast);

            setTimeout(() => {
                $toast.removeClass('toast-enter').addClass('toast-exit');
                setTimeout(() => $toast.remove(), 300);
            }, 3000);
        }

        function addLog(message, level = 'info', link = null) {
            const colors = { success: 'text-green-400', error: 'text-red-400', warning: 'text-yellow-400', info: 'text-blue-400' };
            const time = new Date().toLocaleTimeString('zh-CN');

            let html = '<span class="text-gray-500">[' + time + ']</span> ' + message;

            // 添加链接（优化样式，更醒目）
            if (link && link.url) {
                html += ' <a href="' + link.url + '" target="_blank" class="inline-flex items-center ml-2 px-2 py-0.5 bg-cyan-600/20 text-cyan-400 hover:text-cyan-300 hover:bg-cyan-600/30 rounded border border-cyan-500/30 text-xs font-medium transition">' +
                    '<svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>' +
                    (link.text || '查看') +
                    '</a>';
            }

            const $log = $('<div>', {
                class: colors[level] + ' mb-1',
                html: html
            });

            $logContainer.append($log);
            $logContainer[0].scrollTop = $logContainer[0].scrollHeight;
            if ($logContainer.children().length > 200) $logContainer.children().first().remove();
        }

        function updateStatus(running) {
            isRunning = running;
            if (running) {
                $statusBadge.text('运行中').removeClass('bg-gray-400').addClass('bg-green-500');
                $startRegisterBtn.hide();
                $stopRegisterBtn.show();
            } else {
                $statusBadge.text('闲置中').removeClass('bg-green-500').addClass('bg-gray-400');
                $startRegisterBtn.show();
                $stopRegisterBtn.hide();
            }
        }

        function renderTable() {
            // 根据过滤模式应用过滤
            let displayData = filteredAccounts;
            if (filterMode === 'local') {
                displayData = filteredAccounts.filter(acc => acc.source === 'local');
            }

            const totalPages = Math.ceil(displayData.length / pageSize);
            const startIndex = (currentPage - 1) * pageSize;
            const endIndex = startIndex + pageSize;
            const pageData = displayData.slice(startIndex, endIndex);

            if (pageData.length === 0) {
                $accountTableBody.html('<tr><td colspan="8" class="px-4 py-8 text-center text-gray-400">暂无数据</td></tr>');
            } else {
                const rows = pageData.map((acc, idx) => {
                    const rowId = 'row-' + (startIndex + idx);
                    // 处理APIKEY显示
                    const apikeyDisplay = acc.apikey ?
                        '<code class="bg-indigo-50 text-indigo-700 px-2 py-1 rounded text-xs font-mono">' + acc.apikey.substring(0, 20) + '...</code>' :
                        '<span class="text-gray-400 text-xs italic">未生成</span>';

                    // 处理状态显示
                    const status = acc.status || 'active';
                    const statusDisplay = status === 'active' ?
                        '<span class="px-2 py-1 bg-green-100 text-green-700 rounded-full text-xs font-medium">✓ 正常</span>' :
                        '<span class="px-2 py-1 bg-red-100 text-red-700 rounded-full text-xs font-medium">✗ 失效</span>';

                    return '<tr class="group" id="' + rowId + '">' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm text-gray-700 font-medium">' + (startIndex + idx + 1) + '</td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm text-gray-700 truncate max-w-[200px] clickable-cell" title="点击复制: ' + acc.email + '" data-copy="' + acc.email + '">' + acc.email + '</td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm text-gray-700 hide-mobile clickable-cell" title="点击复制密码" data-copy="' + acc.password + '"><code class="bg-blue-50 text-blue-700 px-2 py-1 rounded text-xs font-mono">' + acc.password + '</code></td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm text-gray-700 hide-mobile clickable-cell" title="点击复制Token" data-copy="' + acc.token + '"><code class="bg-green-50 text-green-700 px-2 py-1 rounded text-xs font-mono">' + acc.token.substring(0, 20) + '...</code></td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm text-gray-700 hide-mobile' + (acc.apikey ? ' clickable-cell' : '') + '"' + (acc.apikey ? ' title="点击复制APIKEY" data-copy="' + acc.apikey + '"' : '') + '>' + apikeyDisplay + '</td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-xs sm:text-sm text-gray-700 hide-mobile">' + new Date(acc.createdAt).toLocaleString('zh-CN') + '</td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3 text-center">' + statusDisplay + '</td>' +
                        '<td class="px-2 sm:px-4 py-2 sm:py-3"><div class="flex gap-1 sm:gap-2 flex-wrap">' +
                            '<button class="copy-full-btn action-btn text-indigo-600 hover:text-indigo-800 text-xs sm:text-sm font-medium whitespace-nowrap" ' +
                            'data-email="' + acc.email + '" ' +
                            'data-password="' + acc.password + '" ' +
                            'data-token="' + acc.token + '" ' +
                            'data-apikey="' + (acc.apikey || '') + '" ' +
                            'data-createdat="' + acc.createdAt + '">复制全部</button>' +
                            (!acc.apikey ? '<button class="refetch-apikey-btn action-btn text-green-600 hover:text-green-800 text-xs sm:text-sm font-medium whitespace-nowrap" ' +
                            'data-email="' + acc.email + '" ' +
                            'data-token="' + acc.token + '">🔑 获取APIKEY</button>' : '') +
                        '</div></td>' +
                    '</tr>';
                });
                $accountTableBody.html(rows.join(''));

                // 绑定单元格点击复制事件
                $('.clickable-cell').on('click', function() {
                    const copyText = $(this).data('copy');
                    if (copyText) {
                        navigator.clipboard.writeText(copyText);
                        const cellContent = $(this).text().trim();
                        const displayText = cellContent.length > 30 ? cellContent.substring(0, 30) + '...' : cellContent;
                        showToast('已复制: ' + displayText, 'success');
                    }
                });

                // 绑定"复制全部"按钮事件
                $('.copy-full-btn').on('click', function() {
                    const email = $(this).data('email');
                    const password = $(this).data('password');
                    const token = $(this).data('token');
                    const apikey = $(this).data('apikey');
                    const createdAt = $(this).data('createdat');

                    // 构建完整的账号信息
                    let fullInfo = '邮箱: ' + email + '\\n密码: ' + password + '\\n';
                    fullInfo += 'Token: ' + token + '\\n';
                    if (apikey) {
                        fullInfo += 'APIKEY: ' + apikey + '\\n';
                    }
                    fullInfo += '创建时间: ' + new Date(createdAt).toLocaleString('zh-CN');

                    navigator.clipboard.writeText(fullInfo);
                    showToast('已复制完整账号信息', 'success');
                });

                // 绑定"获取APIKEY"按钮事件
                $('.refetch-apikey-btn').on('click', async function() {
                    const email = $(this).data('email');
                    const token = $(this).data('token');
                    $(this).prop('disabled', true).text('获取中...');
                    await refetchSingleApikey(email, token);
                    // loadAccounts会重新渲染表格，按钮会自动恢复
                });
            }

            // 更新分页控件
            updatePagination(displayData.length, totalPages);

            // 控制本地操作按钮的显示
            if (filterMode === 'local') {
                $('.local-operation-btn').show();
            } else {
                $('.local-operation-btn').hide();
            }
        }

        function updatePagination(totalItems, totalPages) {
            $('#totalItems').text(totalItems);

            // 更新按钮状态
            $('#firstPageBtn, #prevPageBtn').prop('disabled', currentPage === 1);
            $('#nextPageBtn, #lastPageBtn').prop('disabled', currentPage === totalPages || totalPages === 0);

            // 渲染页码
            const $pageNumbers = $('#pageNumbers');
            $pageNumbers.empty();

            if (totalPages <= 7) {
                for (let i = 1; i <= totalPages; i++) {
                    addPageButton(i, $pageNumbers);
                }
            } else {
                addPageButton(1, $pageNumbers);
                if (currentPage > 3) $pageNumbers.append('<span class="px-2">...</span>');

                let start = Math.max(2, currentPage - 1);
                let end = Math.min(totalPages - 1, currentPage + 1);

                for (let i = start; i <= end; i++) {
                    addPageButton(i, $pageNumbers);
                }

                if (currentPage < totalPages - 2) $pageNumbers.append('<span class="px-2">...</span>');
                addPageButton(totalPages, $pageNumbers);
            }
        }

        function addPageButton(page, container) {
            const isActive = page === currentPage;
            const $btn = $('<button>', {
                text: page,
                class: 'px-3 py-1 border rounded ' + (isActive ? 'bg-indigo-600 text-white border-indigo-600' : 'border-gray-300 hover:bg-gray-100'),
                click: () => {
                    currentPage = page;
                    renderTable();
                }
            });
            container.append($btn);
        }

        async function loadAccounts() {
            const response = await fetch('/api/accounts');
            accounts = await response.json();
            filteredAccounts = accounts;
            $totalAccounts.text(accounts.length);
            currentPage = 1;

            // 加载并合并本地账号
            await loadLocalAccounts();
        }

        $searchInput.on('input', function() {
            const keyword = $(this).val().toLowerCase();
            filteredAccounts = accounts.filter(acc => acc.email.toLowerCase().includes(keyword));
            currentPage = 1;
            renderTable();
        });

        // 分页按钮事件
        $('#firstPageBtn').on('click', () => { currentPage = 1; renderTable(); });
        $('#prevPageBtn').on('click', () => { if (currentPage > 1) { currentPage--; renderTable(); } });
        $('#nextPageBtn').on('click', () => { const totalPages = Math.ceil(filteredAccounts.length / pageSize); if (currentPage < totalPages) { currentPage++; renderTable(); } });
        $('#lastPageBtn').on('click', () => { currentPage = Math.ceil(filteredAccounts.length / pageSize); renderTable(); });
        $('#pageSizeSelect').on('change', function() {
            pageSize = parseInt($(this).val());
            currentPage = 1;
            renderTable();
        });

        async function loadSettings() {
            try {
                const response = await fetch('/api/config');
                if (!response.ok) {
                    if (response.status === 302) {
                        window.location.href = '/login';
                        return;
                    }
                    throw new Error('HTTP ' + response.status);
                }
                const config = await response.json();
                $('#emailTimeout').val(config.emailTimeout);
                $('#emailCheckInterval').val(config.emailCheckInterval || 1);
                $('#registerDelay').val(config.registerDelay);
                $('#retryTimes').val(config.retryTimes);
                $('#concurrency').val(config.concurrency || 1);
                $('#enableNotification').prop('checked', config.enableNotification);
                $('#pushplusToken').val(config.pushplusToken || '');
            } catch (error) {
                console.error('加载配置失败:', error);
                showToast('加载配置失败', 'error');
            }
        }

        $('#refreshBtn').on('click', loadAccounts);

        // 统计卡片点击事件 - 切换过滤模式
        $('#totalAccountsCard').on('click', function() {
            filterMode = 'all';
            $('.stat-card').removeClass('active');
            $(this).addClass('active');
            currentPage = 1;
            renderTable();
        });

        $('#localAccountsCard').on('click', function() {
            filterMode = 'local';
            $('.stat-card').removeClass('active');
            $(this).addClass('active');
            currentPage = 1;
            renderTable();
        });

        // 默认选中总账号卡片
        $('#totalAccountsCard').addClass('active');

        $('#clearLogBtn').on('click', function() {
            $logContainer.html('<div class="text-gray-500">日志已清空</div>');
            addLog('✓ 日志已清空', 'success');
        });

        $('#settingsBtn').on('click', function() {
            $('#settingsPanel').slideToggle();
        });

        $('#cancelSettingsBtn').on('click', function() {
            $('#settingsPanel').slideUp();
        });

        $('#saveSettingsBtn').on('click', async function() {
            try {
                const config = {
                    emailTimeout: parseInt($('#emailTimeout').val()),
                    emailCheckInterval: parseFloat($('#emailCheckInterval').val()),
                    registerDelay: parseInt($('#registerDelay').val()),
                    retryTimes: parseInt($('#retryTimes').val()),
                    concurrency: parseInt($('#concurrency').val()),
                    enableNotification: $('#enableNotification').is(':checked'),
                    pushplusToken: $('#pushplusToken').val().trim()
                };

                const response = await fetch('/api/config', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(config)
                });

                if (!response.ok) {
                    if (response.status === 302) {
                        window.location.href = '/login';
                        return;
                    }
                    throw new Error('HTTP ' + response.status);
                }

                const result = await response.json();
                if (result.success) {
                    showToast('设置已保存', 'success');
                    $('#settingsPanel').slideUp();
                } else {
                    showToast('保存失败: ' + (result.error || '未知错误'), 'error');
                }
            } catch (error) {
                console.error('保存配置失败:', error);
                showToast('保存失败: ' + error.message, 'error');
            }
        });

        $('#logoutBtn').on('click', async function() {
            if (confirm('确定要退出登录吗？')) {
                await fetch('/api/logout', { method: 'POST' });
                document.cookie = 'sessionId=; path=/; max-age=0';
                window.location.href = '/login';
            }
        });

        $('#exportBtn').on('click', async function() {
            try {
                const response = await fetch('/api/export');
                const blob = await response.blob();
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'zai_accounts_' + Date.now() + '.txt';
                a.click();
                showToast('导出成功！', 'success');
            } catch (error) {
                showToast('导出失败: ' + error.message, 'error');
            }
        });

        $('#importBtn').on('click', function() {
            $('#importFileInput').click();
        });

        $('#importFileInput').on('change', async function(e) {
            const file = e.target.files[0];
            if (!file) return;

            try {
                showToast('开始导入，请稍候...', 'info');
                const text = await file.text();
                const lines = text.split('\\n').filter(line => line.trim());

                // 准备批量数据
                const importData = [];
                const emailSet = new Set();

                for (const line of lines) {
                    const parts = line.split('----');
                    let email, password, token, apikey;

                    if (parts.length >= 4) {
                        // 四字段格式：账号----密码----Token----APIKEY
                        email = parts[0].trim();
                        password = parts[1].trim();
                        token = parts[2].trim();
                        apikey = parts[3].trim() || null;
                    } else if (parts.length === 3) {
                        // 三字段格式（旧格式）：账号----密码----Token
                        email = parts[0].trim();
                        password = parts[1].trim();
                        token = parts[2].trim();
                        apikey = null;
                    } else {
                        continue;
                    }

                    // 去重检查
                    if (!emailSet.has(email)) {
                        emailSet.add(email);
                        importData.push({ email, password, token, apikey });
                    }
                }

                // 批量导入
                const response = await fetch('/api/import-batch', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ accounts: importData })
                });

                const result = await response.json();
                if (result.success) {
                    showToast('导入完成！成功: ' + result.imported + ', 跳过重复: ' + result.skipped, 'success');
                    await loadAccounts();
                } else {
                    showToast('导入失败: ' + result.error, 'error');
                }

                $(this).val('');
            } catch (error) {
                showToast('导入失败: ' + error.message, 'error');
            }
        });

        // 本地存储操作事件
        $('#exportLocalBtn').on('click', exportLocalAccounts);

        $('#importLocalBtn').on('click', function() {
            $('#importLocalFileInput').click();
        });

        $('#importLocalFileInput').on('change', async function(e) {
            const file = e.target.files[0];
            if (!file) return;
            await importToLocal(file);
            $(this).val(''); // 清空input，允许重复选择同一文件
        });

        $('#syncToServerBtn').on('click', syncLocalToServer);

        $('#batchRefetchApikeyBtn').on('click', batchRefetchApikey);

        $('#batchCheckAccountsBtn').on('click', batchCheckAccounts);

        $('#deleteInactiveBtn').on('click', deleteInactiveAccounts);

        $startRegisterBtn.on('click', async function() {
            try {
                const count = parseInt($('#registerCount').val());
                if (!count || count < 1) {
                    alert('请输入有效数量');
                    return;
                }

                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ count })
                });

                const result = await response.json();

                if (!response.ok) {
                    if (response.status === 302) {
                        window.location.href = '/login';
                        return;
                    }

                    // 显示详细错误信息
                    if (result.isRunning) {
                        const msg = result.error + '\\n\\n' +
                            '当前进度：' + result.stats.success + ' 成功 / ' + result.stats.failed + ' 失败 / ' + result.stats.total + ' 已完成';
                        showToast(msg, 'warning');
                        addLog('⚠️ ' + result.error, 'warning');
                    } else {
                        showToast(result.error || '启动失败', 'error');
                        addLog('✗ ' + (result.error || '启动失败'), 'error');
                    }
                    return;
                }

                if (!result.success) {
                    addLog('✗ ' + (result.error || '启动失败'), 'error');
                }
            } catch (error) {
                console.error('启动注册失败:', error);
                addLog('✗ 启动失败: ' + error.message, 'error');
                showToast('启动失败: ' + error.message, 'error');
            }
        });

        $stopRegisterBtn.on('click', async function() {
            if (confirm('确定要停止当前注册任务吗？')) {
                const response = await fetch('/api/stop', { method: 'POST' });
                const result = await response.json();
                if (result.success) {
                    addLog('⚠️ 已发送停止信号...', 'warning');
                }
            }
        });

        // ========== IndexedDB 操作库 ==========
        const DB_NAME = 'ZaiAccountsDB';
        const DB_VERSION = 1;
        const STORE_NAME = 'accounts';

        let db = null;

        // 初始化 IndexedDB
        async function initIndexedDB() {
            return new Promise((resolve, reject) => {
                const request = indexedDB.open(DB_NAME, DB_VERSION);

                request.onerror = () => {
                    console.error('IndexedDB初始化失败:', request.error);
                    addLog('⚠️ 本地存储初始化失败', 'warning');
                    reject(request.error);
                };

                request.onsuccess = () => {
                    db = request.result;
                    console.log('✓ IndexedDB初始化成功');
                    // loadAccounts() 会调用 loadLocalAccounts() 合并本地账号
                    resolve(db);
                };

                request.onupgradeneeded = (event) => {
                    const db = event.target.result;

                    if (!db.objectStoreNames.contains(STORE_NAME)) {
                        const store = db.createObjectStore(STORE_NAME, { keyPath: 'id', autoIncrement: true });
                        store.createIndex('email', 'email', { unique: true });
                        store.createIndex('source', 'source', { unique: false });
                        store.createIndex('createdAt', 'createdAt', { unique: false });
                        console.log('✓ 创建IndexedDB表结构');
                    }
                };
            });
        }

        // 保存账号到 IndexedDB
        async function saveToLocal(account) {
            if (!db) {
                console.warn('IndexedDB未初始化');
                return false;
            }

            return new Promise((resolve, reject) => {
                const transaction = db.transaction([STORE_NAME], 'readwrite');
                const store = transaction.objectStore(STORE_NAME);

                const accountData = {
                    email: account.email,
                    password: account.password,
                    token: account.token,
                    apikey: account.apikey || null,
                    source: account.source || 'local', // local/kv/synced
                    createdAt: account.createdAt || new Date().toISOString()
                };

                const request = store.add(accountData);

                request.onsuccess = () => {
                    console.log('✓ 账号已保存到本地:', account.email);
                    resolve(true);
                };

                request.onerror = () => {
                    if (request.error.name === 'ConstraintError') {
                        console.log('⚠️ 账号已存在，跳过:', account.email);
                        resolve(false);
                    } else {
                        console.error('保存失败:', request.error);
                        reject(request.error);
                    }
                };
            });
        }

        // 获取所有本地账号
        async function getAllLocalAccounts() {
            if (!db) return [];

            return new Promise((resolve, reject) => {
                const transaction = db.transaction([STORE_NAME], 'readonly');
                const store = transaction.objectStore(STORE_NAME);
                const request = store.getAll();

                request.onsuccess = () => resolve(request.result);
                request.onerror = () => reject(request.error);
            });
        }

        // 加载本地账号到界面
        async function loadLocalAccounts() {
            try {
                const localAccounts = await getAllLocalAccounts();
                console.log(\`✓ 加载了 \${localAccounts.length} 个本地账号\`);

                // 合并服务端账号和本地账号到accounts数组
                // 使用Map去重（以email为key）
                const accountMap = new Map();

                // 先添加服务器账号
                accounts.forEach(acc => {
                    accountMap.set(acc.email, acc);
                });

                // 再添加本地账号（如果email不存在）
                localAccounts.forEach(acc => {
                    if (!accountMap.has(acc.email)) {
                        // 格式化为统一的账号格式
                        accountMap.set(acc.email, {
                            email: acc.email,
                            password: acc.password,
                            token: acc.token,
                            apikey: acc.apikey || null,
                            source: acc.source || 'local',
                            createdAt: acc.createdAt
                        });
                    }
                });

                // 更新accounts和filteredAccounts
                accounts = Array.from(accountMap.values());
                filteredAccounts = accounts;

                // 更新统计
                $totalAccounts.text(accounts.length);
                $('#localAccountsCount').text(accounts.filter(a => a.source === 'local').length);

                // 重新渲染表格（保持当前过滤模式）
                renderTable();
            } catch (error) {
                console.error('加载本地账号失败:', error);
            }
        }

        // 导出本地账号为TXT
        async function exportLocalAccounts() {
            try {
                const localAccounts = await getAllLocalAccounts();
                if (localAccounts.length === 0) {
                    showToast('没有本地账号可导出', 'warning');
                    return;
                }

                const content = localAccounts.map(acc =>
                    \`\${acc.email}----\${acc.password}----\${acc.token}----\${acc.apikey || ''}\`
                ).join('\\n');

                const blob = new Blob([content], { type: 'text/plain' });
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = \`zai_local_accounts_\${Date.now()}.txt\`;
                a.click();
                URL.revokeObjectURL(url);

                showToast(\`已导出 \${localAccounts.length} 个本地账号\`, 'success');
            } catch (error) {
                console.error('导出失败:', error);
                showToast('导出失败: ' + error.message, 'error');
            }
        }

        // 导入TXT到本地存储
        async function importToLocal(file) {
            try {
                const text = await file.text();
                const lines = text.split('\\n').filter(line => line.trim());

                let imported = 0;
                let skipped = 0;

                for (const line of lines) {
                    const parts = line.split('----').map(p => p.trim());
                    if (parts.length >= 3) {
                        const account = {
                            email: parts[0],
                            password: parts[1],
                            token: parts[2],
                            apikey: parts[3] || null,
                            source: 'local',
                            createdAt: new Date().toISOString()
                        };

                        const success = await saveToLocal(account);
                        if (success) imported++;
                        else skipped++;
                    }
                }

                await loadLocalAccounts();
                showToast(\`导入完成！成功: \${imported}, 跳过: \${skipped}\`, 'success');
            } catch (error) {
                console.error('导入失败:', error);
                showToast('导入失败: ' + error.message, 'error');
            }
        }

        // 同步本地账号到服务器
        async function syncLocalToServer() {
            try {
                const localAccounts = await getAllLocalAccounts();
                const localOnly = localAccounts.filter(a => a.source === 'local');

                if (localOnly.length === 0) {
                    showToast('没有需要同步的本地账号', 'info');
                    return;
                }

                const response = await fetch('/api/sync-local', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ accounts: localOnly })
                });

                const result = await response.json();

                if (response.ok && result.success) {
                    // 同步成功后删除本地已同步的账号
                    const transaction = db.transaction([STORE_NAME], 'readwrite');
                    const store = transaction.objectStore(STORE_NAME);
                    const emailIndex = store.index('email');

                    let deleted = 0;
                    for (const acc of localOnly) {
                        const request = emailIndex.getKey(acc.email);
                        request.onsuccess = () => {
                            if (request.result) {
                                store.delete(request.result);
                                deleted++;
                            }
                        };
                    }

                    // 等待删除完成
                    transaction.oncomplete = async () => {
                        await loadLocalAccounts();
                        showToast(\`同步成功！已同步 \${result.synced} 个账号，已删除 \${deleted} 个本地记录\`, 'success');
                    };
                } else {
                    showToast(result.error || '同步失败', 'error');
                }
            } catch (error) {
                console.error('同步失败:', error);
                showToast('同步失败: ' + error.message, 'error');
            }
        }

        // 清空本地存储
        async function clearLocalAccounts() {
            if (!db) return;

            if (!confirm('确定要清空所有本地账号吗？此操作不可恢复！')) return;

            return new Promise((resolve, reject) => {
                const transaction = db.transaction([STORE_NAME], 'readwrite');
                const store = transaction.objectStore(STORE_NAME);
                const request = store.clear();

                request.onsuccess = () => {
                    loadLocalAccounts();
                    showToast('本地账号已清空', 'success');
                    resolve();
                };
                request.onerror = () => reject(request.error);
            });
        }

        // 重新获取单个账号的APIKEY
        async function refetchSingleApikey(email, token) {
            try {
                const response = await fetch('/api/refetch-apikey', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, token })
                });

                const result = await response.json();

                if (result.success) {
                    showToast('✓ ' + email + ' APIKEY获取成功', 'success');
                    // 更新本地账号数据
                    await loadAccounts();
                    return { success: true, apikey: result.apikey };
                } else {
                    showToast('✗ ' + email + ' ' + result.error, 'error');
                    return { success: false, error: result.error };
                }
            } catch (error) {
                console.error('获取APIKEY失败:', error);
                showToast('✗ ' + email + ' 获取失败: ' + error.message, 'error');
                return { success: false, error: error.message };
            }
        }

        // 批量获取APIKEY
        async function batchRefetchApikey() {
            // 找出所有没有APIKEY的账号
            const accountsWithoutKey = accounts.filter(acc => !acc.apikey);

            if (accountsWithoutKey.length === 0) {
                showToast('所有账号都已有APIKEY', 'info');
                return;
            }

            if (!confirm('发现 ' + accountsWithoutKey.length + ' 个账号缺少APIKEY，确定要批量获取吗？')) {
                return;
            }

            let successCount = 0;
            let failedCount = 0;
            const total = accountsWithoutKey.length;

            showToast('开始批量获取APIKEY，共 ' + total + ' 个账号...', 'info');

            for (let i = 0; i < accountsWithoutKey.length; i++) {
                const acc = accountsWithoutKey[i];
                addLog('[' + (i + 1) + '/' + total + '] 正在为 ' + acc.email + ' 获取APIKEY...', 'info');

                const result = await refetchSingleApikey(acc.email, acc.token);

                if (result.success) {
                    successCount++;
                    addLog('  ✓ ' + acc.email + ' 成功', 'success');
                } else {
                    failedCount++;
                    addLog('  ✗ ' + acc.email + ' 失败: ' + result.error, 'error');
                }

                // 每个请求之间延迟1秒，避免过快
                if (i < accountsWithoutKey.length - 1) {
                    await new Promise(resolve => setTimeout(resolve, 1000));
                }
            }

            showToast('批量获取完成！成功 ' + successCount + ' 个，失败 ' + failedCount + ' 个',
                      successCount > 0 ? 'success' : 'error');
        }

        // 批量检测账号存活性
        async function batchCheckAccounts() {
            if (accounts.length === 0) {
                showToast('暂无账号需要检测', 'info');
                return;
            }

            if (!confirm('确定要检测所有账号的存活性吗？这可能需要一些时间。')) {
                return;
            }

            const emails = accounts.map(acc => acc.email);
            showToast('开始批量检测，共 ' + emails.length + ' 个账号...', 'info');
            addLog('开始批量检测账号存活性...', 'info');

            try {
                const response = await fetch('/api/check-accounts', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ emails: emails })
                });

                const result = await response.json();

                if (result.success) {
                    const activeCount = result.results.filter(r => r.isActive).length;
                    const inactiveCount = result.results.filter(r => !r.isActive).length;

                    addLog('检测完成！正常: ' + activeCount + ' 个，失效: ' + inactiveCount + ' 个', 'success');
                    showToast('检测完成！正常: ' + activeCount + ' 个，失效: ' + inactiveCount + ' 个', 'success');

                    // 刷新账号列表
                    await loadAccounts();
                } else {
                    showToast('检测失败: ' + result.error, 'error');
                }
            } catch (error) {
                console.error('批量检测失败:', error);
                showToast('批量检测失败: ' + error.message, 'error');
            }
        }

        // 删除失效账号
        async function deleteInactiveAccounts() {
            const inactiveCount = accounts.filter(acc => acc.status === 'inactive').length;

            if (inactiveCount === 0) {
                showToast('没有失效账号需要删除', 'info');
                return;
            }

            if (!confirm('发现 ' + inactiveCount + ' 个失效账号，确定要删除吗？此操作不可恢复！')) {
                return;
            }

            try {
                const response = await fetch('/api/delete-inactive', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });

                const result = await response.json();

                if (result.success) {
                    showToast('成功删除 ' + result.deleted + ' 个失效账号', 'success');
                    addLog('成功删除 ' + result.deleted + ' 个失效账号', 'success');
                    await loadAccounts();
                } else {
                    showToast('删除失败: ' + result.error, 'error');
                }
            } catch (error) {
                console.error('删除失效账号失败:', error);
                showToast('删除失败: ' + error.message, 'error');
            }
        }

        function connectSSE() {
            const eventSource = new EventSource('/events');
            eventSource.onmessage = (event) => {
                const data = JSON.parse(event.data);
                switch (data.type) {
                    case 'connected':
                        addLog('✓ 已连接到服务器', 'success');
                        updateStatus(data.isRunning);
                        break;
                    case 'start':
                        updateStatus(true);
                        taskStartTime = Date.now();
                        totalTaskCount = data.config.count;
                        $progressContainer.show();
                        updateProgress(0, totalTaskCount, 0, 0);
                        addLog('🚀 开始注册 ' + data.config.count + ' 个账号', 'info');
                        $successCount.text(0);
                        $failedCount.text(0);
                        break;
                    case 'log':
                        addLog(data.message, data.level, data.link);
                        if (data.stats) {
                            $successCount.text(data.stats.success);
                            $failedCount.text(data.stats.failed);
                            updateProgress(data.stats.total, totalTaskCount, data.stats.success, data.stats.failed);
                        }
                        break;
                    case 'account_added':
                        accounts.unshift(data.account);
                        filteredAccounts = accounts;
                        $totalAccounts.text(accounts.length);
                        renderTable();
                        // KV账号不需要保存到IndexedDB（已在服务器，无需本地备份）
                        break;
                    case 'local_account_added':
                        // KV保存失败，仅保存到IndexedDB
                        data.account.source = 'local'; // 标记为仅本地账号
                        saveToLocal(data.account).then(() => {
                            addLog(\`💾 账号已保存到本地存储: \${data.account.email}\`, 'warning');
                            loadLocalAccounts(); // 更新本地账号统计
                        }).catch(err => {
                            console.error('保存到本地失败:', err);
                            addLog(\`❌ 本地保存失败: \${data.account.email}\`, 'error');
                        });
                        break;
                    case 'complete':
                        updateStatus(false);
                        $successCount.text(data.stats.success);
                        $failedCount.text(data.stats.failed);
                        $timeValue.text(data.stats.elapsedTime + 's');
                        updateProgress(data.stats.total, totalTaskCount, data.stats.success, data.stats.failed);
                        addLog('✓ 注册完成！成功: ' + data.stats.success + ', 失败: ' + data.stats.failed, 'success');
                        setTimeout(() => $progressContainer.fadeOut(), 3000);
                        break;
                }
            };
            eventSource.onerror = () => {
                addLog('✗ 连接断开，5秒后重连...', 'error');
                eventSource.close();
                setTimeout(connectSSE, 5000);
            };
        }

        $(document).ready(async function() {
            await initIndexedDB(); // 初始化IndexedDB
            loadAccounts();
            loadSettings();
            connectSSE();
        });
    </script>
</body>
</html>`;

// HTTP 处理器
async function handler(req: Request): Promise<Response> {
  const url = new URL(req.url);

  // 登录页面（无需鉴权）
  if (url.pathname === "/login") {
    return new Response(LOGIN_PAGE, { headers: { "Content-Type": "text/html; charset=utf-8" } });
  }

  // 登录 API（无需鉴权）
  if (url.pathname === "/api/login" && req.method === "POST") {
    const clientIP = getClientIP(req);

    // 检查 IP 是否被锁定
    const lockCheck = checkIPLocked(clientIP);
    if (lockCheck.locked) {
      console.log(`🚫 IP ${clientIP} 尝试登录但已被锁定，剩余 ${lockCheck.remainingTime} 秒`);
      return new Response(JSON.stringify({
        success: false,
        error: `登录失败次数过多，账号已被锁定`,
        remainingTime: lockCheck.remainingTime,
        code: "ACCOUNT_LOCKED"
      }), {
        status: 429,  // Too Many Requests
        headers: { "Content-Type": "application/json" }
      });
    }

    const body = await req.json();
    if (body.username === AUTH_USERNAME && body.password === AUTH_PASSWORD) {
      // 登录成功，清除失败记录
      clearLoginFailure(clientIP);
      const sessionId = generateSessionId();

      // 保存 session 到 KV，设置 24 小时过期
      const sessionKey = ["sessions", sessionId];
      try {
        await kv.set(sessionKey, { createdAt: Date.now() }, { expireIn: 86400000 }); // 24小时过期
      } catch (error) {
        console.error("❌ Failed to save session to KV:", error);

        // Check if it's a quota exhausted error
        const errorMessage = error instanceof Error ? error.message : String(error);
        if (errorMessage.includes("quota is exhausted")) {
          return new Response(JSON.stringify({
            success: false,
            error: "KV 存储配额已耗尽，请清理数据或升级配额"
          }), {
            status: 507, // Insufficient Storage
            headers: { "Content-Type": "application/json" }
          });
        }

        return new Response(JSON.stringify({
          success: false,
          error: "登录失败: 无法保存会话"
        }), {
          status: 500,
          headers: { "Content-Type": "application/json" }
        });
      }

      console.log(`✅ IP ${clientIP} 登录成功`);
      return new Response(JSON.stringify({ success: true, sessionId }), {
        headers: { "Content-Type": "application/json" }
      });
    }

    // 登录失败，记录失败次数
    recordLoginFailure(clientIP);
    const attempts = loginAttempts.get(clientIP)?.attempts || 0;
    console.log(`❌ IP ${clientIP} 登录失败（第 ${attempts} 次）`);

    return new Response(JSON.stringify({
      success: false,
      error: "用户名或密码错误",
      attemptsRemaining: Math.max(0, MAX_LOGIN_ATTEMPTS - attempts)
    }), {
      status: 401,
      headers: { "Content-Type": "application/json" }
    });
  }

  // 鉴权检查（其他所有路径都需要验证）
  const auth = await checkAuth(req);
  if (!auth.authenticated) {
    // 判断是 API 请求还是页面请求
    const isApiRequest = url.pathname.startsWith('/api/');

    if (isApiRequest) {
      // API 请求返回 401 JSON 响应
      return new Response(JSON.stringify({
        success: false,
        error: "未授权访问，请先登录",
        code: "UNAUTHORIZED"
      }), {
        status: 401,
        headers: { "Content-Type": "application/json" }
      });
    } else {
      // 页面请求返回 302 重定向
      return new Response(null, {
        status: 302,
        headers: { "Location": "/login" }
      });
    }
  }

  // 登出 API
  if (url.pathname === "/api/logout" && req.method === "POST") {
    if (auth.sessionId) {
      // 从 KV 删除 session
      const sessionKey = ["sessions", auth.sessionId];
      await kv.delete(sessionKey);
    }
    return new Response(JSON.stringify({ success: true }), {
      headers: { "Content-Type": "application/json" }
    });
  }

  // 主页
  if (url.pathname === "/" || url.pathname === "/index.html") {
    return new Response(HTML_PAGE, { headers: { "Content-Type": "text/html; charset=utf-8" } });
  }

  // 获取配置
  if (url.pathname === "/api/config" && req.method === "GET") {
    // 从 KV 读取配置，如果不存在则返回默认值
    const configKey = ["config", "register"];
    const savedConfig = await kv.get(configKey);

    const config = savedConfig.value || registerConfig;
    return new Response(JSON.stringify(config), {
      headers: { "Content-Type": "application/json" }
    });
  }

  // 保存配置
  if (url.pathname === "/api/config" && req.method === "POST") {
    const body = await req.json();
    registerConfig = { ...registerConfig, ...body };

    // 保存到 KV 持久化
    const configKey = ["config", "register"];
    await kv.set(configKey, registerConfig);

    return new Response(JSON.stringify({ success: true }), {
      headers: { "Content-Type": "application/json" }
    });
  }

  // SSE
  if (url.pathname === "/events") {
    console.log(`🔌 新的 SSE 连接建立，当前客户端数: ${sseClients.size + 1}`);
    const stream = new ReadableStream({
      start(controller) {
        sseClients.add(controller);
        console.log(`✓ SSE 客户端已添加到连接池，isRunning=${isRunning}`);
        // 发送当前状态
        const message = `data: ${JSON.stringify({ type: 'connected', isRunning })}\n\n`;
        controller.enqueue(new TextEncoder().encode(message));

        // 发送历史日志（最近50条）
        const recentLogs = logHistory.slice(-50);
        for (const log of recentLogs) {
          const logMessage = `data: ${JSON.stringify(log)}\n\n`;
          controller.enqueue(new TextEncoder().encode(logMessage));
        }

        const keepAlive = setInterval(() => {
          try {
            controller.enqueue(new TextEncoder().encode(": keepalive\n\n"));
          } catch {
            clearInterval(keepAlive);
            sseClients.delete(controller);
          }
        }, 30000);
      }
    });

    return new Response(stream, {
      headers: { "Content-Type": "text/event-stream", "Cache-Control": "no-cache", "Connection": "keep-alive" }
    });
  }

  // 获取运行状态（新增 API）
  if (url.pathname === "/api/status") {
    return new Response(JSON.stringify({
      isRunning,
      stats,
      logCount: logHistory.length
    }), {
      headers: { "Content-Type": "application/json" }
    });
  }

  // 账号列表
  if (url.pathname === "/api/accounts") {
    const accounts = [];
    const entries = kv.list({ prefix: ["zai_accounts"] }, { reverse: true });
    for await (const entry of entries) {
      accounts.push(entry.value);
    }
    return new Response(JSON.stringify(accounts), { headers: { "Content-Type": "application/json" } });
  }

  // 导出
  if (url.pathname === "/api/export") {
    const lines: string[] = [];
    const entries = kv.list({ prefix: ["zai_accounts"] });
    for await (const entry of entries) {
      const data = entry.value as any;
      // 支持四字段格式：账号----密码----Token----APIKEY
      if (data.apikey) {
        lines.push(`${data.email}----${data.password}----${data.token}----${data.apikey}`);
      } else {
        // 兼容旧格式，APIKEY为空
        lines.push(`${data.email}----${data.password}----${data.token}----`);
      }
    }
    return new Response(lines.join('\n'), {
      headers: {
        "Content-Type": "text/plain",
        "Content-Disposition": `attachment; filename="zai_accounts_${Date.now()}.txt"`
      }
    });
  }

  // 导入
  if (url.pathname === "/api/import" && req.method === "POST") {
    try {
      const body = await req.json();
      const { email, password, token, apikey } = body;

      if (!email || !password || !token) {
        return new Response(JSON.stringify({ success: false, error: "缺少必要字段" }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }

      // 保存到 KV
      const timestamp = Date.now();
      const key = ["zai_accounts", timestamp, email];
      try {
        await kv.set(key, {
          email,
          password,
          token,
          apikey: apikey || null,  // 支持APIKEY字段
          createdAt: new Date().toISOString()
        });
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : String(error);
        if (errorMessage.includes("quota is exhausted")) {
          return new Response(JSON.stringify({
            success: false,
            error: "KV 存储配额已耗尽，无法导入账号"
          }), {
            status: 507,
            headers: { "Content-Type": "application/json" }
          });
        }
        throw error;
      }

      return new Response(JSON.stringify({ success: true }), {
        headers: { "Content-Type": "application/json" }
      });
    } catch (error: any) {
      const msg = error instanceof Error ? error.message : String(error);
      return new Response(JSON.stringify({ success: false, error: msg }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  // 批量导入（优化性能，支持去重）
  if (url.pathname === "/api/import-batch" && req.method === "POST") {
    try {
      const body = await req.json();
      const { accounts: importAccounts } = body;

      if (!Array.isArray(importAccounts)) {
        return new Response(JSON.stringify({ success: false, error: "数据格式错误" }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }

      // 获取已存在的邮箱
      const existingEmails = new Set();
      const entries = kv.list({ prefix: ["zai_accounts"] });
      for await (const entry of entries) {
        const data = entry.value as any;
        existingEmails.add(data.email);
      }

      // 批量写入（去重）
      let imported = 0;
      let skipped = 0;
      let quotaExhausted = false;
      const timestamp = Date.now();

      for (const [index, acc] of importAccounts.entries()) {
        const { email, password, token, apikey } = acc;

        if (!email || !password || !token) {
          skipped++;
          continue;
        }

        // 检查是否已存在
        if (existingEmails.has(email)) {
          skipped++;
          continue;
        }

        // 使用不同的时间戳避免键冲突
        const key = ["zai_accounts", timestamp + index, email];
        try {
          await kv.set(key, {
            email,
            password,
            token,
            apikey: apikey || null,  // 支持APIKEY字段
            createdAt: new Date().toISOString()
          });

          existingEmails.add(email);
          imported++;
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : String(error);
          if (errorMessage.includes("quota is exhausted")) {
            console.error("❌ KV quota exhausted during batch import");
            quotaExhausted = true;
            break; // Stop importing if quota is exhausted
          }
          // Log other errors but continue
          console.error(`Failed to import account ${email}:`, error);
          skipped++;
        }
      }

      if (quotaExhausted) {
        return new Response(JSON.stringify({
          success: false,
          imported,
          skipped: skipped + (importAccounts.length - imported - skipped),
          error: "KV 存储配额已耗尽，已导入 " + imported + " 个账号"
        }), {
          status: 507,
          headers: { "Content-Type": "application/json" }
        });
      }

      return new Response(JSON.stringify({ success: true, imported, skipped }), {
        headers: { "Content-Type": "application/json" }
      });
    } catch (error: any) {
      const msg = error instanceof Error ? error.message : String(error);
      return new Response(JSON.stringify({ success: false, error: msg }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  // 同步本地账号到服务器
  if (url.pathname === "/api/sync-local" && req.method === "POST") {
    try {
      const body = await req.json();
      const { accounts: localAccounts } = body;

      if (!Array.isArray(localAccounts)) {
        return new Response(JSON.stringify({ success: false, error: "数据格式错误" }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }

      // 获取已存在的邮箱
      const existingEmails = new Set();
      const entries = kv.list({ prefix: ["zai_accounts"] });
      for await (const entry of entries) {
        const data = entry.value as any;
        existingEmails.add(data.email);
      }

      // 批量同步（去重）
      let synced = 0;
      let skipped = 0;
      let quotaExhausted = false;
      const timestamp = Date.now();

      for (const [index, acc] of localAccounts.entries()) {
        const { email, password, token, apikey } = acc;

        if (!email || !password || !token) {
          skipped++;
          continue;
        }

        // 检查是否已存在
        if (existingEmails.has(email)) {
          skipped++;
          continue;
        }

        // 使用不同的时间戳避免键冲突
        const key = ["zai_accounts", timestamp + index, email];
        try {
          await kv.set(key, {
            email,
            password,
            token,
            apikey: apikey || null,
            createdAt: acc.createdAt || new Date().toISOString()
          });

          existingEmails.add(email);
          synced++;
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : String(error);
          if (errorMessage.includes("quota is exhausted")) {
            console.error("❌ KV quota exhausted during sync");
            quotaExhausted = true;
            break;
          }
          console.error(`Failed to sync account ${email}:`, error);
          skipped++;
        }
      }

      if (quotaExhausted) {
        return new Response(JSON.stringify({
          success: false,
          synced,
          skipped: skipped + (localAccounts.length - synced - skipped),
          error: "KV 存储配额已耗尽，已同步 " + synced + " 个账号"
        }), {
          status: 507,
          headers: { "Content-Type": "application/json" }
        });
      }

      return new Response(JSON.stringify({ success: true, synced, skipped }), {
        headers: { "Content-Type": "application/json" }
      });
    } catch (error: any) {
      const msg = error instanceof Error ? error.message : String(error);
      return new Response(JSON.stringify({ success: false, error: msg }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  // 开始注册
  if (url.pathname === "/api/register" && req.method === "POST") {
    if (isRunning) {
      return new Response(JSON.stringify({
        success: false,
        error: "任务正在运行中，请等待当前任务完成或手动停止后再试",
        isRunning: true,
        stats: {
          success: stats.success,
          failed: stats.failed,
          total: stats.success + stats.failed
        }
      }), {
        status: 400,
        headers: { "Content-Type": "application/json" }
      });
    }

    const body = await req.json();
    const count = body.count || 5;

    // 立即启动任务（不等待完成）
    batchRegister(count).catch(err => {
      console.error("注册任务异常:", err);
      broadcast({ type: 'log', level: 'error', message: `✗ 任务异常: ${err.message}` });
    });

    return new Response(JSON.stringify({ success: true }), { headers: { "Content-Type": "application/json" } });
  }

  // 停止注册
  if (url.pathname === "/api/stop" && req.method === "POST") {
    if (!isRunning) {
      return new Response(JSON.stringify({ error: "没有运行中的任务" }), {
        status: 400,
        headers: { "Content-Type": "application/json" }
      });
    }

    shouldStop = true;
    return new Response(JSON.stringify({ success: true }), { headers: { "Content-Type": "application/json" } });
  }

  // 重新获取APIKEY（单个账号）
  if (url.pathname === "/api/refetch-apikey" && req.method === "POST") {
    try {
      const body = await req.json();
      const { email, token } = body;

      if (!email || !token) {
        return new Response(JSON.stringify({
          success: false,
          error: "缺少必需参数: email 或 token"
        }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }

      // 尝试使用Token快速获取APIKEY
      const accessToken = await loginToApi(cleanToken(token));
      if (!accessToken) {
        return new Response(JSON.stringify({
          success: false,
          error: "Token已失效，请使用账号密码重新注册"
        }), {
          status: 401,
          headers: { "Content-Type": "application/json" }
        });
      }

      const { orgId, projectId } = await getCustomerInfo(accessToken);
      if (!orgId || !projectId) {
        return new Response(JSON.stringify({
          success: false,
          error: "获取客户信息失败"
        }), {
          status: 500,
          headers: { "Content-Type": "application/json" }
        });
      }

      const apikey = await createApiKey(accessToken, orgId, projectId);
      if (!apikey) {
        return new Response(JSON.stringify({
          success: false,
          error: "创建APIKEY失败"
        }), {
          status: 500,
          headers: { "Content-Type": "application/json" }
        });
      }

      // 更新KV中的账号APIKEY
      const entries = kv.list({ prefix: ["zai_accounts"] });
      for await (const entry of entries) {
        const account = entry.value as any;
        if (account.email === email) {
          await kv.set(entry.key, {
            ...account,
            apikey: apikey
          });
          break;
        }
      }

      return new Response(JSON.stringify({
        success: true,
        apikey: apikey
      }), {
        headers: { "Content-Type": "application/json" }
      });

    } catch (error: any) {
      return new Response(JSON.stringify({
        success: false,
        error: "请求错误: " + error?.message
      }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  // 批量检测账号存活性
  if (url.pathname === "/api/check-accounts" && req.method === "POST") {
    try {
      const body = await req.json();
      const { emails } = body;

      if (!emails || !Array.isArray(emails)) {
        return new Response(JSON.stringify({
          success: false,
          error: "缺少必需参数: emails"
        }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }

      const results: any[] = [];
      const entries = kv.list({ prefix: ["zai_accounts"] });

      for await (const entry of entries) {
        const account = entry.value as any;
        if (emails.includes(account.email)) {
          const isActive = await checkAccountStatus(account.token);
          const newStatus = isActive ? 'active' : 'inactive';

          // 更新账号状态
          await kv.set(entry.key, {
            ...account,
            status: newStatus
          });

          results.push({
            email: account.email,
            status: newStatus,
            isActive: isActive
          });
        }
      }

      return new Response(JSON.stringify({
        success: true,
        results: results
      }), {
        headers: { "Content-Type": "application/json" }
      });

    } catch (error: any) {
      return new Response(JSON.stringify({
        success: false,
        error: "请求错误: " + error?.message
      }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  // 删除失效账号
  if (url.pathname === "/api/delete-inactive" && req.method === "POST") {
    try {
      let deletedCount = 0;
      const entries = kv.list({ prefix: ["zai_accounts"] });

      for await (const entry of entries) {
        const account = entry.value as any;
        if (account.status === 'inactive') {
          await kv.delete(entry.key);
          deletedCount++;
        }
      }

      return new Response(JSON.stringify({
        success: true,
        deleted: deletedCount
      }), {
        headers: { "Content-Type": "application/json" }
      });

    } catch (error: any) {
      return new Response(JSON.stringify({
        success: false,
        error: "请求错误: " + error?.message
      }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  return new Response("Not Found", { status: 404 });
}

// Initialize KV database before loading config
await initKV();

// 启动时从 KV 加载配置和日志
(async () => {
  // 加载配置
  const configKey = ["config", "register"];
  const savedConfig = await kv.get(configKey);
  if (savedConfig.value) {
    registerConfig = { ...registerConfig, ...savedConfig.value };
    console.log("✓ 已加载保存的配置");
  }

  // 清理历史日志（重启时清空）
  const logKey = ["logs", "recent"];
  try {
    await kv.delete(logKey);
    console.log("✓ 已清理历史日志数据");
  } catch (error) {
    console.log("⚠️ 清理日志失败:", error);
  }
})();

console.log(`🚀 Z.AI 管理系统 V2 启动: http://localhost:${PORT}`);
console.log(`🔐 登录账号: ${AUTH_USERNAME}`);
console.log(`🔑 登录密码: ${AUTH_PASSWORD}`);
console.log(`💡 访问 http://localhost:${PORT}/login 登录`);
await serve(handler, { port: PORT });

/*
  📦 源码地址:
  https://github.com/dext7r/ZtoApi/tree/main/deno/zai/zai_register.ts
  |
  💬 交流讨论: https://linux.do/t/topic/1009939
──────────────────────────────────────────────────
*/
