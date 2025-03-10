package tags

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	tags_remove "gitee.com/jieepre/go-site/pkg/tags-remove"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) GetTags(c *gin.Context) {
	tsgs := make([]model.Tag, 0)
	if err := global.GORM.Table(model.TTagTable).Find(&tsgs).Error; err != nil {
		result.Error(c, "查询失败")
		return
	}
	tags := tags_remove.RemoveDuplicateElement(tsgs)
	result.Ok(c, tags)
}
