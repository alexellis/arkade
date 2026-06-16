package pkg

import (
	"fmt"
	"strconv"
	"time"
)

var (
	Version        string
	GitCommit      string
	BuildTimestamp string
)

func GetVersionInfo() string {
	if len(Version) == 0 {
		return "dev"
	}
	return Version
}

// BuildString returns the version string.
func BuildString() string {
	return GetVersionInfo()
}

// BuildDateString returns a human-readable build date, or empty if not set.
func BuildDateString() string {
	if len(BuildTimestamp) == 0 {
		return ""
	}
	ts64, err := strconv.ParseInt(BuildTimestamp, 10, 64)
	if err != nil {
		return ""
	}
	return time.Unix(ts64, 0).UTC().Format(time.RFC3339)
}

func UserAgent() string {
	return fmt.Sprintf("arkade/%s", Version)
}
