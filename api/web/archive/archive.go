package archive

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	inpkg "gitee.com/jieepre/keepblog/internal/pkg"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/page"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
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

	// 每页显示6个年月分组
	pageSize := 6

	// 使用分页方法获取归档数据
	archivePosts, totalGroups, totalPosts, err := h.Service.PostService.GetArchivePostsPaged(int(pageNum), pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	// 计算总页数（基于年月分组数量）
	totalPages := totalGroups / pageSize
	if totalGroups%pageSize != 0 {
		totalPages++
	}

	//传到模板中需要转换成template.HTML类型，否则html代码会被转义
	index, _ := page.HandleIndexWithPageSize(totalGroups, int(pageNum), pageSize, "archives/page")
	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	pagination := template.HTML(index)
	// 归档翻页页 canonical 指向自身。
	archivePath := "/archives"
	if pageNum > 1 {
		archivePath = fmt.Sprintf("/archives/page/%d", pageNum)
	}

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
		"totalPosts":      totalPosts,
		"currentPage":     pageNum,
		"totalPages":      totalPages,
		"canonical":       inpkg.CanonicalURL(site.URL, archivePath),
	})
}

func (h Handler) ArchivesInfo(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	sidebarInfo := h.Service.SidebarService.Sidebar()
	archivePosts, err := h.Service.PostService.GetArchivePosts(year, month)
	if err != nil {
		slog.Errorf("归档查询错误: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	// 计算当前筛选条件下的文章总数
	totalPosts := 0
	if archivePosts != nil && archivePosts.Archives != nil {
		for _, posts := range archivePosts.Archives {
			totalPosts += len(posts)
		}
	}

	site, err := h.Service.WebSiteService.GetWebSite()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	// 格式化标题
	title := "归档"
	archivePath := "/archives"
	pagedesc := ""
	if year != "" && month != "" {
		title = year + "年" + month + "月 归档"
		archivePath = "/archives/" + year + "/" + month
		pagedesc = fmt.Sprintf("%s年%s月发布的全部文章，共 %d 篇。", year, month, totalPosts)
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
		"title":           title,
		"totalPosts":      totalPosts,
		"currentPage":     1,
		"totalPages":      1,
		"filterYear":      year,
		"filterMonth":     month,
		"canonical":       inpkg.CanonicalURL(site.URL, archivePath),
		"pagedesc":        pagedesc,
	})
}
