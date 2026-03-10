package docker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenCommand_StdoutMode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	dockerfileContent := `FROM golang:1.24 AS builder
FROM alpine:3.19
FROM ghcr.io/openfaas/of-watchdog:0.25.0`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("failed to write dockerfile: %v", err)
	}

	cmd := MakeGen()
	cmd.SetArgs([]string{"-f", dockerfilePath, "--stdout"})

	var stdout, stderr strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	result := stdout.String()

	// Check that images are present in output
	expectedImages := []string{"golang", "alpine", "ghcr.io/openfaas/of-watchdog"}
	for _, img := range expectedImages {
		if !strings.Contains(result, "- "+img) {
			t.Errorf("stdout missing image: %s, stdout=%q, stderr=%q", img, result, stderr.String())
		}
	}

	// Check that it has the images: header
	if !strings.HasPrefix(result, "images:\n") {
		t.Errorf("expected stdout to start with 'images:\n', got: %q", result)
	}
}

func TestGenCommand_WriteToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	dockerfileContent := `FROM golang:1.24 AS builder
FROM alpine:3.19`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("failed to write dockerfile: %v", err)
	}

	cmd := MakeGen()
	cmd.SetArgs([]string{"-f", dockerfilePath})

	var stdout, stderr strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Check output contains success message (goes to stderr)
	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "Generated") {
		t.Errorf("expected success message, stdout=%q, stderr=%q", stdout.String(), stderr.String())
	}

	// Check that arkade.yaml was created
	arkadeYamlPath := filepath.Join(tmpDir, "arkade.yaml")
	content, err := os.ReadFile(arkadeYamlPath)
	if err != nil {
		t.Fatalf("failed to read arkade.yaml: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "images:") {
		t.Errorf("arkade.yaml missing 'images:' header")
	}
	if !strings.Contains(contentStr, "- golang") {
		t.Errorf("arkade.yaml missing golang image")
	}
	if !strings.Contains(contentStr, "- alpine") {
		t.Errorf("arkade.yaml missing alpine image")
	}
}

func TestGenCommand_NoImagesFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	dockerfileContent := `# Just a comment
RUN echo "hello"`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("failed to write dockerfile: %v", err)
	}

	cmd := MakeGen()
	cmd.SetArgs([]string{"-f", dockerfilePath})

	var output strings.Builder
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for no images found, got nil")
	}

	if !strings.Contains(output.String(), "no images found") {
		t.Errorf("expected error about no images, got: %s", output.String())
	}
}

func TestGenCommand_VariableImagesSkipped(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	dockerfileContent := `ARG VERSION=1.0
FROM golang:${VERSION}
FROM alpine:3.19`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("failed to write dockerfile: %v", err)
	}

	cmd := MakeGen()
	cmd.SetArgs([]string{"-f", dockerfilePath, "--stdout"})

	var stdout, stderr strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	result := stdout.String()

	// golang:${VERSION} should be skipped, only alpine should be present
	if strings.Contains(result, "golang") {
		t.Errorf("expected golang to be skipped (has variable tag), got: %s", result)
	}
	if !strings.Contains(result, "- alpine") {
		t.Errorf("expected alpine to be present, got: %s, stderr=%s", result, stderr.String())
	}
}

func TestGenCommand_RegistryWithPort(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	dockerfileContent := `FROM myregistry.com:5000/myimage:1.2.3`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("failed to write dockerfile: %v", err)
	}

	cmd := MakeGen()
	cmd.SetArgs([]string{"-f", dockerfilePath, "--stdout"})

	var stdout, stderr strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	result := stdout.String()

	if !strings.Contains(result, "- myregistry.com:5000/myimage") {
		t.Errorf("expected registry:port/image format, got: %s, stderr=%s", result, stderr.String())
	}
}

func TestGenCommand_DeduplicatesImages(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arkade-docker-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	dockerfileContent := `FROM alpine:3.19 AS builder
FROM alpine:3.19 AS runtime
FROM golang:1.24`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("failed to write dockerfile: %v", err)
	}

	cmd := MakeGen()
	cmd.SetArgs([]string{"-f", dockerfilePath, "--stdout"})

	var stdout, stderr strings.Builder
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	result := stdout.String()

	// Count occurrences of alpine
	count := strings.Count(result, "- alpine")
	if count != 1 {
		t.Errorf("expected alpine to appear once (deduplicated), got %d times, stdout=%q, stderr=%q", count, result, stderr.String())
	}
}
