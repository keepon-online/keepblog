# 项目结构

## 目录树

```
keepblog/
├── api/                          # API 处理器层
│   ├── admin/                    # 后台管理 API (14 个模块)
│   │   ├── about/                # 关于页面管理
│   │   │   └── about.go
│   │   ├── category/             # 分类管理
│   │   │   └── category.go
│   │   ├── common/               # 公共功能（上传等）
│   │   │   └── upload.go
│   │   ├── dashboard/            # 仪表板
│   │   │   └── dashboard.go
│   │   ├── link/                 # 友情链接管理
│   │   │   └── link.go
│   │   ├── log/                  # 日志查看器
│   │   │   └── log.go
│   │   ├── login/                # 登录认证
│   │   │   └── login.go
│   │   ├── monitor/              # 系统监控
│   │   │   ├── types.go          # 监控 DTO
│   │   │   ├── handlers.go       # HTTP 处理器
│   │   │   └── collectors.go     # 系统信息采集
│   │   ├── music/                # 音乐管理
│   │   │   └── music.go
│   │   ├── notice/               # 通知管理
│   │   │   └── notice.go
│   │   ├── post/                 # 文章管理
│   │   │   └── post.go
│   │   ├── system/               # 系统设置
│   │   │   ├── logs.go
│   │   │   └── user.go
│   │   ├── tags/                 # 标签管理
│   │   │   └── tags.go
│   │   └── website/              # 网站配置
│   │       └── website.go
│   └── web/                      # 前台 API (8 个模块)
│       ├── about/                # 关于页面
│       │   └── about.go
│       ├── archive/              # 归档
│       │   └── archive.go
│       ├── category/             # 分类浏览
│       │   └── category.go
│       ├── home/                 # 首页
│       │   └── home.go
│       ├── link/                 # 友情链接
│       │   └── link.go
│       ├── music/                # 音乐播放
│       │   └── music.go
│       ├── post/                 # 文章详情
│       │   └── post.go
│       └── tags/                 # 标签浏览
│           └── tags.go
├── config/                       # 配置管理
│   ├── config.go                 # 配置加载和验证
│   └── config_test.go            # 配置测试
├── internal/                     # 内部包
│   ├── app/                      # 应用入口
│   │   ├── app.go                # 应用程序主体
│   │   └── routers.go            # 路由创建
│   ├── cache/                    # Redis 缓存
│   │   └── redis.go
│   ├── config/                   # 服务器配置
│   │   └── server.go
│   ├── core/                     # 核心初始化
│   │   ├── db.go                 # 数据库初始化
│   │   ├── log.go                # 日志初始化
│   │   ├── migrate.go            # 种子数据初始化
│   │   ├── seed.go               # 嵌入式种子数据
│   │   ├── timer.go              # 定时任务
│   │   └── migrations/           # Goose 版本化迁移
│   ├── errors/                   # 错误处理
│   │   └── errors.go
│   ├── logger/                   # 日志系统
│   │   └── logger.go
│   ├── middleware/               # 中间件 (10 个)
│   │   ├── access.go             # 访问日志
│   │   ├── cache.go              # 缓存中间件
│   │   ├── cors.go               # CORS 跨域
│   │   ├── csrf.go               # CSRF 防护
│   │   ├── error.go              # 错误处理
│   │   ├── jwt.go                # JWT 验证
│   │   ├── log.go                # 日志记录
│   │   ├── memorylimit.go        # 内存限制
│   │   ├── print.go              # 打印中间件
│   │   └── ratelimit.go          # 速率限制
│   ├── model/                    # 数据模型
│   │   ├── post.go               # 文章、分类、标签、评论等
│   │   ├── music.go              # 音乐模型
│   │   ├── dashboard.go          # 仪表板数据
│   │   ├── request/              # 请求模型
│   │   │   ├── category.go
│   │   │   ├── link.go
│   │   │   ├── music.go
│   │   │   ├── post.go
│   │   │   └── user.go
│   │   ├── response/             # 响应模型
│   │   │   ├── category.go
│   │   │   ├── link.go
│   │   │   ├── music.go
│   │   │   ├── post.go
│   │   │   └── user.go
│   │   ├── server/               # 服务器信息模型
│   │   │   ├── cpu.go
│   │   │   ├── disk.go
│   │   │   ├── mem.go
│   │   │   ├── net.go
│   │   │   └── server.go
│   │   └── system/               # 系统模型
│   │       ├── system.go         # 访问日志、登录日志、网站配置
│   │       └── notice.go         # 通知模型
│   ├── monitor/                  # 健康检查和监控
│   │   ├── health.go
│   │   └── metrics.go
│   ├── pkg/                      # 内部工具包
│   │   ├── core/                 # 核心上下文
│   │   │   ├── context.go        # 请求上下文
│   │   │   └── parse.go
│   │   ├── oss/                  # 对象存储
│   │   │   └── minio.go
│   │   └── seo.go                # SEO 优化
│   ├── router/                   # 路由定义
│   │   ├── route.go              # 路由基础定义
│   │   ├── admin.go              # 后台路由注册
│   │   └── web.go                # 前台路由注册
│   ├── service/                  # 业务逻辑层 (12 个服务)
│   │   ├── about/                # 关于服务
│   │   │   └── about.go
│   │   ├── category/             # 分类服务
│   │   │   └── category.go
│   │   ├── dashboard/            # 仪表板服务
│   │   │   └── dashboard.go
│   │   ├── link/                 # 友链服务
│   │   │   └── link.go
│   │   ├── music/                # 音乐服务
│   │   │   └── music.go
│   │   ├── notice/               # 通知服务
│   │   │   └── notice.go
│   │   ├── post/                 # 文章服务
│   │   │   └── post.go
│   │   ├── sidebar/              # 侧边栏服务
│   │   │   └── sidebar.go
│   │   ├── system/               # 系统服务
│   │   │   └── system.go
│   │   ├── tag/                  # 标签服务
│   │   │   └── tag.go
│   │   ├── website/              # 网站服务
│   │   │   └── website.go
│   │   └── service.go            # 服务聚合与依赖注入
│   ├── websocket/                # WebSocket 支持
│   │   ├── handler.go
│   │   └── hub.go
│   └── web/                      # Web 模板
│       └── template.go
├── pkg/                          # 公共工具包 (20+ 个)
│   ├── area/                     # IP 地理位置
│   │   └── area.go
│   ├── bcryt.go                  # 密码加密
│   ├── bing.go                   # Bing API
│   ├── cloudtag/                 # 标签云
│   │   └── cloudtag.go
│   ├── cmd/                      # 命令行工具
│   │   └── cmd.go
│   ├── copier/                   # 对象复制
│   │   └── copier.go
│   ├── daily/                    # 每日数据
│   │   └── daily.go
│   ├── file/                     # 文件操作
│   │   └── file.go
│   ├── hash/                     # Hash 工具
│   │   └── hash.go
│   ├── image.go                  # 图片处理
│   ├── ipconvert.go              # IP 转换
│   ├── iputil/                   # IP 工具
│   │   └── iputil.go
│   ├── jwttoken/                 # JWT 令牌
│   │   ├── token.go
│   │   ├── blacklist.go          # 令牌黑名单
│   │   └── token_test.go
│   ├── math/                     # 数学工具
│   │   └── math.go
│   ├── md/                       # Markdown 处理
│   │   └── md.go
│   ├── page/                     # 分页工具
│   │   └── page.go
│   ├── psutil/                   # 系统工具
│   │   └── psutil.go
│   ├── result/                   # 响应结果
│   │   └── result.go
│   ├── tags-remove/              # 标签移除
│   │   └── tags.go
│   ├── ua/                       # 用户代理
│   │   └── ua.go
│   └── xpack/                    # 打包工具
│       └── xpack.go
├── static/                       # 静态资源
│   ├── static.go                 # 前台静态资源
│   ├── console/                  # 后台前端
│   │   └── static.go
│   ├── css/                      # 样式文件
│   ├── js/                       # JavaScript 文件
│   ├── images/                   # 图片资源
│   ├── plugins/                  # 插件
│   ├── robots.txt                # 搜索引擎爬虫规则
│   └── favicon.ico               # 网站图标
├── global/                       # 全局变量
│   └── global.go                 # 全局 GORM 实例
├── main.go                       # 应用入口
├── go.mod                        # 模块定义
├── go.sum                        # 依赖锁定
├── Makefile                      # 构建脚本
├── Dockerfile                    # Docker 构建
├── docker-compose.yml            # Docker 编排
├── entrypoint.sh                 # 容器启动脚本
├── config-example.yaml           # 配置示例
├── config.yaml                   # 配置文件
├── build.sh                      # 构建脚本
├── README.md                     # 项目说明
└── CHANGELOG.md                  # 更新日志
```

