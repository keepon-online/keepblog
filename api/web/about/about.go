package about

import (
	"net/http"

	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/md"
	"github.com/gin-gonic/gin"
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

	c.HTML(http.StatusOK, "about.html", gin.H{
		"about":           about,
		"note":            string(md.Goldmark2html([]byte(about.Note))),
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "关于",
	})
}
