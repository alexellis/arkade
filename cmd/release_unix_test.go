//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseTagResolution(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	ghScript := `#!/bin/sh
if [ "$1" = "api" ] && [ "$2" = "repos/{owner}/{repo}/releases?per_page=1" ]; then
    if [ "$ARKADE_TEST_RELEASES_EXIT" = "1" ]; then
        exit 1
    fi
    printf '%s' "$ARKADE_TEST_RELEASES"
elif [ "$1" = "api" ] && [ "$2" = "repos/{owner}/{repo}/tags?per_page=1" ]; then
    if [ "$ARKADE_TEST_TAGS_EXIT" = "1" ]; then
        exit 1
    fi
    printf '%s' "$ARKADE_TEST_TAGS"
elif [ "$1" = "api" ]; then
    printf 'PUBLIC\n'
elif [ "$1" = "release" ] && [ "$2" = "create" ]; then
    printf '%s' "$3" > "$ARKADE_TEST_RELEASE_MARKER"
fi
`
	if err := os.WriteFile(gh, []byte(ghScript), 0755); err != nil {
		t.Fatal(err)
	}
	git := filepath.Join(dir, "git")
	if err := os.WriteFile(git, []byte("#!/bin/sh\nprintf 'Initial release\\n'\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	for _, tc := range []struct {
		name         string
		releases     string
		tags         string
		releasesExit string
		tagsExit     string
		wantTag      string
		wantErr      string
	}{
		{name: "bare release", releases: `[{"tag_name":"0.2.4"}]`, wantTag: "0.2.5"},
		{name: "prefixed release", releases: `[{"tag_name":"v1.0.4"}]`, wantTag: "v1.0.5"},
		{name: "new project", releases: "[]", tags: "[]", wantTag: "v0.0.1"},
		{name: "existing project", releases: "[]", tags: `[{"name":"0.11.133"}]`, wantErr: "repository tag \"0.11.133\" exists"},
		{name: "invalid release list", releases: "null", tags: "[]", wantErr: "null release list"},
		{name: "release without tag", releases: `[{}]`, wantErr: "release without a tag"},
		{name: "invalid tag list", releases: "[]", tags: "null", wantErr: "could not read repository tags"},
		{name: "release API failure", releasesExit: "1", wantErr: "gh exited with code 1"},
		{name: "tag API failure", releases: "[]", tagsExit: "1", wantErr: "could not check repository tags"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			marker := filepath.Join(t.TempDir(), "created")
			t.Setenv("ARKADE_TEST_RELEASES", tc.releases)
			t.Setenv("ARKADE_TEST_TAGS", tc.tags)
			t.Setenv("ARKADE_TEST_RELEASES_EXIT", tc.releasesExit)
			t.Setenv("ARKADE_TEST_TAGS_EXIT", tc.tagsExit)
			t.Setenv("ARKADE_TEST_RELEASE_MARKER", marker)

			err := MakeRelease().Execute()
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected %q error, got %v", tc.wantErr, err)
				}
				if _, err := os.Stat(marker); !os.IsNotExist(err) {
					t.Fatalf("gh release create was called: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			tag, err := os.ReadFile(marker)
			if err != nil {
				t.Fatalf("gh release create was not called: %v", err)
			}
			if string(tag) != tc.wantTag {
				t.Fatalf("created tag %q, want %q", tag, tc.wantTag)
			}
		})
	}
}