---

## 核心目录说明

### 1. api/ - API 处理器层

API 层负责处理 HTTP 请求和响应，分为前台和后台两部分。

#### api/admin/ - 后台管理 API
| 模块 | 文件 | 功能 | 主要接口 |
|------|------|------|----------|
| about | about.go | 关于页面管理 | 更新、获取关于页面 |
| category | category.go | 分类管理 | CRUD、状态管理 |
| common | upload.go | 公共功能 | 图片上传 |
| dashboard | dashboard.go | 仪表板 | 数据统计、图表 |
| link | link.go | 友情链接管理 | CRUD、状态管理 |
| log | log.go | 日志查看器 | 日志级别、查看、统计 |
| login | login.go | 登录认证 | 登录、登出、刷新令牌 |
| monitor | monitor.go | 系统监控 | CPU、内存、磁盘、网络 |
| music | music.go | 音乐管理 | CRUD、排序、状态 |
| notice | notice.go | 通知管理 | 列表、已读、删除 |
| post | post.go | 文章管理 | CRUD、发布、置顶 |
| system | logs.go, user.go | 系统设置 | 用户、日志管理 |
| tags | tags.go | 标签管理 | 标签列表 |
| website | website.go | 网站配置 | SEO、ICP、社交链接 |

参考: `api/admin/`

