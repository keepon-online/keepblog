# Gitalk 评论系统快速配置

## 🚀 5 分钟快速配置

### 1️⃣ 创建 GitHub 仓库

访问 https://github.com/new 创建公开仓库：
- 仓库名称：`blog-comments`
- 可见性：**Public**（必须）

### 2️⃣ 创建 OAuth 应用

访问 https://github.com/settings/applications/new

填写信息：
```
Application name: 你的博客名称 - 评论系统
Homepage URL: https://yourdomain.com
Authorization callback URL: https://yourdomain.com
```

获取：
- **Client ID**: `4b719159ee203950a239`
- **Client Secret**: `b9d643cffab418205da3c8c1a1bf762a9d19817b`

### 3️⃣ 修改配置文件

编辑 `config.yaml`：

```yaml
gitalk:
  enable: true
  clientID: '你的-Client-ID'
  clientSecret: '你的-Client-Secret'
  repo: 'blog-comments'
  owner: '你的-GitHub-用户名'
  admin:
    - '你的-GitHub-用户名'
```

### 4️⃣ 重启应用

```bash
# 本地运行
make run

# Docker 部署
docker-compose restart
```

### 5️⃣ 初始化评论

1. 访问任意文章页面
2. 使用管理员账号登录 GitHub
3. 点击 **使用 GitHub 登录**
4. 授权应用
5. 完成！

---

## ✅ 验证

访问文章页面，底部应该显示评论区域。

---

## 🐛 常见问题

### 评论区域不显示
- 检查 `enable: true`
- 检查配置文件格式

### Error: Not Found
- 确保仓库是 **Public**
- 检查仓库名称和所有者

### Error: Unauthorized
- 检查 Client ID 和 Client Secret
- 检查 OAuth 应用的回调 URL

---

## 📖 详细文档

查看完整配置指南：[gitalk-setup-guide.md](./gitalk-setup-guide.md)

---

**配置示例**：

```yaml
# 完整配置示例
gitalk:
  enable: true                                    # 启用 Gitalk
  clientID: '4b719159ee203950a239'               # GitHub OAuth Client ID
  clientSecret: 'b9d643cffab418205da3c8c1a1bf762a9d19817b'  # GitHub OAuth Client Secret
  repo: 'blog-comments'                           # 评论仓库名称
  owner: 'your-github-username'                   # GitHub 用户名
  admin:                                          # 管理员列表
    - 'your-github-username'
```

---

**相关链接**：
- GitHub OAuth 应用：https://github.com/settings/developers
- Gitalk 官方文档：https://github.com/gitalk/gitalk
- 评论仓库示例：https://github.com/your-username/blog-comments
