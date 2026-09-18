# 更新日志

所有此项目的显著更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)，
并且本项目遵循 [Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

> 2026-09 起发版 tag 改用语义化版本号（形如 `v2.1.2`）。此前递增数字 tag 的对应关系：
> v6=2.0.0、v7=2.0.1、v8=2.0.2、v9=2.1.0、v10=2.1.1。

## [未发布]

## [2.1.1] - 2026-09-17

### 修复
- 首页音乐播放器自动播放被浏览器策略拦截：`play()` 仅由真实用户手势触发，消除 `NotAllowedError`。
- 静态资源缓存由 `immutable` 长缓存改为 `no-cache` + 内容哈希 ETag 协商缓存，避免发版后访客端旧资源被钉死最长一年。

### 变更
- 发版脚本同步 GitHub 远端时 branch 与 tag 分两次推送，避免合并为一条 push 命令时 GitHub 只为其中一个 ref 生成事件、镜像构建不触发。

## [2.1.0] - 2026-09-15

### 新增
- GitHub Actions（私有镜像仓库 `keepon-online/keepblog`）：push master 自动跑 `go vet` + `go test`；打 `v*` tag 自动构建多阶段镜像并推送 Docker Hub（`jieepre/keepblog`），支持 GHA 层缓存。

### 变更
- `release.sh` 不再本地构建镜像：打 tag 后同步推送 Gitee 与 GitHub 双远端，镜像构建推送全部交给 Actions，本地无需 Docker 构建环境。

## [2.0.2] - 2026-09-15

### 修复
- 后台登录成功后无法进入控制台：2.0.0 将 `GetAsyncRoutes` 响应统一为 `{code, message, payload}` 信封后，前端 `router/utils.ts` 仍按旧字段解构 `{ data }`，得到 `undefined` 导致路由初始化 Promise 永远挂起。现改用 `payload`。
- `config-example.yaml` 的 `jwt.secret` 恢复为空占位：测试值会被打包进镜像默认配置，未显式配置 `JWT_SECRET` 的部署将以公开弱密钥运行（entrypoint 只自动替换空 secret）。

### 变更
- 移除 Drone CI（`.drone.yml`）：CI 镜像推送长期失效（最后成功推送停留在 2025-12），发布流程统一为 `scripts/release.sh`（tag 驱动 + 本地构建推送镜像）。

## [2.0.1] - 2026-09-14

### 修复
- Docker 首次部署 JWT 密钥自动生成失效：`entrypoint.sh` 原判断条件（配置中不存在 `jwt:` 段才生成）永远为假，导致空 secret 直接启动失败。现改为检测空 secret 并原地替换，或段缺失时追加。
- 初始化失败时进程不退出：gookit/slog 的 `Fatalf` 只记日志不退出，`main.go` 继续执行导致空指针 panic 淹没真实错误。现显式 `os.Exit(1)`。

### 变更
- CI（Drone）Docker 推送仓库由旧名 `jieepre/site` 改为 `jieepre/keepblog`，与镜像更名保持一致。
- `release.sh` 发版时自动推送镜像到 Docker Hub（原先需手动 `make docker-push`，易遗漏）。
- Dockerfile 构建阶段先清理 `static/console` 旧产物再拷贝前端构建结果，避免过期 hash 文件被 embed 进二进制。
- `.dockerignore` 更新：排除 `.git`、`frontend/node_modules`、`config.yaml` 等，移除旧产物名 `go-site`。

## [2.0.0] - 2026-08-22

### 变更
- **项目更名为 KeepBlog**：Go module 路径由 `gitee.com/jieepre/go-site` 改为 `gitee.com/jieepre/keepblog`，二进制产物、Docker 镜像名（`jieepre/keepblog`）、前端包名同步更新。
- **环境变量前缀变更**：`GOSITE_*` 改为 `KEEPBLOG_*`（如 `KEEPBLOG_HTTP_PORT`）；`JWT_SECRET` 不变。使用环境变量覆盖配置的部署需同步修改。
- minio 默认 bucket 名由 `go-site` 改为 `keepblog`（已显式配置 bucketName 的部署不受影响）。
- 默认文章与默认音乐种子数据外置到 `internal/core/seeddata/`。

### 重构
- 完成 service 层数据库依赖注入、版本化 Goose 迁移、归档查询 N+1 治理与阅读数原子自增。
- 升级 Gin 至 1.11.0、go-redis 至 v9.6.1、golang-jwt 至 v5.2.2。
- 移除 gogf/gf 与 pkg/errors，监控模块拆分并抽取共享 GORM 查询构建器。
- 数据库错误使用真实 HTTP 状态码和结构化 ErrorResponse；成功响应信封保持兼容。

### 安全
- 移除 JWT 签名密钥的代码内置默认值：必须在 `config.yaml` 的 `jwt.secret` 或环境变量 `JWT_SECRET` 提供，缺失时启动报错退出；Docker 首次部署自动生成随机密钥。
- hashids salt 可通过 `hashids.salt` 配置覆盖（默认保持历史值，旧文章 URL 不受影响）。
- 管理员初始密码改为首次启动随机生成并打印一次（原为固定弱密码），种子数据不再包含真实手机号/邮箱。
- `config.yaml` 移出 git 追踪（含真实密钥的配置不应入库）；**注意：历史提交中仍可见旧密钥，建议轮换 minio/gitalk/百度推送凭据**。
- `config-example.yaml` 中的真实样例密钥替换为占位符。

## [1.0.0] - 2025-07-08

### 新增
- 项目初始版本。
- 用户认证功能。
- 用于管理文章、分类和标签的后台仪表盘。
- 用于展示博客内容的前端网站。
- 所有主要功能的API端点。

### 变更
- 更新数据库结构以支持新功能。
- 重构路由逻辑以实现更好的组织。

### 修复
- 修正了文章分页逻辑中的一个错误。
- 解决了文件上传过程中的一个安全漏洞。
