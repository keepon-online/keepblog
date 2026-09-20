package tags

import (
	"fmt"
	inpkg "gitee.com/jieepre/keepblog/internal/pkg"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/page"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Tags(c *gin.Context) {
	sidebarInfo := h.Service.SidebarService.Sidebar()
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", nil)

		return
	}
	c.HTML(http.StatusOK, "tags.html", gin.H{
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "标签",
		"canonical":       inpkg.CanonicalURL(site.URL, "/tags"),
	})
}

func (h *Handler) TagsPage(c *gin.Context) {
	tagName := c.Param("tag")
	pageNum, _ := strconv.ParseInt(c.Param("page"), 10, 64)

	tagsPosts, total, _ := h.Service.PostService.GetPostsByTag(tagName, int(pageNum))
	//传到模板中需要转换成template.HTML类型，否则html代码会被转义
	index, _ := page.HandleIndex(int(total), int(pageNum), "tags/"+tagName+"/page")
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	pagination := template.HTML(index)
	sidebarInfo := h.Service.SidebarService.Sidebar()

	// 标签详情页带标签名的标题与描述，翻页页 canonical 指向自身。
	path := "/tags/" + url.PathEscape(tagName)
	if pageNum > 1 {
		path += fmt.Sprintf("/page/%d", pageNum)
	}
	c.HTML(http.StatusOK, "tags-info.html", gin.H{
		"site":            site,
		"tagsPosts":       tagsPosts,
		"tagName":         tagName,
		"tags":            sidebarInfo.Tag,
		"pages":           pagination,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "标签：" + tagName,
		"pagedesc":        fmt.Sprintf("标签「%s」下的全部文章，共 %d 篇。", tagName, total),
		"canonical":       inpkg.CanonicalURL(site.URL, path),
	})
}
