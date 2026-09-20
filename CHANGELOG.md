# 更新日志

所有此项目的显著更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)，
并且本项目遵循 [Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

> 2026-09 起发版 tag 改用语义化版本号（形如 `v2.1.2`）。此前递增数字 tag 的对应关系：
> v6=2.0.0、v7=2.0.1、v8=2.0.2、v9=2.1.0、v10=2.1.1。

## [2.6.0] - 2026-09-20

### 新增
- **阅读进度条**：导航区域新增页面阅读进度指示，随滚动位置实时更新，支持长页面阅读反馈。
- **页脚波浪动画**：新增 SVG 多层波浪视差装饰，并针对移动端与 `prefers-reduced-motion` 做适配。
- **文章页头部视觉兜底**：有封面文章使用真实封面图作为头图；无封面文章使用星空闪烁与流星动画背景。

### 变更
- **搜索弹窗交互升级**：新增搜索图标、关键词输入提示、清空按钮、搜索历史标签、单条/全部历史清理、键盘导航提示及更完整的加载/空结果状态。
- **代码块工具栏样式统一**：复制按钮与代码块容器改为集中式样式管理，避免动态注入样式造成 PJAX 重复初始化问题。
- **Font Awesome 子集同步**：源码扫描列表补充本批次使用的搜索相关图标，并将不可用的图标引用替换为已有子集图标，避免线上显示方框。
- 继续统一搜索弹窗、文章卡片、标签云及响应式页面的圆角、阴影、间距和交互反馈。
- **文章元信息重构**：日期、分类、标签改为语义化胶囊徽章，日期使用 `<time datetime>`，分类缺失时自动隐藏空徽章。
- **代码块工具栏升级**：新增 Mac 风格红黄绿指示灯、语言标签和复制状态提示；复制按钮支持 PJAX 重复初始化且不会叠加。
- **标签云视觉优化**：调整为渐变配色、圆角胶囊、悬停阴影与键盘 `focus-visible` 状态，提升可读性与交互反馈。
- 统一文章卡片、侧边栏、导航、页脚与移动端的间距、圆角、阴影和交互动效。
- 补充文章分类空值守卫、备案链接安全属性及响应式样式。

## [2.5.2] - 2026-09-20

### 修复
- **首页打字机改为原生实现**：typed.js 库经官方 CDN 比对无误、最小页面可用，但在本页多脚本环境下初始化后不推进（真 Chrome 复现：光标出现、文字永远为空）。弃用该库并以约 40 行原生 JS 重写（打字/停留/删除/循环），保留一言 API 取句、失败回退站名、PJAX 返回重建三重行为；移除 typed.js 插件依赖。经真 Chrome 虚拟时间密集采样验证全链路正常。

## [2.5.1] - 2026-09-20

### 修复
- 首页打字机两类失效：PJAX 从其他页返回首页时脚本不重执行导致永远停在静态站名（现监听 `pjax:complete` 对新元素重建，防重复初始化）；一言 API 失败/超时无兜底导致打字机不启动（现 5 秒超时后回退静态站名展示）。

## [2.5.0] - 2026-09-20

### 变更
- **评论系统由 gitalk 改为 Artalk（自托管）**：gitalk 已停更五年、clientSecret 明文暴露在前端、依赖 unpkg 与 api.github.com（大陆双重不可达），且上线以来真实评论为 0，予以移除。改为集成 [Artalk](https://artalk.js.org)（Go + SQLite 同栈，MIT 持续维护）：前端 JS/CSS 从自有 Artalk 服务端加载，不依赖公共 CDN；评论区滚动到可视区附近才加载并初始化（IntersectionObserver 懒加载，不阻塞文章首屏）；初始化脚本置于 PJAX 替换区内并全程 IIFE/var 包裹（修复 gitalk 时代换页后评论不渲染、重复声明崩溃的问题）。新增 `artalk` 配置段（enable/server/site），服务端部署步骤见 `docs/artalk-deploy.md`。
- 顺带修复：文章页 `GLOBAL_CONFIG_SITE.isPost/isToc` 标志此前被评论开关连带门控且位于 PJAX 替换区外——评论关闭时标志不设置、PJAX 换页后不重设；现无条件输出并移入替换区内。

## [2.4.1] - 2026-09-20

### 变更
- **文章页相关推荐卡片视觉与交互全面优化**：
  - 弃用传统 `inline-block` 杂糅 calc 布局，重构为现代响应式 CSS Grid，间距更均衡统一（桌面 3 列，平板 2 列，手机 1 列）；
  - 卡片圆角统一为 `8px` 并补齐卡片阴影与悬浮上浮动效（`translateY(-4px)` 与深度阴影），与全站卡片风格和谐一致；
  - 引入持久渐变遮罩（`linear-gradient`）：解决此前 hover 时图片变亮导致白字反光看不清的体验问题，文字始终具有高对比度与文字阴影；
  - 文本信息移至卡片底部对齐，避免居中遮挡封面图中心视觉焦点，文章发布日期补充日历图标，标题 hover 高亮主题色。

## [2.4.0] - 2026-09-20

### 新增
- SEO 基础设施批次（技术 SEO）：
  - 全站输出 `<link rel="canonical">`：文章页、标签/分类/归档详情页（含翻页页指向自身）、关于/友链页、首页及 `/page/N`，由各页面 handler 注入，公共拼接逻辑收口在 `internal/pkg.CanonicalURL`；
  - 文章页新增 BlogPosting 结构化数据（JSON-LD：headline/description/image/datePublished/dateModified/author/wordCount）与 `article:published_time`、`article:modified_time`、`twitter:title`、`twitter:description` 元信息；
  - 列表页 title 差异化：标签/分类详情页输出 `标签：X | 站名`、归档年月页输出 `2024年07月 归档 | 站名`、首页翻页输出 `第 N 页 | 站名`；标签/分类/归档详情页的 meta description 改为带文章数的页面级描述（`pagedesc`），文章页 description 改用文章摘要而非全站描述；
  - sitemap 扩容：文章条数上限由 50 放开到全量（5 万上限），并纳入标签/分类/归档/关于/友链聚合页；首页 lastmod 改取最新一篇文章时间，不再每次生成都写"今天"；
  - robots.txt 屏蔽 `/api/`、`/console`、`/search/`、`/daily` 等对搜索无价值的路径，节省抓取预算；
  - 首页站点标题由空打字机 span 改为 `h1` 内静态输出站点名（Typed.js 加载后照常替换为一言），无 JS 环境与爬虫可读到有效 H1；
  - 正文图片（goldmark 渲染）统一补 `loading="lazy"` 与 `decoding="async"`，改善 LCP；
  - 新增模板冒烟测试：全部模板随应用同款函数解析，并渲染 head.html 验证 canonical/差异化 title/JSON-LD 合法性。
- **IndexNow 主动收录推送**（Bing/Yandex 等兼容搜索引擎）：`indexnow` 配置段（enable/key/endpoint，见 config-example.yaml），每日 21:00 与百度推送共用一次文章查询推送全部已发布文章；自动在 `/{key}.txt` 提供协议要求的密钥文件；站点域名取 `web_site.url`。
- **Google/Bing 站长平台验证字段**：`web_site` 表新增 `google_site`/`bing_site` 列（启动迁移幂等补列），前台输出 `google-site-verification` 与 `msvalidate.01` meta，后台"备案统计"设置页新增对应输入项；接入 Google Search Console / Bing Webmaster 的前置条件就此齐备。
- **空摘要自动兜底**：文章 `summary` 为空时，文章页 meta description/og/JSON-LD 与 RSS description 改用正文纯文本前 120 字（新增 `md.Excerpt`：解析 AST 取文本，跳过代码块与行内代码，按 rune 截断）——此前空摘要文章的搜索描述完全为空。
- **分页条链接语义修正**：第 1 页统一指向列表基路径（首页 `/`、标签/分类/归档详情页同理，与 canonical 一致，不再生成 `/page/1` 这类重复 URL）；当前页不再链接到自身（输出无 href 锚点 + `aria-current="page"`）；禁用的上一页/下一页不再输出 `href="#"` 假链接；`/page/1` 旧地址 301 到 `/` 兜底；`HandleIndex` 收敛为 `HandleIndexWithPageSize` 的默认 pageSize 封装，消除两份重复的 HTML 构造逻辑。

- **修复分页条上一页/下一页箭头未渲染与视觉样式问题**：
  - `chevron-left`/`chevron-right` 补入 Font Awesome 子集，并修复子集 CSS 缺少 `.fa-solid, .fas` 与 `.fa-regular, .far` 的 `font-family: "Font Awesome 6 Free"` 声明导致的图标显示为方块问号 `[?]` 缺陷；子集脚本补齐全站扫描交叉校验与核心字体族规则兜底；
  - 分页视觉体验全面优化：当前激活页码文字改用纯白加粗（`#ffffff`）并配柔和青色光晕阴影，彻底解决浅灰文字对比度不足问题；
  - 链接容器改为 flex 撑满整个卡片，实现全卡片可点击与数字/图标完美居中；
  - 区分禁用翻页箭头与省略号：禁用箭头保留卡片轮廓与布局对齐（降低不透明度且无阴影），保持左右对称美观；省略号去卡片底与阴影作为文本分隔符展示；
  - 增加可交互页码的平滑 hover 上浮微动效，适配移动端窄屏尺寸。
- 消除软 404：文章 hashids 解析失败或文章不存在时由 200 改为真实 404 状态码；各列表页 handler 的服务端错误由 200 渲染 error.html 改为 500；NoRoute 404 页补传站点数据（此前传 nil 导致页面标题为空）。
- 百度定时推送的切片缺陷：`make([]string, len(content))` 先填满空串再 append，推送 body 前面有整排空行白白消耗配额，改为容量语义 `make([]string, 0, len(content))`。
- 移动端 viewport 移除 `maximum-scale=1.0, user-scalable=no` 缩放禁用（Lighthouse 移动可访问性扣分项）。

## [2.3.0] - 2026-09-20

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
