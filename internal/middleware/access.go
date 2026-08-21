package middleware

import (
	"strings"
	"time"

	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"gitee.com/jieepre/keepblog/pkg"
	"gitee.com/jieepre/keepblog/pkg/area"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
)

var ignoreURIS = []string{
	"/css",
	"/js",
	"/images",
	"/plugins",
	"/robots.txt",
}

func Statistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		for _, uri := range ignoreURIS {
			if strings.Contains(url, uri) {
				c.Next()
				return
			}

		}
		ip := c.ClientIP()
		ipL := pkg.Ip2long(ip)
		slog.Infof("ip:[%s],地区:[ %s ]", ip, area.Area(ipL))
		var accessLog system.AccessLog
		now := time.Now().Format(time.DateOnly)
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
