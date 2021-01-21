// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallNATSConnector() *cobra.Command {
	var natsConnectorApp = &cobra.Command{
		Use:          "nats-connector",
		Short:        "Install OpenFaaS connector for NATS",
		Long:         "Install OpenFaaS connector for NATS to invoke OpenFaaS functions using NATS.",
		Example:      "arkade install nats-connector",
		SilenceUsage: true,
	}

	natsConnectorApp.Flags().StringP("namespace", "n", "openfaas", "The namespace to install NATS connector (default: openfaas")
	natsConnectorApp.Flags().Bool("update-repo", true, "Update the helm repo")
	natsConnectorApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set topics=nats-test,)")

	natsConnectorApp.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		namespace, _ := natsConnectorApp.Flags().GetString("namespace")

		overrides := map[string]string{}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		natsConnectorOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("openfaas/nats-connector").
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(natsConnectorOptions)
		if err != nil {
			return err
		}

		println(NATSConnectorInstallMsg)
		return nil
	}

	return natsConnectorApp
}

const NATSConnectorInfoMsg = `# View the connector logs:

kubectl logs deploy/nats-connector -n openfaas -f

# Find out more on the project homepage:
https://github.com/openfaas/faas-netes/tree/master/chart/nats-connector
`

const NATSConnectorInstallMsg = `=======================================================================
= nats-connector has been installed.                                   =
=======================================================================` +
	"\n\n" + NATSConnectorInfoMsg + "\n\n" + pkg.ThanksForUsing
