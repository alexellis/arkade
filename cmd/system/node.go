package system

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func getLatestNodeVersion(version, channel string) (*string, error) {
	res, err := http.Get(fmt.Sprintf("https://nodejs.org/download/%s/%s", channel, version))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(`(?m)node-(.*)-linux-.*"`)
	result := regex.FindStringSubmatch(string(body))
	if len(result) < 2 {
		return nil, fmt.Errorf("could not find latest version for %s", version)
	}
	return &result[1], nil
}

func MakeInstallNode() *cobra.Command {
	command := &cobra.Command{
		Use:   "node",
		Short: "Install Node.js",
		Long:  `Node.js is a JavaScript runtime built on Chrome's V8 JavaScript engine.`,
		Example: `arkade system install node
  arkade system install node --version v17.9.0`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "latest", "The version for Node.js to install, either a specific version, 'latest' or 'latest-CODENAME' (eg: latest-gallium)")
	command.Flags().StringP("path", "p", "/usr/local/", "Installation path")
	command.Flags().StringP("channel", "c", "release", "The channel to install from, can be 'releases' or 'nightly',")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture for Prometheus, eg: amd64")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")
		channel, _ := cmd.Flags().GetString("channel")

		fmt.Printf("Installing Node.js to %s\n", installPath)

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}
		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "x64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		}

		if (version == "latest" || strings.Contains(version, "latest-")) && channel == "release" {
			v, err := getLatestNodeVersion(version, channel)
			if err != nil {
				return err
			}
			version = *v
		} else if (version == "latest" || strings.Contains(version, "latest-")) && channel == "nightly" {
			return fmt.Errorf("please set a specific version for downloading a nightly builds")
		}

		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)
		filename := fmt.Sprintf("%s/%s.tar.gz", version, fmt.Sprintf("node-%s-linux-%s", version, dlArch))
		dlURL := fmt.Sprintf("https://nodejs.org/download/%s/%s", channel, filename)

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

		tempUnpackPath, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("%s*", "node"))
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempUnpackPath)
		fmt.Printf("Unpacking binaries to: %s\n", tempUnpackPath)
		if err = archive.UntarNested(f, tempUnpackPath); err != nil {
			return err
		}

		fmt.Printf("Copying binaries to: %s\n", installPath)
		nodeDir := fmt.Sprintf("%s/%s", tempUnpackPath, fmt.Sprintf("node-%s-linux-%s", version, dlArch))
		if err := cp.Copy(nodeDir, installPath); err != nil {
			return err
		}
		return nil
	}
	return command
}
