package gitearelease

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReleases_GitHub_Success(t *testing.T) {
	// GitHub JSON format
	mockData := `{"id": 1, "tag_name": "v1.0.0", "name": "Release 1.0.0", "body": "Test release", "url": "https://api.github.com/repos/testuser/testrepo/releases/1", "html_url": "https://github.com/testuser/testrepo/releases/tag/v1.0.0", "tarball_url": "https://api.github.com/repos/testuser/testrepo/tarball/v1.0.0", "zipball_url": "https://api.github.com/repos/testuser/testrepo/zipball/v1.0.0", "draft": false, "prerelease": false, "created_at": "2023-01-01T00:00:00Z", "published_at": "2023-01-01T00:00:00Z", "author": {"login": "testuser", "id": 1, "type": "User"}, "assets": [{"id": 1, "name": "binary.tar.gz", "size": 1024, "download_count": 42, "created_at": "2023-01-01T00:00:00Z", "browser_download_url": "https://github.com/testuser/testrepo/releases/download/v1.0.0/binary.tar.gz", "content_type": "application/gzip"}]}`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/testuser/testrepo/releases/latest" {
			t.Errorf("Expected path /repos/testuser/testrepo/releases/latest, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL,
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   true,
		Provider: "github",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("Expected 1 release, got %d", len(releases))
	}

	rel := releases[0]
	if rel.TagName != "v1.0.0" {
		t.Errorf("Expected tag v1.0.0, got %s", rel.TagName)
	}

	if rel.Name != "Release 1.0.0" {
		t.Errorf("Expected name Release 1.0.0, got %s", rel.Name)
	}

	if len(rel.Assets) != 1 {
		t.Fatalf("Expected 1 asset, got %d", len(rel.Assets))
	}

	asset := rel.Assets[0]
	if asset.Name != "binary.tar.gz" {
		t.Errorf("Expected asset name binary.tar.gz, got %s", asset.Name)
	}

	if asset.Size != 1024 {
		t.Errorf("Expected asset size 1024, got %d", asset.Size)
	}

	if asset.DownloadCount != 42 {
		t.Errorf("Expected download count 42, got %d", asset.DownloadCount)
	}

	// Note: GitHub doesn't have UUID field
	if asset.UUID != "" {
		t.Errorf("Expected empty UUID for GitHub, got %s", asset.UUID)
	}
}

func TestGetReleases_GitHub_All(t *testing.T) {
	mockData := `[{"id": 1, "tag_name": "v1.0.0", "name": "Release 1.0.0", "body": "Test release 1", "url": "https://api.github.com/repos/testuser/testrepo/releases/1", "html_url": "https://github.com/testuser/testrepo/releases/tag/v1.0.0", "tarball_url": "https://api.github.com/repos/testuser/testrepo/tarball/v1.0.0", "zipball_url": "https://api.github.com/repos/testuser/testrepo/zipball/v1.0.0", "draft": false, "prerelease": false, "created_at": "2023-01-01T00:00:00Z", "published_at": "2023-01-01T00:00:00Z", "author": {"login": "testuser", "id": 1, "type": "User"}, "assets": []}, {"id": 2, "tag_name": "v1.1.0", "name": "Release 1.1.0", "body": "Test release 2", "url": "https://api.github.com/repos/testuser/testrepo/releases/2", "html_url": "https://github.com/testuser/testrepo/releases/tag/v1.1.0", "tarball_url": "https://api.github.com/repos/testuser/testrepo/tarball/v1.1.0", "zipball_url": "https://api.github.com/repos/testuser/testrepo/zipball/v1.1.0", "draft": false, "prerelease": true, "created_at": "2023-01-02T00:00:00Z", "published_at": "2023-01-02T00:00:00Z", "author": {"login": "testuser", "id": 1, "type": "User"}, "assets": []}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/testuser/testrepo/releases" {
			t.Errorf("Expected path /repos/testuser/testrepo/releases, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL,
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   false,
		Provider: "github",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(releases) != 2 {
		t.Fatalf("Expected 2 releases, got %d", len(releases))
	}

	if releases[1].Prerelease != true {
		t.Errorf("Expected prerelease to be true for second release")
	}
}

func TestGetRepositories_GitHub_Success(t *testing.T) {
	// GitHub JSON format
	mockData := `[{"id": 1, "name": "repo1", "full_name": "testuser/repo1", "description": "Test repo 1", "private": false, "fork": false, "size": 1024, "language": "Go", "html_url": "https://github.com/testuser/repo1", "clone_url": "https://github.com/testuser/repo1.git", "ssh_url": "git@github.com:testuser/repo1.git", "stargazers_count": 10, "forks_count": 5, "watchers_count": 8, "open_issues_count": 2, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "avatar_url": "https://github.com/testuser.png", "type": "User"}, "permissions": {"admin": true, "push": true, "pull": true}, "has_issues": true, "has_wiki": true, "has_projects": true, "has_releases": true, "has_packages": false}, {"id": 2, "name": "repo2", "full_name": "testuser/repo2", "description": "Test repo 2", "private": false, "fork": false, "size": 2048, "language": "Python", "html_url": "https://github.com/testuser/repo2", "clone_url": "https://github.com/testuser/repo2.git", "ssh_url": "git@github.com:testuser/repo2.git", "stargazers_count": 20, "forks_count": 10, "watchers_count": 15, "open_issues_count": 0, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "avatar_url": "https://github.com/testuser.png", "type": "User"}, "permissions": {"admin": false, "push": false, "pull": true}, "has_issues": true, "has_wiki": false, "has_projects": false, "has_releases": false, "has_packages": false}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/testuser/repos" {
			t.Errorf("Expected path /users/testuser/repos, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL,
		User:         "testuser",
		WithReleases: false,
		Provider:     "github",
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(repos) != 2 {
		t.Fatalf("Expected 2 repos, got %d", len(repos))
	}

	// Check first repo
	if repos[0].Name != "repo1" {
		t.Errorf("Expected repo name repo1, got %s", repos[0].Name)
	}

	if repos[0].StarsCount != 10 {
		t.Errorf("Expected stars count 10, got %d", repos[0].StarsCount)
	}

	// GitHub ReleaseCounter is a placeholder (1 if HasReleases is true)
	if repos[0].ReleaseCounter != 1 {
		t.Errorf("Expected ReleaseCounter 1 (placeholder), got %d", repos[0].ReleaseCounter)
	}

	if repos[1].ReleaseCounter != 0 {
		t.Errorf("Expected ReleaseCounter 0 (no releases), got %d", repos[1].ReleaseCounter)
	}
}

