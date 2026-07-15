GC = go build
APP_NAME = go-site
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION = $(shell go version | cut -d ' ' -f 3)

# 前端目录与构建产物落地目录（供 Go embed）
FRONTEND_DIR := frontend
CONSOLE_DIR := static/console

# 设 SKIP_FRONTEND=1 可跳过前端构建（前端已由别处构建并同步到 static/console）
# 用法示例：make linux SKIP_FRONTEND=1
ifeq ($(SKIP_FRONTEND),1)
FRONTEND_DEP :=
else
FRONTEND_DEP := build-frontend
endif

# 版本注入参数
LDFLAGS = -s -w \
	-X 'gitee.com/jieepre/go-site/internal/version.Version=$(VERSION)' \
	-X 'gitee.com/jieepre/go-site/internal/version.GitCommit=$(GIT_COMMIT)' \
	-X 'gitee.com/jieepre/go-site/internal/version.BuildTime=$(BUILD_TIME)' \
	-X 'gitee.com/jieepre/go-site/internal/version.GoVersion=$(GO_VERSION)'

BUILD_FLAGS = -trimpath -ldflags "$(LDFLAGS)"

.PHONY: all build build-frontend build-go linux linux-arm64 darwin darwin-arm64 windows clean clean-frontend run test docker console version help dev-frontend dev-backend

# 默认目标
all: build

# ================================================
# 前端（管理后台 console）
# 要求：pnpm >=9、Node ^20.19 || >=22.13
# ================================================

# 构建前端并把产物同步到 static/console（保留 static.go，供 Go embed）
build-frontend:
	@echo "==> 构建前端 $(FRONTEND_DIR)..."
	cd $(FRONTEND_DIR) && pnpm install --frozen-lockfile && pnpm build
	@echo "==> 同步产物到 $(CONSOLE_DIR)（保留 static.go）..."
	find $(CONSOLE_DIR) -mindepth 1 -not -name "static.go" -delete 2>/dev/null || true
	cp -r $(FRONTEND_DIR)/dist/. $(CONSOLE_DIR)/

# 仅构建 Go（不重建前端，使用 static/console 中已有产物）
build-go:
	$(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME) main.go

# 本地构建（前端 + Go 一键，确保二进制内嵌最新前端）
build: $(FRONTEND_DEP) build-go

# Linux 构建
linux: $(FRONTEND_DEP)
	GOOS=linux GOARCH=amd64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-linux-amd64 main.go

# Linux ARM64 构建
linux-arm64: $(FRONTEND_DEP)
	GOOS=linux GOARCH=arm64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-linux-arm64 main.go

# macOS 构建
darwin: $(FRONTEND_DEP)
	GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-darwin-amd64 main.go

# macOS ARM64 构建
darwin-arm64: $(FRONTEND_DEP)
	GOOS=darwin GOARCH=arm64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-darwin-arm64 main.go

# Windows 构建
windows: $(FRONTEND_DEP)
	GOOS=windows GOARCH=amd64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-windows-amd64.exe main.go

# 清理构建产物
clean:
	rm -rf bin/*

# 清理前端构建产物（dist 与依赖）
clean-frontend:
	rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules

# 运行项目（后端，本地开发）
run:
	go run main.go

# 运行测试
test:
	go test -v ./...

# 构建Docker镜像（Dockerfile 内已包含前端构建阶段）
docker:
	docker build -t jieepre/$(APP_NAME):$(VERSION) -t jieepre/$(APP_NAME):latest .

# 同步前端到后端（= build-frontend）
console: build-frontend

# ================================================
# 开发模式
# 前端: make dev-frontend  → vite dev server @ 8848（HMR，/api 自动转发到 8589）
# 后端: make dev-backend    → go run @ 8589
# 浏览器访问 http://localhost:8848
# ================================================
dev-frontend:
	cd $(FRONTEND_DIR) && pnpm dev

dev-backend:
	go run main.go

# 推送Docker镜像
docker-push:
	docker push jieepre/$(APP_NAME):$(VERSION)
	docker push jieepre/$(APP_NAME):latest

# 显示版本信息
version:
	@echo "Version:    $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

# 帮助信息
help:
	@echo "Go-Site 构建命令:"
	@echo "  make build           - 本地构建（前端 + Go 一键）"
	@echo "  make build-frontend  - 仅构建前端并同步产物到 static/console"
	@echo "  make build-go        - 仅构建 Go（使用已有前端产物）"
	@echo "  make linux           - Linux AMD64 构建"
	@echo "  make linux-arm64     - Linux ARM64 构建"
	@echo "  make darwin          - macOS AMD64 构建"
	@echo "  make darwin-arm64    - macOS ARM64 构建"
	@echo "  make windows         - Windows AMD64 构建"
	@echo "  make clean           - 清理 Go 构建产物"
	@echo "  make clean-frontend  - 清理前端 dist 与 node_modules"
	@echo "  make run             - 运行项目（后端）"
	@echo "  make test            - 运行测试"
	@echo "  make docker          - 构建Docker镜像"
	@echo "  make docker-push     - 推送Docker镜像"
	@echo "  make console         - 同步前端到后端（= build-frontend）"
	@echo "  make version         - 显示版本信息"
	@echo "开发命令:"
	@echo "  make dev-frontend    - 启动前端 dev server（vite @ 8848，HMR）"
	@echo "  make dev-backend     - 启动后端（go run @ 8589）"
