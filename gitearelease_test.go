package gitearelease

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
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

	tempDir, err := os.CreateTemp("", "download_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir.Name())

	filename := "testfile.bin"
	filePath, err := DownloadBinary(server.URL, os.TempDir(), filename)
	if err != nil {
		t.Fatalf("DownloadBinary() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("DownloadBinary() file %s does not exist", filePath)
	}

	expectedFilePath := filepath.Join(os.TempDir(), filename)
	if filePath != expectedFilePath {
		t.Errorf("DownloadBinary() filePath = %v, want %v", filePath, expectedFilePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %s", err)
	}

	if string(content) != "test content" {
		t.Errorf("DownloadBinary() file content = %v, want %v", string(content), "test content")
	}
}

func setupMockServer(body string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, body)
	}))
}

func TestGetRepositories_Success(t *testing.T) {
	mockData := `[{"id": 1, "name": "Repo1", "ReleaseCounter": 2},{"id": 2, "name": "Repo2", "ReleaseCounter": 0}]`
	mockServer := setupMockServer(mockData, 200)
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL, // Use the URL of the mock server
		User:         "testuser",
		WithReleases: false,
	}

	expectedRepos := []Repository{
		{ID: 1, Name: "Repo1", ReleaseCounter: 2},
		{ID: 2, Name: "Repo2", ReleaseCounter: 0},
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(repos, expectedRepos) {
		t.Errorf("Expected %+v, got %+v", expectedRepos, repos)
	}
}

func TestGetRepositories_APIError(t *testing.T) {
	mockServer := setupMockServer("", 500) // Simulate server error
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL, // Use the URL of the mock server
		User:         "testuser",
		WithReleases: false,
	}

	_, err := GetRepositories(repositoriesToFetch)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestGetRepositories_WithReleases(t *testing.T) {
	mockData := `[{"id": 1, "name": "Repo1", "ReleaseCounter": 2},{"id": 2, "name": "Repo2", "ReleaseCounter": 0}]`
	mockServer := setupMockServer(mockData, 200)
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL, // Use the URL of the mock server
		User:         "testuser",
		WithReleases: true, // Only fetch repos with releases
	}

	expectedRepos := []Repository{
		{ID: 1, Name: "Repo1", ReleaseCounter: 2},
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(repos, expectedRepos) {
		t.Errorf("Expected %+v, got %+v", expectedRepos, repos)
	}
}

func TestFetchData(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the path part of the URL
		if r.URL.Path != "/data" {
			t.Errorf("Expected URL path: /data, got: %s", r.URL.Path)
		}

		switch r.URL.Path {
		case "/data":
			// Set the response status code
			w.WriteHeader(http.StatusOK)

			// Set the response body
			_, err := w.Write([]byte("Test data"))
			if err != nil {
				t.Errorf("Failed to write the response: %v", err)
			}
		case "/error":
			// Simulate an error response
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	data, err := fetchData(server.URL + "/data")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the returned data
	expectedData := "Test data"
	if string(data) != expectedData {
		t.Errorf("Expected data: %s, got: %s", expectedData, data)
	}

	// Test the error case by requesting the error path
	_, err = fetchData(server.URL + "/error")
	if err == nil {
		t.Errorf("Expected an error but got nil")
	}
}

func TestCompareVersionsHelper(t *testing.T) {
	versionstrings := VersionStrings{
		VersionStrings: versionstringstruct{
			Older: "There is a newer release available",
			Equal: "You are up to date",
			Newer: "You are on an unreleased version",
		},
		VersionOptions: versionoptionsstruct{
			DieIfOlder:           true,
			DieIfNewer:           true,
			ShowMessageOnCurrent: true,
		},
	}

	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.1"

	expected := "There is a newer release available"
	versionstrings.VersionStrings.Older = ""
	versionstrings.VersionStrings.UpgradeURL = ""
	versionstrings.VersionOptions.DieIfOlder = false
	result := CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.VersionOptions.DieIfOlder = false
	versionstrings.VersionStrings.UpgradeURL = "https://example.com/upgrade"
	expected = "There is a newer release available at https://example.com/upgrade"
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	expected = "There is a newer release available"
	versionstrings.VersionStrings.UpgradeURL = ""
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.0"
	versionstrings.VersionOptions.ShowMessageOnCurrent = false
	expected = ""
	versionstrings.VersionStrings.Equal = ""
	versionstrings.VersionOptions.DieIfOlder = false
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.0"
	versionstrings.VersionOptions.ShowMessageOnCurrent = true
	expected = "You are up to date"
	versionstrings.VersionStrings.Equal = ""
	versionstrings.VersionOptions.DieIfOlder = false
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.VersionOptions.ShowMessageOnCurrent = true
	versionstrings.VersionOptions.DieIfNewer = false
	expected = "You are on an unreleased version"
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.0"
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.VersionOptions.ShowMessageOnCurrent = true
	versionstrings.VersionOptions.DieIfNewer = false
	versionstrings.VersionStrings.Newer = ""
	expected = "You are on an unreleased version"
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.0"
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

}
func TestCompareVersions(t *testing.T) {
	versionstrings := VersionStrings{
		VersionStrings: versionstringstruct{
			Older:      "Older version",
			Equal:      "Equal version",
			Newer:      "Newer version",
			UpgradeURL: "https://example.com/upgrade",
		},
		VersionOptions: versionoptionsstruct{
			DieIfOlder:           true,
			DieIfNewer:           true,
			ShowMessageOnCurrent: true,
		},
	}

	// Test case 1: Own version is older than the current version
	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.1"
	expected := -1
	result := CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 2: Own version is newer than the current version
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.0"
	expected = 1
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 3: Own version is equal to the current version
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.1"
	expected = 0
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 4: Own version has more version numbers than the current version
	versionstrings.Own = "1.0.1.1"
	versionstrings.Latest = "1.0.1"
	expected = 1
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 5: Own version has fewer version numbers than the current version
	versionstrings.Own = "1.0"
	versionstrings.Latest = "1.0.1"
	expected = -1
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}
}
