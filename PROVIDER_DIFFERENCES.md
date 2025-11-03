# Provider Differences and Limitations

This document outlines the differences and limitations between Gitea, GitHub, and GitLab providers.

## Feature Comparison Matrix

| Feature | Gitea | GitHub | GitLab | Notes |
|---------|-------|--------|--------|-------|
| **GetReleases** | ✅ Full | ✅ Full | ✅ Full | All providers support fetching releases |
| **GetRepositories** | ✅ Full | ✅ Full | ✅ Full | All providers support fetching repositories |
| **Latest Release Endpoint** | ✅ Native | ✅ Native | ⚠️ Workaround | GitLab fetches all and takes first |
| **ReleaseCounter** | ✅ Accurate | ⚠️ Placeholder | ❌ Not Available | See details below |
| **Release ID** | ✅ Real ID | ✅ Real ID | ⚠️ Generated | GitLab uses tag name hash |
| **Draft Releases** | ✅ Supported | ✅ Supported | ❌ Not Supported | GitLab handles differently |
| **Prerelease Flag** | ✅ Supported | ✅ Supported | ❌ Not Supported | GitLab handles differently |
| **Asset Size** | ✅ Available | ✅ Available | ❌ Not Available | GitLab API doesn't provide |
| **Asset Download Count** | ✅ Available | ✅ Available | ❌ Not Available | GitLab API doesn't provide |
| **Asset UUID** | ✅ Available | ❌ Not Available | ❌ Not Available | Gitea-specific |
| **Asset CreatedAt** | ✅ Available | ✅ Available | ❌ Not Available | GitLab API doesn't provide |
| **Tarball/Zipball URLs** | ✅ Direct | ✅ Direct | ⚠️ Conditional | Only if sources available |
| **Repository HasIssues** | ✅ Available | ✅ Available | ❌ Not Available | Not in GitLab API |
| **Repository HasWiki** | ✅ Available | ✅ Available | ❌ Not Available | Not in GitLab API |
| **Repository HasProjects** | ✅ Available | ✅ Available | ❌ Not Available | Not in GitLab API |
| **Repository HasPackages** | ✅ Available | ✅ Available | ❌ Not Available | Not in GitLab API |

## Detailed Differences

### ReleaseCounter (Repository Release Count)

**Gitea**: ✅ Fully supported
- The `ReleaseCounter` field is directly available in the repository API response
- Accurate count of releases for the repository

**GitHub**: ⚠️ Placeholder only
- GitHub's repository API doesn't include release count in the standard response
- Currently set to `1` if `HasReleases` is `true`, otherwise `0`
- To get accurate count, would require an additional API call to `/repos/{owner}/{repo}/releases`

**GitLab**: ❌ Not available
- GitLab's projects API doesn't provide release count
- Currently always set to `0`
- To get accurate count, would require an additional API call to `/projects/{id}/releases`

### Latest Release Endpoint

**Gitea & GitHub**: ✅ Native support
- Both have dedicated `/releases/latest` endpoints
- Efficient single API call

**GitLab**: ⚠️ Workaround
- GitLab doesn't have a `/latest` endpoint
- Implementation fetches all releases and takes the first one
- Less efficient but functionally equivalent

### Release ID

**Gitea & GitHub**: ✅ Real numeric IDs
- Both platforms provide numeric IDs for releases
- Stable and unique identifiers

**GitLab**: ⚠️ Generated hash
- GitLab doesn't provide numeric IDs for releases
- Implementation generates a hash from the tag name
- Not a real ID but provides a consistent identifier

### Draft and Prerelease Flags

**Gitea & GitHub**: ✅ Full support
- Both platforms support draft and prerelease flags
- Accurately reflected in the `Draft` and `Prerelease` fields

**GitLab**: ❌ Not supported
- GitLab handles drafts and prereleases differently
- These fields are always set to `false` for GitLab releases
- GitLab uses different mechanisms (protected tags, milestones) for similar concepts

### Asset Information

**Gitea**: ✅ Complete
- All asset fields available: `ID`, `Name`, `Size`, `DownloadCount`, `CreatedAt`, `UUID`, `BrowserDownloadURL`, `Type`

**GitHub**: ✅ Mostly complete
- Available: `ID`, `Name`, `Size`, `DownloadCount`, `CreatedAt`, `BrowserDownloadURL`, `Type`
- Missing: `UUID` (GitHub-specific limitation)

**GitLab**: ⚠️ Limited
- Available: `ID`, `Name`, `BrowserDownloadURL`, `Type`
- Missing: `Size`, `DownloadCount`, `CreatedAt`, `UUID`
- GitLab's API doesn't provide these fields in the releases endpoint

### Tarball/Zipball URLs

**Gitea & GitHub**: ✅ Direct URLs
- Both provide direct tarball and zipball URLs
- Always available

**GitLab**: ⚠️ Conditional
- GitLab uses a different structure with "sources"
- URLs are only populated if sources are available in the release
- May be empty if no sources are attached

### Repository Metadata Fields

**Gitea**: ✅ Complete
- All metadata fields available: `HasIssues`, `HasWiki`, `HasProjects`, `HasPackages`, `HasReleases`

**GitHub**: ✅ Complete
- All metadata fields available: `HasIssues`, `HasWiki`, `HasProjects`, `HasPackages`, `HasReleases`

**GitLab**: ⚠️ Limited
- Only `HasReleases` is available
- Missing: `HasIssues`, `HasWiki`, `HasProjects`, `HasPackages`
- These fields are not part of GitLab's projects API response

## Workarounds and Future Improvements

### Getting Accurate ReleaseCounter

If you need accurate release counts for GitHub or GitLab:

```go
// For GitHub, you could make an additional call:
// GET /repos/{owner}/{repo}/releases
// Then count the releases

// For GitLab, you could make an additional call:
// GET /projects/{id}/releases
// Then count the releases
```

However, this would require changing the API design to support optional additional data fetching, which is not currently implemented to maintain API simplicity.

### Getting Asset Information for GitLab

GitLab's releases API doesn't provide asset size or download counts. To get this information, you would need to:
1. Make HEAD requests to each asset URL to get size
2. Use GitLab's statistics API (if available)
3. Store this information separately

These are not currently implemented as they would require significant API changes.

## Recommendations

1. **For ReleaseCounter**: If you need accurate counts, consider fetching releases separately for GitHub/GitLab
2. **For Asset Details**: Be aware that GitLab asset information is limited
3. **For Draft/Prerelease**: Use provider-specific logic if you need this for GitLab
4. **For Latest Release**: Performance impact is minimal for GitLab (only fetches until first release)

## Backward Compatibility

All existing Gitea code continues to work unchanged. The differences only affect:
- New code using GitHub/GitLab
- Code that relies on specific fields that may be unavailable
