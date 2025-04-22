// Package gitearelease provides functions to interact with Gitea's API and fetch the releases
package gitearelease

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
func GetReleases(r ReleaseToFetch) ([]Release, error) {
	releaseType := "releases"
	if r.Latest {
		releaseType = "releases/latest"
	}
	apiURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/%s",
		r.BaseURL, r.User, r.Repo, releaseType)

	apiData, err := fetchData(apiURL)
	if err != nil {
		return nil, err
	}

	var releases []Release
	if r.Latest {
		var rel Release
		if err := json.Unmarshal(apiData, &rel); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
		rel.Body = strings.ReplaceAll(rel.Body, "\n", " ")
		releases = append(releases, rel)
		return releases, nil
	}

	if err := json.Unmarshal(apiData, &releases); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	for i := range releases {
		releases[i].Body = strings.ReplaceAll(releases[i].Body, "\n", " ")
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
func GetRepositories(r RepositoriesToFetch) ([]Repository, error) {
	apiURL := fmt.Sprintf("%s/api/v1/users/%s/repos", r.BaseURL, r.User)
	// fmt.Println("APIURL", apiURL)

	apiData, err := fetchData(apiURL)
	if err != nil {
		return nil, err
	}

	//Print the raw API data for debugging
	// fmt.Println("Raw API Data:", string(apiData))

	var allRepos []Repository
	if err := json.Unmarshal(apiData, &allRepos); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	// Print the parsed repositories for debugging
	// fmt.Println("Parsed Repositories:", allRepos)

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

	// Print the filtered repositories for debugging
	// fmt.Println("Filtered Repositories with Releases:", repos)

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
func CompareVersions(v VersionStrings) int {
	v.Own = TrimVersionPrefix(v.Own)
	v.Latest = TrimVersionPrefix(v.Latest)

	ownNumbers := strings.Split(v.Own, ".")
	latestNumbers := strings.Split(v.Latest, ".")

	for i := 0; i < len(ownNumbers) && i < len(latestNumbers); i++ {
		ownNum, err := strconv.Atoi(ownNumbers[i])
		if err != nil {
			fmt.Println("Invalid version:", ownNumbers[i])
			return 0
		}
		latestNum, err := strconv.Atoi(latestNumbers[i])
		if err != nil {
			fmt.Println("Invalid version:", latestNumbers[i])
			return 0
		}
		if ownNum > latestNum {
			return 1
		}
		if ownNum < latestNum {
			return -1
		}
	}

	if len(ownNumbers) > len(latestNumbers) {
		return 1
	}
	if len(ownNumbers) < len(latestNumbers) {
		return -1
	}
	return 0
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
