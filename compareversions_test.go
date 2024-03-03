package gitearelease

import (
	"testing"
)

func TestCompareVersionsHelper(t *testing.T) {
	versionstrings := VersionStrings{
		VersionStrings: versionstringstruct{
			Older: "There is a newer release available",
			Equal: "You are up to date",
			Newer: "You are on an unreleased version",
		},
		VersionOptions: versionoptionsstruct{
			DieIfOlder:           true,
			DieIfNewer:           true,
			ShowMessageOnCurrent: true,
		},
	}

	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.1"

	expected := "There is a newer release available"
	versionstrings.VersionStrings.Older = ""
	versionstrings.VersionStrings.UpgradeURL = ""
	versionstrings.VersionOptions.DieIfOlder = false
	result := CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.VersionOptions.DieIfOlder = false
	versionstrings.VersionStrings.UpgradeURL = "https://example.com/upgrade"
	expected = "There is a newer release available at https://example.com/upgrade"
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	expected = "There is a newer release available"
	versionstrings.VersionStrings.UpgradeURL = ""
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.0"
	versionstrings.VersionOptions.ShowMessageOnCurrent = false
	expected = ""
	versionstrings.VersionStrings.Equal = ""
	versionstrings.VersionOptions.DieIfOlder = false
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.0"
	versionstrings.VersionOptions.ShowMessageOnCurrent = true
	expected = "You are up to date"
	versionstrings.VersionStrings.Equal = ""
	versionstrings.VersionOptions.DieIfOlder = false
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.VersionOptions.ShowMessageOnCurrent = true
	versionstrings.VersionOptions.DieIfNewer = false
	expected = "You are on an unreleased version"
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.0"
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

	versionstrings.VersionOptions.ShowMessageOnCurrent = true
	versionstrings.VersionOptions.DieIfNewer = false
	versionstrings.VersionStrings.Newer = ""
	expected = "You are on an unreleased version"
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.0"
	result = CompareVersionsHelper(versionstrings)
	if result != expected {
		t.Errorf("Expected: %s, got: %s", expected, result)
	}

}
func TestCompareVersions(t *testing.T) {
	versionstrings := VersionStrings{
		VersionStrings: versionstringstruct{
			Older:      "Older version",
			Equal:      "Equal version",
			Newer:      "Newer version",
			UpgradeURL: "https://example.com/upgrade",
		},
		VersionOptions: versionoptionsstruct{
			DieIfOlder:           true,
			DieIfNewer:           true,
			ShowMessageOnCurrent: true,
		},
	}

	// Test case 1: Own version is older than the current version
	versionstrings.Own = "1.0.0"
	versionstrings.Latest = "1.0.1"
	expected := -1
	result := CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 2: Own version is newer than the current version
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.0"
	expected = 1
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 3: Own version is equal to the current version
	versionstrings.Own = "1.0.1"
	versionstrings.Latest = "1.0.1"
	expected = 0
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 4: Own version has more version numbers than the current version
	versionstrings.Own = "1.0.1.1"
	versionstrings.Latest = "1.0.1"
	expected = 1
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}

	// Test case 5: Own version has fewer version numbers than the current version
	versionstrings.Own = "1.0"
	versionstrings.Latest = "1.0.1"
	expected = -1
	result = CompareVersions(versionstrings)
	if result != expected {
		t.Errorf("Expected: %d, got: %d", expected, result)
	}
}
