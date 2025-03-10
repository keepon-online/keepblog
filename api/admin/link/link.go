package link

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

func (h *Handler) SaveLink(c *gin.Context) {
	category := model.FriendLink{}
	err := c.ShouldBindJSON(&category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.LinkService.Save(category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdateLink(c *gin.Context) {
	link := model.FriendLink{}
	err := c.ShouldBindJSON(&link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.LinkService.Update(link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdateLinkState(c *gin.Context) {
	link := model.FriendLink{}
	err := c.ShouldBindJSON(&link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.LinkService.Update(link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) DeleteLink(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.ParseInt(param, 0, 64)

	err := h.Service.LinkService.Delete(int(id))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) DetailLink(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.ParseInt(param, 0, 64)
	linkInfo, err := h.Service.LinkService.GetLink(uint32(id))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, linkInfo)
}
func (h *Handler) GetLinkList(c *gin.Context) {
	links, err := h.Service.LinkService.GetLinkList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, links)
}
