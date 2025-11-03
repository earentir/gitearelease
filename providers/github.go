package providers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GitHubProvider implements the Provider interface for GitHub
type GitHubProvider struct{}

// NewGitHubProvider creates a new GitHub provider instance
func NewGitHubProvider() *GitHubProvider {
	return &GitHubProvider{}
}

// GetReleasesURL constructs the GitHub API URL for fetching releases
func (p *GitHubProvider) GetReleasesURL(baseURL, user, repo string, latest bool) string {
	if latest {
		return fmt.Sprintf("%s/repos/%s/%s/releases/latest", baseURL, user, repo)
	}
	return fmt.Sprintf("%s/repos/%s/%s/releases", baseURL, user, repo)
}

// GetRepositoriesURL constructs the GitHub API URL for fetching repositories
func (p *GitHubProvider) GetRepositoriesURL(baseURL, user string) string {
	return fmt.Sprintf("%s/users/%s/repos", baseURL, user)
}

// githubRelease represents GitHub's release JSON structure
type githubRelease struct {
	ID          int    `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	TarballURL  string `json:"tarball_url"`
	ZipballURL  string `json:"zipball_url"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Author      struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
		Type  string `json:"type"`
	} `json:"author"`
	Assets []struct {
		ID                 int    `json:"id"`
		Name               string `json:"name"`
		Size               int64  `json:"size"`
		DownloadCount      int    `json:"download_count"`
		CreatedAt          string `json:"created_at"`
		BrowserDownloadURL string `json:"browser_download_url"`
		ContentType        string `json:"content_type"`
	} `json:"assets"`
}

// NormalizeRelease converts GitHub JSON to the standard Release struct
func (p *GitHubProvider) NormalizeRelease(data []byte, latest bool) ([]Release, error) {
	var releases []Release

	if latest {
		var ghRel githubRelease
		if err := json.Unmarshal(data, &ghRel); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
		rel := p.convertGitHubRelease(ghRel)
		releases = append(releases, rel)
		return releases, nil
	}

	var ghReleases []githubRelease
	if err := json.Unmarshal(data, &ghReleases); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	for _, ghRel := range ghReleases {
		rel := p.convertGitHubRelease(ghRel)
		releases = append(releases, rel)
	}

	return releases, nil
}

// convertGitHubRelease converts a GitHub release to the standard Release struct
func (p *GitHubProvider) convertGitHubRelease(ghRel githubRelease) Release {
	rel := Release{
		ID:          ghRel.ID,
		TagName:     ghRel.TagName,
		Name:        ghRel.Name,
		Body:        strings.ReplaceAll(ghRel.Body, "\n", " "),
		URL:         ghRel.URL,
		HTMLUrl:     ghRel.HTMLURL,
		TarballURL:  ghRel.TarballURL,
		ZipballURL:  ghRel.ZipballURL,
		Draft:       ghRel.Draft,
		Prerelease:  ghRel.Prerelease,
		CreatedAt:   ghRel.CreatedAt,
		PublishedAt: ghRel.PublishedAt,
		Author: Author{
			Login:    ghRel.Author.Login,
			Username: ghRel.Author.Login,
		},
		Assets: make([]Asset, len(ghRel.Assets)),
	}

	for i, ghAsset := range ghRel.Assets {
		rel.Assets[i] = Asset{
			ID:                 ghAsset.ID,
			Name:               ghAsset.Name,
			Size:               ghAsset.Size,
			DownloadCount:      ghAsset.DownloadCount,
			CreatedAt:          ghAsset.CreatedAt,
			BrowserDownloadURL: ghAsset.BrowserDownloadURL,
			Type:               ghAsset.ContentType,
		}
	}

	return rel
}

// githubRepository represents GitHub's repository JSON structure
type githubRepository struct {
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
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	WatchersCount   int    `json:"watchers_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	DefaultBranch   string `json:"default_branch"`
	Archived        bool   `json:"archived"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Owner           struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Type      string `json:"type"`
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

// NormalizeRepositories converts GitHub JSON to the standard Repository slice
func (p *GitHubProvider) NormalizeRepositories(data []byte) ([]Repository, error) {
	var ghRepos []githubRepository
	if err := json.Unmarshal(data, &ghRepos); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	repos := make([]Repository, len(ghRepos))
	for i, ghRepo := range ghRepos {
		repos[i] = p.convertGitHubRepository(ghRepo)
	}

	return repos, nil
}

// convertGitHubRepository converts a GitHub repository to the standard Repository struct
func (p *GitHubProvider) convertGitHubRepository(ghRepo githubRepository) Repository {
	createdAt, _ := time.Parse(time.RFC3339, ghRepo.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, ghRepo.UpdatedAt)

	repo := Repository{
		ID:              ghRepo.ID,
		Name:            ghRepo.Name,
		FullName:        ghRepo.FullName,
		Description:     ghRepo.Description,
		Private:         ghRepo.Private,
		Fork:            ghRepo.Fork,
		Size:            ghRepo.Size,
		Language:        ghRepo.Language,
		HTMLURL:         ghRepo.HTMLURL,
		CloneURL:        ghRepo.CloneURL,
		SSHURL:          ghRepo.SSHURL,
		StarsCount:      ghRepo.StargazersCount,
		ForksCount:      ghRepo.ForksCount,
		WatchersCount:   ghRepo.WatchersCount,
		OpenIssuesCount: ghRepo.OpenIssuesCount,
		DefaultBranch:   ghRepo.DefaultBranch,
		Archived:        ghRepo.Archived,
		CreatedAt:       createdAt.Format(time.RFC3339),
		UpdatedAt:       updatedAt.Format(time.RFC3339),
		HasIssues:       ghRepo.HasIssues,
		HasWiki:         ghRepo.HasWiki,
		HasProjects:     ghRepo.HasProjects,
		HasReleases:     ghRepo.HasReleases,
		HasPackages:     ghRepo.HasPackages,
	}

	repo.Owner = Owner{
		ID:        ghRepo.Owner.ID,
		Login:     ghRepo.Owner.Login,
		Username:  ghRepo.Owner.Login,
		AvatarURL: ghRepo.Owner.AvatarURL,
	}

	repo.Permissions = Permissions{
		Admin: ghRepo.Permissions.Admin,
		Push:  ghRepo.Permissions.Push,
		Pull:  ghRepo.Permissions.Pull,
	}

	// Note: GitHub doesn't have ReleaseCounter in the standard API response
	// We would need to fetch releases separately (GET /repos/{owner}/{repo}/releases)
	// For now, we'll set it based on HasReleases as a placeholder
	// Value of 1 indicates releases exist, but is not an accurate count
	if ghRepo.HasReleases {
		repo.ReleaseCounter = 1 // Placeholder - not accurate count
	} else {
		repo.ReleaseCounter = 0 // No releases
	}

	return repo
}

// DetectProvider checks if the baseURL is GitHub
func (p *GitHubProvider) DetectProvider(baseURL string) bool {
	lowerURL := strings.ToLower(baseURL)
	return strings.Contains(lowerURL, "github.com") || strings.Contains(lowerURL, "api.github.com")
}
