package middleware

import (
	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"gitee.com/jieepre/go-site/pkg/result"
	"github.com/gin-gonic/gin"
	"strings"
)

func JwtVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		ignoreURI := []string{
			"/api/refreshToken",
			"/api/login",
			"/api/v1/upload/images",
		}
		for i := range ignoreURI {
			if c.Request.RequestURI == ignoreURI[i] {
				c.Next()
				return
			}
		}
		Authorization := c.Request.Header.Get("Authorization")
		if Authorization == "" {
			result.Error(c, "TOKEN 不能为空")
			c.Abort()
			return
		}
		token := strings.ReplaceAll(Authorization, "Bearer ", "")
		_, err := jwttoken.ParseToken(token)
		if err != nil {
			result.Error(c, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
