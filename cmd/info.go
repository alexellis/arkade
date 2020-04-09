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

		appName := args[0]

		switch appName {
		case "openfaas":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.OpenFaaSInfoMsg)
		case "nginx-ingress":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.NginxIngressInfoMsg)
		case "cert-manager":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.CertManagerInfoMsg)
		case "openfaas-ingress":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.OpenfaasIngressInfoMsg)
		case "inlets-operator":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.InletsOperatorInfoMsg)
		case "mongodb":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.MongoDBInfoMsg)
		case "metrics-server":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.MetricsInfoMsg)
		case "linkerd":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.LinkerdInfoMsg)
		case "cron-connector":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.CronConnectorInfoMsg)
		case "kafka-connector":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.KafkaConnectorInfoMsg)
		case "kube-state-metrics":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.KubeStateMetricsInfoMsg)
		case "minio":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.MinioInfoMsg)
		case "postgresql":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.PostgresqlInfoMsg)
		case "kubernetes-dashboard":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.KubernetesDashboardInfoMsg)
		case "istio":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.IstioInfoMsg)
		case "crossplane":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.CrossplanInfoMsg)
		case "docker-registry-ingress":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.RegistryIngressInfoMsg)
		case "traefik2":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.Traefik2InfoMsg)
		case "tekton":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.TektonInfoMsg)
		case "grafana":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.GrafanaInfoMsg)
		case "argocd":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.ArgoCDInfoMsg)
		case "portainer":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.PortainerInfoMsg)
		case "jenkins":
			fmt.Printf("Info for app: %s\n", appName)
			fmt.Println(apps.JenkinsInfoMsg)
		default:
			return fmt.Errorf("no info available for app: %s", appName)
		}

		return nil
	}

	return info
}
