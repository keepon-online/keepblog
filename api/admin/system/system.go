package system

import (
	"gitee.com/jieepre/go-site/internal/model/request"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h Handler) LoginLog(c *gin.Context) {
	query := request.LoginLogQuery{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	list, err := h.Service.SystemService.LoginLogList(query)

	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, list)
}
func (h Handler) AccessLog(c *gin.Context) {
	query := request.AccessLogQuery{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	list, err := h.Service.SystemService.AccessLogList(query)

	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, list)
}
