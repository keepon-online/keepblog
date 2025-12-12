package dashboard

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
)

type Service struct {
}

func NewDashboardService() *Service {
	return &Service{}
}

func (s Service) DashboardData() *model.DashboardData {

	return &model.DashboardData{
		Panel: s.PanelGroup(),
		Pie:   s.Pie(),
		Bar:   s.Bar(),
		Line:  s.Line(),
	}
}

func (s Service) PanelGroup() model.PanelGroup {

	var categoryTotal int64
	var tagTotal int64
	var postTotal int64
	var visit int64

	global.GORM.Table(model.TCategoryTable).Count(&categoryTotal)
	global.GORM.Table(model.TPostsTable).Where("is_published=1 and is_deleted=0").Count(&postTotal)
	global.GORM.Model(system.AccessLog{}).Select("COUNT(DISTINCT ip )").Scan(&visit)
	global.GORM.Model(model.Tag{}).Select("COUNT(DISTINCT tag_name )").Scan(&tagTotal)

	return model.PanelGroup{
		CategoryTotal: uint(categoryTotal),
		TagTotal:      uint(tagTotal),
		PostTotal:     uint(postTotal),
		Visit:         uint(visit),
	}

}

func (s Service) Pie() []map[string]any {
	pie := make([]map[string]any, 0)
	global.GORM.Table(model.TPostsTable).Select("c.category_name `name`,COUNT( post.category_id ) `value` ").
		Joins("left join  category c ON post.category_id = c.category_id ").
		Where("post.is_published = 1 AND post.is_deleted = 0 ").Group("post.category_id").Scan(&pie)

	return pie
}

func (s Service) Bar() []map[string]any {
	pie := make([]map[string]any, 0)
	global.GORM.Model(system.AccessLog{}).
		Select("strftime('%m-%d',create_at,'unixepoch') `name`,SUM(pv) `value`").
		Where("status = 200 AND url LIKE '%/post/%' AND create_at >= strftime('%s', 'now', '-7 days')").Group("strftime('%m-%d',create_at,'unixepoch')").
		Find(&pie)
	return pie
}

func (s Service) Line() []map[string]any {
	pie := make([]map[string]any, 0)
	global.GORM.Model(system.AccessLog{}).
		Select("strftime('%m-%d',create_at,'unixepoch') `name`,SUM(pv) `pv`,SUM(uv) `uv`").
		Where("status = 200 AND url LIKE '%/post/%' AND create_at >= strftime('%s', 'now', '-7 days')").Group("strftime('%m-%d',create_at,'unixepoch')").
		Find(&pie)
	return pie
}
