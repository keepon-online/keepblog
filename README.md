# Go-Site 博客系统

[![Go](https://img.shields.io/github/go-mod/go-version/username/go-site)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 项目简介

Go-Site 是一个基于 Go 语言和 Gin 框架开发的全功能博客系统，提供完整的前后台管理功能。

## 功能特性

### 前台功能
- 文章展示（分类、标签、归档）
- 文章搜索
- 评论系统
- 友情链接
- 响应式设计

### 后台管理
- 文章管理（发布、编辑、删除）
- 分类/标签管理
- 评论审核
- 系统监控
- 用户管理

## 技术栈

- **后端框架**: Gin
- **数据库**: SQLite (通过 GORM)
- **身份验证**: JWT
- **Markdown处理**: Goldmark + Lute
- **对象存储**: Minio
- **定时任务**: Gocron
- **系统监控**: Gopsutil

## 安装指南

### 前置要求
- Go 1.21+
- SQLite3

### 安装步骤
1. 克隆仓库：
   ```bash
   git clone https://github.com/username/go-site.git
   cd go-site
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 配置项目：
   复制 `config.example.yaml` 为 `config.yaml` 并修改配置

4. 启动项目：
   ```bash
   go run main.go
   ```

## 配置说明

主要配置项 (`config.yaml`)：

```yaml
server:
  port: 8589  # 前端端口
  admin_port: 8000  # 后台端口
  console_port: 8890  # 控制台端口

database:
  path: "data/site.db"  # SQLite数据库路径

jwt:
  secret: "your-secret-key"  # JWT密钥
  expire: 720h  # 过期时间
```

## 使用说明

### 启动服务
项目启动后会运行三个服务：
1. 前端网站: http://localhost:8589
2. 管理后台: http://localhost:8000
3. 控制台: http://localhost:8890

### 管理员账号
默认管理员账号：
- 用户名: admin
- 密码: admin123

首次登录后请立即修改密码。

## API文档

[查看API文档](docs/api.md)

## 贡献指南

欢迎贡献代码！请遵循以下步骤：
1. Fork 本项目
2. 创建分支 (`git checkout -b feature/your-feature`)
3. 提交修改 (`git commit -am 'Add some feature'`)
4. 推送分支 (`git push origin feature/your-feature`)
5. 创建 Pull Request

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件