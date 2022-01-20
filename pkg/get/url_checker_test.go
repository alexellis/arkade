//go:build e2e
// +build e2e

package get

import (
	"net/http"
	"testing"
)

// Test_CheckTools runs end to end tests to verify the URLS for various tools using a HTTP head request.

func Test_CheckTools(t *testing.T) {
	tools := MakeTools()
	toolsToSkip := []string{
		"kumactl", // S3 bucket disallow HEAD requests
	}

	os := "linux"
	arch := "x86_64"

	for _, toolV := range tools {
		tool := toolV
		skip := false

		for _, cli := range toolsToSkip {
			if tool.Name == cli {
				skip = true
			}
		}

		if skip {
			continue
		}

		t.Run("Download of "+tool.Name, func(t *testing.T) {
			t.Parallel()
			quiet := true
			url, err := tool.GetURL(os, arch, tool.Version, quiet)
			if err != nil {
				t.Fatalf("Error getting url for %s: %s", tool.Name, err)
			}
			// t.Logf("Checking %s via %s", tool.Name, url)

			status, body, headers, err := tool.Head(url)
			if err != nil {
				t.Fatalf("Error with HTTP HEAD for %s, %s: %s", tool.Name, url, err)
			}

			if status != http.StatusOK {
				t.Fatalf("Error with HTTP HEAD for %s, %s: status code: %d, body: %s", tool.Name, url, status, body)
			}

			if headers.Get("Content-Length") == "" {
				t.Fatalf("Error with HTTP HEAD for %s, %s: content-length zero", tool.Name, url)
			}
		})
	}
}
