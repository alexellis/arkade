// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/cmd/apps"
	"github.com/alexellis/arkade/pkg/get"
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
			renderTable(os.Stdout, appList)
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

		if len(args) == 1 {
			name := args[0]
			app := findAppByName(name, appList)
			if app == nil {
				return errors.New(checkForTool(name, get.MakeTools()))
			}
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
	arkadeApps["argocd"] = NewArkadeApp(apps.MakeInstallArgoCD, apps.ArgoCDInfoMsg)
	arkadeApps["cassandra"] = NewArkadeApp(apps.MakeInstallCassandra, apps.CassandraInfoMsg)
	arkadeApps["cert-manager"] = NewArkadeApp(apps.MakeInstallCertManager, apps.CertManagerInfoMsg)
	arkadeApps["cockroachdb"] = NewArkadeApp(apps.MakeInstallCockroachdb, apps.CockroachdbInfoMsg)
	arkadeApps["consul-connect"] = NewArkadeApp(apps.MakeInstallConsul, apps.ConsulInfoMsg)
	arkadeApps["cron-connector"] = NewArkadeApp(apps.MakeInstallCronConnector, apps.CronConnectorInfoMsg)
	arkadeApps["crossplane"] = NewArkadeApp(apps.MakeInstallCrossplane, apps.CrossplaneInfoMsg)
	arkadeApps["docker-registry"] = NewArkadeApp(apps.MakeInstallRegistry, apps.RegistryInfoMsg)
	arkadeApps["docker-registry-ingress"] = NewArkadeApp(apps.MakeInstallRegistryIngress, apps.RegistryIngressInfoMsg)
	arkadeApps["falco"] = NewArkadeApp(apps.MakeInstallFalco, apps.FalcoInfoMsg)
	arkadeApps["gitea"] = NewArkadeApp(apps.MakeInstallGitea, apps.GiteaInfoMsg)
	arkadeApps["gitlab"] = NewArkadeApp(apps.MakeInstallGitLab, apps.GitlabInfoMsg)
	arkadeApps["grafana"] = NewArkadeApp(apps.MakeInstallGrafana, apps.GrafanaInfoMsg)
	arkadeApps["influxdb"] = NewArkadeApp(apps.MakeInstallinfluxdb, apps.InfluxdbInfoMsg)
	arkadeApps["ingress-nginx"] = NewArkadeApp(apps.MakeInstallNginx, apps.NginxIngressInfoMsg)
	arkadeApps["inlets-operator"] = NewArkadeApp(apps.MakeInstallInletsOperator, apps.InletsOperatorInfoMsg)
	arkadeApps["istio"] = NewArkadeApp(apps.MakeInstallIstio, apps.IstioInfoMsg)
	arkadeApps["jenkins"] = NewArkadeApp(apps.MakeInstallJenkins, apps.JenkinsInfoMsg)
	arkadeApps["kafka"] = NewArkadeApp(apps.MakeInstallConfluentPlatformKafka, apps.KafkaInfoMsg)
	arkadeApps["kafka-connector"] = NewArkadeApp(apps.MakeInstallKafkaConnector, apps.KafkaConnectorInfoMsg)
	arkadeApps["kong-ingress"] = NewArkadeApp(apps.MakeInstallKongIngress, apps.KongIngressInfoMsg)
	arkadeApps["kube-image-prefetch"] = NewArkadeApp(apps.MakeInstallKubeImagePrefetch, apps.KubeImagePrefetchInfoMsg)
	arkadeApps["kube-state-metrics"] = NewArkadeApp(apps.MakeInstallKubeStateMetrics, apps.KubeStateMetricsInfoMsg)
	arkadeApps["kubernetes-dashboard"] = NewArkadeApp(apps.MakeInstallKubernetesDashboard, apps.KubernetesDashboardInfoMsg)
	arkadeApps["kuma"] = NewArkadeApp(apps.MakeInstallKuma, apps.KumaInfoMsg)
	arkadeApps["kyverno"] = NewArkadeApp(apps.MakeInstallKyverno, apps.KyvernoInfoMsg)
	arkadeApps["linkerd"] = NewArkadeApp(apps.MakeInstallLinkerd, apps.LinkerdInfoMsg)
	arkadeApps["loki"] = NewArkadeApp(apps.MakeInstallLoki, apps.LokiInfoMsg)
	arkadeApps["metallb-arp"] = NewArkadeApp(apps.MakeInstallMetalLB, apps.MetalLBInfoMsg)
	arkadeApps["metrics-server"] = NewArkadeApp(apps.MakeInstallMetricsServer, apps.MetricsInfoMsg)
	arkadeApps["minio"] = NewArkadeApp(apps.MakeInstallMinio, apps.MinioInfoMsg)
	arkadeApps["mongodb"] = NewArkadeApp(apps.MakeInstallMongoDB, apps.MongoDBInfoMsg)
	arkadeApps["mqtt-connector"] = NewArkadeApp(apps.MakeInstallMQTTConnector, apps.MQTTConnectorInfoMsg)
	arkadeApps["nats-connector"] = NewArkadeApp(apps.MakeInstallNATSConnector, apps.NATSConnectorInfoMsg)
	arkadeApps["nfs-provisioner"] = NewArkadeApp(apps.MakeInstallNfsProvisioner, apps.NfsClientProvisioneriInfoMsg)
	arkadeApps["opa-gatekeeper"] = NewArkadeApp(apps.MakeInstallOPAGateKeeper, apps.OPAGatekeeperInfoMsg)
	arkadeApps["openfaas"] = NewArkadeApp(apps.MakeInstallOpenFaaS, apps.OpenFaaSInfoMsg)
	arkadeApps["openfaas-ingress"] = NewArkadeApp(apps.MakeInstallOpenFaaSIngress, apps.OpenfaasIngressInfoMsg)
	arkadeApps["openfaas-loki"] = NewArkadeApp(apps.MakeInstallOpenFaaSLoki, apps.LokiOFInfoMsg)
	arkadeApps["portainer"] = NewArkadeApp(apps.MakeInstallPortainer, apps.PortainerInfoMsg)
	arkadeApps["postgresql"] = NewArkadeApp(apps.MakeInstallPostgresql, apps.PostgresqlInfoMsg)
	arkadeApps["prometheus"] = NewArkadeApp(apps.MakeInstallPrometheus, apps.PrometheusInfoMsg)
	arkadeApps["qemu-static"] = NewArkadeApp(apps.MakeInstallQemuStatic, apps.QemuStaticInfoMsg)
	arkadeApps["rabbitmq"] = NewArkadeApp(apps.MakeInstallRabbitmq, apps.RabbitmqInfoMsg)
	arkadeApps["redis"] = NewArkadeApp(apps.MakeInstallRedis, apps.RedisInfoMsg)
	arkadeApps["registry-creds"] = NewArkadeApp(apps.MakeInstallRegistryCredsOperator, apps.RegistryCredsOperatorInfoMsg)
	arkadeApps["sealed-secret"] = NewArkadeApp(apps.MakeInstallSealedSecrets, apps.SealedSecretsInfoMsg)
	arkadeApps["tekton"] = NewArkadeApp(apps.MakeInstallTekton, apps.TektonInfoMsg)
	arkadeApps["traefik2"] = NewArkadeApp(apps.MakeInstallTraefik2, apps.Traefik2InfoMsg)
	arkadeApps["vault"] = NewArkadeApp(apps.MakeInstallVault, apps.VaultInfoMsg)
	arkadeApps["waypoint"] = NewArkadeApp(apps.MakeInstallWaypoint, apps.WaypointInfoMsg)

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

func checkForTool(appName string, tools []get.Tool) string {

	for _, tool := range tools {
		if strings.EqualFold(tool.Name, appName) {
			return fmt.Sprintf("no such app. %s is available as a tool, run \"arkade get %s\" to get it", appName, appName)
		}
	}
	return fmt.Sprintf("no such app: %s, run \"arkade install --help\" for a list of apps", appName)
}

func renderTable(w io.Writer, appMap map[string]ArkadeApp) {

	symbols := tw.NewSymbolCustom("Lines").
		WithRow("-").
		WithColumn("|").
		WithCenter("|").
		WithMidLeft("|").
		WithMidRight("|")

	outline := tw.Border{
		Left:   tw.On,
		Right:  tw.On,
		Top:    tw.Off,
		Bottom: tw.Off,
	}

	table := tablewriter.NewTable(w,
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders: outline,
				Symbols: symbols,
			})),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoWrap: tw.WrapNone},
			},
		}),
	)
	table.Header([]string{"Tool", "Description"})
	appCount := len(appMap)

	appSortedList := make([]string, 0, appCount)

	for a := range appMap {
		appSortedList = append(appSortedList, a)
	}
	sort.Strings(appSortedList)

	for _, k := range appSortedList {
		table.Append([]string{k, appMap[k].Installer().Short})
	}

	table.Render()
	fmt.Fprintf(w, "\nThere are %d apps that you can install on your cluster.\n", appCount)
}

func findAppByName(name string, apps map[string]ArkadeApp) *ArkadeApp {
	if app, exists := apps[name]; exists {
		return &app
	}
	return nil
}
