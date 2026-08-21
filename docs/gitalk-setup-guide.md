# Gitalk 评论系统配置指南

本指南将帮助您在 keepblog 项目中配置和启用 Gitalk 评论系统。

---

## 📋 前提条件

- GitHub 账号
- 一个公开的 GitHub 仓库（用于存储评论）
- 项目已部署并可以通过域名访问

---

## 🎯 配置步骤

### 步骤 1: 创建 GitHub 仓库

创建一个专门用于存储评论的 GitHub 仓库：

1. 访问 https://github.com/new
2. 填写仓库信息：
   - **Repository name**: `blog-comments`（推荐名称）
   - **Description**: `博客评论存储仓库`
   - **Public**: 必须选择公开仓库
   - **Initialize with README**: 可选
3. 点击 **Create repository**

---

### 步骤 2: 创建 GitHub OAuth 应用

1. 访问 GitHub OAuth 应用创建页面：
   ```
   https://github.com/settings/applications/new
   ```

2. 填写应用信息：

   | 字段 | 值 | 说明 |
   |------|-----|------|
   | **Application name** | `你的博客名称 - 评论系统` | 例如：`小助理博客 - 评论系统` |
   | **Homepage URL** | `https://yourdomain.com` | 你的博客域名（必须是 HTTPS） |
   | **Application description** | `博客评论系统` | 可选 |
   | **Authorization callback URL** | `https://yourdomain.com` | 你的博客域名（必须是 HTTPS） |

3. 点击 **Register application**

4. 创建成功后，你会看到：
   - **Client ID**: 类似 `4b719159ee203950a239`
   - **Client Secret**: 点击 **Generate a new client secret** 生成，类似 `b9d643cffab418205da3c8c1a1bf762a9d19817b`

⚠️ **重要**: 请妥善保管 Client Secret，它只会显示一次！

---

### 步骤 3: 配置 config.yaml

编辑项目根目录的 `config.yaml` 文件：

```yaml
gitalk:
  enable: true                                    # 启用 Gitalk
  clientID: '4b719159ee203950a239'               # 替换为你的 Client ID
  clientSecret: 'b9d643cffab418205da3c8c1a1bf762a9d19817b'  # 替换为你的 Client Secret
  repo: 'blog-comments'                           # 评论仓库名称
  owner: 'your-github-username'                   # 你的 GitHub 用户名
  admin:
    - 'your-github-username'                      # 管理员列表（可以多个）
```

**配置说明**：

| 配置项 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `enable` | boolean | 是 | 是否启用 Gitalk，设置为 `true` |
| `clientID` | string | 是 | GitHub OAuth 应用的 Client ID |
| `clientSecret` | string | 是 | GitHub OAuth 应用的 Client Secret |
| `repo` | string | 是 | 存储评论的 GitHub 仓库名称 |
| `owner` | string | 是 | GitHub 仓库所有者的用户名 |
| `admin` | array | 是 | 管理员列表，可以初始化评论 |

---

### 步骤 4: 重启应用

配置完成后，重启应用使配置生效：

```bash
# 本地运行
make run

# Docker 部署
docker-compose restart
```

---

## ✅ 验证配置

### 1. 检查配置是否生效

访问任意文章页面，例如：
```
https://yourdomain.com/post/your-post-slug
```

在文章底部应该能看到评论区域。

### 2. 初始化评论

第一次访问文章时，Gitalk 会显示：
```
未找到相关的 Issues 进行评论
请联系 @your-github-username 初始化创建
```

**初始化步骤**：
1. 使用管理员账号（配置中的 `admin`）登录 GitHub
2. 点击 **使用 GitHub 登录** 按钮
3. 授权应用访问你的 GitHub 账号
4. Gitalk 会自动在评论仓库中创建一个 Issue
5. 初始化完成后，评论区域就可以正常使用了

### 3. 测试评论功能

