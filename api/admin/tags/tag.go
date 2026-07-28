package tags

import (
	"strconv"

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

// GetTags 前台侧边栏 tag cloud 用（去重、带样式）。
func (h *Handler) GetTags(c *gin.Context) {
	tsgs := make([]model.Tag, 0)
	if err := global.GORM.Table(model.TTagTable).Find(&tsgs).Error; err != nil {
		result.Error(c, "查询失败")
		return
	}
	tags := tags_remove.RemoveDuplicateElement(tsgs)
	result.Ok(c, tags)
}

// GetTagList 后台管理用：返回全部标签（不过滤、不去重）。
func (h *Handler) GetTagList(c *gin.Context) {
	tagList, err := h.Service.TagService.GetTagList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, tagList)
}

// SaveTag 新增标签。
func (h *Handler) SaveTag(c *gin.Context) {
	tag := model.Tag{}
	if err := c.ShouldBindJSON(&tag); err != nil {
		result.Error(c, err.Error())
		return
	}
	if _, err := h.Service.TagService.Save(tag); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// UpdateTag 更新标签名。
func (h *Handler) UpdateTag(c *gin.Context) {
	tag := model.Tag{}
	if err := c.ShouldBindJSON(&tag); err != nil {
		result.Error(c, err.Error())
		return
	}
	if err := h.Service.TagService.Update(tag); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// DeleteTag 删除标签（含 post_tag 关联清理）。
func (h *Handler) DeleteTag(c *gin.Context) {
	param := c.Param("tagId")
	tagId, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		result.Error(c, "标签ID无效")
		return
	}
	if err := h.Service.TagService.Delete(uint32(tagId)); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// DetailTag 获取单个标签。
func (h *Handler) DetailTag(c *gin.Context) {
	param := c.Param("tagId")
	tagId, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		result.Error(c, "标签ID无效")
		return
	}
	tagInfo, err := h.Service.TagService.GetTag(uint32(tagId))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, tagInfo)
}
