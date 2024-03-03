package gitearelease

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CompareVersions compares two version strings and returns  -1 if own is older than current, 0 if own is equal to current and 1 if own is newer than current
func CompareVersions(versionstrings VersionStrings) int {
	// Remove the version prefixs from the version strings
	versionstrings.Own = TrimVersionPrefix(versionstrings.Own)
	versionstrings.Latest = TrimVersionPrefix(versionstrings.Latest)

	// Split the version strings into individual version numbers
	ownNumbers := strings.Split(versionstrings.Own, ".")
	currentNumbers := strings.Split(versionstrings.Latest, ".")

	// Convert the version numbers to integers and compare them
	for i := 0; i < len(ownNumbers) && i < len(currentNumbers); i++ {
		ownNum, err := strconv.Atoi(ownNumbers[i])
		if err != nil {
			fmt.Println("Invalid version number:", ownNumbers[i])
			return 0
		}
		currentNum, err := strconv.Atoi(currentNumbers[i])
		if err != nil {
			fmt.Println("Invalid version number:", currentNumbers[i])
			return 0
		}

		if ownNum > currentNum {
			return 1
		} else if ownNum < currentNum {
			return -1
		}
	}

	// If we got here, the version numbers are the same up to this point
	// Check if one of the versions has more version numbers than the other
	if len(ownNumbers) > len(currentNumbers) {
		return 1
	} else if len(ownNumbers) < len(currentNumbers) {
		return -1
	}

	// The version numbers are the same
	return 0
}

// CompareVersionsHelper is a helper function for CompareVersions that returns a string instead of an integer
func CompareVersionsHelper(versionstrings VersionStrings) string {
	switch CompareVersions(versionstrings) {
	case -1:
		if versionstrings.VersionStrings.Older == "" {
			versionstrings.VersionStrings.Older = "There is a newer release available"
		}

		if versionstrings.VersionStrings.UpgradeURL != "" {
			versionstrings.VersionStrings.Older = "There is a newer release available at " + versionstrings.VersionStrings.UpgradeURL
		}

		if versionstrings.VersionOptions.DieIfOlder {
			fmt.Println(versionstrings.VersionStrings.Older)
			os.Exit(125)
		}

		return versionstrings.VersionStrings.Older
	case 0:
		if versionstrings.VersionOptions.ShowMessageOnCurrent {
			if versionstrings.VersionStrings.Equal == "" {
				versionstrings.VersionStrings.Equal = "You are up to date"
			}
			return versionstrings.VersionStrings.Equal
		}
		return ""
	case 1:
		if versionstrings.VersionStrings.Newer == "" {
			versionstrings.VersionStrings.Newer = "You are on an unreleased version"
		}

		if versionstrings.VersionOptions.DieIfNewer {
			fmt.Println(versionstrings.VersionStrings.Newer)
			os.Exit(125)
		}

		return versionstrings.VersionStrings.Newer
	}

	return ""
}
