package system

import (
	"net/http/httptest"
	"testing"

	"gitee.com/jieepre/keepblog/internal/model/request"
	"gitee.com/jieepre/keepblog/internal/testutil"
	"gitee.com/jieepre/keepblog/pkg/jwttoken"

	"github.com/gin-gonic/gin"
)

// loginCtx 构造 Login 所需的最小 gin.Context
func loginCtx(t *testing.T) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest("POST", "/api/login", nil)
	c.Request = req
	return c
}

func TestLogin_Success(t *testing.T) {
	testutil.NewTestDB(t)
	if err := jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456")); err != nil {
		t.Fatalf("Configure 报错: %v", err)
	}
	s := NewSystemService(testutil.NewTestDB(t))

	resp, err := s.Login(request.LoginRequest{
		Username: "admin",
		Password: testutil.TestUserPassword,
	}, loginCtx(t))
	if err != nil {
		t.Fatalf("Login 报错: %v", err)
	}
	if resp.Username != "admin" {
		t.Errorf("Username = %q", resp.Username)
	}
	if len(resp.Roles) != 1 || resp.Roles[0] != "admin" {
		t.Errorf("Roles = %v, want [admin]", resp.Roles)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("应签发 accessToken 与 refreshToken")
	}

	// 签发的 token 可解析且属于该用户
	claims, err := jwttoken.ParseToken(resp.AccessToken)
	if err != nil {
		t.Fatalf("ParseToken 报错: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("token claims.Username = %q", claims.Username)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	testutil.NewTestDB(t)
	_ = jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456"))
	s := NewSystemService(testutil.NewTestDB(t))

	_, err := s.Login(request.LoginRequest{
		Username: "admin",
		Password: "wrong-password",
	}, loginCtx(t))
	if err == nil {
		t.Fatal("错误密码应登录失败")
	}
	if err.Error() != "用户名或密码错误" {
		t.Errorf("错误信息 = %q", err.Error())
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	testutil.NewTestDB(t)
	_ = jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456"))
	s := NewSystemService(testutil.NewTestDB(t))

	_, err := s.Login(request.LoginRequest{
		Username: "nobody",
		Password: "whatever",
	}, loginCtx(t))
	if err == nil {
		t.Fatal("不存在的用户应登录失败")
	}
}

func TestRefreshToken_Flow(t *testing.T) {
	testutil.NewTestDB(t)
	_ = jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456"))
	s := NewSystemService(testutil.NewTestDB(t))

	loginResp, err := s.Login(request.LoginRequest{
		Username: "admin",
		Password: testutil.TestUserPassword,
	}, loginCtx(t))
	if err != nil {
		t.Fatalf("Login 报错: %v", err)
	}

	refreshed, err := s.RefreshToken(loginResp.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken 报错: %v", err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Error("刷新后应有新的 accessToken/refreshToken")
	}

	// 特征（怪癖）：claims 只含秒级时间戳，同一秒内签发的 refreshToken
	// 与旧串完全相同——"令牌轮换"在同一秒内并不产生新令牌。
	// 若后续为 claims 增加唯一性（jti）可把此断言反转为 !=。
	if refreshed.RefreshToken != loginResp.RefreshToken {
		t.Log("refreshToken 发生了真实轮换（同一秒内通常相同）")
	}

	// 无效 token 拒绝
	if _, err := s.RefreshToken("invalid-token"); err == nil {
		t.Error("无效 refreshToken 应报错")
	}
}

func TestChangePassword_InvalidatesOldTokens(t *testing.T) {
	testutil.NewTestDB(t)
	_ = jwttoken.Configure([]byte("test-secret-0123456789abcdef0123456"))
	s := NewSystemService(testutil.NewTestDB(t))

	loginResp, err := s.Login(request.LoginRequest{
		Username: "admin",
		Password: testutil.TestUserPassword,
	}, loginCtx(t))
	if err != nil {
		t.Fatalf("Login 报错: %v", err)
	}

	if err := s.ChangePassword(request.ChangePasswordRequest{
		Username:    "admin",
		OldPassword: testutil.TestUserPassword,
		NewPassword: "new-password-456",
	}); err != nil {
		t.Fatalf("ChangePassword 报错: %v", err)
	}

	// 旧密码失效、新密码可用
	if _, err := s.Login(request.LoginRequest{Username: "admin", Password: testutil.TestUserPassword}, loginCtx(t)); err == nil {
		t.Error("旧密码应已失效")
	}
	if _, err := s.Login(request.LoginRequest{Username: "admin", Password: "new-password-456"}, loginCtx(t)); err != nil {
		t.Errorf("新密码登录失败: %v", err)
	}

	// 修改密码前签发的 refreshToken 因用户级黑名单被拒
	if _, err := s.RefreshToken(loginResp.RefreshToken); err == nil {
		t.Error("改密前的 refreshToken 应已失效")
	}

	// 旧密码再改密应报错
	if err := s.ChangePassword(request.ChangePasswordRequest{
		Username: "admin", OldPassword: testutil.TestUserPassword, NewPassword: "x",
	}); err == nil {
		t.Error("旧密码错误时 ChangePassword 应报错")
	}
}
