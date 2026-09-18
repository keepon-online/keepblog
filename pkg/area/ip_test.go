package area

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
