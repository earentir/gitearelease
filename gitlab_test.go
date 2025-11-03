package gitearelease

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetReleases_GitLab_Success(t *testing.T) {
	// GitLab JSON format
	mockData := `[{"tag_name": "v1.0.0", "name": "Release 1.0.0", "description": "Test release", "created_at": "2023-01-01T00:00:00Z", "released_at": "2023-01-01T00:00:00Z", "author": {"id": 1, "username": "testuser", "name": "Test User", "email": "test@example.com", "avatar_url": "https://gitlab.com/testuser.png"}, "commit": {"id": "abc123", "short_id": "abc123", "title": "Initial commit", "created_at": "2023-01-01T00:00:00Z", "message": "Initial commit", "author_name": "Test User", "author_email": "test@example.com"}, "milestones": [], "commit_path": "/testuser/testrepo/-/commit/abc123", "tag_path": "/testuser/testrepo/-/tags/v1.0.0", "assets": {"count": 1, "links": [{"id": 1, "name": "binary.tar.gz", "url": "https://gitlab.com/testuser/testrepo/-/releases/v1.0.0/downloads/binary.tar.gz", "direct_asset_url": "https://gitlab.com/testuser/testrepo/-/releases/v1.0.0/downloads/binary.tar.gz", "link_type": "other"}], "sources": [{"format": "tar.gz", "url": "https://gitlab.com/testuser/testrepo/-/archive/v1.0.0/testrepo-v1.0.0.tar.gz"}]}, "evidences": [], "_links": {"self": "https://gitlab.com/api/v4/projects/testuser%2Ftestrepo/releases/v1.0.0", "edit_url": "https://gitlab.com/testuser/testrepo/-/releases/v1.0.0/edit"}}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GitLab doesn't have /latest endpoint, so it fetches all and takes first
		// normalizeBaseURL adds /api/v4, then GetReleasesURL adds /projects/...
		// Note: URL may decode %2F to /, so we check for either
		if !strings.Contains(r.URL.Path, "/projects/") || !strings.Contains(r.URL.Path, "/releases") {
			t.Errorf("Expected path containing /projects/.../releases, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL, // Server URL without /api/v4 (will be added by normalizeBaseURL)
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   true,
		Provider: "gitlab",
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

	// GitLab doesn't support Draft/Prerelease flags
	if rel.Draft != false {
		t.Errorf("Expected draft to be false for GitLab, got %v", rel.Draft)
	}

	if rel.Prerelease != false {
		t.Errorf("Expected prerelease to be false for GitLab, got %v", rel.Prerelease)
	}

	// Check assets - GitLab has limited asset information
	if len(rel.Assets) != 1 {
		t.Fatalf("Expected 1 asset, got %d", len(rel.Assets))
	}

	asset := rel.Assets[0]
	if asset.Name != "binary.tar.gz" {
		t.Errorf("Expected asset name binary.tar.gz, got %s", asset.Name)
	}

	// GitLab assets don't have Size, DownloadCount, CreatedAt, UUID
	if asset.Size != 0 {
		t.Errorf("Expected asset size 0 (unavailable), got %d", asset.Size)
	}

	if asset.DownloadCount != 0 {
		t.Errorf("Expected download count 0 (unavailable), got %d", asset.DownloadCount)
	}

	if asset.UUID != "" {
		t.Errorf("Expected empty UUID for GitLab, got %s", asset.UUID)
	}

	// Check tarball URL from sources
	if rel.TarballURL == "" {
		t.Errorf("Expected tarball URL to be populated from sources")
	}
}

func TestGetReleases_GitLab_All(t *testing.T) {
	mockData := `[{"tag_name": "v1.0.0", "name": "Release 1.0.0", "description": "Test release 1", "created_at": "2023-01-01T00:00:00Z", "released_at": "2023-01-01T00:00:00Z", "author": {"id": 1, "username": "testuser", "name": "Test User", "email": "test@example.com", "avatar_url": "https://gitlab.com/testuser.png"}, "commit": {"id": "abc123", "short_id": "abc123", "title": "Initial commit", "created_at": "2023-01-01T00:00:00Z", "message": "Initial commit", "author_name": "Test User", "author_email": "test@example.com"}, "milestones": [], "commit_path": "/testuser/testrepo/-/commit/abc123", "tag_path": "/testuser/testrepo/-/tags/v1.0.0", "assets": {"count": 0, "links": [], "sources": []}, "evidences": [], "_links": {"self": "https://gitlab.com/api/v4/projects/testuser%2Ftestrepo/releases/v1.0.0", "edit_url": "https://gitlab.com/testuser/testrepo/-/releases/v1.0.0/edit"}}, {"tag_name": "v1.1.0", "name": "Release 1.1.0", "description": "Test release 2", "created_at": "2023-01-02T00:00:00Z", "released_at": "2023-01-02T00:00:00Z", "author": {"id": 1, "username": "testuser", "name": "Test User", "email": "test@example.com", "avatar_url": "https://gitlab.com/testuser.png"}, "commit": {"id": "def456", "short_id": "def456", "title": "Second commit", "created_at": "2023-01-02T00:00:00Z", "message": "Second commit", "author_name": "Test User", "author_email": "test@example.com"}, "milestones": [], "commit_path": "/testuser/testrepo/-/commit/def456", "tag_path": "/testuser/testrepo/-/tags/v1.1.0", "assets": {"count": 0, "links": [], "sources": []}, "evidences": [], "_links": {"self": "https://gitlab.com/api/v4/projects/testuser%2Ftestrepo/releases/v1.1.0", "edit_url": "https://gitlab.com/testuser/testrepo/-/releases/v1.1.0/edit"}}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL,
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   false,
		Provider: "gitlab",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(releases) != 2 {
		t.Fatalf("Expected 2 releases, got %d", len(releases))
	}

	if releases[1].TagName != "v1.1.0" {
		t.Errorf("Expected second release tag v1.1.0, got %s", releases[1].TagName)
	}
}

func TestGetRepositories_GitLab_Success(t *testing.T) {
	// GitLab JSON format
	mockData := `[{"id": 1, "name": "repo1", "path": "repo1", "path_with_namespace": "testuser/repo1", "description": "Test repo 1", "visibility": "public", "fork": false, "size": 1024000, "language": "Go", "web_url": "https://gitlab.com/testuser/repo1", "ssh_url_to_repo": "git@gitlab.com:testuser/repo1.git", "http_url_to_repo": "https://gitlab.com/testuser/repo1.git", "star_count": 10, "forks_count": 5, "open_issues_count": 2, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "last_activity_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "username": "testuser", "name": "Test User", "avatar_url": "https://gitlab.com/testuser.png"}, "permissions": {"project_access": {"access_level": 40}, "group_access": null}}, {"id": 2, "name": "repo2", "path": "repo2", "path_with_namespace": "testuser/repo2", "description": "Test repo 2", "visibility": "private", "fork": false, "size": 2048000, "language": "Python", "web_url": "https://gitlab.com/testuser/repo2", "ssh_url_to_repo": "git@gitlab.com:testuser/repo2.git", "http_url_to_repo": "https://gitlab.com/testuser/repo2.git", "star_count": 20, "forks_count": 10, "open_issues_count": 0, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "last_activity_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "username": "testuser", "name": "Test User", "avatar_url": "https://gitlab.com/testuser.png"}, "permissions": {"project_access": {"access_level": 30}, "group_access": null}}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// normalizeBaseURL adds /api/v4, then GetRepositoriesURL adds /users/...
		expectedPath := "/api/v4/users/testuser/projects"
		if !strings.HasSuffix(r.URL.Path, expectedPath) {
			t.Errorf("Expected path ending with %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL, // Server URL without /api/v4 (will be added by normalizeBaseURL)
		User:         "testuser",
		WithReleases: false,
		Provider:     "gitlab",
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

	// GitLab ReleaseCounter is always 0 (not available)
	if repos[0].ReleaseCounter != 0 {
		t.Errorf("Expected ReleaseCounter 0 (unavailable), got %d", repos[0].ReleaseCounter)
	}

	// Check private field
	if repos[0].Private != false {
		t.Errorf("Expected private false (public visibility), got %v", repos[0].Private)
	}

	if repos[1].Private != true {
		t.Errorf("Expected private true (private visibility), got %v", repos[1].Private)
	}

	// GitLab doesn't provide HasIssues, HasWiki, HasProjects, HasPackages
	if repos[0].HasIssues != false {
		t.Errorf("Expected HasIssues false (unavailable), got %v", repos[0].HasIssues)
	}

	if repos[0].HasWiki != false {
		t.Errorf("Expected HasWiki false (unavailable), got %v", repos[0].HasWiki)
	}

	if repos[0].HasProjects != false {
		t.Errorf("Expected HasProjects false (unavailable), got %v", repos[0].HasProjects)
	}

	if repos[0].HasPackages != false {
		t.Errorf("Expected HasPackages false (unavailable), got %v", repos[0].HasPackages)
	}

	// Check permissions (GitLab uses access levels)
	// Access level 40 = Maintainer, which is not >= 50 (Owner), so admin should be false
	if repos[0].Permissions.Admin {
		t.Errorf("Expected admin false (access_level 40 = Maintainer, not Owner >= 50), got %v", repos[0].Permissions.Admin)
	}

	if !repos[0].Permissions.Push {
		t.Errorf("Expected push true (access_level 40 >= 30), got %v", repos[0].Permissions.Push)
	}
}

func TestGetRepositories_GitLab_WithReleases(t *testing.T) {
	mockData := `[{"id": 1, "name": "repo1", "path": "repo1", "path_with_namespace": "testuser/repo1", "description": "Test repo 1", "visibility": "public", "fork": false, "size": 1024000, "language": "Go", "web_url": "https://gitlab.com/testuser/repo1", "ssh_url_to_repo": "git@gitlab.com:testuser/repo1.git", "http_url_to_repo": "https://gitlab.com/testuser/repo1.git", "star_count": 10, "forks_count": 5, "open_issues_count": 2, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "last_activity_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "username": "testuser", "name": "Test User", "avatar_url": "https://gitlab.com/testuser.png"}, "permissions": {"project_access": {"access_level": 40}, "group_access": null}}, {"id": 2, "name": "repo2", "path": "repo2", "path_with_namespace": "testuser/repo2", "description": "Test repo 2", "visibility": "private", "fork": false, "size": 2048000, "language": "Python", "web_url": "https://gitlab.com/testuser/repo2", "ssh_url_to_repo": "git@gitlab.com:testuser/repo2.git", "http_url_to_repo": "https://gitlab.com/testuser/repo2.git", "star_count": 20, "forks_count": 10, "open_issues_count": 0, "default_branch": "main", "archived": false, "created_at": "2023-01-01T00:00:00Z", "last_activity_at": "2023-01-02T00:00:00Z", "owner": {"id": 1, "username": "testuser", "name": "Test User", "avatar_url": "https://gitlab.com/testuser.png"}, "permissions": {"project_access": {"access_level": 30}, "group_access": null}}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	repositoriesToFetch := RepositoriesToFetch{
		BaseURL:      mockServer.URL,
		User:         "testuser",
		WithReleases: true,
		Provider:     "gitlab",
	}

	repos, err := GetRepositories(repositoriesToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// GitLab ReleaseCounter is always 0, so WithReleases filter won't match any
	// (unless we manually set ReleaseCounter > 0, which we don't)
	if len(repos) != 0 {
		t.Fatalf("Expected 0 repos with releases (GitLab ReleaseCounter always 0), got %d", len(repos))
	}
}

func TestGetReleases_GitLab_NoLatestEndpoint(t *testing.T) {
	// Test that GitLab fetches all releases when Latest=true
	// (since GitLab doesn't have a /latest endpoint)
	mockData := `[{"tag_name": "v1.0.0", "name": "Release 1.0.0", "description": "Latest release", "created_at": "2023-01-02T00:00:00Z", "released_at": "2023-01-02T00:00:00Z", "author": {"id": 1, "username": "testuser", "name": "Test User", "email": "test@example.com", "avatar_url": "https://gitlab.com/testuser.png"}, "commit": {"id": "abc123", "short_id": "abc123", "title": "Latest commit", "created_at": "2023-01-02T00:00:00Z", "message": "Latest commit", "author_name": "Test User", "author_email": "test@example.com"}, "milestones": [], "commit_path": "/testuser/testrepo/-/commit/abc123", "tag_path": "/testuser/testrepo/-/tags/v1.0.0", "assets": {"count": 0, "links": [], "sources": []}, "evidences": [], "_links": {"self": "https://gitlab.com/api/v4/projects/testuser%2Ftestrepo/releases/v1.0.0", "edit_url": "https://gitlab.com/testuser/testrepo/-/releases/v1.0.0/edit"}}, {"tag_name": "v0.9.0", "name": "Release 0.9.0", "description": "Older release", "created_at": "2023-01-01T00:00:00Z", "released_at": "2023-01-01T00:00:00Z", "author": {"id": 1, "username": "testuser", "name": "Test User", "email": "test@example.com", "avatar_url": "https://gitlab.com/testuser.png"}, "commit": {"id": "def456", "short_id": "def456", "title": "Older commit", "created_at": "2023-01-01T00:00:00Z", "message": "Older commit", "author_name": "Test User", "author_email": "test@example.com"}, "milestones": [], "commit_path": "/testuser/testrepo/-/commit/def456", "tag_path": "/testuser/testrepo/-/tags/v0.9.0", "assets": {"count": 0, "links": [], "sources": []}, "evidences": [], "_links": {"self": "https://gitlab.com/api/v4/projects/testuser%2Ftestrepo/releases/v0.9.0", "edit_url": "https://gitlab.com/testuser/testrepo/-/releases/v0.9.0/edit"}}]`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should still call the same endpoint (no /latest)
		// normalizeBaseURL adds /api/v4, then GetReleasesURL adds /projects/...
		// Note: URL may decode %2F to /, so we check for either
		if !strings.Contains(r.URL.Path, "/projects/") || !strings.Contains(r.URL.Path, "/releases") {
			t.Errorf("Expected path containing /projects/.../releases, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockData))
	}))
	defer mockServer.Close()

	releaseToFetch := ReleaseToFetch{
		BaseURL:  mockServer.URL, // Server URL without /api/v4 (will be added by normalizeBaseURL)
		User:     "testuser",
		Repo:     "testrepo",
		Latest:   true,
		Provider: "gitlab",
	}

	releases, err := GetReleases(releaseToFetch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should only return the first release (latest)
	if len(releases) != 1 {
		t.Fatalf("Expected 1 release (latest), got %d", len(releases))
	}

	// Should be v1.0.0 (first in array, which is treated as latest)
	if releases[0].TagName != "v1.0.0" {
		t.Errorf("Expected tag v1.0.0, got %s", releases[0].TagName)
	}
}
