// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallCassandra() *cobra.Command {
	var cassandra = &cobra.Command{
		Use:          "cassandra",
		Short:        "Install cassandra",
		Long:         "Install cassandra",
		Example:      "arkade install cassandra",
		SilenceUsage: true,
	}

	cassandra.Flags().StringP("namespace", "n", "cassandra", "The namespace to install cassandra")
	cassandra.Flags().Bool("update-repo", true, "Update the helm repo")
	cassandra.Flags().Bool("persistence", false, "Make cassandra persistent")
	cassandra.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	cassandra.PreRunE = func(command *cobra.Command, args []string) error {
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

		arch := k8s.GetNodeArchitecture()
		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		return nil
	}

	cassandra.RunE = func(command *cobra.Command, args []string) error {
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
		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		cassandraAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("bitnami/cassandra").
			WithHelmURL("https://charts.bitnami.com/bitnami").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err := apps.MakeInstallChart(cassandraAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(cassandraInstallMsg)

		return nil
	}

	return cassandra
}

const CassandraInfoMsg = `
# By default the "cassandra" username is used
# To get your password run:

	export CASSANDRA_PASSWORD=$(kubectl get secret --namespace "cassandra" cassandra -o jsonpath="{.data.cassandra-password}" | base64 --decode)

# Run a Cassandra pod that you can use as a client:

	kubectl run --namespace cassandra cassandra-client --rm --tty -i --restart='Never' \
	--env CASSANDRA_PASSWORD=$CASSANDRA_PASSWORD \
	\
	--image docker.io/bitnami/cassandra:3.11.10-debian-10-r149 -- bash

# Connect using the cqlsh client:

	cqlsh -u cassandra -p $CASSANDRA_PASSWORD cassandra

# To connect to your database from outside the cluster execute the following commands:

	kubectl port-forward --namespace cassandra svc/cassandra 9042:9042 &
	cqlsh -u cassandra -p $CASSANDRA_PASSWORD 127.0.0.1 9042

# Enable persistence:
	arkade install cassandra --persistence
`

var cassandraInstallMsg = `=======================================================================
=                      cassandra has been installed                     =
=======================================================================` +
	"\n\n" + CassandraInfoMsg + "\n\n" + pkg.SupportMessageShort
