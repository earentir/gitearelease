# gitearelease

Access gitea releases over the API

## get releases
```go func gitearelease.GetReleases(repoURL string, owner string, repo string, latest bool) ([]gitearelease.Release, error) ```

You get the owner and repo from the url

https://gitea.repodomain.tld/api/v1/repos/earentir/somerepo/releases

The owner is earentir

The repo is somerepo

```go

repoURL := "https://gitea.repodomain.tld"
owner := "earentir"
repo := "somerepo"

releases, err := gitearelease.GetReleases(repoURL, owner, repo, true) //the last value is a bool, it will instead get the latest release by adding /latest in the URL
if err != nil {
	fmt.Println(err)
	return
}

for _, release := range releases {
	fmt.Println("Release tag:", release.TagName)
	fmt.Println("  Release message:", release.Body)
	fmt.Println("  Author:", release.Author.FullName)
	fmt.Println("  Email:", release.Author.Email)
	fmt.Println("  Release date:", release.CreatedAt)
	fmt.Println("  Assets:")
	for _, asset := range release.Assets {
		fmt.Println("    Name:", asset.Name)
		fmt.Println("    Size:", asset.Size)
		fmt.Println("    Download URL:", asset.BrowserDownloadURL)
		fmt.Println("    UUID:", asset.UUID)
		fmt.Println("    Download count:", asset.DownloadCount)
		fmt.Println("    Type:", asset.Type)
	}
	fmt.Println()
}
```

## get user repos
```go func gitearelease.GetRepositories(baseURL string, user string, withrelease bool) ([]gitearelease.Repository, error)```

```go

repoURL := "https://gitea.repodomain.tld"
owner := "earentir"

rp, err := gitearelease.GetRepositories(repoURL, owner, true) //the last value is a bool, it will only return repos that have releases
if err != nil {
	fmt.Println(err)
	return
}

for _, r := range rp {
	fmt.Println("Repository:", r.Name)
	fmt.Println("Name:", r.FullName)
	fmt.Println("  Description:", r.Description)
	fmt.Println("  Release counter:", r.ReleaseCounter)
	fmt.Println("  Created:", r.Created)
	fmt.Println("  Updated:", r.Updated)
	fmt.Println()
}
```
