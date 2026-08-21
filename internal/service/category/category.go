package category

import (
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/keepblog/internal/model"
)

// Service 分类服务
type Service struct {
	db *gorm.DB
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Save 保存分类
func (service Service) Save(info model.Category) error {
	if err := service.db.Table(model.TCategoryTable).Create(&info).Error; err != nil {
		slog.Errorf("保存分类失败: %s", err.Error())
		return errors.New("保存分类失败: " + err.Error())
	}
	return nil
}

// Update 更新分类
func (service Service) Update(info model.Category) error {
	if err := service.db.Table(model.TCategoryTable).Save(&info).Error; err != nil {
		slog.Errorf("更新分类失败: %s", err.Error())
		return errors.New("更新分类失败: " + err.Error())
	}
	return nil
}

// Delete 删除分类
func (service Service) Delete(categoryId int) error {
	if err := service.db.Table(model.TCategoryTable).Delete(&model.Category{}, categoryId).Error; err != nil {
		slog.Errorf("删除分类失败: %s", err.Error())
		return errors.New("删除分类失败: " + err.Error())
	}
	return nil
}

// GetCategory 获取单个分类
func (service Service) GetCategory(categoryId uint32) (*model.Category, error) {
	var categoryInfo model.Category

	err := NewQueryBuilder(service.db).
		ById(categoryId).
		First(&categoryInfo)

	if err != nil {
		slog.Errorf("获取分类失败: %s", err.Error())
		return nil, errors.New("获取分类失败")
	}

	return &categoryInfo, nil
}

// GetCategoryList 获取所有分类列表（后台管理）
func (service Service) GetCategoryList() ([]model.Category, error) {
	categories := make([]model.Category, 0)

	err := NewQueryBuilder(service.db).Find(&categories)

	if err != nil {
		slog.Errorf("获取分类列表失败: %s", err.Error())
		return nil, errors.New("获取分类列表失败")
	}

	return categories, nil
}

// GetCategories 获取分类及文章数量（前台展示）
func (service Service) GetCategories() ([]model.CategoryCount, error) {
	categories := make([]model.CategoryCount, 0)

	err := NewQueryBuilder(service.db).
		WithPostCount().
		Find(&categories)

	if err != nil {
		slog.Errorf("获取分类统计失败: %s", err.Error())
		return nil, errors.New("获取分类统计失败")
	}

	return categories, nil
}
