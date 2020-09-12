// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/alexellis/arkade/cmd/apps"
	"github.com/spf13/cobra"
)

func makeCommand(use, msg string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          use,
		Short:        fmt.Sprintf("Find info about a %s app", use),
		Long:         fmt.Sprintf("Find info about how to use the installed %s app", use),
		Example:      fmt.Sprintf(`  arkade info %s`, use),
		SilenceUsage: true,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Info for app: %s\n", use)
		fmt.Println(msg)

		return nil
	}

	return cmd
}

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

		return fmt.Errorf("no info available for app: %s", appName)
	}

	addCmd := func(cmd *cobra.Command, use, msg string) {
		cmd.AddCommand(makeCommand(use, msg))
	}

	addCmd(info, "openfaas", apps.OpenFaaSInfoMsg)
	addCmd(info, "nginx-ingress", apps.NginxIngressInfoMsg)
	addCmd(info, "cert-manager", apps.CertManagerInfoMsg)
	addCmd(info, "openfaas-ingress", apps.OpenfaasIngressInfoMsg)
	addCmd(info, "inlets-operator", apps.InletsOperatorInfoMsg)
	addCmd(info, "mongodb", apps.MongoDBInfoMsg)
	addCmd(info, "metrics-server", apps.MetricsInfoMsg)
	addCmd(info, "linkerd", apps.LinkerdInfoMsg)
	addCmd(info, "cron-connector", apps.CronConnectorInfoMsg)
	addCmd(info, "kafka-connector", apps.KafkaConnectorInfoMsg)
	addCmd(info, "kube-state-metrics", apps.KubeStateMetricsInfoMsg)
	addCmd(info, "minio", apps.MinioInfoMsg)
	addCmd(info, "postgresql", apps.PostgresqlInfoMsg)
	addCmd(info, "kubernetes-dashboard", apps.KubernetesDashboardInfoMsg)
	addCmd(info, "istio", apps.IstioInfoMsg)
	addCmd(info, "crossplane", apps.CrossplanInfoMsg)
	addCmd(info, "docker-registry-ingress", apps.RegistryIngressInfoMsg)
	addCmd(info, "traefik2", apps.Traefik2InfoMsg)
	addCmd(info, "tekton", apps.TektonInfoMsg)
	addCmd(info, "grafana", apps.GrafanaInfoMsg)
	addCmd(info, "argocd", apps.ArgoCDInfoMsg)
	addCmd(info, "portainer", apps.PortainerInfoMsg)
	addCmd(info, "jenkins", apps.JenkinsInfoMsg)
	addCmd(info, "loki", apps.LokiInfoMsg)
	addCmd(info, "nats-connector", apps.NATSConnectorInfoMsg)
	addCmd(info, "openfaas-loki", apps.LokiOFInfoMsg)
	addCmd(info, "redis", apps.RedisInfoMsg)
	addCmd(info, "kube-image-prefetch", apps.KubeImagePrefetchInfoMsg)
	addCmd(info, "registry-creds", apps.RegistryCredsOperatorInfoMsg)

	return info
}
