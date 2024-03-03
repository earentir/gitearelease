package gitearelease

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// mockServer simulates an HTTP server for testing the download function.
func mockServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("test content"))
		if err != nil {
			panic(err)
		}
	})
	return httptest.NewServer(handler)
}

func TestDownloadBinary_Success(t *testing.T) {
	server := mockServer()
	defer server.Close()

	tempDir, err := ioutil.TempDir("", "download_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	filename := "testfile.bin"
	filePath, err := DownloadBinary(server.URL, tempDir, filename)
	if err != nil {
		t.Fatalf("DownloadBinary() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("DownloadBinary() file %s does not exist", filePath)
	}

	expectedFilePath := filepath.Join(tempDir, filename)
	if filePath != expectedFilePath {
		t.Errorf("DownloadBinary() filePath = %v, want %v", filePath, expectedFilePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %s", err)
	}

	if string(content) != "test content" {
		t.Errorf("DownloadBinary() file content = %v, want %v", string(content), "test content")
	}
}
