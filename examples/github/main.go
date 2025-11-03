// examples/github_example.go
// Demonstrates usage of the gitearelease package with GitHub.
package main

import (
	"fmt"
	"time"

	"github.com/earentir/gitearelease"
)

func main() {
	// Configure HTTP timeout (optional)
	gitearelease.SetHTTPTimeout(30 * time.Second)

	// Fetch latest release from GitHub
	relCfg := gitearelease.ReleaseToFetch{
		BaseURL: "https://api.github.com", // GitHub API base URL
		User:    "golang",                 // Owner/username
		Repo:    "go",                     // Repository name
		Latest:  true,
		// Provider is auto-detected from BaseURL, but you can explicitly set it:
		// Provider: "github",
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
	fmt.Printf("Latest GitHub Release:\n")
	fmt.Printf("  Tag: %s\n", latest.TagName)
	fmt.Printf("  Name: %s\n", latest.Name)
	fmt.Printf("  Published: %s\n", latest.PublishedAt)
	fmt.Printf("  Draft: %v, Prerelease: %v\n", latest.Draft, latest.Prerelease)
	fmt.Printf("  Assets: %d\n", len(latest.Assets))

	for _, asset := range latest.Assets {
		fmt.Printf("    - %s (%d bytes, %d downloads)\n", asset.Name, asset.Size, asset.DownloadCount)
	}

	// Fetch repositories
	repoCfg := gitearelease.RepositoriesToFetch{
		BaseURL:      "https://api.github.com",
		User:         "golang",
		WithReleases: true,
	}

	repos, err := gitearelease.GetRepositories(repoCfg)
	if err != nil {
		fmt.Printf("Error fetching repositories: %v\n", err)
		return
	}

	fmt.Printf("\nRepositories with releases: %d\n", len(repos))
	for _, repo := range repos {
		// Note: GitHub's ReleaseCounter is a placeholder (1 if HasReleases is true)
		fmt.Printf("  - %s (Releases: %d)\n", repo.Name, repo.ReleaseCounter)
	}
}
