// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// kasten contains a suite of Sponsored Apps for arkade
package kasten

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallK10() *cobra.Command {
	var k10cmd = &cobra.Command{
		Use:          "k10",
		Short:        "Install K10",
		Long:         `Install K10 - backup and restore for Kubernetes.`,
		Example:      `arkade install k10`,
		SilenceUsage: true,
	}

	k10cmd.Flags().StringP("namespace", "n", "kasten-io", "The namespace used for installation")
	k10cmd.Flags().Bool("update-repo", true, "Update the helm repo")
	k10cmd.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	k10cmd.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		wait, _ := command.Flags().GetBool("wait")
		namespace, _ := command.Flags().GetString("namespace")
		overrides := map[string]string{}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		k10cmdOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kasten/k10").
			WithHelmURL("https://charts.kasten.io/").
			WithOverrides(overrides).
			WithWait(wait).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(k10cmdOptions)

		if err != nil {
			return err
		}

		fmt.Println(k10InstallCmd)

		return nil
	}

	return k10cmd
}

const k10InfoMsg = `# The K10 app has been installed

# You may also need to install pre-requisites and configure a 
# CSI drive for your cluster.

https://docs.kasten.io/latest/install/storage.html

# The app may take a few moments to come up, then access the
# dashboard via:

kubectl --namespace kasten-io port-forward service/gateway 8080:8000

http://127.0.0.1:8080/k10/#/

# Find out your next steps here:

https://docs.kasten.io/latest/install/install.html
`

const k10InstallCmd = `=======================================================================
= k10 has been installed.                                   =
=======================================================================` +
	"\n\n" + k10InfoMsg + "\n\n" + pkg.ThanksForUsing

func mergeFlags(existingMap map[string]string, setOverrides []string) error {
	for _, setOverride := range setOverrides {
		flag := strings.Split(setOverride, "=")
		if len(flag) != 2 {
			return fmt.Errorf("incorrect format for custom flag `%s`", setOverride)
		}
		existingMap[flag[0]] = flag[1]
	}
	return nil
}
