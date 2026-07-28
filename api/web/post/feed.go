package post

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// rssFeed 是 RSS 2.0 的根元素，用于 xml.Marshal 输出。
type rssFeed struct {
	XMLName xml.Name    `xml:"rss"`
	Version string      `xml:"version,attr"`
	Channel rssChannel  `xml:"channel"`
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

// feedLimit RSS / sitemap 输出的最大文章条数。
const feedLimit = 50

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
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        link,
			Description: p.Summary,
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

// Sitemap 输出 sitemap 0.9（/sitemap.xml），包含首页与全部已发布文章。
func (h *Handler) Sitemap(c *gin.Context) {
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.String(http.StatusOK, "")
		return
	}
	baseURL := strings.TrimRight(site.URL, "/")

	posts, _ := h.Service.PostService.GetPublishedPostsForFeed(feedLimit)

	urls := make([]sitemapURL, 0, len(posts)+1)
	// 首页
	urls = append(urls, sitemapURL{
		Loc:     baseURL + "/",
		Lastmod: time.Now().Format("2006-01-02"),
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

// Robots 输出 robots.txt（/robots.txt），动态注入站点 sitemap 地址。
// 原为 static/robots.txt 编译期 embed 文件，无法插入站点 URL 变量，故改为动态渲染。
func (h *Handler) Robots(c *gin.Context) {
	site, err := h.Service.WebSiteService.GetWebSite()
	baseURL := ""
	if err == nil {
		baseURL = strings.TrimRight(site.URL, "/")
	}

	var b strings.Builder
	b.WriteString("User-agent: *\n")
	b.WriteString("Allow:\n")
	b.WriteString("\n")
	if baseURL != "" {
		b.WriteString("Sitemap: " + baseURL + "/sitemap.xml\n")
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(b.String()))
}
