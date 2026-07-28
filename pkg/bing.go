package pkg

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetBingImage() string {
	// 获取今天的日期
	today := time.Now().Format("2006-01-02")
	// 构建必应壁纸的 URL
	url := fmt.Sprintf("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN&date=%s", today)
	// 发送 HTTP 请求获取壁纸信息
	req, _ := http.NewRequest("GET", url, nil)
	//设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0")
	//向服务器发送请求
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching wallpaper info:", err)
		return "https://www.bing.com/th?id=OHR.OcalaNF_ZH-CN1112502059_1920x1080.jpg&amp;rf=LaDigue_1920x1080.jpg&amp;pid=hp"
	}
	defer func() { _ = resp.Body.Close() }()
	// 解析 JSON 响应获取壁纸的 URL
	var data struct {
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return "https://www.bing.com/th?id=OHR.OcalaNF_ZH-CN1112502059_1920x1080.jpg&amp;rf=LaDigue_1920x1080.jpg&amp;pid=hp"
	}
	if len(data.Images) == 0 {
		fmt.Println("No wallpaper found for today")
		return "https://www.bing.com/th?id=OHR.OcalaNF_ZH-CN1112502059_1920x1080.jpg&amp;rf=LaDigue_1920x1080.jpg&amp;pid=hp"
	}
	wallpaperURL := "https://www.bing.com" + data.Images[0].URL

	return wallpaperURL
}
