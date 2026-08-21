package post

import (
	"encoding/json"
	"net/http"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/hash"
	"gitee.com/jieepre/keepblog/pkg/md"
	"github.com/gin-gonic/gin"
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
	posts, _ := h.Service.PostService.GetPost(postIds[0])
	if posts == nil {
		c.HTML(http.StatusOK, "500.html", nil)
		return
	}
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	sidebarInfo := h.Service.SidebarService.Sidebar()
	rendered := md.Render([]byte(posts.PostContent))
	// Gitalk admin 是 []string，模板直接渲染会得到字面量字符串，无法作为 JS 数组。
	// 这里预序列化为 JSON 字符串，模板里用 unescaped 输出成合法的 JS 数组字面量。
	gitalkAdminJSON, _ := json.Marshal(config.Get().Gitalk.Admin)
	// 上一篇/下一篇/相关文章：查询失败不中断渲染，传 nil 时模板用 {{if}} 跳过。
	prevPost, nextPost, _ := h.Service.PostService.GetAdjacentPosts(posts.PostId, posts.PubTime)
	relatedPosts, _ := h.Service.PostService.GetRelatedPosts(posts.PostId, posts.Tags, 6)
	c.HTML(http.StatusOK, "post.html", gin.H{
		"posts":           posts,
		"content":         rendered.HTML,
		"toc":             rendered.TOC,
		"stats":           rendered.Stats,
		"gitalk":          config.Get().Gitalk,
		"gitalkAdmin":     string(gitalkAdminJSON),
		"site":            site,
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
