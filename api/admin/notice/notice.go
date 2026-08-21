package notice

import (
	"strconv"

	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	noticeService "gitee.com/jieepre/go-site/internal/service/notice"
	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

// Handler 通知 Handler
type Handler struct {
	*core.Context
}

// GetNotices 获取用户通知列表
func (h *Handler) GetNotices(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		result.Error(c, "请重新登录")
		return
	}

	service := noticeService.NewNoticeService(global.GORM)
	notices, unreadCount, err := service.GetUserNotices(userID)
	if err != nil {
		result.Error(c, "获取通知失败")
		return
	}

	result.Ok(c, gin.H{
		"notices":     notices,
		"unreadCount": unreadCount,
	})
}

// MarkAsRead 标记通知为已读
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		result.Error(c, "请重新登录")
		return
	}

	noticeIDStr := c.Param("id")
	noticeID, err := strconv.ParseInt(noticeIDStr, 10, 64)
	if err != nil {
		result.Error(c, "无效的通知ID")
		return
	}

	service := noticeService.NewNoticeService(global.GORM)
	if err := service.MarkAsRead(noticeID, userID); err != nil {
		result.Error(c, "标记失败")
		return
	}

	result.Ok(c, nil)
}

// MarkAllAsRead 标记所有通知为已读
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		result.Error(c, "请重新登录")
		return
	}

	var req struct {
		Type *uint8 `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, err.Error())
		return
	}

	var noticeType *system.NoticeType
	if req.Type != nil {
		t := system.NoticeType(*req.Type)
		noticeType = &t
	}

	service := noticeService.NewNoticeService(global.GORM)
	if err := service.MarkAllAsRead(userID, noticeType); err != nil {
		result.Error(c, "标记失败")
		return
	}

	result.Ok(c, nil)
}

// DeleteNotice 删除通知
func (h *Handler) DeleteNotice(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		result.Error(c, "请重新登录")
		return
	}

	noticeIDStr := c.Param("id")
	noticeID, err := strconv.ParseInt(noticeIDStr, 10, 64)
	if err != nil {
		result.Error(c, "无效的通知ID")
		return
	}

	service := noticeService.NewNoticeService(global.GORM)
	if err := service.DeleteNotice(noticeID, userID); err != nil {
		result.Error(c, "删除失败")
		return
	}

	result.Ok(c, nil)
}

// ClearNotices 清空通知
func (h *Handler) ClearNotices(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		result.Error(c, "请重新登录")
		return
	}

	var req struct {
		Type *uint8 `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, err.Error())
		return
	}

	var noticeType *system.NoticeType
	if req.Type != nil {
		t := system.NoticeType(*req.Type)
		noticeType = &t
	}

	service := noticeService.NewNoticeService(global.GORM)
	if err := service.ClearNotices(userID, noticeType); err != nil {
		result.Error(c, "清空失败")
		return
	}

	result.Ok(c, nil)
}

// SendNotice 发送通知
func (h *Handler) SendNotice(c *gin.Context) {
	// 获取当前用户ID
	currentUserID := h.getUserID(c)
	if currentUserID == 0 {
		result.Error(c, "请重新登录")
		return
	}

	var req struct {
		Type        uint8  `json:"type" binding:"required,min=1,max=3"`
		Title       string `json:"title" binding:"required,max=200"`
		Description string `json:"description" binding:"required,max=500"`
		UserId      int64  `json:"userId"` // 0 表示发给自己
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, "请求参数错误: "+err.Error())
		return
	}

	// 如果 userId 为 0，则发给当前用户自己
	targetUserId := req.UserId
	if targetUserId == 0 {
		targetUserId = currentUserID
	}

	notice := &system.Notice{
		UserId:      targetUserId,
		Type:        system.NoticeType(req.Type),
		Title:       req.Title,
		Description: req.Description,
	}

	service := noticeService.NewNoticeService(global.GORM)

	// 保存到数据库并推送
	if err := service.CreateNotice(notice); err != nil {
		result.Error(c, "发送失败: "+err.Error())
		return
	}

	result.Ok(c, nil)
}

// getUserID 从 Token 获取用户ID
func (h *Handler) getUserID(c *gin.Context) int64 {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return 0
	}
	tokenStr := authorization[7:] // 移除 "Bearer "
	claims, err := jwttoken.ParseToken(tokenStr)
	if err != nil {
		return 0
	}

	// 从数据库获取用户ID
	var user model.User
	if err := global.GORM.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		return 0
	}
	return user.UserId
}
