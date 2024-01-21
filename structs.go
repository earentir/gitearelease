package gitearelease

import "time"

// VersionStrings represents the strings used to describe the version comparison
type VersionStrings struct {
	Own            string
	Current        string
	VersionStrings versionstringstruct
	VersionOptions versionoptionsstruct
}

type versionstringstruct struct {
	Older           string
	Upgrade         string
	Equal           string
	Newer           string
	OfferUpgradeURL string
}

type versionoptionsstruct struct {
	DieIfOlder bool
}

// ReleaseToFetch represents a release from a repository
type ReleaseToFetch struct {
	BaseURL string
	User    string
	Repo    string
	Latest  bool
}

// RepositoriesToFetch represents a release from a repository
type RepositoriesToFetch struct {
	BaseURL    string
	User       string
	WithReleas bool
}

// Release represents a release from a repository
type Release struct {
	ID          int    `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	URL         string `json:"url"`
	HTMLUrl     string `json:"html_url"`
	TarballURL  string `json:"tarball_url"`
	ZipballURL  string `json:"zipball_url"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Author      author
	Assets      []asset
}

type author struct {
	Login     string `json:"login"`
	LoginName string `json:"login_name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

type asset struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	DownloadCount      int    `json:"download_count"`
	CreatedAt          string `json:"created_at"`
	UUID               string `json:"uuid"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Type               string `json:"type"` //Detect the Type of the asset
}

// Repository represents a repository from a user or organization
type Repository struct {
	ID    int `json:"id"`
	Owner struct {
		ID                int       `json:"id"`
		Login             string    `json:"login"`
		LoginName         string    `json:"login_name"`
		FullName          string    `json:"full_name"`
		Email             string    `json:"email"`
		AvatarURL         string    `json:"avatar_url"`
		Language          string    `json:"language"`
		IsAdmin           bool      `json:"is_admin"`
		LastLogin         time.Time `json:"last_login"`
		Created           time.Time `json:"created"`
		Restricted        bool      `json:"restricted"`
		Active            bool      `json:"active"`
		ProhibitLogin     bool      `json:"prohibit_login"`
		Location          string    `json:"location"`
		Website           string    `json:"website"`
		Description       string    `json:"description"`
		Visibility        string    `json:"visibility"`
		FollowersCount    int       `json:"followers_count"`
		FollowingCount    int       `json:"following_count"`
		StarredReposCount int       `json:"starred_repos_count"`
		Username          string    `json:"username"`
	} `json:"owner"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Empty       bool   `json:"empty"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
	Template    bool   `json:"template"`
	// Parent          interface{} `json:"parent"`
	Mirror          bool      `json:"mirror"`
	Size            int       `json:"size"`
	Language        string    `json:"language"`
	LanguagesURL    string    `json:"languages_url"`
	HTMLURL         string    `json:"html_url"`
	Link            string    `json:"link"`
	SSHURL          string    `json:"ssh_url"`
	CloneURL        string    `json:"clone_url"`
	OriginalURL     string    `json:"original_url"`
	Website         string    `json:"website"`
	StarsCount      int       `json:"stars_count"`
	ForksCount      int       `json:"forks_count"`
	WatchersCount   int       `json:"watchers_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	OpenPrCounter   int       `json:"open_pr_counter"`
	ReleaseCounter  int       `json:"release_counter"`
	DefaultBranch   string    `json:"default_branch"`
	Archived        bool      `json:"archived"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ArchivedAt      time.Time `json:"archived_at"`
	Permissions     struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	HasIssues       bool `json:"has_issues"`
	InternalTracker struct {
		EnableTimeTracker                bool `json:"enable_time_tracker"`
		AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
		EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
	} `json:"internal_tracker"`
	HasWiki                       bool      `json:"has_wiki"`
	HasPullRequests               bool      `json:"has_pull_requests"`
	HasProjects                   bool      `json:"has_projects"`
	HasReleases                   bool      `json:"has_releases"`
	HasPackages                   bool      `json:"has_packages"`
	HasActions                    bool      `json:"has_actions"`
	IgnoreWhitespaceConflicts     bool      `json:"ignore_whitespace_conflicts"`
	AllowMergeCommits             bool      `json:"allow_merge_commits"`
	AllowRebase                   bool      `json:"allow_rebase"`
	AllowRebaseExplicit           bool      `json:"allow_rebase_explicit"`
	AllowSquashMerge              bool      `json:"allow_squash_merge"`
	AllowRebaseUpdate             bool      `json:"allow_rebase_update"`
	DefaultDeleteBranchAfterMerge bool      `json:"default_delete_branch_after_merge"`
	DefaultMergeStyle             string    `json:"default_merge_style"`
	DefaultAllowMaintainerEdit    bool      `json:"default_allow_maintainer_edit"`
	AvatarURL                     string    `json:"avatar_url"`
	Internal                      bool      `json:"internal"`
	MirrorInterval                string    `json:"mirror_interval"`
	MirrorUpdated                 time.Time `json:"mirror_updated"`
	// RepoTransfer                  interface{} `json:"repo_transfer"`
}
