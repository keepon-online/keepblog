package music

import (
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

// GetMusicList get enabled music list for frontend
func (h *Handler) GetMusicList(c *gin.Context) {
	musicList, err := h.Service.MusicService.GetEnabledMusicList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, musicList)
}
