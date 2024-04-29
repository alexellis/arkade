package system

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

func getLatestNodeVersion(version, channel string) (*string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://nodejs.org/download/%s/%s", channel, version), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not find latest version for %s, (%d), body: %s", version, res.StatusCode, string(body))
	}

	regex := regexp.MustCompile(`(?m)node-v(\d+.\d+.\d+)-linux-.*`)
	result := regex.FindStringSubmatch(string(body))

	if len(result) < 2 {
		if v, ok := os.LookupEnv("ARK_DEBUG"); ok && v == "1" {
			fmt.Printf("Body: %s\n", string(body))
		}
		return nil, fmt.Errorf("could not find latest version for %s, (%d), %s", version, res.StatusCode, result)
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

	command.Flags().StringP("version", "v", "latest", "The version for Node.js to install, either a specific version, 'latest' or 'latest-CODENAME' (eg: latest-hydrogen)")
	command.Flags().StringP("path", "p", "/usr/local/", "Installation path")
	command.Flags().StringP("channel", "c", "release", "The channel to install from, can be 'release' or 'nightly',")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")
		channel, _ := cmd.Flags().GetString("channel")

		fmt.Printf("Installing Node.js to: %s\n", installPath)

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
		if err = archive.UntarNested(f, tempUnpackPath, true, false); err != nil {
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
