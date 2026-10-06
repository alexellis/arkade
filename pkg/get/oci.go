// Copyright (c) arkade author(s) 2026. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package get

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

	wantedName := tool.Name
	if osName := strings.ToLower(operatingSystem); strings.Contains(osName, "ming") || osName == "windows" {
		wantedName = wantedName + ".exe"
	}

	// Stream the exported image and extract only the wanted binary,
	// so unrelated entries (including symlinks) are skipped rather
	// than aborting the install.
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

	binaryPath, err := extractFileFromTar(tarFile, wantedName, tempDir)
	tarFile.Close()
	if err != nil {
		return "", "", fmt.Errorf("extracting %q from image %s for %s: %w",
			wantedName, imageName, platform.String(), err)
	}

	finalName := tool.Name
	if osName := strings.ToLower(operatingSystem); strings.Contains(osName, "ming") || osName == "windows" {
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

// extractFileFromTar scans a tar stream and writes the first regular
// file whose basename matches name to dir. Directory and symlink
// entries are skipped, so unrelated image contents never abort the
// install.
func extractFileFromTar(r io.Reader, name, dir string) (string, error) {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return "", fmt.Errorf("file %q not found in image", name)
		}
		if err != nil {
			return "", err
		}
		if hdr.Typeflag != tar.TypeReg || filepath.Base(hdr.Name) != name {
			continue
		}

		outPath := filepath.Join(dir, name)
		out, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return "", err
		}
		if err := out.Close(); err != nil {
			return "", err
		}
		return outPath, nil
	}
}

// ociPlatform converts arkade's OS/arch values into a container
// platform spec, failing fast on values outside the supported set
// rather than coercing them to a platform that installs the wrong
// binary.
func ociPlatform(arch, operatingSystem string) (*v1.Platform, error) {
	osName := strings.ToLower(operatingSystem)
	switch {
	case strings.Contains(osName, "ming") || osName == "windows":
		osName = "windows"
	case strings.HasPrefix(osName, "darwin"):
		osName = "darwin"
	case strings.HasPrefix(osName, "linux"):
		osName = "linux"
	default:
		return nil, fmt.Errorf("OS %q is not supported for OCI-based tools, use linux, darwin or ming (windows)", operatingSystem)
	}

	archName := strings.ToLower(arch)
	switch {
	case archName == "x86_64" || archName == "amd64":
		archName = "amd64"
	case archName == "aarch64" || archName == "arm64":
		archName = "arm64"
	case archName == "armv6l" || archName == "armv7l":
		archName = "arm"
	default:
		return nil, fmt.Errorf("architecture %q is not supported for OCI-based tools, use amd64, arm64 or armv7l", arch)
	}

	return &v1.Platform{OS: osName, Architecture: archName}, nil
}
