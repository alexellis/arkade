package venafi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/spf13/cobra"
)

// MakeCloudIssuer makes an app for the Venafi Cloud issuer
func MakeCloudIssuer() *cobra.Command {

	command := &cobra.Command{
		Use:   "cloud-issuer",
		Short: "Install the cert-manager issuer for Venafi cloud",
		Long: `Install the cert-manager issuer for Venafi cloud to obtain 
TLS certificates from enterprise-grade CAs.

Register and download your secret from your dashboard at:
https://www.venafi.com/venaficloud/devopsaccelerate`,
		Example: `  arkade venafi install cloud-issuer
  arkade venafi install cloud-issuer --help`,
		SilenceUsage: true,
	}

	command.Flags().String("secret", "", "Your Venafi cloud secret")
	command.Flags().StringP("secret-file", "f", "", "Your Venafi cloud secret from a file")
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

		// 	kubectl create secret generic \
		//    cloud-secret \
		//    --namespace='NAMESPACE OF YOUR ISSUER RESOURCE' \
		//    --from-literal=apikey='YOUR_CLOUD_API_KEY_HERE'

		tokenFileName, _ := command.Flags().GetString("secret-file")
		tokenString, _ := command.Flags().GetString("secret")

		var accessKeyFrom, accessKeyValue string
		if len(tokenFileName) > 0 {
			accessKeyFrom = "--from-file"
			accessKeyValue = tokenFileName
		} else if len(tokenString) > 0 {
			accessKeyFrom = "--from-literal"
			accessKeyValue = tokenString
		} else {
			return fmt.Errorf(`--secret or secret-file is a required`)
		}

		fmt.Printf(`Installing the cloud-issuer for you now.
Name: %s
Namespace: %s
Zone: %s

`, name, namespace, zone)
		clusterSecretName := name + "-secret"
		res, err := k8s.KubectlTask("create", "secret", "generic",
			clusterSecretName,
			"--namespace="+namespace,
			accessKeyFrom, "apikey="+accessKeyValue)

		if err != nil {
			return err
		} else if len(res.Stderr) > 0 && strings.Contains(res.Stderr, "AlreadyExists") {
			fmt.Printf("[Warning] secret %s already exists and will be used\n", clusterSecretName)
		} else if len(res.Stderr) > 0 {
			return fmt.Errorf("error from kubectl\n%q", res.Stderr)
		}

		tmpl, err := template.New("yaml").Parse(issuerTemplate)

		if err != nil {
			return err
		}

		var tpl bytes.Buffer

		err = tmpl.Execute(&tpl, struct {
			Name      string
			Namespace string
			Zone      string
		}{
			Name:      name,
			Namespace: namespace,
			Zone:      zone,
		})

		if err != nil {
			return err
		}

		d := os.TempDir()
		p := path.Join(d, "issuer.yaml")

		err = ioutil.WriteFile(p, tpl.Bytes(), os.ModePerm)
		if err != nil {
			return err
		}

		res, err = k8s.KubectlTask("apply", "-f", p)

		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(`unable to apply %s, error: %s`, p, res.Stderr)
		}
		fmt.Println(res.Stdout)

		// Check for error, mention to install cert-manager
		//no matches for kind "Issuer"

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
    zone: "{{.Zone}}"
    cloud:
      apiTokenSecretRef:
        name: {{.Name}}-secret
        key: apikey
`
