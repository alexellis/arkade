package chart

import (
	"testing"
)

func Test_tagIsUpgradable(t *testing.T) {
	tests := []struct {
		title    string
		current  string
		latest   string
		expected bool
	}{
		{
			title:    "Upgradeable",
			current:  "1.0.0",
			latest:   "1.1.0",
			expected: true,
		},
		{
			title:    "Same version",
			current:  "1.0.0",
			latest:   "1.0.0",
			expected: false,
		},
		{
			title:    "latest is RC",
			current:  "1.0.0",
			latest:   "1.0.0-RC",
			expected: false,
		},
		{
			title:    "latest is rc",
			current:  "1.0.0",
			latest:   "1.0.0-rc",
			expected: false,
		},
		{
			title:    "current is rootless",
			current:  "1.0.0-rootless",
			latest:   "1.1.0",
			expected: false,
		},
		{
			title:    "latest is rootless",
			current:  "1.0.0",
			latest:   "1.1.0-rootless",
			expected: false,
		},
		{
			title:    "current is 'latest'",
			current:  "latest",
			latest:   "1.0.0",
			expected: false,
		},
		{
			title:    "both are rootless different version'",
			current:  "1.0.0-rootless",
			latest:   "1.0.1-rootless",
			expected: true,
		},
		{
			title:    "both are rootless same version'",
			current:  "1.0.0-rootless",
			latest:   "1.0.0-rootless",
			expected: false,
		},
		{
			title:    "both are rc same version'",
			current:  "1.0.0-rc",
			latest:   "1.0.0-rc",
			expected: false,
		},
		{
			title:    "both are rc different version'",
			current:  "1.0.0-rc",
			latest:   "1.0.1-rc",
			expected: true,
		},
		{
			title:    "both are rc with suffix & same version'",
			current:  "1.0.0-rc1",
			latest:   "1.0.0-rc2",
			expected: false,
		},
		{
			title:    "both are rc with suffix & different version'",
			current:  "1.0.0-rc1",
			latest:   "1.0.1-rc2",
			expected: false,
		},
	}

	for _, tc := range tests {

		t.Run(tc.title, func(t *testing.T) {

			upgradeableRes := tagIsUpgradeable(tc.current, tc.latest)

			if upgradeableRes != tc.expected {
				t.Fatalf("want: %t\n got: %t\n", tc.expected, upgradeableRes)
			}
		})
	}
}

func TestGetLatestTag(t *testing.T) {
	tests := []struct {
		name             string
		discoveredTags   []string
		expectedTagVal   string
		expectedIsSemVer bool
	}{
		{
			name:             "Valid semantic tags",
			discoveredTags:   []string{"v1.0.0", "v2.1.0", "v2.3.4", "v2.3.3"},
			expectedTagVal:   "v2.3.4",
			expectedIsSemVer: true,
		},
		{
			name:             "No valid semantic tags",
			discoveredTags:   []string{"invalid", "v.a.b", "xyz"},
			expectedTagVal:   "",
			expectedIsSemVer: false,
		},
		{
			name:             "Empty list",
			discoveredTags:   []string{},
			expectedTagVal:   "",
			expectedIsSemVer: false,
		},
		{
			name:             "Mixed valid and invalid tags",
			discoveredTags:   []string{"v1.0.0", "invalid", "v2.1.0-beta", "v1.2.3"},
			expectedTagVal:   "v2.1.0-beta",
			expectedIsSemVer: true,
		},
		{
			name:             "similar tag values",
			discoveredTags:   []string{"17", "17.0", "17.0.0"},
			expectedTagVal:   "17",
			expectedIsSemVer: true,
		},
		{
			name:             "similar tag values different arrival order",
			discoveredTags:   []string{"17.0", "17", "17.0.0"},
			expectedTagVal:   "17.0",
			expectedIsSemVer: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tagVal, isSemVer := getLatestTag(tc.discoveredTags)
			if tagVal != tc.expectedTagVal || isSemVer != tc.expectedIsSemVer {
				t.Fatalf("\nwant: (%s, %v) \n got: (%s, %v)\n", tc.expectedTagVal, tc.expectedIsSemVer, tagVal, isSemVer)
			}
		})
	}
}
