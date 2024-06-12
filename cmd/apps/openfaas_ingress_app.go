// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bytes"
	"errors"
	"fmt"
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
	IssuerName       string
	IssuerAPI        string
	IngressName      string
	ClusterIssuer    bool
	IngressService   string
	Namespace        string
}

// MakeInstallOpenFaaSIngress will install a clusterissuer and request a cert from certmanager for the domain you specify
func MakeInstallOpenFaaSIngress() *cobra.Command {
	var openfaasIngress = &cobra.Command{
		Use:   "openfaas-ingress",
		Short: "Install openfaas ingress with TLS",
		Long:  `Install openfaas ingress. Requires cert-manager 1.0.0+ the cluster.`,
		Example: `  arkade install openfaas-ingress \
  --domain openfaas.example.com \
  --email openfaas@example.com`,
		SilenceUsage: true,
	}

	openfaasIngress.Flags().StringP("namespace", "n", "openfaas", "Give a Kubernetes namespace")
	openfaasIngress.Flags().StringP("domain", "d", "", "Custom Ingress Domain")
	openfaasIngress.Flags().StringP("email", "e", "", "Letsencrypt Email")
	openfaasIngress.Flags().String("ingress-class", "nginx", `Ingress class to be used such as "nginx" or "traefik"`)
	openfaasIngress.Flags().Bool("staging", false, "set --staging to true to use the staging Letsencrypt issuer")
	openfaasIngress.Flags().String("issuer", "", "provide the name of a pre-existing issuer, rather than creating one for LetsEncrypt")
	openfaasIngress.Flags().Bool("cluster-issuer", false, "set to true to create a clusterissuer rather than a namespaces issuer (default: false)")
	openfaasIngress.Flags().String("oidc-plugin-domain", "", "Set to the auth domain for openfaas OIDC installations")

	openfaasIngress.RunE = func(command *cobra.Command, args []string) error {

		email, _ := command.Flags().GetString("email")
		domain, _ := command.Flags().GetString("domain")
		issuer, _ := command.Flags().GetString("issuer")
		namespace, _ := command.Flags().GetString("namespace")
		ingressClass, _ := command.Flags().GetString("ingress-class")

		if email == "" || domain == "" {
			return errors.New("both --email and --domain flags should be set and not empty, please set these values")
		}

		if ingressClass == "" {
			return errors.New("--ingress-class must be set")
		}

		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		staging, _ := command.Flags().GetBool("staging")
		clusterIssuer, _ := command.Flags().GetBool("cluster-issuer")

		if len(issuer) > 0 {
			fmt.Printf("Using existing issuer: %s\n", issuer)
		} else if err := createIssuer(domain, email, ingressClass, "openfaas-gateway", staging, clusterIssuer, namespace); err != nil {
			return err
		}

		caps, err := k8s.GetCapabilities()
		if err != nil {
			return err
		}

		hasNetworking := caps["networking.k8s.io/v1"]

		if err := createIngress(domain, email, ingressClass, "openfaas-gateway", staging, clusterIssuer, issuer, namespace, hasNetworking); err != nil {
			return err
		}

		oidcDomain, _ := command.Flags().GetString("oidc-plugin-domain")

		if len(oidcDomain) > 0 {
			if err := createIngress(oidcDomain, email, ingressClass, "oidc-plugin", staging, clusterIssuer, issuer, namespace, hasNetworking); err != nil {
				return err
			}
		}

		fmt.Println(openfaasIngressInstallMsg)

		return nil
	}

	return openfaasIngress
}

func createIssuer(domain, email, ingressClass, ingressName string, staging bool, clusterIssuer bool, namespace string) error {
	yamlBytes, templateErr := buildIssuerYAML(email, ingressClass, staging, clusterIssuer, namespace)
	if templateErr != nil {
		log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
		return templateErr
	}

	tempFile, tempFileErr := writeTempFile(yamlBytes, fmt.Sprintf("%s-issuer.yaml", ingressName))
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
		return fmt.Errorf(`unable to apply YAML files.
Have you got OpenFaaS running in the openfaas namespace and cert-manager 1.0.0 or higher installed in cert-manager namespace? %s`,
			res.Stderr)
	}
	return nil
}

