// Package gitearelease provides functions to interact with Git hosting platforms (Gitea, GitHub, GitLab) and fetch releases.
// The package automatically detects the provider from the BaseURL, but you can also explicitly specify it.
package gitearelease

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/earentir/gitearelease/providers"
)

// defaultHTTPTimeout is applied to every outbound HTTP operation made by
// this package unless the caller overrides it through SetHTTPTimeout.
const defaultHTTPTimeout = 15 * time.Second

// httpClient is shared by every helper in this package so that the timeout
// applies uniformly. You should rarely need to touch this directly.
var httpClient = &http.Client{Timeout: defaultHTTPTimeout}

// SetHTTPTimeout overrides the package‑level HTTP timeout. Pass zero or a
// negative value to restore the built‑in default (15 s). This call is safe
// to make at any time, even concurrently.
func SetHTTPTimeout(d time.Duration) {
	if d <= 0 {
		d = defaultHTTPTimeout
	}
	httpClient.Timeout = d
}

/* -------------------------------------------------------------------------- */
/*  PUBLIC API – identical function names & signatures as before              */
/* -------------------------------------------------------------------------- */

// GetReleases returns all releases or only the latest release from a repository.
// Supports Gitea, GitHub, and GitLab. Provider is auto-detected from BaseURL if not specified.
func GetReleases(r ReleaseToFetch) ([]Release, error) {
	// Determine provider - auto-detect if not specified
	providerType := providers.ProviderType(r.Provider)
	if providerType == "" {
		// Auto-detect from BaseURL
		provider := providers.GetProvider("", r.BaseURL)
		if provider.DetectProvider(r.BaseURL) {
			// Determine which provider was detected
			if providers.NewGitHubProvider().DetectProvider(r.BaseURL) {
				providerType = providers.ProviderGitHub
			} else if providers.NewGitLabProvider().DetectProvider(r.BaseURL) {
				providerType = providers.ProviderGitLab
			} else {
				providerType = providers.ProviderGitea // Default for backward compatibility
			}
		} else {
			providerType = providers.ProviderGitea // Default for backward compatibility
		}
	}

	// Normalize BaseURL for GitHub (add api.github.com if needed)
	baseURL := normalizeBaseURL(r.BaseURL, providerType)

	// Get the appropriate provider
	provider := providers.GetProvider(providerType, baseURL)

	// Construct API URL using provider
	apiURL := provider.GetReleasesURL(baseURL, r.User, r.Repo, r.Latest)

	// Fetch data
	apiData, err := fetchData(apiURL)
	if err != nil {
		return nil, err
	}

	// Normalize response using provider
	providerReleases, err := provider.NormalizeRelease(apiData, r.Latest)
	if err != nil {
		return nil, err
	}

	// Convert from provider types to main package types
	releases := make([]Release, len(providerReleases))
	for i, pr := range providerReleases {
		releases[i] = convertProviderRelease(pr)
	}

	return releases, nil
}

// DownloadBinary downloads a binary from a URL and saves it to a file.
func DownloadBinary(url, outputDir, filename string) (string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download binary: server returned %s", resp.Status)
	}

	outPath := filepath.Join(outputDir, filename)
	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("create file %q: %w", outPath, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("write file %q: %w", outPath, err)
	}
	return outPath, nil
}

// GetRepositories returns all repositories of a user and can filter by releases.
// Supports Gitea, GitHub, and GitLab. Provider is auto-detected from BaseURL if not specified.
func GetRepositories(r RepositoriesToFetch) ([]Repository, error) {
	// Determine provider - auto-detect if not specified
	providerType := providers.ProviderType(r.Provider)
	if providerType == "" {
		// Auto-detect from BaseURL
		if providers.NewGitHubProvider().DetectProvider(r.BaseURL) {
			providerType = providers.ProviderGitHub
		} else if providers.NewGitLabProvider().DetectProvider(r.BaseURL) {
			providerType = providers.ProviderGitLab
		} else {
			providerType = providers.ProviderGitea // Default for backward compatibility
		}
	}

	// Normalize BaseURL for GitHub (add api.github.com if needed)
	baseURL := normalizeBaseURL(r.BaseURL, providerType)

	// Get the appropriate provider
	provider := providers.GetProvider(providerType, baseURL)

	// Construct API URL using provider
	apiURL := provider.GetRepositoriesURL(baseURL, r.User)

	// Fetch data
	apiData, err := fetchData(apiURL)
	if err != nil {
		return nil, err
	}

	// Normalize response using provider
	providerRepos, err := provider.NormalizeRepositories(apiData)
	if err != nil {
		return nil, err
	}

	// Convert from provider types to main package types
	allRepos := make([]Repository, len(providerRepos))
	for i, pr := range providerRepos {
		allRepos[i] = convertProviderRepository(pr)
	}

	// Filter by releases if requested
	withReleases := r.WithReleases
	if !withReleases {
		return allRepos, nil
	}

	repos := make([]Repository, 0, len(allRepos))
	for _, repo := range allRepos {
		if repo.ReleaseCounter > 0 {
			repos = append(repos, repo)
		}
	}

	return repos, nil
}

