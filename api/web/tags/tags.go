package tags

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/page"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Tags(c *gin.Context) {
	sidebarInfo := h.Service.SidebarService.Sidebar()
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", nil)

		return
	}
	c.HTML(200, "tags.html", gin.H{
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"title":           "标签",
	})
}

func (h *Handler) TagsPage(c *gin.Context) {
	tagName := c.Param("tag")
	pageNum, _ := strconv.ParseInt(c.Param("page"), 10, 64)

	tagsPosts, total, _ := h.Service.PostService.GetPostsByTag(tagName, int(pageNum))
	//传到模板中需要转换成template.HTML类型，否则html代码会被转义
	index, err := page.HandleIndex(int(total), int(pageNum), "tags/"+tagName+"/page")
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	pagination := template.HTML(index)
	sidebarInfo := h.Service.SidebarService.Sidebar()
	c.HTML(http.StatusOK, "tags-info.html", gin.H{
		"site":            site,
		"tagsPosts":       tagsPosts,
		"tagName":         tagName,
		"tags":            sidebarInfo.Tag,
		"pages":           pagination,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"title":           "标签",
	})
}
