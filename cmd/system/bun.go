package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallBun() *cobra.Command {
	command := &cobra.Command{
		Use:   "bun",
		Short: "Install Bun",
		Long:  `Bun is an incredibly fast JavaScript runtime, bundler, transpiler and package manager – all in one.`,
		Example: `arkade system install bun
  arkade system install bun --version v0.1.8`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "latest", "The version for Bun to install (default: latest)")
	command.Flags().StringP("path", "p", "/usr/local/bin", "Installation path")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture for Bun, eg: amd64")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		fmt.Printf("Installing Bun to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) == "windows" {
			return fmt.Errorf("bun does not support Windows")
		}

		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "x64"
		}

		if version == "latest" {
			v, err := get.FindGitHubRelease("oven-sh", "bun")
			if err != nil {
				return err
			}
			version = v
		} else {
			if !strings.HasPrefix(version, "v") {
				version = "v" + version
			}
			version = "bun-" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		filename := fmt.Sprintf("bun-%s-%s.zip", strings.ToLower(osVer), dlArch)
		dlURL := fmt.Sprintf(githubDownloadTemplate, "oven-sh", "bun", version, filename)

		fmt.Printf("Downloading from: %s\n", dlURL)
		outPath, err := get.DownloadFileP(dlURL, progress)
		if err != nil {
			return err
		}
		defer os.Remove(outPath)

		fmt.Printf("Downloaded to: %s\n", outPath)

		f, err := os.OpenFile(outPath, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		tempUnpackPath, err := os.MkdirTemp(os.TempDir(), "bun*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempUnpackPath)

		fmt.Printf("Unpacking binaries to: %s\n", tempUnpackPath)
		fInfo, err := f.Stat()
		if err := archive.Unzip(f, fInfo.Size(), tempUnpackPath, true); err != nil {
			return err
		}

		fmt.Printf("Copying binaries to: %s\n", installPath)
		filesToCopy := map[string]string{
			fmt.Sprintf("%s/%s", tempUnpackPath, "bun"): fmt.Sprintf("%s/%s", installPath, "bun"),
		}
		for src, dst := range filesToCopy {
			if _, err := get.CopyFileP(src, dst, readWriteExecuteEveryone); err != nil {
				return err
			}
		}

		return nil
	}
	return command
}
