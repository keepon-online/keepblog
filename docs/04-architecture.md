# 系统架构

本文档详细介绍 Go-Site 项目的系统架构设计、数据流和核心组件。

---

## 整体架构

### 分层架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP 客户端                              │
│                  (浏览器、移动端、API 客户端)                   │
└─────────────────────────────────────────────────────────────┘
                            ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                    Gin Web 框架                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              全局中间件层                              │   │
│  │  • 日志记录 (GinLogger)                               │   │
│  │  • 恐慌恢复 (GinRecovery)                             │   │
│  │  • CORS 跨域 (Cors)                                   │   │
│  │  • 响应压缩 (Gzip)                                    │   │
│  └──────────────────────────────────────────────────────┘   │
│                            ↓                                 │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              路由层 (Router)                          │   │
│  │  • 前台路由 (web)                                     │   │
│  │  • 后台路由 (admin)                                   │   │
│  │  • 健康检查 (health)                                  │   │
│  │  • 性能指标 (metrics)                                 │   │
│  │  • WebSocket (ws)                                     │   │
│  └──────────────────────────────────────────────────────┘   │
│                            ↓                                 │
│  ┌──────────────────────────────────────────────────────┐   │
│  │            路由中间件层                                │   │
│  │  • 限流 (RateLimit)                                   │   │
│  │  • 缓存 (Cache)                                       │   │
│  │  • JWT 验证 (JwtVerify)                               │   │
│  │  • 访问统计 (Statistics)                              │   │
│  │  • 错误处理 (ErrorHandler)                            │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    API 处理器层 (Handlers)                    │
│  ┌──────────────────┐         ┌──────────────────┐          │
│  │   前台 API       │         │   后台 API       │          │
│  │  • 首页          │         │  • 文章管理      │          │
│  │  • 文章详情      │         │  • 分类管理      │          │
│  │  • 分类浏览      │         │  • 标签管理      │          │
│  │  • 标签浏览      │         │  • 系统监控      │          │
│  │  • 归档          │         │  • 日志管理      │          │
│  │  • 关于          │         │  • 用户管理      │          │
│  │  • 友链          │         │  • 网站配置      │          │
│  │  • 音乐          │         │  • 通知管理      │          │
│  └──────────────────┘         └──────────────────┘          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  服务层 (Service Layer)                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              AppService (服务聚合)                    │   │
│  │  • PostService      (文章服务)                        │   │
│  │  • CategoryService  (分类服务)                        │   │
│  │  • TagService       (标签服务)                        │   │
│  │  • SidebarService   (侧边栏服务)                      │   │
│  │  • AboutService     (关于服务)                        │   │
│  │  • SystemService    (系统服务)                        │   │
│  │  • LinkService      (友链服务)                        │   │
│  │  • WebSiteService   (网站服务)                        │   │
│  │  • Dashboard        (仪表板服务)                      │   │
│  │  • MusicService     (音乐服务)                        │   │
│  │  • NoticeService    (通知服务)                        │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  模型层 (Model Layer)                         │
│  • Post, Category, Tag, PostTag                             │
│  • Comment, CommentReply                                    │
│  • FriendLink, User, About, Music                           │
│  • AccessLog, LoginLog, WebSite, Notice                     │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  数据访问层 (GORM)                            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────┐         ┌──────────────────┐
│   SQLite 数据库   │         │   Redis 缓存     │
│  (主数据存储)     │         │   (可选)         │
└──────────────────┘         └──────────────────┘
```

参考: `internal/app/routers.go:1`, `internal/service/service.go:1`

---

## 核心组件

### 1. 应用入口 (main.go)

应用启动流程：

```go
func main() {
    // 1. 加载配置
    config.Load()

    // 2. 初始化日志
    core.InitLog()

    // 3. 初始化数据库
    core.InitDB()

    // 4. 初始化 Redis (可选)
    cache.InitRedis()

    // 5. 初始化资源
    core.InitResource()

    // 6. 启动定时任务
    core.Timer()

    // 7. 创建应用
    app := app.NewApp()

    // 8. 启动服务器
    app.Run()
}
```

参考: `main.go:1`

### 2. 应用程序 (internal/app/app.go)

```go
type App struct {
    Engine  *gin.Engine
    Service *service.AppService
}

