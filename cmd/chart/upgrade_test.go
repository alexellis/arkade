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
