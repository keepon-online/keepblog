# KeepBlog 博客系统知识库

> 完整的项目文档索引，涵盖架构设计、API 参考、实现指南和开发规范。

---

## 📚 文档导航

### 核心文档

- [项目概览](#项目概览)
- [架构设计](#架构设计)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)

### 功能模块

- [前台功能](#前台功能)
- [后台管理](#后台管理)
- [API 文档](#api-文档)
- [数据模型](#数据模型)
- [中间件](#中间件)

### 开发指南

- [配置说明](#配置说明)
- [部署指南](#部署指南)
- [开发规范](#开发规范)

### 前端项目

- [前端项目概览](#前端项目)
- [技术栈](#技术栈-1)
- [项目结构](#项目结构-1)
- [配置文件](#配置文件-1)
- [核心功能](#核心功能)
- [开发指南](#开发指南-1)
- [与后端集成](#与后端集成)

---

## 项目概览

### 简介

KeepBlog 是一个基于 Go 语言和 Gin 框架开发的全功能博客系统，提供完整的前后台管理功能。项目采用现代化的架构设计，支持高并发访问，内置监控、日志、缓存等企业级特性。

### 核心特性

- ✅ **前后台分离**：完整的前台展示和后台管理系统
- ✅ **响应式设计**：适配多端访问
- ✅ **全文搜索**：支持标题、内容、摘要搜索
- ✅ **系统监控**：实时监控 CPU、内存、磁盘、网络
- ✅ **日志管理**：动态日志级别、日志查看、清理
- ✅ **对象存储**：集成 Minio 对象存储
- ✅ **缓存支持**：可选 Redis 缓存加速
- ✅ **身份认证**：JWT 令牌认证

### 项目统计

- **Go 文件数**：112+
- **代码行数**：约 10,000+
- **模块数量**：15+
- **API 端点**：80+

---

## 架构设计

### 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前台用户界面                           │
│  (HTML模板 + 静态资源 + JavaScript)                          │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP/WebSocket
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Gin Web 框架层                           │
│  路由分发 → 中间件链 → Handler 处理                         │
└──────────────────────────┬──────────────────────────────────┘
                           │
            ┌──────────────┼──────────────┐
            ▼              ▼              ▼
    ┌───────────┐  ┌───────────┐  ┌───────────┐
    │ Service层 │  │ Cache层   │  │ Monitor层  │
    │ 业务逻辑   │  │ Redis缓存  │  │ 健康检查   │
    └─────┬─────┘  └─────┬─────┘  └───────────┘
          │              │
          ▼              ▼
    ┌──────────────────────────────┐
    │         数据访问层 (GORM)      │
    └──────────────┬───────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
   ┌─────────┐          ┌─────────┐
   │ SQLite  │          │  Minio  │
   │ 数据库  │          │ 对象存储 │
   └─────────┘          └─────────┘
```

### 分层架构

#### 1. 表现层 (Presentation Layer)

**位置**：`api/`

- **前台 API**：`api/web/` - 页面渲染和内容展示
- **后台 API**：`api/admin/` - 管理功能接口
- **职责**：
  - HTTP 请求处理
  - 参数验证和绑定
  - 响应格式化
  - WebSocket 连接处理

#### 2. 业务逻辑层 (Service Layer)

**位置**：`internal/service/`

- **核心服务**：
  - `PostService` - 文章管理
  - `CategoryService` - 分类管理
  - `TagService` - 标签管理
  - `SidebarService` - 侧边栏数据
  - `AboutService` - 关于页面
  - `SystemService` - 系统管理
  - `LinkService` - 友情链接
  - `WebSiteService` - 网站配置
  - `Dashboard` - 仪表盘
  - `MusicService` - 音乐管理

- **职责**：
  - 业务逻辑实现
  - 数据转换和聚合
  - 事务管理
  - 缓存策略

#### 3. 数据访问层 (Data Access Layer)

**位置**：`internal/model/` + `internal/core/`

- **ORM**：使用 GORM 进行数据库操作
- **数据库**：SQLite
- **职责**：
  - 数据模型定义
  - CRUD 操作封装
  - 数据库连接管理
  - 查询优化

#### 4. 基础设施层 (Infrastructure Layer)

**位置**：`internal/` 各子目录

- **配置**：`config/` - 配置加载和验证
- **缓存**：`cache/` - Redis 缓存实现
- **日志**：`logger/` - 结构化日志
- **中间件**：`middleware/` - 跨切面关注点
- **监控**：`monitor/` - 健康检查和指标
- **工具**：`pkg/` - 公共工具函数

### 设计模式

#### 1. 分层架构模式

清晰的职责分离：
```
Handler (表现层) → Service (业务层) → Model (数据层)
```

#### 2. 依赖注入模式

通过结构体组合注入依赖：
```go
type Handler struct {
    Context *core.Context  // 注入上下文
}
```

#### 3. 中间件链模式

使用 Gin 中间件实现横切关注点：
```
Request → Logger → Recovery → CORS → RateLimit → JWT → Handler → Response
```

#### 4. 工厂模式

服务实例创建：
```go
func InitAppService() *AppService {
    return &AppService{
        PostService: post.NewPostService(),
        // ...
    }
}
```

---

## 技术栈

### 后端技术

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **语言** | Go | 1.24+ | 主要开发语言 |
| **框架** | Gin | 1.11.0 | HTTP Web 框架 |
| **ORM** | GORM | 1.31.1 | 数据库操作 |
| **数据库** | SQLite | 3.x | 关系型数据库 |
| **缓存** | Redis | 9.x | 分布式缓存（可选） |
| **认证** | JWT | 5.2.2 | 身份认证 |
| **日志** | Slog | 0.6.0 | 结构化日志 |
| **监控** | Gopsutil | 3.23.10 | 系统监控 |
| **对象存储** | Minio | 7.0.63 | 文件存储 |
| **配置** | Viper | 1.17.0 | 配置管理 |
| **Markdown** | Goldmark | 1.7.13 | Markdown 解析 |
| **压缩** | Gzip | 0.0.6 | 响应压缩 |

### 前端技术

| 组件 | 技术 | 用途 |
|------|------|------|
| **模板引擎** | Go Template | 服务端渲染 |
| **样式** | CSS3 | 页面样式 |
| **脚本** | JavaScript | 交互逻辑 |
| **后台管理** | Console (静态资源) | SPA 管理界面 |

### 开发工具

| 工具 | 用途 |
|------|------|
| **Docker** | 容器化部署 |
| **Make** | 构建自动化 |
| **Git** | 版本控制 |
| **Drone CI** | 持续集成 |

---

## 项目结构

```
keepblog/
├── api/                      # API 处理器层
│   ├── admin/                # 后台管理 API
│   │   ├── about/           # 关于管理
│   │   ├── category/        # 分类管理
│   │   ├── common/          # 公共接口
│   │   ├── dashboard/       # 仪表盘
│   │   ├── link/            # 友情链接
│   │   ├── log/             # 日志查看
│   │   ├── login/           # 登录认证
│   │   ├── monitor/         # 系统监控
│   │   ├── music/           # 音乐管理
│   │   ├── notice/          # 通知管理
│   │   ├── post/            # 文章管理
│   │   ├── system/          # 系统管理
│   │   ├── tags/            # 标签管理
│   │   └── website/         # 网站配置
│   └── web/                  # 前台展示 API
│       ├── about/           # 关于页面
│       ├── archive/         # 归档页面
│       ├── category/        # 分类页面
│       ├── home/            # 首页
│       ├── link/            # 友情链接
│       ├── music/           # 音乐播放
│       ├── post/            # 文章详情
│       └── tags/            # 标签页面
├── config/                   # 配置模块
│   ├── config.go           # 配置加载
│   └── config_test.go      # 配置测试
├── data/                     # 运行时数据
│   └── site.db             # SQLite 数据库文件
├── docs/                     # 项目文档
├── global/                   # 全局变量
│   └── global.go           # 全局定义
├── internal/                 # 内部包（不对外暴露）
│   ├── app/                # 应用核心
│   │   ├── app.go         # 应用主结构
│   │   └── routers.go     # 路由配置
│   ├── cache/              # 缓存层
│   │   └── redis.go       # Redis 实现
│   ├── config/             # 内部配置
│   │   └── server.go      # 服务器配置
│   ├── core/               # 核心功能
│   │   ├── db.go          # 数据库初始化
│   │   ├── init.go        # 初始化
│   │   └── indexes.go     # 数据库索引
│   ├── errors/             # 错误处理
│   │   └── errors.go      # 错误定义
│   ├── logger/             # 日志系统
│   │   └── logger.go      # 日志配置
│   ├── middleware/         # 中间件
│   │   ├── access.go      # 访问日志
│   │   ├── cache.go       # 缓存中间件
│   │   ├── cors.go        # 跨域处理
│   │   ├── csrf.go        # CSRF 防护
│   │   ├── error.go       # 错误处理
│   │   ├── jwt.go         # JWT 认证
│   │   ├── log.go         # 日志中间件
│   │   ├── memorylimit.go # 内存限制
│   │   ├── print.go       # 请求打印
│   │   └── ratelimit.go   # 限流中间件
│   ├── model/              # 数据模型
│   │   ├── dashboard.go   # 仪表盘模型
│   │   ├── music.go       # 音乐模型
│   │   ├── post.go        # 文章模型
│   │   ├── request/       # 请求 DTO
│   │   ├── response/      # 响应 DTO
│   │   ├── server/        # 服务器模型
│   │   └── system/        # 系统模型
│   ├── monitor/            # 监控模块
│   │   └── health.go      # 健康检查
│   ├── pkg/                # 内部包
│   │   └── core/          # 核心类型
│   ├── router/             # 路由注册
│   │   ├── admin.go       # 后台路由
│   │   ├── route.go       # 路由工具
│   │   └── web.go         # 前台路由
│   ├── service/            # 业务逻辑层
│   │   ├── about/         # 关于服务
│   │   ├── category/      # 分类服务
│   │   ├── dashboard/     # 仪表盘服务
│   │   ├── link/          # 友情链接服务
│   │   ├── music/         # 音乐服务
│   │   ├── notice/        # 通知服务
│   │   ├── post/          # 文章服务
│   │   ├── sidebar/       # 侧边栏服务
│   │   ├── system/        # 系统服务
│   │   ├── tag/           # 标签服务
│   │   └── website/       # 网站服务
│   ├── version/            # 版本信息
│   │   └── version.go     # 版本号
│   ├── web/                # Web 资源
│   │   └── templates/     # HTML 模板
│   └── websocket/          # WebSocket
│       └── handler.go     # WebSocket 处理
├── pkg/                      # 公共工具包
│   ├── area/               # 地区工具
│   ├── cloudtag/           # 云标签
│   ├── daily/              # 日常工具
│   ├── file/               # 文件工具
│   ├── hash/               # 哈希工具
│   ├── iputil/             # IP 工具
│   ├── jwttoken/           # JWT 工具
│   ├── md/                 # Markdown 工具
│   ├── page/               # 分页工具
│   ├── psutil/             # 性能工具
│   ├── tags-remove/        # 标签处理
│   ├── ua/                 # User-Agent 解析
│   └── xpack/              # 扩展包
├── static/                   # 静态资源
│   ├── console/            # 后台管理界面
│   ├── css/                # 样式文件
│   ├── images/             # 图片资源
│   ├── js/                 # JavaScript 文件
│   ├── plugins/            # 插件
│   ├── favicon.ico         # 网站图标
│   └── robots.txt          # 爬虫配置
├── CHANGELOG.md            # 更新日志
├── Dockerfile              # Docker 镜像
├── Makefile                # 构建脚本
├── README.md               # 项目说明
├── build.sh                # 构建脚本
├── config-example.yaml     # 配置示例
├── config.yaml             # 配置文件
├── docker-compose.yml      # Docker Compose
├── entrypoint.sh           # 容器入口
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖锁定
└── main.go                 # 程序入口
```

### 目录说明

| 目录 | 说明 |
|------|------|
| `api/` | HTTP 处理器层，处理 HTTP 请求和响应 |
| `config/` | 配置加载和验证 |
| `internal/` | 内部代码，不对外暴露 |
| `pkg/` | 可复用的公共工具包 |
| `static/` | 静态资源文件 |
| `data/` | 运行时数据目录 |
| `docs/` | 项目文档 |

---

## 快速开始

### 前置要求

- **Go**：1.21 或更高版本
- **SQLite3**：数据库支持
- **pnpm**：前端构建工具（可选）
- **Docker**：容器化部署（可选）

### 本地开发

#### 1. 克隆项目

```bash
git clone https://github.com/username/keepblog.git
cd keepblog
```

#### 2. 安装依赖

```bash
go mod download
```

#### 3. 配置环境

```bash
# 复制配置文件
cp config-example.yaml config.yaml

# 编辑配置
vim config.yaml
```

#### 4. 运行项目

```bash
# 使用 Make 运行
make run

# 或直接运行
go run main.go
```

#### 5. 访问应用

- **前台**：http://localhost:8589
- **后台**：http://localhost:8589/console
- **健康检查**：http://localhost:8589/health

### 构建部署

#### Linux 构建

```bash
make linux
```

#### Docker 构建

```bash
# 构建镜像
make docker

# 或使用脚本
./build.sh
```

#### 运行容器

```bash
docker run -d \
  -p 8589:8589 \
  -v ./data:/app/data \
  jieepre/keepblog:latest
```

---

## 前台功能

### 主要页面

| 页面 | 路径 | 功能 |
|------|------|------|
| **首页** | `/`, `/page/:page` | 文章列表，分页显示 |
| **文章详情** | `/post/:hashids` | 显示文章内容和评论 |
| **分类** | `/categories`, `/categories/:category` | 按分类浏览文章 |
| **标签** | `/tags`, `/tags/:tag` | 按标签浏览文章 |
| **归档** | `/archives`, `/archives/:year/:month` | 按时间归档 |
| **搜索** | `/search/:keyword` | 全文搜索 |
| **关于** | `/about` | 关于页面 |
| **友链** | `/link` | 友情链接 |
| **每日一句** | `/daily` | 每日一句话 |

### 功能特性

#### 文章展示

- **Markdown 渲染**：支持完整的 Markdown 语法
- **代码高亮**：集成 Chroma 语法高亮
- **图片懒加载**：优化页面加载性能
- **阅读进度**：文章阅读进度显示

#### 全文搜索

- **搜索范围**：标题、内容、摘要
- **关键词高亮**：搜索结果高亮显示
- **相关性排序**：智能排序搜索结果

#### 评论系统

- **Gitalk 集成**：使用 Gitalk 提供评论功能
- **审核机制**：支持评论审核
- **回复通知**：邮件通知回复

#### 侧边栏

- **最新文章**：显示最新发布的文章
- **分类统计**：各分类文章数量
- **标签云**：热门标签展示
- **归档统计**：按年月统计文章
- **网站资讯**：运行天数、文章总数等

---

## 后台管理

### 登录认证

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/login` | POST | 用户登录 |
| `/api/logout` | POST | 用户登出 |
| `/api/change-password` | POST | 修改密码 |
| `/api/refreshToken` | POST | 刷新令牌 |
| `/api/getInfo` | GET | 获取用户信息 |
| `/api/updateProfile` | PUT | 更新用户资料 |

### 文章管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/post/save` | POST | 保存文章 |
| `/api/v1/site/post/update` | PUT | 更新文章 |
| `/api/v1/site/post/update-publish` | PUT | 发布/取消发布 |
| `/api/v1/site/post/update-top` | PUT | 置顶/取消置顶 |
| `/api/v1/site/post/update-cover/:postId` | PUT | 更新封面图 |
| `/api/v1/site/post/delete/:postId` | DELETE | 删除文章 |
| `/api/v1/site/post/detail/:postId` | GET | 文章详情 |
| `/api/v1/site/post/list` | GET | 文章列表 |

### 分类管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/category/save` | POST | 新增分类 |
| `/api/v1/site/category/update` | PUT | 更新分类 |
| `/api/v1/site/category/update-state` | PUT | 启用/禁用 |
| `/api/v1/site/category/delete/:categoryId` | DELETE | 删除分类 |
| `/api/v1/site/category/detail/:categoryId` | GET | 分类详情 |
| `/api/v1/site/category/list` | GET | 分类列表 |

### 标签管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/tags/list` | GET | 标签列表 |

### 友情链接

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/link/save` | POST | 新增链接 |
| `/api/v1/site/link/update` | PUT | 更新链接 |
| `/api/v1/site/link/update-state` | PUT | 启用/禁用 |
| `/api/v1/site/link/delete/:id` | DELETE | 删除链接 |
| `/api/v1/site/link/detail/:id` | GET | 链接详情 |
| `/api/v1/site/link/list` | GET | 链接列表 |

### 音乐管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/music/save` | POST | 新增音乐 |
| `/api/v1/site/music/update` | PUT | 更新音乐 |
| `/api/v1/site/music/update-state` | PUT | 启用/禁用 |
| `/api/v1/site/music/delete/:id` | DELETE | 删除音乐 |
| `/api/v1/site/music/detail/:id` | GET | 音乐详情 |
| `/api/v1/site/music/list` | GET | 音乐列表 |

### 网站配置

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/web/update` | PUT | 更新配置 |
| `/api/v1/site/web/detail` | GET | 获取配置 |

### 关于页面

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/about/update` | PUT | 更新关于页面 |
| `/api/v1/site/about/detail` | GET | 获取关于内容 |

### 仪表盘

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/dashboard/data` | GET | 仪表盘数据 |

### 通知管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/notice/list` | GET | 通知列表 |
| `/api/v1/notice/read/:id` | PUT | 标记已读 |
| `/api/v1/notice/read-all` | PUT | 全部已读 |
| `/api/v1/notice/delete/:id` | DELETE | 删除通知 |
| `/api/v1/notice/clear` | DELETE | 清空通知 |
| `/api/v1/notice/send` | POST | 发送通知 |

### 日志管理

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/log/level` | GET | 获取日志级别 |
| `/api/log/level` | PUT | 设置日志级别 |
| `/api/log/stats` | GET | 日志统计 |
| `/api/log/list` | GET | 日志列表 |
| `/api/log/read` | GET | 读取日志 |

### 系统监控

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/monitor/server` | GET | 服务器信息 |
| `/api/monitor/realtime` | GET | 实时监控 |
| `/api/monitor/general` | GET | 通用信息 |
| `/api/monitor/loadavg` | GET | 负载平均 |
| `/api/monitor/ram` | GET | 内存信息 |
| `/api/monitor/cpu` | GET | CPU 信息 |
| `/api/monitor/net` | GET | 网络信息 |
| `/api/monitor/diskUsage` | GET | 磁盘使用 |
| `/api/monitor/diskIOStat` | GET | 磁盘 I/O |
| `/api/monitor/base/os` | GET | 操作系统信息 |

### 日志查询

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/site/logs/logon` | GET | 登录日志 |
| `/api/v1/site/logs/access` | GET | 访问日志 |

### 公共接口

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/upload/images` | POST | 上传图片 |

### WebSocket

| 端点 | 方法 | 说明 |
|------|------|------|
| `/ws` | WebSocket | WebSocket 连接 |

---

## API 文档

### 认证方式

后台管理 API 使用 JWT (JSON Web Token) 进行身份认证。

#### 请求头

```http
Authorization: Bearer <token>
```

#### 令牌刷新

访问令牌有效期为 24 小时，过期后需要使用刷新令牌获取新的访问令牌。

### 响应格式

#### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "error message",
  "data": null
}
```

### 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

---

## 数据模型

### 核心模型

#### Post (文章)

```go
type Post struct {
    PostId           uint64  // 文章ID（主键）
    Title            string  // 标题
    PostSlug         string  // URL别名
    Author           string  // 作者
    CoverImage       string  // 封面图片
    PostContent      string  // Markdown内容
    PostContentHtml  string  // HTML内容
    Summary          string  // 摘要
    Type             uint8   // 类型
    CommentEnabled   uint8   // 是否启用评论
    Top              uint8   // 是否置顶
    ReadCount        uint32  // 阅读量
    WordCount        uint32  // 字数
    IsPublished      uint8   // 是否发布
    IsDeleted        uint8   // 是否删除
    Status           uint8   // 状态
    CreateTime       uint64  // 创建时间
    PubTime          uint64  // 发布时间
    LastModifiedTime uint64  // 最后修改时间
    CategoryId       uint32  // 分类ID
    CategoryName     string  // 分类名称（只读）
    TagName          string  // 标签（只读）
    Tags             []string // 标签列表
}
```

**相关表**：`post`

#### Category (分类)

```go
type Category struct {
    CategoryId   uint32  // 分类ID（主键）
    CreateTime   uint64  // 创建时间
    CategoryName string  // 分类名称
    Note         string  // 备注
    State        uint8   // 状态（1-启用，0-禁用）
}
```

**相关表**：`category`

#### Tag (标签)

```go
type Tag struct {
    TagId        uint32  // 标签ID（主键）
    CreateTime   uint64  // 创建时间
    TagName      string  // 标签名称
    TagStyle     string  // 标签样式
    CardTagStyle string  // 卡片标签样式
}
```

**相关表**：`tag`

#### Comment (评论)

```go
type Comment struct {
    CommentId      uint64  // 评论ID（主键）
    Username       string  // 用户名
    Email          string  // 邮箱
    IPAddress      string  // IP地址
    CreateTime     uint64  // 创建时间
    CommentContent string  // 评论内容
    PostId         uint64  // 文章ID
    IsApproved     uint8   // 是否审核通过
}
```

**相关表**：`comment`

#### FriendLink (友情链接)

```go
type FriendLink struct {
    Id               uint32  // 链接ID（主键）
    Title            string  // 标题
    LinkUrl          string  // 链接地址
    LinkIcon         string  // 链接图标
    State            uint8   // 状态
    Type             uint8   // 类型
    LinkDesc         string  // 描述
    CreateTime       uint64  // 创建时间
    LastModifiedTime uint64  // 最后修改时间
}
```

**相关表**：`friend_link`

#### User (用户)

```go
type User struct {
    UserId     int64   // 用户ID（主键）
    Username   string  // 用户名
    Password   string  // 密码（加密）
    NickName   string  // 昵称
    Email      string  // 邮箱
    Phonenumber string  // 手机号
    Sex        uint8   // 性别
    Avatar     string  // 头像
}
```

#### Music (音乐)

```go
type Music struct {
    Id          uint32  // 音乐ID（主键）
    Title       string  // 标题
    Artist      string  // 艺术家
    Url         string  // 音频地址
    Cover       string  // 封面图
    State       uint8   // 状态
    CreateTime  uint64  // 创建时间
}
```

**相关表**：`music`

### 数据库表

| 表名 | 说明 |
|------|------|
| `post` | 文章表 |
| `category` | 分类表 |
| `tag` | 标签表 |
| `post_tag` | 文章标签关联表 |
| `comment` | 评论表 |
| `comment_reply` | 评论回复表 |
| `friend_link` | 友情链接表 |
| `user` | 用户表 |
| `music` | 音乐表 |
| `about` | 关于表 |

---

## 中间件

### 中间件列表

#### 1. GinLogger

**功能**：请求日志记录

**位置**：`internal/middleware/log.go`

**特点**：
- 记录请求方法、路径、状态码、耗时
- 结构化日志输出
- 支持 JSON 格式

#### 2. GinRecovery

**功能**：异常恢复

**位置**：`internal/middleware/log.go`

**特点**：
- 捕获 panic 恢复
- 记录错误堆栈
- 返回 500 错误

#### 3. Cors

**功能**：跨域处理

**位置**：`internal/middleware/cors.go`

**特点**：
- 支持自定义允许的域名
- 支持预检请求
- 配置允许的 HTTP 方法

#### 4. CacheMiddleware

**功能**：缓存中间件

**位置**：`internal/middleware/cache.go`

**特点**：
- 基于路径的缓存
- 可配置缓存时间
- 支持动态刷新

#### 5. JwtVerify

**功能**：JWT 认证

**位置**：`internal/middleware/jwt.go`

**特点**：
- 验证 JWT 令牌
- 提取用户信息
- 支持令牌刷新

#### 6. RateLimit

**功能**：限流保护

**位置**：`internal/middleware/ratelimit.go`

**特点**：
- 基于 IP 的限流
- 基于用户名的限流
- 可配置请求频率

#### 7. ErrorHandler

**功能**：统一错误处理

**位置**：`internal/middleware/error.go`

**特点**：
- 捕获所有错误
- 统一错误格式
- 记录错误日志

#### 8. Csrf

**功能**：CSRF 防护

**位置**：`internal/middleware/csrf.go`

**特点**：
- 生成 CSRF Token
- 验证 Token
- 双重 Cookie 提交

#### 9. Statistics

**功能**：统计中间件

**位置**：`internal/middleware/access.go`

**特点**：
- 记录访问日志
- 统计 PV/UV
- 分析用户行为

#### 10. MemoryLimit

**功能**：内存限制

**位置**：`internal/middleware/memorylimit.go`

**特点**：
- 监控内存使用
- 防止内存泄漏
- 自动清理

### 中间件配置

#### 前台中间件

```go
webGroup := engine.Group("/")
{
    webGroup.Use(middleware.IPBasedRateLimit(60, time.Minute))
    webGroup.Use(middleware.Statistics())
    webGroup.Use(middleware.CacheMiddleware(5 * time.Minute))
}
```

#### 后台中间件

```go
adminGroup := engine.Group("/api")
{
    adminGroup.Use(middleware.ErrorHandler())
    adminGroup.Use(monitor.MetricsMiddleware(metrics))
    adminGroup.Use(middleware.APIRateLimit(100, time.Minute))
    adminGroup.Use(middleware.LoginRateLimit())
    adminGroup.Use(middleware.JwtVerify())
}
```

---

## 配置说明

### 配置文件

项目使用 `config.yaml` 进行配置，配置示例文件为 `config-example.yaml`。

### 配置项

#### HTTP 配置

```yaml
http:
  port: 8589  # HTTP 服务端口
```

#### Redis 配置

```yaml
redis:
  host: localhost     # Redis 主机
  port: 6379         # Redis 端口
  password: ""       # 密码（留空表示无密码）
  database: 0        # 数据库编号
  enable: false      # 是否启用 Redis 缓存
```

#### Minio 配置

```yaml
minio:
  serverUrl: http://192.168.2.166:9000  # Minio 服务地址
  endpoint: 192.168.2.166:9000          # Minio 端点
  accessKeyID: your_access_key         # 访问密钥 ID
  secretAccessKey: your_secret_key      # 密钥
  useSSL: false                         # 是否使用 SSL
  bucketName: test                      # 存储桶名称
```

#### 百度推送配置

```yaml
baidu:
  push: true                              # 是否启用百度推送
  url: https://www.keepon.online         # 网站地址
  token: VzXKVoNrqnSAWINA                # 百度推送 Token
```

#### Gitalk 评论配置

```yaml
gitalk:
  enable: true                          # 是否启用 Gitalk
  clientID: '4b719059ee103950a239'    # GitHub 应用 Client ID
  clientSecret: 'b9d642cffab418205da3c8c1a0bf762a9d10817b'  # Client Secret
  repo: 'gitalk'                        # GitHub 仓库名
  owner: 'keepon-online'               # 仓库所有者
  admin:                               # 管理员列表
    - 'keepon-online'
```

### 环境变量

支持通过环境变量覆盖配置文件：

| 环境变量 | 说明 |
|----------|------|
| `PORT` | HTTP 服务端口 |
| `REDIS_HOST` | Redis 主机 |
| `REDIS_PORT` | Redis 端口 |
| `REDIS_PASSWORD` | Redis 密码 |
| `MINIO_ENDPOINT` | Minio 端点 |

---

## 部署指南

### Docker 部署

#### 构建镜像

```bash
# 使用 Makefile
make docker

# 或使用构建脚本
./build.sh
```

#### 运行容器

```bash
docker run -d \
  --name keepblog \
  -p 8589:8589 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/config.yaml:/app/config.yaml \
  jieepre/keepblog:latest
```

#### Docker Compose

```yaml
version: '3.8'
services:
  keepblog:
    image: jieepre/keepblog:latest
    ports:
      - "8589:8589"
    volumes:
      - ./data:/app/data
      - ./config.yaml:/app/config.yaml
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped

  minio:
    image: minio/minio:latest
    command: server /data
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - ./minio-data:/data
    restart: unless-stopped
```

### 传统部署

#### Linux

```bash
# 下载构建产物
wget https://github.com/username/keepblog/releases/latest/keepblog-linux-amd64

# 赋予执行权限
chmod +x keepblog-linux-amd64

# 运行
./keepblog-linux-amd64
```

#### Windows

```powershell
# 下载构建产物
# keepblog-windows-amd64.exe

# 运行
.\keepblog-windows-amd64.exe
```

#### macOS

```bash
# 下载构建产物
wget https://github.com/username/keepblog/releases/latest/keepblog-darwin-amd64

# 赋予执行权限
chmod +x keepblog-darwin-amd64

# 运行
./keepblog-darwin-amd64
```

### Nginx 反向代理

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8589;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### Systemd 服务

```ini
[Unit]
Description=KeepBlog Blog System
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/keepblog
ExecStart=/opt/keepblog/keepblog
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl start keepblog
sudo systemctl enable keepblog
```

---

## 开发规范

### 代码风格

#### 命名规范

- **包名**：小写，简短，有意义的单词
- **常量**：大驼峰（PascalCase）
- **变量**：小驼峰（camelCase）
- **函数**：小驼峰（camelCase），导出函数首字母大写
- **结构体**：大驼峰（PascalCase）

#### 注释规范

```go
// Package model 定义数据模型
package model

// Post 文章模型
type Post struct {
    PostId uint64 // 文章ID
}

// GetPostById 根据 ID 获取文章
func (s *Service) GetPostById(id uint64) (*Post, error) {
    // 实现
}
```

### Git 提交规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type)**：
- `feat`：新功能
- `fix`：bug 修复
- `docs`：文档更新
- `style`：代码格式（不影响功能）
- `refactor`：重构（不是新功能也不是修复）
- `test`：添加测试
- `chore`：构建过程或辅助工具的变动

**示例**：

```
feat(post): 添加文章置顶功能

- 在文章模型中添加 Top 字段
- 实现置顶/取消置顶接口
- 前台列表按置顶排序

Closes #123
```

### 错误处理

```go
// 使用 fmt.Errorf 包装错误
if err != nil {
    return fmt.Errorf("failed to save post: %w", err)
}

// 使用自定义错误类型
type NotFoundError struct {
    Resource string
    ID       uint64
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s not found: %d", e.Resource, e.ID)
}
```

### 日志规范

```go
// 使用结构化日志
slog.Infow("post saved",
    "postId", post.PostId,
    "title", post.Title,
    "author", post.Author,
)

// 使用上下文
slog.WithField("postId", post.PostId).Info("post saved")

// 错误日志
slog.Errorf("failed to save post: %v", err)
```

### 测试规范

```go
func TestPostService_GetPostById(t *testing.T) {
    tests := []struct {
        name    string
        postId  uint64
        want    *Post
        wantErr bool
    }{
        {
            name:   "success",
            postId: 1,
            want:   &Post{PostId: 1},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.GetPostById(tt.postId)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetPostById() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetPostById() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 附录

### Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 本地构建 |
| `make linux` | Linux AMD64 构建 |
| `make linux-arm64` | Linux ARM64 构建 |
| `make darwin` | macOS AMD64 构建 |
| `make windows` | Windows AMD64 构建 |
| `make docker` | 构建 Docker 镜像 |
| `make console` | 构建前端并同步 |
| `make version` | 显示版本信息 |
| `make run` | 运行应用 |
| `make help` | 查看所有命令 |

### 常见问题

#### Q1: 如何重置管理员密码？

A: 直接修改数据库中的用户记录，或通过登录接口的修改密码功能。

#### Q2: 如何启用 Redis 缓存？

A: 在 `config.yaml` 中设置 `redis.enable: true`，并配置正确的连接信息。

#### Q3: 如何更换数据库？

A: 修改 `internal/core/db.go` 中的数据库驱动，并更新 GORM 配置。

#### Q4: 如何自定义后台管理界面？

A: 修改 `static/console/` 目录下的静态文件。

#### Q5: 如何添加新的中间件？

A: 在 `internal/middleware/` 目录下创建新文件，并在 `routers.go` 中注册。

---

## 前端项目

### 项目概览

KeepBlog 前端项目采用 Vue 3 + Vite + Element Plus 技术栈，基于 [vue-pure-admin](https://github.com/pure-admin/vue-pure-admin) 精简版构建，提供完整的管理后台界面。

### 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **框架** | Vue | 3.5.22 | 核心框架 |
| **构建工具** | Vite | 7.1.12 | 开发构建 |
| **UI 库** | Element Plus | 2.11.5 | UI 组件 |
| **语言** | TypeScript | 5.9.3 | 类型安全 |
| **状态管理** | Pinia | 3.0.3 | 状态管理 |
| **路由** | Vue Router | 4.6.3 | 路由管理 |
| **HTTP 客户端** | Axios | 1.12.2 | API 请求 |
| **图表** | ECharts | 6.0.0 | 数据可视化 |
| **Markdown** | md-editor-v3 | 5.8.5 | Markdown 编辑 |

### 项目结构

```
console/                    # 前端项目根目录
├── build/                  # 构建配置
│   ├── plugins.ts         # Vite 插件配置
│   └── utils.ts           # 构建工具函数
├── dist/                   # 构建产物（自动生成）
│   ├── index.html         # 入口文件
│   ├── static/            # 静态资源
│   │   ├── js/           # JavaScript 文件
│   │   └── css/          # 样式文件
│   └── platform-config.json # 平台配置
├── mock/                   # Mock 数据
├── public/                 # 公共资源
├── src/                    # 源代码
│   ├── api/               # API 接口
│   │   ├── about.ts      # 关于页面 API
│   │   ├── category.ts   # 分类管理 API
│   │   ├── common.ts     # 公共 API
│   │   ├── dashboard.ts  # 仪表盘 API
│   │   ├── link.ts       # 友情链接 API
│   │   ├── music.ts      # 音乐管理 API
│   │   ├── notice.ts     # 通知管理 API
│   │   ├── post.ts       # 文章管理 API
│   │   ├── profile.ts    # 个人资料 API
│   │   ├── routes.ts     # 路由 API
│   │   ├── setting.ts    # 设置 API
│   │   ├── system.ts     # 系统管理 API
│   │   ├── tag.ts        # 标签管理 API
│   │   ├── user.ts       # 用户管理 API
│   │   └── monitor.ts    # 系统监控 API
│   ├── assets/            # 静态资源
│   │   └── iconfont/     # 图标字体
│   ├── components/        # 公共组件
│   │   ├── ReAuth/       # 权限组件
│   │   ├── ReCol/        # 列格组件
│   │   ├── ReDialog/     # 对话框组件
│   │   ├── ReIcon/       # 图标组件
│   │   └── ...
│   ├── config/            # 配置文件
│   │   ├── index.ts      # 平台配置
│   │   └── theme.ts      # 主题配置
│   ├── directives/        # 自定义指令
│   ├── layout/            # 布局组件
│   │   ├── components/   # 布局子组件
│   │   └── index.ts      # 布局入口
│   ├── plugins/           # 插件
│   │   ├── echarts.ts    # ECharts 配置
│   │   └── elementPlus.ts # Element Plus 配置
│   ├── router/            # 路由配置
│   │   ├── index.ts      # 路由入口
│   │   ├── modules/      # 路由模块
│   │   │   ├── login.ts  # 登录路由
│   │   │   ├── manage/   # 管理路由
│   │   │   └── system/   # 系统路由
│   │   └── utils.ts      # 路由工具
│   ├── store/             # 状态管理
│   │   ├── modules/      # Store 模块
│   │   └── types.ts      # 类型定义
│   ├── style/             # 样式文件
│   │   ├── index.scss    # 主样式
│   │   ├── reset.scss    # 重置样式
│   │   └── tailwind.css # Tailwind CSS
│   ├── utils/             # 工具函数
│   │   ├── auth.ts       # 认证工具
│   │   ├── tree.ts       # 树形数据处理
│   │   └── ...
│   ├── views/             # 页面组件
│   │   ├── base/         # 基础页面
│   │   ├── error/        # 错误页面
│   │   ├── login/        # 登录页面
│   │   ├── manage/       # 管理页面
│   │   │   ├── post/     # 文章管理
│   │   │   ├── category/ # 分类管理
│   │   │   ├── tag/      # 标签管理
│   │   │   ├── link/     # 友情链接
│   │   │   ├── music/    # 音乐管理
│   │   │   └── system/   # 系统管理
│   │   └── welcome/      # 欢迎页
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── types/                  # TypeScript 类型定义
├── .env                    # 环境变量（默认）
├── .env.development        # 开发环境配置
├── .env.production         # 生产环境配置
├── .env.staging            # 预发布环境配置
├── .eslintrc.js           # ESLint 配置
├── .prettierrc.js          # Prettier 配置
├── .stylelintrc.js         # Stylelint 配置
├── index.html              # HTML 模板
├── package.json            # 项目依赖
├── pnpm-lock.yaml          # 依赖锁定
├── tsconfig.json           # TypeScript 配置
├── vite.config.ts          # Vite 配置
└── tailwind.config.ts      # Tailwind CSS 配置
```

### 配置文件

#### 环境变量配置

**开发环境** (`.env.development`)：

```env
# 平台本地运行端口号
VITE_PORT = 8848

# 开发环境读取配置文件路径
VITE_PUBLIC_PATH = /

# 开发环境路由历史模式（hash模式）
VITE_ROUTER_HISTORY = "hash"

# 后端 API 地址
VITE_BASE_URL = "http://192.168.3.6:8589"
```

**生产环境** (`.env.production`)：

```env
# 线上环境平台打包路径
VITE_PUBLIC_PATH = /console/

# 线上环境路由历史模式（Hash模式）
VITE_ROUTER_HISTORY = "hash"

# 是否在打包时使用cdn替换本地库
VITE_CDN = false

# 是否启用gzip压缩或brotli压缩
VITE_COMPRESSION = "none"

# 后端 API 地址
VITE_BASE_URL = "https://www.keepon.online"
```

#### Vite 配置 (vite.config.ts)

```typescript
export default ({ mode }: ConfigEnv): UserConfigExport => {
  const { VITE_CDN, VITE_PORT, VITE_COMPRESSION, VITE_PUBLIC_PATH, VITE_BASE_URL } =
    wrapperEnv(loadEnv(mode, root));

  return {
    base: VITE_PUBLIC_PATH,  // 公共路径
    server: {
      port: VITE_PORT,       // 开发服务器端口
      host: "0.0.0.0",
      proxy: {
        "^/api": {
          target: VITE_BASE_URL + "/api",  // 后端服务地址
          changeOrigin: true,
          rewrite: path => path.replace(/^\/api/, "")
        }
      }
    },
    build: {
      target: "es2015",
      sourcemap: false,
      chunkSizeWarningLimit: 4000,
      rollupOptions: {
        output: {
          chunkFileNames: "static/js/[name]-[hash].js",
          entryFileNames: "static/js/[name]-[hash].js",
          assetFileNames: "static/[ext]/[name]-[hash].[ext]"
        }
      }
    }
  };
};
```

### 核心功能

#### 1. 路由管理

**路由模式**：Hash 模式（避免刷新问题）

**路由结构**：

```typescript
// 静态路由（自动导入 src/router/modules/**/*.ts）
const constantRoutes: Array<RouteRecordRaw>

// 动态路由（从后端获取）
const asyncRoutes: Array<RouteRecordRaw>
```

**路由配置示例**：

```typescript
// src/router/modules/manage/post.ts
export default {
  path: "/manage/post",
  name: "PostManage",
  component: () => import("@/layout/index.vue"),
  redirect: "/manage/post/list",
  meta: {
    title: "文章管理",
    icon: "article",
    showLink: true,
    rank: 1
  },
  children: [
    {
      path: "/manage/post/list",
      name: "PostList",
      component: () => import("@/views/manage/post/index.vue"),
      meta: {
        title: "文章列表",
        showLink: true
      }
    }
  ]
};
```

#### 2. API 封装

**请求拦截器**：

```typescript
// src/utils/http/index.ts
import axios from "axios";

const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
  timeout: 10000
});

// 请求拦截
service.interceptors.request.use(
  config => {
    const token = useUserStoreHook().getToken;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截
service.interceptors.response.use(
  response => {
    const { data } = response;
    if (data.code === 200) {
      return data.data;
    } else {
      ElMessage.error(data.message);
      return Promise.reject(new Error(data.message));
    }
  },
  error => {
    ElMessage.error(error.message);
    return Promise.reject(error);
  }
);
```

**API 调用示例**：

```typescript
// src/api/post.ts
import { http } from "@/utils/http";

export const getPostList = (params: any) => {
  return http.request<PostListResponse>("get", "/api/v1/site/post/list", {
    params
  });
};

export const savePost = (data: PostData) => {
  return http.request("post", "/api/v1/site/post/save", {
    data
  });
};
```

#### 3. 状态管理

**Store 模块**：

```typescript
// src/store/modules/user.ts
import { defineStore } from "pinia";

export const useUserStoreHook = defineStore("user", {
  state: (): UserState => ({
    username: "",
    roles: [],
    token: ""
  }),
  actions: {
    SET_USERNAME(username: string) {
      this.username = username;
    },
    SET_ROLES(roles: string[]) {
      this.roles = roles;
    },
    SET_TOKEN(token: string) {
      this.token = token;
    }
  }
});
```

**使用示例**：

```typescript
import { useUserStoreHook } from "@/store/modules/user";

const userStore = useUserStoreHook();

// 设置用户信息
userStore.SET_USERNAME("admin");

// 获取 Token
const token = userStore.token;
```

#### 4. 权限控制

**路由权限**：

```typescript
// src/router/index.ts
router.beforeEach((to, _from, next) => {
  const userInfo = storageLocal().getItem<DataInfo>(userKey);

  // 检查权限
  if (to.meta?.roles && !isOneOfArray(to.meta?.roles, userInfo?.roles)) {
    next({ path: "/error/403" });
  } else {
    next();
  }
});
```

**按钮权限**：

```vue
<template>
  <!-- 只有 admin 角色可以看到删除按钮 -->
  <Auth value="admin">
    <el-button @click="handleDelete">删除</el-button>
  </Auth>
</template>
```

### 开发指南

#### 环境准备

**前置要求**：
- Node.js：^20.19.0 或 >=22.13.0
- pnpm：>=9

**安装依赖**：

```bash
cd console
pnpm install
```

#### 本地开发

```bash
# 启动开发服务器
pnpm dev

# 或
pnpm serve
```

访问：http://localhost:8848

#### 代码规范

**Lint 检查**：

```bash
# ESLint 检查
pnpm lint:eslint

# Prettier 格式化
pnpm lint:prettier

# Stylelint 检查
pnpm lint:stylelint

# 全部检查
pnpm lint
```

**类型检查**：

```bash
pnpm typecheck
```

#### 构建部署

**生产构建**：

```bash
# 构建生产版本
pnpm build

# 构建预发布版本
pnpm build:staging
```

**预览构建结果**：

```bash
pnpm preview
```

### 与后端集成

#### 开发环境

1. **启动后端服务**：

```bash
cd keepblog
make run
```

后端运行在：http://localhost:8589

2. **启动前端开发服务器**：

```bash
cd console
pnpm dev
```

前端运行在：http://localhost:8848

3. **API 代理**：

前端通过 Vite 代理将 `/api` 请求转发到后端：

```typescript
proxy: {
  "^/api": {
    target: "http://localhost:8589/api",
    changeOrigin: true,
    rewrite: path => path.replace(/^\/api/, "")
  }
}
```

#### 生产环境

1. **构建前端**：

```bash
cd console
pnpm build
```

构建产物在 `console/dist/` 目录。

2. **同步到后端**：

```bash
cd ../keepblog
./sync-console.sh
```

或使用 Make：

```bash
make console
```

3. **构建后端**：

```bash
make build
```

4. **运行**：

```bash
./bin/keepblog
```

访问后台管理：http://localhost:8589/console

### 常见问题

#### Q1: 如何修改 API 地址？

A: 编辑对应环境的 `.env` 文件，修改 `VITE_BASE_URL`：

```env
VITE_BASE_URL = "http://your-backend-url"
```

#### Q2: 如何添加新页面？

A: 在 `src/views/` 下创建页面组件，然后在 `src/router/modules/` 添加路由配置。

#### Q3: 如何添加 API 接口？

A: 在 `src/api/` 下创建接口文件，封装 API 调用函数。

#### Q4: 如何使用组件？

A: 全局组件已注册，可以直接使用：

```vue
<template>
  <el-button type="primary">按钮</el-button>
  <IconifyIconOffline icon="ep:edit" />
</template>
```

#### Q5: 如何调试构建问题？

A: 使用构建分析工具：

```bash
pnpm report
```

### 相关资源

- [Vue 3 文档](https://cn.vuejs.org/)
- [Vite 文档](https://cn.vitejs.dev/)
- [Element Plus 文档](https://element-plus.org/zh-CN/)
- [Pinia 文档](https://pinia.vuejs.org/zh/)
- [vue-pure-admin 文档](https://pure-admin.cn/)

---

### 相关资源

- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Go 官方文档](https://go.dev/doc/)
- [JWT 规范](https://jwt.io/)
- [Vue 3 文档](https://cn.vuejs.org/)
- [Vite 文档](https://cn.vitejs.dev/)
- [Element Plus 文档](https://element-plus.org/zh-CN/)

### 许可证

MIT License

### 贡献指南

欢迎提交 Issue 和 Pull Request！

---

**文档版本**：1.1.0
**最后更新**：2026-01-26
**维护者**：KeepBlog Team
