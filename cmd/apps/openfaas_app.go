// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/sethvargo/go-password/password"

	"github.com/spf13/cobra"
)

func MakeInstallOpenFaaS() *cobra.Command {
	var openfaas = &cobra.Command{
		Use:          "openfaas",
		Short:        "Install openfaas",
		Long:         `Install openfaas`,
		Example:      `  arkade install openfaas --load-balancer`,
		SilenceUsage: true,
	}

	openfaas.Flags().BoolP("basic-auth", "a", true, "Enable authentication")
	openfaas.Flags().String("basic-auth-password", "", "Overide the default random basic-auth-password if this is set")
	openfaas.Flags().BoolP("load-balancer", "l", false, "Add a loadbalancer")
	openfaas.Flags().StringP("namespace", "n", "openfaas", "The namespace for the core services")
	openfaas.Flags().Bool("update-repo", true, "Update the helm repo")
	openfaas.Flags().String("pull-policy", "IfNotPresent", "Pull policy for OpenFaaS core services")
	openfaas.Flags().String("function-pull-policy", "Always", "Pull policy for functions")

	openfaas.Flags().Bool("operator", false, "Create OpenFaaS Operator")
	openfaas.Flags().Bool("clusterrole", false, "Create a ClusterRole for OpenFaaS instead of a limited scope Role")
	openfaas.Flags().Bool("direct-functions", false, "Invoke functions directly from the gateway, or load-balance via endpoint IPs when set to false")
	openfaas.Flags().Bool("autoscaler", false, "Deploy OpenFaaS with the autoscaler enabled")
	openfaas.Flags().Bool("jetstream", false, "Deploy OpenFaaS with jetstream queue mode")
	openfaas.Flags().Bool("dashboard", false, "Deploy OpenFaaS with the dashboard enabled")

	openfaas.Flags().Int("queue-workers", 1, "Replicas of queue-worker for HA")
	openfaas.Flags().Int("max-inflight", 1, "Max tasks for queue-worker to process in parallel")
	openfaas.Flags().Int("gateways", 1, "Replicas of gateway")

	openfaas.Flags().Bool("ingress-operator", false, "Get custom domains and Ingress records via the ingress-operator component")

	openfaas.Flags().String("log-provider-url", "", "Set a log provider url for OpenFaaS")

	openfaas.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set gateway.replicas=2)")

	openfaas.Flags().String("license-file", "", "Path to OpenFaaS Pro license file")

	openfaas.RunE = func(command *cobra.Command, args []string) error {
		appOpts := types.DefaultInstallOptions()

		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")
		basicAuthEnabled, _ := command.Flags().GetBool("basic-auth")
		updateRepo, _ := openfaas.Flags().GetBool("update-repo")
		logUrl, _ := command.Flags().GetString("log-provider-url")
		licenseFile, _ := command.Flags().GetString("license-file")
		pullPolicy, _ := command.Flags().GetString("pull-policy")
		functionPullPolicy, _ := command.Flags().GetString("function-pull-policy")
		createOperator, _ := command.Flags().GetBool("operator")
		clusterRole, _ := command.Flags().GetBool("clusterrole")
		directFunctions, _ := command.Flags().GetBool("direct-functions")
		autoscaler, _ := command.Flags().GetBool("autoscaler")
		jetstream, _ := command.Flags().GetBool("jetstream")
		dashboard, _ := command.Flags().GetBool("dashboard")
		gateways, _ := command.Flags().GetInt("gateways")
		maxInflight, _ := command.Flags().GetInt("max-inflight")
		queueWorkers, _ := command.Flags().GetInt("queue-workers")
		ingressOperator, _ := command.Flags().GetBool("ingress-operator")
		lb, _ := command.Flags().GetBool("load-balancer")

		overrides := map[string]string{}

		arch := k8s.GetNodeArchitecture()
		valuesSuffix := getValuesSuffix(arch)

		_, err := k8s.KubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml")

		if err != nil {
			return err
		}

		if basicAuthEnabled {
			pass, _ := command.Flags().GetString("basic-auth-password")

			if len(pass) == 0 {
				var err error
				pass, err = password.Generate(25, 10, 0, false, true)
				if err != nil {
					return err
				}
			}
			secretData := []types.SecretsData{
				{Type: types.StringLiteralSecret, Key: "basic-auth-user", Value: "admin"},
				{Type: types.StringLiteralSecret, Key: "basic-auth-password", Value: pass},
			}

			basicAuthSecret := types.NewGenericSecret("basic-auth", namespace, secretData)
			appOpts.WithSecret(basicAuthSecret)
		}

		if logUrl != "" {
			overrides["gateway.logsProviderURL"] = logUrl
		}

		// If license file is sent, then we assume to set the --pro flag and create the secret
		if len(licenseFile) != 0 {
			overrides["openfaasPro"] = "true"
			secretData := []types.SecretsData{
				{Type: types.FromFileSecret, Key: "license", Value: licenseFile},
			}

			proLicense := types.NewGenericSecret("openfaas-license", namespace, secretData)
			appOpts.WithSecret(proLicense)
		}

		if dashboard {
			privateKey, publicKey, err := generateJWTKeyPair()
			if err != nil {
				return fmt.Errorf("failed to create JWT key-pair: %s", err)
			}

			secretData := []types.SecretsData{
				{Type: types.StringLiteralSecret, Key: "key", Value: string(privateKey)},
				{Type: types.StringLiteralSecret, Key: "key.pub", Value: string(publicKey)},
			}

			dashboardJWT := types.NewGenericSecret("dashboard-jwt", namespace, secretData)
			appOpts.WithSecret(dashboardJWT)
		}

		overrides["clusterRole"] = strconv.FormatBool(clusterRole)
		overrides["gateway.directFunctions"] = strconv.FormatBool(directFunctions)
		overrides["operator.create"] = strconv.FormatBool(createOperator)
		overrides["openfaasImagePullPolicy"] = pullPolicy
		overrides["faasnetes.imagePullPolicy"] = functionPullPolicy
		overrides["basicAuthPlugin.replicas"] = "1"
		overrides["gateway.replicas"] = fmt.Sprintf("%d", gateways)
		overrides["ingressOperator.create"] = strconv.FormatBool(ingressOperator)
		overrides["queueWorker.replicas"] = fmt.Sprintf("%d", queueWorkers)
		overrides["queueWorker.maxInflight"] = fmt.Sprintf("%d", maxInflight)
		overrides["autoscaler.enabled"] = strconv.FormatBool(autoscaler)
		overrides["dashboard.enabled"] = strconv.FormatBool(dashboard)
		overrides["dashboard.publicURL"] = "http://127.0.0.1:8080"

		// the value in the template is "basic_auth" not the more usual basicAuth
		overrides["basic_auth"] = strconv.FormatBool(basicAuthEnabled)

		overrides["serviceType"] = "NodePort"

		if lb {
			overrides["serviceType"] = "LoadBalancer"
		}

		if jetstream {
			overrides["queueMode"] = "jetstream"
		}

		customFlags, _ := command.Flags().GetStringArray("set")
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		appOpts.
			WithKubeconfigPath(kubeConfigPath).
			WithOverrides(overrides).
			WithValuesFile(fmt.Sprintf("values%s.yaml", valuesSuffix)).
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithHelmRepo("openfaas/openfaas").
			WithHelmUpdateRepo(updateRepo).
			WithNamespace(namespace).
			WithInstallNamespace(false).
			WithWait(wait)

		if _, err := apps.MakeInstallChart(appOpts); err != nil {
			return err
		}

		fmt.Println(openfaasPostInstallMsg)

		if basicAuthEnabled == false {
			fmt.Println(
				`Warning: It is not recommended to disable authentication for OpenFaaS.`)
		}
		return nil
	}

	openfaas.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetBool("wait")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("basic-auth")
		if err != nil {
			return err
		}

		_, err = openfaas.Flags().GetBool("update-repo")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("operator")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("clusterrole")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("direct-functions")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("autoscaler")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("jetstream")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("dashboard")
		if err != nil {
			return err
		}

		pullPolicy, _ := cmd.Flags().GetString("pull-policy")
		if len(pullPolicy) == 0 {
			return fmt.Errorf("you must give a value for pull-policy such as IfNotPresent or Always")
		}

		functionPullPolicy, _ := cmd.Flags().GetString("function-pull-policy")
		if len(functionPullPolicy) == 0 {
			return fmt.Errorf("you must give a value for function-pull-policy such as IfNotPresent or Always")
		}

		_, err = cmd.Flags().GetString("license-file")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("ingress-operator")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("load-balancer")
		if err != nil {
			return err
		}

		return nil
	}

	return openfaas
}

