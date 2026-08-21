# 项目概览

## 项目简介

Go-Site 是一个基于 Go 语言开发的全功能博客系统 (CMS)，采用现代化的技术栈和架构设计，提供完整的内容管理、系统监控和性能优化功能。

## 核心特性

### 内容管理
- **文章系统**: 支持 Markdown 编辑、代码高亮、TOC 生成
- **分类管理**: 多级分类、状态控制
- **标签系统**: 标签云、标签页面、样式自定义
- **评论功能**: Gitalk 集成（可选）
- **友情链接**: 链接管理、状态控制
- **音乐播放**: 音乐列表、播放器集成

### 系统功能
- **用户认证**: JWT 令牌、刷新机制、黑名单
- **权限控制**: 基于路由的权限管理
- **系统监控**: CPU、内存、磁盘、网络实时监控
- **日志管理**: 结构化日志、级别控制、日志查看
- **访问统计**: PV/UV 统计、IP 地理位置、用户代理分析
- **性能指标**: Prometheus 格式的性能指标

### 技术特性
- **高性能**: 连接池、缓存、压缩、限流
- **安全性**: JWT、CSRF、SQL 注入防护、密码加密
- **可扩展**: 模块化设计、依赖注入、显式依赖注入
- **容器化**: Docker 支持、多平台构建
- **实时通信**: WebSocket 支持
- **对象存储**: Minio 集成

## 技术栈

### 后端框架
| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.24.11 | 编程语言 |
| Gin | 1.11.0 | Web 框架 |
| GORM | 1.31.1 | ORM 框架 |
| SQLite | - | 数据库 |
| Redis | 9.6.1 | 缓存（可选） |

### 核心依赖
| 依赖 | 版本 | 用途 |
|------|------|------|
| golang-jwt/jwt | 4.5.2 | JWT 认证 |
| gookit/slog | 0.6.0 | 结构化日志 |
| minio-go | 7.0.63 | 对象存储 |
| gocron | 1.36.0 | 定时任务 |
| goldmark | 1.7.13 | Markdown 渲染 |
| gopsutil | 3.23.10 | 系统监控 |
| viper | 1.17.0 | 配置管理 |
| bcrypt | - | 密码加密 |
| hashids | - | ID 混淆 |

### 工具链
- **构建工具**: Make, Go Build
- **容器化**: Docker, Docker Compose
- **版本控制**: Git
- **包管理**: Go Modules

## 架构特点

### 分层架构
```
┌─────────────────────────────────────┐
│         API Layer (Handlers)        │  ← HTTP 请求处理
├─────────────────────────────────────┤
│       Service Layer (Business)      │  ← 业务逻辑
├─────────────────────────────────────┤
│        Model Layer (Data)           │  ← 数据访问
├─────────────────────────────────────┤
│      Database (SQLite/Redis)        │  ← 数据存储
└─────────────────────────────────────┘
```

### 模块化设计
- **API 模块**: 前台 (web) + 后台 (admin)
- **服务模块**: 12 个独立服务
- **中间件**: 10 个可组合中间件
- **工具包**: 20+ 个公共工具

### 依赖注入
```go
// 服务聚合
type AppService struct {
    PostService     *post.Service
    CategoryService *category.Service
    TagService      *tag.Service
    // ... 其他服务
}

// 请求上下文
type Context struct {
    Engine  *gin.Engine
    Service *AppService
}
```

## 项目结构

```
go-site/
├── api/                    # API 处理器层
│   ├── admin/             # 后台管理 API (14 个模块)
│   └── web/               # 前台展示 API (8 个模块)
├── internal/              # 内部包
│   ├── app/               # 应用入口和路由
│   ├── service/           # 业务逻辑层 (12 个服务)
│   ├── model/             # 数据模型
│   ├── middleware/        # 中间件 (10 个)
│   ├── core/              # 核心初始化
│   └── pkg/               # 内部工具包
├── pkg/                   # 公共工具包 (20+ 个)
├── config/                # 配置管理
├── static/                # 静态资源
├── global/                # 全局变量
└── main.go                # 应用入口
```

详细结构请参阅 [项目结构](./03-project-structure.md)。

## 核心流程

### 应用启动流程
```
1. 加载配置 (config.yaml + 环境变量)
   ↓
2. 初始化日志系统
   ↓
3. 连接数据库 (SQLite)
   ↓
4. 初始化 Redis (可选)
   ↓
5. 自动迁移数据表
   ↓
6. 创建初始数据
   ↓
7. 启动定时任务
   ↓
8. 注册路由和中间件
   ↓
9. 启动 HTTP 服务器
```

参考: `main.go:1`, `internal/app/app.go:1`

### 请求处理流程
```
HTTP 请求
   ↓
全局中间件 (日志、恢复、CORS、压缩)
   ↓
路由匹配
   ↓
路由中间件 (限流、缓存、JWT)
   ↓
API 处理器
   ↓
服务层 (业务逻辑)
   ↓
模型层 (数据访问)
   ↓
数据库
   ↓
响应返回
```

