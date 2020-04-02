// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func MakeInstallRancher() *cobra.Command {
	var command = &cobra.Command{
		Use:          "rancher",
		Short:        "Install rancher to import and manage Kubernetes",
		Long:         `Install rancher to import and manage Kubernetes clsuter.`,
		Example:      `  arkade install rancher`,
		SilenceUsage: true,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch && arch != "arm" {
			return fmt.Errorf(`only Intel and "arm" is supported for this app`)
		}

		_, err := kubectlTask("create", "ns",
			"cattle-system")
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}

		req, err := http.NewRequest(http.MethodGet,
			"https://raw.githubusercontent.com/saiyam1814/rancher-deploy/master/rancher.yaml",
			nil)

		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		manifest := string(body)
		if arch == "arm" {
			manifest = strings.Replace(manifest, "linux-amd64", "linux-arm", -1)
		}

		tmp := os.TempDir()
		joined := path.Join(tmp, "rancher.yaml")
		err = ioutil.WriteFile(joined, []byte(manifest), 0644)
		if err != nil {
			return err
		}

		_, err = kubectlTask("apply", "-f", joined, "-n", "cattle-system")
		if err != nil {
			return err
		}

		fmt.Println(RancherInstallMsg)

		return nil
	}

	return command
}

const RancherInfoMsg = `

# Access Rancher Dashboard via NodePort on https://node-ip:30801

Find out more at https://www.rancher.com/
`

const RancherInstallMsg = `=======================================================================
= Rancher has been installed                                        =
=======================================================================` +
	"\n\n" + RancherInfoMsg + "\n\n" + pkg.ThanksForUsing
