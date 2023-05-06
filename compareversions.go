package gitearelease

import (
	"fmt"
	"strconv"
	"strings"
)

func CompareVersions(v1, v2 string) int {
	// Remove the version prefixs from the version strings
	v1 = trimVersionPrefix(v1)
	v2 = trimVersionPrefix(v2)

	// Split the version strings into individual version numbers
	v1Numbers := strings.Split(v1, ".")
	v2Numbers := strings.Split(v2, ".")

	// Convert the version numbers to integers and compare them
	for i := 0; i < len(v1Numbers) && i < len(v2Numbers); i++ {
		v1Num, err := strconv.Atoi(v1Numbers[i])
		if err != nil {
			fmt.Println("Invalid version number:", v1Numbers[i])
			return 0
		}
		v2Num, err := strconv.Atoi(v2Numbers[i])
		if err != nil {
			fmt.Println("Invalid version number:", v2Numbers[i])
			return 0
		}

		if v1Num > v2Num {
			return 1
		} else if v1Num < v2Num {
			return -1
		}
	}

	// If we got here, the version numbers are the same up to this point
	// Check if one of the versions has more version numbers than the other
	if len(v1Numbers) > len(v2Numbers) {
		return 1
	} else if len(v1Numbers) < len(v2Numbers) {
		return -1
	}

	// The version numbers are the same
	return 0
}

func trimVersionPrefix(version string) string {
	version = strings.ToLower(version)

	var verstrings []string = []string{"v", "version", "ver", "release", "rel", "r", "v."}
	for _, verstring := range verstrings {
		version = strings.TrimPrefix(version, verstring)
	}

	return version
}
