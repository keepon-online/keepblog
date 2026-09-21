package post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/internal/service"
	"gitee.com/jieepre/keepblog/internal/testutil"
	"github.com/gin-gonic/gin"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *Handler, func()) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewTestDB(t)
	oldGorm := global.GORM
	global.GORM = db

	ctx := &core.Context{
		Service: service.InitAppService(db),
	}
	handler := &Handler{Context: ctx}

	r := gin.New()
	return r, handler, func() {
		global.GORM = oldGorm
	}
}

func TestSavePost_PreservesAuthorAndCover(t *testing.T) {
	r, handler, cleanup := setupTestRouter(t)
	defer cleanup()

	r.POST("/save", func(c *gin.Context) {
		c.Set("username", "admin_user")
		handler.SavePost(c)
	})

	// Case 1: Custom author and custom cover image provided
	body := map[string]any{
		"title":       "测试文章1",
		"author":      "custom_author",
		"coverImage":  "https://example.com/cover1.jpg",
		"postContent": "正文内容1",
		"tags":        []string{"Go", "Vue"},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP Status = %d, body = %s", w.Code, w.Body.String())
	}

	var saved model.Post
	if err := global.GORM.Where("title = ?", "测试文章1").First(&saved).Error; err != nil {
		t.Fatalf("查询数据库失败: %v", err)
	}

	if saved.Author != "custom_author" {
		t.Errorf("Author = %s, want custom_author", saved.Author)
	}
	if saved.CoverImage != "https://example.com/cover1.jpg" {
		t.Errorf("CoverImage = %s, want https://example.com/cover1.jpg", saved.CoverImage)
	}

	// Case 2: Empty author and empty cover image (for starry night mode)
	body2 := map[string]any{
		"title":       "测试文章2",
		"author":      "",
		"coverImage":  "",
		"postContent": "正文内容2",
		"tags":        []string{"Gin"},
	}
	data2, _ := json.Marshal(body2)
	req2 := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(data2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("HTTP Status = %d, body = %s", w2.Code, w2.Body.String())
	}

	var saved2 model.Post
	if err := global.GORM.Where("title = ?", "测试文章2").First(&saved2).Error; err != nil {
		t.Fatalf("查询数据库失败: %v", err)
	}

	// Should fallback to username from gin context ("admin_user")
	if saved2.Author != "admin_user" {
		t.Errorf("Author = %s, want admin_user", saved2.Author)
	}
	// CoverImage should remain empty ("") for starry night fallback!
	if saved2.CoverImage != "" {
		t.Errorf("CoverImage = %s, want empty string for starry night mode", saved2.CoverImage)
	}
}

func TestGetRandomCover(t *testing.T) {
	r, handler, cleanup := setupTestRouter(t)
	defer cleanup()

	r.GET("/random-cover", handler.GetRandomCover)

	req := httptest.NewRequest(http.MethodGet, "/random-cover", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP Status = %d, body = %s", w.Code, w.Body.String())
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Payload string `json:"payload"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("反序列化响应失败: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("resp.Code = %d, want 200", resp.Code)
	}
	if resp.Payload == "" {
		t.Error("resp.Payload 不应为空，应有随机图片 URL 或兜底壁纸")
	}
}
