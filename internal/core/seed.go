package core

import (
	"crypto/rand"
	"encoding/json"
	"math/big"

	_ "embed"

	"gitee.com/jieepre/keepblog/global"
	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/model/system"
	"gitee.com/jieepre/keepblog/pkg"
	"gitee.com/jieepre/keepblog/pkg/hash"
	"github.com/gookit/slog"
)

// 首次初始化的种子数据外置为文件，避免大段内容硬编码在 Go 源码里。

//go:embed seeddata/welcome_post.md
var welcomePostContent string

//go:embed seeddata/musics.json
var defaultMusicsJSON []byte

// defaultMusics 解析嵌入的默认音乐列表；数据错误属打包期问题，直接放弃写入。
func defaultMusics() []model.Music {
	var musics []model.Music
	if err := json.Unmarshal(defaultMusicsJSON, &musics); err != nil {
		return nil
	}
	return musics
}

// generatePassword 生成 n 位随机字母数字密码，用于首次启动的管理员账号。
// 密码只在初始化时打印一次，不落任何持久化明文。
func generatePassword(n int) string {
	const charset = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
	password := make([]byte, 0, n)
	for len(password) < n {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// crypto/rand 失败属于系统级故障，回退到固定长度不可取，直接终止初始化
			panic("生成随机密码失败: " + err.Error())
		}
		password = append(password, charset[idx.Int64()])
	}
	return string(password)
}

// initAdmin 首次启动时创建管理员账号
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

// initData 首次启动的基础数据：关于页、站点配置、默认分类与示例文章
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
