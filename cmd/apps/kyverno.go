// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallKyverno() *cobra.Command {
	var kyverno = &cobra.Command{
		Use:          "kyverno",
		Short:        "Install Kyverno",
		Long:         `Install Kyverno, which is a Kubernetes Native Policy Management engine`,
		Example:      `arkade install kyverno --set helmKey=value`,
		SilenceUsage: true,
	}

	kyverno.Flags().StringP("namespace", "n", "kyverno", "The namespace used for installation")
	kyverno.Flags().Bool("update-repo", true, "Update the helm repo")
	kyverno.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set serviceMonitor.enabled=true)")

	kyverno.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetBool("update-repo")
		if err != nil {
			return fmt.Errorf("error with --update-repo usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		return nil
	}

	kyverno.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		customFlags, _ := command.Flags().GetStringArray("set")
		namespace, _ := command.Flags().GetString("namespace")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		overrides := map[string]string{}
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kyvernoOptions := types.DefaultInstallOptions().
			WithHelmRepo("kyverno/kyverno").
			WithHelmURL("https://kyverno.github.io/kyverno/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithInstallNamespace(true).
			WithNamespace(namespace).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kyvernoOptions)
		if err != nil {
			return err
		}

		fmt.Println(kyvernoInstallMsg)
		return nil
	}

	return kyverno
}

var KyvernoInfoMsg = `
Kyverno is a Kubernetes Native Policy Management engine. It allows you to:

Manage policies as Kubernetes resources (no new language required.)
Validate, mutate, and generate resource configurations.
Select resources based on labels and wildcards.
View policy enforcement as events.
Scan existing resources for violations.
Access the complete user documentation and guides at: https://kyverno.io/
`

var kyvernoInstallMsg = `=======================================================================
= Kyverno has been installed.                                         =
=======================================================================` +
	"\n\n" + KyvernoInfoMsg + "\n\n" + pkg.SupportMessageShort
