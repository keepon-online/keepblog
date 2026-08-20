package core

import (
	"encoding/json"
	"fmt"
	"time"

	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/global"
	"gitee.com/jieepre/go-site/internal/logger"
	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/system"
	inpkg "gitee.com/jieepre/go-site/internal/pkg"
	postService "gitee.com/jieepre/go-site/internal/service/post"
	"gitee.com/jieepre/go-site/pkg"
	"gitee.com/jieepre/go-site/pkg/hash"
	"github.com/go-co-op/gocron"
	"github.com/gookit/slog"
)

func init() {
	global.GORM = InitDB()
}

func InitLog() {
	// 使用新的结构化日志系统
	logConfig := logger.LogConfig{
		Level:      "info",
		Format:     "text",
		Output:     "both", // 同时输出到控制台和文件
		Path:       "./logs",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	if err := logger.InitLogger(logConfig); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		// 退回到默认日志配置
		slog.Configure(func(logger *slog.SugaredLogger) {
			f := logger.Formatter.(*slog.TextFormatter)
			f.EnableColor = true
		})
	}
}

func InitResource() {
	initTable()
	if !hasData() {
		initAdmin()
		initData()
	}
	// 默认音乐独立判断：仅当 music 表为空时写入，不依赖 user 表是否已有数据，
	// 这样已初始化的旧库升级后也能补上默认音乐。
	initDefaultMusic()
	// 初始化数据库索引
	InitIndexes()
}

func initAdmin() {
	// 初始密码随机生成，仅在首次初始化时打印一次，登录后应立即修改
	password := generatePassword(16)
	hashed, err := pkg.HashPassword(password)
	if err != nil {
		slog.Errorf("初始化管理员账号失败（密码加密失败）:%s", err.Error())
		return
	}
	user := model.User{
		UserId:      1,
		Username:    "admin",
		Password:    hashed,
		NickName:    "管理员",
		Email:       "admin@example.com",
		Phonenumber: "",
		Sex:         1,
		Avatar:      "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Royalty%20Free%20Music%20Pack%20Cover.png",
	}
	if err := global.GORM.Save(&user).Error; err != nil {
		slog.Errorf("初始化管理员账号失败:%s", err.Error())
		return
	}
	slog.Infof("初始化管理员账号 admin，初始密码：%s（仅显示此次，请立即登录后台修改）", password)
}

func hasData() bool {
	user := model.User{}
	global.GORM.Model(user).First(&user)
	return user.UserId != 0
}

func initTable() {

	err := global.GORM.AutoMigrate(
		&model.Post{},
		&model.Category{},
		&model.FriendLink{},
		&model.Tag{},
		&model.PostTag{},
		&model.About{},
		&model.User{},
		&model.Music{},
		&system.AccessLog{},
		&system.WebSite{},
		&system.LoginLog{},
		&system.Notice{},
	)
	if err != nil {
		slog.Errorf("初始化Table失败:%s", err.Error())
	}
	slog.Info("初始化Table")
}

func initData() {
	slog.Info("初始化Data")
	tx := global.GORM.Begin()
	if err := tx.Save(&model.About{Id: 1, Note: "# 关于博客"}).Error; err != nil {
		tx.Rollback()
		slog.Errorf("初始化Data失败:%s", err.Error())
		return
	}
	if err := tx.Save(&system.WebSite{
		Id:          1,
		Icp:         "XICP备12939483号",
		Notice:      "<p>欢迎访问我的博客</p>",
		Title:       "XXX网站",
		Description: "XXX网站描述",
		URL:         "https://www.xxx.com",
		Keywords:    "这个网站的一些关键词",
		Copyright:   "Copyright © 2022-2023 XXX. All Rights Reserved.",
		BaiduStat:   "xxssssd",
		BaiduSite:   "xxxxddd",
		Github:      "https://github.com/xxx",
		Gitee:       "https://gitee.com/xxx",
	}).Error; err != nil {
		tx.Rollback()
		slog.Errorf("初始化Data失败:%s", err.Error())
		return
	}
	if err := tx.Save(&model.Category{
		CategoryId:   1,
		CategoryName: "测试",
		Note:         "测试分类",
		State:        1,
	}).Error; err != nil {
		tx.Rollback()
		slog.Errorf("初始化Data失败:%s", err.Error())
		return
	}
	hid, _ := hash.New().HashidsEncode([]int{1})

	if err := tx.Save(&model.Post{
		PostId:          1,
		Status:          1,
		CategoryId:      1,
		PostSlug:        hid,
		Author:          "佚名",
		Title:           "MySQL 统计实战：从基础到性能优化",
		CoverImage:      "https://ts1.tc.mm.bing.net/th/id/R-C.ccd320596cb9b0499c2d9e89079c7990?rik=bo30tkANeNk4Aw&riu=http%3a%2f%2fwww.finebornchina.cn%2fuploads%2fallimg%2f140430%2f1-140430150445413.jpg&ehk=Hjpp13uPkWtPTUVLZH%2f7V3MKAnXYJJNjmjRq1TE136k%3d&risl=&pid=ImgRaw&r=0",
		PostContent:     welcomePostContent,
		PostContentHtml: "#这是",
		Summary:         "MySQL 统计实战：从基础到性能优化",
		Type:            1,
		IsPublished:     1,
	}).Error; err != nil {
		tx.Rollback()
		slog.Errorf("初始化Data失败:%s", err.Error())
		return
	}

	if err := tx.Save(&model.Tag{
		TagId:   1,
		TagName: "mysql",
	}).Error; err != nil {
		tx.Rollback()
		slog.Errorf("初始化Data失败:%s", err.Error())
		return
	}
	if err := tx.Create(&model.PostTag{
		TagId:  1,
		PostId: 1,
	}).Error; err != nil {
		tx.Rollback()
		slog.Errorf("初始化Data失败:%s", err.Error())
		return
	}

	tx.Commit()
}

// initDefaultMusic 在 music 表为空时写入默认轻音乐。
// CC0 协议（免版权、免署名），托管于 jsDelivr CDN 直链，支持 CORS 与 Range。
// 曲目清单维护在 seeddata/musics.json。管理员可在后台随意增删改。
func initDefaultMusic() {
	var count int64
	global.GORM.Table(model.TMusicTable).Count(&count)
	if count > 0 {
		return
	}

	defaultMusics := defaultMusics()
	if len(defaultMusics) == 0 {
		slog.Error("默认音乐种子数据解析失败，跳过写入")
		return
	}
	for i := range defaultMusics {
		if err := global.GORM.Table(model.TMusicTable).Save(&defaultMusics[i]).Error; err != nil {
			slog.Errorf("初始化默认音乐失败:%s", err.Error())
			return
		}
	}
	slog.Info("初始化默认轻音乐完成")
}

func Timer() {
	slog.Info("定时任务启动")
	s := gocron.NewScheduler(time.Local)
	_, _ = s.Every(1).Day().At("21:00").Do(task)
	_, _ = s.Every(1).Day().At("23:30").Do(updateCoverTask)
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
