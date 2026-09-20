package post

import (
	"regexp"
	"sort"
	"strings"

	"errors"
	"gorm.io/gorm"

	"gitee.com/jieepre/keepblog/internal/model"
)

// Service 文章服务
type Service struct {
	db *gorm.DB
}

// NewPostService 创建文章服务实例
func NewPostService(db *gorm.DB) *Service {
	return &Service{db: db}
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
	qb := NewQueryBuilder(service.db).
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

	// 查询未命中时保留历史兼容行为：返回零值文章，不继续查询标签或更新阅读数。
	if postInfo.PostId == 0 {
		return &postInfo, nil
	}

	// 获取标签列表
	tags, err := service.getPostTags(id, onlyPublished, onlyNotDeleted)
	if err != nil {
		return nil, err
	}
	postInfo.Tags = tags

	// 更新阅读数（仅前台）。更新使用 SQL 原子自增，避免并发请求丢失计数。
	if onlyPublished {
		if err := service.UpdatePostReadCount(postInfo.PostId); err != nil {
			return nil, err
		}
	}

	return &postInfo, nil
}

// getPostTags 获取文章标签。使用 INNER JOIN + 非空条件，
// 无标签文章返回空切片而不是让 NULL 扫描到 string 失败。
func (service Service) getPostTags(postId int, onlyPublished, onlyNotDeleted bool) ([]string, error) {
	tags := make([]string, 0)
	query := service.db.Table(model.TTagTable+" t").
		Select("t.tag_name").
		Joins("JOIN post_tag pt ON t.tag_id = pt.tag_id").
		Joins("JOIN post p ON pt.post_id = p.post_id").
		Where("pt.post_id = ? AND t.tag_name IS NOT NULL", postId)
	if onlyPublished {
		query = query.Where("p.is_published = ?", 1)
	}
	if onlyNotDeleted {
		query = query.Where("p.is_deleted = ?", 0)
	}
	if err := query.Scan(&tags).Error; err != nil {
		return nil, errors.New("获取标签失败: " + err.Error())
	}
	return tags, nil
}