func getValuesSuffix(arch string) string {
	var valuesSuffix string
	switch arch {
	case "arm":
		valuesSuffix = "-armhf"

	case "arm64", "aarch64":
		valuesSuffix = "-arm64"

	default:
		valuesSuffix = ""

	}

	return valuesSuffix
}

func generateJWTKeyPair() ([]byte, []byte, error) {
	// Private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	ecder, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}
	privOut := bytes.Buffer{}
	pem.Encode(&privOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: ecder})

	// Public key
	pub := &priv.PublicKey
	pubder, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}
	pubOut := bytes.Buffer{}
	pem.Encode(&pubOut, &pem.Block{Type: "PUBLIC KEY", Bytes: pubder})

	return privOut.Bytes(), pubOut.Bytes(), nil
}

const OpenFaaSInfoMsg = `# Get the faas-cli
curl -SLsf https://cli.openfaas.com | sudo sh

# Forward the gateway to your machine
kubectl rollout status -n openfaas deploy/gateway
kubectl port-forward -n openfaas svc/gateway 8080:8080 &

# If basic auth is enabled, you can now log into your gateway:
PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode; echo)
echo -n $PASSWORD | faas-cli login --username admin --password-stdin

faas-cli store deploy figlet
faas-cli list

# For Raspberry Pi
faas-cli store list \
 --platform armhf

faas-cli store deploy figlet \
 --platform armhf

# Find out more at:
# https://github.com/openfaas/faas`

const openfaasPostInstallMsg = `=======================================================================
= OpenFaaS has been installed.                                        =
=======================================================================` +
	"\n\n" + OpenFaaSInfoMsg + "\n\n" + pkg.SupportMessageShort
