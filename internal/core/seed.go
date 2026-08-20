package core

import (
	"crypto/rand"
	"encoding/json"
	"math/big"

	_ "embed"

	"gitee.com/jieepre/go-site/internal/model"
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
