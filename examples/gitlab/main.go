// examples/gitlab_example.go
// Demonstrates usage of the gitearelease package with GitLab.
// Note: GitLab has some limitations compared to Gitea/GitHub (see PROVIDER_DIFFERENCES.md).
package main

import (
	"fmt"
	"time"

	"github.com/earentir/gitearelease"
)

func main() {
	// Configure HTTP timeout (optional)
	gitearelease.SetHTTPTimeout(30 * time.Second)

	// Fetch latest release from GitLab
	// Note: GitLab uses project path encoding (owner%2Frepo)
	relCfg := gitearelease.ReleaseToFetch{
		BaseURL: "https://gitlab.com/api/v4", // GitLab API base URL with /api/v4
		User:    "gitlab-org",                // Owner/username
		Repo:    "gitlab",                    // Repository name
		Latest:  true,
		// Provider is auto-detected from BaseURL, but you can explicitly set it:
		// Provider: "gitlab",
	}

	releases, err := gitearelease.GetReleases(relCfg)
	if err != nil {
		fmt.Printf("Error fetching releases: %v\n", err)
		return
	}

	if len(releases) == 0 {
		fmt.Println("No releases found.")
		return
	}

	latest := releases[0]
	fmt.Printf("Latest GitLab Release:\n")
	fmt.Printf("  Tag: %s\n", latest.TagName)
	fmt.Printf("  Name: %s\n", latest.Name)
	fmt.Printf("  Published: %s\n", latest.PublishedAt)
	// Note: GitLab doesn't support Draft/Prerelease flags the same way
	fmt.Printf("  Draft: %v, Prerelease: %v (always false for GitLab)\n", latest.Draft, latest.Prerelease)
	fmt.Printf("  Assets: %d\n", len(latest.Assets))

	for _, asset := range latest.Assets {
		// Note: GitLab assets don't have Size, DownloadCount, or CreatedAt
		fmt.Printf("    - %s (URL: %s)\n", asset.Name, asset.BrowserDownloadURL)
		fmt.Printf("      Note: Size and download count not available for GitLab\n")
	}

	// Fetch repositories
	repoCfg := gitearelease.RepositoriesToFetch{
		BaseURL:      "https://gitlab.com/api/v4",
		User:         "gitlab-org",
		WithReleases: true,
	}

	repos, err := gitearelease.GetRepositories(repoCfg)
	if err != nil {
		fmt.Printf("Error fetching repositories: %v\n", err)
		return
	}

	fmt.Printf("\nRepositories with releases: %d\n", len(repos))
	for _, repo := range repos {
		// Note: GitLab's ReleaseCounter is always 0 (not available in API)
		fmt.Printf("  - %s (Releases: %d - always 0 for GitLab)\n", repo.Name, repo.ReleaseCounter)
		fmt.Printf("    Note: HasIssues, HasWiki, HasProjects, HasPackages not available for GitLab\n")
	}
}
