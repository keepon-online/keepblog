package home

import (
	"html/template"
	"net/http"
	"strconv"

	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg"
	"gitee.com/jieepre/keepblog/pkg/daily"
	"gitee.com/jieepre/keepblog/pkg/page"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Home(c *gin.Context) {
	pageNum := c.Param("page")
	num, _ := strconv.ParseInt(pageNum, 10, 64)
	coverPosts, _, err := h.Service.PostService.GetCoverPosts(int(num))
	if err != nil {
		slog.Errorf("首页文章错误: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	total, err := h.Service.PostService.Total()
	if err != nil {
		slog.Errorf("首页统计错误: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	sidebarInfo := h.Service.SidebarService.Sidebar()
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		slog.Errorf("首页错误: %s", err.Error())
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	//传到模板中需要转换成template.HTML类型，否则html代码会被转义
	paginationTpl, err := page.HandleIndex(int(total), int(num), "page")
	if err != nil {
		slog.Errorf("首页错误: %s", err.Error())
		c.HTML(http.StatusOK, "error.html", nil)
		return
	}
	pagination := template.HTML(paginationTpl)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"coverPosts":      coverPosts,
		"site":            site,
		"pages":           pagination,
		"tags":            sidebarInfo.Tag,
		"categories":      sidebarInfo.Category,
		"cardInfo":        sidebarInfo.CardInfo,
		"latestPosts":     sidebarInfo.LatestPosts,
		"sidebarArchives": sidebarInfo.SidebarArchives,
		"webInfo":         sidebarInfo.WebInfo,
		"wallpaperURL":    pkg.GetBingImage(),
		"title":           "首页",
	})
}

func (h *Handler) Search(c *gin.Context) {
	keyword := c.Param("keyword")

	// 输入验证：限制搜索关键词长度，防止潜在的安全问题
	if len(keyword) > 100 {
		result.Error(c, "搜索关键词过长")
		return
	}
	if len(keyword) < 1 {
		result.Error(c, "请输入搜索关键词")
		return
	}

	res, err := h.Service.PostService.Search(keyword)
	if err != nil {
		result.Error(c, "暂无记录")
		return
	}
	result.Ok(c, res)
}

func (h *Handler) Daily(c *gin.Context) {

	report := daily.GetDailyReport()

	c.String(200, report)
}