func createIngress(domain, email, ingressClass, ingressName string, staging bool, clusterIssuer bool, issuerName, namespace string, hasNetworking bool) error {
	yamlBytes, templateErr := buildOpenfaasIngressYAML(domain, email, ingressClass, ingressName, staging, clusterIssuer, issuerName, namespace, hasNetworking)
	if templateErr != nil {
		log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
		return templateErr
	}

	tempFile, tempFileErr := writeTempFile(yamlBytes, fmt.Sprintf("%s-ingress.yaml", ingressName))
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
Have you got OpenFaaS running in the openfaas namespace and cert-manager 1.0.0 or higher installed in cert-manager namespace? %s`,
			res.Stderr)
	}
	return nil
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

	err := os.WriteFile(filename, input, 0744)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func buildOpenfaasIngressYAML(domain, email, ingressClass, ingressName string, staging, clusterIssuer bool, issuerName, namespace string, hasNetworking bool) ([]byte, error) {
	tmplString := openfaasIngressExtensionTemplate
	if hasNetworking {
		tmplString = openfaasIngressNetworkingTemplate
	}

	templ, err := template.New("yaml").Parse(tmplString)
	if err != nil {
		return nil, err
	}

	ingressService := "gateway"
	if ingressName == "oidc-plugin" {
		ingressService = ingressName
	}

	inputData := inputData{
		Namespace:      namespace,
		IngressDomain:  domain,
		IngressClass:   ingressClass,
		IngressName:    ingressName,
		IngressService: ingressService,
		ClusterIssuer:  clusterIssuer,
	}

	if len(issuerName) > 0 {
		inputData.IssuerName = issuerName
	} else if staging {
		inputData.IssuerName = "letsencrypt-staging"
	} else {
		inputData.IssuerName = "letsencrypt-prod"
	}

	var tpl bytes.Buffer

	if err = templ.Execute(&tpl, inputData); err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}

func buildIssuerYAML(email, ingressClass string, staging, clusterIssuer bool, namespace string) ([]byte, error) {
	templ, err := template.New("issuer-yaml").Parse(http01IssuerTemplate)

	if err != nil {
		return nil, err
	}

	inputData := inputData{
		CertmanagerEmail: email,
		Namespace:        namespace,
		IssuerName:       "letsencrypt-prod",
		IssuerAPI:        "https://acme-v02.api.letsencrypt.org/directory",
		ClusterIssuer:    clusterIssuer,
		IngressClass:     ingressClass,
	}

	if staging {
		inputData.IssuerName = "letsencrypt-staging"
		inputData.IssuerAPI = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	var tpl bytes.Buffer
	if err = templ.Execute(&tpl, inputData); err != nil {
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
	"\n\n" + OpenfaasIngressInfoMsg + "\n\n" + pkg.SupportMessageShort

// Ingress in extensions/v1beta1 are removed in k8s 1.22+, July 2021
var openfaasIngressExtensionTemplate = `
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{.IngressName}}
  namespace: {{.Namespace}}
  annotations:
{{- if .ClusterIssuer }}
    cert-manager.io/cluster-issuer: {{.IssuerName}}
{{- else }}
    cert-manager.io/issuer: {{.IssuerName}}
{{- end }}
    kubernetes.io/ingress.class: {{.IngressClass}}
    cert-manager.io/common-name: {{.IngressDomain}}
spec:
  rules:
  - host: {{.IngressDomain}}
    http:
      paths:
      - backend:
          serviceName: {{.IngressService}}
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - {{.IngressDomain}}
    secretName: {{.IngressName}}
`

// Ingress in networking.k8s.io/v1 was added in k8s 1.19+
// this includes the pathType change added in 1.18
var openfaasIngressNetworkingTemplate = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{.IngressName}}
  namespace: {{.Namespace}}
  annotations:
{{- if .ClusterIssuer }}
    cert-manager.io/cluster-issuer: {{.IssuerName}}
{{- else }}
    cert-manager.io/issuer: {{.IssuerName}}
{{- end }}
    kubernetes.io/ingress.class: {{.IngressClass}}
    cert-manager.io/common-name: {{.IngressDomain}}
spec:
  rules:
  - host: {{.IngressDomain}}
    http:
      paths:
      - path: /
        pathType: ImplementationSpecific
        backend:
          service:
            name: {{.IngressService}}
            port:
              number: 8080
  tls:
  - hosts:
    - {{.IngressDomain}}
    secretName: {{.IngressName}}
`

var http01IssuerTemplate = `
apiVersion: cert-manager.io/v1
{{- if .ClusterIssuer }}
kind: ClusterIssuer
{{- else }}
kind: Issuer
{{- end }}
metadata:
  name: {{.IssuerName}}
{{- if not .ClusterIssuer }}
  namespace: {{.Namespace}}
{{- end }}
spec:
  acme:
    email: {{.CertmanagerEmail}}
    server: {{.IssuerAPI}}
    privateKeySecretRef:
      name: example-issuer-account-key
    solvers:
    - selector: {}
      http01:
        ingress:
          class: {{.IngressClass}}`