func TestGetRepositories_GitHub_WithReleases(t *testing.T) {
	mockData := `[{"id": 1, "name": "repo1", "full_name": "testuser/repo1", "description": "Test repo 1", "private": false, "fork": false, "size": 1024, "language": "Go", "html_url": "https://github.com/testuser/repo1", "clone_url": "https://github.com/testuser/repo1.git", "ssh_url": "git@github.com:testuser/repo1.git", "stargazers_count": 10, "forks_count": 5, "watchers_count": 8, "open_issues_count": 2, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "avatar_url": "https://github.com/testuser.png", "type": "User"}, "permissions": {"admin": true, "push": true, "pull": true}, "has_issues": true, "has_wiki": true, "has_projects": true, "has_releases": true, "has_packages": false}, {"id": 2, "name": "repo2", "full_name": "testuser/repo2", "description": "Test repo 2", "private": false, "fork": false, "size": 2048, "language": "Python", "html_url": "https://github.com/testuser/repo2", "clone_url": "https://github.com/testuser/repo2.git", "ssh_url": "git@github.com:testuser/repo2.git", "stargazers_count": 20, "forks_count": 10, "watchers_count": 15, "open_issues_count": 0, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "login": "testuser", "avatar_url": "https://github.com/testuser.png", "type": "User"}, "permissions": {"admin": false, "push": false, "pull": true}, "has_issues": true, "has_wiki": false, "has_projects": false, "has_releases": false, "has_packages": false}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL,
		User:         "testuser",
		WithReleases: true,
		Provider:     "github",
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should only return repo1 which has releases (ReleaseCounter = 1)
	if len(repos) != 1 {
		t.Fatalf("Expected 1 repo with releases, got %d", len(repos))
	}

	if repos[0].Name != "repo1" {
		t.Errorf("Expected repo name repo1, got %s", repos[0].Name)
	}
}

func TestGetReleases_GitHub_DraftAndPrerelease(t *testing.T) {
	mockData := `{"id": 1, "tag_name": "v1.0.0", "name": "Draft Release", "body": "Draft release", "url": "https://api.github.com/repos/testuser/testrepo/releases/1", "html_url": "https://github.com/testuser/testrepo/releases/tag/v1.0.0", "tarball_url": "https://api.github.com/repos/testuser/testrepo/tarball/v1.0.0", "zipball_url": "https://api.github.com/repos/testuser/testrepo/zipball/v1.0.0", "draft": true, "prerelease": true, "created_at": "2023-01-01T00:00:00Z", "published_at": "2023-01-01T00:00:00Z", "author": {"login": "testuser", "id": 1, "type": "User"}, "assets": []}`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL,
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   true,
		Provider: "github",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("Expected 1 release, got %d", len(releases))
	}

	if !releases[0].Draft {
		t.Errorf("Expected draft to be true")
	}

	if !releases[0].Prerelease {
		t.Errorf("Expected prerelease to be true")
	}
}
