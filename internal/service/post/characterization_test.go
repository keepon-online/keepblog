package post

import (
	"strings"
	"testing"

	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/testutil"
)

// 以下为特征测试：钉住 post 服务当前的对外行为，重构（依赖注入、
// N+1 修复、错误契约改造）过程中任何非预期的行为变化都会在此暴露。

func TestGetPost_OnlyPublishedVisible(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	p, err := s.GetPost(1)
	if err != nil {
		t.Fatalf("GetPost(1) 报错: %v", err)
	}
	if p.Title != "Go 快速入门" {
		t.Errorf("标题 = %q, want %q", p.Title, "Go 快速入门")
	}
	if p.CategoryName != "Go" {
		t.Errorf("CategoryName = %q, want %q", p.CategoryName, "Go")
	}
	if len(p.Tags) != 1 || p.Tags[0] != "go" {
		t.Errorf("Tags = %v, want [go]", p.Tags)
	}
}

func TestGetPost_IncrementsReadCount(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	before, err := s.GetPost(1)
	if err != nil {
		t.Fatalf("GetPost 报错: %v", err)
	}
	if before.ReadCount != 1 {
		t.Fatalf("第一次读取后 ReadCount = %d, want 1（种子 0 + 1）", before.ReadCount)
	}
	after, err := s.GetPost(1)
	if err != nil {
		t.Fatalf("第二次 GetPost 报错: %v", err)
	}
	if after.ReadCount != 2 {
		t.Errorf("第二次读取后 ReadCount = %d, want 2", after.ReadCount)
	}
}

func TestGetPost_DraftDeletedMissingReturnZeroValue(t *testing.T) {
	// 特征（怪癖）：GetPost 对草稿/已删除/不存在的文章不报错，而是返回零值文章，
	// 由调用方（api/web 渲染层）自行判断 PostId==0。重构若收紧此契约需同步调用方。
	testutil.NewTestDB(t)
	s := NewPostService()

	for _, id := range []int{90, 91, 9999} {
		p, err := s.GetPost(id)
		if err != nil {
			t.Errorf("GetPost(%d) 报错 %v（现状应返回零值且无错误）", id, err)
			continue
		}
		if p.PostId != 0 {
			t.Errorf("GetPost(%d) PostId = %d, want 0", id, p.PostId)
		}
	}
}

func TestGetPostDetail_HasTags(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	p, err := s.GetPostDetail(2)
	if err != nil {
		t.Fatalf("GetPostDetail(2) 报错: %v", err)
	}
	if p.Title != "MySQL 统计实战" {
		t.Errorf("Title = %q", p.Title)
	}
	if len(p.Tags) != 1 || p.Tags[0] != "mysql" {
		t.Errorf("Tags = %v, want [mysql]", p.Tags)
	}

	// 不存在的 id 同样是零值 + 无错误
	p99, err := s.GetPostDetail(9999)
	if err != nil || p99.PostId != 0 {
		t.Errorf("GetPostDetail(9999) = (%d, %v), want (0, nil)", p99.PostId, err)
	}
}

func TestGetPostDetail_PostWithoutTagsErrors(t *testing.T) {
	// 特征（已知怪癖）：无标签文章的详情查询因 LEFT JOIN 产生 NULL tag_name 行、
	// 扫描 []string 失败而报错。生产环境后台编辑无标签文章同样会失败，
	// 数据层治理阶段修复后此断言应反转为"正常返回空标签列表"。
	testutil.NewTestDB(t)
	s := NewPostService()

	for _, id := range []int{90, 12} { // 90=草稿、12=无标签的填充文章
		if _, err := s.GetPostDetail(id); err == nil {
			t.Errorf("GetPostDetail(%d) 现状应报错（NULL tag_name 扫描失败）", id)
		}
	}
}