1. 使用 GitHub 账号登录
2. 在评论框中输入评论内容
3. 点击 **Comment** 按钮
4. 评论会显示在页面上，同时在 GitHub 仓库的 Issues 中也能看到

---

## 🎨 前端集成说明

系统已经完全集成了 Gitalk，相关代码位置：

### 1. 模板文件 (`internal/web/post/post.html`)

**CSS 引入**（第 4 行）：
```html
<link href="https://unpkg.com/gitalk/dist/gitalk.css" rel="stylesheet">
```

**评论区域**（第 62-74 行）：
```html
{{if .gitalk.Enable}}
<div id="post-comment">
    <div class="comment-head">
        <div class="comment-headline">
            <i class="fas fa-comments fa-fw"></i>
            <span> 评论</span>
        </div>
    </div>
    <div class="comment-wrap">
        <div>
            <div id="gitalk-container"></div>
        </div>
    </div>
</div>
{{ end }}
```

**Gitalk 初始化**（第 87-106 行）：
```html
{{if .gitalk.Enable}}
<script src="https://unpkg.com/gitalk/dist/gitalk.min.js"></script>
<script>
    let gitalk = new Gitalk({
        clientID: '{{.gitalk.ClientID}}',
        clientSecret: '{{.gitalk.ClientSecret}}',
        repo: '{{.gitalk.Repo}}',
        owner: '{{.gitalk.Owner}}',
        admin: '{{.gitalk.Admin}}',
        language: "zh-CN",
        title: "{{.posts.Title}}",
        id: location.pathname,
        distractionFreeMode: false
    })
    gitalk.render('gitalk-container')
</script>
{{ end }}
```

### 2. 后端配置传递 (`api/web/post/posts.go:40`)

```go
c.HTML(http.StatusOK, "post.html", gin.H{
    "posts":   posts,
    "gitalk":  config.Get().Gitalk,  // 传递 Gitalk 配置
    // ... 其他数据
})
```

### 3. 配置管理 (`config/config.go`)

```go
type gitalk struct {
    Enable       bool     `yaml:"enable"`
    ClientID     string   `yaml:"clientID"`
    ClientSecret string   `yaml:"clientSecret"`
    Repo         string   `yaml:"repo"`
    Owner        string   `yaml:"owner"`
    Admin        []string `yaml:"admin"`
}
```

---

## 🔧 高级配置

### 自定义 Gitalk 选项

如需自定义 Gitalk 的更多选项，可以修改 `internal/web/post/post.html` 中的初始化代码：

```javascript
let gitalk = new Gitalk({
    clientID: '{{.gitalk.ClientID}}',
    clientSecret: '{{.gitalk.ClientSecret}}',
    repo: '{{.gitalk.Repo}}',
    owner: '{{.gitalk.Owner}}',
    admin: '{{.gitalk.Admin}}',
    language: "zh-CN",                    // 语言：zh-CN, en, es-ES 等
    title: "{{.posts.Title}}",            // Issue 标题
    id: location.pathname,                // Issue 唯一标识
    distractionFreeMode: false,           // 无干扰模式
    pagerDirection: 'last',               // 评论排序：last（最新在前）, first（最旧在前）
    createIssueManually: false,           // 是否手动创建 Issue
    proxy: 'https://cors-anywhere.azm.workers.dev/https://github.com/login/oauth/access_token',  // 代理（可选）
    perPage: 10,                          // 每页评论数
    enableHotKey: true,                   // 启用快捷键
})
```

完整选项列表请参考：https://github.com/gitalk/gitalk#options

---

## 🐛 常见问题

### 1. 评论区域不显示

**原因**：
- `enable` 未设置为 `true`
- 配置文件格式错误

**解决方案**：
```bash
# 检查配置文件
cat config.yaml | grep -A 6 "gitalk:"

# 确保 enable: true
```

### 2. 提示 "Error: Not Found"

**原因**：
- 仓库名称错误
- 仓库不是公开的
- 仓库所有者用户名错误

