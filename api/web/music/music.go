package music

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

// GetMusicList 获取启用的音乐列表 (前台接口)
func (h *Handler) GetMusicList(c *gin.Context) {
	musicList, err := h.Service.MusicService.GetEnabledMusicList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, musicList)
}
