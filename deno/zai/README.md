# OpenAI兼容API代理 for Z.ai GLM-4.5 (Deno版本)

这是一个为Z.ai GLM-4.5模型提供OpenAI兼容API接口的代理服务器（Deno版本）。它允许你使用标准的OpenAI API格式与Z.ai的GLM-4.5模型进行交互，支持流式和非流式响应。

> **注意**: 本项目基于 [OpenAI-Compatible-API-Proxy-for-Z](https://github.com/kbykb/OpenAI-Compatible-API-Proxy-for-Z) 二次开发，并改造为 Deno 版本

## ✨ 主要功能

- 🔄 **OpenAI API兼容**: 完全兼容OpenAI的API格式，无需修改客户端代码
- 🌊 **流式响应支持**: 支持实时流式输出,提供更好的用户体验
- 🔐 **身份验证**: 支持API密钥验证，确保服务安全
- 🛠️ **灵活配置**: 通过环境变量进行灵活配置
- 🐳 **Docker支持**: 提供Docker镜像，便于部署
- 🌍 **CORS支持**: 支持跨域请求，便于前端集成
- 📝 **思考过程展示**: 智能处理并展示模型的思考过程
- 📊 **实时监控仪表板**: 提供Web仪表板，实时显示API转发情况和统计信息
- 🔑 **账号管理系统**: 支持批量导入导出Z.ai账号（含APIKEY），基于Session认证
- 🎮 **Playground测试**: 在线测试API请求和响应（需登录）
- 🤖 **批量注册工具**: 自动批量注册Z.ai账号，支持并发和实时监控
- 🦕 **Deno运行时**: 使用现代化的Deno运行时，安全且高效

## 🚀 快速开始

### 环境要求

- Deno 1.40 或更高版本
- Z.ai 的访问令牌（可选，不提供将自动获取匿名token）

### 安装 Deno

**macOS/Linux:**
```bash
curl -fsSL https://deno.land/install.sh | sh
```

**Windows:**
```powershell
irm https://deno.land/install.ps1 | iex
```

### 本地部署

1. **克隆仓库**
   ```bash
   git clone https://github.com/your-username/ZtoApi.git
   cd ZtoApi
   ```

2. **配置环境变量**
   ```bash
   cp config.env .env.local
   # 编辑 .env.local 文件，设置你的 ZAI_TOKEN（可选）
   ```

3. **启动服务**
   ```bash
   # 使用启动脚本（推荐）
   ./start.sh          # macOS/Linux
   start.bat           # Windows

   # 或直接运行
   deno task start

   # 开发模式（自动重启）
   deno task dev
   ```

4. **测试服务**
    ```bash
    curl http://localhost:9090/v1/models
    ```

5. **访问API文档**

   启动服务后，可以通过浏览器访问以下地址查看API文档：
    ```
    http://localhost:9090/docs
    ```

6. **访问Dashboard**

   启动服务后，可以通过浏览器访问以下地址查看实时监控仪表板：
   ```
   http://localhost:9090/dashboard
   ```

### Docker部署

1. **构建镜像**
   ```bash
   docker build -f Dockerfile.deno -t zto-api-deno .
   ```

2. **运行容器**
   ```bash
   docker run -p 9090:9090 \
     -e ZAI_TOKEN=your_z_ai_token \
     -e DEFAULT_KEY=your_api_key \
     zto-api-deno
   ```

## ⚙️ 环境变量配置

### 🚀 快速开始

#### 1. 使用启动脚本（推荐）

**macOS/Linux:**
```bash
./start.sh
```

**Windows:**
```cmd
start.bat
```

#### 2. 手动设置环境变量

**macOS/Linux:**
```bash
export ZAI_TOKEN="your_z_ai_token_here"
export DEFAULT_KEY="sk-your-custom-key"
export PORT="9090"
deno task start
```

**Windows:**
```cmd
set ZAI_TOKEN=your_z_ai_token_here
set DEFAULT_KEY=sk-your-custom-key
set PORT=9090
deno task start
```

#### 3. Docker运行

```bash
docker run -p 9090:9090 \
  -e ZAI_TOKEN=your_z_ai_token_here \
  -e DEFAULT_KEY=sk-your-custom-key \
  -e PORT=9090 \
  zto-api-deno
```

### 📋 环境变量列表

#### 🔑 必需配置

无必需配置。所有配置都有合理的默认值。

#### ⚙️ 可选配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `ZAI_TOKEN` | Z.ai 访问令牌 | 空（自动获取随机匿名token） | `eyJhbGciOiJFUzI1NiIs...` |
| `DEFAULT_KEY` | 客户端API密钥 | `sk-your-key` | `sk-my-api-key` |
| `MODEL_NAME` | 显示模型名称 | `GLM-4.5` | `GLM-4.5-Pro` |
| `PORT` | 服务监听端口 | `9090` | `9000` |
| `DEBUG_MODE` | 调试模式开关 | `true` | `false` |
| `DEFAULT_STREAM` | 默认流式响应 | `true` | `false` |
| `DASHBOARD_ENABLED` | Dashboard功能开关 | `true` | `false` |
| `ENABLE_THINKING` | 思考功能开关 | `false` | `true` |

#### 🔧 高级配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `UPSTREAM_URL` | 上游API地址 | `https://chat.z.ai/api/chat/completions` | 自定义URL |
| `ADMIN_USERNAME` | 管理员用户名 | `admin` | `myadmin` |
| `ADMIN_PASSWORD` | 管理员密码 | `123456` | `strong_password_123` |
| `ADMIN_ENABLED` | 管理面板开关 | `true` | `false` |

⚠️ **安全提示**：生产环境务必修改 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD`，避免使用默认值。

### 📁 配置文件

#### 支持的配置文件（按优先级排序）

1. `.env.local` - 本地环境配置（推荐）
2. `.env` - 环境配置
3. `config.env` - 配置模板

#### 配置文件示例

```bash
# 复制配置文件
cp config.env .env.local

# 编辑配置文件
nano .env.local
```

### 🔐 获取 Z.ai Token

#### 方法1：浏览器开发者工具

1. 登录 [Z.ai](https://chat.z.ai)
2. 打开浏览器开发者工具（F12）
3. 切换到 Network 标签页
4. 发送一条消息
5. 在请求中找到 `Authorization` 头部的 Bearer token

#### 方法2：匿名Token

本项目支持自动获取匿名token，无需手动配置。当 `ZAI_TOKEN` 环境变量未设置时，系统会自动为每次对话获取不同的随机匿名token，避免共享记忆。

## 📖 API使用示例

### Python示例

```python
import openai

# 配置客户端
client = openai.OpenAI(
    api_key="your-api-key",  # 对应 DEFAULT_KEY
    base_url="http://localhost:9090/v1"
)

# 非流式请求
response = client.chat.completions.create(
    model="GLM-4.5",
    messages=[{"role": "user", "content": "你好，请介绍一下自己"}]
)

print(response.choices[0].message.content)

# 流式请求
response = client.chat.completions.create(
    model="GLM-4.5",
    messages=[{"role": "user", "content": "请写一首关于春天的诗"}],
    stream=True
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### curl示例

```bash
# 非流式请求
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "model": "GLM-4.5",
    "messages": [{"role": "user", "content": "你好"}],
    "stream": false
  }'

# 流式请求
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "model": "GLM-4.5",
    "messages": [{"role": "user", "content": "你好"}],
    "stream": true
  }'
```

## 🦕 Deno 特性

### 优势

- ✅ **安全第一**: 默认安全，显式权限控制
- ✅ **现代标准**: 原生支持 TypeScript，无需配置
- ✅ **简单部署**: 单一可执行文件，无需 node_modules
- ✅ **内置工具**: 格式化、测试、打包等工具开箱即用
- ✅ **Web标准**: 使用现代 Web API（fetch、streams等）

### 安全增强

- 🔒 **管理面板认证**: 基于 Session 的登录系统，24小时自动过期
- 🔐 **Playground 访问控制**: 需登录后才能使用 Playground 测试页面
- 🛡️ **敏感信息保护**: 登录页面和前端不暴露默认密码、API Key 等敏感配置
- 📋 **账号数据隔离**: 使用 Deno KV 本地存储，支持导入导出和批量管理（含 APIKEY）

### 任务命令

```bash
# 启动服务
deno task start

# 开发模式（自动重启）
deno task dev

# 缓存依赖
deno task cache
```

### 权限说明

本项目需要以下权限：
- `--allow-net`: 网络访问（必需，用于HTTP服务和上游API调用）
- `--allow-env`: 环境变量读取（必需，用于配置）
- `--allow-read`: 文件读取（可选，用于读取配置文件）

## 🔧 故障排除

### 常见问题

1. **Deno未安装**
   - 参考上面的安装说明安装 Deno

2. **连接失败**
   - 检查服务是否正常运行：`curl http://localhost:9090/v1/models`
   - 访问API文档：`http://localhost:9090/docs`

3. **认证失败**
   - 检查 `DEFAULT_KEY` 环境变量设置
   - 确认请求头中的 `Authorization` 格式正确

4. **端口被占用**
   - 修改 `PORT` 环境变量
   - 或停止占用端口的服务

5. **权限错误**
   - 确保使用正确的权限标志运行 Deno
   - 参考上面的"权限说明"部分

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

## 📚 相关文档

- **[ADMIN_GUIDE.md](./ADMIN_GUIDE.md)** - 账号管理系统使用指南
- **[zai_register.md](./zai_register.md)** - 批量注册工具使用文档
- **[CLAUDE.md](../../CLAUDE.md)** - 项目架构和开发指南

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## ⚠️ 免责声明

本项目与 Z.ai 官方无关，使用前请确保遵守 Z.ai 的服务条款。开发者不对因使用本项目而产生的任何问题负责。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 [Issue](https://github.com/hulisang/ZtoApi/issues)
