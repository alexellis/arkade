// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/alexellis/arkade/cmd/apps"
	"github.com/spf13/cobra"
)

func MakeInstall() *cobra.Command {
	var command = &cobra.Command{
		Use:   "install",
		Short: "Install Kubernetes apps from helm charts or YAML files",
		Long: `Install Kubernetes apps from helm charts or YAML files using the "install"
command. Helm 3 is used by default unless you pass --helm3=false, then helm 2
will be used to generate YAML files which are applied without tiller.

You can also find the post-install message for each app with the "info"
command.`,
		Example: `  arkade install
  arkade install openfaas --helm3 --gateways=2
  arkade install inlets-operator --token-file $HOME/do-token`,
		SilenceUsage: false,
	}

	command.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")
	command.PersistentFlags().Bool("wait", false, "If we should wait for the resource to be ready before returning (helm3 only, default false)")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf(
				`To see a complete list of apps run:

  arkade install --help

And to see options for a specific app before installing, run:

  arkade install APP --help
`)
			return nil
		}

		return nil
	}

	command.AddCommand(apps.MakeInstallOpenFaaS())
	command.AddCommand(apps.MakeInstallMetricsServer())
	command.AddCommand(apps.MakeInstallInletsOperator())
	command.AddCommand(apps.MakeInstallCertManager())
	command.AddCommand(apps.MakeInstallOpenFaaSIngress())
	command.AddCommand(apps.MakeInstallNginx())
	command.AddCommand(apps.MakeInstallChart())
	command.AddCommand(apps.MakeInstallLinkerd())
	command.AddCommand(apps.MakeInstallCronConnector())
	command.AddCommand(apps.MakeInstallKafkaConnector())
	command.AddCommand(apps.MakeInstallKubeStateMetrics())
	command.AddCommand(apps.MakeInstallMinio())
	command.AddCommand(apps.MakeInstallPostgresql())
	command.AddCommand(apps.MakeInstallKubernetesDashboard())
	command.AddCommand(apps.MakeInstallIstio())
	command.AddCommand(apps.MakeInstallCrossplane())
	command.AddCommand(apps.MakeInstallMongoDB())
	command.AddCommand(apps.MakeInstallRegistry())
	command.AddCommand(apps.MakeInstallRegistryIngress())
	command.AddCommand(apps.MakeInstallTraefik2())
	command.AddCommand(apps.MakeInstallGrafana())
	command.AddCommand(apps.MakeInstallArgoCD())
	command.AddCommand(apps.MakeInstallPortainer())
	command.AddCommand(apps.MakeInstallTekton())
	command.AddCommand(apps.MakeInstallJenkins())
	command.AddCommand(apps.MakeInstallLoki())
	command.AddCommand(apps.MakeInstallNATSConnector())
	command.AddCommand(apps.MakeInstallOpenFaaSLoki())
	command.AddCommand(apps.MakeInstallNfsProvisioner())
	command.AddCommand(apps.MakeInstallRedis())
	command.AddCommand(apps.MakeInstallOSM())
	command.AddCommand(apps.MakeInstallKubeImagePrefetch())
	command.AddCommand(apps.MakeInstallRegistryCredsOperator())

	command.AddCommand(MakeInfo())

	return command
}
