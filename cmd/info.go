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
			fmt.Println("You can get info about: openfaas, inlets-operator, mongodb")
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
		default:
			return fmt.Errorf("no info available for app: %s", appName)
		}

		return nil
	}

	return info
}
