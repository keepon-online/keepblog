package music

import (
	"strconv"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

// SaveMusic 添加音乐
func (h *Handler) SaveMusic(c *gin.Context) {
	music := model.Music{}
	if err := c.ShouldBindJSON(&music); err != nil {
		result.Error(c, "参数错误: "+err.Error())
		return
	}
	if err := h.Service.MusicService.Save(music); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// UpdateMusic 更新音乐
func (h *Handler) UpdateMusic(c *gin.Context) {
	music := model.Music{}
	if err := c.ShouldBindJSON(&music); err != nil {
		result.Error(c, "参数错误: "+err.Error())
		return
	}
	if err := h.Service.MusicService.Update(music); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// UpdateMusicState 更新音乐状态
func (h *Handler) UpdateMusicState(c *gin.Context) {
	var req struct {
		Id    uint32 `json:"id"`
		State uint8  `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, "参数错误: "+err.Error())
		return
	}
	if err := h.Service.MusicService.UpdateState(req.Id, req.State); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// DeleteMusic 删除音乐
func (h *Handler) DeleteMusic(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		result.Error(c, "无效的ID")
		return
	}
	if err := h.Service.MusicService.Delete(int(id)); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// DetailMusic 获取音乐详情
func (h *Handler) DetailMusic(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		result.Error(c, "无效的ID")
		return
	}
	music, err := h.Service.MusicService.GetMusic(uint32(id))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, music)
}

// GetMusicList 获取音乐列表
func (h *Handler) GetMusicList(c *gin.Context) {
	musicList, err := h.Service.MusicService.GetMusicList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, musicList)
}
