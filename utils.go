package gitearelease

import "strings"

// trimVersionPrefix removes common version prefixes from a version string
func trimVersionPrefix(version string) string {
	version = strings.ToLower(version)

	var verstrings = []string{"v", "version", "ver", "release", "rel", "r", "v."}
	for _, verstring := range verstrings {
		version = strings.TrimPrefix(version, verstring)
	}

	return version
}
