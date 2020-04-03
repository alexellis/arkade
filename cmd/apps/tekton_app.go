// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallTekton() *cobra.Command {
	var tekton = &cobra.Command{
		Use:          "tekton",
		Short:        "Install Tekton pipelines and dashboard",
		Long:         `Install Tekton pipelines and dashboard`,
		Example:      `  arkade install tekton`,
		SilenceUsage: true,
	}

	tekton.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(`only Intel and AMD (i.e. PC) architecture is supported for this app`)
		}

		fmt.Println("Installing Tekton pipelines...")
		_, err := kubectlTask("apply", "-f",
			"https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml")
		if err != nil {
			return err
		}

		fmt.Println("Installing Tekton dashboard...")
		_, err = kubectlTask("apply", "-f",
			"https://github.com/tektoncd/dashboard/releases/download/v0.5.1/tekton-dashboard-release.yaml")
		if err != nil {
			return err
		}

		fmt.Println(TektonInstallMsg)

		return nil
	}

	return tekton
}

const TektonInfoMsg = `For more information on accessing the
Tekton dashboard: https://github.com/tektoncd/dashboard

Want to launch a Tekton pipeline and see it in action?
Get one here: https://github.com/tektoncd/pipeline/tree/master/examples`

const TektonInstallMsg = `=======================================================================
= Tekton pipelines and dashboard have been installed. =
=======================================================================` +
	"\n\n" + TektonInfoMsg + "\n\n" + pkg.ThanksForUsing
