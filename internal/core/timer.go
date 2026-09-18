package core

import (
	"encoding/json"
	"fmt"
	"time"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/model"
	inpkg "gitee.com/jieepre/keepblog/internal/pkg"
	postService "gitee.com/jieepre/keepblog/internal/service/post"
	"gitee.com/jieepre/keepblog/pkg/area"
	"github.com/go-co-op/gocron"
	"github.com/gookit/slog"
)

// Timer 启动定时任务（百度推送 21:00、封面更新 23:30），
// 返回调度器供应用优雅退出时停止。
func Timer() *gocron.Scheduler {
	slog.Info("定时任务启动")
	s := gocron.NewScheduler(time.Local)
	_, _ = s.Every(1).Day().At("21:00").Do(task)
	_, _ = s.Every(1).Day().At("23:30").Do(updateCoverTask)
	_, _ = s.Every(1).Day().At("03:00").Do(updateIPDBTask)
	s.StartAsync()
	return s
}

func task() {
	baidu := config.Get().Baidu
	if baidu.Push {
		slog.Infof("定时开始运行")
		var content []model.Post
		global.GORM.Table(model.TPostsTable).Where("is_deleted=0 and is_published=1").Find(&content)
		urls := make([]string, len(content))
		for _, post := range content {
			urls = append(urls, fmt.Sprintf("%s/post/%s", baidu.Url, post.PostSlug))
		}
		api := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", baidu.Url, baidu.Token)
		site := inpkg.PushSite(urls, api)
		marshal, _ := json.Marshal(site)

		slog.Infof("定时任务百度收录推送记录:\n%s", string(marshal))
	}

}

func updateIPDBTask() {
	if err := area.UpdateIPDB(); err != nil {
		slog.Warnf("定时更新 IP 数据库失败: %v", err)
		return
	}
	slog.Infof("定时更新 IP 数据库成功")
}
func updateCoverTask() {
	slog.Infof("定时开始运行")
	var content []model.Post
	global.GORM.Table(model.TPostsTable).Where("is_deleted=0 and is_published=1").Find(&content)
	for _, post := range content {
		_ = postService.NewPostService(global.GORM).UpdatePostCoverImag(int(post.PostId))
	}
	slog.Infof("定时任务更新cover结束")
}
