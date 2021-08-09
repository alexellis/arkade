package apps

import (
	"fmt"
	"strconv"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallConfluentPlatformKafka() *cobra.Command {
	kafka := &cobra.Command{
		Use:   "kafka",
		Short: "Install Confluent Platform Kafka",
		Long: `This will install Kafka provided by the Confluent Platform by using the following official Helm chart:
        https://github.com/confluentinc/cp-helm-charts`,
		Example:      "arkade install kafka",
		SilenceUsage: true,
	}

	kafka.Flags().Bool("zookeeper", true, "enable Zookeeper")
	kafka.Flags().Int("zookeeper-server-count", 1, "server count of Zookeeper")

	kafka.Flags().Bool("kafka", true, "enable Kafka")
	kafka.Flags().Int("kafka-broker-count", 1, "broker count of Kafka")

	kafka.Flags().Bool("schema-registry", false, "enable Schema Registry")
	kafka.Flags().Bool("kafka-rest", false, "enable Kafka Rest")

	kafka.Flags().Bool("kafka-connect", false, "enable Kafka Connect")

	kafka.Flags().Bool("ksql-server", false, "enable KSQL Server")

	kafka.Flags().Bool("control-center", false, "enable KSQL Server")

	kafka.Flags().Bool("update-repo", true, "Update the helm repo")

	kafka.RunE = func(command *cobra.Command, args []string) error {
		appOpts := types.DefaultInstallOptions()

		wait, err := command.Flags().GetBool("wait")
		if err != nil {
			return err
		}

		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")

		updateRepo, err := command.Flags().GetBool("update-repo")
		if err != nil {
			return err
		}
		appOpts.WithHelmUpdateRepo(updateRepo)

		overrides := map[string]string{}

		enableZookeeper, err := command.Flags().GetBool("zookeeper")
		if err != nil {
			return err
		}
		overrides["cp-zookeeper.enabled"] = strconv.FormatBool(enableZookeeper)

		zookeeperServerCount, err := command.Flags().GetInt("zookeeper-server-count")
		if err != nil {
			return err
		}
		overrides["cp-zookeeper.servers"] = fmt.Sprintf("%d", zookeeperServerCount)

		enableKafka, err := command.Flags().GetBool("kafka")
		if err != nil {
			return err
		}
		overrides["cp-kafka.enabled"] = strconv.FormatBool(enableKafka)

		kafkaBrokerCount, err := command.Flags().GetInt("kafka-broker-count")
		if err != nil {
			return err
		}
		overrides["cp-kafka.brokers"] = fmt.Sprintf("%d", kafkaBrokerCount)

		enableSchemaRegistry, err := command.Flags().GetBool("schema-registry")
		if err != nil {
			return err
		}
		overrides["cp-schema-registry.enabled"] = strconv.FormatBool(enableSchemaRegistry)

		enableKafkaRest, err := command.Flags().GetBool("kafka-rest")
		if err != nil {
			return err
		}
		overrides["cp-kafka-rest.enabled"] = strconv.FormatBool(enableKafkaRest)

		enableKafkaConnect, err := command.Flags().GetBool("kafka-connect")
		if err != nil {
			return err
		}
		overrides["cp-kafka-connect.enabled"] = strconv.FormatBool(enableKafkaConnect)

		enableKSQLServer, err := command.Flags().GetBool("ksql-server")
		if err != nil {
			return err
		}
		overrides["cp-ksql-server.enabled"] = strconv.FormatBool(enableKSQLServer)

		enableControlCenter, err := command.Flags().GetBool("control-center")
		if err != nil {
			return err
		}
		overrides["cp-control-center.enabled"] = strconv.FormatBool(enableControlCenter)

		appOpts.
			WithKubeconfigPath(kubeConfigPath).
			WithOverrides(overrides).
			WithValuesFile("values.yaml").
			WithHelmURL("https://confluentinc.github.io/cp-helm-charts/").
			WithHelmRepo("confluentinc/cp-helm-charts").
			WithNamespace(namespace).
			WithInstallNamespace(false).
			WithWait(wait)

		if _, err := apps.MakeInstallChart(appOpts); err != nil {
			return err
		}

		fmt.Println(kafkaPostInstallMsg)

		return nil
	}

	return kafka
}

const KafkaInfoMsg = `You can visit the official Helm Chart repository to get more detail about the installation:
https://github.com/confluentinc/cp-helm-charts
`

const kafkaPostInstallMsg = `=======================================================================
= Kafka has been installed.                                        =
=======================================================================` +
	"\n\n" + KafkaInfoMsg + "\n\n" + pkg.ThanksForUsing
