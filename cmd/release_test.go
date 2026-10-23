package cmd

import "testing"

func TestBumpVersionPreservesTagStyle(t *testing.T) {
	for _, tc := range []struct {
		name string
		tag  string
		bump bumpStrategy
		want string
	}{
		{name: "prefixed patch", tag: "v1.0.4", bump: bumpPatchStrategy, want: "v1.0.5"},
		{name: "prefixed minor", tag: "v1.0.4", bump: bumpMinorStrategy, want: "v1.1.0"},
		{name: "prefixed major", tag: "v1.0.4", bump: bumpMajorStrategy, want: "v2.0.0"},
		{name: "bare patch", tag: "0.2.4", bump: bumpPatchStrategy, want: "0.2.5"},
		{name: "bare minor", tag: "0.2.4", bump: bumpMinorStrategy, want: "0.3.0"},
		{name: "bare major", tag: "0.2.4", bump: bumpMajorStrategy, want: "1.0.0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bumpVersion(tc.tag, tc.bump)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("bumpVersion(%q, %s) = %q, want %q", tc.tag, tc.bump, got, tc.want)
			}
		})
	}
}
