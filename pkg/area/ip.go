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

const ipdbPath = "data/ip2region.xdb"

var (
	searcherMu  sync.RWMutex
	searcher    *xdb.Searcher
	searcherErr error
	updateMu    sync.Mutex
)

func downloadFile(url, filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	out, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".*.tmp")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpPath := out.Name()
	defer os.Remove(tmpPath)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		_ = out.Close()
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_ = out.Close()
		return fmt.Errorf("服务器返回错误状态: %s", resp.Status)
	}

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		_ = out.Close()
		return fmt.Errorf("写入文件失败: %w", err)
	}
	if written < 1024 {
		_ = out.Close()
		return fmt.Errorf("下载文件过小，可能不完整: %d bytes", written)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}
	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("重命名文件失败: %w", err)
	}
	return nil
}

func downloadWithFallback(urls []string, filePath string) error {
	var lastErr error
	for i, url := range urls {
		slog.Infof("尝试下载 IP 数据库 (%d/%d): %s", i+1, len(urls), url)
		if err := downloadFile(url, filePath); err != nil {
			lastErr = err
			slog.Errorf("下载失败: %v", err)
			continue
		}
		slog.Infof("IP 数据库下载成功")
		return nil
	}
	return fmt.Errorf("所有下载地址都失败: %w", lastErr)
}

// EnsureIPDB 确保 IP 数据库文件存在，缺失时尝试下载。
func EnsureIPDB() {
	if _, err := os.Stat(ipdbPath); os.IsNotExist(err) {
		if err := downloadWithFallback(ipDbUrls, ipdbPath); err != nil {
			slog.Warnf("IP数据库下载失败，IP地理位置功能将不可用: %v", err)
		}
	}
}

func loadSearcher(filePath string) (*xdb.Searcher, error) {
	content, err := xdb.LoadContentFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("加载IP数据库失败: %w", err)
	}
	result, err := xdb.NewWithBuffer(content)
	if err != nil {
		return nil, fmt.Errorf("创建IP数据库查询器失败: %w", err)
	}
	return result, nil
}

func initSearcher() (*xdb.Searcher, error) {
	searcherMu.Lock()
	defer searcherMu.Unlock()
	if searcher != nil || searcherErr != nil {
		return searcher, searcherErr
	}

	searcher, searcherErr = loadSearcher(ipdbPath)
	return searcher, searcherErr
}

// UpdateIPDB 下载并热切换 IP 数据库。失败时保留当前文件和查询器。
func UpdateIPDB() error {
	updateMu.Lock()
	defer updateMu.Unlock()

	if err := os.MkdirAll(filepath.Dir(ipdbPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(ipdbPath), filepath.Base(ipdbPath)+".*.tmp")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(tmpPath)

	var lastErr error
	for _, url := range ipDbUrls {
		if err := downloadFile(url, tmpPath); err != nil {
			lastErr = err
			continue
		}
		lastErr = nil
		break
	}
	if lastErr != nil {
		return fmt.Errorf("下载 IP 数据库失败: %w", lastErr)
	}

	loaded, err := loadSearcher(tmpPath)
	if err != nil {
		return err
	}
	if err := os.Rename(tmpPath, ipdbPath); err != nil {
		loaded.Close()
		return fmt.Errorf("替换 IP 数据库失败: %w", err)
	}

	searcherMu.Lock()
	old := searcher
	searcher = loaded
	searcherErr = nil
	searcherMu.Unlock()
	if old != nil {
		old.Close()
	}
	return nil
}

func Area(intIP uint32) string {
	if _, err := initSearcher(); err != nil {
		return ""
	}
	searcherMu.RLock()
	defer searcherMu.RUnlock()
	if searcher == nil {
		return ""
	}
	result, err := searcher.Search(intIP)
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
