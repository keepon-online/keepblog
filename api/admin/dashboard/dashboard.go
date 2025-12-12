package dashboard

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) DashboardData(c *gin.Context) {
	group := h.Service.Dashboard.DashboardData()
	result.Ok(c, group)
}
