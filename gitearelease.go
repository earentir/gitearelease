// Package gitearelease contains the functions to fetch the releases from a Gitea repository
package gitearelease

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetReleases will return the all the releases or just the latest release of a repository
func GetReleases(releasetofetch ReleaseToFetch) ([]Release, error) {
	releaseType := "releases"
	if releasetofetch.Latest {
		releaseType = "releases/latest"
	}
	apiURL := fmt.Sprintf("%s/api/v1/repos/%s/%s/%s", releasetofetch.BaseURL, releasetofetch.User, releasetofetch.Repo, releaseType)

	apiData, err := fetchData(apiURL)
	if err != nil {
		return nil, err
	}

	var releases []Release

	if releasetofetch.Latest {
		var release Release
		err = json.Unmarshal(apiData, &release)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %s", err)
		}
		release.Body = strings.ReplaceAll(release.Body, "\n", " ")

		//Detect Type of Asset
		//Return Binary and executable type

		releases = append(releases, release)
	} else {
		err = json.Unmarshal(apiData, &releases)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %s", err)
		}
		for i := range releases {
			releases[i].Body = strings.ReplaceAll(releases[i].Body, "\n", " ")
		}
	}

	return releases, nil
}
