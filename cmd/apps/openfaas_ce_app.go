// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
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

func MakeInstallOpenFaaSCE() *cobra.Command {
	var openfaasCE = &cobra.Command{
		Use:          "openfaas-ce",
		Short:        "Install openfaas-ce",
		Long:         `Install the Community Edition of openfaas`,
		Example:      `  arkade install openfaas-ce`,
		SilenceUsage: true,
	}

	openfaasCE.Flags().BoolP("basic-auth", "a", true, "Enable authentication")
	openfaasCE.Flags().String("basic-auth-password", "", "Overide the default random basic-auth-password if this is set")
	openfaasCE.Flags().BoolP("load-balancer", "l", false, "Add a loadbalancer")
	openfaasCE.Flags().StringP("namespace", "n", "openfaas", "The namespace for the core services")
	openfaasCE.Flags().Bool("update-repo", true, "Update the helm repo")

	openfaasCE.Flags().Int("queue-workers", 1, "Replicas of queue-worker for HA")
	openfaasCE.Flags().Int("max-inflight", 1, "Max tasks for queue-worker to process in parallel")
	openfaasCE.Flags().Int("gateways", 1, "Replicas of gateway")

	openfaasCE.Flags().Bool("ingress-operator", false, "Get custom domains and Ingress records via the ingress-operator component")

	openfaasCE.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set gateway.replicas=2)")

	openfaasCE.RunE = func(command *cobra.Command, args []string) error {
		appOpts := types.DefaultInstallOptions()

		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")
		basicAuthEnabled, _ := command.Flags().GetBool("basic-auth")
		updateRepo, _ := openfaasCE.Flags().GetBool("update-repo")
		gateways, _ := command.Flags().GetInt("gateways")
		queueWorkers, _ := command.Flags().GetInt("queue-workers")
		lb, _ := command.Flags().GetBool("load-balancer")
		maxInflight, _ := command.Flags().GetInt("max-inflight")

		overrides := map[string]string{}

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

		overrides["basicAuthPlugin.replicas"] = "1"
		overrides["gateway.replicas"] = fmt.Sprintf("%d", gateways)
		overrides["queueWorker.replicas"] = fmt.Sprintf("%d", queueWorkers)
		overrides["queueWorker.maxInflight"] = fmt.Sprintf("%d", maxInflight)

		// the value in the template is "basic_auth" not the more usual basicAuth
		overrides["basic_auth"] = strconv.FormatBool(basicAuthEnabled)

		overrides["serviceType"] = "NodePort"

		if lb {
			overrides["serviceType"] = "LoadBalancer"
		}

		customFlags, _ := command.Flags().GetStringArray("set")
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		appOpts.
			WithKubeconfigPath(kubeConfigPath).
			WithOverrides(overrides).
			WithValuesFiles([]string{"values.yaml"}).
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithHelmRepo("openfaas/openfaas").
			WithHelmUpdateRepo(updateRepo).
			WithNamespace(namespace).
			WithInstallNamespace(false).
			WithWait(wait)

		if _, err := apps.MakeInstallChart(appOpts); err != nil {
			return err
		}

		fmt.Println(openfaasCEPostInstallMsg)

		if basicAuthEnabled == false {
			fmt.Println(
				`Warning: It is not recommended to disable authentication for OpenFaaS.`)
		}
		return nil
	}

	openfaasCE.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetBool("wait")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("basic-auth")
		if err != nil {
			return err
		}

		_, err = openfaasCE.Flags().GetBool("update-repo")
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

	return openfaasCE
}

const OpenFaaSCEInfoMsg = `# OpenFaaS CE is licensed for personal use only, or a single 60 day
# evaluation: https://github.com/openfaas/faas/blob/master/EULA.md

# Get the faas-cli
arkade get faas-cli

# Forward the gateway to your machine
kubectl rollout status -n openfaas deploy/gateway
kubectl port-forward -n openfaas svc/gateway 8080:8080 &

# If basic auth is enabled, you can now log into your gateway:
PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode; echo)
echo -n $PASSWORD | faas-cli login --username admin --password-stdin

faas-cli store deploy nodeinfo
faas-cli list

# Find out more at:
# https://docs.openfaas.com/`

const openfaasCEPostInstallMsg = `=======================================================================
= OpenFaaS CE has been installed.                                        =
=======================================================================` +
	"\n\n" + OpenFaaSInfoMsg + "\n\n" + pkg.SupportMessageShort
