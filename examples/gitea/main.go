// examples/gitea_example.go
// Demonstrates usage of the gitearelease package with Gitea.
// This shows the original Gitea functionality and what's different from other providers.
package main

import (
	"fmt"
	"time"

	"github.com/earentir/gitearelease"
	"github.com/earentir/identifybin"
)

var (
	appVersion = "0.0.1"
)

func main() {
	// Configure HTTP timeout (optional)
	gitearelease.SetHTTPTimeout(30 * time.Second)

	// 1. List repositories (filter for those with releases)
	// Gitea-specific: ReleaseCounter is accurate and available
	repoCfg := gitearelease.RepositoriesToFetch{
		BaseURL:      "https://gitea.com", // Your Gitea instance URL
		User:         "earentir",          // Owner/username
		WithReleases: true,
		// Provider is optional - defaults to Gitea for backward compatibility
		// Provider: "gitea", // Explicitly set if needed
	}
	repos, err := gitearelease.GetRepositories(repoCfg)
	if err != nil {
		fmt.Println("Error fetching repos:", err)
		return
	}

	if len(repos) == 0 {
		fmt.Println("No repositories found.")
		return
	}

	for _, r := range repos {
		fmt.Printf("Repository: %s (Full: %s)\n", r.Name, r.FullName)
		// Gitea advantage: Accurate ReleaseCounter (not placeholder like GitHub/GitLab)
		fmt.Printf("Releases available: %d\n", r.ReleaseCounter)
		fmt.Println("---------------------------------")

		// 2. Fetch latest release
		// Gitea-specific: Has native /releases/latest endpoint (efficient)
		relCfg := gitearelease.ReleaseToFetch{
			BaseURL: "https://gitea.com",
			User:    "earentir",
			Repo:    r.Name,
			Latest:  true,
		}

		rels, err := gitearelease.GetReleases(relCfg)
		if err != nil {
			fmt.Println("Error fetching latest release:", err)
			continue
		}

		latest := rels[0]
		fmt.Printf("Latest tag: %s (published %s)\n", latest.TagName, latest.CreatedAt)

		// 3. Compare app version vs. latest tag
		var vs gitearelease.VersionStrings
		vs.Own = appVersion
		vs.Latest = latest.TagName
		vs.VersionStrings.Older = "Upgrade this ASAP"
		vs.VersionStrings.Newer = "You are ahead of the game"
		vs.VersionStrings.Equal = "You are up to date"
		vs.VersionOptions.ShowMessageOnCurrent = true

		msg := gitearelease.CompareVersionsHelper(vs)
		fmt.Println(msg)

		// 4. List all releases and inspect assets
		// Gitea-specific: Full asset information available
		relCfg.Latest = false
		rels, _ = gitearelease.GetReleases(relCfg)
		fmt.Println("All Releases:")
		for _, rel := range rels {
			fmt.Printf("- %s: %s\n", rel.TagName, rel.Body)
			for _, asset := range rel.Assets {
				// Gitea advantage: Full asset details including Size, DownloadCount, UUID
				fmt.Printf("  * %s (%d bytes, %d downloads)\n", asset.Name, asset.Size, asset.DownloadCount)
				fmt.Printf("    UUID: %s\n", asset.UUID) // Gitea-specific field

				// Detect OS/Arch by downloading first N bytes
				data, err := identifybin.DownloadFirstNBytes(asset.BrowserDownloadURL, 256)
				if err != nil {
					fmt.Println("    Error downloading bytes:", err)
					continue
				}
				typeInfo, err := identifybin.DetectOSAndArch(data)
				if err != nil {
					fmt.Println("    Error detecting OS/Arch:", err)
					continue
				}
				asset.Type = fmt.Sprintf("%s %s %s", typeInfo.OperatingSystem, typeInfo.Arch, typeInfo.Endianess)
				fmt.Printf("    Asset Type: %s\n", asset.Type)
				fmt.Printf("    Created: %s\n", asset.CreatedAt) // Available in Gitea
			}
			fmt.Println()
		}
		fmt.Println()
	}

	fmt.Println("\n=== Gitea-Specific Advantages ===")
	fmt.Println("✓ Accurate ReleaseCounter (not placeholder)")
	fmt.Println("✓ Native /releases/latest endpoint (efficient)")
	fmt.Println("✓ Full asset information (Size, DownloadCount, UUID, CreatedAt)")
	fmt.Println("✓ Draft and Prerelease flags fully supported")
	fmt.Println("✓ Complete repository metadata (HasIssues, HasWiki, HasProjects, HasPackages)")
	fmt.Println("✓ Real numeric Release IDs")
	fmt.Println("\nSee PROVIDER_DIFFERENCES.md for full comparison with GitHub/GitLab")
}
