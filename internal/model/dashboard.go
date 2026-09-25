package model

type PanelGroup struct {
	// 文章
	PostTotal uint `json:"postTotal"`
	// 分类
	CategoryTotal uint `json:"categoryTotal"`
	// 标签
	TagTotal uint `json:"tagTotal"`
	// 访问
	Visit uint `json:"visit"`
	// 总字数
	TotalWords uint `json:"totalWords"`
	// 总阅读量
	TotalReadCount uint `json:"totalReadCount"`
	// 今日访问
	TodayVisit uint `json:"todayVisit"`
	// 昨日访问
	YesterdayVisit uint `json:"yesterdayVisit"`
	// 今日访问环比增长率(%)
	VisitGrowth float64 `json:"visitGrowth"`
	// 本周发布文章数
	WeekPostTotal uint `json:"weekPostTotal"`
	// 音乐数量
	TotalMusic uint `json:"totalMusic"`
}

type TopPostItem struct {
	ID           uint64 `json:"id"`
	Title        string `json:"title"`
	ReadCount    uint32 `json:"readCount"`
	CategoryName string `json:"categoryName"`
	CreateTime   uint64 `json:"createTime"`
}

type DashboardData struct {
	Panel    PanelGroup       `json:"panel"`
	Pie      []map[string]any `json:"pie"`
	Bar      []map[string]any `json:"bar"`
	Line     []map[string]any `json:"line"`
	Map      []map[string]any `json:"map"`
	TopPosts []TopPostItem    `json:"topPosts"`
	Days     int              `json:"days"`
}

