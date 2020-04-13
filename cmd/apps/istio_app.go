// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallIstio() *cobra.Command {
	var istio = &cobra.Command{
		Use:          "istio",
		Short:        "Install istio",
		Long:         `Install istio`,
		Example:      `  arkade install istio --loadbalancer`,
		SilenceUsage: true,
	}
	istio.Flags().Bool("update-repo", true, "Update the helm repo")
	istio.Flags().String("namespace", "istio-system", "Namespace for the app")
	istio.Flags().Bool("init", true, "Run the Istio init to add CRDs etc")
	istio.Flags().Bool("helm3", true, "Use Helm 3")

	istio.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set=prometheus.enabled=false)")

	istio.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()
		wait, _ := command.Flags().GetBool("wait")

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		namespace, _ := command.Flags().GetString("namespace")

		if namespace != "istio-system" {
			return fmt.Errorf(`to override the "istio-system" namespace, install Istio via helm manually`)
		}

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		helm3, _ := command.Flags().GetBool("helm3")

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		istioVer := "1.4.5"

		err = addHelmRepo("istio", "https://storage.googleapis.com/istio-release/releases/"+istioVer+"/charts", helm3)
		if err != nil {
			return fmt.Errorf("unable to add repo %s", err)
		}

		updateRepo, _ := istio.Flags().GetBool("update-repo")

		if updateRepo {
			err = updateHelmRepos(helm3)
			if err != nil {
				return fmt.Errorf("unable to update repos %s", err)
			}
		}

		_, err = kubectlTask("create", "ns", "istio-system")

		if err != nil {
			return fmt.Errorf("unable to create namespace %s", err)
		}

		chartPath := path.Join(os.TempDir(), "charts")

		err = fetchChart(chartPath, "istio/istio", defaultVersion, helm3)

		if err != nil {
			return fmt.Errorf("unable fetch chart %s", err)
		}

		overrides := map[string]string{}

		valuesFile, writeErr := writeIstioValues()
		if writeErr != nil {
			return writeErr
		}

		outputPath := path.Join(chartPath, "istio")

		if initIstio, _ := command.Flags().GetBool("init"); initIstio {
			// Waiting for the crds to appear
			err = helm3Upgrade(outputPath, "istio/istio-init", namespace, "", defaultVersion, overrides, true)
			if err != nil {
				return fmt.Errorf("unable to istio-init install chart with helm %s", err)
			}
		}

		fmt.Printf("Waiting for Istio init jobs to create CRDs\n")

		_, err = kubectlTask("wait", "-n", "istio-system", "--for=condition=complete", "job", "--all")
		if err != nil {
			fmt.Printf("error waiting for init jobs")
		}

		fmt.Printf("Giving Istio a few moments to propagate its CRDs.\n")
		time.Sleep(time.Second * 5)

		fmt.Printf("Istio init jobs in completed state or timed-out waiting.\n")

		customFlags, customFlagErr := command.Flags().GetStringArray("set")
		if customFlagErr != nil {
			return fmt.Errorf("error with --set usage: %s", customFlagErr)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		err = helm3Upgrade(outputPath, "istio/istio", namespace, valuesFile, defaultVersion, overrides, wait)
		if err != nil {
			return fmt.Errorf("unable to istio install chart with helm %s", err)
		}

		fmt.Println(istioPostInstallMsg)

		return nil
	}

	return istio
}

const IstioInfoMsg = `# Find out more at:
# https://github.com/istio/`

const istioPostInstallMsg = `=======================================================================
= Istio has been installed.                                        =
=======================================================================` +
	"\n\n" + IstioInfoMsg + "\n\n" + pkg.ThanksForUsing

func writeIstioValues() (string, error) {
	out := `#
# Minimal Istio Configuration taken from https://github.com/weaveworks/flagger

# pilot configuration
pilot:
  enabled: true
  sidecar: true
  resources:
    requests:
      cpu: 10m
      memory: 128Mi

gateways:
  enabled: true
  istio-ingressgateway:
    autoscaleMax: 1

# sidecar-injector webhook configuration
sidecarInjectorWebhook:
  enabled: true

# galley configuration
galley:
  enabled: false

# mixer configuration
mixer:
  policy:
    enabled: false
  telemetry:
    enabled: true
    replicaCount: 1
    autoscaleEnabled: false
  resources:
    requests:
      cpu: 10m
      memory: 128Mi

# addon prometheus configuration
prometheus:
  enabled: true
  scrapeInterval: 5s

# addon jaeger tracing configuration
tracing:
  enabled: false

# Common settings.
global:
  proxy:
    # Resources for the sidecar.
    resources:
      requests:
        cpu: 10m
        memory: 64Mi
      limits:
        cpu: 1000m
        memory: 256Mi
  useMCP: false`

	writeTo := path.Join(os.TempDir(), "istio-values.yaml")
	return writeTo, ioutil.WriteFile(writeTo, []byte(out), 0600)
}
