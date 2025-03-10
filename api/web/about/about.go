package about

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/md"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	*core.Context
}

func (h Handler) About(c *gin.Context) {
	sidebarInfo := h.Service.SidebarService.Sidebar()
	about, _ := h.Service.AboutService.Info()
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}

	note := md.LuneTet(about.Note)
	c.HTML(http.StatusOK, "about.html", gin.H{
		"about":           about,
		"note":            note,
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"title":           "关于",
	})
}
