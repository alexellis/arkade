// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSplitImage(t *testing.T) {
	tests := []struct {
		input    string
		wantRepo string
		wantTag  string
		hasTag   bool
	}{
		{"ghcr.io/me/app:0.1.0", "ghcr.io/me/app", "0.1.0", true},
		{"ghcr.io/me/app", "ghcr.io/me/app", "", false},
		{"ttl.sh/me/app:1h", "ttl.sh/me/app", "1h", true},
		{"ghcr.io/me/app:latest", "ghcr.io/me/app", "latest", true},
	}
	for _, tc := range tests {
		repo, tag, hasTag, err := splitImage(tc.input)
		if err != nil {
			t.Fatalf("splitImage(%q) error: %s", tc.input, err)
		}
		if repo != tc.wantRepo || tag != tc.wantTag || hasTag != tc.hasTag {
			t.Errorf("splitImage(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tc.input, repo, tag, hasTag, tc.wantRepo, tc.wantTag, tc.hasTag)
		}
	}
}

func TestPublishBundlesInvalidPlatform(t *testing.T) {
	cases := []string{
		"amd64=foo.tgz",          // no os/arch separator
		"linux/=foo.tgz",         // missing arch
		"/amd64=foo.tgz",         // missing os
		"linux/amd64/v2=foo.tgz", // variant not supported
	}
	for _, c := range cases {
		if _, err := publishBundles([]string{c}); err == nil {
			t.Errorf("want error for invalid platform bundle %q", c)
		}
	}
}

func TestPublishBundlesMissingSeparator(t *testing.T) {
	_, err := publishBundles([]string{"foo.tgz"})
	if err == nil {
		t.Fatal("want error for bundle without '=' separator")
	}
}

func TestTarDir(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/hello.txt", "world"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("hello.txt", dir+"/link.txt"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()
	if len(data) == 0 {
		t.Fatal("want non-empty tar")
	}
	if !bytesContains(data, "hello.txt") {
		t.Fatal("want hello.txt in tar")
	}
	if !bytesContains(data, "link.txt") {
		t.Fatal("want link.txt in tar")
	}
}

func TestTarDirSymlinkTarget(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/real.txt", "x"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("real.txt", dir+"/link.txt"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()

	tr := tar.NewReader(bytes.NewReader(data))
	var linkName string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if hdr.Name == "link.txt" {
			linkName = hdr.Linkname
		}
	}
	if linkName != "real.txt" {
		t.Fatalf("want symlink target real.txt, got %q", linkName)
	}
}

func TestBuildImageFromDirRejectsFile(t *testing.T) {
	file := t.TempDir() + "/notadir"
	if err := writeFile(file, "x"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := buildImageFromDir(file); err == nil {
		t.Fatal("want error when source is a regular file, not a directory")
	}
}

func TestBuildImageFromDirPlatform(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/f", "x"); err != nil {
		t.Fatal(err)
	}
	img, cleanup, err := buildImageFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	cfg, err := img.ConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.OS == "" || cfg.Architecture == "" {
		t.Fatalf("want OS and Architecture set on dir-built image, got os=%q arch=%q", cfg.OS, cfg.Architecture)
	}
}

func TestBuildImageFromBundlePlatform(t *testing.T) {
	img, err := buildImageFromBundle("", "linux/arm64")
	if err == nil {
		t.Fatal("want error reading missing bundle file")
	}
	_ = img
}

func writeFile(path, contents string) error {
	return os.WriteFile(path, []byte(contents), 0600)
}

func bytesContains(data []byte, substr string) bool {
	return strings.Contains(string(data), substr)
}
