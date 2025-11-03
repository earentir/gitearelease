// Package providers defines interfaces and implementations for different Git hosting platforms.
package providers

// Release represents a normalized release structure used by providers
type Release struct {
	ID          int
	TagName     string
	Name        string
	Body        string
	URL         string
	HTMLUrl     string
	TarballURL  string
	ZipballURL  string
	Draft       bool
	Prerelease  bool
	CreatedAt   string
	PublishedAt string
	Author      Author
	Assets      []Asset
}

// Author represents the author of a release
type Author struct {
	Login     string
	LoginName string
	FullName  string
	Email     string
	Username  string
}

// Asset represents a release asset
type Asset struct {
	ID                 int
	Name               string
	Size               int64
	DownloadCount      int
	CreatedAt          string
	UUID               string
	BrowserDownloadURL string
	Type               string
}

// Repository represents a normalized repository structure used by providers
type Repository struct {
	ID              int
	Name            string
	FullName        string
	Description     string
	Private         bool
	Fork            bool
	Size            int
	Language        string
	HTMLURL         string
	CloneURL        string
	SSHURL          string
	StarsCount      int
	ForksCount      int
	WatchersCount   int
	OpenIssuesCount int
	ReleaseCounter  int
	DefaultBranch   string
	Archived        bool
	CreatedAt       string
	UpdatedAt       string
	Owner           Owner
	Permissions     Permissions
	HasIssues       bool
	HasWiki         bool
	HasProjects     bool
	HasReleases     bool
	HasPackages     bool
}

// Owner represents repository owner information
type Owner struct {
	ID        int
	Login     string
	Username  string
	FullName  string
	Email     string
	AvatarURL string
}

// Permissions represents repository permissions
type Permissions struct {
	Admin bool
	Push  bool
	Pull  bool
}
