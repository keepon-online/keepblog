package category

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

type Service struct {
}

func NewCategoryService() *Service {
	return &Service{}
}

func (service Service) Save(info model.Category) error {
	if err := global.GORM.Table(model.TCategoryTable).Create(&info).Error; err != nil {
		slog.Errorf("save category error %s", err.Error())
		return errors.New("save category error " + err.Error())
	}
	return nil
}

func (service Service) Update(info model.Category) error {
	if err := global.GORM.Table(model.TCategoryTable).Save(&info).Error; err != nil {
		slog.Errorf("update category error %s", err.Error())
		return errors.New("update category error " + err.Error())
	}
	return nil
}

func (service Service) Delete(categoryId int) error {
	if err := global.GORM.Table(model.TCategoryTable).Delete(&model.Category{}, categoryId).Error; err != nil {
		slog.Errorf("delete category error %s", err.Error())
		return errors.New("delete category error " + err.Error())
	}
	return nil
}

func (service Service) GetCategory(tagId uint32) (*model.Category, error) {
	var categoryInfo model.Category
	if err := global.GORM.Table(model.TCategoryTable).Where("category_id", tagId).First(&categoryInfo).Error; err != nil {
		slog.Errorf("get category error %s", err.Error())
		return nil, errors.New("get category error")
	}
	return &categoryInfo, nil
}

func (service Service) GetCategoryList() ([]model.Category, error) {
	categories := make([]model.Category, 0)
	if err := global.GORM.Table(model.TCategoryTable).Find(&categories).Error; err != nil {
		slog.Errorf("get categories error %s", err.Error())
		return nil, errors.New("get categories error")
	}
	return categories, nil
}

func (service Service) GetCategories() ([]model.CategoryCount, error) {
	categories := make([]model.CategoryCount, 0)
	tx := global.GORM.Table(model.TPostsTable).
		Select("category.category_name,COUNT(category.category_id) total").
		Joins("LEFT JOIN category ON post.category_id = category.category_id").
		Where("post.is_published", 1).
		Where("post.is_deleted", 0).
		Group("category.category_name").
		Find(&categories)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return categories, nil
}
