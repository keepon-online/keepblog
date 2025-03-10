package router

import (
	"github.com/gin-gonic/gin"
)

// RouteGroup 路由分组定义
type RouteGroup struct {
	Name   string // 路由名称
	Prefix string // 路由前缀
}

// Route 路由定义
type Route struct {
	Path       string            // 路由
	Method     string            // 请求方法
	Middleware []gin.HandlerFunc // 路由中间件
	Handler    gin.HandlerFunc   //请求接口
}

func RegisterRouter(group *gin.RouterGroup, routes []Route) {
	for _, route := range routes {
		handlers := append(route.Middleware, route.Handler)
		group.Handle(route.Method, route.Path, handlers...)
	}
}
