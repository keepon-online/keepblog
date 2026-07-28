package pkg

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gookit/slog"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func CoverImage() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(2))
	if n.Int64() == 1 {
		//美图
		return "https://cdn.seovx.com/?mom=302"
	} else {
		//古风
		return "https://cdn.seovx.com/ha/?mom=302"
	}
}

func GetPixabayImage() string {
	pixabay := Pixabay()
	if pixabay != nil {
		if pixabay.Total > 0 {
			//随机获取一个图片
			i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(pixabay.Hits))))
			return pixabay.Hits[i.Int64()].LargeImageURL
		}
	}
	return ""
}

// Pixabay 获取图片
func Pixabay() *PixabayResp {
	//设置请求参数
	pixabayRequest := PixabayRequest{
		Category:      "nature",
		Colors:        "",
		EditorsChoice: true,
		ImageType:     "photo",
		Key:           "36601714-193669790ec12a6f183afc0f1",
		Lang:          "zh",
		MinHeight:     "800",
		MinWidth:      "1280",
		Order:         "popular",
		Orientation:   "horizontal",
		Page:          "1",
		PerPage:       "200",
		Pretty:        true,
		Q:             "",
		Safesearch:    true,
	}

	v := url.Values{}
	val := reflect.ValueOf(pixabayRequest)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("json")
		value := fmt.Sprintf("%v", val.Field(i))
		if value != "" {
			v.Add(strings.ToLower(tag), value)
		}
	}

	query := v.Encode()
	//创建请求
	req, _ := http.NewRequest("GET", "https://pixabay.com/api?"+query, nil)
	//设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0")
	//向服务器发送请求
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Errorf("请求失败：%s", err.Error())
		return nil
	}
	defer func() { _ = resp.Body.Close() }()
	//响应内容
	statusCode := resp.StatusCode
	if statusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Errorf("请求失败：%s", err.Error())
			return nil
		}
		var pixabay PixabayResp
		err = json.Unmarshal(body, &pixabay)
		if err != nil {
			slog.Errorf("请求失败：%s", err.Error())
			return nil
		}
		return &pixabay
	}
	return nil
}

type PixabayRequest struct {
	Category      string `json:"category"`
	Colors        string `json:"colors,omitempty"`
	EditorsChoice bool   `json:"editors_choice"`
	ImageType     string `json:"image_type"`
	Key           string `json:"key"`
	Lang          string `json:"lang"`
	MinHeight     string `json:"min_height"`
	MinWidth      string `json:"min_width"`
	Order         string `json:"order"`
	Orientation   string `json:"orientation"`
	Page          string `json:"page"`     // 页数
	PerPage       string `json:"per_page"` // 每页默认20
	Pretty        bool   `json:"pretty"`
	Q             string `json:"q"`
	Safesearch    bool   `json:"safesearch"`
}
type PixabayResp struct {
	Total     int    `json:"total"`
	TotalHits int    `json:"totalHits"`
	Hits      []Hits `json:"hits"`
}

type Hits struct {
	Id              int    `json:"id"`
	PageURL         string `json:"pageURL"`
	Type            string `json:"type"`
	Tags            string `json:"tags"`
	PreviewURL      string `json:"previewURL"`
	PreviewWidth    int    `json:"previewWidth"`
	PreviewHeight   int    `json:"previewHeight"`
	WebformatURL    string `json:"webformatURL"`
	WebformatWidth  int    `json:"webformatWidth"`
	WebformatHeight int    `json:"webformatHeight"`
	LargeImageURL   string `json:"largeImageURL"`
	ImageWidth      int    `json:"imageWidth"`
	ImageHeight     int    `json:"imageHeight"`
	ImageSize       int    `json:"imageSize"`
	Views           int    `json:"views"`
	Downloads       int    `json:"downloads"`
	Collections     int    `json:"collections"`
	Likes           int    `json:"likes"`
	Comments        int    `json:"comments"`
	UserId          int    `json:"user_id"`
	User            string `json:"user"`
	UserImageURL    string `json:"userImageURL"`
}
