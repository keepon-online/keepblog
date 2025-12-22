# Go-Site 博客系统

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](Dockerfile)

## 项目简介

Go-Site 是一个基于 Go 语言和 Gin 框架开发的全功能博客系统，提供完整的前后台管理功能。

## 功能特性

### 前台功能
- 文章展示（分类、标签、归档）
- 全文搜索（标题+内容+摘要）
- 评论系统
- 友情链接
- 响应式设计
- 音乐播放器

### 后台管理
- 文章管理（发布、编辑、删除）
- 分类/标签管理
- 评论审核
- 系统监控（CPU/内存/磁盘/网络）
- 日志管理（动态级别、查看、清理）
- 用户管理

## 技术栈

| 类别 | 技术 |
|------|------|
| 后端框架 | Gin |
| 数据库 | SQLite (GORM) |
| 缓存 | Redis (可选) |
| 身份验证 | JWT |
| 日志 | Slog (结构化日志) |
| 监控 | Gopsutil |
| 对象存储 | Minio |

## 快速开始

### 前置要求
- Go 1.21+
- SQLite3
- pnpm (前端构建)

### 本地开发

```bash
# 克隆项目
git clone https://github.com/username/go-site.git
cd go-site

# 安装依赖
go mod download

# 配置
cp config-example.yaml config.yaml

# 运行
make run
```

### 构建命令

```bash
make build         # 本地构建（带版本注入）
make linux         # Linux AMD64 构建
make linux-arm64   # Linux ARM64 构建
make darwin        # macOS AMD64 构建
make windows       # Windows AMD64 构建
make docker        # 构建 Docker 镜像
make console       # 构建前端并同步
make version       # 显示版本信息
make help          # 查看所有命令
```

### Docker 部署

```bash
# 构建镜像
./build.sh

# 或使用 make
make docker

# 运行容器
docker run -d -p 8589:8589 -v ./data:/app/data jieepre/go-site:latest
```

## 配置说明

主要配置项 (`config.yaml`)：

```yaml
http:
  port: "8589"

redis:
  host: localhost
  port: 6379
  enable: false  # 可选缓存

minio:
  endpoint: "127.0.0.1:9000"
  bucketName: "go-site"
```

## API 端点

| 端点 | 说明 |
|------|------|
| `/health` | 健康检查 |
| `/health/live` | 存活检查 |
| `/metrics` | 性能指标 |
| `/api/log/level` | 日志级别管理 |
| `/api/monitor/server` | 系统监控 |

## 项目结构

```
go-site/
├── api/              # API 处理器
├── config/           # 配置
├── internal/         # 内部包
│   ├── app/          # 应用入口
│   ├── cache/        # Redis 缓存
│   ├── logger/       # 日志系统
│   ├── middleware/   # 中间件
│   ├── monitor/      # 健康检查
│   └── version/      # 版本信息
├── pkg/              # 公共包
├── static/           # 静态资源
│   └── console/      # 后台前端
└── Makefile          # 构建脚本
```

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件