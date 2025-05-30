# gitearelease

`gitearelease` is a Go package for fetching repository and release metadata from a Gitea instance, with built‑in version comparison utilities and configurable HTTP timeouts.

---

## Installation

```bash
go get github.com/earentir/gitearelease
```

---

## Public API

### `SetHTTPTimeout(d time.Duration)`
Override the package‑wide HTTP client timeout. If `d <= 0`, resets to the default of 15 s.

```go
// set a 30 s timeout
gitearelease.SetHTTPTimeout(30 * time.Second)
```

---

### `GetRepositories(cfg RepositoriesToFetch) ([]Repository, error)`
Fetches the list of repositories for a user.

**Parameters**:
- `cfg.BaseURL string` – your Gitea server base (e.g. `https://gitea.example.com`).
- `cfg.User string` – the username or organization.
- `cfg.WithReleas bool` – legacy filter; if true, only repos with releases.
- `cfg.WithReleases bool` – preferred filter; if true, only repos with releases.

**Returns**:
- `[]Repository` – each repo has at least:
  - `ID int`
  - `ReleaseCounter int`
  - `Owner` metadata (login, timestamps, etc.)
  - other JSON fields are ignored on unmarshal.
- `error` on network, JSON, or HTTP status failures.

```go
repos, err := gitearelease.GetRepositories(
    gitearelease.RepositoriesToFetch{
        BaseURL:      "http://gitea.local",
        User:         "alice",
        WithReleases: true,
    },
)
```

---

### `GetReleases(cfg ReleaseToFetch) ([]Release, error)`
Retrieves all releases or just the latest for a given repository.

**Parameters**:
- `cfg.BaseURL string` – Gitea server base URL.
- `cfg.User string` – owner of the repo.
- `cfg.Repo string` – repository name.
- `cfg.Latest bool` – if true, only the latest release is fetched.

**Returns**:
- `[]Release` – each entry includes:
  - `ID, TagName, Name, Body, URL, HTMLURL, TarballURL, ZipballURL`
  - `Draft, Prerelease, CreatedAt, PublishedAt`
  - `Author` (login, email, full name)
  - `Assets` (ID, Name, Size, DownloadCount, CreatedAt, UUID, BrowserDownloadURL, Type)
- `error` on network, JSON, or HTTP status failures.

```go
rels, err := gitearelease.GetReleases(
    gitearelease.ReleaseToFetch{
        BaseURL: "http://gitea.local",
        User:    "alice",
        Repo:    "project",
        Latest:  true,
    },
)
```

---

### `TrimVersionPrefix(v string) string`
Removes common prefixes (`v`, `version`, `rel`, etc.) from a version string.

```go
clean := gitearelease.TrimVersionPrefix("v1.2.3") // "1.2.3"
```

---

### `CompareVersions(v VersionStrings) int`
Compares two version strings (after trimming prefixes).

**Parameters**:
- `v.Own string` – your current version.
- `v.Latest string` – the target version to compare against.
- prefixes are stripped before numeric comparison.

**Returns**:
- `-1` if `Own < Latest`
- `0`  if `Own == Latest`
- `1`  if `Own > Latest`

```go
res := gitearelease.CompareVersions(
    gitearelease.VersionStrings{Own: "1.0.0", Latest: "1.2.0"},
)
// res == -1
```

---

### `CompareVersionsHelper(v VersionStrings) string`
A convenience wrapper around `CompareVersions` that returns a custom message.

**Parameters** (in addition to `Own` and `Latest`):
- `v.VersionStrings.Older` – message when own is older.
- `v.VersionStrings.Equal` – message when equal (only if `ShowMessageOnCurrent` is true).
- `v.VersionStrings.Newer` – message when own is newer.
- `v.VersionStrings.UpgradeURL` – optional URL to include in `Older` message.
- `v.VersionOptions.DieIfOlder`/`DieIfNewer` – if set, prints the message and exits with code 125.
- `v.VersionOptions.ShowMessageOnCurrent` – whether to return the `Equal` message.

```go
vs := gitearelease.VersionStrings{
    Own:    "1.0.0",
    Latest: "1.2.0",
    VersionStrings: gitearelease.VersionStrings{
        Older:   "Please update!",
        Equal:   "You are up to date.",
        Newer:   "You are ahead.",
    },
    VersionOptions: gitearelease.VersionOptions{
        ShowMessageOnCurrent: true,
    },
}
msg := gitearelease.CompareVersionsHelper(vs)
```

---

## Example

See [examples/main.go](examples/main.go) for a full sample that:

1. Lists repos
1. Fetches the latest release
1. Compares versions
1. Iterates all releases and inspects assets

---

## Dependancies & Documentation
[![Go Mod](https://img.shields.io/github/go-mod/go-version/earentir/gitearelease)]()

[![Go Reference](https://pkg.go.dev/badge/github.com/earentir/gitearelease.svg)](https://pkg.go.dev/github.com/earentir/gitearelease)

[![Dependancies](https://img.shields.io/librariesio/github/earentir/gitearelease)](https://libraries.io/github/earentir/gitearelease)

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8581/badge)](https://www.bestpractices.dev/projects/8581)

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/earentir/gitearelease/badge)](https://securityscorecards.dev/viewer/?uri=github.com/earentir/gitearelease)

![Code Climate issues](https://img.shields.io/codeclimate/tech-debt/earentir/gitearelease)

---

## Contributing

Contributions are always welcome!
All contributions are required to follow the https://google.github.io/styleguide/go/

All code contributed must include its tests in (_test) and have a minimum of 80% coverage

---

## Vulnerability Reporting

Please report any security vulnerabilities to the project using issues or directly to the owner.

---

## Code of Conduct
 This project follows the go project code of conduct, please refer to https://go.dev/conduct for more details

 ---

## Roadmap

- Create a finalised version 1 of the package
- Add support for downloading binaries from releases
- Add support for github releases
- Add support for gitlab releases

---

## Authors

- [@earentir](https://www.github.com/earentir)

---

## License

I will always follow the Linux Kernel License as primary, if you require any other OPEN license please let me know and I will try to accomodate it.

[![License](https://img.shields.io/github/license/earentir/gitearelease)](https://opensource.org/license/gpl-2-0)
