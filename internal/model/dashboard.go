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
}
