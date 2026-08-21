package post

import (
	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/model/request"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg"
	"gitee.com/jieepre/keepblog/pkg/hash"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) SavePost(c *gin.Context) {
	var postObj model.Post
	err := c.ShouldBindJSON(&postObj)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	keywords := postObj.Tags
	postObj.Author = "佚名"
	postObj.CoverImage = pkg.GetPixabayImage()
	id, err := h.Service.PostService.SavePost(postObj)
	for _, key := range keywords {
		tag := model.Tag{TagName: key}
		tagId, _ := h.Service.TagService.Save(tag)
		postTag := model.PostTag{
			TagId:  tagId,
			PostId: id,
		}
		_ = global.GORM.Table(model.TPostTagTable).Create(&postTag)
	}
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdatePost(c *gin.Context) {
	var postObj model.Post
	err := c.ShouldBindJSON(&postObj)
	if err != nil {
		result.Error(c, "参数错误")
		return
	}
	err = h.Service.PostService.UpdatePost(postObj)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	_ = global.GORM.Table(model.TPostTagTable).Where("post_id", postObj.PostId).Delete(model.PostTag{})
	keywords := postObj.Tags
	for _, key := range keywords {
		tag := model.Tag{TagName: key}
		tagId, _ := h.Service.TagService.Save(tag)
		postTag := model.PostTag{
			TagId:  tagId,
			PostId: postObj.PostId,
		}
		_ = global.GORM.Table(model.TPostTagTable).Create(&postTag)
	}

	result.Ok(c, nil)
}
func (h *Handler) PublishPost(c *gin.Context) {
	var postObj model.Post
	err := c.ShouldBindJSON(&postObj)
	if err != nil {
		result.Error(c, "参数错误")
		return
	}
	err = h.Service.PostService.PublishArticle(postObj)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) TopPost(c *gin.Context) {
	var postObj request.TopRequest
	err := c.ShouldBindJSON(&postObj)
	if err != nil {
		result.Error(c, "参数错误")
		return
	}
	err = h.Service.PostService.TopPost(postObj.PostId, postObj.Top)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) GetList(c *gin.Context) {
	var req request.PostRequest
	err := c.ShouldBind(&req)
	if err != nil {
		result.Error(c, "参数错误")
		return
	}
	posts, err := h.Service.PostService.GetList(req)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, posts)
}
func (h *Handler) DetailPost(c *gin.Context) {
	param := c.Param("postId")
	ids, _ := hash.New().HashidsDecode(param)
	post, err := h.Service.PostService.GetPostDetail(ids[0])
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, post)
}
func (h *Handler) DeletePost(c *gin.Context) {
	param := c.Param("postId")
	ids, err := hash.New().HashidsDecode(param)
	if err != nil {
		result.Error(c, "参数错误")
		return
	}
	err = h.Service.PostService.DeletePost(ids[0])
	if err != nil {
		result.Error(c, "删除失败")
		return
	}
	result.Ok(c, "ok")
}

func (h *Handler) UpdatePostCoverImag(c *gin.Context) {
	param := c.Param("postId")
	ids, err := hash.New().HashidsDecode(param)
	if err != nil {
		result.Error(c, "参数错误")
		return
	}
	err = h.Service.PostService.UpdatePostCoverImag(ids[0])
	if err != nil {
		result.Error(c, "更新失败,请稍后重试")
		return
	}
	result.Ok(c, "ok")
}

func (h *Handler) UpdatePostAllCoverImag(c *gin.Context) {
	err := h.Service.PostService.UpdatePostAllCoverImag()
	if err != nil {
		result.Error(c, "更新失败,请稍后重试")
		return
	}
	result.Ok(c, "ok")
}
