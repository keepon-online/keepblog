package area

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gookit/slog"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// IP数据库下载地址列表（按优先级排序）
var ipDbUrls = []string{
	"https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v4.xdb",
	"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb",
}

// 全局缓存
var (
	searcher     *xdb.Searcher
	searcherOnce sync.Once
	searcherErr  error
)

// 下载文件函数
func downloadFile(url, filePath string) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建临时文件
	tmpPath := filePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer func() { _ = out.Close() }()

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("下载失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("服务器返回错误状态: %s", resp.Status)
	}

	// 复制内容到文件
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("写入文件失败: %w", err)
	}

	// 验证文件大小
	if written < 1024 {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("下载文件过小，可能不完整: %d bytes", written)
	}

	// 关闭文件后重命名
	_ = out.Close()
	if err := os.Rename(tmpPath, filePath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	return nil
}

// downloadWithFallback 带备用地址的下载
func downloadWithFallback(urls []string, filePath string) error {
	var lastErr error
	for i, url := range urls {
		slog.Infof("尝试下载 IP 数据库 (%d/%d): %s\n", i+1, len(urls), url)
		if err := downloadFile(url, filePath); err != nil {
			lastErr = err
			slog.Errorf("下载失败: %v\n", err)
			continue
		}
		slog.Infof("IP 数据库下载成功")
		return nil
	}
	return fmt.Errorf("所有下载地址都失败: %w", lastErr)
}

// EnsureIPDB 确保 IP 数据库文件存在，缺失时尝试下载。
// 由应用启动时显式调用（原先在 init() 里执行，import 即可能触发网络下载）；
// 下载失败仅告警，Area 查询会降级返回空结果。
func EnsureIPDB() {
	ipdbPath := filepath.Join("data", "ip2region.xdb")

	// 如果文件不存在则下载
	if _, err := os.Stat(ipdbPath); os.IsNotExist(err) {
		if err := downloadWithFallback(ipDbUrls, ipdbPath); err != nil {
			slog.Warnf("IP数据库下载失败，IP地理位置功能将不可用: %v", err)
		}
	}
}

// initSearcher 初始化搜索器（懒加载）
func initSearcher() (*xdb.Searcher, error) {
	searcherOnce.Do(func() {
		ipdbPath := "./data/ip2region.xdb"

		// 检查文件是否存在
		if _, err := os.Stat(ipdbPath); os.IsNotExist(err) {
			searcherErr = fmt.Errorf("IP数据库文件不存在: %s", ipdbPath)
			return
		}

		// 加载到内存
		cBuff, err := xdb.LoadContentFromFile(ipdbPath)
		if err != nil {
			searcherErr = fmt.Errorf("加载IP数据库失败: %w", err)
			return
		}

		// 创建搜索器
		searcher, searcherErr = xdb.NewWithBuffer(cBuff)
	})

	return searcher, searcherErr
}

func Area(intIP uint32) string {
	s, err := initSearcher()
	if err != nil {
		// 静默失败，避免日志刷屏
		return ""
	}

	result, err := s.Search(intIP)
	if err != nil {
		return ""
	}
	return result
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
