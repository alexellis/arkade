// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallPortainer() *cobra.Command {
	var command = &cobra.Command{
		Use:          "portainer",
		Short:        "Install portainer to visualise and manage containers",
		Long:         `Install portainer to visualise and manage containers, now in beta for Kubernetes.`,
		Example:      `  arkade install portainer`,
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
			"portainer")
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}

		req, err := http.NewRequest(http.MethodGet,
			"https://raw.githubusercontent.com/portainer/portainer-k8s/master/portainer-nodeport.yaml",
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
		joined := path.Join(tmp, "portainer.yaml")
		err = ioutil.WriteFile(joined, []byte(manifest), 0644)
		if err != nil {
			return err
		}

		_, err = kubectlTask("apply", "-f", joined, "-n", "portainer")
		if err != nil {
			return err
		}

		fmt.Println(PortainerInstallMsg)

		return nil
	}

	return command
}

const PortainerInfoMsg = `
# Open the UI:

kubectl port-forward -n portainer svc/portainer 9000:9000 &

# http://127.0.0.1:9000

# Or access via NodePort on http://node-ip:30777

Find out more at https://www.portainer.io/
`

const PortainerInstallMsg = `=======================================================================
= Portainer has been installed                                        =
=======================================================================` +
	"\n\n" + PortainerInfoMsg + "\n\n" + pkg.ThanksForUsing
