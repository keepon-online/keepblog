package page

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Info struct {
	/**
	 * 当前页
	 */
	PageNum int `json:"pageNum"`
	/**
	 * 每页的数量
	 */
	PageSize int `json:"pageSize"`

	/**
	 * 总共多少条数据
	 */
	Total int `json:"total"`

	List any `json:"list"`
}

func PaginationInfo(list any, pageNum, pageSize, total int) *Info {
	return &Info{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}
}

// Page 标准分页结构体，接收最原始的DO
// 建议在外部再建一个字段一样的结构体，用以将DO转换成DTO或VO
// https://github.com/yafeng-Soong/gorm-paginator
type Page[T any] struct {
	CurrentPage int64 `json:"currentPage"`
	PageSize    int64 `json:"pageSize"`
	Total       int64 `json:"total"`
	Pages       int64 `json:"pages"`
	Data        []T   `json:"data"`
}

// SelectPages 各种查询条件先在query设置好后再放进来
func (page *Page[T]) SelectPages(query *gorm.DB) (e error) {
	var model T
	query.Model(&model).Count(&page.Total)
	if page.Total == 0 {
		page.Data = []T{}
		return
	}
	e = query.Model(&model).Scopes(Paginate(page)).Find(&page.Data).Error
	return
}

func Paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.CurrentPage <= 0 {
			page.CurrentPage = 1
		}
		switch {
		case page.PageSize > 10000:
			page.PageSize = 10000 // 限制一下分页大小
		case page.PageSize <= 0:
			page.PageSize = 10
		}
		page.Pages = page.Total / page.PageSize
		if page.Total%page.PageSize != 0 {
			page.Pages++
		}
		p := page.CurrentPage
		if page.CurrentPage > page.Pages {
			p = page.Pages
		}
		size := page.PageSize
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}

func HandleIndex(totalRecords, page int, temp string) (string, error) {
	if page == 0 {
		page = 1
	}

	// 每页显示的记录数
	pageSize := 10

	// 计算总页数
	totalPages := totalRecords / pageSize
	if totalRecords%pageSize != 0 {
		totalPages++
	}

	if page > totalPages {

		return "", errors.New("invalid page num")
	}

	// 计算当前页的起始记录和结束记录
	startRecord := (page - 1) * pageSize
	endRecord := startRecord + pageSize
	if endRecord > totalRecords {
		endRecord = totalRecords
	}

	// 构造分页数据
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	// 构造分页 HTML
	var paginationHTML string
	if totalPages > 1 {
		paginationHTML += `<div><ul class="pagination">`

		// 上一页
		if page > 1 {
			paginationHTML += fmt.Sprintf(`<li><a class="extend prev" href="/%s/%d" data-pjax-content="true"><i class="fas fa-chevron-left fa-fw"></i></a></li>`, temp, page-1)
		} else {
			paginationHTML += `<li class="disabled" style="pointer-events: none;cursor: not-allowed;" ><a href="#" style="cursor: not-allowed;pointer-events:none"><i class="fas fa-chevron-left fa-fw"></i></a></li>`
		}

		// 页码
		var start, end int
		if totalPages <= 7 {
			start, end = 1, totalPages
		} else {
			if page <= 4 {
				start, end = 1, 7
			} else if page >= totalPages-3 {
				start, end = totalPages-6, totalPages
			} else {
				start, end = page-3, page+3
			}
		}

		if start > 1 {
			paginationHTML += fmt.Sprintf(`<li><a href="/%s/%d" data-pjax-content="true">1</a></li>`, temp, 1)
			if start > 2 {
				paginationHTML += `<li class="disabled"><a  style="pointer-events: none;cursor: not-allowed;" >...</a></li>`
			}
		}

		for _, p := range pages[start-1 : end] {
			if p == page {
				paginationHTML += fmt.Sprintf(`<li class="page-number current"><a href="/%s/%d" data-pjax-content="true">%d</a></li>`, temp, p, p)
			} else {
				paginationHTML += fmt.Sprintf(`<li><a href="/%s/%d" data-pjax-content="true">%d</a></li>`, temp, p, p)
			}
		}

		if end < totalPages {
			if end < totalPages-1 {
				paginationHTML += `<li class="disabled"><a  style="pointer-events: none;cursor: not-allowed;" >...</a></li>`
			}
			paginationHTML += fmt.Sprintf(`<li><a href="/%s/%d" data-pjax-content="true">%d</a></li>`, temp, totalPages, totalPages)
		}

		// 下一页
		if page < totalPages {
			paginationHTML += fmt.Sprintf(`<li><a class="extend next" href="/%s/%d" data-pjax-content="true"><i class="fas fa-chevron-right fa-fw"></i></a></li>`, temp, page+1)
		} else {
			paginationHTML += `<li class="disabled" style="pointer-events: none;cursor: not-allowed;"><a href="#" style="cursor: not-allowed;pointer-events:none"><i class="fas fa-chevron-right fa-fw"></i></a></li>`
		}

		paginationHTML += `</ul></div>`
	}

	return paginationHTML, nil

}
