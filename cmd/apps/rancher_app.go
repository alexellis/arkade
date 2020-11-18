// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallRancher() *cobra.Command {
	var rancher = &cobra.Command{
		Use:          "rancher",
		Short:        "Install rancher",
		Long:         "Install rancher to manage kubernetes clusters. Requires cert-manager.",
		Example:      "arkade install rancher --hostname mycluster.example.com ",
		SilenceUsage: true,
	}

	rancher.Flags().StringP("namespace", "n", "cattle-system", "The namespace to install rancher")
	rancher.Flags().Bool("update-repo", true, "Update the helm repo")
	rancher.Flags().String("hostname", "", "public DNS record to access rancher")
	rancher.Flags().String("ingress-source", "letsEncrypt", "source of TLS cert: letsEncrypt, rancher, secret")
	rancher.Flags().String("letsencrypt-email", "", "email address used for communication about your certificate")
	rancher.Flags().Int("replicas", 3, "Number of replicas of rancher pods")
	rancher.Flags().String("tls-cert", "", "Path to PEM encoded public key certificate")
	rancher.Flags().String("tls-key", "", "Path to private key associated with given certificate")
	rancher.Flags().String("tls-ca", "", "Path to private CA certificate")

	rancher.RunE = func(command *cobra.Command, args []string) error {
		// Get all flags
		namespace, _ := command.Flags().GetString("namespace")
		hostname, _ := command.Flags().GetString("hostname")
		ingressSource, _ := command.Flags().GetString("ingress-source")
		letsencryptEmail, _ := command.Flags().GetString("letsencrypt-email")
		replicas, _ := command.Flags().GetInt("replicas")
		tlsCert, _ := command.Flags().GetString("tls-cert")
		tlsKey, _ := command.Flags().GetString("tls-key")
		tlsCa, _ := command.Flags().GetString("tls-ca")

		// Flags validations
		if hostname == "" {
			return fmt.Errorf(`--hostname is required by rancher`)
		}

		if ingressSource != "rancher" && ingressSource != "letsEncrypt" && ingressSource != "secret" {
			return fmt.Errorf(`--ingress-source only accepts one of: letsEncrypt, rancher, secret`)
		}

		if ingressSource == "letsEncrypt" && letsencryptEmail == "" {
			return fmt.Errorf(`If you are using Let's Encrypt --letsencrypt-email is required`)
		}

		if (tlsCert == "" && tlsKey != "") || (tlsCert != "" && tlsKey == "") {
			return fmt.Errorf(`To create the TLS secret you need to give both flags: --tls-cert and --tls-key`)
		}

		// initialize client env
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		log.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		// exit on arm
		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(`ARM is experimental and not supported by Rancher https://rancher.com/docs/rancher/v2.x/en/installation/options/arm64-platform/`)
		}

		// create the namespace
		nsRes, nsErr := k8s.KubectlTask("create", "namespace", namespace)
		if nsErr != nil {
			return nsErr
		}

		// ignore errors
		if nsRes.ExitCode != 0 {
			log.Printf("[Warning] unable to create namespace %s, may already exist: %s", namespace, nsRes.Stderr)
		}

		// Deploy TLS cert secret if given
		if tlsCert != "" || tlsKey != "" {
			certOption := fmt.Sprintf("--cert=%s", tlsCert)
			keyOption := fmt.Sprintf("--key=%s", tlsKey)
			tlsRes, tlsErr := k8s.KubectlTask("--namespace", namespace, "create", "secret", "tls", "tls-rancher-ingress", certOption, keyOption)
			if tlsErr != nil {
				return tlsErr
			}

			if tlsRes.ExitCode != 0 {
				log.Printf("[Warning] unable to create tls secret for rancher, you will have to do it manually later")
			}
		}

		if tlsCa != "" {
			caOption := fmt.Sprintf("--from-file=cacerts.pem=%s", tlsCa)

			caRes, caErr := k8s.KubectlTask("--namespace", namespace, "create", "secret", "generic", "tls-ca", caOption)
			if caErr != nil {
				return caErr
			}

			if caRes.ExitCode != 0 {
				log.Printf("[Warning] unable to create tls-ca secret for rancher, you will have to do it manually later")
			}
		}

		// set overrides for chart install
		overrides := map[string]string{}
		overrides["hostname"] = hostname
		overrides["replicas"] = fmt.Sprintf("%d", replicas)
		overrides["ingress.tls.source"] = ingressSource

		if ingressSource == "letsEncrypt" {
			overrides["letsEncrypt.email"] = letsencryptEmail
		}

		if tlsCa != "" {
			overrides["privateCA"] = "true"
		}

		// install the chart
		helmHome := path.Join(userPath, ".helm")
		updateRepo, _ := rancher.Flags().GetBool("update-repo")

		rancherOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(helmHome).
			WithHelmRepo("rancher-stable/rancher").
			WithHelmURL("https://releases.rancher.com/server-charts/stable").
			WithHelmUpdateRepo(updateRepo).
			WithOverrides(overrides)

		if command.Flags().Changed("kubeconfig") {
			kubeconfigPath, _ := command.Flags().GetString("kubeconfig")
			rancherOptions.WithKubeconfigPath(kubeconfigPath)
		}

		os.Setenv("HELM_HOME", helmHome)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(rancherOptions)
		if err != nil {
			return err
		}

		fmt.Println(rancherInstallMsg)

		return nil
	}

	return rancher
}

const RancherInfoMsg = `
# Get started with Rancher here:
# https://rancher.com/docs/rancher/v2.x/en/quick-start-guide/

# Wait for Rancher to be rolled out:

	kubectl -n cattle-system rollout status deploy/rancher
`

var rancherInstallMsg = `=======================================================================
=                      rancher has been installed                     =
=======================================================================` +
	"\n\n" + RancherInfoMsg + "\n\n" + pkg.ThanksForUsing
