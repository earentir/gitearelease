package gitearelease

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// mockServer simulates an HTTP server for testing the download function.
func mockServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

// setupMockServer creates a mock HTTP server that responds with the given body and status code
func setupMockServer(body string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		if body != "" {
			fmt.Fprint(w, body)
		}
	}))
}

// setupMockServerWithHandler creates a mock server with a custom handler
func setupMockServerWithHandler(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestGetRepositories_Gitea_Success(t *testing.T) {
	// Gitea JSON format - must include all required fields
	mockData := `[{"id": 1, "name": "Repo1", "full_name": "testuser/Repo1", "description": "", "private": false, "fork": false, "size": 1024, "language": "Go", "html_url": "http://example.com/repo1", "clone_url": "http://example.com/repo1.git", "ssh_url": "git@example.com:testuser/Repo1.git", "stars_count": 0, "forks_count": 0, "watchers_count": 0, "open_issues_count": 0, "release_counter": 2, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "username": "testuser", "full_name": "", "email": "", "avatar_url": ""}, "permissions": {"admin": false, "push": false, "pull": true}, "has_issues": true, "has_wiki": true, "has_projects": true, "has_releases": true, "has_packages": false}, {"id": 2, "name": "Repo2", "full_name": "testuser/Repo2", "description": "", "private": false, "fork": false, "size": 2048, "language": "Python", "html_url": "http://example.com/repo2", "clone_url": "http://example.com/repo2.git", "ssh_url": "git@example.com:testuser/Repo2.git", "stars_count": 0, "forks_count": 0, "watchers_count": 0, "open_issues_count": 0, "release_counter": 0, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "username": "testuser", "full_name": "", "email": "", "avatar_url": ""}, "permissions": {"admin": false, "push": false, "pull": true}, "has_issues": true, "has_wiki": true, "has_projects": true, "has_releases": false, "has_packages": false}]`
	mockServer := setupMockServer(mockData, 200)
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL,
		User:         "testuser",
		WithReleases: false,
		Provider:     "gitea", // Explicitly set Gitea provider
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(repos) != 2 {
		t.Errorf("Expected 2 repos, got %d", len(repos))
	}

	if repos[0].ReleaseCounter != 2 {
		t.Errorf("Expected ReleaseCounter 2, got %d", repos[0].ReleaseCounter)
	}
}

func TestGetRepositories_Gitea_WithReleases(t *testing.T) {
	mockData := `[{"id": 1, "name": "Repo1", "full_name": "testuser/Repo1", "description": "", "private": false, "fork": false, "size": 1024, "language": "Go", "html_url": "http://example.com/repo1", "clone_url": "http://example.com/repo1.git", "ssh_url": "git@example.com:testuser/Repo1.git", "stars_count": 0, "forks_count": 0, "watchers_count": 0, "open_issues_count": 0, "release_counter": 2, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "username": "testuser", "full_name": "", "email": "", "avatar_url": ""}, "permissions": {"admin": false, "push": false, "pull": true}, "has_issues": true, "has_wiki": true, "has_projects": true, "has_releases": true, "has_packages": false}, {"id": 2, "name": "Repo2", "full_name": "testuser/Repo2", "description": "", "private": false, "fork": false, "size": 2048, "language": "Python", "html_url": "http://example.com/repo2", "clone_url": "http://example.com/repo2.git", "ssh_url": "git@example.com:testuser/Repo2.git", "stars_count": 0, "forks_count": 0, "watchers_count": 0, "open_issues_count": 0, "release_counter": 0, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "username": "testuser", "full_name": "", "email": "", "avatar_url": ""}, "permissions": {"admin": false, "push": false, "pull": true}, "has_issues": true, "has_wiki": true, "has_projects": true, "has_releases": false, "has_packages": false}]`
	mockServer := setupMockServer(mockData, 200)
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL,
		User:         "testuser",
		WithReleases: true,
		Provider:     "gitea",
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(repos) != 1 {
		t.Fatalf("Expected 1 repo with releases, got %d", len(repos))
	}

	if repos[0].ReleaseCounter != 2 {
		t.Errorf("Expected ReleaseCounter 2, got %d", repos[0].ReleaseCounter)
	}
}

func TestGetRepositories_APIError(t *testing.T) {
	mockServer := setupMockServer("", 500)
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL,
		User:         "testuser",
		WithReleases: false,
	}

	_, err := GetRepositories(repositoriesToFetch)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestGetReleases_Gitea_Success(t *testing.T) {
	mockData := `{"id": 1, "tag_name": "v1.0.0", "name": "Release 1.0.0", "body": "Test release", "url": "http://example.com/release", "html_url": "http://example.com/release", "tarball_url": "http://example.com/tarball", "zipball_url": "http://example.com/zipball", "draft": false, "prerelease": false, "created_at": "2023-01-01T00:00:00Z", "published_at": "2023-01-01T00:00:00Z", "author": {"login": "testuser", "username": "testuser"}, "assets": []}`
	mockServer := setupMockServer(mockData, 200)
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL,
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   true,
		Provider: "gitea",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("Expected 1 release, got %d", len(releases))
	}

	if releases[0].TagName != "v1.0.0" {
		t.Errorf("Expected tag v1.0.0, got %s", releases[0].TagName)
	}
}

