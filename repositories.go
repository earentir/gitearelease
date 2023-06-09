package gitearelease

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetRepositories returns all repositories of a user from a gitea instance can be filtered by release
// if withrelease is true only repositories with releases will be returned
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

	if withrelease == false {
		return allRepos, nil
	}

	reposWithRelease := []Repository{}
	for _, repo := range allRepos {
		if repo.ReleaseCounter > 0 && withrelease {
			reposWithRelease = append(reposWithRelease, repo)
		}
	}
	return reposWithRelease, nil
}
