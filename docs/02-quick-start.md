# 快速开始

本指南将帮助您快速安装、配置和运行 KeepBlog 项目。

---

## 环境要求

### 必需
- **Go**: 1.24.11 或更高版本
- **Git**: 用于克隆项目
- **Make**: 用于构建命令（可选）

### 可选
- **Redis**: 用于缓存（可选，默认禁用）
- **Minio**: 用于对象存储（可选）
- **Docker**: 用于容器化部署

---

## 安装步骤

### 1. 克隆项目

```bash
git clone <repository-url>
cd keepblog
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置文件

复制配置示例文件：

```bash
cp config-example.yaml config.yaml
```

编辑 `config.yaml` 配置文件：

```yaml
http:
  port: 8589                    # HTTP 服务端口

redis:
  host: localhost
  port: 6379
  password: ""
  database: 0
  enable: false                 # 是否启用 Redis

minio:
  serverUrl: http://localhost:9000
  endpoint: localhost:9000
  accessKeyID: minioadmin
  secretAccessKey: minioadmin
  useSSL: false
  bucketName: test

baidu:
  push: false                   # 是否启用百度推送
  url: https://www.example.com
  token: ""

gitalk:
  enable: false                 # 是否启用 Gitalk 评论
  clientID: ""
  clientSecret: ""
  repo: ""
  owner: ""
  admin: []
```

### 4. 初始化数据库

首次运行时，应用会自动创建 SQLite 数据库并初始化数据：

```bash
# 数据库文件位置
./data/site.db
```

### 5. 运行项目

#### 使用 Make（推荐）

```bash
make run
```

#### 直接运行

```bash
go run main.go
```

### 6. 访问应用

- **前台**: http://localhost:8589
- **后台**: http://localhost:8589/console
- **健康检查**: http://localhost:8589/health
- **性能指标**: http://localhost:8589/metrics

---

## 默认账户

### 管理员账户
- **用户名**: `admin`
- **密码**: `admin123`

⚠️ **重要**: 首次登录后请立即修改密码！

---

## 构建项目

### 本地构建

```bash
make build
```

构建产物位于 `./bin/` 目录。

### 多平台构建

```bash
# Linux AMD64
make linux

# Linux ARM64
make linux-arm64

# macOS AMD64
make darwin

# macOS ARM64
make darwin-arm64

# Windows AMD64
make windows
```

### 查看版本信息

```bash
make version
```

---

## Docker 部署

### 1. 构建镜像

```bash
make docker
```

或手动构建：

```bash
docker build -t keepblog:latest .
```

### 2. 使用 Docker Compose

创建 `docker-compose.yml`（项目已包含）：

```yaml
services:
  site:
    image: jieepre/site:v1
    restart: always
    container_name: site
    working_dir: /app
    volumes:
      - ./logs:/app/logs
      - ./config:/app/config
      - ./data:/app/data
    environment:
      TZ: Asia/Shanghai
    ports:
      - "8589:8589"
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

启动容器：

```bash
docker-compose up -d
```

查看日志：

```bash
docker-compose logs -f
```

停止容器：

```bash
docker-compose down
```

### 3. 直接运行容器

```bash
docker run -d \
  --name keepblog \
  -p 8589:8589 \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/config:/app/config \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  keepblog:latest
```

---

## 环境变量配置

除了配置文件，还可以使用环境变量覆盖配置：

```bash
# HTTP 端口
export KEEPBLOG_HTTP_PORT=8589

# Redis 配置
export KEEPBLOG_REDIS_HOST=localhost
export KEEPBLOG_REDIS_PORT=6379
export KEEPBLOG_REDIS_PASSWORD=
export KEEPBLOG_REDIS_DATABASE=0
export KEEPBLOG_REDIS_ENABLE=false

# 基础 URL
export KEEPBLOG_BASE_URL=https://example.com
```

---

## 开发模式

### 1. 启用热重载

使用 `air` 工具实现热重载：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 运行
air
```

### 2. 调试模式

设置 Gin 模式为 debug：

```go
// internal/config/server.go
Mode: "debug"  // 默认是 "release"
```

### 3. 查看日志

日志文件位于 `./logs/` 目录：

```bash
# 查看最新日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log
```

---

## 常用命令

### Make 命令

```bash
make build          # 构建项目
make run            # 运行项目
make test           # 运行测试
make clean          # 清理构建产物
make docker         # 构建 Docker 镜像
make docker-push    # 推送 Docker 镜像
make console        # 构建前端并同步
make version        # 显示版本信息
```

### Go 命令

```bash
# 运行项目
go run main.go

# 构建项目
go build -o bin/keepblog main.go

# 运行测试
go test ./...

# 查看依赖
go mod graph

# 更新依赖
go mod tidy
```

### Docker 命令

```bash
# 构建镜像
docker build -t keepblog:latest .

# 运行容器
docker run -d -p 8589:8589 keepblog:latest

# 查看日志
docker logs -f <container-id>

