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

const (
	MetalLBManifest = "https://raw.githubusercontent.com/metallb/metallb/v0.13.10/config/manifests/metallb-native.yaml"
)

func MakeInstallMetalLB() *cobra.Command {
	var command = &cobra.Command{
		Use:          "metallb-arp",
		Short:        "Install MetalLB in L2 (ARP) mode",
		Long:         `Install a network load-balancer implementation for Kubernetes using standard routing protocols`,
		Example:      `arkade install metallb-arp --address-range=<cidr>`,
		SilenceUsage: true,
	}

	command.Flags().String("address-range", "192.168.0.0/24", "Address range for LoadBalancer services")
	command.Flags().String("memberlist-secretkey", "", "A predefined memberlist secretkey, a random key is generated if omitted")

	command.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}
		_, err = command.Flags().GetString("address-range")
		if err != nil {
			return fmt.Errorf("error with --address-range usage: %s", err)
		}
		_, err = command.Flags().GetString("memberlist-secretkey")
		if err != nil {
			return fmt.Errorf("error with --memberlist-secretkey usage: %s", err)
		}

		return nil
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		addressRange, _ := command.Flags().GetString("address-range")

		if err := k8s.Kubectl("apply", "-f", MetalLBManifest); err != nil {
			return err
		}

		var token string

		token, _ = command.Flags().GetString("memberlist-secretkey")

		if len(token) < 1 {
			randomToken := make([]byte, 32)
			rand.Read(randomToken)
			token = b64.StdEncoding.EncodeToString(randomToken)
		}

		secret := types.K8sSecret{
			Type:      "generic",
			Name:      "memberlist",
			Namespace: "metallb-system",
			SecretData: []types.SecretsData{{
				Type:  "string-literal",
				Key:   "secretkey",
				Value: token,
			}},
		}

		if err := k8s.CreateSecret(secret); err != nil {
			return fmt.Errorf("create secret error: %+v", err)
		}

		configMap := fmt.Sprintf(metalLBConfigMap, addressRange)

		if err := k8s.KubectlIn(strings.NewReader(configMap), "apply", "-f", "-"); err != nil {
			return fmt.Errorf("create configmap error: %+v", err)
		}

		fmt.Println(MetalLBInstallMsg)

		return nil
	}

	return command
}

const MetalLBInfoMsg = `
# Get the memberlist secretkey:
export SECRET_KEY=$(kubectl get secret -n metallb-system memberlist \
	-o jsonpath="{.data.secretkey}" | base64 --decode)

echo "Secret Key: $SECRET_KEY"

# Review the generated configuration:

kubectl get configmap -n metallb-system config --template={{.data.config}}

# Find out more at: https://metallb.universe.tf/
`

const MetalLBInstallMsg = `=======================================================================
= MetalLB has been installed.                                         =
=======================================================================` +
	"\n\n" + MetalLBInfoMsg + "\n\n" + pkg.SupportMessageShort

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
