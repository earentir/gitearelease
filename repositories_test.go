package gitearelease

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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
		BaseURL:    mockServer.URL, // Use the URL of the mock server
		User:       "testuser",
		WithReleas: false,
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
		BaseURL:    mockServer.URL, // Use the URL of the mock server
		User:       "testuser",
		WithReleas: false,
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
		BaseURL:    mockServer.URL, // Use the URL of the mock server
		User:       "testuser",
		WithReleas: true, // Only fetch repos with releases
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