func TestGetLatestPosts_LimitFiveByTimeDesc(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	posts, err := s.GetLatestPosts()
	if err != nil {
		t.Fatalf("GetLatestPosts 报错: %v", err)
	}
	if len(posts) != 5 {
		t.Fatalf("len = %d, want 5", len(posts))
	}
	wantFirst := "填充文章G" // 2024-05-07 发布，时间最新
	if posts[0].Title != wantFirst {
		t.Errorf("最新一篇 = %q, want %q", posts[0].Title, wantFirst)
	}
	for i := 1; i < len(posts); i++ {
		if posts[i].PubTime > posts[i-1].PubTime {
			t.Errorf("应按发布时间倒序: posts[%d].PubTime=%d > posts[%d].PubTime=%d",
				i, posts[i].PubTime, i-1, posts[i-1].PubTime)
		}
	}
}

func TestTotal_OnlyPublishedNotDeleted(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	if got := s.Total(); got != 12 {
		t.Errorf("Total = %d, want 12（草稿/已删除不计）", got)
	}
}

func TestGetCoverPosts_PaginationAndTopFirst(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	page1, total, err := s.GetCoverPosts(1)
	if err != nil {
		t.Fatalf("GetCoverPosts(1) 报错: %v", err)
	}
	if total != 12 {
		t.Errorf("total = %d, want 12", total)
	}
	if len(page1) != 10 {
		t.Fatalf("第 1 页条数 = %d, want 10（DefaultPageSize）", len(page1))
	}
	if page1[0].Title != "Go 快速入门" {
		t.Errorf("置顶文章应在首位, got %q", page1[0].Title)
	}

	page2, _, err := s.GetCoverPosts(2)
	if err != nil {
		t.Fatalf("GetCoverPosts(2) 报错: %v", err)
	}
	if len(page2) != 2 {
		t.Fatalf("第 2 页条数 = %d, want 2", len(page2))
	}
}

func TestGetPostsByCategory(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	posts, total, err := s.GetPostsByCategory("MySQL", 1)
	if err != nil {
		t.Fatalf("GetPostsByCategory 报错: %v", err)
	}
	if total != 2 || len(posts) != 2 {
		t.Errorf("MySQL 分类 total=%d len=%d, want 2/2", total, len(posts))
	}
	for _, p := range posts {
		if p.Title == "" {
			t.Error("返回了空标题文章")
		}
	}

	if _, total, _ := s.GetPostsByCategory("不存在", 1); total != 0 {
		t.Errorf("未知分类 total = %d, want 0", total)
	}
}

func TestGetPostsByTag(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	posts, total, err := s.GetPostsByTag("go", 1)
	if err != nil {
		t.Fatalf("GetPostsByTag 报错: %v", err)
	}
	if total != 2 || len(posts) != 2 {
		t.Errorf("go 标签 total=%d len=%d, want 2/2", total, len(posts))
	}
}

func TestGetAdjacentPosts(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	// P2 发布于 2024-02-10：上一篇应为 P1（2024-01-10），下一篇应为 P3（2024-03-10）
	p2, err := s.GetPostDetail(2)
	if err != nil {
		t.Fatalf("GetPostDetail(2) 报错: %v", err)
	}
	prev, next, err := s.GetAdjacentPosts(p2.PostId, p2.PubTime)
	if err != nil {
		t.Fatalf("GetAdjacentPosts 报错: %v", err)
	}
	if prev == nil || prev.Title != "Go 快速入门" {
		t.Errorf("prev = %+v, want Go 快速入门", prev)
	}
	if next == nil || next.Title != "Go 并发模式" {
		t.Errorf("next = %+v, want Go 并发模式", next)
	}
}

func TestGetAdjacentPosts_AtBoundaries(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := NewPostService()

	// 填充文章 P12（2024-05-07）无标签，不能走 GetPostDetail，直接读库取时间基准
	var latest model.Post
	if err := db.First(&latest, 12).Error; err != nil {
		t.Fatalf("读取 P12 失败: %v", err)
	}
	prev, next, err := s.GetAdjacentPosts(latest.PostId, latest.PubTime)
	if err != nil {
		t.Fatalf("GetAdjacentPosts 报错: %v", err)
	}
	if prev == nil {
		t.Error("最新文章应有上一篇")
	}
	if next != nil {
		t.Errorf("最新文章不应有下一篇, got %+v", next)
	}
}

