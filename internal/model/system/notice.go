package system

import "time"

// NoticeType 通知类型
type NoticeType uint8

const (
	NoticeTypeSystem  NoticeType = 1 // 系统通知
	NoticeTypeMessage NoticeType = 2 // 消息
	NoticeTypeTodo    NoticeType = 3 // 待办
)

// NoticeStatus 通知状态
type NoticeStatus uint8

const (
	NoticeStatusUnread NoticeStatus = 0 // 未读
	NoticeStatusRead   NoticeStatus = 1 // 已读
)

// Notice 通知实体
type Notice struct {
	Id          int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId      int64        `json:"userId" gorm:"index;comment:接收用户ID"`
	Type        NoticeType   `json:"type" gorm:"default:1;comment:通知类型 1系统 2消息 3待办"`
	Title       string       `json:"title" gorm:"size:200;comment:标题"`
	Description string       `json:"description" gorm:"size:500;comment:描述内容"`
	Avatar      string       `json:"avatar" gorm:"size:500;comment:头像URL"`
	Extra       string       `json:"extra" gorm:"size:100;comment:额外标签"`
	ExtraStatus string       `json:"extraStatus" gorm:"size:20;comment:标签状态 primary/success/warning/danger/info"`
	Status      NoticeStatus `json:"status" gorm:"default:0;comment:状态 0未读 1已读"`
	Link        string       `json:"link" gorm:"size:500;comment:跳转链接"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (Notice) TableName() string {
	return "system_notice"
}
