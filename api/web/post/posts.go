package post

import (
	"net/http"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/model/system"
	inpkg "gitee.com/jieepre/keepblog/internal/pkg"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/hash"
	"gitee.com/jieepre/keepblog/pkg/md"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

// renderNotFound 输出真正的 404（此前用 200 渲染 404 页构成软 404，会被
// 搜索引擎当作有效页面收录）。站点信息尽力获取，仅用于补齐页面标题。
func (h *Handler) renderNotFound(c *gin.Context) {
	site, _ := h.Service.WebSiteService.GetWebSite()
	if site == nil {
		site = &system.WebSite{}
	}
	c.HTML(http.StatusNotFound, "404.html", gin.H{"site": site})
}

func (h *Handler) Post(c *gin.Context) {
	id := c.Param("hashids")
	postIds, err := hash.New().HashidsDecode(id)
	if err != nil || len(postIds) == 0 {
		h.renderNotFound(c)
		return
	}
	posts, _ := h.Service.PostService.GetPost(postIds[0])
	if posts == nil {
		h.renderNotFound(c)
		return
	}
	// 摘要为空时以正文纯文本兜底：meta description / og / JSON-LD / RSS
	// 都依赖它，空摘要等于把搜索摘要完全交给搜索引擎自行截取。
	if posts.Summary == "" {
		posts.Summary = md.Excerpt([]byte(posts.PostContent), 120)
	}
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	sidebarInfo := h.Service.SidebarService.Sidebar()
	rendered := md.Render([]byte(posts.PostContent))
	// 上一篇/下一篇/相关文章：查询失败不中断渲染，传 nil 时模板用 {{if}} 跳过。
	prevPost, nextPost, _ := h.Service.PostService.GetAdjacentPosts(posts.PostId, posts.PubTime)
	relatedPosts, _ := h.Service.PostService.GetRelatedPosts(posts.PostId, posts.Tags, 6)
	c.HTML(http.StatusOK, "post.html", gin.H{
		"posts":           posts,
		"content":         rendered.HTML,
		"toc":             rendered.TOC,
		"stats":           rendered.Stats,
		"artalk":          config.Get().Artalk,
		"site":            site,
		"canonical":       inpkg.CanonicalURL(site.URL, "/post/"+posts.PostSlug),
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"prevPost":        prevPost,
		"nextPost":        nextPost,
		"relatedPosts":    relatedPosts,
	})
}
