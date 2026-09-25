package tag

import (
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/querybuilder"
	"gitee.com/jieepre/keepblog/pkg/cloudtag"
	tags_remove "gitee.com/jieepre/keepblog/pkg/tags-remove"
)

type Service struct {
	db *gorm.DB
}

func NewTagService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (service Service) Save(tagInfo model.Tag) (uint32, error) {
	var t model.Tag
	err := service.db.Table(model.TTagTable).Where("tag_name", tagInfo.TagName).Last(&t).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {

		return t.TagId, nil
	}
	if err := service.db.Table(model.TTagTable).Create(&tagInfo).Error; err != nil {
		return 0, errors.New("save tag error " + err.Error())
	}
	return tagInfo.TagId, nil
}

func (service Service) GetTag(tagId uint32) (*model.Tag, error) {
	var tagInfo model.Tag
	if err := service.db.Table(model.TTagTable).Where("tag_id", tagId).First(&tagInfo).Error; err != nil {
		slog.Errorf("save tag error %s", err.Error())
		return nil, errors.New("save tag error")
	}
	return &tagInfo, nil
}

// Update 更新标签名（后台编辑）。仅改 tag_name，保留 tag_id。
func (service Service) Update(tag model.Tag) error {
	if err := service.db.Table(model.TTagTable).
		Model(&model.Tag{}).
		Where("tag_id = ?", tag.TagId).
		Update("tag_name", tag.TagName).Error; err != nil {
		slog.Errorf("更新标签失败: %s", err.Error())
		return errors.New("更新标签失败: " + err.Error())
	}
	return nil
}

// Delete 删除标签，并清理 post_tag 关联表，避免脏数据。
// 不校验是否被文章引用——标签被删后，文章的 Tags 查询会自动少一项（走 JOIN）。
func (service Service) Delete(tagId uint32) error {
	// 先清理关联表
	if err := service.db.Table(model.TPostTagTable).
		Where("tag_id = ?", tagId).
		Delete(&model.PostTag{}).Error; err != nil {
		slog.Errorf("清理 post_tag 关联失败: %s", err.Error())
		return errors.New("清理标签关联失败: " + err.Error())
	}
	// 再删标签本身
	if err := service.db.Table(model.TTagTable).
		Delete(&model.Tag{}, tagId).Error; err != nil {
		slog.Errorf("删除标签失败: %s", err.Error())
		return errors.New("删除标签失败: " + err.Error())
	}
	return nil
}

// GetTagList 获取全部标签列表（后台管理用，附带文章篇数统计）。
func (service Service) GetTagList() ([]model.Tag, error) {
	tagList := make([]model.Tag, 0)
	err := service.db.Table(model.TTagTable+" as t").
		Select("t.*, COUNT(pt.post_id) as post_count").
		Joins("LEFT JOIN post_tag pt ON t.tag_id = pt.tag_id").
		Group("t.tag_id").
		Order("t.tag_id DESC").
		Find(&tagList).Error
	if err != nil {
		slog.Errorf("获取标签列表失败: %s", err.Error())
		return nil, errors.New("获取标签列表失败")
	}
	return tagList, nil
}

// GetTags 前台 tag cloud 用：返回带样式、去重的标签列表。
func (service Service) GetTags() ([]model.Tag, error) {
	tagList := make([]model.Tag, 0)

	sql := `
		select tag.*
		from tag
				 left join post_tag pt on tag.tag_id = pt.tag_id
				 left join post p on pt.post_id = p.post_id
		where p.is_deleted = 0
		  and ` + querybuilder.VisibleWhere("p") + `
	`

	if err := service.db.Raw(sql, querybuilder.VisibleNow()).
		Scan(&tagList).Error; err != nil {
		return nil, err
	}
	tagCounts := make([]model.TagCount, 0)
	tagsCloud := make([]model.Tag, 0)
	tags := tags_remove.RemoveDuplicateElement(tagList)
	service.db.Raw("SELECT tag_name,COUNT(tag_name) ratio FROM tag GROUP BY tag_name").Scan(&tagCounts)
	for _, tsg := range tagCounts {
		for _, tag := range tags {
			if tag.TagName == tsg.TagName {
				tagsCloud = append(tagsCloud, model.Tag{
					TagId:        tag.TagId,
					CreateTime:   tag.CreateTime,
					TagName:      tag.TagName,
					TagStyle:     cloudtag.CloudTags(1.15, 1.5, tsg.Ratio),
					CardTagStyle: cloudtag.CloudTags(1.1, 1.35, 1),
				})
			}
		}
	}
	return tagsCloud, nil
}
