package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallTCRedirectTap() *cobra.Command {
	command := &cobra.Command{
		Use:   "tc-redirect-tap",
		Short: "Install tc-redirect-tap",
		Long:  `Install tc-redirect-tap cni plugin for use with faasd, CNI, Kubernetes, etc.`,
		Example: `  arkade system install tc-redirect-tap
  arkade system install tc-redirect-tap --version 2022-04-01-1337`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", githubLatest, "The version for tc-redirect-tap to install")
	command.Flags().StringP("path", "p", "/opt/cni/bin/", "Installation path, where the binary will be added")
	command.Flags().Bool("progress", true, "Show download progress")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		owner := "alexellis"
		repo := "tc-tap-redirect-builder"

		fmt.Printf("Installing tc-redirect-tap to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		if arch != "x86_64" && arch != "aarch64" {
			return fmt.Errorf("this app only supports x86_64 and aarch64 and not %s", arch)
		}

		dlArch := arch
		if arch == "aarch64" {
			dlArch = "arm64"
		}

		if version == githubLatest {
			v, err := get.FindGitHubRelease(owner, repo)
			if err != nil {
				return err
			}

			version = v
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		filename := fmt.Sprintf("tc-redirect-tap-%s", dlArch)
		if dlArch == "x86_64" {
			filename = "tc-redirect-tap"
		}
		dlURL := fmt.Sprintf(githubDownloadTemplate, owner, repo, version, filename)

		fmt.Printf("Downloading from: %s\n", dlURL)
		outPath, err := get.DownloadFileP(dlURL, progress)
		if err != nil {
			return err
		}
		defer os.Remove(outPath)

		fmt.Printf("Downloaded to: %s\n", outPath)

		dst := fmt.Sprintf("%s/%s", installPath, filename)
		fmt.Printf("Copying binary to: %s\n", installPath)
		if _, err := get.CopyFileP(outPath, dst, readWriteExecuteEveryone); err != nil {
			return err
		}

		return nil
	}

	return command
}
