package about

import (
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/go-site/internal/model"
)

type Service struct {
	db *gorm.DB
}

func NewAboutService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (service Service) Save(info model.About) error {
	if err := service.db.Table(model.TAboutTable).Save(&info).Error; err != nil {
		slog.Errorf("save about info error %s", err.Error())
		return errors.New("save about error ")
	}
	return nil
}

func (service Service) Info() (*model.About, error) {
	var aboutInfo model.About
	if err := service.db.Table(model.TAboutTable).First(&aboutInfo).Error; err != nil {
		slog.Errorf("get about info error: %s", err.Error())
		return nil, errors.New("get info error ")
	}
	return &aboutInfo, nil
}
