// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/arkade/pkg/update"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

func MakeUpdate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "update",
		Short: "Replace the running binary with an updated version",
		Long: `The latest release version of arkade will be downloaded from GitHub.

If that text is not found by running "arkade version", then the release
will be downloaded, along with its .sha256 checksum file.

If the checksum matches the downloaded file then the running binary will be
replaced with the new binary.

This command can be run as often as you require, won't download the same 
version twice.`,
		Example:       `  arkade update`,
		Aliases:       []string{"u"},
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	command.Flags().Bool("verify", true, "Verify the checksum of the downloaded binary")
	command.Flags().Bool("force", false, "Force a download of the latest binary, even if up to date, the --verify flag still applies")

	command.RunE = func(cmd *cobra.Command, args []string) error {

		if runtime.GOOS == "windows" {
			return fmt.Errorf("update is not supported on Windows at this time")
		}

		verifyDigest, _ := cmd.Flags().GetBool("verify")
		forceDownload, _ := cmd.Flags().GetBool("force")

		u := update.NewUpdater().
			WithForce(forceDownload).
			WithVerify(verifyDigest).
			WithVerifier(update.DefaultVerifier{}).
			WithVersionCheck(update.DefaultVersionCheck{}).
			WithResolver(&urlResolver{})

		if err := u.Do(); err != nil {
			return err
		}

		fmt.Println("\n", aec.Bold.Apply(pkg.SupportMessageShort))

		return nil
	}
	return command
}

type urlResolver struct {
}

func (u *urlResolver) GetRelease() (string, error) {
	return get.FindGitHubRelease("alexellis", "arkade")
}

func (u *urlResolver) GetDownloadURL(release string) (string, error) {
	arch, operatingSystem := env.GetClientArch()
	arch = strings.ToLower(arch)
	operatingSystem = strings.ToLower(operatingSystem)

	if arch == "x86_64" {
		arch = "amd64"
	}

	name := "arkade"
	toolList := get.MakeTools()
	var tool *get.Tool
	for _, t := range toolList {
		if t.Name == name {
			tool = &t
			break
		}
	}

	downloadUrl, err := get.GetDownloadURL(tool, operatingSystem, arch, release, false)
	if err != nil {
		return "", err
	}

	return downloadUrl, nil
}
