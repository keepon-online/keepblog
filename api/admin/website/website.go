package website

import (
	"gitee.com/jieepre/go-site/internal/model/system"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) SaveWebSite(c *gin.Context) {
	website := system.WebSite{}
	err := c.ShouldBindJSON(&website)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.WebSiteService.Save(website)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

func (h *Handler) GetWebSite(c *gin.Context) {
	info, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		return
	}
	result.Ok(c, info)

}
