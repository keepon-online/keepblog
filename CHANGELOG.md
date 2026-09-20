# 更新日志

所有此项目的显著更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)，
并且本项目遵循 [Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

> 2026-09 起发版 tag 改用语义化版本号（形如 `v2.1.2`）。此前递增数字 tag 的对应关系：
> v6=2.0.0、v7=2.0.1、v8=2.0.2、v9=2.1.0、v10=2.1.1。

## [未发布]

### 新增
- 网站资讯卡片数据补全，总字数/访客数/访问量/最后更新时间全部改为服务端真实数据：
  - 字数在文章保存/更新时按正文计算入库（新增 `md.CountWords`：CJK 逐字计数，空白与标点不计），启动时自动回填存量文章（幂等，不触碰 `last_modified_time`）；
  - 访客数/访问量改用本站 `system_access_log` 访问统计（PV 求和、UV 去重 IP，与后台 dashboard 口径一致），移除已失效的第三方 busuanzi 脚本及其全部引用；
  - "最后更新时间"由服务端输出 RFC3339 时间到 `data-lastPushDate`，前端据此显示相对时间。

### 修复
- 修复侧边栏"最后更新时间"被前端覆盖为 1970-1-1 的问题：模板此前缺失 `data-lastPushDate` 属性，`main.js` 以空值计算相对时间得到纪元日期；现属性缺失时保留服务端渲染文本。
- 修复访问统计中间件 PV 自增的并发竞态：由"读旧值 +1 回写"改为 SQL 原子自增（`pv + 1`）。
- 文章页"阅读时长"此前直接显示字数，现按 300 字/分钟向上取整换算为分钟（`readtime` 模板函数）。

## [2.2.1] - 2026-09-19

### 修复
- **修复前台/后台分组中间件从未生效的存量缺陷**：路由此前各自注册在 Engine 的新建分组上，`webGroup`（IP 限流/访问统计/页面缓存）与 `adminGroup`（错误渲染/指标/限流）的中间件全部空转。现页面路由与后台路由真正挂在对应分组下继承中间件；后台鉴权统一由分组级 JwtVerify 承担（自带 login/refreshToken 白名单），并顺带堵住此前未挂鉴权的少量后台路由。
- 修复 `result.Error/Fail` 与 JWT 中间件因提前写响应头导致 ErrorHandler 跳过渲染、错误响应 body 为空的问题；错误响应现完整输出结构化 JSON（code/request_id/timestamp/path）。
- GitHub 镜像仓库设为公开：私有仓库 Release 资产无法被服务器未认证下载（每日 03:00 更新 404 的根因），公开后链路恢复。

## [2.2.0] - 2026-09-18

> 2.1.3 曾记录于日志但未实际发版，其内容并入本版。

### 新增
- 增加 IP 地区数据库的定时更新能力：应用每天自动检查并更新数据，配套提供独立的校验与上传流程。

### 变更
- 清理 `pkg/md/` 下与 Markdown 渲染无关的历史测试文档和示例文件，保留渲染实现与必要测试。
- IP 查询样本基线统计（ipdb 数据源评估阶段一启动）：命中/未命中/无库/内网/云 CDN 出口计数与国家·省份·城市·ISP 字段空置率，每日 03:10 输出 JSON 快照到日志（`ipdb_stats` 前缀）；纯内存原子计数，不保存任何原始 IP。
- **IP 数据库更新闭环**：运行时每日更新只从自管 GitHub Release（`ipdb-latest`，SHA-256 校验）下载，不再跟踪上游 master；官方直链仅作冷启动兜底。CI 每周构建校验后自动发布归档与滚动 Release。
- 新增 `ipdb.updateUrl` 配置：指向完整 xdb URL（国内镜像或 `ipdb-v*` 归档回滚直链），支持热重载；环境变量 `KEEPBLOG_IPDB_UPDATE_URL`。
- `data/ip2region.xdb`（约 11MB）移出 git 追踪并加入 `.gitignore`，由冷启动下载维护；已有部署的 data 卷不受影响。

### 修复
- 移除已失效的 `goldmark-stats` 外部依赖，将 Markdown 字数与阅读时长统计逻辑迁入项目内部，避免依赖仓库不可用导致构建或维护受阻。

## [2.1.2] - 2026-09-18

### 变更
- 发版版本号由递增数字（v1…v10）改用标准语义化版本 `vX.Y.Z`，CHANGELOG 历史条目按变更幅度同步映射（对应关系见文件头注记），`release.sh` 校验与提示同步更新。
- 评估确认代码高亮自始由后端 goldmark+chroma 服务端实现，据此清理前端死代码：删除从未被模板加载的 `static/plugins/prismjs/` 整目录、5 个未引用的 canvas 特效，以及 main.js 中 Hexo 时代的 `addHighlightTool` 死分支及其专用配置（`highlight`/`copy`/`isHighlightShrink`），消除依赖扫描对 prismjs≤1.29 的 CVE-2024-53382 误报。

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
