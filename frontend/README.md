# keepblog 管理后台

keepblog 博客系统的管理后台前端，基于 [Vue 3](https://vuejs.org/) + [Vite](https://vitejs.dev/) + [Element Plus](https://element-plus.org/) + [TypeScript](https://www.typescriptlang.org/) 构建，脚手架源自 [pure-admin-thin](https://github.com/pure-admin/pure-admin-thin)。

构建产物通过 Go 的 `embed` 嵌入后端二进制，经 `/console` 路径提供服务。

## 目录说明

本目录是主仓库 `keepblog` 的前端子项目，与后端共同构成 monorepo：

```
keepblog/
├── frontend/          ← 本目录：管理后台前端源码
│   ├── src/           应用源码
│   ├── build/         vite 构建辅助
│   ├── public/        静态资源
│   └── types/         全局类型声明
└── static/console/    ← 构建产物落地处（供 Go embed）
```

## 开发

```bash
# 安装依赖（要求 pnpm >=9、Node ^20.19 || >=22.13）
pnpm install

# 启动开发服务器（默认端口 8848）
# /api 请求会自动代理到本地后端 http://localhost:8589
pnpm dev
```

开发时需同时启动后端（在仓库根目录执行 `make dev-backend`），浏览器访问 `http://localhost:8848`。

## 构建

```bash
# 构建到 dist/
pnpm build
```

在 monorepo 语境下，推荐在仓库根目录使用 `make build-frontend`（构建后自动同步产物到 `static/console/`）或 `make build`（前端 + Go 一键构建并 embed）。

## 致谢

本项目基于 [pure-admin](https://github.com/pure-admin) 系列脚手架，感谢其开源贡献。

## 许可证

[MIT](./LICENSE)