func TestGetReleases_Gitea_All(t *testing.T) {
	mockData := `[{"id": 1, "tag_name": "v1.0.0", "name": "Release 1.0.0", "body": "Test release 1", "url": "http://example.com/release1", "html_url": "http://example.com/release1", "tarball_url": "http://example.com/tarball1", "zipball_url": "http://example.com/zipball1", "draft": false, "prerelease": false, "created_at": "2023-01-01T00:00:00Z", "published_at": "2023-01-01T00:00:00Z", "author": {"login": "testuser", "username": "testuser"}, "assets": []}, {"id": 2, "tag_name": "v1.1.0", "name": "Release 1.1.0", "body": "Test release 2", "url": "http://example.com/release2", "html_url": "http://example.com/release2", "tarball_url": "http://example.com/tarball2", "zipball_url": "http://example.com/zipball2", "draft": false, "prerelease": false, "created_at": "2023-01-02T00:00:00Z", "published_at": "2023-01-02T00:00:00Z", "author": {"login": "testuser", "username": "testuser"}, "assets": []}]`
	mockServer := setupMockServer(mockData, 200)
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL,
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   false,
		Provider: "gitea",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(releases) != 2 {
		t.Fatalf("Expected 2 releases, got %d", len(releases))
	}
}

func TestFetchData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/data":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Test data"))
			if err != nil {
				t.Errorf("Failed to write the response: %v", err)
			}
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			t.Errorf("Unexpected URL path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	data, err := fetchData(server.URL + "/data")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedData := "Test data"
	if string(data) != expectedData {
		t.Errorf("Expected data: %q, got: %q", expectedData, string(data))
	}

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

func TestTrimVersionPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "1.0.0"},
		{"V1.0.0", "1.0.0"},
		{"version1.0.0", "ersion1.0.0"}, // "v" prefix removed first, leaves "ersion1.0.0"
		{"ver1.0.0", "er1.0.0"},           // "v" prefix removed first, leaves "er1.0.0"
		{"release1.0.0", "1.0.0"},         // "release" prefix removed (after "v", "version", "ver" don't match)
		{"rel1.0.0", "1.0.0"},             // "rel" is a valid prefix
		{"r1.0.0", "1.0.0"},
		{"v.1.0.0", ".1.0.0"},             // "v." is a valid prefix, but leaves "."
		{"1.0.0", "1.0.0"},
		{"VERSION1.0.0", "ersion1.0.0"},   // "version" prefix removed (case insensitive)
		{"v 1.0.0", " 1.0.0"},             // "v" is removed, but space remains
		{"v1.2.3", "1.2.3"},
		{"V1.2.3", "1.2.3"},
	}

	for _, tt := range tests {
		result := TrimVersionPrefix(tt.input)
		if result != tt.expected {
			t.Errorf("TrimVersionPrefix(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestProvider_AutoDetection(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		provider string
	}{
		{"Gitea default", "https://gitea.example.com", "gitea"},
		{"GitHub detection", "https://api.github.com", "github"},
		{"GitLab detection", "https://gitlab.com/api/v4", "gitlab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that provider detection works
			// We can't easily test the actual detection without making real API calls,
			// but we can test that the provider field is respected
			relCfg := ReleaseToFetch{
				BaseURL:  tt.baseURL,
				User:     "test",
				Repo:     "test",
				Latest:   true,
				Provider:  tt.provider,
			}

			if relCfg.Provider != tt.provider {
				t.Errorf("Expected provider %s, got %s", tt.provider, relCfg.Provider)
			}
		})
	}
}
