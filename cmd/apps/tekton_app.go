// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"

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
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		fmt.Println("Installing Tekton pipelines...")
		_, err := k8s.KubectlTask("apply", "-f",
			"https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml")
		if err != nil {
			return err
		}

		fmt.Println("Installing Tekton dashboard...")
		_, err = k8s.KubectlTask("apply", "-f",
			"https://storage.googleapis.com/tekton-releases/dashboard/latest/tekton-dashboard-release.yaml")
		if err != nil {
			return err
		}

		fmt.Println(TektonInstallMsg)

		return nil
	}

	return tekton
}

const TektonInfoMsg = `Want to launch a Tekton pipeline and see it in action?
Get one here: https://github.com/tektoncd/pipeline/tree/master/examples

For more information...
  https://github.com/tektoncd/pipeline
  https://github.com/tektoncd/dashboard`

const TektonInstallMsg = `=======================================================================
= Tekton pipelines and dashboard have been installed. =
=======================================================================` +
	"\n\n" + TektonInfoMsg + "\n\n" + pkg.SupportMessageShort
