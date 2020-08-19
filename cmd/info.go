// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/alexellis/arkade/cmd/apps"
	"github.com/spf13/cobra"
)

func MakeInfo() *cobra.Command {

	info := &cobra.Command{
		Use:   "info",
		Short: "Find info about a Kubernetes app",
		Long:  "Find info about how to use the installed Kubernetes app",
		Example: `  arkade info [APP]
arkade info openfaas
arkade info inlets-operator
arkade info mongodb
arkade info
arkade info --help`,
		SilenceUsage: true,
	}

	info.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("Run arkade info APP_NAME for more")
			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf("you can only get info about exactly one installed app")
		}

		var (
			appName = args[0]
			msg     string
		)

		switch appName {
		case "openfaas":
			msg = apps.OpenFaaSInfoMsg
		case "nginx-ingress":
			msg = apps.NginxIngressInfoMsg
		case "cert-manager":
			msg = apps.CertManagerInfoMsg
		case "openfaas-ingress":
			msg = apps.OpenfaasIngressInfoMsg
		case "inlets-operator":
			msg = apps.InletsOperatorInfoMsg
		case "mongodb":
			msg = apps.MongoDBInfoMsg
		case "metrics-server":
			msg = apps.MetricsInfoMsg
		case "linkerd":
			msg = apps.LinkerdInfoMsg
		case "cron-connector":
			msg = apps.CronConnectorInfoMsg
		case "kafka-connector":
			msg = apps.KafkaConnectorInfoMsg
		case "kube-state-metrics":
			msg = apps.KubeStateMetricsInfoMsg
		case "minio":
			msg = apps.MinioInfoMsg
		case "postgresql":
			msg = apps.PostgresqlInfoMsg
		case "kubernetes-dashboard":
			msg = apps.KubernetesDashboardInfoMsg
		case "istio":
			msg = apps.IstioInfoMsg
		case "crossplane":
			msg = apps.CrossplanInfoMsg
		case "docker-registry-ingress":
			msg = apps.RegistryIngressInfoMsg
		case "traefik2":
			msg = apps.Traefik2InfoMsg
		case "tekton":
			msg = apps.TektonInfoMsg
		case "grafana":
			msg = apps.GrafanaInfoMsg
		case "argocd":
			msg = apps.ArgoCDInfoMsg
		case "portainer":
			msg = apps.PortainerInfoMsg
		case "jenkins":
			msg = apps.JenkinsInfoMsg
		case "loki":
			msg = apps.LokiInfoMsg
		case "nats-connector":
			msg = apps.NATSConnectorInfoMsg
		case "openfaas-loki":
			msg = apps.LokiOFInfoMsg
		case "redis":
			msg = apps.RedisInfoMsg
		case "kube-image-prefetch":
			msg = apps.KubeImagePrefetchInfoMsg
		case "registry-creds":
			msg = apps.RegistryCredsOperatorInfoMsg
		default:
			return fmt.Errorf("no info available for app: %s", appName)
		}

		fmt.Printf("Info for app: %s\n", appName)
		fmt.Println(msg)

		return nil
	}

	return info
}
