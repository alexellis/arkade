package docker

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_resolveDockerfilePath_FileIsUnchanged(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-paths-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	p := filepath.Join(tmpDir, "Dockerfile.custom")
	if err := os.WriteFile(p, []byte("FROM alpine:3.19\n"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	got, err := resolveDockerfilePath(p)
	if err != nil {
		t.Fatalf("resolveDockerfilePath returned error: %v", err)
	}
	if got != p {
		t.Fatalf("want %q, got %q", p, got)
	}
}

func Test_resolveDockerfilePath_DirImpliesDockerfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-paths-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	got, err := resolveDockerfilePath(tmpDir)
	if err != nil {
		t.Fatalf("resolveDockerfilePath returned error: %v", err)
	}

	want := filepath.Join(tmpDir, "Dockerfile")
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func Test_resolveDockerfilePath_MissingPathIsReturnedAsIs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-paths-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	missing := filepath.Join(tmpDir, "does-not-exist")
	got, err := resolveDockerfilePath(missing)
	if err != nil {
		t.Fatalf("resolveDockerfilePath returned error: %v", err)
	}
	if got != missing {
		t.Fatalf("want %q, got %q", missing, got)
	}
}