func TestGetRelatedPosts(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	// P2(mysql) 的相关文章：同为 mysql 标签的 P5
	related, err := s.GetRelatedPosts(2, []string{"mysql"}, 6)
	if err != nil {
		t.Fatalf("GetRelatedPosts 报错: %v", err)
	}
	if len(related) != 1 || related[0].Title != "Redis 缓存实践" {
		t.Errorf("related = %+v, want [Redis 缓存实践]", related)
	}

	// 无标签时返回空
	empty, err := s.GetRelatedPosts(2, nil, 6)
	if err != nil || len(empty) != 0 {
		t.Errorf("无标签 related = %v (err=%v), want 空", empty, err)
	}
}

func TestGetArchivePostsPaged(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	// 12 篇已发布文章分布在 5 个年月：2024-05(7) 2024-04(1) 2024-03(2) 2024-02(1) 2024-01(1)
	archives, totalGroups, totalPosts, err := s.GetArchivePostsPaged(1, 3)
	if err != nil {
		t.Fatalf("GetArchivePostsPaged 报错: %v", err)
	}
	if totalGroups != 5 || totalPosts != 12 {
		t.Errorf("totalGroups=%d totalPosts=%d, want 5/12", totalGroups, totalPosts)
	}
	if len(archives.Archives) != 3 {
		t.Fatalf("第 1 页分组数 = %d, want 3", len(archives.Archives))
	}
	if got := len(archives.Archives["2024-05"]); got != 7 {
		t.Errorf("2024-05 分组文章数 = %d, want 7", got)
	}
	// 分组按时间倒序，第 1 页应包含最新的 2024-05
	if _, ok := archives.Archives["2024-01"]; ok {
		t.Error("第 1 页不应包含最旧的 2024-01 分组")
	}

	page2, _, _, err := s.GetArchivePostsPaged(2, 3)
	if err != nil {
		t.Fatalf("GetArchivePostsPaged(2) 报错: %v", err)
	}
	if len(page2.Archives) != 2 {
		t.Errorf("第 2 页分组数 = %d, want 2", len(page2.Archives))
	}
}

func TestGetArchivePosts_SpecificMonth(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	archives, err := s.GetArchivePosts("2024", "03")
	if err != nil {
		t.Fatalf("GetArchivePosts 报错: %v", err)
	}
	if got := len(archives.Archives["2024-03"]); got != 2 {
		t.Errorf("2024-03 文章数 = %d, want 2", got)
	}
}

