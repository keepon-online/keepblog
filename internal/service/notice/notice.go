package notice

import (
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/go-site/internal/model/system"
	ws "gitee.com/jieepre/go-site/internal/websocket"
)

// Service 通知服务
type Service struct {
	db *gorm.DB
}

// NewNoticeService 创建通知服务
func NewNoticeService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// NoticeListResponse 通知列表响应
type NoticeListResponse struct {
	Key       string       `json:"key"`
	Name      string       `json:"name"`
	List      []NoticeItem `json:"list"`
	EmptyText string       `json:"emptyText"`
}

// NoticeItem 通知项
type NoticeItem struct {
	Id          int64  `json:"id"`
	Avatar      string `json:"avatar"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Datetime    string `json:"datetime"`
	Type        string `json:"type"`
	Extra       string `json:"extra,omitempty"`
	Status      string `json:"status,omitempty"`
}

// GetUserNotices 获取用户通知（分类展示）
func (s *Service) GetUserNotices(userID int64) ([]NoticeListResponse, int, error) {
	var notices []system.Notice
	if err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notices).Error; err != nil {
		return nil, 0, err
	}

	// 按类型分类
	typeMap := map[system.NoticeType]*NoticeListResponse{
		system.NoticeTypeSystem: {
			Key:       "1",
			Name:      "通知",
			List:      []NoticeItem{},
			EmptyText: "暂无通知",
		},
		system.NoticeTypeMessage: {
			Key:       "2",
			Name:      "消息",
			List:      []NoticeItem{},
			EmptyText: "暂无消息",
		},
		system.NoticeTypeTodo: {
			Key:       "3",
			Name:      "待办",
			List:      []NoticeItem{},
			EmptyText: "暂无待办",
		},
	}

	unreadCount := 0
	for _, n := range notices {
		if n.Status == system.NoticeStatusUnread {
			unreadCount++
		}
		item := NoticeItem{
			Id:          n.Id,
			Avatar:      n.Avatar,
			Title:       n.Title,
			Description: n.Description,
			Datetime:    n.CreatedAt.Format("2006-01-02 15:04"),
			Type:        string(rune('0' + n.Type)),
			Extra:       n.Extra,
			Status:      n.ExtraStatus,
		}
		if resp, ok := typeMap[n.Type]; ok {
			resp.List = append(resp.List, item)
		}
	}

	result := []NoticeListResponse{
		*typeMap[system.NoticeTypeSystem],
		*typeMap[system.NoticeTypeMessage],
		*typeMap[system.NoticeTypeTodo],
	}

	return result, unreadCount, nil
}

// CreateNotice 创建通知
func (s *Service) CreateNotice(notice *system.Notice) error {
	if err := s.db.Create(notice).Error; err != nil {
		slog.Errorf("Failed to create notice: %v", err)
		return err
	}

	// 通过 WebSocket 推送通知
	hub := ws.GetHub()
	if hub.IsUserOnline(notice.UserId) {
		hub.SendToUser(notice.UserId, ws.Message{
			Type:    "notice",
			Payload: notice,
		})
	}

	return nil
}

// MarkAsRead 标记通知为已读
func (s *Service) MarkAsRead(noticeID int64, userID int64) error {
	return s.db.Model(&system.Notice{}).
		Where("id = ? AND user_id = ?", noticeID, userID).
		Update("status", system.NoticeStatusRead).Error
}

// MarkAllAsRead 标记所有通知为已读
func (s *Service) MarkAllAsRead(userID int64, noticeType *system.NoticeType) error {
	query := s.db.Model(&system.Notice{}).Where("user_id = ?", userID)
	if noticeType != nil {
		query = query.Where("type = ?", *noticeType)
	}
	return query.Update("status", system.NoticeStatusRead).Error
}

// DeleteNotice 删除通知
func (s *Service) DeleteNotice(noticeID int64, userID int64) error {
	return s.db.Where("id = ? AND user_id = ?", noticeID, userID).
		Delete(&system.Notice{}).Error
}

// ClearNotices 清空通知
func (s *Service) ClearNotices(userID int64, noticeType *system.NoticeType) error {
	query := s.db.Where("user_id = ?", userID)
	if noticeType != nil {
		query = query.Where("type = ?", *noticeType)
	}
	return query.Delete(&system.Notice{}).Error
}

// BroadcastSystemNotice 广播系统通知给所有在线用户
func (s *Service) BroadcastSystemNotice(title, description string) {
	hub := ws.GetHub()
	hub.Broadcast(ws.Message{
		Type: "system_notice",
		Payload: map[string]string{
			"title":       title,
			"description": description,
		},
	})
}
