# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install SQLite development libraries for CGO
RUN apk --no-cache add sqlite-dev gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Enable CGO for SQLite support
RUN CGO_ENABLED=1 go build -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /app
COPY --from=builder /app/main .

# Create data directory for SQLite database
RUN mkdir -p /app/data
VOLUME ["/app/data"]

# Environment variables
# 基础配置
ENV DEFAULT_KEY=sk-your-key
ENV PORT=9090
ENV UPSTREAM_URL=https://chat.z.ai/api/chat/completions

# 功能开关
ENV DEBUG_MODE=true
ENV DEFAULT_STREAM=true
ENV ENABLE_THINKING=true
ENV DASHBOARD_ENABLED=true

# 管理系统配置
ENV REGISTER_ENABLED=true
ENV REGISTER_DB_PATH=/app/data/zai2api.db
ENV ADMIN_ENABLED=true
ENV ADMIN_USERNAME=admin
ENV ADMIN_PASSWORD=123456
ENV ZAI_USERNAME=admin
ENV ZAI_PASSWORD=123456

# Labels
LABEL maintainer="ZtoApi Contributors"
LABEL description="ZAI to GLM-4.6 API with Registration System"
LABEL version="2.0.0"

# Expose port
EXPOSE 9090

# Run the application
CMD ["./main"]