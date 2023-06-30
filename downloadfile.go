package gitearelease

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadBinary downloads a binary from a URL and saves it to a file
func DownloadBinary(url, outputDir, filename string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download binary: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download binary: server returned status %d", resp.StatusCode)
	}

	outputFile, err := os.Create(filepath.Join(outputDir, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write output file: %s", err)
	}

	return filepath.Join(outputDir, filename), nil
}
