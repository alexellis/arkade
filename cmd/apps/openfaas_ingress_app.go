// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"

	"text/template"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

type inputData struct {
	IngressDomain    string
	CertmanagerEmail string
	IngressClass     string
	IssuerType       string
	IssuerAPI        string
	ClusterIssuer    bool
}

//MakeInstallOpenFaaSIngess will install a clusterissuer and request a cert from certmanager for the domain you specify
func MakeInstallOpenFaaSIngress() *cobra.Command {
	var openfaasIngress = &cobra.Command{
		Use:          "openfaas-ingress",
		Short:        "Install openfaas ingress with TLS",
		Long:         `Install openfaas ingress. Requires cert-manager 0.11.0 or higher installation in the cluster. Please set --domain to your custom domain and set --email to your email - this email is used by letsencrypt for domain expiry etc.`,
		Example:      `  arkade install openfaas-ingress --domain openfaas.example.com --email openfaas@example.com`,
		SilenceUsage: true,
	}

	openfaasIngress.Flags().StringP("domain", "d", "", "Custom Ingress Domain")
	openfaasIngress.Flags().StringP("email", "e", "", "Letsencrypt Email")
	openfaasIngress.Flags().String("ingress-class", "nginx", "Ingress class to be used such as nginx or traefik")
	openfaasIngress.Flags().Bool("staging", false, "set --staging to true to use the staging Letsencrypt issuer")
	openfaasIngress.Flags().Bool("cluster-issuer", false, "set to true to create a clusterissuer rather than a namespaces issuer (default: false)")

	openfaasIngress.RunE = func(command *cobra.Command, args []string) error {

		email, _ := command.Flags().GetString("email")
		domain, _ := command.Flags().GetString("domain")
		ingressClass, _ := command.Flags().GetString("ingress-class")

		if email == "" || domain == "" {
			return errors.New("both --email and --domain flags should be set and not empty, please set these values")
		}

		if ingressClass == "" {
			return errors.New("--ingress-class must be set")
		}

		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		staging, _ := command.Flags().GetBool("staging")
		clusterIssuer, _ := command.Flags().GetBool("cluster-issuer")

		yamlBytes, templateErr := buildYAML(domain, email, ingressClass, staging, clusterIssuer)
		if templateErr != nil {
			log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
			return templateErr
		}

		tempFile, tempFileErr := writeTempFile(yamlBytes, "temp_openfaas_ingress.yaml")
		if tempFileErr != nil {
			log.Print("Unable to save generated yaml file into the temporary directory")
			return tempFileErr
		}

		res, err := k8s.KubectlTask("apply", "-f", tempFile)

		if err != nil {
			log.Print(err)
			return err
		}

		if res.ExitCode != 0 {
			return fmt.Errorf(`Unable to apply YAML files.
Have you got OpenFaaS running in the openfaas namespace and cert-manager 0.11.0 or higher installed in cert-manager namespace? %s`,
				res.Stderr)
		}

		fmt.Println(openfaasIngressInstallMsg)

		return nil
	}

	return openfaasIngress
}

func createTempDirectory(directory string) (string, error) {
	tempDirectory := filepath.Join(os.TempDir(), directory)
	if _, err := os.Stat(tempDirectory); os.IsNotExist(err) {
		log.Printf(tempDirectory)
		errr := os.Mkdir(tempDirectory, 0744)
		if errr != nil {
			log.Printf("couldnt make dir %s", err)
			return "", err
		}
	}

	return tempDirectory, nil
}

func writeTempFile(input []byte, fileLocation string) (string, error) {
	var tempDirectory, dirErr = createTempDirectory(".arkade/")
	if dirErr != nil {
		return "", dirErr
	}

	filename := filepath.Join(tempDirectory, fileLocation)

	err := ioutil.WriteFile(filename, input, 0744)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func buildYAML(domain, email, ingressClass string, staging, clusterIssuer bool) ([]byte, error) {
	tmpl, err := template.New("yaml").Parse(ingressYamlTemplate)

	if err != nil {
		return nil, err
	}

	inputData := inputData{
		IngressDomain:    domain,
		CertmanagerEmail: email,
		IngressClass:     ingressClass,
		IssuerType:       "letsencrypt-prod",
		IssuerAPI:        "https://acme-v02.api.letsencrypt.org/directory",
		ClusterIssuer:    clusterIssuer,
	}

	if staging {
		inputData.IssuerType = "letsencrypt-staging"
		inputData.IssuerAPI = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, inputData)

	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

const OpenfaasIngressInfoMsg = `# You will need to ensure that your domain points to your cluster and is
# accessible through ports 80 and 443. 
#
# This is used to validate your ownership of this domain by LetsEncrypt
# and then you can use https with your installation. 

# Ingress to your domain has been installed for OpenFaaS
# to see the ingress record run
kubectl get -n openfaas ingress openfaas-gateway

# Check the cert-manager logs with:
kubectl logs -n cert-manager deploy/cert-manager

# A cert-manager ClusterIssuer has been installed into the default
# namespace - to see the resource run
kubectl describe ClusterIssuer letsencrypt-prod

# To check the status of your certificate you can run
kubectl describe -n openfaas Certificate openfaas-gateway

# It may take a while to be issued by LetsEncrypt, in the meantime a 
# self-signed cert will be installed`

const openfaasIngressInstallMsg = `=======================================================================
= OpenFaaS Ingress and cert-manager ClusterIssuer have been installed =
=======================================================================` +
	"\n\n" + OpenfaasIngressInfoMsg + "\n\n" + pkg.ThanksForUsing

var ingressYamlTemplate = `
apiVersion: extensions/v1beta1 
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
{{- if .ClusterIssuer }}
    cert-manager.io/cluster-issuer: {{.IssuerType}}
{{- else }}
    cert-manager.io/issuer: {{.IssuerType}}
{{- end }}
    kubernetes.io/ingress.class: {{.IngressClass}}
spec:
  rules:
  - host: {{.IngressDomain}}
    http:
      paths:
      - backend:
          serviceName: gateway
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - {{.IngressDomain}}
    secretName: openfaas-gateway
---
apiVersion: cert-manager.io/v1alpha2
{{- if .ClusterIssuer }}
kind: ClusterIssuer
{{- else }}
kind: Issuer
{{- end }}
metadata:
  name: {{.IssuerType}}
{{- if not .ClusterIssuer }}
  namespace: openfaas
{{- end }}
spec:
  acme:
    email: {{.CertmanagerEmail}}
    server: {{.IssuerAPI}}
    privateKeySecretRef:
      name: example-issuer-account-key
    solvers:
    - http01:
        ingress:
          class: {{.IngressClass}}`
