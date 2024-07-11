// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/cmd/apps"
)

type ArkadeApp struct {
	Name        string
	Installer   func() *cobra.Command
	InfoMessage string
}

func MakeInstall() *cobra.Command {
	appList := GetApps()
	var command = &cobra.Command{
		Use:     "install",
		Short:   "Install Kubernetes apps from helm charts or YAML files",
		Aliases: []string{"i"},
		Long: `Install Kubernetes apps from helm charts or YAML files using the "install"
command.

You can also find the post-install message for each app with the "info"
command.`,
		Example: `  arkade install
  arkade install openfaas  --gateways=2
  arkade install inlets-operator --token-file $HOME/do-token`,
		SilenceUsage: true,
	}

	command.PersistentFlags().String("kubeconfig", "", "Local path for your kubeconfig file")
	command.PersistentFlags().Bool("wait", false, "If we should wait for the resource to be ready before returning (helm3 only, default false)")
	command.Flags().Bool("print-table", false, "print a table in markdown format")

	command.RunE = func(command *cobra.Command, args []string) error {

		printTable, _ := command.Flags().GetBool("print-table")

		if printTable {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Tool", "Description"})

			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetAutoWrapText(false)

			appSortedList := make([]string, 0, len(appList))

			for a := range appList {
				appSortedList = append(appSortedList, a)
			}
			sort.Strings(appSortedList)

			for _, k := range appSortedList {
				table.Append([]string{k, appList[k].Installer().Short})
			}

			table.Render()

			fmt.Printf("\nThere are %d apps that you can install on your cluster.\n", len(appList))
			return nil
		}

		if len(args) == 0 {
			fmt.Printf(
				`You can install %d apps to your Kubernetes cluster.

Run the following to see a list of all available apps:
  arkade install --help

To see options for a specific app before installing, run:

  arkade install APP --help
  arkade install openfaas --help
  arkade install grafana --help

To request a new app, raise a GitHub issue at:
  https://arkade.dev/
`, len(command.Commands()))
			return nil
		}

		name := args[0]
		var app *ArkadeApp
		if len(args) == 1 {
			for _, a := range appList {
				if a.Name == name {
					app = &a
					break
				}
			}
		}
		if app == nil {
			return fmt.Errorf("no such app: %s, run \"arkade install --help\" for a list of apps", name)
		}

		return nil
	}

	for _, app := range appList {
		command.AddCommand(app.Installer())
	}

	command.AddCommand(MakeInfo())

	return command
}

