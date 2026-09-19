package result

import (
	"net/http"

	"gitee.com/jieepre/keepblog/internal/errors"
	"github.com/gin-gonic/gin"
)

// Result 成功响应信封。HTTP 状态恒为 200，业务码在 body 的 code 字段。
// 错误响应走 HTTP 状态码 + ErrorResponse（见 internal/middleware/error.go）。
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"payload"`
}

// Ok 成功响应（HTTP 200）
func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &Result{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 简易错误响应，源码级兼容已有 200+ 处调用点。
// 原实现返回 HTTP 200 + body code=400；现改为真实 HTTP 400 状态码
// + ErrorResponse 结构，由 ErrorHandler 中间件统一渲染。
// 调用方若需要特定错误码，应改用 Fail（见 internal/errors）。
func Error(c *gin.Context, message string) {
	_ = c.Error(errors.BadRequest(message))
	// 只记录错误并中断,响应由 ErrorHandler 中间件统一渲染;
	// 此处 AbortWithStatus 会提前写响应头,导致 ErrorHandler 跳过渲染输出空 body
	c.Abort()
}

// Fail 返回带特定错误码的错误响应。语义敏感的调用点应使用此函数。
func Fail(c *gin.Context, appErr *errors.AppError) {
	_ = c.Error(appErr)
	c.Abort()
}

// With 返回自定义状态码的响应。迁移期间保留，逐步替换为 Fail 或直接 Ok。
// Deprecated: 请使用 Ok 或 Fail
func With(c *gin.Context, code int, message string, data any) {
	c.JSON(http.StatusOK, &Result{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
