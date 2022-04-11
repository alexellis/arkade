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

func MakeInstallPortainer() *cobra.Command {
	var portainer = &cobra.Command{
		Use:          "portainer",
		Short:        "Install portainer to visualise and manage containers",
		Long:         `Install portainer to visualise and manage containers, now in beta for Kubernetes.`,
		Example:      `  arkade install portainer`,
		SilenceUsage: true,
	}

	portainer.Flags().String("namespace", "default", "Namespace for the app")
	portainer.Flags().Bool("persistence", false, "Use a 10Gi Persistent Volume to store data")
	portainer.Flags().String("service-type", "ClusterIP", "Service Type for the main Portainer Service; ClusterIP, NodePort or LoadBalancer")

	portainer.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set tls.enabled=false)")

	portainer.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, err := command.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		persistence, err := command.Flags().GetBool("persistence")
		if err != nil {
			return err
		}

		serviceType, err := command.Flags().GetString("service-type")
		if err != nil {
			return err
		}

		if serviceType != "ClusterIP" && serviceType != "NodePort" && serviceType != "LoadBalancer" {
			return fmt.Errorf("the service-type must be one of: ClusterIP, NodePort or LoadBalancer")
		}

		overrides := map[string]string{}

		overrides["service.type"] = serviceType

		if persistence {
			overrides["persistence.enabled"] = "true"
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		portainerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("portainer/portainer").
			WithHelmURL("https://portainer.github.io/k8s/").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(portainerOptions)

		if err != nil {
			return err
		}

		println(portainerInstallMsg)
		return nil
	}
	return portainer
}

const PortainerInfoMsg = `
# Open the UI:

kubectl port-forward -n default svc/portainer 9000:9000 &

# http://127.0.0.1:9000

# If service type was NodePort, you can access it on http://node-ip:30777 as well

Find out more at https://www.portainer.io/
`

const portainerInstallMsg = `=======================================================================
= Portainer has been installed                                        =
=======================================================================` +
	"\n\n" + PortainerInfoMsg + "\n\n" + pkg.SupportMessageShort
