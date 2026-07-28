package link

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Links(c *gin.Context) {
	sidebarInfo := h.Service.SidebarService.Sidebar()
	links, _ := h.Service.LinkService.GetLinks()
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "link.html", gin.H{
		"tags":            sidebarInfo.Tag,
		"links":           links,
		"site":            site,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "友链",
	})
}
