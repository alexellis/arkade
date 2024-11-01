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

func TestGetCandidateTag(t *testing.T) {
	tests := []struct {
		name             string
		discoveredTags   []string
		currentTag       string
		expectedTagVal   string
		expectedIsSemVer bool
	}{
		{
			name:             "Valid semantic tags",
			discoveredTags:   []string{"v1.0.0", "v2.1.0", "v2.3.4", "v2.3.3"},
			currentTag:       "v2.3.3",
			expectedTagVal:   "v2.3.4",
			expectedIsSemVer: true,
		},
		{
			name:             "No valid semantic tags",
			discoveredTags:   []string{"invalid", "v.a.b", "xyz"},
			currentTag:       "v2.3.3",
			expectedTagVal:   "",
			expectedIsSemVer: false,
		},
		{
			name:             "Empty list",
			discoveredTags:   []string{},
			currentTag:       "v2.3.3",
			expectedTagVal:   "",
			expectedIsSemVer: false,
		},
		{
			name:             "Mixed valid and invalid tags",
			discoveredTags:   []string{"v1.2", "v1.0.0", "invalid", "v2.1.0-beta", "v1.2.4"},
			currentTag:       "v1.2.3",
			expectedTagVal:   "v1.2.4",
			expectedIsSemVer: true,
		},
		{
			name:             "Mixed valid and invalid tags",
			discoveredTags:   []string{"v1.2", "v1.0.0", "invalid", "v2.1.0-beta", "v1.3.3"},
			currentTag:       "v1.2.3",
			expectedTagVal:   "v1.3.3",
			expectedIsSemVer: true,
		},
		{
			name:             "similar tag values",
			discoveredTags:   []string{"17", "17.0", "17.0.0"},
			currentTag:       "16",
			expectedTagVal:   "17",
			expectedIsSemVer: true,
		},
		{
			name:             "similar tag values different arrival order",
			discoveredTags:   []string{"17.0", "17", "17.0.0"},
			currentTag:       "16",
			expectedTagVal:   "17",
			expectedIsSemVer: true,
		},
		{
			name:             "similar tag values different current format",
			discoveredTags:   []string{"17.0", "17", "17.0.0"},
			currentTag:       "16.0",
			expectedTagVal:   "17.0",
			expectedIsSemVer: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tagVal, isSemVer := getCandidateTag(tc.discoveredTags, tc.currentTag)
			if tagVal != tc.expectedTagVal || isSemVer != tc.expectedIsSemVer {
				t.Fatalf("\nwant: (%s, %v) \n got: (%s, %v)\n", tc.expectedTagVal, tc.expectedIsSemVer, tagVal, isSemVer)
			}
		})
	}
}

func TestAttributesMatch(t *testing.T) {
	tests := []struct {
		name     string
		c        tagAttributes
		n        tagAttributes
		expected bool
	}{
		{
			name: "All matching attributes",
			c: tagAttributes{
				hasSuffix: true,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3-beta",
			},
			n: tagAttributes{
				hasSuffix: true,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3-beta",
			},
			expected: true,
		},
		{
			name: "Different hasSuffix",
			c: tagAttributes{
				hasSuffix: true,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3-beta",
			},
			n: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3",
			},
			expected: false,
		},
		{
			name: "Different hasMajor",
			c: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3",
			},
			n: tagAttributes{
				hasSuffix: false,
				hasMajor:  false,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v0.2.3",
			},
			expected: false,
		},
		{
			name: "All attributes false",
			c: tagAttributes{
				hasSuffix: false,
				hasMajor:  false,
				hasMinor:  false,
				hasPatch:  false,
				original:  "",
			},
			n: tagAttributes{
				hasSuffix: false,
				hasMajor:  false,
				hasMinor:  false,
				hasPatch:  false,
				original:  "any",
			},
			expected: true,
		},
		{
			name: "Different hasMinor and hasPatch",
			c: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  false,
				original:  "v1.2",
			},
			n: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  false,
				hasPatch:  true,
				original:  "v1.0.1",
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.c.attributesMatch(tc.n)
			if result != tc.expected {
				t.Fatalf("\nwant: %t \n got: %t\n", tc.expected, result)
			}
		})
	}
}

func TestGetTagAttributes(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected tagAttributes
	}{
		{
			name: "Full semantic version with suffix",
			tag:  "v1.2.3-beta",
			expected: tagAttributes{
				hasSuffix: true,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3-beta",
			},
		},
		{
			name: "Full semantic version without suffix",
			tag:  "v1.2.3",
			expected: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  true,
				original:  "v1.2.3",
			},
		},
		{
			name: "Major and minor version without suffix",
			tag:  "v1.2",
			expected: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  true,
				hasPatch:  false,
				original:  "v1.2",
			},
		},
		{
			name: "Only major version without suffix",
			tag:  "v1",
			expected: tagAttributes{
				hasSuffix: false,
				hasMajor:  true,
				hasMinor:  false,
				hasPatch:  false,
				original:  "v1",
			},
		},
		{
			name: "Empty string",
			tag:  "",
			expected: tagAttributes{
				hasSuffix: false,
				hasMajor:  false,
				hasMinor:  false,
				hasPatch:  false,
				original:  "",
			},
		},
		{
			name: "Version with suffix only",
			tag:  "v1-beta",
			expected: tagAttributes{
				hasSuffix: true,
				hasMajor:  true,
				hasMinor:  false,
				hasPatch:  false,
				original:  "v1-beta",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getTagAttributes(tc.tag)
			if result != tc.expected {
				t.Fatalf("\nwant: %v \n got: %v\n", tc.expected, result)
			}
		})
	}
}
