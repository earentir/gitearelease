package providers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GiteaProvider implements the Provider interface for Gitea instances
type GiteaProvider struct{}

// NewGiteaProvider creates a new Gitea provider instance
func NewGiteaProvider() *GiteaProvider {
	return &GiteaProvider{}
}

// GetReleasesURL constructs the Gitea API URL for fetching releases
func (p *GiteaProvider) GetReleasesURL(baseURL, user, repo string, latest bool) string {
	releaseType := "releases"
	if latest {
		releaseType = "releases/latest"
	}
	return fmt.Sprintf("%s/api/v1/repos/%s/%s/%s", baseURL, user, repo, releaseType)
}

// GetRepositoriesURL constructs the Gitea API URL for fetching repositories
func (p *GiteaProvider) GetRepositoriesURL(baseURL, user string) string {
	return fmt.Sprintf("%s/api/v1/users/%s/repos", baseURL, user)
}

// NormalizeRelease converts Gitea JSON to the standard Release struct
func (p *GiteaProvider) NormalizeRelease(data []byte, latest bool) ([]Release, error) {
	// Gitea JSON matches our Release structure, but we need to handle the conversion
	type giteaRelease struct {
		ID          int     `json:"id"`
		TagName     string  `json:"tag_name"`
		Name        string  `json:"name"`
		Body        string  `json:"body"`
		URL         string  `json:"url"`
		HTMLUrl     string  `json:"html_url"`
		TarballURL  string  `json:"tarball_url"`
		ZipballURL  string  `json:"zipball_url"`
		Draft       bool    `json:"draft"`
		Prerelease  bool    `json:"prerelease"`
		CreatedAt   string  `json:"created_at"`
		PublishedAt string  `json:"published_at"`
		Author      Author  `json:"author"`
		Assets      []Asset `json:"assets"`
	}

	var releases []Release

	if latest {
		var giteaRel giteaRelease
		if err := json.Unmarshal(data, &giteaRel); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
		rel := Release{
			ID:          giteaRel.ID,
			TagName:     giteaRel.TagName,
			Name:        giteaRel.Name,
			Body:        strings.ReplaceAll(giteaRel.Body, "\n", " "),
			URL:         giteaRel.URL,
			HTMLUrl:     giteaRel.HTMLUrl,
			TarballURL:  giteaRel.TarballURL,
			ZipballURL:  giteaRel.ZipballURL,
			Draft:       giteaRel.Draft,
			Prerelease:  giteaRel.Prerelease,
			CreatedAt:   giteaRel.CreatedAt,
			PublishedAt: giteaRel.PublishedAt,
			Author:      giteaRel.Author,
			Assets:      giteaRel.Assets,
		}
		releases = append(releases, rel)
		return releases, nil
	}

	var giteaReleases []giteaRelease
	if err := json.Unmarshal(data, &giteaReleases); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	for _, giteaRel := range giteaReleases {
		rel := Release{
			ID:          giteaRel.ID,
			TagName:     giteaRel.TagName,
			Name:        giteaRel.Name,
			Body:        strings.ReplaceAll(giteaRel.Body, "\n", " "),
			URL:         giteaRel.URL,
			HTMLUrl:     giteaRel.HTMLUrl,
			TarballURL:  giteaRel.TarballURL,
			ZipballURL:  giteaRel.ZipballURL,
			Draft:       giteaRel.Draft,
			Prerelease:  giteaRel.Prerelease,
			CreatedAt:   giteaRel.CreatedAt,
			PublishedAt: giteaRel.PublishedAt,
			Author:      giteaRel.Author,
			Assets:      giteaRel.Assets,
		}
		releases = append(releases, rel)
	}

	return releases, nil
}

// NormalizeRepositories converts Gitea JSON to the standard Repository slice
func (p *GiteaProvider) NormalizeRepositories(data []byte) ([]Repository, error) {
	// For Gitea, we'll need to parse and convert - this is complex, so we'll use a simpler approach
	// by parsing the JSON and mapping fields
	var repos []Repository
	type giteaRepo struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		FullName        string `json:"full_name"`
		Description     string `json:"description"`
		Private         bool   `json:"private"`
		Fork            bool   `json:"fork"`
		Size            int    `json:"size"`
		Language        string `json:"language"`
		HTMLURL         string `json:"html_url"`
		CloneURL        string `json:"clone_url"`
		SSHURL          string `json:"ssh_url"`
		StarsCount      int    `json:"stars_count"`
		ForksCount      int    `json:"forks_count"`
		WatchersCount   int    `json:"watchers_count"`
		OpenIssuesCount int    `json:"open_issues_count"`
		ReleaseCounter  int    `json:"release_counter"`
		DefaultBranch   string `json:"default_branch"`
		Archived        bool   `json:"archived"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
		Owner           struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Username  string `json:"username"`
			FullName  string `json:"full_name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		} `json:"owner"`
		Permissions struct {
			Admin bool `json:"admin"`
			Push  bool `json:"push"`
			Pull  bool `json:"pull"`
		} `json:"permissions"`
		HasIssues   bool `json:"has_issues"`
		HasWiki     bool `json:"has_wiki"`
		HasProjects bool `json:"has_projects"`
		HasReleases bool `json:"has_releases"`
		HasPackages bool `json:"has_packages"`
	}

	var giteaRepos []giteaRepo
	if err := json.Unmarshal(data, &giteaRepos); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	for _, gr := range giteaRepos {
		repos = append(repos, Repository{
			ID:              gr.ID,
			Name:            gr.Name,
			FullName:        gr.FullName,
			Description:     gr.Description,
			Private:         gr.Private,
			Fork:            gr.Fork,
			Size:            gr.Size,
			Language:        gr.Language,
			HTMLURL:         gr.HTMLURL,
			CloneURL:        gr.CloneURL,
			SSHURL:          gr.SSHURL,
			StarsCount:      gr.StarsCount,
			ForksCount:      gr.ForksCount,
			WatchersCount:   gr.WatchersCount,
			OpenIssuesCount: gr.OpenIssuesCount,
			ReleaseCounter:  gr.ReleaseCounter,
			DefaultBranch:   gr.DefaultBranch,
			Archived:        gr.Archived,
			CreatedAt:       gr.CreatedAt,
			UpdatedAt:       gr.UpdatedAt,
			Owner: Owner{
				ID:        gr.Owner.ID,
				Login:     gr.Owner.Login,
				Username:  gr.Owner.Username,
				FullName:  gr.Owner.FullName,
				Email:     gr.Owner.Email,
				AvatarURL: gr.Owner.AvatarURL,
			},
			Permissions: Permissions{
				Admin: gr.Permissions.Admin,
				Push:  gr.Permissions.Push,
				Pull:  gr.Permissions.Pull,
			},
			HasIssues:   gr.HasIssues,
			HasWiki:     gr.HasWiki,
			HasProjects: gr.HasProjects,
			HasReleases: gr.HasReleases,
			HasPackages: gr.HasPackages,
		})
	}

	return repos, nil
}

// DetectProvider checks if the baseURL is a Gitea instance
func (p *GiteaProvider) DetectProvider(baseURL string) bool {
	// Gitea instances typically have /api/v1 in the path or are explicitly not github/gitlab
	// Default provider, so we return true if it doesn't match other providers
	lowerURL := strings.ToLower(baseURL)
	return !strings.Contains(lowerURL, "github.com") &&
		!strings.Contains(lowerURL, "gitlab.com") &&
		!strings.Contains(lowerURL, "api.github.com") &&
		!strings.Contains(lowerURL, "gitlab")
}
