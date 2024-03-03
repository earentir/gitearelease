package gitearelease

import (
	"encoding/json"
	"fmt"
)

// GetRepositories returns all repositories of a user from a gitea instance can be filtered by release
// if withrelease is true only repositories with releases will be returned
func GetRepositories(repositoriestofetch RepositoriesToFetch) ([]Repository, error) {

	apiURL := fmt.Sprintf("%s/api/v1/users/%s/repos", repositoriestofetch.BaseURL, repositoriestofetch.User)

	apiData, err := fetchData(apiURL)
	if err != nil {
		return nil, err
	}

	var allRepos []Repository
	err = json.Unmarshal(apiData, &allRepos)
	if err != nil {
		return nil, err
	}

	if !repositoriestofetch.WithReleas {
		return allRepos, nil
	}

	reposWithRelease := []Repository{}
	for _, repo := range allRepos {
		if repo.ReleaseCounter > 0 && repositoriestofetch.WithReleas {
			reposWithRelease = append(reposWithRelease, repo)
		}
	}
	return reposWithRelease, nil
}
