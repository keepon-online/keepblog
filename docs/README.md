# Go-Site 项目文档索引

> 全功能博客系统 (CMS) - 基于 Go + Gin + SQLite

## 📚 文档导航

### 快速开始
- [项目概览](./01-overview.md) - 项目介绍、技术栈、核心特性
- [快速开始](./02-quick-start.md) - 安装、配置、运行指南
- [项目结构](./03-project-structure.md) - 目录结构详解

### 架构设计
- [系统架构](./04-architecture.md) - 整体架构、分层设计、数据流
- [数据模型](./05-data-models.md) - 数据库表结构、关系图
- [API 设计](./06-api-design.md) - RESTful API 规范、端点列表

### 核心模块
- [前台功能](./07-frontend-features.md) - 首页、文章、分类、标签、归档
- [后台管理](./08-admin-features.md) - 内容管理、系统配置、监控
- [认证授权](./09-authentication.md) - JWT 认证、权限控制
- [中间件系统](./10-middleware.md) - 日志、限流、缓存、安全

### 服务层
- [服务架构](./11-service-layer.md) - 服务设计、依赖注入
- [文章服务](./12-post-service.md) - 文章 CRUD、发布流程
- [分类标签服务](./13-category-tag-service.md) - 分类和标签管理
- [系统服务](./14-system-service.md) - 配置、日志、监控

### 工具包
- [公共工具](./15-utilities.md) - 工具函数、辅助类
- [Markdown 处理](./16-markdown.md) - Markdown 渲染、代码高亮
- [对象存储](./17-object-storage.md) - Minio 集成、文件上传

### 部署运维
- [配置管理](./18-configuration.md) - 配置文件、环境变量
- [构建部署](./19-build-deploy.md) - 编译、Docker、CI/CD
- [监控日志](./20-monitoring.md) - 系统监控、日志管理
- [性能优化](./21-performance.md) - 缓存策略、数据库优化

### 开发指南
- [开发规范](./22-development-guide.md) - 代码规范、最佳实践
- [测试指南](./23-testing.md) - 单元测试、集成测试
- [故障排查](./24-troubleshooting.md) - 常见问题、调试技巧

### API 参考
- [前台 API](./api/web-api.md) - 前台接口文档
- [后台 API](./api/admin-api.md) - 后台接口文档
- [WebSocket API](./api/websocket-api.md) - 实时通信接口

---

## 🎯 项目概览

### 技术栈
- **语言**: Go 1.24.11
- **框架**: Gin 1.11.0
- **数据库**: SQLite (GORM 1.31.1)
- **缓存**: Redis 9.6.1 (可选)
- **认证**: JWT (golang-jwt/jwt v5.2.2)
- **对象存储**: Minio 7.0.63
- **日志**: Slog (gookit/slog 0.6.0)

### 核心特性
- ✅ 文章管理 (发布、编辑、分类、标签)
- ✅ Markdown 编辑器 (代码高亮、TOC)
- ✅ 分类和标签系统
- ✅ 友情链接管理
- ✅ 音乐播放器
- ✅ 系统监控 (CPU、内存、磁盘、网络)
- ✅ 访问统计和日志
- ✅ JWT 认证和权限控制
- ✅ Redis 缓存支持
- ✅ Minio 对象存储
- ✅ WebSocket 实时通信
- ✅ 百度收录推送
- ✅ SEO 优化
- ✅ Docker 容器化

### 项目统计
- **Go 文件**: 112 个
- **API 端点**: 80+ 个
- **数据表**: 13 个
- **中间件**: 10 个
- **服务模块**: 12 个

---

## 🚀 快速链接

### 开发
```bash
# 克隆项目
git clone <repository-url>

# 安装依赖
go mod download

# 运行项目
make run

# 构建项目
make build
```

### Docker
```bash
# 构建镜像
make docker

# 运行容器
docker-compose up -d
```

### 访问
- **前台**: http://localhost:8589
- **后台**: http://localhost:8589/console
- **健康检查**: http://localhost:8589/health
- **性能指标**: http://localhost:8589/metrics

---

## 📖 文档约定

### 代码引用格式
文档中的代码位置使用 `file_path:line_number` 格式标注，例如：
- `internal/app/app.go:45` - 应用初始化
- `api/admin/post/post.go:123` - 文章保存接口

### 目录结构图标
- 📁 目录
- 📄 文件
- 🔧 配置文件
- 🐳 Docker 相关
- 📝 文档文件

### 重要性标记
- ⚠️ 重要提示
- 💡 最佳实践
- 🔒 安全相关
- ⚡ 性能优化
- 🐛 已知问题

---

## 🤝 贡献指南

请参阅 [开发规范](./22-development-guide.md) 了解如何为项目做出贡献。

---

## 📝 更新日志

查看 [CHANGELOG.md](../CHANGELOG.md) 了解版本更新历史。

---

## 📄 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

---

**文档生成时间**: 2026-01-28
**项目版本**: 基于 master 分支 (commit: d520557)
