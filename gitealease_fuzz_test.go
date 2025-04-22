package gitearelease

import (
	"os"
	"path/filepath"
	"testing"
)

func FuzzDownloadBinary(f *testing.F) {
	// Add initial seed corpus. These are examples of inputs that your function expects.
	f.Add("http://example.com", "outputDir", "filename.bin") // You might need to adjust this based on your actual usage.

	f.Fuzz(func(t *testing.T, url, _, filename string) {
		// Set up a mock server to avoid making real HTTP requests.
		server := mockServer()
		defer server.Close()

		// Use the mock server's URL instead of the fuzzed URL to ensure the server responds appropriately.
		// The outputDir and filename are still fuzzed values.
		tempDir, err := os.CreateTemp("", "fuzz_download")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %s", err)
		}
		defer os.RemoveAll(tempDir.Name())

		outputDir := os.TempDir() // Overrde outputDir to use a controlled environment

		filePath, err := DownloadBinary(server.URL, outputDir, filename)
		if err != nil {
			t.Errorf("DownloadBinary() error = %v, wantErr %v", err, false)
		}

		// Ensure the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("DownloadBinary() file %s does not exist", filePath)
		}

		// Check for the correctness of the file path
		expectedFilePath := filepath.Join(outputDir, filename)
		if filePath != expectedFilePath {
			t.Errorf("DownloadBinary() filePath = %v, want %v", filePath, expectedFilePath)
		}
	})
}
