package apps

import (
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallMetalLB() *cobra.Command {
	var command = &cobra.Command{
		Use:          "metallb",
		Short:        "Install metallb",
		Long:         `Install metallb for service type:LoadBalancer`,
		Example:      `arkade install metallb`,
		SilenceUsage: true,
	}

	command.Flags().String("address-range", "192.168.0.0/24", "Address range for loadBalancer services")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		addressRange, _ := command.Flags().GetString("address-range")

		err := k8s.Kubectl("apply", "-f",
			"https://raw.githubusercontent.com/metallb/metallb/v0.9.5/manifests/namespace.yaml")
		if err != nil {
			return err
		}

		err = k8s.Kubectl("apply", "-f",
			"https://raw.githubusercontent.com/metallb/metallb/v0.9.5/manifests/metallb.yaml")
		if err != nil {
			return err
		}

		token := make([]byte, 32)
		rand.Read(token)

		secret := types.K8sSecret{
			Type:      "generic",
			Name:      "memberlist",
			Namespace: "metallb-system",
			SecretData: []types.SecretsData{{
				Type:  "string-literal",
				Key:   "secretkey",
				Value: b64.StdEncoding.EncodeToString(token),
			}},
		}

		err = k8s.CreateSecret(secret)
		if err != nil {
			return fmt.Errorf("Create secret error: %+v", err)
		}

		configMap := fmt.Sprintf(metalLBConfigMap, addressRange)

		err = k8s.KubectlIn(strings.NewReader(configMap), "apply", "-f", "-")
		if err != nil {
			return fmt.Errorf("Create configmap error: %+v", err)
		}

		fmt.Println(MetalLBInstallMsg)

		return nil
	}

	return command
}

const MetalLBInfoMsg = `# Find out more at:
# https://metallb.universe.tf/
`

const MetalLBInstallMsg = `=======================================================================
= metalLB has been installed.                                  =
=======================================================================` +
	"\n\n" + MetalLBInfoMsg + "\n\n" + pkg.ThanksForUsing

const metalLBConfigMap = `
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: metallb-system
  name: config
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - %s
`
