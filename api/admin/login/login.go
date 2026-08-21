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
	result.Ok(c, make([]response.Routes, 0))
}
func (h *Handler) GetInfo(c *gin.Context) {
	// 从 Token 获取用户名
	authorization := c.GetHeader("Authorization")
	tokenStr := strings.ReplaceAll(authorization, "Bearer ", "")
	token, err := jwttoken.ParseToken(tokenStr)
	if err != nil {
		result.Error(c, "请重新登录")
		return
	}

	userInfo, err := h.Service.SystemService.GetUserInfo(token.Username)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, userInfo)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	// 从 Token 获取用户名
	authorization := c.GetHeader("Authorization")
	tokenStr := strings.ReplaceAll(authorization, "Bearer ", "")
	token, err := jwttoken.ParseToken(tokenStr)
	if err != nil {
		result.Error(c, "请重新登录")
		return
	}

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, "请求参数错误")
		return
	}

	if err := h.Service.SystemService.UpdateProfile(token.Username, req); err != nil {
		result.Error(c, err.Error())
		return
	}
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
