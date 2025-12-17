package archive

import (
	"html/template"
	"net/http"
	"strconv"

	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/page"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h Handler) Archives(c *gin.Context) {
	sidebarInfo := h.Service.SidebarService.Sidebar()
	pageStr := c.Param("page")
	pageNum, _ := strconv.ParseInt(pageStr, 10, 64)
	if pageNum == 0 {
		pageNum = 1
	}
	archivePosts, err := h.Service.PostService.GetArchivePosts("", "")

	// 计算归档中的文章总数用于分页
	totalPosts := 0
	if archivePosts != nil && archivePosts.Archives != nil {
		for _, posts := range archivePosts.Archives {
			totalPosts += len(posts)
		}
	}

	//传到模板中需要转换成template.HTML类型，否则html代码会被转义
	index, err := page.HandleIndex(totalPosts, int(pageNum), "archives/page")
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	pagination := template.HTML(index)
	c.HTML(http.StatusOK, "archives.html", gin.H{
		"archives":        archivePosts,
		"site":            site,
		"pages":           pagination,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "归档",
	})
}

func (h Handler) ArchivesInfo(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	sidebarInfo := h.Service.SidebarService.Sidebar()
	archivePosts, _ := h.Service.PostService.GetArchivePosts(year, month)
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "archives.html", gin.H{
		"archives":        archivePosts,
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "归档",
	})
}