# 进入容器
docker exec -it <container-id> sh

# 停止容器
docker stop <container-id>

# 删除容器
docker rm <container-id>
```

---

## 目录结构

运行后会自动创建以下目录：

```
keepblog/
├── data/                   # 数据目录
│   └── site.db            # SQLite 数据库
├── logs/                   # 日志目录
│   ├── app.log            # 应用日志
│   └── error.log          # 错误日志
├── config/                 # 配置目录
│   └── config.yaml        # 配置文件
└── uploads/                # 上传文件目录（如果不使用 Minio）
```

---

## 验证安装

### 1. 健康检查

```bash
curl http://localhost:8589/health
```

预期响应：

```json
{
  "status": "healthy",
  "timestamp": "2026-01-28T10:00:00Z"
}
```

### 2. 就绪检查

```bash
curl http://localhost:8589/health/ready
```

### 3. 存活检查

```bash
curl http://localhost:8589/health/live
```

### 4. 性能指标

```bash
curl http://localhost:8589/metrics
```

### 5. 前台访问

浏览器访问: http://localhost:8589

### 6. 后台登录

浏览器访问: http://localhost:8589/console

使用默认账户登录：
- 用户名: `admin`
- 密码: `admin123`

---

## 可选配置

### 启用 Redis 缓存

1. 安装 Redis：

```bash
# macOS
brew install redis
brew services start redis

# Ubuntu
sudo apt-get install redis-server
sudo systemctl start redis

# Docker
docker run -d -p 6379:6379 redis:alpine
```

2. 修改配置文件：

```yaml
redis:
  host: localhost
  port: 6379
  password: ""
  database: 0
  enable: true          # 启用 Redis
```

### 启用 Minio 对象存储

1. 安装 Minio：

```bash
# Docker
docker run -d \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  -v $(pwd)/minio-data:/data \
  minio/minio server /data --console-address ":9001"
```

2. 创建存储桶：

访问 http://localhost:9001，登录后创建名为 `test` 的存储桶。

3. 修改配置文件：

```yaml
minio:
  serverUrl: http://localhost:9000
  endpoint: localhost:9000
  accessKeyID: minioadmin
  secretAccessKey: minioadmin
  useSSL: false
  bucketName: test
```

### 启用百度推送

1. 获取百度推送令牌：

访问 [百度站长平台](https://ziyuan.baidu.com/)，获取推送令牌。

2. 修改配置文件：

```yaml
baidu:
  push: true
  url: https://www.example.com
  token: "your-baidu-token"
```

### 启用 Gitalk 评论

1. 创建 GitHub OAuth 应用：

访问 https://github.com/settings/applications/new

2. 修改配置文件：

```yaml
gitalk:
  enable: true
  clientID: "your-client-id"
  clientSecret: "your-client-secret"
  repo: "your-repo"
  owner: "your-github-username"
  admin: ["your-github-username"]
```

---

## 故障排查

### 端口被占用

```bash
# 查看端口占用
lsof -i :8589

# 杀死进程
kill -9 <pid>

# 或修改配置文件中的端口
```

### 数据库连接失败

```bash
# 检查数据库文件是否存在
ls -la data/site.db

# 检查文件权限
chmod 644 data/site.db

# 删除数据库重新初始化
rm data/site.db
```

### Redis 连接失败

```bash
# 检查 Redis 是否运行
redis-cli ping

# 检查配置文件中的 Redis 配置
# 如果不需要 Redis，设置 enable: false
```

### 日志查看

```bash
# 查看应用日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log

# Docker 日志
docker logs -f <container-id>
```

### 权限问题

```bash
# 确保目录有写权限
chmod -R 755 data logs uploads

# Docker 容器权限
docker exec -it <container-id> ls -la /app
```

---

## 下一步

- [项目结构](./03-project-structure.md) - 了解项目目录结构
- [系统架构](./04-architecture.md) - 了解系统架构设计
- [API 设计](./06-api-design.md) - 查看 API 接口文档
- [配置管理](./18-configuration.md) - 详细配置说明

---

## 常见问题

### Q: 如何修改管理员密码？

A: 登录后台后，访问"个人中心"修改密码，或使用 API：

```bash
curl -X POST http://localhost:8589/api/change-password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"oldPassword":"admin123","newPassword":"new-password"}'
```

### Q: 如何备份数据？

A: 备份 SQLite 数据库文件：

```bash
cp data/site.db data/site.db.backup
```

### Q: 如何迁移数据？

A: 复制整个 `data/` 目录到新环境。

### Q: 如何更新项目？

A: 拉取最新代码并重新构建：

```bash
git pull
go mod download
make build
```

### Q: 如何查看版本信息？

A: 使用 Make 命令：

```bash
make version
```

或运行二进制文件：

```bash
./bin/keepblog --version
```

---

**相关文件**:
- `config-example.yaml` - 配置示例
- `Makefile` - 构建命令
- `Dockerfile` - Docker 构建
- `docker-compose.yml` - Docker 编排
