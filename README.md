# gitearelease

`gitearelease` is a Go package for fetching repository and release metadata from Git hosting platforms (Gitea, GitHub, GitLab), with built‑in version comparison utilities and configurable HTTP timeouts.

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
Removes common prefixes (`v`, `version`, `rel`, etc.) from a version string. Used internally by `CompareVersions` before parsing version numbers and suffixes.

```go
clean := gitearelease.TrimVersionPrefix("v1.2.3") // "1.2.3"
clean := gitearelease.TrimVersionPrefix("v1.0.0-commithash") // "1.0.0-commithash"
```

---

### `CompareVersions(v VersionStrings) int`
Compares two version strings (after trimming prefixes). Supports semantic versions, version suffixes, and date-based versions.

**Parameters**:
- `v.Own string` – your current version.
- `v.Latest string` – the target version to compare against.
- prefixes are stripped before comparison.
- supports semantic versions, suffixes, and dates.

**Returns**:
- `-1` if `Own < Latest`
- `0`  if `Own == Latest`
- `1`  if `Own > Latest`

**Supported Version Formats**:
- **Semantic versions**: `1.0.0`, `2.1.3`, `0.1.33-c350f37`
- **Date-based versions**:
  - `YYYY-MM-DD` (e.g., `2024-01-15`)
  - `YYYY.MM.DD` (e.g., `2024.01.15`)
  - `YYYYMMDD` (e.g., `20240115`)
  - `YY-MM-DD` (e.g., `24-01-15`)

**Examples**:
- `v1.0.0` vs `v1.0.1` → -1
- `v1.0.0-abc123` vs `v1.0.0-def456` → -1 ("abc123" < "def456")
- `2024-01-15` vs `2024-01-10` → 1 (chronological comparison)
- `2023.12.25` vs `2024.01.01` → -1 (chronological comparison)

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

## Examples

### Provider-Specific Examples

Each example demonstrates the provider's capabilities and documents any limitations:

- **[Gitea Example](examples/gitea/main.go)**: Full-featured example showing Gitea-specific advantages
  - Accurate `ReleaseCounter`
  - Complete asset information (Size, DownloadCount, UUID, CreatedAt)
  - All repository metadata fields
  - Native `/releases/latest` endpoint

- **[GitHub Example](examples/github/main.go)**: Example using GitHub API
  - Near-complete feature support
  - Note: `ReleaseCounter` is a placeholder (1 if releases exist)

- **[GitLab Example](examples/gitlab/main.go)**: Example using GitLab API with limitations noted
  - Basic release/repository functionality
  - Note: Missing `ReleaseCounter`, asset details, draft/prerelease flags

### Original Example

See [examples/main.go](examples/main.go) for the original full sample that:

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

## Multi-Provider Support

This package now supports **Gitea**, **GitHub**, and **GitLab**! Provider detection is automatic based on the `BaseURL`, but you can also explicitly specify the provider.

### Auto-Detection

The package automatically detects the provider from the `BaseURL`:

```go
// GitHub - automatically detected
releases, err := gitearelease.GetReleases(gitearelease.ReleaseToFetch{
    BaseURL: "https://api.github.com",
    User:    "golang",
    Repo:    "go",
    Latest:  true,
})

// GitLab - automatically detected
releases, err := gitearelease.GetReleases(gitearelease.ReleaseToFetch{
    BaseURL: "https://gitlab.com",
    User:    "gitlab-org",
    Repo:    "gitlab",
    Latest:  true,
})

// Gitea - automatically detected (or default)
releases, err := gitearelease.GetReleases(gitearelease.ReleaseToFetch{
    BaseURL: "https://gitea.com",
    User:    "earentir",
    Repo:    "gitearelease",
    Latest:  true,
})
```

### Explicit Provider Selection

You can explicitly specify the provider:

```go
// Explicitly specify GitHub
releases, err := gitearelease.GetReleases(gitearelease.ReleaseToFetch{
    BaseURL:  "https://api.github.com",
    User:     "golang",
    Repo:     "go",
    Latest:   true,
    Provider: "github", // "gitea", "github", or "gitlab"
})
```

### BaseURL Format

- **GitHub**: Use `https://api.github.com` or `https://github.com` (auto-converted)
- **GitLab**: Use `https://gitlab.com/api/v4` or `https://gitlab.com` (auto-converted)
- **Gitea**: Use your Gitea instance URL (e.g., `https://gitea.com`)

### Backward Compatibility

All existing code continues to work without changes! The package defaults to Gitea behavior when no provider is specified, ensuring full backward compatibility.

### Provider Differences

⚠️ **Important**: Not all features are available on all providers. See [PROVIDER_DIFFERENCES.md](PROVIDER_DIFFERENCES.md) for detailed information.

**Quick Summary**:
- **Gitea**: Full feature support (most complete)
- **GitHub**: Near-complete support (ReleaseCounter is placeholder)
- **GitLab**: Limited support (missing ReleaseCounter, asset details, draft/prerelease flags)

**Key Differences**:
- `ReleaseCounter`: Only accurate for Gitea (placeholder for GitHub, unavailable for GitLab)
- Asset information: Complete for Gitea/GitHub, limited for GitLab (missing size, download count)
- Draft/Prerelease: Fully supported for Gitea/GitHub, not supported for GitLab
- Repository metadata: Complete for Gitea/GitHub, limited for GitLab

See the examples directory for provider-specific usage examples.

## Roadmap

- Create a finalised version 1 of the package
- Add support for downloading binaries from releases

---

## Authors

- [@earentir](https://www.github.com/earentir)

---

## License

I will always follow the Linux Kernel License as primary, if you require any other OPEN license please let me know and I will try to accomodate it.

[![License](https://img.shields.io/github/license/earentir/gitearelease)](https://opensource.org/license/gpl-2-0)
