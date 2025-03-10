package category

import (
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Handler struct {
	*core.Context
}

func (h *Handler) SaveCategory(c *gin.Context) {
	category := model.Category{}
	err := c.ShouldBindJSON(&category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.CategoryService.Save(category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdateCategory(c *gin.Context) {
	category := model.Category{}
	err := c.ShouldBindJSON(&category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.CategoryService.Update(category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdateCategoryState(c *gin.Context) {
	category := model.Category{}
	err := c.ShouldBindJSON(&category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.CategoryService.Update(category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) DeleteCategory(c *gin.Context) {
	param := c.Param("categoryId")
	categoryId, _ := strconv.ParseInt(param, 0, 64)

	posts, err := h.Service.PostService.GetPostsByCategoryId(categoryId)
	if len(posts) > 0 {
		result.Error(c, "不能删除该栏目")
		return
	}
	err = h.Service.CategoryService.Delete(int(categoryId))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) DetailCategory(c *gin.Context) {
	param := c.Param("categoryId")
	categoryId, _ := strconv.ParseInt(param, 0, 64)
	categoryInfo, err := h.Service.CategoryService.GetCategory(uint32(categoryId))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, categoryInfo)
}
func (h *Handler) GetCategoryList(c *gin.Context) {
	categories, err := h.Service.CategoryService.GetCategoryList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, categories)
}
