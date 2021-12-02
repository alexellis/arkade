// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallCockroachdb() *cobra.Command {
	var cockroachdb = &cobra.Command{
		Use:          "cockroachdb",
		Short:        "Install CockroachDB",
		Long:         "Install CockroachDB",
		Example:      ` arkade app install cockroachdb`,
		SilenceUsage: true,
	}
	cockroachdb.Flags().String("namespace", "default", "Namespace for the app")
	cockroachdb.Flags().Bool("persistence", false, "Use a 100Gi Persistent Volume to store data")

	cockroachdb.Flags().Bool("single-node", false, "Run CockroachDB instances in standalone mode with replication disabled, so the StatefulSet does NOT FORM A CLUSTER")
	cockroachdb.Flags().Int64("replicas", 1, "Statefulset replica count")

	cockroachdb.Flags().Bool("tls", false, "Whether to run securely using TLS certificates")

	cockroachdb.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set tls.enabled=false)")

	cockroachdb.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, _ := command.Flags().GetString("namespace")
		persistence, _ := command.Flags().GetBool("persistence")
		singleNode, _ := command.Flags().GetBool("single-node")
		enableTls, _ := command.Flags().GetBool("enable-tls")
		replicas, _ := command.Flags().GetInt64("replicas")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if suffix := getValuesSuffix(arch); suffix == "-armhf" {
			return fmt.Errorf(`CockroachDB is currently not supported on armhf architectures`)
		}

		overrides := map[string]string{}

		if singleNode {
			overrides["conf.single-node"] = "true"
		}
		if enableTls {
			overrides["tls.enabled"] = "false"
		}
		if persistence {
			overrides["storage.persistentVolume.enabled"] = "true"
		}
		overrides["statefulset.replicas"] = fmt.Sprintf("%d", replicas)

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		cockroachdbOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("cockroachdb/cockroachdb").
			WithHelmURL("https://charts.cockroachdb.com/").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(cockroachdbOptions)
		if err != nil {
			return err
		}

		println(cockroachdbInstallMsg)
		return nil
	}

	return cockroachdb
}

const CockroachdbInfoMsg = `# Get started at: https://www.cockroachlabs.com/docs/stable

# If you used tls check that the secrets were created on the cluster

kubectl describe secrets/crdb-cockroachdb-ca-secret
kubectl describe secrets/crdb-cockroachdb-client-secret
kubectl describe secrets/crdb-cockroachdb-node-secret

# If you used persistence confirm that the volumes are created

kubectl get pv

# Visit the CockroachDB dashboard:

kubectl port-forward service/cockroachdb 8080:8080`

const cockroachdbInstallMsg = `=======================================================================
=                  CockroachDB has been installed                     =
=======================================================================
 ` + pkg.ThanksForUsing + CockroachdbInfoMsg
