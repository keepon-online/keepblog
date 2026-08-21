package website

import (
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/keepblog/internal/model/system"
)

type Service struct {
	db *gorm.DB
}

func NewWebSiteService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (service *Service) Save(info system.WebSite) error {
	if err := service.db.Save(&info).Error; err != nil {
		slog.Errorf("save web info error %s", err.Error())
		return errors.New("save error ")
	}
	return nil
}

func (service *Service) GetWebSite() (info *system.WebSite, err error) {
	if err := service.db.Model(system.WebSite{}).First(&info).Error; err != nil {
		slog.Errorf("get web info error: %s", err.Error())
		return nil, errors.New("get info error ")
	}
	return
}
