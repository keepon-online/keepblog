package common

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitee.com/jieepre/keepblog/internal/middleware"
	"github.com/gin-gonic/gin"
)

// 3.6.1 上传修复的回归测试：该接口曾被 JWT 白名单前缀匿名放行，
// HTML/SVG 可同域托管构成存储型 XSS；内容层必须再做
// 扩展名白名单 + 文件头真实类型校验 + 大小上限三道收敛。

// pngHead 是合法 PNG 文件头，http.DetectContentType 据此识别 image/png
var pngHead = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

func newUploadTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// 复用生产错误渲染（result.Error 只记录错误，由 ErrorHandler 统一输出）
	r.Use(middleware.ErrorHandler())
	r.POST("/api/v1/upload/images", (&Handler{}).UploadImage)
	return r
}

// doUpload 构造携带单个 file 字段的上传请求；filename 为空表示不带文件
func doUpload(t *testing.T, r *gin.Engine, filename string, content []byte) *httptest.ResponseRecorder {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if filename != "" {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			t.Fatalf("构造 multipart 失败: %v", err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatalf("写入文件内容失败: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 writer 失败: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload/images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// HTML/SVG 曾是存储型 XSS 的载荷本体，必须在扩展名层拒绝
func TestUploadImageRejectsDangerousExtensions(t *testing.T) {
	r := newUploadTestRouter()

	for _, filename := range []string{"shell.html", "evil.svg", "cmd.php", "no-ext"} {
		w := doUpload(t, r, filename, pngHead)
		if !strings.Contains(w.Body.String(), "仅支持") {
			t.Errorf("%s: 应被扩展名白名单拒绝: %s", filename, w.Body.String())
		}
	}
}

// 伪装成 png 的 HTML：客户端自报类型不可信，须按文件头识别真实类型
func TestUploadImageRejectsDisguisedContent(t *testing.T) {
	r := newUploadTestRouter()

	html := []byte("<!DOCTYPE html><html><script>alert(1)</script></html>")
	w := doUpload(t, r, "avatar.png", html)
	if !strings.Contains(w.Body.String(), "不是有效的图片") {
		t.Fatalf("伪装图片应被文件头校验拒绝: %s", w.Body.String())
	}
}

func TestUploadImageRejectsOversize(t *testing.T) {
	r := newUploadTestRouter()

	oversize := make([]byte, maxUploadImageSize+1)
	copy(oversize, pngHead)
	w := doUpload(t, r, "big.png", oversize)
	if !strings.Contains(w.Body.String(), "10MB") {
		t.Fatalf("超限文件应被大小校验拒绝: %s", w.Body.String())
	}
}

func TestUploadImageRejectsMissingFile(t *testing.T) {
	r := newUploadTestRouter()

	w := doUpload(t, r, "", nil)
	if !strings.Contains(w.Body.String(), "请选择要上传的图片文件") {
		t.Fatalf("缺文件应返回明确错误: %s", w.Body.String())
	}
}

// 合法 PNG 必须通过全部校验层：到达对象存储层并因未配置而快速失败
// （错误信息是“对象存储未配置”，而非任何校验层消息）
func TestUploadImageValidPNGBypassesValidationLayers(t *testing.T) {
	r := newUploadTestRouter()

	content := append(append([]byte{}, pngHead...), bytes.Repeat([]byte{0x00}, 64)...)
	w := doUpload(t, r, "avatar.png", content)

	body := w.Body.String()
	for _, validationMsg := range []string{"请选择", "仅支持", "不是有效的图片", "10MB"} {
		if strings.Contains(body, validationMsg) {
			t.Fatalf("合法 PNG 不应被 %q 校验层拒绝: %s", validationMsg, body)
		}
	}
	if !strings.Contains(body, "对象存储未配置") {
		t.Fatalf("合法 PNG 应到达对象存储层（未配置快速失败）: %s", body)
	}
}
