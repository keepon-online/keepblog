package about

import (
	"net/http"

	inpkg "gitee.com/jieepre/keepblog/internal/pkg"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/md"
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
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "about.html", gin.H{
		"about":           about,
		"note":            md.ToHTML([]byte(about.Note)),
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "关于",
		"canonical":       inpkg.CanonicalURL(site.URL, "/about"),
	})
}
