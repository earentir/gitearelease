package gitearelease

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetRepositories(baseURL, user string, withrelease bool) ([]Repository, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+"/api/v1/users/"+user+"/repos", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repositories: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var allRepos []Repository
	err = json.Unmarshal(body, &allRepos)
	if err != nil {
		return nil, err
	}

	reposWithReleases := []Repository{}
	for _, repo := range allRepos {
		if repo.ReleaseCounter > 0 && withrelease {
			reposWithReleases = append(reposWithReleases, repo)
		}
	}

	return reposWithReleases, nil
}
