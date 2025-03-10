package middleware

import (
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model/system"
	"gitee.com/jieepre/go-site/pkg"
	"gitee.com/jieepre/go-site/pkg/area"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"strings"
	"time"
)

func Statistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		ignoreURIS := []string{
			"/css",
			"/js",
			"/images",
			"/plugins",
			"/robots.txt",
		}
		for _, uri := range ignoreURIS {

			if strings.Contains(url, uri) {
				return
			}

		}
		ip := c.ClientIP()
		ipL := pkg.Ip2long(ip)
		slog.Infof("ip:[%s],地区:[ %s ]", ip, area.Area(ipL))
		var accessLog system.AccessLog
		now := time.Now().Format("2006-01-02")
		global.GORM.
			Where("ip = ? and url = ? and strftime('%Y-%m-%d', create_at, 'unixepoch') = ?", pkg.Ip2long(ip), url, now).
			First(&accessLog)

		//记录pv请求次数
		if accessLog.Id != nil {
			global.GORM.Model(system.AccessLog{}).
				Where("id = ?", *accessLog.Id).
				Update("pv", *accessLog.PV+1)
		}
		c.Next()
		//记录access日志
		if accessLog.Id == nil {
			defaultView := 1
			log := system.AccessLog{
				Ip:      &ipL,
				URL:     url,
				Status:  c.Writer.Status(),
				PV:      &defaultView,
				UV:      &defaultView,
				UA:      c.Request.UserAgent(),
				Referer: c.Request.Referer(),
				Area:    area.Area(ipL),
			}
			global.GORM.Create(&log)
		}
	}
}
