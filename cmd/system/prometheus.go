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
		Use:     "prometheus",
		Short:   "Install Prometheus",
		Long:    `Install the Prometheus monitoring system and time series database.`,
		Aliases: []string{"prom"},
		Example: `  # Install Prometheus to the default path
  arkade system install prometheus

  # Install a specific version
  arkade system install prometheus --version v2.34.0

  # Install on macOS
  arkade system install prometheus --os darwin --arch arm64

  # Print the resolved version to stdout and exit
  arkade system install prometheus --print-version`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "latest", "The version for Prometheus to install")
	command.Flags().StringP("path", "p", "/usr/local/bin", "Installation path for the prometheus and promtool binaries")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. x86_64 or arm64")
	command.Flags().String("os", "", "Operating system i.e. linux or darwin, leave blank to detect the client OS")
	command.Flags().Bool("print-version", false, "Print the resolved version to stdout and exit")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		if version == "latest" {
			v, err := get.FindGitHubRelease("prometheus", "prometheus")
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

		fmt.Printf("Installing Prometheus to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if cmd.Flags().Changed("os") {
			osVer, _ = cmd.Flags().GetString("os")
		}
		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlOS := strings.ToLower(osVer)
		if dlOS != "linux" && dlOS != "darwin" {
			return fmt.Errorf("unsupported operating system: %q, use linux or darwin (macOS)", osVer)
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		}

		fmt.Printf("Installing version: %s for: %s / %s\n", version, dlOS, dlArch)

		filename := fmt.Sprintf("prometheus-%s.%s-%s.tar.gz", strings.TrimPrefix(version, "v"), dlOS, dlArch)
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
		if err := spinWhile("Unpacking Prometheus", func() error {
			return archive.Untar(f, tempUnpackPath, true, true)
		}); err != nil {
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
