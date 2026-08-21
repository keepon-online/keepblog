package tags_remove

import "gitee.com/jieepre/keepblog/internal/model"

func RemoveDuplicateElement(tags []model.Tag) []model.Tag {
	result := make([]model.Tag, 0, len(tags))
	temp := map[string]struct{}{}
	for _, item := range tags {
		if _, ok := temp[item.TagName]; !ok {
			temp[item.TagName] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
