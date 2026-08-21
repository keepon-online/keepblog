package category

import (
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/page"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Categories(c *gin.Context) {
	sidebarInfo := h.Service.SidebarService.Sidebar()
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {

		c.HTML(200, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "categories.html", gin.H{
		"site":            site,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "分类",
	})
}

func (h *Handler) CategoriesPage(c *gin.Context) {
	categoryName := c.Param("category")
	pageStr := c.Param("page")
	pageNum, _ := strconv.ParseInt(pageStr, 10, 64)
	if pageNum == 0 {
		pageNum = 1
	}
	categoriesPosts, total, _ := h.Service.PostService.GetPostsByCategory(categoryName, int(pageNum))
	sidebarInfo := h.Service.SidebarService.Sidebar()
	//传到模板中需要转换成template.HTML类型，否则html代码会被转义
	index, _ := page.HandleIndex(int(total), int(pageNum), "categories/"+categoryName+"/page")
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	pagination := template.HTML(index)
	c.HTML(http.StatusOK, "categories-info.html", gin.H{
		"site":            site,
		"categoriesPosts": categoriesPosts,
		"categoryName":    categoryName,
		"tags":            sidebarInfo.Tag,
		"pages":           pagination,
		"categories":      sidebarInfo.Category,
		"latestPosts":     sidebarInfo.LatestPosts,
		"cardInfo":        sidebarInfo.CardInfo,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"title":           "分类",
	})
}