// GetLatestPosts 最新文章
func (service Service) GetLatestPosts() ([]model.LatestPosts, error) {
	var latestPosts []model.LatestPosts

	err := NewQueryBuilder(service.db).
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

// GetPublishedPostsForFeed 取已发布文章列表（用于 RSS/sitemap），按发布时间倒序。
// 返回 model.Post 的部分字段（title/post_slug/summary/cover_image/post_content/pub_time/last_modified_time）；
// post_content 供 RSS 在 summary 为空时生成纯文本兜底摘要。
func (service Service) GetPublishedPostsForFeed(limit int) ([]model.Post, error) {
	var posts []model.Post

	err := NewQueryBuilder(service.db).
		Select("post.title, post.post_slug, post.summary, post.cover_image, post.post_content, post.pub_time, post.last_modified_time").
		WithPublished().
		WithNotDeleted().
		Order("post.pub_time DESC").
		Limit(limit).
		Find(&posts)

	if err != nil {
		return nil, errors.New("查询失败: " + err.Error())
	}

	return posts, nil
}

// GetAdjacentPosts 取上一篇/下一篇（按发布时间相邻，已发布未删除，排除当前文章）。
// 上一篇：发布时间早于当前文章的最近一篇；下一篇：发布时间晚于当前文章的最早一篇。
// 当前文章的 pub_time 作为基准传入，避免再查一次。
func (service Service) GetAdjacentPosts(postId uint64, currentPubTime uint64) (prev, next *model.LatestPosts, err error) {
	const baseWhere = "is_published = 1 AND is_deleted = 0 AND post_id <> ?"

	// 上一篇：pub_time < 当前，按 pub_time DESC 取第一条
	var prevPost model.LatestPosts
	if err = service.db.Table(model.TPostsTable).
		Select("title, post_slug, cover_image, pub_time").
		Where(baseWhere+" AND pub_time < ?", postId, currentPubTime).
		Order("pub_time DESC").
		Limit(1).
		Scan(&prevPost).Error; err == nil && prevPost.Title != "" {
		prev = &prevPost
	}

	// 下一篇：pub_time > 当前，按 pub_time ASC 取第一条
	var nextPost model.LatestPosts
	if err = service.db.Table(model.TPostsTable).
		Select("title, post_slug, cover_image, pub_time").
		Where(baseWhere+" AND pub_time > ?", postId, currentPubTime).
		Order("pub_time ASC").
		Limit(1).
		Scan(&nextPost).Error; err == nil && nextPost.Title != "" {
		next = &nextPost
	}

	return prev, next, nil
}

// GetRelatedPosts 按当前文章的标签取相关文章（排除自身，去重，limit）。
// tags 为当前文章的标签名列表，通过 post_tag/tag JOIN 反查命中这些标签的其他文章。
func (service Service) GetRelatedPosts(postId uint64, tags []string, limit int) ([]model.LatestPosts, error) {
	if len(tags) == 0 {
		return []model.LatestPosts{}, nil
	}
	if limit <= 0 {
		limit = 6
	}

	var posts []model.LatestPosts
	err := service.db.Table(model.TPostsTable).
		Select("post.title, post.post_slug, post.cover_image, post.pub_time").
		Joins("JOIN post_tag pt ON post.post_id = pt.post_id").
		Joins("JOIN tag t ON pt.tag_id = t.tag_id").
		Where("t.tag_name IN ?", tags).
		Where("post.is_published = 1 AND post.is_deleted = 0").
		Where("post.post_id <> ?", postId).
		Group("post.post_id").
		Order("post.pub_time DESC").
		Limit(limit).
		Scan(&posts).Error

	if err != nil {
		return nil, errors.New("查询相关文章失败: " + err.Error())
	}

	return posts, nil
}

// GetCoverPosts 首页文章列表
func (service Service) GetCoverPosts(pageNum int) ([]model.LatestPosts, int64, error) {
	var coverPosts []model.LatestPosts
	var total int64

	// 先统计总数（不需要 JOIN）
	countQb := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted()

	if err := countQb.Count(&total); err != nil {
		return nil, 0, errors.New("统计失败: " + err.Error())
	}

	// 查询列表（需要 JOIN）
	qb := NewQueryBuilder(service.db).
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
func (service Service) Total() (int64, error) {
	var total int64
	if err := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted().
		Count(&total); err != nil {
		return 0, errors.New("统计文章失败: " + err.Error())
	}
	return total, nil
}

// UpdatePostReadCount 原子增加文章阅读数。
// 参数不再接收调用方读到的旧值，避免并发请求的读-改-写覆盖彼此更新。
func (service Service) UpdatePostReadCount(postId uint64) error {
	result := service.db.Table(model.TPostsTable).
		Where("post_id = ?", postId).
		UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if result.Error != nil {
		return errors.New("更新失败: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New("更新失败: 文章不存在")
	}
	return nil
}

// archiveRow 是归档查询的扁平结果。SQLite 的 strftime 是当前 SQLite 存储
// 格式的一部分，查询集中在这里，避免调用方按年月循环发起 N+1 查询。
type archiveRow struct {
	YearMonth  string `gorm:"column:year_month"`
	Title      string `gorm:"column:title"`
	PostSlug   string `gorm:"column:post_slug"`
	PubTime    uint64 `gorm:"column:pub_time"`
	CoverImage string `gorm:"column:cover_image"`
}

// fetchArchiveRows 一次查询指定年月的全部文章，再在内存中分组。
func (service Service) fetchArchiveRows(yearMonths []string) ([]archiveRow, error) {
	if len(yearMonths) == 0 {
		return []archiveRow{}, nil
	}
	rows := make([]archiveRow, 0)
	err := service.db.Table(model.TPostsTable).
		Select("strftime('%Y-%m', pub_time, 'unixepoch') AS year_month, title, post_slug, pub_time, cover_image").
		Where("is_published = ? AND is_deleted = ?", 1, 0).
		Where("strftime('%Y-%m', pub_time, 'unixepoch') IN ?", yearMonths).
		Order("pub_time DESC").
		Scan(&rows).Error
	return rows, err
}

func groupArchiveRows(rows []archiveRow, yearMonths []string) map[string][]model.ArchivePosts {
	archives := make(map[string][]model.ArchivePosts, len(yearMonths))
	for _, ym := range yearMonths {
		archives[ym] = make([]model.ArchivePosts, 0)
	}
	for _, row := range rows {
		archives[row.YearMonth] = append(archives[row.YearMonth], model.ArchivePosts{
			Title: row.Title, PostSlug: row.PostSlug, PubTime: row.PubTime, CoverImage: row.CoverImage,
		})
	}
	return archives
}

// GetArchivePosts 获取归档文章。所有年月模式固定为两条查询（分组 + 文章），
// 指定月份模式只发起一条文章查询，不再按年月逐组查询。
func (service Service) GetArchivePosts(year, month string) (*model.ArchivesPosts, error) {
	if year != "" && month != "" {
		ym := year + "-" + month
		rows, err := service.fetchArchiveRows([]string{ym})
		if err != nil {
			return nil, errors.New("查询归档文章失败: " + err.Error())
		}
		return &model.ArchivesPosts{Archives: groupArchiveRows(rows, []string{ym})}, nil
	}

	var yearMonths []string
	if err := service.db.Table(model.TPostsTable).
		Select("strftime('%Y-%m', pub_time, 'unixepoch') AS year_month").
		Where("is_published = ? AND is_deleted = ?", 1, 0).
		Group("year_month").
		Order("year_month DESC").
		Scan(&yearMonths).Error; err != nil {
		return nil, errors.New("查询归档分组失败: " + err.Error())
	}
	rows, err := service.fetchArchiveRows(yearMonths)
	if err != nil {
		return nil, errors.New("查询归档文章失败: " + err.Error())
	}
	return &model.ArchivesPosts{Archives: groupArchiveRows(rows, yearMonths)}, nil
}

// GetArchivePostsPaged 分页获取归档文章（按年月分组）。年月分组页内的
// 文章通过 IN 一次取回，将原先每个年月一次查询的 N+1 降为固定 3 条查询。
func (service Service) GetArchivePostsPaged(pageNum, pageSize int) (*model.ArchivesPosts, int, int, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	var yearMonths []string
	if err := service.db.Table(model.TPostsTable).
		Select("strftime('%Y-%m', pub_time, 'unixepoch') AS year_month").
		Where("is_published = ? AND is_deleted = ?", 1, 0).
		Group("year_month").
		Order("year_month DESC").
		Scan(&yearMonths).Error; err != nil {
		return nil, 0, 0, errors.New("查询归档分组失败: " + err.Error())
	}

	var totalPosts int64
	if err := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted().
		Count(&totalPosts); err != nil {
		return nil, len(yearMonths), 0, errors.New("统计文章失败: " + err.Error())
	}

	totalGroups := len(yearMonths)
	start := (pageNum - 1) * pageSize
	if start >= totalGroups {
		return &model.ArchivesPosts{Archives: make(map[string][]model.ArchivePosts)}, totalGroups, int(totalPosts), nil
	}
	end := start + pageSize
	if end > totalGroups {
		end = totalGroups
	}
	pagedYearMonths := yearMonths[start:end]
	rows, err := service.fetchArchiveRows(pagedYearMonths)
	if err != nil {
		return nil, totalGroups, int(totalPosts), errors.New("查询归档文章失败: " + err.Error())
	}
	return &model.ArchivesPosts{Archives: groupArchiveRows(rows, pagedYearMonths)}, totalGroups, int(totalPosts), nil
}

// GetPostsByCategory 根据分类获取文章
func (service Service) GetPostsByCategory(category string, pageNum int) ([]model.TagCategoryPosts, int64, error) {
	var posts []model.TagCategoryPosts
	var total int64

	// 先统计总数
	countQb := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted().
		WithCategory().
		ByCategory(category)

	if err := countQb.Count(&total); err != nil {
		return nil, 0, err
	}

	// 查询列表
	qb := NewQueryBuilder(service.db).
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
	countQb := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted().
		WithTags().
		ByTag(tagName)

	if err := countQb.Count(&total); err != nil {
		return nil, 0, err
	}

	// 查询列表
	qb := NewQueryBuilder(service.db).
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

	err := NewQueryBuilder(service.db).
		ByCategoryId(id).
		Find(&posts)

	if err != nil {
		return nil, errors.New("查询失败")
	}

	return posts, nil
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

	err := NewQueryBuilder(service.db).
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
	if err := NewQueryBuilder(service.db).
		WithPublished().
		WithNotDeleted().
		WithSearch(keyword).
		Count(&total); err != nil {
		return nil, errors.New("统计搜索结果失败: " + err.Error())
	}

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
