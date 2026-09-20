# Artalk 评论服务端部署

go-site 前台已集成 [Artalk](https://artalk.js.org) 评论（替换已停更的 gitalk）。
评论后端独立部署，与 go-site 通过 HTTP 交互；前端 JS/CSS 直接从 Artalk
服务端加载，不依赖 unpkg/github.com 等公共 CDN（大陆可达性由自己服务器决定）。

## 1. 部署 Artalk 服务端

在生产服务器（已有 Docker 环境）新增一个 compose 服务，例如并入现有
`docker-compose.yml` 或独立成文件：

```yaml
services:
  artalk:
    image: artalk/artalk-go
    container_name: artalk
    restart: unless-stopped
    ports:
      - "127.0.0.1:23366:23366"   # 建议只监听本机，由反代暴露 HTTPS
    volumes:
      - ./data:/data              # SQLite 与配置持久化
    environment:
      - TZ=Asia/Shanghai
      - ATK_LOCALE=zh-CN
      - ATK_SITE_DEFAULT=小助理    # 站点名，与 go-site 配置的 artalk.site 一致
      - ATK_SITE_URL=https://www.keepon.online
      - ATK_APP_KEY=<openssl rand -hex 16 生成>
```

启动后创建管理员：

```bash
docker compose up -d
docker exec -it artalk artalk admin
```

## 2. 反向代理

为评论 API 挂一个 HTTPS 域名（如 `comment.keepon.online`）指向 `127.0.0.1:23366`，
并在反代层启用 WebSocket 支持（站内通知需要）：

```nginx
location / {
    proxy_pass http://127.0.0.1:23366;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

在 Artalk 控制中心（管理员在评论框登录后右下角入口）确认站点
**受信任域名**包含 `https://www.keepon.online`（跨域白名单，默认会自动加
ATK_SITE_URL）。

## 3. go-site 侧配置

`config.yaml`：

```yaml
artalk:
  enable: true
  server: "https://comment.keepon.online"
  site: "小助理"
```

重启 go-site 后，文章页滚动到评论区附近才会加载 Artalk 资源并初始化
（懒加载，不阻塞文章首屏；PJAX 换页安全）。

## 4. 邮件通知等

邮件 SMTP、多渠道推送（Telegram 等）、验证码与审核策略均在 Artalk
控制中心图形界面配置，无需改 go-site。图片上传可配置自定义图床
（兼容 S3/MinIO）。

## 迁移说明

原 gitalk 评论数为 0，无需迁移数据。旧 `gitalk` 配置段在升级后可删除
（未知 yaml 键会被忽略，留着也不影响）。