**解决方案**：
1. 检查仓库是否存在：`https://github.com/your-username/blog-comments`
2. 确保仓库是 **Public**
3. 检查 `owner` 配置是否正确

### 3. 提示 "Error: Validation Failed"

**原因**：
- Issue 标题过长（超过 50 个字符）
- `id` 参数过长

**解决方案**：
修改 `internal/web/post/post.html` 中的 `id` 参数：

```javascript
// 使用 MD5 哈希缩短 ID
id: md5(location.pathname),  // 需要引入 MD5 库

// 或使用文章 ID
id: '{{.posts.PostId}}',
```

### 4. 提示 "Error: Unauthorized"

**原因**：
- Client ID 或 Client Secret 错误
- OAuth 应用配置错误

**解决方案**：
1. 检查 Client ID 和 Client Secret 是否正确
2. 检查 OAuth 应用的 **Authorization callback URL** 是否正确
3. 重新生成 Client Secret

### 5. 无法初始化评论

**原因**：
- 未使用管理员账号登录
- 管理员列表配置错误

**解决方案**：
```yaml
gitalk:
  admin:
    - 'your-github-username'  # 确保用户名正确
```

### 6. 评论显示乱码

**原因**：
- 语言设置错误

**解决方案**：
```javascript
language: "zh-CN",  // 中文
// 或
language: "en",     // 英文
```

---

## 🔒 安全建议

### 1. 保护 Client Secret

⚠️ **不要将 Client Secret 提交到公开仓库**

**方法 1：使用环境变量**

```bash
# 设置环境变量
export KEEPBLOG_GITALK_CLIENT_SECRET="your-client-secret"
```

修改 `config/config.go` 添加环境变量绑定：
```go
viper.BindEnv("gitalk.clientSecret", "KEEPBLOG_GITALK_CLIENT_SECRET")
```

**方法 2：使用 .gitignore**

```bash
# .gitignore
config.yaml
```

只提交 `config-example.yaml` 作为示例。

### 2. 限制 OAuth 应用权限

在 GitHub OAuth 应用设置中：
- 只授予必要的权限
- 定期检查授权的应用

### 3. 定期更新 Client Secret

建议每 6-12 个月更新一次 Client Secret。

---

## 📊 评论管理

### 查看评论

所有评论都存储在 GitHub 仓库的 Issues 中：
```
https://github.com/your-username/blog-comments/issues
```

### 管理评论

作为管理员，你可以：
- **编辑评论**：在 GitHub Issue 中编辑
- **删除评论**：关闭或删除 Issue
- **回复评论**：在 Issue 中回复
- **标记垃圾评论**：添加标签或关闭 Issue

### 评论通知

GitHub 会自动发送邮件通知：
- 有新评论时
- 有人回复你的评论时
- 有人 @ 你时

---

## 🎉 完成

配置完成后，您的博客就拥有了完整的评论功能！

**功能特性**：
- ✅ GitHub 账号登录
- ✅ Markdown 支持
- ✅ 表情支持
- ✅ 评论回复
- ✅ 邮件通知
- ✅ 评论管理
- ✅ 反垃圾评论（GitHub 自带）

---

## 📚 相关资源

- **Gitalk 官方文档**: https://github.com/gitalk/gitalk
- **GitHub OAuth 文档**: https://docs.github.com/en/developers/apps/building-oauth-apps
- **Gitalk 在线演示**: https://gitalk.github.io/

---

## 🆘 需要帮助？

如果遇到问题：
1. 查看本文档的 **常见问题** 部分
2. 查看 Gitalk 官方 Issues: https://github.com/gitalk/gitalk/issues
3. 在项目中提交 Issue

---

**配置文件位置**:
- 配置文件: `config.yaml`
- 配置示例: `config-example.yaml`
- 配置代码: `config/config.go:53-60`
- 模板文件: `internal/web/post/post.html:62-106`
- 后端传递: `api/web/post/posts.go:40`

**最后更新**: 2026-01-28
