package model

type PanelGroup struct {
	// 文章
	PostTotal uint `json:"postTotal"`
	//分类
	CategoryTotal uint `json:"categoryTotal"`
	//标签
	TagTotal uint `json:"tagTotal"`
	//访问
	Visit uint `json:"visit"`
	// 总字数
	TotalWords uint `json:"totalWords"`
	// 总阅读量
	TotalReadCount uint `json:"totalReadCount"`
	// 今日访问
	TodayVisit uint `json:"todayVisit"`
	// 音乐数量
	TotalMusic uint `json:"totalMusic"`
}

type DashboardData struct {
	Panel PanelGroup       `json:"panel"`
	Pie   []map[string]any `json:"pie"`
	Bar   []map[string]any `json:"bar"`
	Line  []map[string]any `json:"line"`
	Map   []map[string]any `json:"map"`
}
