package music

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewMusicService() *Service {
	return &Service{}
}

// Save 保存音乐
func (s *Service) Save(music model.Music) error {
	if err := global.GORM.Table(model.TMusicTable).Create(&music).Error; err != nil {
		slog.Errorf("save music error: %s", err.Error())
		return errors.New("保存音乐失败: " + err.Error())
	}
	return nil
}

// Update 更新音乐
func (s *Service) Update(music model.Music) error {
	if err := global.GORM.Table(model.TMusicTable).Save(&music).Error; err != nil {
		slog.Errorf("update music error: %s", err.Error())
		return errors.New("更新音乐失败: " + err.Error())
	}
	return nil
}

// UpdateState 更新音乐状态
func (s *Service) UpdateState(id uint32, state uint8) error {
	if err := global.GORM.Table(model.TMusicTable).Where("id = ?", id).Update("state", state).Error; err != nil {
		slog.Errorf("update music state error: %s", err.Error())
		return errors.New("更新音乐状态失败: " + err.Error())
	}
	return nil
}

// Delete 删除音乐
func (s *Service) Delete(id int) error {
	if err := global.GORM.Table(model.TMusicTable).Delete(&model.Music{}, id).Error; err != nil {
		slog.Errorf("delete music error: %s", err.Error())
		return errors.New("删除音乐失败: " + err.Error())
	}
	return nil
}

// GetMusic 获取单个音乐
func (s *Service) GetMusic(id uint32) (*model.Music, error) {
	var music model.Music
	if err := global.GORM.Table(model.TMusicTable).Where("id = ?", id).First(&music).Error; err != nil {
		slog.Errorf("get music error: %s", err.Error())
		return nil, errors.New("获取音乐失败")
	}
	return &music, nil
}

// GetMusicList 获取全部音乐列表 (后台管理)
func (s *Service) GetMusicList() ([]model.Music, error) {
	musicList := make([]model.Music, 0)
	if err := global.GORM.Table(model.TMusicTable).Order("sort ASC, id DESC").Find(&musicList).Error; err != nil {
		slog.Errorf("get music list error: %s", err.Error())
		return nil, errors.New("获取音乐列表失败")
	}
	return musicList, nil
}

// GetEnabledMusicList 获取启用的音乐列表 (前台展示)
func (s *Service) GetEnabledMusicList() ([]model.Music, error) {
	musicList := make([]model.Music, 0)
	if err := global.GORM.Table(model.TMusicTable).Where("state = ?", 1).Order("sort ASC, id DESC").Find(&musicList).Error; err != nil {
		slog.Errorf("get enabled music list error: %s", err.Error())
		return nil, errors.New("获取音乐列表失败")
	}
	return musicList, nil
}
