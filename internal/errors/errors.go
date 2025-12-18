package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// ErrorCode 错误码定义
type ErrorCode int

const (
	// 系统错误 (1000-1999)
	ErrInternal   ErrorCode = 1000
	ErrDatabase   ErrorCode = 1001
	ErrRedis      ErrorCode = 1002
	ErrConfig     ErrorCode = 1003
	ErrFileSystem ErrorCode = 1004
	ErrNetwork    ErrorCode = 1005

	// 认证授权错误 (2000-2999)
	ErrUnauthorized     ErrorCode = 2000
	ErrForbidden        ErrorCode = 2001
	ErrTokenExpired     ErrorCode = 2002
	ErrTokenInvalid     ErrorCode = 2003
	ErrLoginFailed      ErrorCode = 2004
	ErrPermissionDenied ErrorCode = 2005

	// 参数验证错误 (3000-3999)
	ErrInvalidParam   ErrorCode = 3000
	ErrMissingParam   ErrorCode = 3001
	ErrInvalidFormat  ErrorCode = 3002
	ErrInvalidRange   ErrorCode = 3003
	ErrDuplicateValue ErrorCode = 3004

	// 业务逻辑错误 (4000-4999)
	ErrResourceNotFound    ErrorCode = 4000
	ErrResourceExists      ErrorCode = 4001
	ErrResourceLocked      ErrorCode = 4002
	ErrOperationNotAllowed ErrorCode = 4003
	ErrQuotaExceeded       ErrorCode = 4004

	// 外部服务错误 (5000-5999)
	ErrThirdPartyService ErrorCode = 5000
	ErrPaymentService    ErrorCode = 5001
	ErrEmailService      ErrorCode = 5002
	ErrSMSService        ErrorCode = 5003
)

// AppError 应用错误结构
type AppError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	// 内部字段
	httpStatus int
	cause      error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.cause
}

// HTTPStatus 获取HTTP状态码
func (e *AppError) HTTPStatus() int {
	if e.httpStatus != 0 {
		return e.httpStatus
	}

	// 根据错误码映射HTTP状态码
	switch {
	case e.Code >= 2000 && e.Code < 3000:
		return http.StatusUnauthorized
	case e.Code >= 3000 && e.Code < 4000:
		return http.StatusBadRequest
	case e.Code >= 4000 && e.Code < 5000:
		return http.StatusNotFound
	case e.Code >= 5000:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithRequestID 添加请求ID
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithMetadata 添加元数据
func (e *AppError) WithMetadata(key string, value interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithStack 添加堆栈信息
func (e *AppError) WithStack() *AppError {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	e.Stack = string(buf[:n])
	return e
}

// WithHTTPStatus 设置HTTP状态码
func (e *AppError) WithHTTPStatus(status int) *AppError {
	e.httpStatus = status
	return e
}

// 错误构造函数

// New 创建新错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// Wrap 包装现有错误
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		cause:     err,
	}
}

// Internal 内部服务器错误
func Internal(message string) *AppError {
	return New(ErrInternal, message)
}

// InvalidParam 参数错误
func InvalidParam(message string) *AppError {
	return New(ErrInvalidParam, message)
}

// Unauthorized 未授权错误
func Unauthorized(message string) *AppError {
	return New(ErrUnauthorized, message)
}

// NotFound 资源未找到错误
func NotFound(message string) *AppError {
	return New(ErrResourceNotFound, message)
}

// BadRequest 请求错误
func BadRequest(message string) *AppError {
	return New(ErrInvalidParam, message).WithHTTPStatus(http.StatusBadRequest)
}

// 错误码到消息的映射
var errorMessages = map[ErrorCode]string{
	ErrInternal:            "内部服务器错误",
	ErrDatabase:            "数据库连接错误",
	ErrRedis:               "缓存服务错误",
	ErrConfig:              "配置错误",
	ErrFileSystem:          "文件系统错误",
	ErrNetwork:             "网络连接错误",
	ErrUnauthorized:        "未授权访问",
	ErrForbidden:           "禁止访问",
	ErrTokenExpired:        "令牌已过期",
	ErrTokenInvalid:        "无效令牌",
	ErrLoginFailed:         "登录失败",
	ErrPermissionDenied:    "权限不足",
	ErrInvalidParam:        "参数错误",
	ErrMissingParam:        "缺少必要参数",
	ErrInvalidFormat:       "格式错误",
	ErrInvalidRange:        "数值超出范围",
	ErrDuplicateValue:      "数据重复",
	ErrResourceNotFound:    "资源不存在",
	ErrResourceExists:      "资源已存在",
	ErrResourceLocked:      "资源被锁定",
	ErrOperationNotAllowed: "操作不被允许",
	ErrQuotaExceeded:       "配额已超限",
	ErrThirdPartyService:   "第三方服务错误",
	ErrPaymentService:      "支付服务错误",
	ErrEmailService:        "邮件服务错误",
	ErrSMSService:          "短信服务错误",
}

// GetMessage 获取错误码对应的默认消息
func GetMessage(code ErrorCode) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "未知错误"
}

// NewWithCode 使用错误码创建错误（使用默认消息）
func NewWithCode(code ErrorCode) *AppError {
	return New(code, GetMessage(code))
}
