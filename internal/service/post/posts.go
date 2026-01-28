package post

import (
	"regexp"
	"sort"
	"strings"

	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"github.com/pkg/errors"
)

// Service 文章服务
type Service struct {
}

// NewPostService 创建文章服务实例
func NewPostService() *Service {
	return &Service{}
}

// GetPost 前台获取文章（已发布、未删除）
func (service Service) GetPost(id int) (*model.Post, error) {
	return service.getPostWithConditions(id, true, true)
}

// GetPostDetail 后台获取文章（所有状态）
func (service Service) GetPostDetail(id int) (*model.Post, error) {
	return service.getPostWithConditions(id, false, false)
}

// getPostWithConditions 通用的获取文章方法
func (service Service) getPostWithConditions(id int, onlyPublished, onlyNotDeleted bool) (*model.Post, error) {
	var postInfo model.Post

	// 构建查询
	qb := NewQueryBuilder().
		Select("post.*, category.category_name, t.tag_name").
		ById(id).
		WithCategory().
		WithTags()

	// 根据条件添加过滤
	if onlyPublished {
		qb = qb.WithPublished()
	}
	if onlyNotDeleted {
		qb = qb.WithNotDeleted()
	}

	// 执行查询
	if err := qb.Scan(&postInfo); err != nil {
		return nil, errors.New("找不到记录")
	}

	// 获取标签列表
	tags, err := service.getPostTags(id, onlyPublished, onlyNotDeleted)
	if err != nil {
		return nil, err
	}
	postInfo.Tags = tags

	// 更新阅读数（仅前台）
	if onlyPublished {
		_ = service.UpdatePostReadCount(postInfo.PostId, postInfo.ReadCount)
	}

	return &postInfo, nil
}

// getPostTags 获取文章标签（提取公共逻辑）
func (service Service) getPostTags(postId int, onlyPublished, onlyNotDeleted bool) ([]string, error) {
	var tags []string

	qb := NewQueryBuilder().
		Select("t.tag_name").
		ById(postId).
		WithTags()

	if onlyPublished {
		qb = qb.WithPublished()
	}
	if onlyNotDeleted {
		qb = qb.WithNotDeleted()
	}

	if err := qb.Scan(&tags); err != nil {
		return nil, errors.New("获取标签失败")
	}

	return tags, nil
}

// GetLatestPosts 最新文章
func (service Service) GetLatestPosts() ([]model.LatestPosts, error) {
	var latestPosts []model.LatestPosts

	err := NewQueryBuilder().
		Select("post.title, post.author, post.pub_time, post.cover_image, post.post_slug, post.pub_time, category.category_name, category.category_id").
		WithPublished().
		WithNotDeleted().
		WithCategory().
		WithOrderByTime(true).
		Limit(5).
		Find(&latestPosts)

	if err != nil {
		return nil, errors.New("查询失败: " + err.Error())
	}

	return latestPosts, nil
}

// GetCoverPosts 首页文章列表
func (service Service) GetCoverPosts(pageNum int) ([]model.LatestPosts, int64, error) {
	var coverPosts []model.LatestPosts
	var total int64

	// 先统计总数（不需要 JOIN）
	countQb := NewQueryBuilder().
		WithPublished().
		WithNotDeleted()

	if err := countQb.Count(&total); err != nil {
		return nil, 0, errors.New("统计失败: " + err.Error())
	}

	// 查询列表（需要 JOIN）
	qb := NewQueryBuilder().
		Select("post.title, post.author, post.pub_time, post.cover_image, post.top, post.post_slug, category.category_name, category.category_id, t.tag_name").
		WithPublished().
		WithNotDeleted().
		WithCategory().
		WithTags().
		WithOrderByTopAndTime().
		WithPagination(pageNum, DefaultPageSize)

	err := qb.Find(&coverPosts)

	if err != nil {
		return nil, 0, errors.New("查询失败: " + err.Error())
	}

	return coverPosts, total, nil
}

// Total 统计文章总数
func (service Service) Total() (total int64) {
	NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		Count(&total)
	return
}

// UpdatePostReadCount 更新文章阅读数
func (service Service) UpdatePostReadCount(postId uint64, readCount uint32) error {
	err := global.GORM.Table(model.TPostsTable).
		Where("post_id", postId).
		Update("read_count", readCount+1).Error
	if err != nil {
		return errors.New("更新失败: " + err.Error())
	}
	return nil
}

