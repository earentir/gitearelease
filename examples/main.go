// examples/main.go
// Demonstrates usage of the gitearelease package: listing repos, fetching releases, and comparing versions.
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
	repoCfg := gitearelease.RepositoriesToFetch{
		BaseURL:      "https://gitea.com",
		User:         "earentir",
		WithReleases: true,
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
		fmt.Printf("Releases available: %d\n", r.ReleaseCounter)
		fmt.Println("---------------------------------")

		// 2. Fetch latest release
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
		// versionstrings.VersionOptions.DieIfOlder = true
		vs.VersionOptions.ShowMessageOnCurrent = true

		msg := gitearelease.CompareVersionsHelper(vs)
		fmt.Println(msg)

		// 4. List all releases and inspect assets
		relCfg.Latest = false
		rels, _ = gitearelease.GetReleases(relCfg)
		fmt.Println("All Releases:")
		for _, rel := range rels {
			fmt.Printf("- %s: %s\n", rel.TagName, rel.Body)
			for _, asset := range rel.Assets {
				fmt.Printf("  * %s (%d bytes)\n", asset.Name, asset.Size)
				// Detect OS/Arch by downloading first N bytes
				data, err := identifybin.DownloadFirstNBytes(asset.BrowserDownloadURL, 256)
				if err != nil {
					fmt.Println("Error downloading bytes:", err)
					continue
				}
				typeInfo, err := identifybin.DetectOSAndArch(data)
				if err != nil {
					fmt.Println("Error detecting OS/Arch:", err)
					continue
				}
				asset.Type = fmt.Sprintf("%s %s %s", typeInfo.OperatingSystem, typeInfo.Arch, typeInfo.Endianess)
				fmt.Println("    Asset Type:", asset.Type)
				fmt.Println("    UUID:", asset.UUID)
				fmt.Println("    Download count:", asset.DownloadCount)
				fmt.Println("    Type:", asset.Type)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