参考: `internal/app/routers.go:1`, `internal/middleware/`

## 性能指标

### 数据库连接池
- **最大连接数**: 25
- **最小空闲连接**: 10
- **连接生命周期**: 2 小时
- **空闲超时**: 30 分钟

参考: `internal/core/db.go:1`

### 限流配置
- **前台限流**: 60 次/分钟 (基于 IP)
- **后台限流**: 100 次/分钟 (基于 IP)
- **登录限流**: 5 次/分钟 (基于 IP)

参考: `internal/middleware/ratelimit.go:1`

### 缓存策略
- **前台页面缓存**: 5 分钟
- **Redis 缓存**: 可选启用
- **静态资源**: 浏览器缓存

参考: `internal/middleware/cache.go:1`

## 安全机制

### 认证授权
- **JWT 令牌**: 访问令牌 + 刷新令牌
- **令牌黑名单**: 登出时加入黑名单
- **密码加密**: bcrypt 加密存储

参考: [认证授权](./09-authentication.md)

### 安全防护
- **CSRF 防护**: 令牌验证
- **SQL 注入防护**: 参数化查询
- **XSS 防护**: 输出转义
- **路径规范化**: 防止路径遍历

参考: [中间件系统](./10-middleware.md)

## 监控能力

### 系统监控
- **CPU**: 使用率、核心数、型号
- **内存**: 使用率、总量、可用量
- **磁盘**: 使用率、IO 统计
- **网络**: 流量统计、连接数

参考: `internal/model/server/`, `api/admin/monitor/`

### 日志管理
- **日志级别**: DEBUG, INFO, WARN, ERROR
- **日志输出**: 文件 + 控制台
- **日志轮转**: 按大小和时间
- **日志查看**: Web 界面查看

参考: [监控日志](./20-monitoring.md)

### 性能指标
- **Prometheus 格式**: `/metrics` 端点
- **健康检查**: `/health` 端点
- **就绪检查**: `/health/ready` 端点
- **存活检查**: `/health/live` 端点

参考: `internal/monitor/`

## 部署方式

### 本地运行
```bash
# 安装依赖
go mod download

# 运行项目
make run

# 访问
http://localhost:8589
```

### Docker 部署
```bash
# 构建镜像
make docker

# 运行容器
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 多平台构建
```bash
make linux          # Linux AMD64
make linux-arm64    # Linux ARM64
make darwin         # macOS AMD64
make darwin-arm64   # macOS ARM64
make windows        # Windows AMD64
```

详细部署请参阅 [构建部署](./19-build-deploy.md)。

## 配置管理

### 配置文件
- **位置**: `config.yaml`
- **示例**: `config-example.yaml`
- **格式**: YAML

### 环境变量
支持通过环境变量覆盖配置:
```bash
GOSITE_HTTP_PORT=8589
GOSITE_REDIS_ENABLE=true
GOSITE_BASE_URL=https://example.com
```

详细配置请参阅 [配置管理](./18-configuration.md)。

## 扩展功能

### 百度推送
- **定时推送**: 每日 21:00
- **自动收录**: 新文章自动推送
- **配置**: `config.yaml` 中启用

### Gitalk 评论
- **GitHub 集成**: OAuth 认证
- **配置**: `config.yaml` 中配置
- **可选功能**: 可禁用

### Minio 对象存储
- **图片上传**: 支持 Minio 存储
- **配置**: `config.yaml` 中配置
- **可选功能**: 可使用本地存储

## 开发工具

### Make 命令
```bash
make build          # 构建项目
make run            # 运行项目
make test           # 运行测试
make clean          # 清理构建产物
make docker         # 构建 Docker 镜像
make console        # 构建前端
make version        # 显示版本信息
```

### 版本信息
编译时注入版本信息:
- **Version**: Git 标签或 "dev"
- **GitCommit**: 短提交哈希
- **BuildTime**: 构建时间戳
- **GoVersion**: Go 版本

参考: `Makefile:1`

## 项目统计

| 指标 | 数量 |
|------|------|
| Go 文件 | 112 个 |
| API 端点 | 80+ 个 |
| 数据表 | 13 个 |
| 中间件 | 10 个 |
| 服务模块 | 12 个 |
| 工具包 | 20+ 个 |
| 代码行数 | ~15,000 行 |

## 下一步

- [快速开始](./02-quick-start.md) - 安装和运行项目
- [项目结构](./03-project-structure.md) - 详细的目录结构
- [系统架构](./04-architecture.md) - 架构设计和数据流
- [API 设计](./06-api-design.md) - API 接口文档

---

**相关链接**:
- [GitHub 仓库](https://github.com/your-repo/go-site)
- [问题反馈](https://github.com/your-repo/go-site/issues)
- [更新日志](../CHANGELOG.md)
