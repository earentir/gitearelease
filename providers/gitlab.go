package providers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GitLabProvider implements the Provider interface for GitLab
type GitLabProvider struct{}

// NewGitLabProvider creates a new GitLab provider instance
func NewGitLabProvider() *GitLabProvider {
	return &GitLabProvider{}
}

// GetReleasesURL constructs the GitLab API URL for fetching releases
func (p *GitLabProvider) GetReleasesURL(baseURL, user, repo string, latest bool) string {
	// GitLab uses project path encoding (owner%2Frepo)
	// Use PathEscape to properly encode the slash
	projectPath := fmt.Sprintf("%s%%2F%s", user, repo)
	// baseURL already includes /api/v4 from normalizeBaseURL
	if latest {
		// GitLab doesn't have a /latest endpoint, we'll fetch all and take first
		return fmt.Sprintf("%s/projects/%s/releases", baseURL, projectPath)
	}
	return fmt.Sprintf("%s/projects/%s/releases", baseURL, projectPath)
}

// GetRepositoriesURL constructs the GitLab API URL for fetching repositories
func (p *GitLabProvider) GetRepositoriesURL(baseURL, user string) string {
	// baseURL already includes /api/v4 from normalizeBaseURL
	return fmt.Sprintf("%s/users/%s/projects", baseURL, user)
}

