package core

import (
	"encoding/json"
	"fmt"
	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
	inpkg "gitee.com/jieepre/go-site/internal/pkg"
	postService "gitee.com/jieepre/go-site/internal/service/post"
	"gitee.com/jieepre/go-site/pkg"
	"github.com/go-co-op/gocron"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
	"time"
)

func init() {
	global.GORM = InitDB()
}

func InitLog() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	fileHandler, err := handler.NewRotateFileHandler("./logs/site.log", rotatefile.EveryDay, handler.WithLogLevels(slog.AllLevels))
	if err != nil {
		return
	}
	slog.PushHandler(fileHandler)
}

func InitResource() {
	initTable()
	if !hasData() {
		initAdmin()
		initData()
	}
}

func initAdmin() {
	password, _ := pkg.HashPassword("Aa123456")
	user := model.User{
		UserId:      1,
		Username:    "admin",
		Password:    password,
		NickName:    "管理员",
		Email:       "jieepre@outlook.com",
		Phonenumber: "19928902099",
		Sex:         1,
		Avatar:      "https://img0.baidu.com/it/u=3138581320,1425802456&fm=253&fmt=auto&app=138&f=GIF?w=480&h=480",
	}
	global.GORM.Save(&user)
	slog.Info("初始化管理员账号")
}

func hasData() bool {
	user := model.User{}
	global.GORM.Model(user).First(&user)
	return user.UserId != 0
}

func initTable() {

	err := global.GORM.AutoMigrate(
		&model.Comment{},
		&model.CommentReply{},
		&model.Post{},
		&model.Category{},
		&model.FriendLink{},
		&model.Tag{},
		&model.PostTag{},
		&model.About{},
		&model.User{},
		&system.AccessLog{},
		&system.WebSite{},
		&system.LoginLog{},
	)
	if err != nil {
		slog.Errorf("初始化Table失败:%s", err.Error())
	}
	slog.Info("初始化Table")
}

func initData() {
	global.GORM.Table(model.TAboutTable).Save(&model.About{Id: 1, Note: "# 关于博客"})
	slog.Info("初始化Data")
}

func Timer() {
	slog.Info("定时任务启动")
	s := gocron.NewScheduler(time.Local)
	s.Every(1).Day().At("21:00").Do(task)
	s.Every(1).Day().At("23:30").Do(updateCoverTask)
	s.StartAsync()
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

func updateCoverTask() {
	slog.Infof("定时开始运行")
	var content []model.Post
	global.GORM.Table(model.TPostsTable).Where("is_deleted=0 and is_published=1").Find(&content)
	for _, post := range content {
		_ = postService.NewPostService().UpdatePostCoverImag(int(post.PostId))
	}
	slog.Infof("定时任务更新cover结束")
}
