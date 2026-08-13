// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallActionsRunner() *cobra.Command {

	command := &cobra.Command{
		Use:   "actions-runner",
		Short: "Install GitHub Actions Runner",
		Long:  `Install GitHub Actions Runner for self-hosted CI.`,
		Example: `  # Install actions-runner to the default directory
  arkade system install actions-runner

  # Install to an alternate directory, from a specific version
  arkade system install actions-runner --version 2.290.1 --path /opt/

  # Install on macOS
  arkade system install actions-runner --os darwin --arch arm64

  # Download archive, do not install
  arkade system install actions-runner --archive-only

  # Print the resolved version to stdout and exit
  arkade system install actions-runner --print-version`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "$HOME/actions-runner/", "Installation path, where the Actions Runner files will be extracted")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. x86_64 or arm64")
	command.Flags().String("os", "", "Operating system i.e. linux or darwin, leave blank to detect the client OS")
	command.Flags().BoolP("archive-only", "a", false, "Only download the archive and do not unpack it, prints the archive path to stdout")
	command.Flags().Bool("print-version", false, "Print the resolved version to stdout and exit")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {

		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")

		archiveOnly, _ := cmd.Flags().GetBool("archive-only")
		out := io.Writer(os.Stdout)
		if archiveOnly {
			out = os.Stderr
		}

		if version == "" {
			v, err := get.FindGitHubRelease("actions", "runner")
			if err != nil {
				return err
			}
			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		if printVersion, _ := cmd.Flags().GetBool("print-version"); printVersion {
			fmt.Println(strings.TrimPrefix(version, "v"))
			return nil
		}

		fmt.Fprintf(out, "Installing Actions Runner to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Fprintf(out, "Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if cmd.Flags().Changed("os") {
			osVer, _ = cmd.Flags().GetString("os")
		}

		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlOS := strings.ToLower(osVer)
		switch dlOS {
		case "darwin":
			dlOS = "osx"
		case "linux":
		default:
			return fmt.Errorf("unsupported operating system: %q, use linux or darwin (macOS)", osVer)
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "x64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		} else if arch == "armv7" || arch == "armv7l" {
			dlArch = "arm"
		}

		fmt.Fprintf(out, "Installing version: %s for: %s / %s\n", version, dlOS, dlArch)

		filename := fmt.Sprintf("actions-runner-%s-%s-%s.tar.gz", dlOS, dlArch, strings.TrimPrefix(version, "v"))
		dlURL := fmt.Sprintf(githubDownloadTemplate, "actions", "runner", version, filename)
		fmt.Fprintf(out, "Downloading from: %s\n", dlURL)

		progress, _ := cmd.Flags().GetBool("progress")
		outPath, err := get.DownloadFileP(dlURL, progress)
		if err != nil {
			return err
		}
		defer os.Remove(outPath)

		fmt.Fprintf(out, "Downloaded to: %s\n", outPath)

		if archiveOnly {
			dest := path.Join(installPath, filename)
			if _, err := get.CopyFileP(outPath, dest, readWriteExecuteEveryone); err != nil {
				return err
			}
			fmt.Println(dest)
			return nil
		}

		f, err := os.OpenFile(outPath, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Printf("Unpacking Actions Runner to: %s\n", path.Join(installPath, "actions-runner"))

		if err := spinWhile("Unpacking Actions Runner", func() error {
			return archive.UntarNested(f, installPath, true, true, true, false)
		}); err != nil {
			return err
		}

		return nil
	}

	return command
}
