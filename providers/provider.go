// Package providers defines interfaces and implementations for different Git hosting platforms.
package providers

// Provider defines the interface that all Git hosting providers must implement.
type Provider interface {
	// GetReleasesURL constructs the API URL for fetching releases
	GetReleasesURL(baseURL, user, repo string, latest bool) string

	// GetRepositoriesURL constructs the API URL for fetching repositories
	GetRepositoriesURL(baseURL, user string) string

	// NormalizeRelease converts provider-specific JSON to the standard Release struct
	NormalizeRelease(data []byte, latest bool) ([]Release, error)

	// NormalizeRepositories converts provider-specific JSON to the standard Repository slice
	NormalizeRepositories(data []byte) ([]Repository, error)

	// DetectProvider checks if a given BaseURL matches this provider
	DetectProvider(baseURL string) bool
}

// ProviderType represents the type of Git hosting provider
type ProviderType string

const (
	ProviderGitea  ProviderType = "gitea"
	ProviderGitHub ProviderType = "github"
	ProviderGitLab ProviderType = "gitlab"
)

// GetProvider returns the appropriate provider implementation based on the type or auto-detection
func GetProvider(providerType ProviderType, baseURL string) Provider {
	// If provider type is explicitly set, use it
	if providerType != "" {
		switch providerType {
		case ProviderGitHub:
			return &GitHubProvider{}
		case ProviderGitLab:
			return &GitLabProvider{}
		case ProviderGitea:
			return &GiteaProvider{}
		default:
			// Fallback to auto-detection
		}
	}

	// Auto-detect provider from baseURL
	if NewGitHubProvider().DetectProvider(baseURL) {
		return NewGitHubProvider()
	}
	if NewGitLabProvider().DetectProvider(baseURL) {
		return NewGitLabProvider()
	}

	// Default to Gitea for backward compatibility
	return NewGiteaProvider()
}
