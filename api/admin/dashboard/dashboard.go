package dashboard

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) PanelGroup(c *gin.Context) {
	group := h.Service.Dashboard.PanelGroup()
	result.Ok(c, group)
}
func (h *Handler) Pie(c *gin.Context) {
	group := h.Service.Dashboard.Pie()
	result.Ok(c, group)
}
func (h *Handler) Bar(c *gin.Context) {
	group := h.Service.Dashboard.Bar()
	result.Ok(c, group)
}
func (h *Handler) Line(c *gin.Context) {
	group := h.Service.Dashboard.Line()
	result.Ok(c, group)
}
