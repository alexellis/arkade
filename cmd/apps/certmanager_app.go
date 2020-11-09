// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bytes"
	"fmt"
	"log"

	"text/template"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type IssuerInput struct {
	CertmanagerEmail string
	IngressClass     string
	IssuerName       string
	IssuerAPI        string
}

func MakeInstallCertManager() *cobra.Command {
	var certManager = &cobra.Command{
		Use:          "cert-manager",
		Short:        "Install cert-manager",
		Long:         "Install cert-manager for TLS certificates management",
		Example:      "arkade install cert-manager",
		SilenceUsage: true,
	}

	certManager.Flags().StringP("namespace", "n", "cert-manager", "The namespace to install cert-manager")
	certManager.Flags().StringP("version", "v", "v1.5.4", "The version of cert-manager to install, has to be >= v1.0.0")
	certManager.Flags().Bool("update-repo", true, "Update the helm repo")
	// Acme options
	certManager.Flags().Bool("acme", false, "Enable automated usage of kubernetes.io/tls-acme: \"true\" annotation (implies wait)")
	certManager.Flags().Bool("staging", false, "set --staging to true to generate a staging Letsencrypt cluster issuer")
	certManager.Flags().StringP("email", "e", "", "Letsencrypt Email")
	certManager.Flags().String("ingress-class", "nginx", "Ingress class to be used such as nginx or traefik")
	// Override flags
	certManager.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set key=value)")

	certManager.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		wait, _ := command.Flags().GetBool("wait")

		namespace, _ := command.Flags().GetString("namespace")
		version, _ := command.Flags().GetString("version")

		if !semver.IsValid(version) {
			return fmt.Errorf("%q is not a valid semver version", version)
		}

		updateRepo, _ := certManager.Flags().GetBool("update-repo")

		overrides := map[string]string{}
		overrides["installCRDs"] = "true"

		// Allow usage of lego behavior, see https://cert-manager.io/docs/usage/ingress/#optional-configuration
		useAcme, _ := certManager.Flags().GetBool("acme")
		if useAcme {
			// Configure default issuer name in helm according to --staging flag
			staging, _ := certManager.Flags().GetBool("staging")
			issuerName := "letsencrypt-production-issuer"
			if staging {
				issuerName = "letsencrypt-staging-issuer"
			}

			overrides["ingressShim.defaultIssuerName"] = issuerName
			overrides["ingressShim.defaultIssuerKind"] = "ClusterIssuer"
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return err
		}
		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		certmanagerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("jetstack/cert-manager").
			WithHelmURL("https://charts.jetstack.io").
			WithOverrides(overrides).
			WithWait(wait || useAcme).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(certmanagerOptions)
		if err != nil {
			return err
		}

		// Install default issuer if --acme is given
		if useAcme {
			email, _ := command.Flags().GetString("email")
			ingressClass, _ := command.Flags().GetString("ingress-class")
			staging, _ := certManager.Flags().GetBool("staging")

			err = WriteDefaultIssuer(email, ingressClass, staging)
			if err != nil {
				return err
			}
			// Print ACME installation message
			fmt.Println(AcmeInfoMsg)
		}

		fmt.Println(certManagerInstallMsg)

		return nil
	}

	return certManager
}

const CertManagerInfoMsg = `# Get started with cert-manager here:
# https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

const AcmeInfoMsg = `# A default ClusterIssuer was installed and configured for acme usage (with kubernetes.io/tls-acme: "true"):
# https://cert-manager.io/docs/usage/ingress/#optional-configuration`

const certManagerInstallMsg = `=======================================================================
= cert-manager  has been installed.                                   =
=======================================================================` +
	"\n\n" + CertManagerInfoMsg + "\n\n" + pkg.ThanksForUsing

func WriteDefaultIssuer(email string, ingressClass string, staging bool) (err error) {
	log.Printf("Installing default letsencrypt issuer\n")

	yamlBytes, templateErr := buildDefaultIssuerYAML(email, ingressClass, staging)
	if templateErr != nil {
		fmt.Errorf("unable to install the application. Could not build the templated yaml file for the resources")
		return templateErr
	}

	tempFile, tempFileErr := writeTempFile(yamlBytes, "temp_registry_ingress.yaml")
	if tempFileErr != nil {
		fmt.Errorf("unable to save generated yaml file into the temporary directory")
		return tempFileErr
	}

	res, err := k8s.KubectlTask("apply", "-f", tempFile)

	if err != nil {
		return err
	}

	if res.ExitCode > 0 {
		return fmt.Errorf("error installing cluster default issuer: error: %s", res.Stderr)
	}

	return nil
}

func buildDefaultIssuerYAML(email string, ingressClass string, staging bool) ([]byte, error) {
	tmpl, err := template.New("yaml").Parse(clusterIssuerTemplate)

	if err != nil {
		return nil, err
	}

	inputData := IssuerInput{
		CertmanagerEmail: email,
		IngressClass:     ingressClass,
		IssuerName:       "letsencrypt-production-issuer",
		IssuerAPI:        "https://acme-v02.api.letsencrypt.org/directory",
	}

	if staging {
		inputData.IssuerName = "letsencrypt-staging-issuer"
		inputData.IssuerAPI = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	if len(email) > 0 {
		inputData.CertmanagerEmail = fmt.Sprintf("    email: %s", email)
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, inputData)

	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

var clusterIssuerTemplate = `
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: {{.IssuerName}}
  namespace: cert-manager
spec:
  acme:
{{.CertmanagerEmail}}
    # The ACME server URL for production
    server: {{.IssuerAPI}}
    # Name of a secret used to store the ACME account private key
    privateKeySecretRef:
      name: {{.IssuerName}}-key
    # Enable the HTTP-01 challenge provider
    solvers:
      # An empty 'selector' means that this solver matches all domains
      - selector: {}
        http01:
          ingress:
            class: {{.IngressClass}}`
