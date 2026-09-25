package link

import (
	"net/http"
	"strconv"
	"time"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
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

func (h *Handler) CheckLink(c *gin.Context) {
	targetUrl := c.Query("url")
	if targetUrl == "" {
		result.Error(c, "url 不能为空")
		return
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		result.Ok(c, gin.H{"status": 0, "msg": "无效链接"})
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		result.Ok(c, gin.H{"status": 0, "msg": "访问超时或无法连接"})
		return
	}
	defer resp.Body.Close()
	result.Ok(c, gin.H{"status": resp.StatusCode, "msg": resp.Status})
}

