# 更新日志

所有此项目的显著更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)，
并且本项目遵循 [Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

## [未发布]

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

### 变更
- 默认文章与默认音乐种子数据外置到 `internal/core/seeddata/`。

### 新增
-

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
