// Package gitearelease contains the functions to fetch the releases from a Gitea repository
package gitearelease

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetReleases will return the all the releases or just the latest release of a repository
func GetReleases(repoURL, owner, repo string, latest bool) ([]Release, error) {
	releaseType := "releases"
	if latest {
		releaseType = "releases/latest"
	}
	apiURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/%s", repoURL, owner, repo, releaseType)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest releases: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	var releases []Release

	if latest {
		var release Release
		err = json.Unmarshal(body, &release)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %s", err)
		}
		release.Body = strings.ReplaceAll(release.Body, "\n", " ")

		//Detect Type of Asset
		//Return Binary and executable type

		releases = append(releases, release)
	} else {
		err = json.Unmarshal(body, &releases)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %s", err)
		}
		for i := range releases {
			releases[i].Body = strings.ReplaceAll(releases[i].Body, "\n", " ")
		}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch latest releases: %s", string(body))
	}

	return releases, nil
}
