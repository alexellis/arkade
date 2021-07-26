// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallKeda() *cobra.Command {
	var command = &cobra.Command{
		Use:          "keda",
		Short:        "Install keda",
		Long:         `Install keda, the Kubernetes Event-driven Autoscaling from the official chart https://kedacore.github.io/charts`,
		Example:      `arkade install keda`,
		SilenceUsage: true,
	}

	command.Flags().Bool("update-repo", true, "Update the helm repo")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}
		namespace, _ := command.Flags().GetString("namespace")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		kedaOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kedacore/keda").
			WithHelmUpdateRepo(updateRepo).
			WithHelmRepoVersion("2.3.2").
			WithHelmURL("https://kedacore.github.io/charts").
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kedaOptions)
		if err != nil {
			return err
		}

		fmt.Println(KedaInfoMsgInstallMsg)

		return nil
	}

	return command
}

const KedaInfoMsg = `
KEDA

KEDA is a Kubernetes-based Event Driven Autoscaler. With KEDA, you can drive the scaling of any container in Kubernetes based on the number of events needing to be processed.

KEDA is a single-purpose and lightweight component that can be added into any Kubernetes cluster. KEDA works alongside standard Kubernetes components like the Horizontal Pod Autoscaler and can extend functionality without overwriting or duplication. With KEDA you can explicitly map the apps you want to use event-driven scale, with other apps continuing to function. This makes KEDA a flexible and safe option to run alongside any number of any other Kubernetes applications or frameworks.

Features

- Event-driven 
- Autoscaling Made Simple
- Built-in Scalers
- Multiple Workload Types
- Vendor-Agnostic
- Azure Functions Support

More details

See https://keda.sh/docs/2.3/concepts/ for the usage of the different event sources and scalers.`

const KedaInfoMsgInstallMsg = `=======================================================================
= KEDA has been installed                                             =
=======================================================================` +
	"\n\n" + KedaInfoMsg + "\n\n" + pkg.ThanksForUsing
