package area

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

// 下载文件函数
func downloadFile(url, filePath string) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态: %s", resp.Status)
	}

	// 复制内容到文件
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

func init() {
	ipdbPath := filepath.Join("data", "ip2region.xdb")

	// 如果文件不存在则下载
	if _, err := os.Stat(ipdbPath); os.IsNotExist(err) {
		url := "https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"
		fmt.Printf("下载 IP 数据库: %s\n", url)

		if err := downloadFile(url, ipdbPath); err != nil {
			fmt.Printf("下载失败: %v\n", err)
			// 创建空文件防止程序崩溃
			_ = os.WriteFile(ipdbPath, []byte{}, 0644)
		} else {
			fmt.Println("IP 数据库下载成功")
		}
	}
}

func Area(intIP uint32) string {
	// 1、从 dbPath 加载整个 xdb 到内存
	cBuff, err := xdb.LoadContentFromFile("./data/ip2region.xdb")
	if err != nil {
		fmt.Printf("failed: %s\n", err.Error())
		return ""
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		fmt.Printf("failed to create searcher with content: %s\n", err)
		return ""
	}

	// 备注：并发使用，用整个 xdb 缓存创建的 searcher 对象可以安全用于并发。
	search, err := searcher.Search(intIP)

	if err != nil {
		fmt.Printf("failed to search with content:%s\n", err)
		return ""
	}
	return search
}

// GetCurrAbPath 获取当前文件绝对路径
func GetCurrAbPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
