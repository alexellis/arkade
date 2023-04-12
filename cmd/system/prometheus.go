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

func MakeInstallPrometheus() *cobra.Command {
	command := &cobra.Command{
		Use:   "prometheus",
		Short: "Install Prometheus",
		Long:  `Install the Prometheus monitoring system and time series database.`,
		Example: `  arkade system install prometheus
  arkade system install prometheus --version v2.34.0`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "latest", "The version for Prometheus to install")
	command.Flags().StringP("path", "p", "/usr/local/bin", "Installation path, where a go subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		fmt.Printf("Installing Prometheus to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}
		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		}

		if version == "latest" {
			v, err := get.FindGitHubRelease("prometheus", "prometheus")
			if err != nil {
				return err
			}
			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		filename := fmt.Sprintf("prometheus-%s.linux-%s.tar.gz", strings.TrimPrefix(version, "v"), dlArch)
		dlURL := fmt.Sprintf(githubDownloadTemplate, "prometheus", "prometheus", version, filename)

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

		tempUnpackPath, err := os.MkdirTemp(os.TempDir(), "prometheus*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempUnpackPath)

		fmt.Printf("Unpacking binaries to: %s\n", tempUnpackPath)
		if err := archive.Untar(f, tempUnpackPath, true, true); err != nil {
			return err
		}

		fmt.Printf("Copying binaries to: %s\n", installPath)
		filesToCopy := map[string]string{
			fmt.Sprintf("%s/%s", tempUnpackPath, "prometheus"): fmt.Sprintf("%s/%s", installPath, "prometheus"),
			fmt.Sprintf("%s/%s", tempUnpackPath, "promtool"):   fmt.Sprintf("%s/%s", installPath, "promtool"),
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
