# ================================
# 多阶段构建 - KeepBlog
# ================================

# 构建参数
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# ================================
# 阶段0: 构建前端（管理后台 console）
# 产物拷贝到 Go builder 的 static/console 供 embed
# ================================
FROM node:20-alpine AS frontend
RUN corepack enable && corepack prepare pnpm@9.15.9 --activate
WORKDIR /fe
# 先复制依赖描述与 pnpm 配置，利用 Docker 层缓存
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/.npmrc frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
# 复制前端源码并构建
COPY frontend/ .
RUN pnpm build

# ================================
# 阶段1: 构建 Go
# ================================
FROM golang:1.24.11-alpine AS builder

# 构建参数
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

# 安装 CGO 依赖（SQLite 需要）
RUN apk add --no-cache --update gcc musl-dev sqlite-dev

# 设置环境变量（启用 CGO）
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=1

WORKDIR /app

# 先复制依赖文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 用前端阶段产物覆盖 static/console（保留 static.go，供 embed）
COPY --from=frontend /fe/dist/. static/console/

# 版本注入编译
RUN go build -trimpath \
    -ldflags="-s -w \
    -X 'gitee.com/jieepre/keepblog/internal/version.Version=${VERSION}' \
    -X 'gitee.com/jieepre/keepblog/internal/version.GitCommit=${GIT_COMMIT}' \
    -X 'gitee.com/jieepre/keepblog/internal/version.BuildTime=${BUILD_DATE}' \
    -X 'gitee.com/jieepre/keepblog/internal/version.GoVersion=$(go version | cut -d\" \" -f3)'" \
    -o keepblog .

# ================================
# 阶段2: 运行
# ================================
FROM alpine:3.19

# 安装运行时依赖（SQLite 需要）
RUN apk --no-cache add ca-certificates tzdata sqlite-libs bash && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建非root用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /app/keepblog /app/keepblog
COPY --from=builder /app/config-example.yaml /app/default/config.yaml

# 创建数据目录并设置权限
RUN mkdir -p /app/data /app/logs /app/config && \
    chown -R appuser:appgroup /app

# 入口点脚本（转换Windows换行符）
COPY entrypoint.sh /app/entrypoint.sh
RUN sed -i 's/\r$//' /app/entrypoint.sh && chmod +x /app/entrypoint.sh && chown appuser:appgroup /app/entrypoint.sh

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8589

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8589/health/live || exit 1

# 启动
ENTRYPOINT ["/app/entrypoint.sh"]
