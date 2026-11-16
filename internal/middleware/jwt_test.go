package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/pkg/jwttoken"
	"github.com/gin-gonic/gin"
)

// 3.6.1 鉴权链路修复的回归测试：refreshToken 不得当 accessToken 用、
// 登出/改密后令牌即时失效、白名单不得被路径穿越绕过。

// newJwtTestRouter 按生产装配方式（ErrorHandler + JwtVerify）构建探针路由
func newJwtTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	r.Use(JwtVerify())
	r.Any("/api/v1/probe", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{"username": username})
	})
	return r
}

func doJwtProbe(r *gin.Engine, auth string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/probe", nil)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	r.ServeHTTP(w, req)
	return w
}

func TestJwtVerifyAcceptsAccessToken(t *testing.T) {
	if err := jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456")); err != nil {
		t.Fatalf("配置密钥失败: %v", err)
	}
	token, err := jwttoken.AccessToken("admin")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}

	w := doJwtProbe(newJwtTestRouter(), "Bearer "+token)
	if w.Code != http.StatusOK {
		t.Fatalf("access token 应放行，状态码 = %d, body = %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "admin") {
		t.Fatalf("上下文应写入用户名: %s", w.Body.String())
	}
}

// 核心回归：refreshToken 直接调用管理 API 必须被 token_type 检查拦下
func TestJwtVerifyRejectsRefreshToken(t *testing.T) {
	if err := jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456")); err != nil {
		t.Fatalf("配置密钥失败: %v", err)
	}
	refresh, err := jwttoken.RefreshToken("admin")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}

	w := doJwtProbe(newJwtTestRouter(), "Bearer "+refresh)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("refresh token 当 access 用应 401，状态码 = %d, body = %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Token 类型错误") {
		t.Fatalf("应提示类型错误: %s", w.Body.String())
	}
}

func TestJwtVerifyRejectsMissingOrMalformedHeader(t *testing.T) {
	if err := jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456")); err != nil {
		t.Fatalf("配置密钥失败: %v", err)
	}
	r := newJwtTestRouter()

	cases := []struct {
		name string
		auth string
		want string
	}{
		{"missing", "", "未提供认证信息"},
		{"not_bearer", "Basic dXNlcjpwYXNz", "Token格式错误"},
		{"empty_token", "Bearer ", "Token不能为空"},
		{"garbage", "Bearer not-a-token", ""},
	}
	for _, c := range cases {
		w := doJwtProbe(r, c.auth)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s: 应 401，状态码 = %d", c.name, w.Code)
		}
		if c.want != "" && !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("%s: 应包含 %q: %s", c.name, c.want, w.Body.String())
		}
	}
}

// 登出/轮换后的令牌（token 级黑名单）必须被拒
func TestJwtVerifyRejectsBlacklistedToken(t *testing.T) {
	if err := jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456")); err != nil {
		t.Fatalf("配置密钥失败: %v", err)
	}
	token, err := jwttoken.AccessToken("logout-user")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if err := jwttoken.InvalidateToken(token); err != nil {
		t.Fatalf("注销失败: %v", err)
	}

	w := doJwtProbe(newJwtTestRouter(), "Bearer "+token)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("已注销令牌应 401，状态码 = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Token已失效") {
		t.Fatalf("应提示令牌已失效: %s", w.Body.String())
	}
}

// 改密后的用户级撤销：撤销前签发的 access token 必须被拒
func TestJwtVerifyRejectsUserRevokedToken(t *testing.T) {
	if err := jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456")); err != nil {
		t.Fatalf("配置密钥失败: %v", err)
	}
	token, err := jwttoken.AccessToken("pw-change-user")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	// 撤销时间取签发时间 +1，保证严格晚于签发
	jwttoken.InvalidateUserTokens("pw-change-user", tokenIssuedAt(t, token)+1)

	w := doJwtProbe(newJwtTestRouter(), "Bearer "+token)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("用户级撤销后的令牌应 401，状态码 = %d", w.Code)
	}
}

// 白名单匹配：登录/健康检查/静态资源放行，管理 API 与穿越路径不白名单
func TestIsPathWhitelisted(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/api/login", true},
		{"/api/refreshToken", true},
		{"/health", true},
		{"/metrics", true},
		{"/console/index.html", true},
		{"/static/js/app.js", true},
		// 上传接口曾因前缀匹配被匿名放行（3.6.1 高危），必须不在白名单
		{"/api/v1/upload/images", false},
		{"/api/v1/post/list", false},
		// 路径穿越规范化后不得命中白名单
		{"/api/login/../v1/upload/images", false},
		{"/console/../api/v1/post/list", false},
	}
	for _, c := range cases {
		if got := isPathWhitelisted(c.path); got != c.want {
			t.Errorf("isPathWhitelisted(%q) = %v, 期望 %v", c.path, got, c.want)
		}
	}
}

// tokenIssuedAt 从令牌中取出签发时间（Unix 秒）
func tokenIssuedAt(t *testing.T, token string) int64 {
	t.Helper()
	claims, err := jwttoken.ParseToken(token)
	if err != nil {
		t.Fatalf("解析令牌失败: %v", err)
	}
	if claims.IssuedAt == nil {
		t.Fatal("令牌缺少 iat 声明")
	}
	return claims.IssuedAt.Unix()
}
