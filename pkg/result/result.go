package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"payload"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &Result{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, &Result{
		Code:    400,
		Message: message,
		Data:    nil,
	})
}

func With(c *gin.Context, code int, message string, data any) {
	c.JSON(http.StatusOK, &Result{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
