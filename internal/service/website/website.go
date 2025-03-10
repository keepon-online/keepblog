package website

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model/system"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewWebSiteService() *Service {
	return &Service{}
}

func (service *Service) Save(info system.WebSite) error {
	if err := global.GORM.Save(&info).Error; err != nil {
		slog.Errorf("save web info error %s", err.Error())
		return errors.New("save error ")
	}
	return nil
}

func (service *Service) GetWebSite() (info *system.WebSite, err error) {
	if err := global.GORM.Model(system.WebSite{}).First(&info).Error; err != nil {
		slog.Errorf("get web info error: %s", err.Error())
		return nil, errors.New("get info error ")
	}
	return
}
