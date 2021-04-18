// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallKyverno() *cobra.Command {
	var kyverno = &cobra.Command{
		Use:          "kyverno",
		Aliases:      []string{"kyverno"},
		Short:        "Install kyverno",
		Long:         `Install kyverno`,
		Example:      `  arkade install kyverno --namespace default`,
		SilenceUsage: true,
	}

	kyverno.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	kyverno.Flags().Bool("update-repo", true, "Update the helm repo")
	kyverno.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	kyverno.RunE = func(command *cobra.Command, args []string) error {
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

		kyvernoOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kyverno/kyverno").
			WithHelmURL("https://kyverno.github.io/kyverno").
			WithOverrides(overrides).
			WithWait(wait).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kyvernoOptions)

		if err != nil {
			return err
		}

		fmt.Println(KyvernoInstallMsg)

		return nil
	}

	return kyverno
}

const KyvernoInfoMsg = `Thank you for installing kyverno 😀

Your release is named kyverno.

We have installed the "default" profile of Pod Security Standards and set them in audit mode.

Visit https://kyverno.io/policies/ to find more sample policies.`

const KyvernoInstallMsg = `=======================================================================
= kyverno has been installed.                                   =
=======================================================================` +
	"\n\n" + KyvernoInfoMsg + "\n\n" + pkg.ThanksForUsing