func NewApp() *App {
    // 创建 Gin 引擎
    engine := gin.New()

    // 初始化服务
    appService := service.NewAppService()

    // 创建应用
    app := &App{
        Engine:  engine,
        Service: appService,
    }

    // 注册路由
    app.setupRouters()

    return app
}
```

参考: `internal/app/app.go:1`

### 3. 路由系统 (internal/router/)

#### 路由注册流程

```go
func (a *App) setupRouters() {
    // 1. 全局中间件
    a.Engine.Use(
        middleware.GinLogger(),
        middleware.GinRecovery(),
        middleware.Cors(),
        gzip.Gzip(gzip.DefaultCompression),
    )

    // 2. 健康检查和监控
    a.Engine.GET("/health", monitor.HealthCheck)
    a.Engine.GET("/health/ready", monitor.ReadyCheck)
    a.Engine.GET("/health/live", monitor.LiveCheck)
    a.Engine.GET("/metrics", monitor.Metrics)

    // 3. 前台路由
    router.RegisterWebRoutes(a.Engine, a.Service)

    // 4. 后台路由
    router.RegisterAdminRoutes(a.Engine, a.Service)

    // 5. WebSocket
    a.Engine.GET("/ws", websocket.HandleWebSocket)
}
```

参考: `internal/app/routers.go:1`

#### 前台路由 (internal/router/web.go)

```go
func RegisterWebRoutes(engine *gin.Engine, service *service.AppService) {
    // 前台中间件
    web := engine.Group("/")
    web.Use(
        middleware.IPBasedRateLimit(60),  // 限流
        middleware.Statistics(),           // 统计
        middleware.CacheMiddleware(5),     // 缓存
    )

    // 注册路由
    home.RegisterRoutes(web, service)
    post.RegisterRoutes(web, service)
    category.RegisterRoutes(web, service)
    tags.RegisterRoutes(web, service)
    archive.RegisterRoutes(web, service)
    about.RegisterRoutes(web, service)
    link.RegisterRoutes(web, service)
    music.RegisterRoutes(web, service)
}
```

参考: `internal/router/web.go:1`

#### 后台路由 (internal/router/admin.go)

```go
func RegisterAdminRoutes(engine *gin.Engine, service *service.AppService) {
    // 后台中间件
    api := engine.Group("/api")
    api.Use(
        middleware.ErrorHandler(),
        middleware.MetricsMiddleware(),
        middleware.APIRateLimit(100),
        middleware.JwtVerify(),  // JWT 验证
    )

    // 注册路由
    post.RegisterRoutes(api, service)
    category.RegisterRoutes(api, service)
    tags.RegisterRoutes(api, service)
    link.RegisterRoutes(api, service)
    music.RegisterRoutes(api, service)
    about.RegisterRoutes(api, service)
    website.RegisterRoutes(api, service)
    dashboard.RegisterRoutes(api, service)
    monitor.RegisterRoutes(api, service)
    log.RegisterRoutes(api, service)
    system.RegisterRoutes(api, service)
    notice.RegisterRoutes(api, service)
}
```

参考: `internal/router/admin.go:1`

### 4. 服务层 (internal/service/)

#### 服务聚合

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

func NewAppService() *AppService {
    return &AppService{
        PostService:     post.NewService(),
        CategoryService: category.NewService(),
        TagService:      tag.NewService(),
        // ... 其他服务
    }
}
```

参考: `internal/service/service.go:1`

#### 服务接口

```go
type PostServiceInterface interface {
    CreatePost(req *request.PostRequest) error
    UpdatePost(postId int, req *request.PostRequest) error
    DeletePost(postId int) error
    GetPost(postId int) (*response.PostResponse, error)
    ListPosts(page, pageSize int) (*response.PostListResponse, error)
    PublishPost(postId int, isPublished bool) error
    TopPost(postId int, top bool) error
}
```

服务接口文件已移除；当前通过 `InitAppService(db *gorm.DB)` 进行显式依赖注入。

### 5. 中间件系统 (internal/middleware/)

#### 中间件链

```
请求 → 日志 → 恢复 → CORS → 压缩 → 限流 → 缓存 → JWT → 统计 → 处理器 → 响应
```

#### 主要中间件

**日志中间件** (`log.go`):
```go
func GinLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)

        logger.Info("HTTP Request",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "duration", duration,
        )
    }
}
```