// TrimVersionPrefix removes common version prefixes from a version string.
func TrimVersionPrefix(v string) string {
	v = strings.ToLower(v)
	for _, prefix := range []string{"v", "version", "ver", "release", "rel", "r", "v."} {
		v = strings.TrimPrefix(v, prefix)
	}
	return v
}

/* -------------------------------------------------------------------------- */
/*  INTERNALS                                                                 */
/* -------------------------------------------------------------------------- */

// convertProviderRelease converts a provider Release to the main package Release
func convertProviderRelease(pr providers.Release) Release {
	rel := Release{
		ID:          pr.ID,
		TagName:     pr.TagName,
		Name:        pr.Name,
		Body:        pr.Body,
		URL:         pr.URL,
		HTMLUrl:     pr.HTMLUrl,
		TarballURL:  pr.TarballURL,
		ZipballURL:  pr.ZipballURL,
		Draft:       pr.Draft,
		Prerelease:  pr.Prerelease,
		CreatedAt:   pr.CreatedAt,
		PublishedAt: pr.PublishedAt,
		Author: Author{
			Login:     pr.Author.Login,
			LoginName: pr.Author.LoginName,
			FullName:  pr.Author.FullName,
			Email:     pr.Author.Email,
			Username:  pr.Author.Username,
		},
		Assets: make([]Asset, len(pr.Assets)),
	}

	for i, pa := range pr.Assets {
		rel.Assets[i] = Asset{
			ID:                 pa.ID,
			Name:               pa.Name,
			Size:               pa.Size,
			DownloadCount:      pa.DownloadCount,
			CreatedAt:          pa.CreatedAt,
			UUID:               pa.UUID,
			BrowserDownloadURL: pa.BrowserDownloadURL,
			Type:               pa.Type,
		}
	}

	return rel
}

// convertProviderRepository converts a provider Repository to the main package Repository
func convertProviderRepository(pr providers.Repository) Repository {
	createdAt, _ := time.Parse(time.RFC3339, pr.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pr.UpdatedAt)

	repo := Repository{
		ID:              pr.ID,
		Name:            pr.Name,
		FullName:        pr.FullName,
		Description:     pr.Description,
		Private:         pr.Private,
		Fork:            pr.Fork,
		Size:            pr.Size,
		Language:        pr.Language,
		HTMLURL:         pr.HTMLURL,
		CloneURL:        pr.CloneURL,
		SSHURL:          pr.SSHURL,
		StarsCount:      pr.StarsCount,
		ForksCount:      pr.ForksCount,
		WatchersCount:   pr.WatchersCount,
		OpenIssuesCount: pr.OpenIssuesCount,
		ReleaseCounter:  pr.ReleaseCounter,
		DefaultBranch:   pr.DefaultBranch,
		Archived:        pr.Archived,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		HasIssues:       pr.HasIssues,
		HasWiki:         pr.HasWiki,
		HasProjects:     pr.HasProjects,
		HasReleases:     pr.HasReleases,
		HasPackages:     pr.HasPackages,
	}

	repo.Owner.ID = pr.Owner.ID
	repo.Owner.Login = pr.Owner.Login
	repo.Owner.Username = pr.Owner.Username
	repo.Owner.FullName = pr.Owner.FullName
	repo.Owner.Email = pr.Owner.Email
	repo.Owner.AvatarURL = pr.Owner.AvatarURL

	repo.Permissions.Admin = pr.Permissions.Admin
	repo.Permissions.Push = pr.Permissions.Push
	repo.Permissions.Pull = pr.Permissions.Pull

	return repo
}

