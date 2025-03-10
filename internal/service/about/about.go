package about

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewAboutService() *Service {
	return &Service{}
}

func (service Service) Save(info model.About) error {
	if err := global.GORM.Table(model.TAboutTable).Save(&info).Error; err != nil {
		slog.Errorf("save about info error %s", err.Error())
		return errors.New("save about error ")
	}
	return nil
}

func (service Service) Info() (*model.About, error) {
	var aboutInfo model.About
	if err := global.GORM.Table(model.TAboutTable).First(&aboutInfo).Error; err != nil {
		slog.Errorf("get about info error: %s", err.Error())
		return nil, errors.New("get info error ")
	}
	return &aboutInfo, nil
}
