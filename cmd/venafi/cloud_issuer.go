package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeCloudIssuer makes an app for the Venafi Cloud issuer
func MakeCloudIssuer() *cobra.Command {

	command := &cobra.Command{
		Use:   "cloud-issuer",
		Short: "Install the cert-manager issuer for Venafi cloud",
		Long: `Install the cert-manager issuer for Venafi cloud to obtain 
TLS certificates from enterprise-grade CAs.`,
		Example: `  arkade venafi install cloud-issuer
  arkade venafi install cloud-issuer --help`,
		SilenceUsage: true,
	}

	command.Flags().String("cloud-secret", "", "Your Venafi cloud secret")
	command.Flags().StringP("cloud-secret-file", "f", "", "Your Venafi cloud secret from a file")
	command.Flags().String("namespace", "default", "Namespace for the issuer")
	command.Flags().String("name", "cloud-venafi-issuer", "Name for the issuer")
	command.Flags().StringP("zone", "z", "", "The zone for Venafi cloud")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		name, err := command.Flags().GetString("name")
		if err != nil {
			return err
		}
		namespace, err := command.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		zone, err := command.Flags().GetString("zone")
		if err != nil {
			return err
		}

		if len(zone) == 0 {
			return fmt.Errorf("a zone is required")
		}

		fmt.Printf(`Installing the cloud-issuer for you now.
Name: %s
Namespace: %s
Zone: %s

`, name, namespace, zone)
		// 	kubectl create secret generic \
		//    cloud-secret \
		//    --namespace='NAMESPACE OF YOUR ISSUER RESOURCE' \
		//    --from-literal=apikey='YOUR_CLOUD_API_KEY_HERE'

		fmt.Println(`# Query the status of the issuer:
kubectl get issuer ` + name + ` -n ` + namespace + ` -o wide

# Find out how to issue a certificate with cert-manager:
# https://cert-manager.io/docs/usage/certificate/
`)

		return nil
	}

	return command
}

const CloudIssuerInfo = `# Check the status of the issuer:
kubectl get issuer name -n namespace -o wide

# Find out how to issue a certificate with cert-manager:
# https://cert-manager.io/docs/usage/certificate/`

const issuerTemplate = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  venafi:
    zone: "{{.Zone}}" # Set this to the Venafi policy zone you want to use
    cloud:
      apiTokenSecretRef:
        name: cloud-secret
		key: apikey
`