// GetArchivePosts 获取归档文章
func (service Service) GetArchivePosts(year, month string) (*model.ArchivesPosts, error) {
	var years []string
	global.GORM.Table(model.TPostsTable).
		Raw(`SELECT strftime('%Y-%m', pub_time, 'unixepoch') year FROM post WHERE is_published = 1 AND is_deleted = 0 GROUP BY year`).
		Scan(&years)

	m := make(map[string][]model.ArchivePosts, 0)
	archives := model.ArchivesPosts{}

	if year != "" && month != "" {
		posts := make([]model.ArchivePosts, 0)
		global.GORM.Table(model.TPostsTable).
			Raw("SELECT title, post_slug, pub_time, cover_image FROM post WHERE is_published = 1 AND is_deleted = 0 AND strftime('%Y/%m', pub_time, 'unixepoch') = ?", year+"/"+month).
			Scan(&posts)
		m[year+"-"+month] = posts
		archives.Archives = m
		return &archives, nil
	}

	for _, year := range years {
		posts := make([]model.ArchivePosts, 0)
		global.GORM.Table(model.TPostsTable).
			Raw("SELECT title, post_slug, pub_time, cover_image FROM post WHERE is_published = 1 AND is_deleted = 0 AND strftime('%Y-%m', pub_time, 'unixepoch') = ?", year).
			Scan(&posts)
		m[year] = posts
	}

	archives.Archives = m
	return &archives, nil
}

// GetArchivePostsPaged 分页获取归档文章（按年月分组）
func (service Service) GetArchivePostsPaged(pageNum, pageSize int) (*model.ArchivesPosts, int, int, error) {
	// 获取所有年月分组，按时间倒序
	var yearMonths []string
	global.GORM.Raw(`SELECT strftime('%Y-%m', pub_time, 'unixepoch') as year_month
		FROM post
		WHERE is_published = 1 AND is_deleted = 0
		GROUP BY year_month
		ORDER BY year_month DESC`).Scan(&yearMonths)

	totalGroups := len(yearMonths)

	// 获取总文章数
	var totalPosts int64
	NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		Count(&totalPosts)

	// 计算分页
	if pageNum <= 0 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize
	end := start + pageSize

	if start >= totalGroups {
		return &model.ArchivesPosts{Archives: make(map[string][]model.ArchivePosts)}, totalGroups, int(totalPosts), nil
	}

	if end > totalGroups {
		end = totalGroups
	}

	// 获取当前页的年月分组
	pagedYearMonths := yearMonths[start:end]

	// 为每个年月分组获取文章
	m := make(map[string][]model.ArchivePosts)
	for _, ym := range pagedYearMonths {
		posts := make([]model.ArchivePosts, 0)
		global.GORM.Table(model.TPostsTable).
			Raw("SELECT title, post_slug, pub_time, cover_image FROM post WHERE is_published = 1 AND is_deleted = 0 AND strftime('%Y-%m', pub_time, 'unixepoch') = ? ORDER BY pub_time DESC", ym).
			Scan(&posts)
		m[ym] = posts
	}

	archives := &model.ArchivesPosts{Archives: m}
	return archives, totalGroups, int(totalPosts), nil
}

