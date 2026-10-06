// Copyright (c) arkade author(s) 2022. All rights reserved.
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

func MakeInstallConfluentPlatformKafka() *cobra.Command {
	kafka := &cobra.Command{
		Use:   "kafka",
		Short: "Install Kafka",
		Long: `This will install Kafka using the Bitnami chart, distributed
as an OCI image: oci://registry-1.docker.io/bitnamicharts/kafka`,
		Example:      "arkade install kafka",
		SilenceUsage: true,
	}

	kafka.Flags().Int("replicas", 1, "number of Kafka brokers")
	kafka.Flags().Int("controller-replicas", 1, "number of Kafka controllers (KRaft mode)")
	kafka.Flags().String("storage-class", "", "override the storage class for Kafka data")
	kafka.Flags().String("heap", "1g", "JVM heap size for the Kafka broker, e.g. 1g or 2g")
	kafka.Flags().Bool("kafka", true, "enable Kafka")
	kafka.Flags().Bool("update-repo", true, "Update the helm repo")

	kafka.RunE = func(command *cobra.Command, args []string) error {
		appOpts := types.DefaultInstallOptions()

		wait, err := command.Flags().GetBool("wait")
		if err != nil {
			return err
		}

		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")

		if _, err := command.Flags().GetBool("update-repo"); err != nil {
			return err
		}

		overrides := map[string]string{}

		enableKafka, err := command.Flags().GetBool("kafka")
		if err != nil {
			return err
		}
		overrides["kafka.enabled"] = fmt.Sprintf("%v", enableKafka)

		replicas, err := command.Flags().GetInt("replicas")
		if err != nil {
			return err
		}
		overrides["controller.replicaCount"] = fmt.Sprintf("%d", replicas)

		controllerReplicas, err := command.Flags().GetInt("controller-replicas")
		if err != nil {
			return err
		}
		overrides["broker.replicaCount"] = fmt.Sprintf("%d", controllerReplicas)

		heap, err := command.Flags().GetString("heap")
		if err != nil {
			return err
		}
		if len(heap) > 0 {
			overrides["heapOpts"] = fmt.Sprintf("-Xmx%s -Xms%s", heap, heap)
		}

		storageClass, err := command.Flags().GetString("storage-class")
		if err != nil {
			return err
		}
		if len(storageClass) > 0 {
			overrides["persistence.storageClass"] = storageClass
			overrides["logPersistence.storageClass"] = storageClass
		}

		customFlags, _ := command.Flags().GetStringArray("set")
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		// Bitnami moved its free container images to the frozen
		// "bitnamilegacy" registry in Aug 2025, the chart's default
		// bitnami/kafka tags are gone. Point at the legacy images so the
		// install works out of the box; users can override via --set.
		overrides["global.imageRegistry"] = ""
		overrides["image.registry"] = "docker.io"
		overrides["image.repository"] = "bitnamilegacy/kafka"
		overrides["controller.image.registry"] = "docker.io"
		overrides["controller.image.repository"] = "bitnamilegacy/kafka"
		overrides["controller-els.image.registry"] = "docker.io"
		overrides["controller-els.image.repository"] = "bitnamilegacy/kafka"

		appOpts.
			WithKubeconfigPath(kubeConfigPath).
			WithOverrides(overrides).
			WithHelmURL("oci://registry-1.docker.io/bitnamicharts/kafka").
			WithHelmRepo("oci://registry-1.docker.io/bitnamicharts/kafka").
			WithNamespace(namespace).
			WithInstallNamespace(false).
			WithWait(wait)

		// The default options includes the `values.yaml` file but this is
		// already implied when using the OCI chart.
		appOpts.Helm.ValuesFiles = []string{}

		if _, err := apps.MakeInstallChart(appOpts); err != nil {
			return err
		}

		fmt.Println(kafkaPostInstallMsg)

		return nil
	}

	return kafka
}

const KafkaInfoMsg = `You can visit the official Helm Chart to get more detail about the installation:
https://artifacthub.io/packages/helm/bitnami/kafka
`

var kafkaPostInstallMsg = `=======================================================================
= Kafka has been installed.                                        =
=======================================================================` +
	"\n\n" + KafkaInfoMsg + "\n\n" + pkg.SupportMessageShort
