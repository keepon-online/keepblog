package area

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gitee.com/jieepre/keepblog/config"
	"github.com/gookit/slog"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// 默认更新源：本仓库 GitHub Release 的滚动发布 ipdb-latest。
// CI（.github/workflows/ipdb.yml）每周从官方源数据构建、bench 校验后重建该
// Release，运行时更新只认此源（SHA256SUMS 校验），不跟踪上游 master。
// 该地址同时是版本锁定锚点：回滚时把 config 的 ipdb.updateUrl 指向某个
// ipdb-v* 归档 Release 的资产直链即可。
const (
	defaultReleaseXdbURL = "https://github.com/keepon-online/keepblog/releases/download/ipdb-latest/ip2region_v4.xdb"
	releaseSumName       = "SHA256SUMS"
)

// officialIPDBUrls 官方直链，仅用于冷启动兜底（本地无任何 xdb 文件且自管
// Release 不可达时），版本未锁定；日常定时更新绝不使用。
var officialIPDBUrls = []string{
	"https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v4.xdb",
	"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb",
}

// ipdbPath 数据库落盘位置；var 以便测试重定向到临时目录。
var ipdbPath = "data/ip2region.xdb"

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
	defer func() { _ = os.Remove(tmpPath) }()

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		_ = out.Close()
		return fmt.Errorf("下载失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
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

// releaseSource 返回 xdb 与 SHA256SUMS 的下载 URL。
// ipdb.updateUrl 配置为完整 xdb URL（国内镜像或归档回滚直链），校验文件
// 从同级目录推导。注意不能用 path.Dir：它会把 URL 的 "//" 折叠成 "/"。
func releaseSource() (xdbURL, sumURL string) {
	xdbURL = defaultReleaseXdbURL
	if cfg := config.Get().Ipdb; cfg != nil && cfg.UpdateUrl != "" {
		xdbURL = strings.TrimRight(cfg.UpdateUrl, "/")
	}
	return xdbURL, xdbURL[:strings.LastIndex(xdbURL, "/")+1] + releaseSumName
}

// fetchChecksum 下载并解析 SHA256SUMS（格式：<hex64>  <文件名>）。
func fetchChecksum(sumURL string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(sumURL)
	if err != nil {
		return "", fmt.Errorf("下载 SHA256SUMS 失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载 SHA256SUMS 失败: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", fmt.Errorf("读取 SHA256SUMS 失败: %w", err)
	}
	fields := strings.Fields(string(body))
	if len(fields) == 0 || len(fields[0]) != 64 {
		return "", fmt.Errorf("SHA256SUMS 格式无效")
	}
	return strings.ToLower(fields[0]), nil
}

// fileSHA256 计算文件哈希
func fileSHA256(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// EnsureIPDB 确保 IP 数据库文件存在（冷启动）。优先自管 Release（带哈希
// 校验），失败时回退官方直链并明确告警——这是官方源唯一的使用场景。
func EnsureIPDB() {
	if _, err := os.Stat(ipdbPath); !os.IsNotExist(err) {
		return
	}
	xdbURL, sumURL := releaseSource()
	ensureIPDBFrom(xdbURL, sumURL)
}

func ensureIPDBFrom(xdbURL, sumURL string) {
	if err := updateFrom(xdbURL, sumURL); err != nil {
		slog.Warnf("自管 Release 下载失败: %v", err)
	} else {
		return
	}
	slog.Warn("冷启动回退官方源（版本未锁定，仅本次）")
	if err := downloadWithFallback(officialIPDBUrls, ipdbPath); err != nil {
		slog.Warnf("IP数据库下载失败，IP地理位置功能将不可用: %v", err)
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

// UpdateIPDB 从自管 Release 下载并热切换 IP 数据库（含 SHA-256 校验）。
// 失败时保留当前文件和查询器，绝不回退官方源。
func UpdateIPDB() error {
	xdbURL, sumURL := releaseSource()
	return updateFrom(xdbURL, sumURL)
}

// updateFrom 从指定源下载并热切换：取哈希 → 未变化则跳过 → 下载校验
// → 加载验证 → 原子替换 → 锁内切换查询器。任一步失败保留现状。
func updateFrom(xdbURL, sumURL string) error {
	updateMu.Lock()
	defer updateMu.Unlock()

	want, err := fetchChecksum(sumURL)
	if err != nil {
		return err
	}

	// 哈希与当前文件一致则无需更新，避免每日无效下载与替换
	if current, err := fileSHA256(ipdbPath); err == nil && current == want {
		slog.Infof("IP 数据库已是最新（sha256 %s…）", want[:12])
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(ipdbPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(ipdbPath), filepath.Base(ipdbPath)+".*.tmp")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	if err := downloadFile(xdbURL, tmpPath); err != nil {
		return fmt.Errorf("下载 IP 数据库失败: %w", err)
	}
	got, err := fileSHA256(tmpPath)
	if err != nil {
		return fmt.Errorf("计算下载数据哈希失败: %w", err)
	}
	if got != want {
		return fmt.Errorf("SHA-256 校验失败: 期望 %s，实际 %s，拒绝替换", want, got)
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
	slog.Infof("IP 数据库已更新（sha256 %s…）", want[:12])
	return nil
}

func Area(intIP uint32) string {
	if _, err := initSearcher(); err != nil {
		recordNoDB()
		return ""
	}
	searcherMu.RLock()
	if searcher == nil {
		searcherMu.RUnlock()
		recordNoDB()
		return ""
	}
	result, err := searcher.Search(intIP)
	searcherMu.RUnlock()

	// 阶段一基线统计：纯内存原子计数，不保存原始 IP
	if err != nil {
		recordQuery(intIP, "")
		return ""
	}
	recordQuery(intIP, result)
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
