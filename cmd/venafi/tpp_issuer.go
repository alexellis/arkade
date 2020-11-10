// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ./arkade venafi install tpp-issuer --url https://tpp.venafidemo.com/vedsdk/
// --ca-bundle ~/Downloads/Trust\ Bundle.pem --zone arkade --custom-field cost-center=venafi.com --custom-field org=DE

// MakeTPPIssuer makes the app for the TPP issuer
func MakeTPPIssuer() *cobra.Command {

	command := &cobra.Command{
		Use:   "tpp-issuer",
		Short: "Install the cert-manager issuer for Venafi TPP",
		Long: `Install the cert-manager issuer for Venafi TPP to obtain 
TLS certificates from enterprise-grade CAs from self-hosted Venafi 
instances.`,
		Example: `  arkade venafi install cloud-issuer
  arkade venafi install tpp-issuer --help`,
		SilenceUsage: true,
	}

	command.Flags().String("name", "tpp-venafi-issuer", "The name for the Issuer")
	command.Flags().String("namespace", "default", "The Kubernetes namespace for the Issuer")
	command.Flags().Bool("cluster-issuer", false, "Use a ClusterIssuer instead of an Issuer for the given namespace")
	command.Flags().String("url", "", "The URL for your TPP server including the \"/vedsdk\" suffix")
	command.Flags().StringP("username", "u", "", "Your TPP username")
	command.Flags().StringP("zone", "z", "", "The zone for the issuer")
	command.Flags().StringP("password", "p", "", "Your TPP password")
	command.Flags().String("ca-bundle", "", "The path to a ca-bundle file")
	command.Flags().StringArray("custom-fields", []string{""}, "A number of custom fields for the TPP issuer and its policy")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		if len(username) == 0 {
			return fmt.Errorf("username is required")
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		if len(password) == 0 {
			return fmt.Errorf("password is required")
		}
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}
		if len(url) == 0 {
			return fmt.Errorf("url is required")
		}
		_, err = cmd.Flags().GetBool("cluster-issuer")
		if err != nil {
			return err
		}
		bundlePath, err := cmd.Flags().GetString("ca-bundle")
		if err != nil {
			return err
		}
		if len(bundlePath) > 0 {
			if _, err = os.Stat(bundlePath); err != nil {
				return errors.Wrapf(err, "ca-bundle %q not found", bundlePath)
			}
		}
		zone, err := command.Flags().GetString("zone")
		if err != nil {
			return err
		}
		if len(zone) == 0 {
			return fmt.Errorf("a zone is required")
		}
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing the TPP issuer for you now.")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		name, _ := cmd.Flags().GetString("name")
		zone, _ := cmd.Flags().GetString("zone")
		namespace, _ := cmd.Flags().GetString("namespace")
		caBundlePath, _ := cmd.Flags().GetString("ca-bundle")

		clusterIssuer, _ := command.Flags().GetBool("cluster-issuer")

		kind := "Issuer"
		if clusterIssuer {
			kind = "ClusterIssuer"
		}
		encodedBundle := ""
		if len(caBundlePath) > 0 {
			data, err := ioutil.ReadFile(caBundlePath)
			if err != nil {
				return err
			}
			encodedBundle = base64.StdEncoding.EncodeToString(data)
		}

		clusterSecretName := name + "-secret"
		res, err := k8s.KubectlTask("create", "secret", "generic",
			clusterSecretName,
			"--namespace="+namespace,
			"--from-literal",
			"username="+username,
			"--from-literal",
			"password="+password)

		if err != nil {
			return err
		} else if len(res.Stderr) > 0 && strings.Contains(res.Stderr, "AlreadyExists") {
			fmt.Printf("[Warning] secret %s already exists and will be used\n", clusterSecretName)
		} else if len(res.Stderr) > 0 {
			return fmt.Errorf("error from kubectl\n%q", res.Stderr)
		}

		// fmt.Println(encodedBundle)
		fmt.Println(url, username, password, name, clusterIssuer)

		manifest, err := templateManifest(tppIssuerTemplate, struct {
			Name      string
			Namespace string
			Zone      string
			Kind      string
			URL       string
			CABundle  string
		}{
			Name:      name,
			Namespace: namespace,
			Zone:      zone,
			Kind:      kind,
			URL:       url,
			CABundle:  encodedBundle,
		})

		if err != nil {
			return err
		}

		p, err := writeFile("tpp-issuer.yaml", manifest)

		res, err = k8s.KubectlTask("apply", "-f", p)

		if err != nil {
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(`unable to apply %s, error: %s`, p, res.Stderr)
		}

		fmt.Println(res.Stdout)

		fmt.Println(`# Query the status of the issuer:
kubectl get issuer ` + name + ` -n ` + namespace + ` -o wide

# Find out how to issue a certificate with cert-manager:
# https://cert-manager.io/docs/usage/certificate/
		`)

		return nil
	}

	return command
}

const TPPIssuerInfo = `# Check the status of the issuer:
kubectl get issuer name -n namespace -o wide

# Find out how to issue a certificate with cert-manager:
# https://cert-manager.io/docs/usage/certificate/`

const tppIssuerTemplate = `apiVersion: cert-manager.io/v1
kind: {{.Kind}}
metadata:
  name: {{.Name}}
{{- if eq .Kind "Issuer" }}
  namespace: {{.Namespace}}
{{- end }}
spec:
  venafi:
    zone: "{{.Zone}}"
    tpp:
      url: {{.URL}}
{{- if ne .CABundle "" }}
      caBundle: "{{.CABundle}}"
{{- end }}
      credentialsRef:
        name: {{.Name}}-secret
`
