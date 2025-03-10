package about

import (
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Save(c *gin.Context) {
	about := model.About{}
	err := c.ShouldBindJSON(&about)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.AboutService.Save(about)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

func (h *Handler) Detail(c *gin.Context) {
	info, err := h.Service.AboutService.Info()
	if err != nil {
		return
	}
	result.Ok(c, info)

}
