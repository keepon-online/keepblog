package dashboard

import (
	"strconv"

	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) DashboardData(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val > 0 {
			days = val
		}
	}
	group := h.Service.Dashboard.DashboardData(days)
	result.Ok(c, group)
}

