// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallKafkaConnector() *cobra.Command {
	var command = &cobra.Command{
		Use:          "kafka-connector",
		Short:        "Install kafka-connector for OpenFaaS",
		Long:         `Install kafka-connector for OpenFaaS`,
		Example:      `  arkade install kafka-connector`,
		SilenceUsage: true,
	}

	command.Flags().StringP("namespace", "n", "openfaas", "The namespace used for installation")
	command.Flags().Bool("update-repo", true, "Update the helm repo")
	command.Flags().StringP("topics", "t", "faas-request", "The topics for the connector to bind to")
	command.Flags().String("broker-host", "kafka", "The host for the Kafka broker")
	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		updateRepo, _ := command.Flags().GetBool("update-repo")

		namespace, _ := command.Flags().GetString("namespace")

		topicsVal, err := command.Flags().GetString("topics")
		if err != nil {
			return err
		}

		brokerHostVal, err := command.Flags().GetString("broker-host")
		if err != nil {
			return err
		}

		overrides := map[string]string{
			"topics":      topicsVal,
			"broker_host": brokerHostVal,
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)
		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		kafkaConnectorAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("openfaas/kafka-connector").
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(kafkaConnectorAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(kafkaConnectorInstallMsg)

		return nil
	}

	return command
}

const KafkaConnectorInfoMsg = `# View the connector's logs:

kubectl logs deploy/kafka-connector -n openfaas -f

# Find out more on the project homepage:

# https://github.com/openfaas-incubator/kafka-connector/`

const kafkaConnectorInstallMsg = `=======================================================================
= kafka-connector has been installed.                                   =
=======================================================================` +
	"\n\n" + KafkaConnectorInfoMsg + "\n\n" + pkg.ThanksForUsing
