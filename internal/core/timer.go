package core

import (
	"encoding/json"
	"fmt"
	"strings"
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

// Timer 启动定时任务（搜索引擎推送 21:00、封面更新 23:30、IP 库更新 03:00、
// IP 查询基线统计快照 03:10），返回调度器供应用优雅退出时停止。
func Timer() *gocron.Scheduler {
	slog.Info("定时任务启动")
	s := gocron.NewScheduler(time.Local)
	_, _ = s.Every(1).Day().At("21:00").Do(pushSearchEngines)
	_, _ = s.Every(1).Day().At("23:30").Do(updateCoverTask)
	_, _ = s.Every(1).Day().At("03:00").Do(updateIPDBTask)
	_, _ = s.Every(1).Day().At("03:10").Do(area.LogStatsSnapshot)
	s.StartAsync()
	return s
}

// pushSearchEngines 每日向已配置的搜索引擎推送全部已发布文章 URL：
// 百度收录 API（普通收录）与 IndexNow（Bing/Yandex 等），二者共用一次文章查询。
func pushSearchEngines() {
	var content []model.Post
	global.GORM.Table(model.TPostsTable).Where("is_deleted=0 and is_published=1").Find(&content)

	if baidu := config.Get().Baidu; baidu != nil && baidu.Push {
		slog.Infof("定时开始运行")
		// 注意容量切片：make([]string, len(content)) 会先填满 len 个空串再 append，
		// 推送 body 前面会出现整排空行，白白消耗百度推送配额。
		urls := make([]string, 0, len(content))
		for _, post := range content {
			urls = append(urls, fmt.Sprintf("%s/post/%s", baidu.Url, post.PostSlug))
		}
		api := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", baidu.Url, baidu.Token)
		site := inpkg.PushSite(urls, api)
		marshal, _ := json.Marshal(site)

		slog.Infof("定时任务百度收录推送记录:\n%s", string(marshal))
	}

	if indexNow := config.Get().IndexNow; indexNow != nil && indexNow.Enable && indexNow.Key != "" {
		// 站点域名取 web_site.url（与 sitemap/robots 同源）。
		var baseURL string
		global.GORM.Raw("SELECT url FROM web_site LIMIT 1").Scan(&baseURL)
		baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
		if baseURL == "" {
			slog.Warnf("indexnow 推送跳过: web_site.url 为空")
			return
		}
		host := strings.TrimPrefix(strings.TrimPrefix(baseURL, "https://"), "http://")

		urls := make([]string, 0, len(content))
		for _, post := range content {
			urls = append(urls, fmt.Sprintf("%s/post/%s", baseURL, post.PostSlug))
		}
		if err := inpkg.PushIndexNow(urls, host, indexNow.Key, indexNow.Endpoint); err != nil {
			slog.Errorf("indexnow 推送失败: %v", err)
		}
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
