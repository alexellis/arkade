// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
)

var linkerdVersion = "stable-2.6.1"

func MakeInstallLinkerd() *cobra.Command {
	var linkerd = &cobra.Command{
		Use:          "linkerd",
		Short:        "Install linkerd",
		Long:         `Install linkerd`,
		Example:      `  arkade install linkerd`,
		SilenceUsage: true,
	}

	linkerd.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)
		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		userPath, err := getUserPath()
		if err != nil {
			return err
		}

		_, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q\n", clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		err = downloadLinkerd(userPath, clientOS)
		if err != nil {
			return err
		}
		fmt.Println("Running linkerd check, this may take a few moments.")

		_, err = linkerdCli("check", "--pre")
		if err != nil {
			return err
		}

		res, err := linkerdCli("install")
		if err != nil {
			return err
		}
		file, err := ioutil.TempFile("", "linkerd")
		if err != nil {
			return err
		}

		w := bufio.NewWriter(file)
		_, err = w.WriteString(res.Stdout)
		if err != nil {
			return err
		}
		w.Flush()

		err = k8s.Kubectl("apply", "-R", "-f", file.Name())
		if err != nil {
			return err
		}

		defer os.Remove(file.Name())

		_, err = linkerdCli("check")
		if err != nil {
			return err
		}

		fmt.Println(`=======================================================================
= Linkerd has been installed.                                        =
=======================================================================

# Get the linkerd-cli
curl -sL https://run.linkerd.io/install | sh

# Find out more at:
# https://linkerd.io

# To use the Linkerd CLI set this path:

export PATH=$PATH:` + path.Join(userPath, "bin/") + `
linkerd --help

` + pkg.ThanksForUsing)
		return nil
	}

	return linkerd
}

func getLinkerdURL(os, version string) string {
	osSuffix := strings.ToLower(os)
	return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-%s", version, version, osSuffix)
}

func downloadLinkerd(userPath, clientOS string) error {
	filePath := path.Join(path.Join(userPath, "bin"), "linkerd")
	if _, statErr := os.Stat(filePath); statErr != nil {
		linkerdURL := getLinkerdURL(clientOS, linkerdVersion)
		fmt.Println(linkerdURL)
		parsedURL, _ := url.Parse(linkerdURL)

		res, err := http.DefaultClient.Get(parsedURL.String())
		if err != nil {
			return err
		}

		defer res.Body.Close()
		out, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, res.Body)
		if err != nil {
			return err
		}

		err = os.Chmod(filePath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func linkerdCli(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s", env.LocalBinary("linkerd", "")),
		Args:        parts,
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return res, err
	}

	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	return res, nil
}

func getUserPath() (string, error) {
	userPath, err := config.InitUserDir()
	return userPath, err
}

func getExportPath() string {
	userPath := config.GetUserDir()
	return path.Join(userPath, "bin/")
}

var LinkerdInfoMsg = `# Get the linkerd-cli
curl -sL https://run.linkerd.io/install | sh

# Find out more at:
# https://linkerd.io

# To use the Linkerd CLI set this path:

export PATH=$PATH:` + getExportPath() + `
linkerd --help`

var linkerdInstallMsg = `=======================================================================
= Linkerd has been installed.                                         =
=======================================================================` +
	"\n\n" + LinkerdInfoMsg + "\n\n" + pkg.ThanksForUsing
