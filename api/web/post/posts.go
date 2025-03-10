package post

import (
	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/hash"
	"gitee.com/jieepre/go-site/pkg/md"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Post(c *gin.Context) {
	id := c.Param("hashids")
	postIds, err := hash.New().HashidsDecode(id)
	if err != nil {
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	posts, err := h.Service.PostService.GetPost(postIds[0])
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	sidebarInfo := h.Service.SidebarService.Sidebar()
	c.HTML(http.StatusOK, "post.html", gin.H{
		"posts":           posts,
		"content":         string(md.Goldmark2html([]byte(posts.PostContent))),
		"toc":             md.Goldmark2htmlToc([]byte(posts.PostContent)),
		"stats":           md.Goldmarkstats([]byte(posts.PostContent)),
		"gitalk":          config.Get().Gitalk,
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
	})
}
