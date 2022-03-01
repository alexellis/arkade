// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"
	"strconv"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallMQTTConnector() *cobra.Command {
	var command = &cobra.Command{
		Use:   "mqtt-connector",
		Short: "Install mqtt-connector for OpenFaaS",
		Long: `Install mqtt-connector for OpenFaaS so that you can invoke functions when 
messages are received on a given topic on an MQTT broker.`,
		Example:      `  arkade install mqtt-connector`,
		SilenceUsage: true,
	}

	command.Flags().StringP("namespace", "n", "openfaas", "The namespace used for installation")
	command.Flags().Bool("update-repo", true, "Update the helm repo")
	command.Flags().StringP("topics", "t", "", "The topics for the connector to bind to, currently supports one topic")
	command.Flags().String("broker-host", "tcp://test.mosquitto.org:1883", "The host for the MQTT broker")
	command.Flags().String("client-id", "mqtt-connector-1", "The client ID for the MQTT broker")
	command.Flags().Bool("async", false, "Invoke functions asynchronously from events")
	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		updateRepo, _ := command.Flags().GetBool("update-repo")
		async, err := command.Flags().GetBool("async")
		if err != nil {
			return err
		}

		namespace, _ := command.Flags().GetString("namespace")
		clientID, _ := command.Flags().GetString("client-id")
		topicsVal, err := command.Flags().GetString("topics")
		if err != nil {
			return err
		}

		brokerHostVal, err := command.Flags().GetString("broker-host")
		if err != nil {
			return err
		}

		overrides := map[string]string{
			"topic":       topicsVal,
			"broker":      brokerHostVal,
			"clientID":    clientID,
			"asyncInvoke": strconv.FormatBool(async),
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		if len(topicsVal) == 0 {
			return fmt.Errorf("--topics is required")
		}

		if len(brokerHostVal) == 0 {
			return fmt.Errorf("--broker-host is required")
		}

		mqttConnectorAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("openfaas/mqtt-connector").
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(mqttConnectorAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(mqttConnectorInstallMsg)

		return nil
	}

	return command
}

const MQTTConnectorInfoMsg = `# View the connector's logs:

kubectl logs deploy/mqtt-connector -n openfaas -f

# Find out more on the project homepage:

# https://github.com/openfaas/mqtt-connector/`

const mqttConnectorInstallMsg = `=======================================================================
= mqtt-connector has been installed.                                   =
=======================================================================` +
	"\n\n" + MQTTConnectorInfoMsg + "\n\n" + pkg.SupportMessageShort
