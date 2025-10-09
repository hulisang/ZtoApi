#!/bin/bash
# ========================================
# ZtoApi Docker 启动脚本
# 同时运行 API 服务和注册工具
# ========================================

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}    ZtoApi Docker Container${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""

# 显示配置信息
echo -e "${BLUE}📋 配置信息:${NC}"
echo "  API 服务端口: ${PORT:-9090}"
echo "  注册工具端口: 8001"
echo "  模型: ${MODEL_NAME:-GLM-4.5}"
echo "  KV 数据库: /app/.deno-kv/ (共享)"
echo "  Admin 账号: ${ADMIN_USERNAME:-admin}"
echo "  注册工具账号: ${ZAI_USERNAME:-admin}"
echo ""

# 创建 KV 数据库目录（确保权限）
mkdir -p /app/.deno-kv
echo -e "${GREEN}✓${NC} KV 数据库目录已创建"

# 定义日志文件
API_LOG="/app/logs/api.log"
REGISTER_LOG="/app/logs/register.log"
mkdir -p /app/logs

# 信号处理函数
cleanup() {
    echo ""
    echo -e "${YELLOW}收到停止信号，正在关闭服务...${NC}"
    kill $API_PID 2>/dev/null || true
    kill $REGISTER_PID 2>/dev/null || true
    wait $API_PID 2>/dev/null || true
    wait $REGISTER_PID 2>/dev/null || true
    echo -e "${GREEN}✓${NC} 服务已停止"
    exit 0
}

trap cleanup SIGTERM SIGINT

# 启动 API 服务（后台运行）
echo -e "${BLUE}🚀 启动 API 服务 (main.ts)...${NC}"
DENO_KV_PATH=/app/.deno-kv/kv.db deno run \
    --allow-net \
    --allow-env \
    --allow-read \
    --allow-write=/app/.deno-kv \
    --unstable-kv \
    main.ts > "$API_LOG" 2>&1 &
API_PID=$!
echo -e "${GREEN}✓${NC} API 服务已启动 (PID: $API_PID)"
echo "   日志: $API_LOG"
echo "   访问: http://localhost:${PORT:-9090}"

# 等待 API 服务就绪
sleep 2

# 启动注册工具（后台运行）
echo ""
echo -e "${BLUE}🔧 启动账号注册工具 (zai_register.ts)...${NC}"
DENO_KV_PATH=/app/.deno-kv/kv.db deno run \
    --allow-net \
    --allow-env \
    --allow-read \
    --allow-write=/app/.deno-kv \
    --unstable-kv \
    zai_register.ts > "$REGISTER_LOG" 2>&1 &
REGISTER_PID=$!
echo -e "${GREEN}✓${NC} 注册工具已启动 (PID: $REGISTER_PID)"
echo "   日志: $REGISTER_LOG"
echo "   访问: http://localhost:8001"

# 等待注册工具就绪
sleep 2

echo ""
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}    所有服务已启动成功！${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "${BLUE}📌 服务地址:${NC}"
echo "  • API 文档:    http://localhost:${PORT:-9090}/docs"
echo "  • Dashboard:   http://localhost:${PORT:-9090}/dashboard"
echo "  • Admin 面板:  http://localhost:${PORT:-9090}/admin"
echo "  • 注册工具:    http://localhost:8001"
echo ""
echo -e "${BLUE}💾 KV 数据库:${NC}"
echo "  • 路径: /app/.deno-kv/kv.db"
echo "  • 两个服务共享同一数据库"
echo ""
echo -e "${BLUE}📊 查看日志:${NC}"
echo "  • API:    docker exec -it <container> tail -f $API_LOG"
echo "  • 注册:   docker exec -it <container> tail -f $REGISTER_LOG"
echo ""
echo -e "${YELLOW}按 Ctrl+C 停止所有服务${NC}"
echo ""

# 持续监控进程状态
while true; do
    # 检查 API 服务
    if ! kill -0 $API_PID 2>/dev/null; then
        echo -e "${RED}✗${NC} API 服务异常退出"
        echo "最后日志:"
        tail -n 20 "$API_LOG"
        kill $REGISTER_PID 2>/dev/null || true
        exit 1
    fi
    
    # 检查注册工具
    if ! kill -0 $REGISTER_PID 2>/dev/null; then
        echo -e "${RED}✗${NC} 注册工具异常退出"
        echo "最后日志:"
        tail -n 20 "$REGISTER_LOG"
        kill $API_PID 2>/dev/null || true
        exit 1
    fi
    
    sleep 5
done

