package post

import (
	"net/http"
	"regexp"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/model"
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
	// GetPost 对不存在/草稿/未到点定时文章返回零值结构体而非 nil，
	// 只判 nil 会把空白文章页以 200 渲染出去（软 404）。
	if posts == nil || posts.PostId == 0 {
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
	// og:image 兜底：无封面文章分享到社交平台时没有卡片图，取正文第一张
	// 插图顶上；正文也无图则保持为空，模板按 summary 卡片渲染。
	// 只改内存值，封面字段落库仍以用户在后台设置的为准。
	if posts.CoverImage == "" {
		posts.CoverImage = firstImageSrc(rendered.HTML)
	}
	// 上一篇/下一篇/相关文章：查询失败不中断渲染，传 nil 时模板用 {{if}} 跳过。
	prevPost, nextPost, _ := h.Service.PostService.GetAdjacentPosts(posts.PostId, posts.PubTime)
	relatedPosts, _ := h.Service.PostService.GetRelatedPosts(posts.PostId, posts.Tags, 6)
	// 系列导航：同系列文章按发布时间升序列出，查询失败不中断渲染。
	var seriesPosts []model.SeriesPost
	if posts.Series != "" {
		seriesPosts, _ = h.Service.PostService.GetSeriesPosts(posts.Series)
	}
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
		"seriesPosts":     seriesPosts,
	})
}

// imgSrcRegexp 提取渲染后 HTML 中第一张 <img> 的 src。goldmark 渲染的属性
// 顺序固定（src 最前），这里只做兜底取图，容错优先于严格解析。
var imgSrcRegexp = regexp.MustCompile(`<img\s[^>]*src="([^"]+)"`)

// firstImageSrc 返回 HTML 中第一张图片的 src，无图返回空串。
func firstImageSrc(html string) string {
	m := imgSrcRegexp.FindStringSubmatch(html)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}