// gitlabRelease represents GitLab's release JSON structure
type gitlabRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	ReleasedAt  string `json:"released_at"`
	Author      struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	} `json:"author"`
	Commit struct {
		ID          string `json:"id"`
		ShortID     string `json:"short_id"`
		Title       string `json:"title"`
		CreatedAt   string `json:"created_at"`
		Message     string `json:"message"`
		AuthorName  string `json:"author_name"`
		AuthorEmail string `json:"author_email"`
	} `json:"commit"`
	Milestones []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		State       string `json:"state"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		DueDate     string `json:"due_date"`
		StartDate   string `json:"start_date"`
	} `json:"milestones"`
	CommitPath string              `json:"commit_path"`
	TagPath    string              `json:"tag_path"`
	Assets     gitlabReleaseAssets `json:"assets"`
	Evidences  []interface{}       `json:"evidences"`
	Links      struct {
		Self    string `json:"self"`
		EditURL string `json:"edit_url"`
	} `json:"_links"`
}

type gitlabReleaseAssets struct {
	Count int `json:"count"`
	Links []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		URL            string `json:"url"`
		DirectAssetURL string `json:"direct_asset_url"`
		LinkType       string `json:"link_type"`
	} `json:"links"`
	Sources []struct {
		Format string `json:"format"`
		URL    string `json:"url"`
	} `json:"sources"`
}

// NormalizeRelease converts GitLab JSON to the standard Release struct
func (p *GitLabProvider) NormalizeRelease(data []byte, latest bool) ([]Release, error) {
	var releases []Release

	if latest {
		var gitlabReleases []gitlabRelease
		if err := json.Unmarshal(data, &gitlabReleases); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
		if len(gitlabReleases) == 0 {
			return releases, nil
		}
		// GitLab doesn't have a /latest endpoint, so we take the first one
		rel := p.convertGitLabRelease(gitlabReleases[0])
		releases = append(releases, rel)
		return releases, nil
	}

	var gitlabReleases []gitlabRelease
	if err := json.Unmarshal(data, &gitlabReleases); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	for _, glRel := range gitlabReleases {
		rel := p.convertGitLabRelease(glRel)
		releases = append(releases, rel)
	}

	return releases, nil
}

// convertGitLabRelease converts a GitLab release to the standard Release struct
func (p *GitLabProvider) convertGitLabRelease(glRel gitlabRelease) Release {
	// Generate a simple ID from tag name hash (GitLab doesn't provide a numeric ID)
	id := 0
	for _, char := range glRel.TagName {
		id = id*31 + int(char)
	}
	if id < 0 {
		id = -id
	}

	rel := Release{
		ID:          id,
		TagName:     glRel.TagName,
		Name:        glRel.Name,
		Body:        strings.ReplaceAll(glRel.Description, "\n", " "),
		URL:         glRel.Links.Self,
		HTMLUrl:     glRel.TagPath,
		TarballURL:  "",    // GitLab uses different structure
		ZipballURL:  "",    // GitLab uses different structure
		Draft:       false, // GitLab doesn't have draft releases in the same way
		Prerelease:  false, // GitLab doesn't mark prereleases the same way
		CreatedAt:   glRel.CreatedAt,
		PublishedAt: glRel.ReleasedAt,
		Author: Author{
			Login:    glRel.Author.Username,
			Username: glRel.Author.Username,
			FullName: glRel.Author.Name,
			Email:    glRel.Author.Email,
		},
		Assets: make([]Asset, len(glRel.Assets.Links)),
	}

	// Convert GitLab assets
	// Note: GitLab API doesn't provide Size, DownloadCount, CreatedAt, or UUID for assets
	for i, glAsset := range glRel.Assets.Links {
		rel.Assets[i] = Asset{
			ID:                 glAsset.ID,
			Name:               glAsset.Name,
			BrowserDownloadURL: glAsset.DirectAssetURL,
			Type:               glAsset.LinkType,
			// Size, DownloadCount, CreatedAt, UUID are not available in GitLab API
			Size:          0,
			DownloadCount: 0,
			CreatedAt:     "",
			UUID:          "",
		}
	}

	// Add source links as tarball/zipball if available
	for _, source := range glRel.Assets.Sources {
		if source.Format == "tar.gz" || source.Format == "tar" {
			rel.TarballURL = source.URL
		} else if source.Format == "zip" {
			rel.ZipballURL = source.URL
		}
	}

	return rel
}

// gitlabRepository represents GitLab's repository JSON structure
type gitlabRepository struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	Description       string `json:"description"`
	Visibility        string `json:"visibility"`
	Fork              bool   `json:"fork"`
	Size              int64  `json:"size"`
	Language          string `json:"language"`
	WebURL            string `json:"web_url"`
	SSHURL            string `json:"ssh_url_to_repo"`
	HTTPURL           string `json:"http_url_to_repo"`
	StarCount         int    `json:"star_count"`
	ForksCount        int    `json:"forks_count"`
	OpenIssuesCount   int    `json:"open_issues_count"`
	DefaultBranch     string `json:"default_branch"`
	Archived          bool   `json:"archived"`
	CreatedAt         string `json:"created_at"`
	LastActivityAt    string `json:"last_activity_at"`
	Owner             struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	Permissions struct {
		ProjectAccess struct {
			AccessLevel int `json:"access_level"`
		} `json:"project_access"`
		GroupAccess struct {
			AccessLevel int `json:"access_level"`
		} `json:"group_access"`
	} `json:"permissions"`
}

// NormalizeRepositories converts GitLab JSON to the standard Repository slice
func (p *GitLabProvider) NormalizeRepositories(data []byte) ([]Repository, error) {
	var glRepos []gitlabRepository
	if err := json.Unmarshal(data, &glRepos); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	repos := make([]Repository, len(glRepos))
	for i, glRepo := range glRepos {
		repos[i] = p.convertGitLabRepository(glRepo)
	}

	return repos, nil
}

// convertGitLabRepository converts a GitLab repository to the standard Repository struct
func (p *GitLabProvider) convertGitLabRepository(glRepo gitlabRepository) Repository {
	createdAt, _ := time.Parse(time.RFC3339, glRepo.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, glRepo.LastActivityAt)

	// Convert size from bytes to KB (approximate)
	sizeKB := int(glRepo.Size / 1024)

	repo := Repository{
		ID:              glRepo.ID,
		Name:            glRepo.Name,
		FullName:        glRepo.PathWithNamespace,
		Description:     glRepo.Description,
		Private:         glRepo.Visibility == "private",
		Fork:            glRepo.Fork,
		Size:            sizeKB,
		Language:        glRepo.Language,
		HTMLURL:         glRepo.WebURL,
		CloneURL:        glRepo.HTTPURL,
		SSHURL:          glRepo.SSHURL,
		StarsCount:      glRepo.StarCount,
		ForksCount:      glRepo.ForksCount,
		OpenIssuesCount: glRepo.OpenIssuesCount,
		DefaultBranch:   glRepo.DefaultBranch,
		Archived:        glRepo.Archived,
		CreatedAt:       createdAt.Format(time.RFC3339),
		UpdatedAt:       updatedAt.Format(time.RFC3339),
	}

	repo.Owner = Owner{
		ID:        glRepo.Owner.ID,
		Login:     glRepo.Owner.Username,
		Username:  glRepo.Owner.Username,
		FullName:  glRepo.Owner.Name,
		AvatarURL: glRepo.Owner.AvatarURL,
	}

	// Convert permissions (GitLab uses access levels, we approximate)
	accessLevel := 0
	if glRepo.Permissions.ProjectAccess.AccessLevel > 0 {
		accessLevel = glRepo.Permissions.ProjectAccess.AccessLevel
	} else if glRepo.Permissions.GroupAccess.AccessLevel > 0 {
		accessLevel = glRepo.Permissions.GroupAccess.AccessLevel
	}
	// GitLab access levels: 10=Guest, 20=Reporter, 30=Developer, 40=Maintainer, 50=Owner
	repo.Permissions = Permissions{
		Admin: accessLevel >= 50,
		Push:  accessLevel >= 30,
		Pull:  accessLevel >= 10,
	}

	// Note: GitLab doesn't provide ReleaseCounter directly in the projects API
	// To get accurate count, would require fetching /projects/{id}/releases separately
	// Setting to 0 to indicate unavailability (not a real count)
	repo.ReleaseCounter = 0

	// Note: GitLab projects API doesn't provide these boolean flags:
	// HasIssues, HasWiki, HasProjects, HasPackages are not in the response
	// These fields will be false (zero value) for GitLab repositories
	repo.HasIssues = false
	repo.HasWiki = false
	repo.HasProjects = false
	repo.HasPackages = false

	return repo
}

// DetectProvider checks if the baseURL is GitLab
func (p *GitLabProvider) DetectProvider(baseURL string) bool {
	lowerURL := strings.ToLower(baseURL)
	return strings.Contains(lowerURL, "gitlab.com") || strings.Contains(lowerURL, "gitlab")
}
