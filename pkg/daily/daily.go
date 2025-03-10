package daily

import (
	"encoding/json"
	"fmt"
	"github.com/6tail/lunar-go/calendar"
	"github.com/gookit/slog"
	"io"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Date      string   `json:"date"`
		News      []string `json:"news"`
		Weiyu     string   `json:"weiyu"`
		Image     string   `json:"image"`
		HeadImage string   `json:"head_image"`
	} `json:"data"`
	Time  int    `json:"time"`
	Usage int    `json:"usage"`
	LogId string `json:"log_id"`
}

func GetDailyReport() string {

	resp, err := http.Get("https://v2.alapi.cn/api/zaobao?token=OM8PzGLpWHOkNyV7&format=json")

	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Errorf("GetDailyReport response body:%s", err.Error())
		return ""
	}
	var result Result
	_ = json.Unmarshal(body, &result)

	var newsStr strings.Builder

	newsStr.WriteString(GetTitle(result.Data.Date) + "\n")
	newsStr.WriteString("每天60秒读懂世界\n")

	if result.Code == 200 {
		for _, news := range result.Data.News {
			newsStr.WriteString(news + "\n")
		}
		newsStr.WriteString(result.Data.Weiyu)
	}

	return newsStr.String()
}

func GetTitle(dateStr string) string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	date, _ := time.ParseInLocation("2006-01-02", dateStr, location)
	solar := calendar.NewSolarFromDate(date)
	year := solar.GetYear()
	month := solar.GetMonth()
	day := solar.GetDay()
	week := solar.GetWeekInChinese()
	lunar := solar.GetLunar()
	monthInChinese := lunar.GetMonthInChinese()
	dayInChinese := lunar.GetDayInChinese()
	title := fmt.Sprintf("今天是%d年%d月%d日 星期%s 农历:%s月%s", year, month, day, week, monthInChinese, dayInChinese)
	return title
}
