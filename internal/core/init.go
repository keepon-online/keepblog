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

	content := `
# MySQL 统计实战：从基础到性能优化

## 一、统计的基石：核心函数
MySQL 提供丰富的聚合函数来满足统计需求：

1. ​**​COUNT​**​：统计行数
    SELECT COUNT(*) FROM orders; -- 统计总订单量

2. ​**​SUM/AVG​**​：数值计算
    SELECT SUM(amount) AS total_sales, AVG(amount) AS avg_price 
    FROM orders WHERE create_date > '2024-01-01';

3. ​**​MAX/MIN​**​：极值查询
    SELECT MAX(temperature), MIN(humidity) FROM sensor_data;

## 二、进阶统计：分组与过滤
通过 GROUP BY 实现多维统计：

    -- 按日期统计销售额
    SELECT DATE(create_time) AS day, 
           SUM(amount) AS daily_sales,
           COUNT(DISTINCT user_id) AS active_users 
    FROM orders 
    GROUP BY day 
    HAVING daily_sales > 10000;

注意：HAVING 用于分组后过滤，WHERE 用于分组前过滤

## 三、性能优化三板斧
### 1. 索引策略
- 为 WHERE/GROUP BY/ORDER BY 涉及的列创建复合索引
- 优先选择区分度高的字段作为索引前导列

### 2. 统计信息管理
    ANALYZE TABLE orders; -- 手动更新统计信息
    SHOW TABLE STATUS LIKE 'orders'; -- 查看估算值

建议在低峰期执行统计信息更新

## 四、典型应用场景
1. ​**​用户行为分析​**​  
       -- 统计7日留存率
       SELECT reg_date,
              COUNT(DISTINCT user_id) AS reg_users,
              COUNT(DISTINCT CASE WHEN login_date = reg_date + INTERVAL 7 DAY THEN user_id END)/COUNT(DISTINCT user_id) AS retention_rate
       FROM user_events 
       GROUP BY reg_date;


**引用说明**
: 基础聚合函数与应用场景
: 分组统计与函数详解
: 统计信息管理与性能优化
: 典型业务场景示例
: 高性能统计实现方案
`
	if err := tx.Save(&model.Post{
		PostId:          1,
		Status:          1,
		CategoryId:      1,
		PostSlug:        hid,
		Author:          "佚名",
		Title:           "MySQL 统计实战：从基础到性能优化",
		CoverImage:      "https://ts1.tc.mm.bing.net/th/id/R-C.ccd320596cb9b0499c2d9e89079c7990?rik=bo30tkANeNk4Aw&riu=http%3a%2f%2fwww.finebornchina.cn%2fuploads%2fallimg%2f140430%2f1-140430150445413.jpg&ehk=Hjpp13uPkWtPTUVLZH%2f7V3MKAnXYJJNjmjRq1TE136k%3d&risl=&pid=ImgRaw&r=0",
		PostContent:     content,
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
// 仓库：effacestudios/Royalty-Free-Music-Pack。管理员可在后台随意增删改。
func initDefaultMusic() {
	var count int64
	global.GORM.Table(model.TMusicTable).Count(&count)
	if count > 0 {
		return
	}

	defaultMusics := []model.Music{
		{Name: "Bubbles", Artist: "Royalty Free", Url: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Bubbles.mp3", Cover: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Royalty%20Free%20Music%20Pack%20Cover.png", Sort: 1, State: 1},
		{Name: "Happy Life", Artist: "Royalty Free", Url: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Happy%20Life.mp3", Cover: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Royalty%20Free%20Music%20Pack%20Cover.png", Sort: 2, State: 1},
		{Name: "Newness", Artist: "Royalty Free", Url: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Newness.mp3", Cover: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Royalty%20Free%20Music%20Pack%20Cover.png", Sort: 3, State: 1},
		{Name: "Planning", Artist: "Royalty Free", Url: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Planning.mp3", Cover: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Royalty%20Free%20Music%20Pack%20Cover.png", Sort: 4, State: 1},
		{Name: "Mysterious", Artist: "Royalty Free", Url: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Mysterious.mp3", Cover: "https://cdn.jsdelivr.net/gh/effacestudios/Royalty-Free-Music-Pack@master/Royalty%20Free%20Music%20Pack%20Cover.png", Sort: 5, State: 1},
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
