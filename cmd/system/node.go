package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

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

		resolver := &get.NodeVersionResolver{
			Channel: channel,
			Version: version,
		}

		if (version == "latest" || strings.Contains(version, "latest-")) && channel == "release" {
			v, err := resolver.GetVersion()
			if err != nil {
				return err
			}
			version = v
			resolver.Version = v
		} else if (version == "latest" || strings.Contains(version, "latest-")) && channel == "nightly" {
			return fmt.Errorf("please set a specific version for downloading a nightly builds")
		}

		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		tools := get.MakeTools()
		var tool *get.Tool
		for _, t := range tools {
			if t.Name == "node" {
				tool = &t
				break
			}
		}

		if tool == nil {
			return fmt.Errorf("unable to find node definition")
		}

		tool.VersionResolver = &get.NodeVersionResolver{
			Channel: channel,
			Version: version,
		}

		tempPath, err := get.DownloadNested(tool, arch, osVer, version, installPath, progress, !progress)
		defer os.RemoveAll(tempPath)
		if err != nil {
			return err
		}
		fmt.Printf("Temp Path: %s \n", tempPath)

		err = get.MoveFromInternalDir(fmt.Sprintf("node-%s-linux-%s", version, dlArch), tempPath, installPath)
		if err != nil {
			return err
		}

		return nil
	}
	return command
}