**限流中间件** (`ratelimit.go`):
```go
func IPBasedRateLimit(requestsPerMinute int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(requestsPerMinute/60), requestsPerMinute)

    return func(c *gin.Context) {
        ip := c.ClientIP()
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**JWT 中间件** (`jwt.go`):
```go
func JwtVerify() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 白名单检查
        if isWhitelisted(c.Request.URL.Path) {
            c.Next()
            return
        }

        // 提取令牌
        token := extractToken(c)
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // 验证令牌
        claims, err := jwttoken.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 黑名单检查
        if jwttoken.IsBlacklisted(token) {
            c.JSON(401, gin.H{"error": "Token revoked"})
            c.Abort()
            return
        }

        c.Set("userId", claims.UserId)
        c.Next()
    }
}
```

参考: `internal/middleware/`

---

## 数据流

### 1. 前台请求流程

```
用户访问首页
    ↓
GET /
    ↓
全局中间件 (日志、恢复、CORS、压缩)
    ↓
前台中间件 (限流 60/分钟、统计、缓存 5 分钟)
    ↓
home.Index 处理器
    ↓
PostService.ListPosts()
    ↓
GORM 查询数据库
    ↓
返回文章列表
    ↓
渲染 HTML 模板
    ↓
响应返回给用户
```

参考: `api/web/home/home.go:1`, `internal/service/post/post.go:1`

### 2. 后台请求流程

```
管理员创建文章
    ↓
POST /api/v1/site/post/save
    ↓
全局中间件 (日志、恢复、CORS、压缩)
    ↓
后台中间件 (错误处理、指标、限流 100/分钟、JWT 验证)
    ↓
JWT 验证通过
    ↓
post.Save 处理器
    ↓
PostService.CreatePost()
    ↓
1. 验证请求参数
2. 处理 Markdown 内容
3. 生成 HTML
4. 提取摘要
5. 计算字数
6. 保存到数据库
7. 关联标签
8. 百度推送 (可选)
    ↓
返回成功响应
```

参考: `api/admin/post/post.go:1`, `internal/service/post/post.go:1`

### 3. 认证流程

```
用户登录
    ↓
POST /api/login
    ↓
login.Login 处理器
    ↓
1. 验证用户名和密码
2. 生成访问令牌 (15 分钟)
3. 生成刷新令牌 (7 天)
4. 记录登录日志
    ↓
返回令牌
    ↓
用户携带令牌访问受保护资源
    ↓
JWT 中间件验证令牌
    ↓
1. 检查白名单
2. 提取令牌
3. 验证签名
4. 检查过期时间
5. 检查黑名单
    ↓
验证通过，继续处理
```

参考: `api/admin/login/login.go:1`, `internal/middleware/jwt.go:1`, `pkg/jwttoken/token.go:1`

### 4. 缓存流程

```
用户访问首页
    ↓
缓存中间件检查
    ↓
Redis 中是否有缓存？
    ↓ 是
返回缓存内容
    ↓ 否
继续处理请求
    ↓
生成响应
    ↓
存入 Redis (5 分钟)
    ↓
返回响应
```

参考: `internal/middleware/cache.go:1`

---

## 依赖注入

### 请求上下文

```go
type Context struct {
    Engine  *gin.Engine
    Service *AppService
}

// API 处理器嵌入 Context
type PostHandler struct {
    *core.Context
}

func NewPostHandler(ctx *core.Context) *PostHandler {
    return &PostHandler{Context: ctx}
}

