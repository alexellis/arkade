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

func MakeInstallKafkaConnector() *cobra.Command {
	var command = &cobra.Command{
		Use:   "kafka-connector",
		Short: "Install kafka-connector for OpenFaaS",
		Long: `Install OpenFaaS Pro kafka-connector for OpenFaaS so that you can invoke 
functions when messages are received on a given topic on a Kafka broker.`,
		Example:      `  arkade install kafka-connector`,
		SilenceUsage: true,
	}

	command.Flags().StringP("namespace", "n", "openfaas", "The namespace used for installation")
	command.Flags().Bool("update-repo", true, "Update the helm repo")
	command.Flags().StringP("topics", "t", "faas-request", "The topics for the connector to bind to")
	command.Flags().String("broker-hosts", "kafka:9092", "The server address or multiple addresses separated by a comma for the Kafka broker(s)")
	command.Flags().String("license-file", "", "The path to your license for OpenFaaS Pro")
	command.Flags().String("image", "", "The container image for the connector")

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

		brokerHostsVal, err := command.Flags().GetString("broker-hosts")
		if err != nil {
			return err
		}

		imageVal := ""
		if command.Flags().Changed("image") {
			imageVal, err = command.Flags().GetString("image")
			if err != nil {
				return nil
			}
		}

		overrides := map[string]string{
			"topics":      topicsVal,
			"brokerHosts": brokerHostsVal,
		}
		if len(imageVal) > 0 {
			overrides["image"] = imageVal
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kafkaConnectorAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("openfaas/kafka-connector").
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		// If license file is sent, then we assume to set the --pro flag and create the secret
		licenseFile, err := command.Flags().GetString("license-file")
		if err != nil {
			return err
		}

		if len(licenseFile) == 0 {
			return fmt.Errorf("--license-file is required for OpenFaaS Pro")
		}

		secretData := []types.SecretsData{
			{Type: types.FromFileSecret, Key: "license", Value: licenseFile},
		}

		proLicense := types.NewGenericSecret("openfaas-license", namespace, secretData)
		kafkaConnectorAppOptions.WithSecret(proLicense)

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
	"\n\n" + KafkaConnectorInfoMsg + "\n\n" + pkg.SupportMessageShort