// GetPostsByCategory 根据分类获取文章
func (service Service) GetPostsByCategory(category string, pageNum int) ([]model.TagCategoryPosts, int64, error) {
	var posts []model.TagCategoryPosts
	var total int64

	// 先统计总数
	countQb := NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		WithCategory().
		ByCategory(category)

	if err := countQb.Count(&total); err != nil {
		return nil, 0, err
	}

	// 查询列表
	qb := NewQueryBuilder().
		Select("post.title, post.post_slug, post.pub_time, post.cover_image").
		WithPublished().
		WithNotDeleted().
		WithCategory().
		ByCategory(category).
		WithPagination(pageNum, DefaultPageSize)

	err := qb.Find(&posts)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPostsByTag 根据标签获取文章
func (service Service) GetPostsByTag(tagName string, pageNum int) ([]model.TagCategoryPosts, int64, error) {
	var posts []model.TagCategoryPosts
	var total int64

	// 先统计总数
	countQb := NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		WithTags().
		ByTag(tagName)

	if err := countQb.Count(&total); err != nil {
		return nil, 0, err
	}

	// 查询列表
	qb := NewQueryBuilder().
		Select("post.title, post.post_slug, post.pub_time, post.cover_image").
		WithPublished().
		WithNotDeleted().
		WithTags().
		ByTag(tagName).
		WithPagination(pageNum, DefaultPageSize)

	err := qb.Find(&posts)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPostsByCategoryId 根据分类 ID 获取文章
func (service Service) GetPostsByCategoryId(id int64) ([]model.Post, error) {
	posts := make([]model.Post, 0)

	err := NewQueryBuilder().
		ByCategoryId(id).
		Find(&posts)

	if err != nil {
		return nil, errors.New("查询失败")
	}

	return posts, nil
}

// GetPostsArchive 获取归档数据
func (service Service) GetPostsArchive(num int64) ([]map[string]any, int64, error) {
	var archives []struct {
		Year  string
		Month string
		Count int
	}
	global.GORM.Raw("SELECT strftime('%Y', create_time, 'unixepoch') AS year, strftime('%m', create_time, 'unixepoch') AS month, COUNT(*) AS count FROM post GROUP BY year, month ORDER BY year DESC, month DESC").Scan(&archives)

	// 生成归档数据
	archiveData := make([]map[string]any, len(archives))
	for i, archive := range archives {
		// 查询归档时间段内的文章
		var articles []model.ArchivePosts
		global.GORM.Table("post").
			Where("strftime('%Y', create_time, 'unixepoch') = ? AND strftime('%m', create_time, 'unixepoch') = ?", archive.Year, archive.Month).
			Order("create_time DESC").
			Find(&articles)

		// 生成文章列表数据
		articleData := make([]map[string]interface{}, len(articles))
		for j, article := range articles {
			articleData[j] = map[string]interface{}{
				"title":      article.Title,
				"created_at": article.PubTime,
			}
		}

		// 生成归档数据
		archiveData[i] = map[string]any{
			"year":     archive.Year,
			"month":    archive.Month,
			"count":    archive.Count,
			"articles": articleData,
		}
	}

	return archiveData, 1, nil
}

// Search 搜索文章
func (service Service) Search(keyword string) (posts []model.SearchPost, err error) {
	return service.SearchPaged(keyword, 1, 20)
}

// SearchPaged 分页搜索
func (service Service) SearchPaged(keyword string, pageNum, pageSize int) ([]model.SearchPost, error) {
	if keyword == "" {
		return []model.SearchPost{}, nil
	}

	// 查询匹配的文章
	var rawPosts []struct {
		PostId     uint32 `gorm:"column:post_id"`
		Title      string `gorm:"column:title"`
		PostSlug   string `gorm:"column:post_slug"`
		Summary    string `gorm:"column:summary"`
		Content    string `gorm:"column:post_content"`
		CoverImage string `gorm:"column:cover_image"`
		CreateTime int64  `gorm:"column:create_time"`
	}

	err := NewQueryBuilder().
		Select("post_id, title, post_slug, summary, post_content, cover_image, create_time").
		WithPublished().
		WithNotDeleted().
		WithSearch(keyword).
		WithOrderByTime(true).
		WithPagination(pageNum, pageSize).
		Find(&rawPosts)

	if err != nil {
		return nil, err
	}

	// 处理结果：计算相关度和高亮
	posts := make([]model.SearchPost, 0, len(rawPosts))
	for _, raw := range rawPosts {
		relevance := 0
		highlight := ""

		// 计算相关度分数
		lowerKeyword := strings.ToLower(keyword)
		if strings.Contains(strings.ToLower(raw.Title), lowerKeyword) {
			relevance += 10 // 标题匹配权重最高
		}
		if strings.Contains(strings.ToLower(raw.Summary), lowerKeyword) {
			relevance += 5 // 摘要匹配
		}
		if strings.Contains(strings.ToLower(raw.Content), lowerKeyword) {
			relevance += 1 // 内容匹配
		}

		// 生成高亮片段
		highlight = extractHighlight(raw.Content, keyword, 100)
		if highlight == "" {
			highlight = extractHighlight(raw.Summary, keyword, 100)
		}
		if highlight == "" && len(raw.Summary) > 0 {
			if len(raw.Summary) > 100 {
				highlight = raw.Summary[:100] + "..."
			} else {
				highlight = raw.Summary
			}
		}

		posts = append(posts, model.SearchPost{
			Id:         raw.PostId,
			Title:      highlightKeyword(raw.Title, keyword),
			PostSlug:   raw.PostSlug,
			Summary:    raw.Summary,
			Highlight:  highlight,
			CoverImage: raw.CoverImage,
			CreateAt:   raw.CreateTime,
			Relevance:  relevance,
		})
	}

	// 按相关度排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Relevance > posts[j].Relevance
	})

	return posts, nil
}

// SearchWithResult 带统计的搜索
func (service Service) SearchWithResult(keyword string, pageNum, pageSize int) (*model.SearchResult, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	posts, err := service.SearchPaged(keyword, pageNum, pageSize)
	if err != nil {
		return nil, err
	}

	// 统计总数
	var total int64
	NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		WithSearch(keyword).
		Count(&total)

	return &model.SearchResult{
		Posts:    posts,
		Total:    total,
		Keyword:  keyword,
		PageNum:  pageNum,
		PageSize: pageSize,
	}, nil
}

// extractHighlight 提取包含关键词的高亮片段
func extractHighlight(content, keyword string, maxLen int) string {
	if content == "" || keyword == "" {
		return ""
	}

	lowerContent := strings.ToLower(content)
	lowerKeyword := strings.ToLower(keyword)
	idx := strings.Index(lowerContent, lowerKeyword)

	if idx == -1 {
		return ""
	}

	// 计算片段起始位置
	start := idx - maxLen/2
	if start < 0 {
		start = 0
	}

	// 计算片段结束位置
	end := idx + len(keyword) + maxLen/2
	if end > len(content) {
		end = len(content)
	}

	// 提取片段
	snippet := content[start:end]

	// 添加省略号
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	return highlightKeyword(snippet, keyword)
}

// highlightKeyword 高亮关键词
func highlightKeyword(text, keyword string) string {
	if text == "" || keyword == "" {
		return text
	}

	// 使用正则进行大小写不敏感替换
	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(keyword))
	return re.ReplaceAllString(text, "<mark>$0</mark>")
}