func TestSearchPaged_RelevanceAndHighlight(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	results, err := s.SearchPaged("mysql", 1, 20)
	if err != nil {
		t.Fatalf("SearchPaged 报错: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("搜索 mysql 应有结果")
	}
	top := results[0]
	if top.Title != "<mark>MySQL</mark> 统计实战" {
		t.Errorf("标题高亮 = %q", top.Title)
	}
	// 标题+摘要+内容均命中的相关度应为 10+5+1=16，排在首位
	if top.Relevance != 16 {
		t.Errorf("Relevance = %d, want 16", top.Relevance)
	}
	if !strings.Contains(top.Highlight, "<mark>") {
		t.Errorf("Highlight 应包含 <mark>: %q", top.Highlight)
	}
	// 相关度非升序
	for i := 1; i < len(results); i++ {
		if results[i].Relevance > results[i-1].Relevance {
			t.Error("结果应按相关度降序")
		}
	}
}

func TestSearch_EmptyKeyword(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	posts, err := s.Search("")
	if err != nil || len(posts) != 0 {
		t.Errorf("空关键词应返回空结果, got %v (err=%v)", posts, err)
	}
}

func TestSearchWithResult(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewPostService()

	result, err := s.SearchWithResult("mysql", 1, 10)
	if err != nil {
		t.Fatalf("SearchWithResult 报错: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("Total = %d, want 1", result.Total)
	}
	if result.Keyword != "mysql" || result.PageNum != 1 || result.PageSize != 10 {
		t.Errorf("回显字段错误: %+v", result)
	}
}

func TestHighlightKeyword(t *testing.T) {
	cases := []struct {
		name, text, keyword, want string
	}{
		{"大小写不敏感", "hello World", "world", "hello <mark>World</mark>"},
		{"多次出现", "go and Go and GO", "go", "<mark>go</mark> and <mark>Go</mark> and <mark>GO</mark>"},
		{"空关键词原样返回", "text", "", "text"},
		{"正则元字符安全", "a.b axb", "a.b", "<mark>a.b</mark> axb"},
		{"中文", "并发编程", "并发", "<mark>并发</mark>编程"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := highlightKeyword(c.text, c.keyword); got != c.want {
				t.Errorf("highlightKeyword(%q,%q) = %q, want %q", c.text, c.keyword, got, c.want)
			}
		})
	}
}

func TestExtractHighlight(t *testing.T) {
	if got := extractHighlight("", "kw", 10); got != "" {
		t.Errorf("空内容应返回空, got %q", got)
	}
	if got := extractHighlight("content", "", 10); got != "" {
		t.Errorf("空关键词应返回空, got %q", got)
	}
	if got := extractHighlight("没有命中词的内容", "kw", 10); got != "" {
		t.Errorf("未命中应返回空, got %q", got)
	}

	long := strings.Repeat("x", 200) + "keyword" + strings.Repeat("y", 200)
	got := extractHighlight(long, "keyword", 60)
	if !strings.Contains(got, "<mark>keyword</mark>") {
		t.Errorf("片段应包含高亮关键词: %q", got)
	}
	if !strings.HasPrefix(got, "...") || !strings.HasSuffix(got, "...") {
		t.Errorf("长文两侧应带省略号: %q", got)
	}

	// 关键词在开头：只有尾部省略号
	head := "keyword" + strings.Repeat("z", 200)
	got = extractHighlight(head, "keyword", 60)
	if strings.HasPrefix(got, "...") {
		t.Errorf("片段起点在开头不应有前导省略号: %q", got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("片段截断应有尾部省略号: %q", got)
	}
}

func TestQueryBuilder_Integration(t *testing.T) {
	testutil.NewTestDB(t)

	var total int64
	if err := NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		Count(&total); err != nil {
		t.Fatalf("Count 报错: %v", err)
	}
	if total != 12 {
		t.Errorf("Count = %d, want 12", total)
	}

	// ByCategoryId 不带发布过滤（后台语义）
	var posts []struct{ PostId uint64 }
	if err := NewQueryBuilder().
		ByCategoryId(1).
		Find(&posts); err != nil {
		t.Fatalf("Find 报错: %v", err)
	}
	// cat1: P1,P3,P4 + P6..P12(7) + P90 + P91 = 12
	if len(posts) != 12 {
		t.Errorf("cat1 全部文章 = %d, want 12", len(posts))
	}
}

func TestPaginationScope_EdgeValues(t *testing.T) {
	// pageNum<=0 归一为 1；pageSize<=0 归一为 DefaultPageSize
	scope := PaginationScope(0, 0)
	if scope == nil {
		t.Fatal("PaginationScope(0,0) 返回 nil")
	}

	var total int64
	err := NewQueryBuilder().
		WithPublished().
		WithNotDeleted().
		WithPagination(0, 0).
		Count(&total)
	if err != nil {
		t.Fatalf("Count 报错: %v", err)
	}
	if total != 12 {
		t.Errorf("边界分页下 Count = %d, want 12（Offset/Limit 不影响 Count）", total)
	}
}
