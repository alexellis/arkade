// Copyright (c) arkade author(s) 2026. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package get

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// downloadFromOCI pulls the tool's OCI image for the requested
// platform, extracts the binary and installs it into movePath (or
// the arkade binary dir).
func downloadFromOCI(tool *Tool, arch, operatingSystem, version, movePath string, quiet bool, _ ProgressCallback) (string, string, error) {

	tag := GetToolVersion(tool, version)
	if tag == "" {
		tag = "latest"
	}

	platform, err := ociPlatform(arch, operatingSystem)
	if err != nil {
		return "", "", err
	}

	opts := []crane.Option{
		crane.WithPlatform(platform),
		crane.WithAuth(authn.Anonymous),
	}

	imageName := fmt.Sprintf("%s:%s", tool.OCIImage, tag)

	if !quiet {
		fmt.Printf("Resolving %s\n", imageName)
	}

	img, err := crane.Pull(imageName, opts...)
	if err != nil {
		return "", "", err
	}

	tempDir, err := os.MkdirTemp("", "arkade-oci-")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tempDir)

	tarFile, err := os.CreateTemp(tempDir, "image*.tar")
	if err != nil {
		return "", "", err
	}

	if err := crane.Export(img, tarFile); err != nil {
		tarFile.Close()
		return "", "", fmt.Errorf("exporting %s: %w", imageName, err)
	}
	_, err = tarFile.Seek(0, 0)
	if err != nil {
		tarFile.Close()
		return "", "", err
	}

	extractDir := filepath.Join(tempDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		tarFile.Close()
		return "", "", err
	}

	// Symlink extraction is disabled: the image is extracted into a
	// temp dir and we only need the plain binary.
	if err := archive.UntarNested(tarFile, extractDir, false, true, false, false); err != nil {
		tarFile.Close()
		return "", "", fmt.Errorf("untarring image for %s: %w", tool.Name, err)
	}
	tarFile.Close()

	wantedName := tool.Name
	if strings.Contains(strings.ToLower(operatingSystem), "mingw") {
		wantedName = wantedName + ".exe"
	}

	binaryPath, err := findFile(extractDir, wantedName)
	if err != nil {
		return "", "", fmt.Errorf("no binary %q found in image %s for %s: %w",
			wantedName, imageName, platform.String(), err)
	}

	finalName := tool.Name
	if strings.Contains(strings.ToLower(operatingSystem), "mingw") {
		finalName = finalName + ".exe"
	}

	var localPath string
	if movePath == "" {
		if _, err := config.InitUserDir(); err != nil {
			return "", "", err
		}
		localPath = env.LocalBinary(finalName, "")
	} else {
		localPath = filepath.Join(movePath, finalName)
	}

	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return "", "", err
	}

	if _, err := atomicCopyFile(binaryPath, localPath, 0755); err != nil {
		return "", "", err
	}

	return localPath, finalName, nil
}

// ociPlatform converts arkade's OS/arch values into a container
// platform spec.
func ociPlatform(arch, operatingSystem string) (*v1.Platform, error) {
	osName := strings.ToLower(operatingSystem)
	switch {
	case strings.Contains(osName, "mingw"):
		osName = "windows"
	case strings.HasPrefix(osName, "darwin"):
		osName = "darwin"
	case osName == "linux" || strings.HasPrefix(osName, "linux"):
		osName = "linux"
	default:
		return nil, fmt.Errorf("OS %q is not supported for OCI-based tools, use linux, darwin or mingw (windows)", operatingSystem)
	}

	archName := strings.ToLower(arch)
	switch {
	case archName == "x86_64" || archName == "amd64":
		archName = "amd64"
	case archName == "aarch64" || archName == "arm64":
		archName = "arm64"
	default:
		return nil, fmt.Errorf("architecture %q is not supported for OCI-based tools, use amd64 or arm64", arch)
	}

	return &v1.Platform{OS: osName, Architecture: archName}, nil
}

// findFile walks dir looking for a file with the given name,
// at any depth.
func findFile(dir, name string) (string, error) {
	var found string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == name {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("file %q not found in archive", name)
	}
	return found, nil
}
