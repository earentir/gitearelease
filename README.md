# gitearelease
Access gitea releases over the API

## Example
You get the owner and repo from the url

https://gitea.repodomain.tld/api/v1/repos/earentir/somerepo/releases

The owner is earentir

The repo is somerepo


```go

repoURL := "https://gitea.repodomain.tld"
owner := "earentir"
repo := "somerepo"

releases, err := gitearelease.GetLatestReleases(repoURL, owner, repo, true) //the last value is a bool, it will instead get the latest release by adding /latest in the URL
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
