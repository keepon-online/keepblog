package area

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadFileRejectsSmallResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not an xdb"))
	}))
	defer server.Close()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "ip2region.xdb")
	if err := downloadFile(server.URL, filePath); err == nil {
		t.Fatal("downloadFile accepted a response that is too small")
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("downloadFile left a destination file after failure: %v", err)
	}
}

func TestDownloadFileKeepsExistingFileOnFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "ip2region.xdb")
	original := []byte("existing database")
	if err := os.WriteFile(filePath, original, 0644); err != nil {
		t.Fatal(err)
	}
	if err := downloadFile(server.URL, filePath); err == nil {
		t.Fatal("downloadFile accepted an HTTP failure")
	}
	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("existing file changed after failed download: %q", got)
	}
}

// redirectIPDBPath 把包级 ipdbPath 重定向到临时目录，测试后恢复。
func redirectIPDBPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	old := ipdbPath
	ipdbPath = filepath.Join(dir, "ip2region.xdb")
	t.Cleanup(func() { ipdbPath = old })
	return ipdbPath
}

// resetSearcher 清理测试污染的全局查询器状态。
func resetSearcher(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		searcherMu.Lock()
		if searcher != nil {
			searcher.Close()
		}
		searcher, searcherErr = nil, nil
		searcherMu.Unlock()
	})
}

// fakeRelease 模拟自管 Release 端点：xdb 内容、其真实哈希或指定篡改哈希。
type fakeRelease struct {
	xdbBody    []byte
	sumSha     string // 为空时使用 xdbBody 的真实哈希
	sumFailure bool
	hits       int
}

func (f *fakeRelease) server(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ip2region_v4.xdb":
			f.hits++
			_, _ = w.Write(f.xdbBody)
		case "/SHA256SUMS":
			if f.sumFailure {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			sha := f.sumSha
			if sha == "" {
				sum := sha256.Sum256(f.xdbBody)
				sha = hex.EncodeToString(sum[:])
			}
			_, _ = w.Write([]byte(sha + "  ip2region_v4.xdb\n"))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestUpdateFromRejectsChecksumMismatch(t *testing.T) {
	dbPath := redirectIPDBPath(t)
	resetSearcher(t)
	old := []byte(strings.Repeat("current database ", 128))
	if err := os.WriteFile(dbPath, old, 0644); err != nil {
		t.Fatal(err)
	}

	rel := &fakeRelease{
		xdbBody: []byte(strings.Repeat("tampered database ", 128)),
		sumSha:  strings.Repeat("ab", 32), // 格式合法但与内容不符
	}
	server := rel.server(t)
	defer server.Close()

	err := updateFrom(server.URL+"/ip2region_v4.xdb", server.URL+"/SHA256SUMS")
	if err == nil || !strings.Contains(err.Error(), "SHA-256") {
		t.Fatalf("期望哈希校验失败，实际: %v", err)
	}
	got, _ := os.ReadFile(dbPath)
	if string(got) != string(old) {
		t.Fatal("校验失败后旧库被改动")
	}
}

func TestUpdateFromKeepsOldWhenReleaseUnavailable(t *testing.T) {
	dbPath := redirectIPDBPath(t)
	old := []byte(strings.Repeat("keep me ", 256))
	if err := os.WriteFile(dbPath, old, 0644); err != nil {
		t.Fatal(err)
	}

	// 官方源装上哨兵：任何请求都视为违规（UpdateIPDB 不得回退官方源）
	sentinel := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("UpdateIPDB 不应请求官方源")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer sentinel.Close()
	prevOfficial := officialIPDBUrls
	officialIPDBUrls = []string{sentinel.URL}
	t.Cleanup(func() { officialIPDBUrls = prevOfficial })

	rel := &fakeRelease{sumFailure: true}
	server := rel.server(t)
	defer server.Close()

	if err := updateFrom(server.URL+"/ip2region_v4.xdb", server.URL+"/SHA256SUMS"); err == nil {
		t.Fatal("Release 不可用时应当报错")
	}
	got, _ := os.ReadFile(dbPath)
	if string(got) != string(old) {
		t.Fatal("失败后旧库被改动")
	}
}

func TestEnsureIPDBFallsBackToOfficialOnColdStart(t *testing.T) {
	dbPath := redirectIPDBPath(t)

	rel := &fakeRelease{sumFailure: true} // 自管 Release 不可达
	server := rel.server(t)
	defer server.Close()

	officialBody := []byte(strings.Repeat("official fallback ", 128))
	official := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(officialBody)
	}))
	defer official.Close()
	prevOfficial := officialIPDBUrls
	officialIPDBUrls = []string{official.URL}
	t.Cleanup(func() { officialIPDBUrls = prevOfficial })

	ensureIPDBFrom(server.URL+"/ip2region_v4.xdb", server.URL+"/SHA256SUMS")

	got, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatal("冷启动回退后应存在数据库文件")
	}
	if !bytes.Equal(got, officialBody) {
		t.Fatal("回退源内容不符")
	}
}

func TestUpdateFromHotSwapsAndSkipsUnchanged(t *testing.T) {
	dbPath := redirectIPDBPath(t)
	resetSearcher(t)

	sample, err := os.ReadFile("../../data/ip2region.xdb")
	if err != nil {
		t.Skip("需要本地 data/ip2region.xdb 样本（CI 出库后跳过）")
	}
	rel := &fakeRelease{xdbBody: sample}
	server := rel.server(t)
	defer server.Close()

	if err := updateFrom(server.URL+"/ip2region_v4.xdb", server.URL+"/SHA256SUMS"); err != nil {
		t.Fatalf("首次更新失败: %v", err)
	}
	if rel.hits != 1 {
		t.Fatalf("xdb 端点请求数 = %d, want 1", rel.hits)
	}
	searcherMu.RLock()
	swapped := searcher
	searcherMu.RUnlock()
	if swapped == nil {
		t.Fatal("热切换后查询器为空")
	}
	got, _ := os.ReadFile(dbPath)
	if !bytes.Equal(got, sample) {
		t.Fatal("落盘内容与源不符")
	}

	// 第二次：哈希未变应跳过下载
	if err := updateFrom(server.URL+"/ip2region_v4.xdb", server.URL+"/SHA256SUMS"); err != nil {
		t.Fatalf("无变化更新报错: %v", err)
	}
	if rel.hits != 1 {
		t.Fatalf("无变化时不应重复下载，hits = %d", rel.hits)
	}
}

func TestReleaseSourceDerivesSumURL(t *testing.T) {
	xdbURL, sumURL := releaseSource()
	// 两个 URL 都必须是合法绝对地址（path.Dir 折叠 "//" 的回归防护）
	if !strings.HasPrefix(xdbURL, "https://") {
		t.Fatalf("xdbURL 非法: %s", xdbURL)
	}
	if !strings.HasPrefix(sumURL, "https://") {
		t.Fatalf("sumURL 非法(双斜杠被折叠的回归): %s", sumURL)
	}
	if want := xdbURL[:strings.LastIndex(xdbURL, "/")+1] + releaseSumName; sumURL != want {
		t.Fatalf("sumURL = %s, want %s", sumURL, want)
	}
}
