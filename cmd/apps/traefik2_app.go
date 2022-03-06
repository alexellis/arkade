// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallTraefik2() *cobra.Command {
	var traefik2 = &cobra.Command{
		Use:          "traefik2",
		Short:        "Install traefik2",
		Long:         "Install traefik2",
		Example:      `  arkade app install traefik2`,
		SilenceUsage: true,
	}

	traefik2.Flags().StringP("namespace", "n", "kube-system", "The namespace used for installation")
	traefik2.Flags().Bool("update-repo", true, "Update the helm repo")
	traefik2.Flags().Bool("load-balancer", true, "Use a load-balancer for the IngressController")
	traefik2.Flags().Bool("dashboard", false, "Expose dashboard if you want access to dashboard from the browser")
	traefik2.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set key=value)")
	traefik2.Flags().Bool("wait", false, "Wait for the chart to be installed")
	traefik2.Flags().Bool("ingress-provider", true, "Add Traefik's ingressprovider along with the CRD provider")

	traefik2.RunE = func(command *cobra.Command, args []string) error {

		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		updateRepo, _ := traefik2.Flags().GetBool("update-repo")
		namespace, _ := traefik2.Flags().GetString("namespace")
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q\n", clientOS)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		err = helm.AddHelmRepo("traefik", "https://helm.traefik.io/traefik", updateRepo)
		if err != nil {
			return fmt.Errorf("Unable to add repo %s", err)
		}

		err = helm.FetchChart("traefik/traefik", "")
		if err != nil {
			return fmt.Errorf("Unable fetch chart: %s", err)
		}

		overrides := map[string]string{}
		lb, _ := command.Flags().GetBool("load-balancer")
		dashboard, _ := command.Flags().GetBool("dashboard")
		wait, _ := command.Flags().GetBool("wait")
		ingressProvider, _ := command.Flags().GetBool("ingress-provider")

		svc := "NodePort"
		if lb {
			svc = "LoadBalancer"
		}
		overrides["service.type"] = svc

		overrides["additional.checkNewVersion"] = "false"
		overrides["additional.sendAnonymousUsage"] = "false"

		if dashboard {
			overrides["dashboard.ingressRoute"] = "true"
		}

		if ingressProvider {
			overrides["additionalArguments"] = `{--providers.kubernetesingress}`
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		err = helm.Helm3Upgrade("traefik/traefik", namespace,
			"values.yaml",
			"",
			overrides,
			wait)

		if err != nil {
			return err
		}

		fmt.Println(traefikInstallMsg)
		return nil
	}

	return traefik2
}

const Traefik2InfoMsg = `# Get started at: https://doc.traefik.io/traefik/

# Install with an optional dashboard

arkade install traefik2 --dashboard

# Find your LoadBalancer IP:

kubectl get svc -n kube-system traefik`

const traefikInstallMsg = `=======================================================================
=                  traefik2 has been installed                        =
=======================================================================
 ` + pkg.ThanksForUsing + Traefik2InfoMsg
