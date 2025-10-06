# Two2Api - Two Chat Sutra OpenAI 兼容代理

> 更新时间：2025-10-02 18:45 (UTC+8) · 执行者 Codex

Two2Api 依托仓库提供的 Deno 模板，将 Two Chat 前端所用的 `Sutra` 接口封装为标准的 `/v1/chat/completions` 与 `/v1/models`。配置好会话令牌后，可直接以 OpenAI SDK 方式访问，也保留了首页、文档、部署、监控与 Playground 页面。

## 功能亮点

- ✅ 上游指向 `https://chatsutra-server.account-2b0.workers.dev/v2/chat/completions`
- ✅ 透传 `extra_body.online_search` 等参数，保持 Two Chat 原生能力
- ✅ 支持流式/非流式模式（如上游开启 SSE，会自动转换为 OpenAI chunk）
- ✅ 内建监控面板与多语言页面，便于观测请求统计
- ✅ `.env` 即可配置 `x-session-token` 与默认模型

## 目录结构

```
two2api/
├── main.ts            # Two Chat 定制主程序
├── .env.example       # 环境变量示例（含 X_SESSION_TOKEN）
├── deno.json          # Deno 任务与编译配置
├── lib/               # 类型、工具、SEO、i18n 逻辑
├── pages/             # 文档、部署、Playground 页面
├── README.md          # 当前说明文档
└── start.sh           # 本地启动脚本（支持 .env 自动加载）
```

## 快速开始

1. 复制示例配置

   ```bash
   cp .env.example .env
   ```

2. 编辑 `.env`

   - `X_SESSION_TOKEN`：在 Two Chat 浏览器中获取的 `x-session-token`
   - `MODEL_NAME`：默认 `sutra-v2`，可根据前端返回的可用模型调整
   - `DEFAULT_KEY`：本地代理鉴权用密钥，调用时需携带

3. 启动服务

   ```bash
   deno task dev   # 开发模式
   # 或
   deno task start # 生产模式
   ```

   浏览器访问 `http://localhost:9090` 查看首页，Playground 页面可直接测试请求。

## 环境变量说明

| 变量名 | 用途 | 默认值 |
| --- | --- | --- |
| `PORT` | 监听端口 | `9090` |
| `DEBUG_MODE` | 是否打印调试日志 | `false` |
| `DEFAULT_STREAM` | 未指定时使用流式 | `true` |
| `DASHBOARD_ENABLED` | 是否开启监控页 | `true` |
| `UPSTREAM_URL` | Two Chat Sutra 接口地址 | `https://chatsutra-server.account-2b0.workers.dev/v2/chat/completions` |
| `X_SESSION_TOKEN` | Two Chat 上游身份令牌 | 空（必填） |
| `X_SESSION_COOKIE` | 会话 Cookie（如 `authjs.session-token=...`） | 空（建议同步） |
| `DEFAULT_KEY` | 本地代理访问密钥 | `sk-two-demo` |
| `MODEL_NAME` | 默认模型 | `sutra-v2` |
| `DEFAULT_TEMPERATURE` | 默认温度 | `0.6` |
| `DEFAULT_MAX_TOKENS` | 默认最大 tokens | `2048` |
| `DEFAULT_EXTRA_BODY` | 默认 `extra_body` JSON（自动填充） | `{"online_search":false,...}` |
| `SERVICE_NAME` | 页面展示名称 | `Two2Api` |
| `SERVICE_EMOJI` | 页面展示表情 | `🌀` |
| `FOOTER_TEXT` | 页脚文案 | `Two 体验转接 OpenAI 接口` |
| `DISCUSSION_URL` | 讨论入口 | `https://github.com/dext7r/ZtoApi` |
| `GITHUB_REPO` | 仓库链接 | `https://github.com/dext7r/ZtoApi` |

> 建议从浏览器开发者工具复制最新的 `x-session-token`，一旦过期需重新更新。

## 接口调用示例

```bash
# 获取模型列表
curl http://localhost:9090/v1/models \
  -H "Authorization: Bearer sk-two-demo"

# 非流式对话
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-two-demo" \
  -d '{
    "model": "sutra-v2",
    "messages": [
      {"role": "user", "content": "介绍下 Two2Api"}
    ],
    "extra_body": {"online_search": true},
    "stream": false
  }'

# 流式对话
curl -N -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-two-demo" \
  -d '{
    "model": "sutra-v2",
    "messages": [
      {"role": "user", "content": "来点实时搜索内容"}
    ],
    "extra_body": {"online_search": true},
    "stream": true
  }'
```

在 SDK 中只需将 `base_url` 指向 `http://localhost:9090/v1`，并把密钥替换为 `DEFAULT_KEY` 即可。

## Two Chat 定制要点

- `transformToUpstream` 透传 `extra_body`、`temperature` 等参数，确保可用前端附加能力
- `transformFromUpstream` 针对 Two 返回的 `text` / `output` 字段做兼容补全
- `getUpstreamHeaders` 自动注入 `x-session-token`、Origin、Referer，模拟浏览器调用

## 部署建议

- **Deno Deploy**：`deployctl deploy --project=<name> main.ts`，在环境变量中写入 `X_SESSION_TOKEN`、`DEFAULT_KEY`
- **容器/自托管**：使用 `start.sh` 或自定义脚本加载 `.env` 后运行 `deno run --allow-net --allow-env main.ts`
- **监控**：访问 `/dashboard` 观察实时请求与耗时

## 许可证

MIT License
