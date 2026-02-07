package system

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallZvolSnapshotter() *cobra.Command {
	command := &cobra.Command{
		Use:   "zvol-snapshotter",
		Short: "Install containerd zvol snapshotter",
		Long:  "Install the ZFS Volume snapshotter plugin for containerd",
		Example: `  arkade system install zvol-snapshotter \
    --dataset tank/snapshots
  arkade system install zvol-snapshotter --version <version> \
    --dataset tank/snapshots`,
		SilenceUsage: true,
		RunE:         runInstallZvolSnapshotter,
	}

	command.Flags().StringP("version", "v", githubLatest, "Zvol snapshotter version to install")
	command.Flags().StringP("path", "p", "/usr/local/bin", "Installation path where the binary will be installed")
	command.Flags().String("dataset", "", "ZFS dataset that will be used for snapshots")
	command.Flags().String("size", "20G", "Space to allocate when creating volumes")
	command.Flags().Bool("systemd", true, "Add and enable systemd service")
	command.Flags().Bool("progress", true, "Show download progress")

	return command
}

func runInstallZvolSnapshotter(cmd *cobra.Command, args []string) error {
	installPath, _ := cmd.Flags().GetString("path")
	version, _ := cmd.Flags().GetString("version")
	progress, _ := cmd.Flags().GetBool("progress")
	systemd, _ := cmd.Flags().GetBool("systemd")

	volumeSize, _ := cmd.Flags().GetString("size")
	dataset, _ := cmd.Flags().GetString("dataset")
	if dataset == "" {
		return fmt.Errorf("please configure a ZFS dataset for snapshots using the --dataset flag")
	}

	owner := "welteki"
	repo := "zvol-snapshotter"

	arch, osVer := env.GetClientArch()

	if cmd.Flags().Changed("os") {
		osVer, _ = cmd.Flags().GetString("os")
	}
	if cmd.Flags().Changed("arch") {
		arch, _ = cmd.Flags().GetString("arch")
	}

	if strings.ToLower(osVer) != "linux" {
		return fmt.Errorf("this app only supports Linux")
	}

	if arch != "x86_64" && arch != "aarch64" {
		return fmt.Errorf("this app only supports x86_64 and aarch64 and not %s", arch)
	}

	dlArch := arch
	if arch == "x86_64" {
		dlArch = "amd64"
	}
	if arch == "aarch64" {
		dlArch = "arm64"
	}

	if version == githubLatest {
		v, err := get.FindGitHubRelease(owner, repo)
		if err != nil {
			return err
		}

		version = v
	} else if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	dlVersion := strings.TrimPrefix(version, "v")

	fmt.Printf("Installing zvol snapshotter version: %s for: %s, to: %s\n", version, arch, installPath)

	filename := fmt.Sprintf("zvol-snapshotter-%s-linux-%s.tar.gz", dlVersion, dlArch)
	dlURL := fmt.Sprintf(githubDownloadTemplate, owner, repo, version, filename)

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

	tempUnpackPath, err := os.MkdirTemp(os.TempDir(), "firecracker-*")
	if err != nil {
		return err
	}
	fmt.Printf("Unpacking zvol-snapshotter to: %s\n", tempUnpackPath)
	if err := spinWhile("Unpacking zvol-snapshotter", func() error {
		return archive.Untar(f, tempUnpackPath, true, true)
	}); err != nil {
		return err
	}

	fmt.Printf("Copying zvol-snapshotter binaries to: %s\n", installPath)
	filesToCopy := map[string]string{
		fmt.Sprintf("%s/containerd-zvol-grpc", tempUnpackPath): fmt.Sprintf("%s/containerd-zvol-grpc", installPath),
	}
	for src, dst := range filesToCopy {
		if _, err := get.CopyFileP(src, dst, readWriteExecuteEveryone); err != nil {
			return err
		}
	}

	confDir := "/etc/containerd-zvol-grpc"
	if err := createZvolSnapshotterConf(confDir, dataset, volumeSize); err != nil {
		return fmt.Errorf("failed to create zvol-snapshotter configuration: %w", err)
	}

	if systemd {
		systemdUnitName := "zvol-snapshotter.service"
		systemdUnitUrl := fmt.Sprintf("https://raw.githubusercontent.com/welteki/zvol-snapshotter/refs/tags/%s/scripts/config/zvol-snapshotter.service", version)
		fmt.Printf("Downloading zvol-snapshotter.service file from %s\n", systemdUnitUrl)

		svcTmpPath, err := get.DownloadFileP(systemdUnitUrl, false)
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded zvol-snapshotter.service file to %s\n", svcTmpPath)
		defer os.Remove(svcTmpPath)

		// Overwrite the binary path in systemd unit file if the binary has been installed
		// to a different location.
		if cmd.Flags().Changed("path") {
			content, err := os.ReadFile(svcTmpPath)
			if err != nil {
				return err
			}

			oldPath := "/usr/local/bin/containerd-zvol-grpc"
			newPath := path.Join(installPath, "containerd-zvol-grpc")
			updatedContent := strings.ReplaceAll(string(content), oldPath, newPath)

			if err := os.WriteFile(svcTmpPath, []byte(updatedContent), 0644); err != nil {
				return err
			}
		}

		systemdUnitFile := path.Join("/etc/systemd/system", systemdUnitName)
		if _, err = get.CopyFile(svcTmpPath, systemdUnitFile); err != nil {
			return err
		}
		fmt.Printf("Copied zvol-snapshotter.service file to %s\n", systemdUnitFile)

		if _, err = executeShellCmd(context.Background(), "systemctl", "daemon-reload"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "systemctl", "enable", "zvol-snapshotter"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "systemctl", "enable", "zvol-snapshotter"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "systemctl", "start", "zvol-snapshotter"); err != nil {
			return err
		}
	}

	return nil
}

func createZvolSnapshotterConf(confDir string, dataset string, volumeSize string) error {
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}

	configFile, err := os.Create(filepath.Join(confDir, "config.toml"))
	if err != nil {
		return err
	}
	defer configFile.Close()

	configTmpl, err := template.New("zvol-snapshotter-config").Parse(`# Snapshotter root directory for metadata
root_path="/var/lib/containerd-zvol-grpc"
# ZFS dataset that will be used for snapshots
dataset="{{.Dataset}}"
# Space to allocate when creating volumes
volume_size="{{.VolumeSize}}"
# File system to use for snapshot device mounts
fs_type="ext4"
`)
	if err != nil {
		return err
	}

	if err := configTmpl.Execute(configFile, struct {
		Dataset    string
		VolumeSize string
	}{
		Dataset:    dataset,
		VolumeSize: volumeSize,
	}); err != nil {
		return err
	}

	return nil
}
