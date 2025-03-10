package tag

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/pkg/cloudtag"
	tags_remove "gitee.com/jieepre/go-site/pkg/tags-remove"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
}

func NewTagService() *Service {
	return &Service{}
}

func (service Service) Save(tagInfo model.Tag) (error, uint32) {
	var t model.Tag
	err := global.GORM.Table(model.TTagTable).Where("tag_name", tagInfo.TagName).Last(&t).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, t.TagId
	}
	if err := global.GORM.Table(model.TTagTable).Create(&tagInfo).Error; err != nil {
		return errors.New("save tag error " + err.Error()), 0
	}
	return nil, tagInfo.TagId
}

func (service Service) GetTag(tagId uint32) (*model.Tag, error) {
	var tagInfo model.Tag
	if err := global.GORM.Table(model.TTagTable).Where("tag_id", tagId).First(&tagInfo).Error; err != nil {
		slog.Errorf("save tag error %s", err.Error())
		return nil, errors.New("save tag error")
	}
	return &tagInfo, nil
}

func (service Service) GetTags() ([]model.Tag, error) {
	tagList := make([]model.Tag, 0)

	sql := `
		select tag.*
		from tag
				 left join post_tag pt on tag.tag_id = pt.tag_id
				 left join post p on pt.post_id = p.post_id
		where p.is_deleted = 0
		  and p.is_published = 1
	`

	if err := global.GORM.Raw(sql).
		Scan(&tagList).Error; err != nil {
		return nil, err
	}
	tagCounts := make([]model.TagCount, 0)
	tagsCloud := make([]model.Tag, 0)
	tags := tags_remove.RemoveDuplicateElement(tagList)
	global.GORM.Raw("SELECT tag_name,COUNT(tag_name) ratio FROM tag GROUP BY tag_name").Scan(&tagCounts)
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