#### api/web/ - 前台展示 API
| 模块 | 文件 | 功能 | 主要接口 |
|------|------|------|----------|
| about | about.go | 关于页面 | 获取关于内容 |
| archive | archive.go | 归档 | 按时间归档文章 |
| category | category.go | 分类浏览 | 分类列表、分类文章 |
| home | home.go | 首页 | 文章列表、搜索、每日推荐 |
| link | link.go | 友情链接 | 友链列表 |
| music | music.go | 音乐播放 | 音乐列表 |
| post | post.go | 文章详情 | 单篇文章展示 |
| tags | tags.go | 标签浏览 | 标签列表、标签页面 |

参考: `api/web/`

---

### 2. internal/ - 内部包

内部包包含应用的核心逻辑，不对外暴露。

#### internal/app/ - 应用入口
- **app.go**: 应用程序主体，初始化和启动
- **routers.go**: 路由创建和中间件注册

参考: `internal/app/app.go:1`, `internal/app/routers.go:1`

#### internal/service/ - 业务逻辑层
| 服务 | 文件 | 职责 |
|------|------|------|
| PostService | post/post.go | 文章 CRUD、发布、置顶、搜索 |
| CategoryService | category/category.go | 分类 CRUD、状态管理 |
| TagService | tag/tag.go | 标签管理、标签云 |
| SidebarService | sidebar/sidebar.go | 侧边栏数据聚合 |
| AboutService | about/about.go | 关于页面管理 |
| SystemService | system/system.go | 系统配置、日志管理 |
| LinkService | link/link.go | 友链 CRUD、状态管理 |
| WebSiteService | website/website.go | 网站配置管理 |
| Dashboard | dashboard/dashboard.go | 仪表板数据统计 |
| MusicService | music/music.go | 音乐 CRUD、排序 |
| NoticeService | notice/notice.go | 通知管理 |