func GetApps() map[string]ArkadeApp {
	arkadeApps := map[string]ArkadeApp{}
	arkadeApps["mongodb"] = NewArkadeApp(apps.MakeInstallMongoDB, apps.MongoDBInfoMsg)
	arkadeApps["metrics-server"] = NewArkadeApp(apps.MakeInstallMetricsServer, apps.MetricsInfoMsg)
	arkadeApps["linkerd"] = NewArkadeApp(apps.MakeInstallLinkerd, apps.LinkerdInfoMsg)
	arkadeApps["cron-connector"] = NewArkadeApp(apps.MakeInstallCronConnector, apps.CronConnectorInfoMsg)
	arkadeApps["kafka-connector"] = NewArkadeApp(apps.MakeInstallKafkaConnector, apps.KafkaConnectorInfoMsg)
	arkadeApps["kube-state-metrics"] = NewArkadeApp(apps.MakeInstallKubeStateMetrics, apps.KubeStateMetricsInfoMsg)
	arkadeApps["kubernetes-dashboard"] = NewArkadeApp(apps.MakeInstallKubernetesDashboard, apps.KubernetesDashboardInfoMsg)
	arkadeApps["istio"] = NewArkadeApp(apps.MakeInstallIstio, apps.IstioInfoMsg)
	arkadeApps["crossplane"] = NewArkadeApp(apps.MakeInstallCrossplane, apps.CrossplaneInfoMsg)
	arkadeApps["docker-registry-ingress"] = NewArkadeApp(apps.MakeInstallRegistryIngress, apps.RegistryIngressInfoMsg)
	arkadeApps["postgresql"] = NewArkadeApp(apps.MakeInstallPostgresql, apps.PostgresqlInfoMsg)
	arkadeApps["minio"] = NewArkadeApp(apps.MakeInstallMinio, apps.MinioInfoMsg)
	arkadeApps["openfaas"] = NewArkadeApp(apps.MakeInstallOpenFaaS, apps.OpenFaaSInfoMsg)
	arkadeApps["ingress-nginx"] = NewArkadeApp(apps.MakeInstallNginx, apps.NginxIngressInfoMsg)
	arkadeApps["cert-manager"] = NewArkadeApp(apps.MakeInstallCertManager, apps.CertManagerInfoMsg)
	arkadeApps["openfaas-ingress"] = NewArkadeApp(apps.MakeInstallOpenFaaSIngress, apps.OpenfaasIngressInfoMsg)
	arkadeApps["openfaas-loki"] = NewArkadeApp(apps.MakeInstallOpenFaaSLoki, apps.LokiOFInfoMsg)
	arkadeApps["loki"] = NewArkadeApp(apps.MakeInstallLoki, apps.LokiInfoMsg)
	arkadeApps["redis"] = NewArkadeApp(apps.MakeInstallRedis, apps.RedisInfoMsg)
	arkadeApps["nats-connector"] = NewArkadeApp(apps.MakeInstallNATSConnector, apps.NATSConnectorInfoMsg)
	arkadeApps["jenkins"] = NewArkadeApp(apps.MakeInstallJenkins, apps.JenkinsInfoMsg)
	arkadeApps["portainer"] = NewArkadeApp(apps.MakeInstallPortainer, apps.PortainerInfoMsg)
	arkadeApps["argocd"] = NewArkadeApp(apps.MakeInstallArgoCD, apps.ArgoCDInfoMsg)
	arkadeApps["grafana"] = NewArkadeApp(apps.MakeInstallGrafana, apps.GrafanaInfoMsg)
	arkadeApps["tekton"] = NewArkadeApp(apps.MakeInstallTekton, apps.TektonInfoMsg)
	arkadeApps["traefik2"] = NewArkadeApp(apps.MakeInstallTraefik2, apps.Traefik2InfoMsg)
	arkadeApps["inlets-operator"] = NewArkadeApp(apps.MakeInstallInletsOperator, apps.InletsOperatorInfoMsg)
	arkadeApps["nfs-provisioner"] = NewArkadeApp(apps.MakeInstallNfsProvisioner, apps.NfsClientProvisioneriInfoMsg)
	arkadeApps["docker-registry"] = NewArkadeApp(apps.MakeInstallRegistry, apps.RegistryInfoMsg)
	arkadeApps["kube-image-prefetch"] = NewArkadeApp(apps.MakeInstallKubeImagePrefetch, apps.KubeImagePrefetchInfoMsg)
	arkadeApps["registry-creds"] = NewArkadeApp(apps.MakeInstallRegistryCredsOperator, apps.RegistryCredsOperatorInfoMsg)
	arkadeApps["gitea"] = NewArkadeApp(apps.MakeInstallGitea, apps.GiteaInfoMsg)
	arkadeApps["kong-ingress"] = NewArkadeApp(apps.MakeInstallKongIngress, apps.KongIngressInfoMsg)
	arkadeApps["sealed-secret"] = NewArkadeApp(apps.MakeInstallSealedSecrets, apps.SealedSecretsInfoMsg)
	arkadeApps["consul-connect"] = NewArkadeApp(apps.MakeInstallConsul, apps.ConsulInfoMsg)
	arkadeApps["sealed-secret"] = NewArkadeApp(apps.MakeInstallSealedSecrets, apps.SealedSecretsInfoMsg)
	arkadeApps["gitlab"] = NewArkadeApp(apps.MakeInstallGitLab, apps.GitlabInfoMsg)
	arkadeApps["opa-gatekeeper"] = NewArkadeApp(apps.MakeInstallOPAGateKeeper, apps.OPAGatekeeperInfoMsg)
	arkadeApps["mqtt-connector"] = NewArkadeApp(apps.MakeInstallMQTTConnector, apps.MQTTConnectorInfoMsg)
	arkadeApps["falco"] = NewArkadeApp(apps.MakeInstallFalco, apps.FalcoInfoMsg)
	arkadeApps["influxdb"] = NewArkadeApp(apps.MakeInstallinfluxdb, apps.InfluxdbInfoMsg)
	arkadeApps["kafka"] = NewArkadeApp(apps.MakeInstallConfluentPlatformKafka, apps.KafkaInfoMsg)
	arkadeApps["kyverno"] = NewArkadeApp(apps.MakeInstallKyverno, apps.KyvernoInfoMsg)
	arkadeApps["rabbitmq"] = NewArkadeApp(apps.MakeInstallRabbitmq, apps.RabbitmqInfoMsg)
	arkadeApps["cassandra"] = NewArkadeApp(apps.MakeInstallCassandra, apps.CassandraInfoMsg)
	arkadeApps["metallb-arp"] = NewArkadeApp(apps.MakeInstallMetalLB, apps.MetalLBInfoMsg)
	arkadeApps["cockroachdb"] = NewArkadeApp(apps.MakeInstallCockroachdb, apps.CockroachdbInfoMsg)
	arkadeApps["prometheus"] = NewArkadeApp(apps.MakeInstallPrometheus, apps.PrometheusInfoMsg)
	arkadeApps["waypoint"] = NewArkadeApp(apps.MakeInstallWaypoint, apps.WaypointInfoMsg)
	arkadeApps["kuma"] = NewArkadeApp(apps.MakeInstallKuma, apps.KumaInfoMsg)
	arkadeApps["qemu-static"] = NewArkadeApp(apps.MakeInstallQemuStatic, apps.QemuStaticInfoMsg)
	arkadeApps["vault"] = NewArkadeApp(apps.MakeInstallVault, apps.VaultInfoMsg)

	// Special "chart" app - let a user deploy any helm chart
	arkadeApps["chart"] = NewArkadeApp(apps.MakeInstallChart, "")
	return arkadeApps
}

func NewArkadeApp(cmd func() *cobra.Command, msg string) ArkadeApp {
	return ArkadeApp{
		Installer:   cmd,
		InfoMessage: msg,
	}
}