// normalizeBaseURL ensures the BaseURL is properly formatted for the provider
func normalizeBaseURL(baseURL string, providerType providers.ProviderType) string {
	baseURL = strings.TrimSuffix(baseURL, "/")

	// For GitHub, if user provided github.com, convert to api.github.com
	if providerType == providers.ProviderGitHub {
		if strings.Contains(baseURL, "github.com") && !strings.Contains(baseURL, "api.github.com") {
			baseURL = strings.Replace(baseURL, "github.com", "api.github.com", 1)
		}
		// If no domain specified, assume api.github.com
		if baseURL == "" || baseURL == "github.com" {
			return "https://api.github.com"
		}
		// Ensure https:// prefix
		if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
			return "https://" + baseURL
		}
	}

	// For GitLab, ensure proper API path
	if providerType == providers.ProviderGitLab {
		if strings.Contains(baseURL, "gitlab.com") && !strings.Contains(baseURL, "/api/v4") {
			// gitlab.com -> gitlab.com/api/v4
			baseURL = strings.TrimSuffix(baseURL, "/api/v4")
		}
		// If no domain specified, assume gitlab.com
		if baseURL == "" || baseURL == "gitlab.com" {
			return "https://gitlab.com/api/v4"
		}
		// Ensure https:// prefix
		if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
			return "https://" + baseURL
		}
		// Ensure /api/v4 suffix for GitLab
		if !strings.HasSuffix(baseURL, "/api/v4") && !strings.Contains(baseURL, "/api/v4/") {
			baseURL = strings.TrimSuffix(baseURL, "/") + "/api/v4"
		}
	}

	return baseURL
}

func fetchData(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build GET %q: %w", url, err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %q: server returned %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return body, nil
}

/* -------------------------------------------------------------------------- */
/*  CompareVersions and helpers                                               */
/* -------------------------------------------------------------------------- */

// CompareVersions compares two version strings.
// Returns -1 if own is older than latest, 0 if equal, and 1 if newer.
// Supports version suffixes like "v1.0.0-commithash".
func CompareVersions(v VersionStrings) int {
	v.Own = TrimVersionPrefix(v.Own)
	v.Latest = TrimVersionPrefix(v.Latest)

	ownNumbers := strings.Split(v.Own, ".")
	latestNumbers := strings.Split(v.Latest, ".")

	for i := 0; i < len(ownNumbers) && i < len(latestNumbers); i++ {
		ownNum, ownSuffix := extractNumberAndSuffix(ownNumbers[i])
		latestNum, latestSuffix := extractNumberAndSuffix(latestNumbers[i])

		// Compare numeric parts first
		if ownNum > latestNum {
			return 1
		}
		if ownNum < latestNum {
			return -1
		}

		// If numeric parts are equal, compare suffixes lexicographically
		if ownSuffix != latestSuffix {
			if ownSuffix > latestSuffix {
				return 1
			}
			if ownSuffix < latestSuffix {
				return -1
			}
		}
	}

	// If all compared components are equal, longer version is considered newer
	if len(ownNumbers) > len(latestNumbers) {
		return 1
	}
	if len(ownNumbers) < len(latestNumbers) {
		return -1
	}
	return 0
}

// extractNumberAndSuffix extracts the numeric prefix and any suffix from a version component.
// For example, "123-abc" returns (123, "-abc"), and "456" returns (456, "").
func extractNumberAndSuffix(component string) (int, string) {
	// Find the first non-digit character
	for i, char := range component {
		if char < '0' || char > '9' {
			// Parse the numeric part
			num, err := strconv.Atoi(component[:i])
			if err != nil {
				// If we can't parse the number, treat it as 0
				return 0, component
			}
			return num, component[i:]
		}
	}
	// No suffix, parse the whole thing as a number
	num, err := strconv.Atoi(component)
	if err != nil {
		return 0, component
	}
	return num, ""
}

// CompareVersionsHelper wraps CompareVersions and returns descriptive strings.
func CompareVersionsHelper(v VersionStrings) string {
	switch CompareVersions(v) {
	case -1:
		if v.VersionStrings.Older == "" {
			v.VersionStrings.Older = "There is a newer release available"
		}
		if v.VersionStrings.UpgradeURL != "" {
			v.VersionStrings.Older = "There is a newer release available at " + v.VersionStrings.UpgradeURL
		}
		if v.VersionOptions.DieIfOlder {
			fmt.Println(v.VersionStrings.Older)
			os.Exit(125)
		}
		return v.VersionStrings.Older
	case 0:
		if v.VersionOptions.ShowMessageOnCurrent {
			if v.VersionStrings.Equal == "" {
				v.VersionStrings.Equal = "You are up to date"
			}
			return v.VersionStrings.Equal
		}
		return ""
	case 1:
		if v.VersionStrings.Newer == "" {
			v.VersionStrings.Newer = "You are on an unreleased version"
		}
		if v.VersionOptions.DieIfNewer {
			fmt.Println(v.VersionStrings.Newer)
			os.Exit(125)
		}
		return v.VersionStrings.Newer
	}
	return ""
}
