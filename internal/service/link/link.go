package link

import (
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/go-site/internal/model"
)

// Service 友链服务
type Service struct {
	db *gorm.DB
}

// NewLinkService 创建友链服务实例
func NewLinkService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Save 保存友链
func (service Service) Save(info model.FriendLink) error {
	if err := service.db.Table(model.TFriendLinkTable).Create(&info).Error; err != nil {
		slog.Errorf("保存友链失败: %s", err.Error())
		return errors.New("保存友链失败: " + err.Error())
	}
	return nil
}

// Update 更新友链
func (service Service) Update(info model.FriendLink) error {
	if err := service.db.Table(model.TFriendLinkTable).Save(&info).Error; err != nil {
		slog.Errorf("更新友链失败: %s", err.Error())
		return errors.New("更新友链失败: " + err.Error())
	}
	return nil
}

// Delete 删除友链
func (service Service) Delete(id int) error {
	if err := service.db.Table(model.TFriendLinkTable).Delete(&model.FriendLink{}, id).Error; err != nil {
		slog.Errorf("删除友链失败: %s", err.Error())
		return errors.New("删除友链失败: " + err.Error())
	}
	return nil
}

// GetLink 获取单个友链
func (service Service) GetLink(id uint32) (*model.FriendLink, error) {
	var linkInfo model.FriendLink

	err := NewQueryBuilder(service.db).
		ById(id).
		First(&linkInfo)

	if err != nil {
		slog.Errorf("获取友链失败: %s", err.Error())
		return nil, errors.New("获取友链失败")
	}

	return &linkInfo, nil
}

// GetLinkList 获取所有友链列表（后台管理）
func (service Service) GetLinkList() ([]model.FriendLink, error) {
	links := make([]model.FriendLink, 0)

	err := NewQueryBuilder(service.db).
		WithOrderByTime(true).
		Find(&links)

	if err != nil {
		slog.Errorf("获取友链列表失败: %s", err.Error())
		return nil, errors.New("获取友链列表失败")
	}

	return links, nil
}

// GetLinks 获取启用状态的友链（前台展示）
func (service Service) GetLinks() ([]model.FriendLink, error) {
	links := make([]model.FriendLink, 0)

	err := NewQueryBuilder(service.db).
		WithActive().
		WithOrderByTime(true).
		Find(&links)

	if err != nil {
		slog.Errorf("获取友链失败: %s", err.Error())
		return nil, errors.New("获取友链失败")
	}

	return links, nil
}
