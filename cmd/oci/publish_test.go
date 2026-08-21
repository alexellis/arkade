// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"strings"
	"syscall"
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

func TestSplitImageRejectsDigest(t *testing.T) {
	if _, _, _, err := splitImage("ghcr.io/me/app@sha256:abc123"); err == nil {
		t.Fatal("want error for digest reference")
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

func TestTarDirSymlinkDereferenced(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/real.txt", "target-content"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("real.txt", dir+"/link.txt"); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err != nil {
		t.Fatal(err)
	}

	tr := tar.NewReader(bytes.NewReader(buf.Bytes()))
	found := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if hdr.Name == "link.txt" {
			found = true
			if hdr.Typeflag != tar.TypeReg {
				t.Fatalf("want link.txt dereferenced to a regular file, got type %d", hdr.Typeflag)
			}
			body, err := io.ReadAll(tr)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != "target-content" {
				t.Fatalf("want dereferenced link.txt to carry target content, got %q", body)
			}
		}
	}
	if !found {
		t.Fatal("want link.txt in tar")
	}
}

func TestTarDirBrokenSymlinkErrors(t *testing.T) {
	dir := t.TempDir()
	if err := os.Symlink("does-not-exist", dir+"/dangling"); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err == nil {
		t.Fatal("want error for broken symlink")
	}
}

func TestTarDirSymlinkOutsideRootErrors(t *testing.T) {
	outside := t.TempDir()
	if err := writeFile(outside+"/secret.txt", "sensitive"); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Symlink(outside, dir+"/leak"); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err == nil {
		t.Fatal("want error when symlink resolves outside the source directory")
	}
}

func TestTarDirSymlinkCycleErrors(t *testing.T) {
	dir := t.TempDir()
	if err := os.Symlink(".", dir+"/loop"); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err == nil {
		t.Fatal("want error when a directory symlink cycles back into the source tree")
	}
}

func TestTarDirSkipsSpecialFile(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/ok.txt", "x"); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(dir+"/pipe", 0600); err != nil {
		t.Skipf("mkfifo unavailable: %s", err)
	}
	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err != nil {
		t.Fatalf("tarDir should not error or hang on a FIFO: %s", err)
	}
	if bytesContains(buf.Bytes(), "pipe") {
		t.Fatal("want FIFO to be skipped, not written to the tar")
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
