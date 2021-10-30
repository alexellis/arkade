package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallKuma() *cobra.Command {
	var kuma = &cobra.Command{
		Use:          "kuma",
		Short:        "Install Kuma",
		Long:         "Install Kuma",
		Example:      ` arkade app install kuma`,
		SilenceUsage: true,
	}

	kuma.Flags().String("namespace", "kuma-system", "Namespace for the app")

	kuma.Flags().String("control-plane-mode", "standalone", "Kuma CP modes: one of standalone,zone,global")
	kuma.Flags().Bool("auto-scale", false, "Enable Horizontal Pod Autoscaling (requires the Metrics Server)")

	kuma.Flags().Bool("use-cni", false, "Use CNI instead of proxy init container")
	kuma.Flags().Bool("ingress", false, "Deploy Ingress for cross cluster communication")

	kuma.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set ingress.enabled=false)")

	kuma.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		controlPlaneMode, err := command.Flags().GetString("control-plane-mode")
		if err != nil {
			return err
		}
		// control-plane-mode  standalone,zone,global
		if controlPlaneMode != "standalone" && controlPlaneMode != "zone" && controlPlaneMode != "global" {
			return fmt.Errorf(`kuma's control-plane mode must be one of global, zone or standalone`)
		}

		_, err = command.Flags().GetBool("auto-scale")
		if err != nil {
			return err
		}

		_, err = command.Flags().GetBool("use-cni")
		if err != nil {
			return err
		}

		_, err = command.Flags().GetBool("ingress")
		if err != nil {
			return err
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		return nil
	}

	kuma.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)
		if suffix := getValuesSuffix(arch); suffix == "-armhf" {
			return fmt.Errorf(`kuma is currently not supported on armhf architectures`)
		}

		namespace, _ := command.Flags().GetString("namespace")
		controlPlaneMode, _ := command.Flags().GetString("control-plane-mode")
		autoScale, _ := command.Flags().GetBool("auto-scale")
		useCNI, _ := command.Flags().GetBool("use-cni")
		ingress, _ := command.Flags().GetBool("ingress")
		customFlags, _ := command.Flags().GetStringArray("set")

		overrides := map[string]string{}

		overrides["controlPlane.mode"] = controlPlaneMode

		if autoScale {
			overrides["controlPlane.autoscaling.enabled"] = "true"
		}

		if useCNI {
			overrides["cni.enabled"] = "true"
		}

		if ingress {
			overrides["ingress.enabled"] = "true"
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kumaOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kuma/kuma").
			WithHelmURL("https://kumahq.github.io/charts").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kumaOptions)
		if err != nil {
			return err
		}

		println(kumaInstallMsg)
		return nil
	}

	return kuma
}

const KumaInfoMsg = `
Kuma has been installed, you can access the control-plane via either the GUI, kubectl, the HTTP API, or the CLI:

# You can use kubectl without forwarding the port:
kubectl get meshes

# Forward the port with:
kubectl port-forward svc/kuma-control-plane -n kuma-system 5681:5681

## You can access the GUI on: http://127.0.0.1:5681/gui

## You can access the API on: http://127.0.0.1:5681

## You can use kumactl: kumactl get meshes

# Find out more on the project homepage:
# https://kuma.io/docs/1.4.x/documentation/overview/#kubernetes-mode
`

const kumaInstallMsg = `=======================================================================
=                      Kuma has been installed                        =
=======================================================================
 ` + pkg.ThanksForUsing + KumaInfoMsg
