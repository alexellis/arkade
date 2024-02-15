package system

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/exp/slices"

	"github.com/Masterminds/semver"
	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

type ReferenceObject struct {
	Type string `json:"tag,omitempty"`
}

type Reference struct {
	Ref    string          `json:"ref,omitempty"`
	Url    string          `json:"url,omitempty"`
	Object ReferenceObject `json:"object,omitempty"`
}

func MakeInstallAWSCLI() *cobra.Command {
	command := &cobra.Command{
		Use:   "aws-cli",
		Short: "Install AWS CLI",
		Long:  `Install AWS CLI for interacting with Amazon Web Services APIs.`,
		Example: `  arkade system install aws-cli
  arkade system install aws-cli --version <version>`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", githubLatest, "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "/usr/local/bin", "Installation path, where a aws cli subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().StringP("work-dir", "w", "", "Working directory that installer files should be copied to (current directory if not supplied)")
	command.Flags().Bool("run-installer", true, "Whether or not arkade should run the downloaded installer")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")
		workDir, _ := cmd.Flags().GetString("work-dir")
		runInstaller, _ := cmd.Flags().GetBool("run-installer")
		fmt.Printf("Installing AWS CLI to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		if version == githubLatest {
			v, err := getAWSCLIVersion("aws", "aws-cli")
			if err != nil {
				return err
			}
			version = v
		}

		fmt.Printf("Installing version: %s for: %s\n", version, arch)

		awsCliTool := get.Tool{
			Owner:   "amazonaws",
			Repo:    "awscli",
			Name:    "awscli",
			Version: version,
			URLTemplate: `
			{{$version := .Version}}
			{{$base := printf "https://%s.%s.com" .Repo .Owner}}

			{{$ext := "zip"}}
			{{$uri := printf "%s-exe-linux-%s" .Repo .Arch}}

			{{ if HasPrefix .OS "Ming" }}
			{{$uri = "AWSCLIV2"}}
			{{$ext = "msi"}}
			{{ else if eq .OS "Darwin" -}}
			{{$uri = "AWSCLIV2"}}
			{{$ext = "pkg"}}
			{{ end -}}

			{{$base}}/{{$uri}}-{{$version}}.{{$ext}}
		`,
		}

		dUrl, err := awsCliTool.GetURL(osVer, arch, version, !progress)
		if err != nil {
			return err
		}
		fmt.Printf("Downloading from: %s\n", dUrl)

		outPath, err := get.DownloadFileP(dUrl, progress)
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded to: %s\n", outPath)

		f, err := os.OpenFile(outPath, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		if workDir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			workDir = cwd
		}

		fmt.Printf("Copying file to: %s\n", workDir)

		filename := filepath.Base(outPath)
		if _, err = get.CopyFileP(
			outPath,
			fmt.Sprintf("%s/%s", workDir, filename), readWriteExecuteEveryone,
		); err != nil {
			return err
		}

		isArchive, err := awsCliTool.IsArchive(true)
		if err != nil {
			return err
		}

		if isArchive {
			unpackPath := fmt.Sprintf("%s/awscli", workDir)
			fmt.Printf("Unpacking AWS CLI to: %s\n", unpackPath)

			fInfo, err := f.Stat()
			if err != nil {
				return err
			}
			if err := archive.Unzip(f, fInfo.Size(), unpackPath, true); err != nil {
				return err
			}

			workDir = unpackPath
		}

		if runInstaller {
			if err := runBundledInstaller(osVer, workDir, filename, installPath); err != nil {
				return err
			}
		} else {
			tpl, err := installationInstructions(osVer, workDir, filename, installPath)
			if err != nil {
				return err
			}
			fmt.Printf("\n%s", tpl)
		}

		return nil
	}

	return command
}

func getAWSCLIVersion(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/tags", owner, repo)

	client := http.Client{Timeout: time.Second * 10}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return "", err
	}

	var references []Reference
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &references); err != nil {
		return "", err
	}

	tags := make([]*semver.Version, 0)
	for _, reference := range references {
		if reference.Object.Type == "tag" {
			trimmed := strings.TrimPrefix(reference.Ref, "refs/tags/")
			v, err := semver.NewVersion(trimmed)
			if err != nil && errors.Is(err, semver.ErrInvalidSemVer) {
				continue
			}

			tags = append(tags, v)
		}
	}

	var comparer = func(a, b *semver.Version) int {
		return a.Compare(b)
	}

	slices.SortFunc(tags, comparer)
	latest := tags[len(tags)-1]

	return latest.String(), nil
}

func runBundledInstaller(osVer string, workDir string, filename string, installPath string) error {
	fmt.Printf("Running bundled installer from download URL\n")

	cmd := ""
	args := make([]string, 0)

	switch osVer {
	case "Darwin":
		cmd = "installer"
		pkgDir := fmt.Sprintf(" -pkg %s/%s", workDir, filename)
		args = append(args, pkgDir, "-target /")
	case "Linux":
		cmd = fmt.Sprintf(".%s/install", workDir)
		args = append(args, fmt.Sprintf("--bin-dir %s", installPath))
	default:
		if strings.HasPrefix(osVer, "Ming") {
			cmd = "msiexec"
			msiDir := fmt.Sprintf("%s/%s", workDir, filename)
			args = append(args, msiDir, fmt.Sprintf("INSTALLDIR=%s", installPath))
		}
	}

	installTask := execute.ExecTask{
		Command:     cmd,
		Args:        args,
		StreamStdio: false,
	}

	result, err := installTask.Execute(context.Background())
	if err != nil {
		return err
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("error running bundled installer for platform, stderr: %s", result.Stderr)
	}

	return nil
}

func installationInstructions(osVer string, workDir string, filename string, installPath string) ([]byte, error) {
	t := template.New("Installation Instructions")

	switch osVer {
	case "Darwin":
		t.Parse(`# Run the downloaded .pkg installer, you will be prompted for authorisation

sudo installer -pkg {{.WorkDir}}/{{.Filename}} -target /

# Test the binary:
aws --version
`)
	case "Linux":
		t.Parse(`# Run the downloaded script installer, you will be prompted for authorisation

sudo ./{{.WorkDir}}/awscli/install --bin-dir {{.InstallPath}}

# Test the binary:
aws --version
`)
	default:
		if strings.HasPrefix(osVer, "Ming") {
			t.Parse(`# Run the downloaded .msi installer

msiexec {{.WorkDir}}/{{.Filename}} INSTALLDIR={{.InstallPath}}

# Test the binary:
aws --version
`)
		}
	}

	var tpl bytes.Buffer
	var data = struct {
		Filename    string
		WorkDir     string
		InstallPath string
	}{
		Filename:    filename,
		WorkDir:     workDir,
		InstallPath: installPath,
	}

	if err := t.Execute(&tpl, data); err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}
