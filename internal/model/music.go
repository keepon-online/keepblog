package model

const TMusicTable = "music"

// Music 音乐配置模型
type Music struct {
	Id         uint32 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" gorm:"default:''"`   // 歌曲名称
	Artist     string `json:"artist" gorm:"default:''"` // 艺术家/歌手
	Url        string `json:"url" gorm:"default:''"`    // 音乐文件URL
	Cover      string `json:"cover" gorm:"default:''"`  // 封面图片URL
	Lrc        string `json:"lrc" gorm:"default:''"`    // 歌词URL (可选)
	Sort       int    `json:"sort" gorm:"default:0"`    // 排序 (值越小越靠前)
	State      uint8  `json:"state" gorm:"default:1"`   // 状态 1:启用 0:禁用
	CreateTime uint64 `json:"createTime" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (Music) TableName() string {
	return TMusicTable
}
