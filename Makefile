GC = go build
APP_NAME = go-site
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION = $(shell go version | cut -d ' ' -f 3)

# 版本注入参数
LDFLAGS = -s -w \
	-X 'gitee.com/jieepre/go-site/internal/version.Version=$(VERSION)' \
	-X 'gitee.com/jieepre/go-site/internal/version.GitCommit=$(GIT_COMMIT)' \
	-X 'gitee.com/jieepre/go-site/internal/version.BuildTime=$(BUILD_TIME)' \
	-X 'gitee.com/jieepre/go-site/internal/version.GoVersion=$(GO_VERSION)'

BUILD_FLAGS = -trimpath -ldflags "$(LDFLAGS)"

.PHONY: all build linux darwin windows clean run test docker version help

# 默认目标
all: build

# 本地构建
build:
	$(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME) main.go

# Linux 构建
linux:
	GOOS=linux GOARCH=amd64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-linux-amd64 main.go

# Linux ARM64 构建
linux-arm64:
	GOOS=linux GOARCH=arm64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-linux-arm64 main.go

# macOS 构建
darwin:
	GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-darwin-amd64 main.go

# macOS ARM64 构建
darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-darwin-arm64 main.go

# Windows 构建
windows:
	GOOS=windows GOARCH=amd64 $(GC) $(BUILD_FLAGS) -o bin/$(APP_NAME)-windows-amd64.exe main.go

# 清理构建产物
clean:
	rm -rf bin/*

# 运行项目
run:
	go run main.go

# 运行测试
test:
	go test -v ./...

# 构建Docker镜像
docker:
	docker build -t jieepre/$(APP_NAME):$(VERSION) -t jieepre/$(APP_NAME):latest .

# 同步前端到后端
console:
	@echo "构建前端并同步到后端..."
	@bash sync-console.sh

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
	@echo "  make build        - 本地构建"
	@echo "  make linux        - Linux AMD64 构建"
	@echo "  make linux-arm64  - Linux ARM64 构建"
	@echo "  make darwin       - macOS AMD64 构建"
	@echo "  make darwin-arm64 - macOS ARM64 构建"
	@echo "  make windows      - Windows AMD64 构建"
	@echo "  make clean        - 清理构建产物"
	@echo "  make run          - 运行项目"
	@echo "  make test         - 运行测试"
	@echo "  make docker       - 构建Docker镜像"
	@echo "  make docker-push  - 推送Docker镜像"
	@echo "  make version      - 显示版本信息"
