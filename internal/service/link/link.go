package link

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewLinkService() *Service {
	return &Service{}
}

func (service Service) Save(info model.FriendLink) error {
	if err := global.GORM.Table(model.TFriendLinkTable).Create(&info).Error; err != nil {
		slog.Errorf("save link error %s", err.Error())
		return errors.New("save link error " + err.Error())
	}
	return nil
}

func (service Service) Update(info model.FriendLink) error {
	if err := global.GORM.Table(model.TFriendLinkTable).Save(&info).Error; err != nil {
		slog.Errorf("update link error %s", err.Error())
		return errors.New("update link error " + err.Error())
	}
	return nil
}

func (service Service) Delete(id int) error {
	if err := global.GORM.Table(model.TFriendLinkTable).Delete(&model.FriendLink{}, id).Error; err != nil {
		slog.Errorf("delete link error %s", err.Error())
		return errors.New("delete link error " + err.Error())
	}
	return nil
}

func (service Service) GetLink(id uint32) (*model.FriendLink, error) {
	var linkInfo model.FriendLink
	if err := global.GORM.Table(model.TFriendLinkTable).Where("id", id).First(&linkInfo).Error; err != nil {
		slog.Errorf("get link error %s", err.Error())
		return nil, errors.New("get link error")
	}
	return &linkInfo, nil
}

func (service Service) GetLinkList() ([]model.FriendLink, error) {
	links := make([]model.FriendLink, 0)
	if err := global.GORM.Table(model.TFriendLinkTable).Find(&links).Error; err != nil {
		slog.Errorf("get link error %s", err.Error())
		return nil, errors.New("get link error")
	}
	return links, nil
}

func (service Service) GetLinks() ([]model.FriendLink, error) {
	links := make([]model.FriendLink, 0)
	if err := global.GORM.Table(model.TFriendLinkTable).Where("state", 1).Find(&links).Error; err != nil {
		slog.Errorf("get link error %s", err.Error())
		return nil, errors.New("get link error")
	}
	return links, nil

}
