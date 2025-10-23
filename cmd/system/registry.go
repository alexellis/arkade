package system

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	execute "github.com/alexellis/go-execute/v2"
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
	cmd.Flags().Bool("overwrite", false, "Overwrite existing binary if found")
	cmd.Flags().String("arch", "", "CPU architecture i.e. amd64")
	cmd.Flags().String("type", "", "Type of registry - '' means binary only. 'mirror' creates a pull through cache.")
	cmd.Flags().String("docker-password", "", "Password for Docker Hub registry authentication (only for 'mirror' type)")
	cmd.Flags().String("docker-password-file", "", "Path to a file containing access token for registry authentication (only for 'mirror' type)")
	cmd.Flags().String("docker-username", "", "Username for registry authentication (only for 'mirror' type)")
	cmd.Flags().String("bind-addr", "0.0.0.0", "Bind address for the registry server (only for 'mirror' type)")
	cmd.Flags().String("storage", "/var/lib/registry", "Path to registry storage (only for 'mirror' type)")
	cmd.Flags().String("tls", "", "Give \"actuated\" or leave empty.")

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

		installType, _ := cmd.Flags().GetString("type")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

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

		if _, err := os.Stat(filepath.Join(installPath, toolName)); os.IsNotExist(err) || overwrite {
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
			if err := archive.UntarNested(f, tempDirName, true, false); err != nil {
				return err
			}

			fmt.Printf("Copying %s binary to: %s\n", toolName, installPath)

			if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
				fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
			}

			if _, err := get.CopyFile(fmt.Sprintf("%s/%s", tempDirName, toolName), fmt.Sprintf("%s/%s", installPath, toolName)); err != nil {
				if strings.Contains(err.Error(), "text file busy") {
					return fmt.Errorf("stop any running \"%s\" processes, error: %w", toolName, err)
				}
				return err
			}
		} else {
			fmt.Printf("%s already exists in %s, skipping download and install.\n", toolName, installPath)
		}

		if installType == "mirror" {
			accessToken, _ := cmd.Flags().GetString("access-token")
			accessTokenFile, _ := cmd.Flags().GetString("access-token-file")
			bindAddr, _ := cmd.Flags().GetString("bind-addr")
			storage, _ := cmd.Flags().GetString("storage")
			username, _ := cmd.Flags().GetString("username")
			tls, _ := cmd.Flags().GetString("tls")

			fmt.Printf("Setting up registry mirror service\n")
			if err := setupRegistryMirrorService(installPath, accessToken, accessTokenFile, bindAddr, storage, username, tls); err != nil {
				return err
			}
			fmt.Printf(`View logs:

  sudo journalctl -u registry-mirror.service -f
`)
		}

		return nil
	}

	return cmd
}

func setupRegistryMirrorService(installPath, accessToken, accessTokenFile, bindAddr, storagePath, username, tls string) error {

	registryEtc := "/etc/registry"
	os.MkdirAll(registryEtc, 0755)
	os.MkdirAll(storagePath, 0755)

	cfg := template.New("registry-config")
	tmpl, err := cfg.Parse(configTmp)
	if err != nil {
		return err
	}

	token := accessToken
	if accessTokenFile != "" {
		data, err := os.ReadFile(accessTokenFile)
		if err != nil {
			return err
		}
		token = strings.TrimSpace(string(data))
	}

	if len(token) > 0 && len(username) == 0 {
		return fmt.Errorf("username must be provided when access token is set")
	}

	bindAddr, port, found := strings.Cut(bindAddr, ":") // remove port if present
	if !found {
		port = "5000"
	}

	bindAddr = fmt.Sprintf("%s:%s", strings.TrimRight(bindAddr, ":"), port)

	buf := &strings.Builder{}
	if err := tmpl.Execute(buf, map[string]string{
		"USERNAME":   username,
		"TOKEN":      token,
		"BRIDGE":     bindAddr,
		"TLS":        tls,
		"REMOTE_URL": "https://registry-1.docker.io",
	}); err != nil {
		return err
	}

	if err := os.WriteFile(fmt.Sprintf("%s/config.yml", registryEtc), []byte(buf.String()), 0644); err != nil {
		return err
	}

	fmt.Printf("Wrote: %s\n", fmt.Sprintf("%s/config.yml", registryEtc))

	reg := template.New("registry-service")
	svcTmpl, err := reg.Parse(registrySvcTmpl)
	if err != nil {
		return err
	}

	svcBuf := &strings.Builder{}
	if err := svcTmpl.Execute(svcBuf, map[string]string{
		"ConfigPath": fmt.Sprintf("%s/config.yml", registryEtc),
		"TLS":        tls,
	}); err != nil {
		return err
	}

	servicePath := "/etc/systemd/system/registry-mirror.service"
	if err := os.WriteFile(servicePath, []byte(svcBuf.String()), 0644); err != nil {
		return err
	}

	fmt.Printf("Wrote: %s\n", servicePath)

	fmt.Printf("Starting \"registry-mirror\" service\n")
	taskReload := execute.ExecTask{
		Command:     "systemctl",
		Args:        []string{"daemon-reload"},
		StreamStdio: false,
	}
	if _, err := taskReload.Execute(context.Background()); err != nil {
		return err
	}

	taskEnable := execute.ExecTask{
		Command:     "systemctl",
		Args:        []string{"enable", "registry-mirror.service", "--now"},
		StreamStdio: true,
	}
	if _, err := taskEnable.Execute(context.Background()); err != nil {
		return err
	}

	viewLogLines := 10
	viewLogs := execute.ExecTask{
		Command:     "journalctl",
		Args:        []string{"-u", "registry-mirror.service", "--lines", fmt.Sprintf("%d", viewLogLines)},
		StreamStdio: true,
	}
	if _, err := viewLogs.Execute(context.Background()); err != nil {
		return err
	}

	return nil

}

var configTmp = `version: 0.1
log:
  accesslog:
    disabled: true
  level: warn
  formatter: text

storage:
  filesystem:
    rootdirectory: /var/lib/registry

proxy:
  remoteurl: {{ .REMOTE_URL }}
{{- if ne .USERNAME "" }}
  username: {{ .USERNAME }}

  # A Docker Hub Personal Access token created with "Public repos only" scope
  password: {{ .TOKEN }}
{{- end }}

http:
  addr: {{ .BRIDGE }}
  relativeurls: false
  draintimeout: 60s

{{- if eq .TLS "actuated" }}
  # Enable self-signed TLS from the TLS certificate and key
  # managed by actuated for server <> microVM communication
  tls:
    certificate: /var/lib/actuated/certs/server.crt
    key: /var/lib/actuated/certs/server.key
{{- end }}
`

var registrySvcTmpl = `[Unit]
Description=Registry
After=network.target {{ if eq .TLS "actuated" }}actuated.service{{ end }}

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/usr/local/bin/registry serve {{.ConfigPath}}

[Install]
WantedBy=multi-user.target
`
