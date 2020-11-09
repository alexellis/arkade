// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
	command.Flags().String("url", "", "The URL for your TPP server")
	command.Flags().StringP("username", "u", "", "Your TPP username")
	command.Flags().StringP("password", "p", "", "Your TPP password")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
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
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing the TPP issuer for you now.")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		name, _ := cmd.Flags().GetString("name")
		clusterIssuer, _ := cmd.Flags().GetBool("cluster-issuer")

		fmt.Println(url, username, password, name, clusterIssuer)
		return nil
	}

	return command
}

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
      caBundle: {{.CABundle}}
      credentialsRef:
        name: {{.Name}}-secret
`
