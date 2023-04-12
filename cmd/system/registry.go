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

func MakeInstallRegistry() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "registry",
		Short:        "Install registry",
		Long:         "Install registry Open Source Registry implementation for storing and distributing container images using the OCI Distribution Specification.",
		Example:      `arkade system install registry`,
		SilenceUsage: true,
	}

	cmd.Flags().StringP("version", "v", "", "Version of the registry binary pack, leave blank for latest")
	cmd.Flags().String("path", "/usr/local/bin", "Install path, where the distribution binaries will installed")
	cmd.Flags().Bool("progress", true, "Show download progress")
	cmd.Flags().String("arch", "", "CPU architecture i.e. amd64")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("progress")
		if err != nil {
			return err
		}

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		arch, osVer := env.GetClientArch()
		if cmd.Flags().Changed("arch") {
			archFlag, _ := cmd.Flags().GetString("arch")
			arch = archFlag
		}

		toolName := "registry"
		fmt.Printf("Installing %s to %s\n", toolName, installPath)

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app currently only supports Linux")
		}

		if version == "" {
			latestVerison, err := get.FindGitHubRelease("distribution", "distribution")
			if err != nil {
				return err
			}
			version = latestVerison
		}

		downloadArch := ""

		if arch == "x86_64" {
			downloadArch = "amd64"
		} else if arch == "aarch64" {
			downloadArch = "arm64"
		} else {
			return fmt.Errorf("this app currently only supports arm64 and amd64 archs")
		}

		containerdTool := get.Tool{
			Name:    toolName,
			Repo:    "distribution",
			Owner:   "distribution",
			Version: version,
			BinaryTemplate: `
			{{$archStr := .Arch}}
			{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{$archStr = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "armv7"}}
			{{- else if eq .Arch "armv6l" -}}
			{{$arch = "armv6"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$archStr = "amd64"}}
			{{- end -}}
			{{.Name}}_{{.VersionNumber}}_{{.OS}}_{{$archStr}}.tar.gz
			`,
		}

		url, err := containerdTool.GetURL(osVer, downloadArch, containerdTool.Version, !progress)
		if err != nil {
			return err
		}

		outPath, err := get.DownloadFileP(url, progress)
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

		tempDirName := fmt.Sprintf("%s/%s", os.TempDir(), toolName)
		defer os.RemoveAll(tempDirName)
		if err := archive.UntarNested(f, tempDirName); err != nil {
			return err
		}

		fmt.Printf("Copying %s binary to: %s\n", toolName, installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		_, err = get.CopyFile(fmt.Sprintf("%s/%s", tempDirName, toolName), fmt.Sprintf("%s/%s", installPath, toolName))
		if err != nil {
			return err
		}

		return nil
	}

	return cmd
}
