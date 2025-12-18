package login

import (
	"strings"

	"gitee.com/jieepre/go-site/internal/model/request"
	"gitee.com/jieepre/go-site/internal/model/response"
	"gitee.com/jieepre/go-site/internal/pkg/core"
	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	loginResponse, err := h.Service.SystemService.Login(req, c)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, loginResponse)
}
func (h *Handler) RefreshToken(c *gin.Context) {
	s := struct {
		RefreshToken string `json:"refreshToken"`
	}{}
	err := c.ShouldBindJSON(&s)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	token, err := h.Service.SystemService.RefreshToken(s.RefreshToken)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, token)
}
func (h *Handler) Logout(c *gin.Context) {
	// 获取当前 Token 并加入黑名单
	authorization := c.GetHeader("Authorization")
	if authorization != "" {
		tokenStr := strings.ReplaceAll(authorization, "Bearer ", "")
		if tokenStr != "" {
			_ = jwttoken.InvalidateToken(tokenStr)
		}
	}

	result.Ok(c, nil)
}
func (h *Handler) GetAsyncRoutes(c *gin.Context) {
	routers := make([]response.Routes, 0)
	c.JSON(200, gin.H{
		"success": true,
		"data":    routers,
	})
}
func (h *Handler) GetInfo(c *gin.Context) {

	result.Ok(c, nil)
}
func (h *Handler) ChangePassword(c *gin.Context) {

	var req request.ChangePasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	Authorization := c.GetHeader("Authorization")
	tokenStr := strings.ReplaceAll(Authorization, "Bearer ", "")
	token, err := jwttoken.ParseToken(tokenStr)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	req.Username = token.Username

	err = h.Service.SystemService.ChangePassword(req)
	if err != nil {
		result.Error(c, err.Error())
		return
	}

	result.Ok(c, nil)
}
