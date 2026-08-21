package music

import (
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/go-site/internal/model"
)

// Service 音乐服务
type Service struct {
	db *gorm.DB
}

// NewMusicService 创建音乐服务实例
func NewMusicService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Save 保存音乐
func (s *Service) Save(music model.Music) error {
	if err := s.db.Table(model.TMusicTable).Create(&music).Error; err != nil {
		slog.Errorf("保存音乐失败: %s", err.Error())
		return errors.New("保存音乐失败: " + err.Error())
	}
	return nil
}

// Update 更新音乐
func (s *Service) Update(music model.Music) error {
	if err := s.db.Table(model.TMusicTable).Save(&music).Error; err != nil {
		slog.Errorf("更新音乐失败: %s", err.Error())
		return errors.New("更新音乐失败: " + err.Error())
	}
	return nil
}

// UpdateState 更新音乐状态
func (s *Service) UpdateState(id uint32, state uint8) error {
	if err := s.db.Table(model.TMusicTable).Where("id = ?", id).Update("state", state).Error; err != nil {
		slog.Errorf("更新音乐状态失败: %s", err.Error())
		return errors.New("更新音乐状态失败: " + err.Error())
	}
	return nil
}

// Delete 删除音乐
func (s *Service) Delete(id int) error {
	if err := s.db.Table(model.TMusicTable).Delete(&model.Music{}, id).Error; err != nil {
		slog.Errorf("删除音乐失败: %s", err.Error())
		return errors.New("删除音乐失败: " + err.Error())
	}
	return nil
}

// GetMusic 获取单个音乐
func (s *Service) GetMusic(id uint32) (*model.Music, error) {
	var music model.Music

	err := NewQueryBuilder(s.db).
		ById(id).
		First(&music)

	if err != nil {
		slog.Errorf("获取音乐失败: %s", err.Error())
		return nil, errors.New("获取音乐失败")
	}

	return &music, nil
}

// GetMusicList 获取全部音乐列表（后台管理）
func (s *Service) GetMusicList() ([]model.Music, error) {
	musicList := make([]model.Music, 0)

	err := NewQueryBuilder(s.db).
		WithOrderBySort().
		Find(&musicList)

	if err != nil {
		slog.Errorf("获取音乐列表失败: %s", err.Error())
		return nil, errors.New("获取音乐列表失败")
	}

	return musicList, nil
}

// GetEnabledMusicList 获取启用的音乐列表（前台展示）
func (s *Service) GetEnabledMusicList() ([]model.Music, error) {
	musicList := make([]model.Music, 0)

	err := NewQueryBuilder(s.db).
		WithActive().
		WithOrderBySort().
		Find(&musicList)

	if err != nil {
		slog.Errorf("获取音乐列表失败: %s", err.Error())
		return nil, errors.New("获取音乐列表失败")
	}

	return musicList, nil
}
