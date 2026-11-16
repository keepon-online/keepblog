package jwttoken

import (
	"testing"
	"time"
)

// 3.6.1 令牌体系修复的回归测试：access/refresh 靠 token_type 声明区分用途，
// 此前 refreshToken 可直接当 accessToken 调用管理 API。

// 各签发函数必须写入正确的 token_type 声明
func TestTokenTypeClaims(t *testing.T) {
	withTestSecret(t)

	cases := []struct {
		name string
		sign func(string) (string, error)
		want string
	}{
		{"AccessToken", AccessToken, TokenTypeAccess},
		{"RefreshToken", RefreshToken, TokenTypeRefresh},
		{"CreateToken", CreateToken, TokenTypeAccess},
	}
	for _, c := range cases {
		token, err := c.sign("admin")
		if err != nil {
			t.Fatalf("%s: 签发失败: %v", c.name, err)
		}
		claims, err := ParseToken(token)
		if err != nil {
			t.Fatalf("%s: 解析失败: %v", c.name, err)
		}
		if claims.TokenType != c.want {
			t.Errorf("%s: token_type = %q, 期望 %q", c.name, claims.TokenType, c.want)
		}
	}
}

// 空密钥必须 fail-fast，不允许回退到空密钥签发（可被任意伪造）
func TestConfigureRejectsEmptySecret(t *testing.T) {
	if err := Configure(nil); err == nil {
		t.Fatal("空密钥应返回错误")
	}
	if err := Configure([]byte{}); err == nil {
		t.Fatal("空切片密钥应返回错误")
	}
}

// 单 token 黑名单（登出/轮换注销）：InvalidateToken 后同名令牌必须失效。
// 测试进程未启用 Redis，天然覆盖内存降级路径（cache.Enable=false）。
func TestTokenBlacklistInvalidation(t *testing.T) {
	withTestSecret(t)

	token, err := RefreshToken("blacklist-user")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if IsTokenBlacklisted(token) {
		t.Fatal("未注销的令牌不应在黑名单中")
	}

	if err := InvalidateToken(token); err != nil {
		t.Fatalf("注销失败: %v", err)
	}
	if !IsTokenBlacklisted(token) {
		t.Fatal("已注销的令牌必须命中黑名单")
	}

	other, _ := AccessToken("blacklist-user")
	if IsTokenBlacklisted(other) {
		t.Fatal("黑名单按 token 串隔离，不应牵连同用户其它令牌")
	}
}

// 用户级撤销（改密/登出批量失效）：撤销时间点之后签发的令牌仍可用，
// 之前的全部失效。
func TestUserTokenInvalidation(t *testing.T) {
	username := "user-revoke-test"
	if !IsUserTokenValid(username, time.Now().Unix()) {
		t.Fatal("从未撤销的用户，令牌应有效")
	}

	revokeAt := time.Now().Unix()
	InvalidateUserTokens(username, revokeAt)

	if IsUserTokenValid(username, revokeAt-1) {
		t.Fatal("撤销前签发的令牌必须失效")
	}
	if IsUserTokenValid(username, revokeAt) {
		t.Fatal("撤销时刻签发的令牌必须失效（非严格大于撤销时间即无效）")
	}
	if !IsUserTokenValid(username, revokeAt+1) {
		t.Fatal("撤销后重新登录签发的令牌应有效")
	}
}
