package gitearelease

import (
	"fmt"
	"strconv"
	"strings"
)

// CompareVersions compares two version strings and returns  -1 if own is older than current, 0 if own is equal to current and 1 if own is newer than current
func CompareVersions(own, current string) int {
	// Remove the version prefixs from the version strings
	own = TrimVersionPrefix(own)
	current = TrimVersionPrefix(current)

	// Split the version strings into individual version numbers
	ownNumbers := strings.Split(own, ".")
	currentNumbers := strings.Split(current, ".")

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
