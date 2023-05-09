package gitearelease

type Release struct {
	ID          int    `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	URL         string `json:"url"`
	HTMLUrl     string `json:"html_url"`
	TarballUrl  string `json:"tarball_url"`
	ZipballUrl  string `json:"zipball_url"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Author      Author
	Assets      []Asset
}

type Author struct {
	Login     string `json:"login"`
	LoginName string `json:"login_name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

type Asset struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	DownloadCount      int    `json:"download_count"`
	CreatedAt          string `json:"created_at"`
	UUID               string `json:"uuid"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Type               string `json:"type"`
}

type Repository struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	FullName       string `json:"full_name"`
	Description    string `json:"description"`
	ReleaseCounter int    `json:"release_counter"`
	Created        string `json:"created_at"`
	Updated        string `json:"updated_at"`
}