**服务聚合** (`service.go`):
```go
type AppService struct {
    PostService     *post.Service
    CategoryService *category.Service
    TagService      *tag.Service
    SidebarService  *sidebar.Service
    AboutService    *about.Service
    SystemService   *system.Service
    LinkService     *link.Service
    WebSiteService  *website.Service
    Dashboard       *dashboard.Service
    MusicService    *music.Service
    NoticeService   *notice.Service
}
```

参考: `internal/service/service.go:1`

#### internal/model/ - 数据模型
| 模型文件 | 包含的模型 |
|----------|-----------|
| post.go | Post, Category, Tag, PostTag, Comment, CommentReply, FriendLink, User, About |
| music.go | Music |
| dashboard.go | DashboardData, PostStatistics, CategoryStatistics |
| system/system.go | AccessLog, LoginLog, WebSite |
| system/notice.go | Notice |
| server/*.go | ServerInfo, CPUInfo, MemInfo, DiskInfo, NetInfo |

**请求/响应模型**:
- `model/request/`: 请求参数模型
- `model/response/`: 响应数据模型

参考: `internal/model/`

#### internal/middleware/ - 中间件
| 中间件 | 文件 | 功能 |
|--------|------|------|
| GinLogger | log.go | 请求日志记录 |
| GinRecovery | log.go | 恐慌恢复 |
| Cors | cors.go | CORS 跨域 |
| Gzip | - | 响应压缩 (gin-contrib) |
| IPBasedRateLimit | ratelimit.go | IP 限流 (60/分钟) |
| APIRateLimit | ratelimit.go | API 限流 (100/分钟) |
| LoginRateLimit | ratelimit.go | 登录限流 (5/分钟) |
| Statistics | access.go | 访问统计 |
| CacheMiddleware | cache.go | 缓存 (5 分钟) |
| ErrorHandler | error.go | 错误处理 |
| JwtVerify | jwt.go | JWT 验证 |
| CSRFProtection | csrf.go | CSRF 防护 |
| MemoryLimit | memorylimit.go | 内存限制 |

参考: `internal/middleware/`

#### internal/core/ - 核心初始化
- **db.go**: SQLite 数据库连接和连接池配置
- **log.go**: 日志初始化
- **migrate.go**: 种子数据初始化
- **seed.go**: 嵌入式文章、音乐种子数据与随机管理员密码
- **timer.go**: 定时任务
- **migrations/**: Goose 版本化 schema/index 迁移；启动时幂等执行

参考: `internal/app/app.go:Initialize`, `internal/core/migrations/runner.go:Run`

#### internal/router/ - 路由定义
- **route.go**: 路由基础定义和接口
- **admin.go**: 后台路由注册
- **web.go**: 前台路由注册

参考: `internal/router/`

#### internal/monitor/ - 监控
- **health.go**: 健康检查端点
- **metrics.go**: Prometheus 性能指标

参考: `internal/monitor/health.go:1`, `internal/monitor/metrics.go:1`

#### internal/pkg/ - 内部工具包
- **core/context.go**: 请求上下文封装
- **oss/minio.go**: Minio 对象存储客户端
- **seo.go**: SEO 优化工具

参考: `internal/pkg/`

---

### 3. pkg/ - 公共工具包

公共工具包提供可复用的工具函数和辅助类。

| 工具包 | 文件 | 功能 |
|--------|------|------|
| area | area.go | IP 地理位置解析 (ip2region) |
| bcrypt | bcryt.go | 密码加密和验证 |
| bing | bing.go | Bing API 集成 |
| cloudtag | cloudtag.go | 标签云生成 |
| cmd | cmd.go | 命令行工具 |
| copier | copier.go | 对象复制 |
| daily | daily.go | 每日数据推荐 |
| file | file.go | 文件操作 |
| hash | hash.go | Hash 工具 (hashids) |
| image | image.go | 图片处理 |
| iputil | iputil.go | IP 工具函数 |
| jwttoken | token.go, blacklist.go | JWT 令牌生成和黑名单 |
| math | math.go | 数学工具 |
| md | md.go | Markdown 渲染 (goldmark) |
| page | page.go | 分页工具 |
| psutil | psutil.go | 系统信息获取 |
| result | result.go | 统一响应格式 |
| tags-remove | tags.go | HTML 标签移除 |
| ua | ua.go | 用户代理解析 |
| xpack | xpack.go | 打包工具 |

参考: `pkg/`

---

### 4. config/ - 配置管理

- **config.go**: 配置加载、验证、环境变量支持
- **config_test.go**: 配置测试

**配置结构**:
```go
type Config struct {
    HTTP  HTTPConfig
    Redis RedisConfig
    Minio MinioConfig
    Baidu BaiduConfig
    Gitalk GitalkConfig
}
```

参考: `config/config.go:1`

---

### 5. static/ - 静态资源

- **static.go**: 前台静态资源嵌入
- **console/static.go**: 后台前端嵌入
- **css/**: 样式文件
- **js/**: JavaScript 文件
- **images/**: 图片资源
- **plugins/**: 插件
- **robots.txt**: 搜索引擎爬虫规则
- **favicon.ico**: 网站图标

使用 `embed` 包将静态资源嵌入到二进制文件中。

参考: `static/static.go:1`

---

### 6. global/ - 全局变量

- **global.go**: 全局 GORM 数据库实例

```go
var DB *gorm.DB
```

参考: `global/global.go:1`

---

## 文件命名规范

### Go 文件
- **小写 + 下划线**: `user_service.go`, `post_handler.go`
- **测试文件**: `*_test.go`
- **接口文件**: `interfaces.go`

### 配置文件
- **YAML 格式**: `config.yaml`, `config-example.yaml`
- **环境变量**: `KEEPBLOG_*`

### 构建文件
- **Makefile**: 构建命令
- **Dockerfile**: Docker 镜像构建
- **docker-compose.yml**: Docker 编排
- **entrypoint.sh**: 容器启动脚本

---

## 模块依赖关系

```
main.go
  ↓
internal/app/app.go (应用初始化)
  ↓
internal/app/routers.go (路由注册)
  ↓
internal/router/*.go (路由定义)
  ↓
api/admin/*.go, api/web/*.go (API 处理器)
  ↓
internal/service/*.go (业务逻辑)
  ↓
internal/model/*.go (数据模型)
  ↓
global/global.go (数据库实例)
```

**中间件流程**:
```
HTTP 请求
  ↓
全局中间件 (日志、恢复、CORS、压缩)
  ↓
路由中间件 (限流、缓存、JWT)
  ↓
API 处理器
  ↓
响应返回
```

---

## 代码组织原则

### 1. 分层架构
- **API 层**: 处理 HTTP 请求和响应
- **Service 层**: 业务逻辑处理
- **Model 层**: 数据访问和持久化

### 2. 依赖注入
- 通过 `Context` 注入服务依赖
- 避免全局变量（除数据库实例）

### 3. 接口抽象
- 定义服务接口 (`interfaces.go`)
- 便于测试和扩展

### 4. 模块化
- 每个功能模块独立目录
- 清晰的模块边界

### 5. 可测试性
- 单元测试文件 (`*_test.go`)
- 接口 mock 支持

---

## 下一步

- [系统架构](./04-architecture.md) - 了解整体架构设计
- 迁移文件：`internal/core/migrations/sql/`
- API 路由：`internal/router/`

---

**相关文件**:
- `main.go:1` - 应用入口
- `internal/app/app.go:1` - 应用初始化
- `internal/service/service.go:1` - 服务聚合