// 使用服务
func (h *PostHandler) Save(c *gin.Context) {
    // 访问服务
    err := h.Service.PostService.CreatePost(req)
    // ...
}
```

参考: `internal/pkg/core/context.go:1`, `api/admin/post/post.go:1`

---

## 数据库设计

### 连接池配置

```go
func InitDB() {
    db, err := gorm.Open(sqlite.Open("./data/site.db"), &gorm.Config{
        PrepareStmt: true,  // 预编译语句
    })

    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)           // 最大连接数
    sqlDB.SetMaxIdleConns(10)           // 最小空闲连接
    sqlDB.SetConnMaxLifetime(2 * time.Hour)  // 连接生命周期
    sqlDB.SetConnMaxIdleTime(30 * time.Minute)  // 空闲超时

    global.DB = db
}
```

参考: `internal/core/db.go:1`

### 版本化数据库迁移

数据库 schema 和索引由 Goose 管理，迁移文件嵌入二进制：

```go
if err := migrations.Run(global.GORM); err != nil {
    return fmt.Errorf("run database migrations failed: %w", err)
}
core.InitResource() // 只负责种子数据
```

首次接入的旧库会执行一次兼容性 AutoMigrate 以补齐基线缺失字段，随后由
`internal/core/migrations/sql/` 中的版本化 SQL 管理变更。迁移使用
`IF NOT EXISTS`，重复启动幂等且不删除已有数据。

参考：`internal/core/migrations/runner.go:Run`

### 1. 连接池

- **数据库连接池**: 最大 25 个连接，最小 10 个空闲连接
- **Redis 连接池**: 最大 100 个连接，最小 10 个空闲连接

### 2. 缓存策略

- **前台页面缓存**: 5 分钟 (Redis)
- **静态资源缓存**: 浏览器缓存
- **数据库查询缓存**: GORM 预编译语句

### 3. 限流

- **前台限流**: 60 次/分钟 (基于 IP)
- **后台限流**: 100 次/分钟 (基于 IP)
- **登录限流**: 5 次/分钟 (基于 IP)

### 4. 压缩

- **Gzip 压缩**: 默认压缩级别
- **静态资源**: 预压缩

### 5. 数据库优化

- **索引**: 关键字段创建索引
- **预编译语句**: 减少解析开销
- **批量操作**: 减少数据库往返

参考: `internal/middleware/ratelimit.go:1`, `internal/middleware/cache.go:1`

---

## 安全机制

### 1. 认证授权

- **JWT 令牌**: 访问令牌 + 刷新令牌
- **令牌黑名单**: 登出时加入黑名单
- **密码加密**: bcrypt 加密存储

### 2. 安全防护

- **CSRF 防护**: 令牌验证
- **SQL 注入防护**: 参数化查询
- **XSS 防护**: 输出转义
- **路径规范化**: 防止路径遍历

### 3. 访问控制

- **JWT 白名单**: 公开接口无需认证
- **路由权限**: 基于路由的权限控制
- **IP 限流**: 防止暴力攻击

参考: `internal/middleware/jwt.go:1`, `internal/middleware/csrf.go:1`

---

## 监控和日志

### 1. 健康检查

- **/health**: 基本健康检查
- **/health/ready**: 就绪检查（数据库连接）
- **/health/live**: 存活检查

### 2. 性能指标

- **/metrics**: Prometheus 格式的性能指标
- **请求计数**: 按路由统计
- **响应时间**: 按路由统计
- **错误率**: 按路由统计

### 3. 日志系统

- **结构化日志**: JSON 格式
- **日志级别**: DEBUG, INFO, WARN, ERROR
- **日志输出**: 文件 + 控制台
- **日志轮转**: 按大小和时间

参考: `internal/monitor/health.go:1`, `internal/monitor/metrics.go:1`, `internal/logger/logger.go:1`

---

## 扩展性

### 1. 服务扩展

添加新服务：

```go
// 1. 定义服务接口
type NewServiceInterface interface {
    DoSomething() error
}

// 2. 实现服务
type NewService struct{}

func (s *NewService) DoSomething() error {
    // 实现逻辑
    return nil
}

// 3. 添加到 AppService
type AppService struct {
    // ... 现有服务
    NewService *NewService
}
```

### 2. 中间件扩展

添加新中间件：

```go
func NewMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}

// 注册中间件
engine.Use(NewMiddleware())
```

### 3. 路由扩展

添加新路由：

```go
func RegisterNewRoutes(group *gin.RouterGroup, service *service.AppService) {
    handler := NewHandler(service)
    group.GET("/new", handler.Get)
    group.POST("/new", handler.Post)
}
```

---

## 下一步

- 数据库迁移：`internal/core/migrations/`
- API 路由：`internal/router/`
- 中间件：`internal/middleware/`

---

**相关文件**:
- `main.go:1` - 应用入口
- `internal/app/app.go:1` - 应用程序
- `internal/app/routers.go:1` - 路由注册
- `internal/service/service.go:1` - 服务聚合
- `internal/middleware/` - 中间件目录
