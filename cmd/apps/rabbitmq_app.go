// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallRabbitmq() *cobra.Command {
	var rabbitmq = &cobra.Command{
		Use:          "rabbitmq",
		Short:        "Install rabbitmq",
		Long:         "Install rabbitmq",
		Example:      "arkade install rabbitmq",
		SilenceUsage: true,
	}

	rabbitmq.Flags().StringP("namespace", "n", "rabbitmq", "The namespace to install rabbitmq")
	rabbitmq.Flags().Bool("update-repo", true, "Update the helm repo")
	rabbitmq.Flags().Bool("persistence", false, "Make rabbitmq persistent")
	rabbitmq.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	rabbitmq.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetBool("wait")
		if err != nil {
			return fmt.Errorf("error with --wait usage: %s", err)
		}

		_, err = command.Flags().GetBool("persistence")
		if err != nil {
			return fmt.Errorf("error with --persistence usage: %s", err)
		}

		_, err = command.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error with --namespace usage: %s", err)
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

	rabbitmq.RunE = func(command *cobra.Command, args []string) error {
		// Get all flags
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		wait, _ := command.Flags().GetBool("wait")
		persistence, _ := command.Flags().GetBool("persistence")
		namespace, _ := command.Flags().GetString("namespace")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		customFlags, _ := command.Flags().GetStringArray("set")
		overrides := map[string]string{}

		if persistence {
			overrides["persistence.enabled"] = "true"
			overrides["persistence.size"] = "2Gi"
		}

		// set custom flags
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		rabbitmqAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("bitnami/rabbitmq").
			WithHelmURL("https://charts.bitnami.com/bitnami").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err := apps.MakeInstallChart(rabbitmqAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(rabbitmqInstallMsg)

		return nil
	}

	return rabbitmq
}

const RabbitmqInfoMsg = `
# By default the "user" username is used
# To get your password run:

  export RABBITMQ_PASSWORD=$(kubectl get secret --namespace rabbitmq rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d)

# To Access the RabbitMQ Management interface:

  kubectl port-forward --namespace rabbitmq svc/rabbitmq 15672:15672

# To Access the RabbitMQ AMQP port:

  kubectl port-forward --namespace rabbitmq svc/rabbitmq 5672:5672

# Enable persistence:

  arkade install rabbitmq --persistence
`

var rabbitmqInstallMsg = `=======================================================================
=                      rabbitmq has been installed                     =
=======================================================================` +
	"\n\n" + RabbitmqInfoMsg + "\n\n" + pkg.SupportMessageShort
