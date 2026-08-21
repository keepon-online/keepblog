package core

import (
	"gitee.com/jieepre/keepblog/internal/service"
	"github.com/gin-gonic/gin"
)

type Context struct {
	Engine  *gin.Engine
	Service *service.AppService
}

type Handler struct {
	*Context
}
