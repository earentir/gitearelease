# gitearelease

A package to access details for repositories and releases in gitea, it includes various helpers including the comparison of the current and leatest released version

## Usage/Examples

### Check version of current and compare to latest release
```go
import (
	"fmt"

	"github.com/earentir/gitearelease"
)

var (
	appversion = "1.1.14"
)

func checkVersion() {
	// Setup the release to fetch
	var releasetofetch gitearelease.ReleaseToFetch
	releasetofetch.BaseURL = "https://gitea.earentir.dev"
	releasetofetch.User = "earentir"
	releasetofetch.Repo = "dns"

	// Latest release
	releasetofetch.Latest = true

	rels, err := gitearelease.GetReleases(releasetofetch)
	if err != nil {
		fmt.Println(err)
	}

	var versionstrings gitearelease.VersionStrings
	versionstrings.Own = appversion
	versionstrings.Current = rels[0].TagName
        // Example of custom comparison messages
	versionstrings.VersionStrings.Older = "Upgrade this ASAP"
	versionstrings.VersionStrings.Newer = "You are ahead of the game"
        versionstrings.VersionStrings.Equal = "You are up to date"

        // Optionally we could terminate the applicaiton if we are older than latest release
	// versionstrings.VersionOptions.DieIfOlder = true

	fmt.Println(gitearelease.CompareVersionsHelper(versionstrings))
}
```


## Func Reference

### GetRepositories
```go
func GetRepositories(repositoriestofetch RepositoriesToFetch) ([]Repository, error)
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `repositoriestofetch` | `struct` | **Required**  |

GetRepositories returns all repositories of a user from a gitea instance can be filtered by release if withrelease is true only repositories with releases will be returned

### GetReleases
```go
func GetReleases(releasetofetch ReleaseToFetch) ([]Release, error)
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `releasetofetch` | `struct` | **Required**  |

GetReleases will return the all the releases or just the latest release of a repository

### CompareVersions
```go
  func CompareVersions(versionstrings VersionStrings) int
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `versionstrings` | `struct` | **Required**  |

CompareVersions compares two version strings and returns  -1 if own is older than current, 0 if own is equal to current and 1 if own is newer than current

### CompareVersionsHelper
```go
  func CompareVersionsHelper(versionstrings VersionStrings) string
```

| Parameter | Type  | Description   |
| :-------- | :----- | :----------- |
| `versionstrings`  | `struct` | **Required** |

CompareVersionsHelper is a helper function for CompareVersions that returns a string instead of an integer

#### Helpers
```go
  func DownloadBinary(url, outputDir, filename string) (string, error)
```
Simple Binary downloader to fetch the selected release (Not Tested)

```go
  func TrimVersionPrefix(version string) string
```
Trivial Version String Cleaner


#### Exported structs
```go
type Release
type ReleaseToFetch
type RepositoriesToFetch
type Repository
type VersionStrings
```
## Dependancies & Documentation
[![Go Mod](https://img.shields.io/github/go-mod/go-version/earentir/gitearelease)]()

[![Go Reference](https://pkg.go.dev/badge/github.com/earentir/gitearelease.svg)](https://pkg.go.dev/github.com/earentir/gitearelease)

[![Dependancies](https://img.shields.io/librariesio/github/earentir/gitearelease)](https://libraries.io/github/earentir/gitearelease)

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8581/badge)](https://www.bestpractices.dev/projects/8581)

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/earentir/gitearelease/badge)](https://securityscorecards.dev/viewer/?uri=github.com/earentir/gitearelease)


## Contributing

Contributions are always welcome!
All contributions are required to follow the https://google.github.io/styleguide/go/

All code contributed must include its tests in (_test) and have a minimum of 80% coverage

## Vulnerability Reporting

Please report any security vulnerabilities to the project using issues or directly to the owner.

## Code of Conduct
 This project follows the go project code of conduct, please refer to https://go.dev/conduct for more details

## Roadmap

- Create a finalised version 1 of the package
- Add support for downloading binaries from releases
- Add support for github releases
- Add support for gitlab releases

## Authors

- [@earentir](https://www.github.com/earentir)


## License

I will always follow the Linux Kernel License as primary, if you require any other OPEN license please let me know and I will try to accomodate it.

[![License](https://img.shields.io/github/license/earentir/gitearelease)](https://opensource.org/license/gpl-2-0)
