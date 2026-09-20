package post

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/pkg/md"
	"github.com/gin-gonic/gin"
)

// rssFeed 是 RSS 2.0 的根元素，用于 xml.Marshal 输出。
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

// sitemapURLSet 是 sitemap 0.9 的根元素。
type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod,omitempty"`
}

// feedLimit RSS 输出的最大文章条数。
const feedLimit = 50

// sitemapLimit sitemap 输出的最大文章条数。sitemap 协议单文件上限 5 万 URL，
// 远大于博客文章规模，相当于全量收录。
const sitemapLimit = 50000

// Feed 输出 RSS 2.0 订阅源（/rss.xml）。
func (h *Handler) Feed(c *gin.Context) {
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.String(http.StatusOK, "")
		return
	}
	baseURL := strings.TrimRight(site.URL, "/")

	posts, _ := h.Service.PostService.GetPublishedPostsForFeed(feedLimit)

	items := make([]rssItem, 0, len(posts))
	for _, p := range posts {
		link := baseURL + "/post/" + p.PostSlug
		// 摘要为空时以正文纯文本兜底，避免订阅器里出现空描述。
		description := p.Summary
		if description == "" {
			description = md.Excerpt([]byte(p.PostContent), 120)
		}
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        link,
			Description: description,
			PubDate:     time.Unix(int64(p.PubTime), 0).Format(time.RFC1123Z),
			GUID:        link,
		})
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         site.Title,
			Link:          baseURL,
			Description:   site.Description,
			Language:      "zh-cn",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		c.String(http.StatusOK, "")
		return
	}
	// xml.MarshalIndent 不会带 XML 声明，补上（RSS 订阅器要求）
	output := append([]byte(xml.Header), data...)
	c.Data(http.StatusOK, "application/xml; charset=utf-8", output)
}

// Sitemap 输出 sitemap 0.9（/sitemap.xml）。包含首页、全部已发布文章
// （不设 50 条上限）、标签/分类/归档/关于/友链等聚合页。
func (h *Handler) Sitemap(c *gin.Context) {
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.String(http.StatusOK, "")
		return
	}
	baseURL := strings.TrimRight(site.URL, "/")

	posts, _ := h.Service.PostService.GetPublishedPostsForFeed(sitemapLimit)

	urls := make([]sitemapURL, 0, len(posts)+16)
	// 首页。lastmod 取最新一篇文章的时间，而非生成时刻——每次生成都写"今天"
	// 会让搜索引擎无法分辨站点是否真的有更新。
	urls = append(urls, sitemapURL{
		Loc:     baseURL + "/",
		Lastmod: latestPostDate(posts),
	})
	// 文章页
	for _, p := range posts {
		u := sitemapURL{Loc: baseURL + "/post/" + p.PostSlug}
		if p.LastModifiedTime > 0 {
			u.Lastmod = time.Unix(int64(p.LastModifiedTime), 0).Format("2006-01-02")
		} else if p.PubTime > 0 {
			u.Lastmod = time.Unix(int64(p.PubTime), 0).Format("2006-01-02")
		}
		urls = append(urls, u)
	}
	// 聚合页（无 lastmod：内容随文章变化，交给搜索引擎按抓取频率发现）
	for _, path := range []string{"/about", "/link", "/tags", "/categories", "/archives"} {
		urls = append(urls, sitemapURL{Loc: baseURL + path})
	}
	// 标签详情页
	if tags, err := h.Service.TagService.GetTags(); err == nil {
		for _, t := range tags {
			urls = append(urls, sitemapURL{Loc: baseURL + "/tags/" + url.PathEscape(t.TagName)})
		}
	}
	// 分类详情页
	if categories, err := h.Service.CategoryService.GetCategories(); err == nil {
		for _, ct := range categories {
			urls = append(urls, sitemapURL{Loc: baseURL + "/categories/" + url.PathEscape(ct.CategoryName)})
		}
	}

	urlSet := sitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	data, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		c.String(http.StatusOK, "")
		return
	}
	output := append([]byte(xml.Header), data...)
	c.Data(http.StatusOK, "application/xml; charset=utf-8", output)
}

// latestPostDate 取文章列表（按发布时间倒序）中最新一篇的日期；
// 列表为空时回退到今天。
func latestPostDate(posts []model.Post) string {
	if len(posts) > 0 {
		if p := posts[0]; p.LastModifiedTime > 0 {
			return time.Unix(int64(p.LastModifiedTime), 0).Format("2006-01-02")
		} else if p.PubTime > 0 {
			return time.Unix(int64(p.PubTime), 0).Format("2006-01-02")
		}
	}
	return time.Now().Format("2006-01-02")
}

// Robots 输出 robots.txt（/robots.txt），动态注入站点 sitemap 地址，
// 并屏蔽对搜索无价值、会浪费抓取预算的路径（API/后台/搜索/日报）。
// 原为 static/robots.txt 编译期 embed 文件，无法插入站点 URL 变量，故改为动态渲染。
func (h *Handler) Robots(c *gin.Context) {
	site, err := h.Service.WebSiteService.GetWebSite()
	baseURL := ""
	if err == nil {
		baseURL = strings.TrimRight(site.URL, "/")
	}

	var b strings.Builder
	b.WriteString("User-agent: *\n")
	for _, path := range []string{"/api/", "/console", "/search/", "/daily"} {
		b.WriteString("Disallow: " + path + "\n")
	}
	b.WriteString("\n")
	if baseURL != "" {
		b.WriteString("Sitemap: " + baseURL + "/sitemap.xml\n")
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(b.String()))
}
