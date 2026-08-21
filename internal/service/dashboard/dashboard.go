package dashboard

import (
	"strings"

	"gorm.io/gorm"

	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
)

type Service struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s Service) DashboardData() *model.DashboardData {

	return &model.DashboardData{
		Panel: s.PanelGroup(),
		Pie:   s.Pie(),
		Bar:   s.Bar(),
		Line:  s.Line(),
		Map:   s.MapData(),
	}
}

func (s Service) PanelGroup() model.PanelGroup {

	var categoryTotal int64
	var tagTotal int64
	var postTotal int64
	var visit int64
	var totalWords int64
	var totalReadCount int64
	var todayVisit int64
	var totalMusic int64

	s.db.Table(model.TCategoryTable).Count(&categoryTotal)
	s.db.Table(model.TPostsTable).Where("is_published=1 and is_deleted=0").Count(&postTotal)
	s.db.Model(system.AccessLog{}).Select("COUNT(DISTINCT ip )").Scan(&visit)
	s.db.Model(model.Tag{}).Select("COUNT(DISTINCT tag_name )").Scan(&tagTotal)

	// 新增统计
	// 总字数
	s.db.Table(model.TPostsTable).
		Where("is_published=1 and is_deleted=0").
		Select("COALESCE(SUM(word_count), 0)").
		Scan(&totalWords)

	// 总阅读量
	s.db.Table(model.TPostsTable).
		Where("is_published=1 and is_deleted=0").
		Select("COALESCE(SUM(read_count), 0)").
		Scan(&totalReadCount)

	// 今日访问（基于今天的访问记录）
	s.db.Model(system.AccessLog{}).
		Where("create_at >= strftime('%s', 'now', 'start of day')").
		Select("COUNT(DISTINCT ip)").
		Scan(&todayVisit)

	// 音乐数量
	s.db.Table("music").Count(&totalMusic)

	return model.PanelGroup{
		CategoryTotal:  uint(categoryTotal),
		TagTotal:       uint(tagTotal),
		PostTotal:      uint(postTotal),
		Visit:          uint(visit),
		TotalWords:     uint(totalWords),
		TotalReadCount: uint(totalReadCount),
		TodayVisit:     uint(todayVisit),
		TotalMusic:     uint(totalMusic),
	}

}

func (s Service) Pie() []map[string]any {
	pie := make([]map[string]any, 0)
	s.db.Table(model.TPostsTable).Select("c.category_name `name`,COUNT( post.category_id ) `value` ").
		Joins("left join  category c ON post.category_id = c.category_id ").
		Where("post.is_published = 1 AND post.is_deleted = 0 ").Group("post.category_id").Scan(&pie)

	return pie
}

func (s Service) Bar() []map[string]any {
	pie := make([]map[string]any, 0)
	s.db.Model(system.AccessLog{}).
		Select("strftime('%m-%d',create_at,'unixepoch') `name`,SUM(pv) `value`").
		Where("status = 200 AND url LIKE '%/post/%' AND create_at >= strftime('%s', 'now', '-7 days')").Group("strftime('%m-%d',create_at,'unixepoch')").
		Find(&pie)
	return pie
}

func (s Service) Line() []map[string]any {
	pie := make([]map[string]any, 0)
	s.db.Model(system.AccessLog{}).
		Select("strftime('%m-%d',create_at,'unixepoch') `name`,SUM(pv) `pv`,SUM(uv) `uv`").
		Where("status = 200 AND url LIKE '%/post/%' AND create_at >= strftime('%s', 'now', '-7 days')").Group("strftime('%m-%d',create_at,'unixepoch')").
		Find(&pie)
	return pie
}

// MapData 获取访客地理分布数据（按省份聚合）
// ip2region格式: 中国|0|江苏省|南京市|电信 或 美国|0|0|0|...
func (s Service) MapData() []map[string]any {
	// 先获取原始数据
	var rawData []struct {
		Area  string
		Count int64
	}

	s.db.Model(system.AccessLog{}).
		Select("area as Area, COUNT(DISTINCT ip) as Count").
		Where("area != '' AND area IS NOT NULL").
		Group("area").
		Find(&rawData)

	// 解析省份并聚合（中国省份和国外国家）
	provinceMap := make(map[string]int64)

	for _, item := range rawData {
		location := parseLocation(item.Area)
		if location != "" && location != "0" {
			provinceMap[location] += item.Count
		}
	}

	// 转换为返回格式
	mapData := make([]map[string]any, 0)
	for name, value := range provinceMap {
		mapData = append(mapData, map[string]any{
			"name":  name,
			"value": value,
		})
	}

	return mapData
}

// parseLocation 从ip2region格式解析地理位置
// 格式: 中国|0|江苏省|南京市|电信 或 美国|0|华盛顿|西雅图|...
// 中国IP返回省份名称，国外IP返回国家名称
func parseLocation(area string) string {
	if area == "" {
		return ""
	}

	parts := strings.Split(area, "|")
	if len(parts) < 1 {
		return area
	}

	country := parts[0]

	// 中国IP：提取省份
	if country == "中国" {
		if len(parts) >= 3 {
			province := parts[2]
			if province == "0" || province == "" {
				return "" // 无效省份数据
			}
			// 移除"省"、"自治区"等后缀
			province = strings.TrimSuffix(province, "省")
			province = strings.TrimSuffix(province, "自治区")
			province = strings.TrimSuffix(province, "特别行政区")
			province = strings.TrimSuffix(province, "市") // 直辖市
			return province
		}
		return ""
	}

	// 国外IP：返回国家名称
	if country != "" && country != "0" {
		return country
	}

	return ""
}
