package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallKongIngress() *cobra.Command {
	var command = &cobra.Command{
		Use:          "kong-ingress",
		Short:        "Install kong-ingress for OpenFaaS",
		Long:         `Install kong-ingress for OpenFaaS`,
		Example:      `arkade install kong-ingress`,
		SilenceUsage: true,
	}

	command.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	command.Flags().Bool("update-repo", true, "Update the helm repo")

	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		updateRepo, _ := command.Flags().GetBool("update-repo")

		namespace, _ := command.Flags().GetString("namespace")

		overrides := map[string]string{}
		// always need to be set to  false
		overrides["ingressController.installCRDs"] = "false"

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kongIngressAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kong/kong").
			WithHelmURL("https://charts.konghq.com/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(kongIngressAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(kongIngressInstallMsg)

		return nil
	}

	return command
}

const KongIngressInfoMsg = `# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# Find out more on the project homepage:
# https://github.com/Kong/kubernetes-ingress-controller`

const kongIngressInstallMsg = `=======================================================================
= kong-ingress has been installed.                                  =
=======================================================================` +
	"\n\n" + KongIngressInfoMsg + "\n\n" + pkg.SupportMessageShort
