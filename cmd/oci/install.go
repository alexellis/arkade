// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/spf13/cobra"
)

func MakeOciInstall() *cobra.Command {
	command := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Install the contents of an OCI image to a given path",
		Long: `Use this command to install binaries or packages distributed within an 
OCI image.`,
		Example: `  # Install vmmeter to /usr/local/bin
  arkade oci install ghcr.io/openfaasltd/vmmeter

  # Install a specific version of vmmeter to /tmp/
  arkade oci install ghcr.io/openfaasltd/vmmeter --path /tmp --version 0.1.0

  # Install vmmeter for arm64 as an architecture override, instead of using uname
  arkade oci install ghcr.io/openfaasltd/vmmeter --arch arm64
`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "latest", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "/usr/local/bin", "Installation path, where a buildkitd subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")
	command.Flags().BoolP("gzipped", "g", false, "Is this a gzipped tarball?")
	command.Flags().Bool("quiet", false, "Suppress progress output")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		gzipped, _ := cmd.Flags().GetBool("gzipped")
		quiet, _ := cmd.Flags().GetBool("quiet")

		if len(args) < 1 {
			return fmt.Errorf("please provide an image name")
		}

		imageName := args[0]

		if !strings.Contains(imageName, ":") {
			imageName = imageName + ":" + version
		}

		st := time.Now()

		fmt.Printf("Installing %s to %s\n", imageName, installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		clientArch, clientOS := env.GetClientArch()

		if cmd.Flags().Changed("arch") {
			clientArch, _ = cmd.Flags().GetString("arch")
		}

		tempFile, err := os.CreateTemp(os.TempDir(), "arkade-oci-")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tempFile.Name())

		f, err := os.Create(tempFile.Name())
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", tempFile.Name(), err)
		}
		defer f.Close()

		var img v1.Image

		downloadArch, downloadOS := getDownloadArch(clientArch, clientOS)

		img, err = crane.Pull(imageName, crane.WithPlatform(&v1.Platform{Architecture: downloadArch, OS: downloadOS}))
		if err != nil {
			return fmt.Errorf("pulling %s: %w", imageName, err)
		}

		if err := crane.Export(img, f); err != nil {
			return fmt.Errorf("exporting %s: %w", imageName, err)
		}

		tarFile, err := os.Open(tempFile.Name())
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", tempFile.Name(), err)
		}
		defer tarFile.Close()

		if err := archive.UntarNested(tarFile, installPath, gzipped, quiet); err != nil {
			return fmt.Errorf("failed to untar %s: %w", tempFile.Name(), err)
		}

		fmt.Printf("Took %s\n", time.Since(st).Round(time.Millisecond))

		return nil
	}

	return command
}

func getDownloadArch(clientArch, clientOS string) (arch string, os string) {
	downloadArch := strings.ToLower(clientArch)
	downloadOS := strings.ToLower(clientOS)

	if downloadArch == "x86_64" {
		downloadArch = "amd64"
	} else if downloadArch == "aarch64" {
		downloadArch = "arm64"
	}

	return downloadArch, downloadOS
}
