// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallSealedSecrets() *cobra.Command {
	var command = &cobra.Command{
		Use:          "sealed-secrets",
		Short:        "Install sealed-secrets",
		Long:         `Install sealed-secrets`,
		Example:      `arkade install sealed-secrets`,
		SilenceUsage: true,
	}
	command.Flags().String("namespace", "default", "Namespace for the app")

	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set secretName=secret-data)")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		wait, _ := command.Flags().GetBool("wait")

		namespace, _ := command.Flags().GetString("namespace")

		updateRepo, _ := command.Flags().GetBool("update-repo")

		overrides := map[string]string{}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		sealedSecretAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("sealed-secrets/sealed-secrets").
			WithHelmURL("https://bitnami-labs.github.io/sealed-secrets").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithWait(wait).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(sealedSecretAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(SealedSecretsPostInstallMsg)
		return nil
	}
	return command
}

const SealedSecretsPostInstallMsg = `=======================================================================
=                 The SealedSecrets app has been installed.           =
=======================================================================` +
	"\n\n" + pkg.SupportMessageShort

var SealedSecretsInfoMsg = `# Find out more on the project homepage:
# https://github.com/bitnami-labs/sealed-secrets#usage

# You can install the "kubeseal" CLI via:

arkade get kubeseal
`
