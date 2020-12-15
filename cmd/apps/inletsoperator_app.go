// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallInletsOperator() *cobra.Command {
	var inletsOperator = &cobra.Command{
		Use:          "inlets-operator",
		Short:        "Install inlets-operator",
		Long:         `Install inlets-operator to get public IPs for your cluster`,
		Example:      `  arkade install inlets-operator --namespace default`,
		SilenceUsage: true,
	}

	inletsOperator.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	inletsOperator.Flags().StringP("license", "l", "", "The license key if using inlets-pro")
	inletsOperator.Flags().StringP("license-file", "f", "", "Text file containing license key, used for inlets-pro")
	inletsOperator.Flags().StringP("provider", "p", "digitalocean", "Your infrastructure provider - 'equinix-metal', 'digitalocean', 'scaleway', 'linode', 'civo', 'gce', 'ec2', 'azure'")
	inletsOperator.Flags().StringP("zone", "z", "us-central1-a", "The zone to provision the exit node (GCE)")
	inletsOperator.Flags().String("project-id", "", "Project ID to be used (for GCE and Equinix Metal)")
	inletsOperator.Flags().StringP("region", "r", "lon1", "The default region to provision the exit node (DigitalOcean, Equinix Metal and Scaleway)")
	inletsOperator.Flags().String("organization-id", "", "The organization id (Scaleway)")
	inletsOperator.Flags().String("subscription-id", "", "The subscription id (Azure)")
	inletsOperator.Flags().StringP("token-file", "t", "", "Text file containing token or a service account JSON file")
	inletsOperator.Flags().StringP("token", "k", "", "The API access token")
	inletsOperator.Flags().StringP("secret-key-file", "s", "", "Text file containing secret key, used for providers like ec2")
	inletsOperator.Flags().Bool("update-repo", true, "Update the helm repo")

	inletsOperator.Flags().String("pro-client-image", "", "Docker image for inlets-pro's client")
	inletsOperator.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	inletsOperator.PreRunE = func(command *cobra.Command, args []string) error {
		tokenFileName, _ := command.Flags().GetString("token-file")
		tokenString, _ := command.Flags().GetString("token")

		if len(tokenFileName) > 0 && len(tokenString) > 0 {
			return fmt.Errorf(`--token-file or --access-key is a required field for your cloud API token or service account JSON`)
		}
		if len(tokenFileName) > 0 {
			if _, err := os.Stat(tokenFileName); err != nil {
				return err
			}
		}

		secretKeyFile, _ := command.Flags().GetString("secret-key-file")
		if len(secretKeyFile) > 0 {
			if _, err := os.Stat(secretKeyFile); err != nil {
				return err
			}
		}

		return nil
	}

	inletsOperator.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		wait, _ := command.Flags().GetBool("wait")

		namespace, _ := command.Flags().GetString("namespace")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)
		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		updateRepo, _ := inletsOperator.Flags().GetBool("update-repo")
		err = helm.AddHelmRepo("inlets", "https://inlets.github.io/inlets-operator/", updateRepo)
		if err != nil {
			return err
		}

		err = helm.FetchChart("inlets/inlets-operator", defaultVersion)
		if err != nil {
			return err
		}
		overrides, err := getInletsOperatorOverrides(command)

		if err != nil {
			return err
		}

		_, err = k8s.KubectlTask("apply", "-f", "https://raw.githubusercontent.com/inlets/inlets-operator/master/artifacts/crds/inlets.inlets.dev_tunnels.yaml")
		if err != nil {
			return err
		}

		tokenFileName, _ := command.Flags().GetString("token-file")
		tokenString, _ := command.Flags().GetString("token")

		s := Secret{
			Namespace: namespace,
			Name:      "inlets-access-key",
		}
		if len(tokenFileName) > 0 {
			s.Literals = append(s.Literals, SecretLiteral{
				Name:     "inlets-access-key",
				FromFile: tokenFileName,
			})
		} else {
			s.Literals = append(s.Literals, SecretLiteral{
				Name:      "inlets-access-key",
				FromValue: tokenString,
			})
		}

		err = applySecret(s)
		if err != nil {
			return err
		}

		secretKeyFile, _ := command.Flags().GetString("secret-key-file")
		if len(secretKeyFile) > 0 {
			s := Secret{
				Namespace: namespace,
				Name:      "inlets-secret-key",
			}
			s.Literals = append(s.Literals, SecretLiteral{
				Name:     "inlets-access-key",
				FromFile: secretKeyFile,
			})

			err = applySecret(s)
			if err != nil {
				return err
			}
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		region, _ := command.Flags().GetString("region")
		overrides["region"] = region

		if val, _ := command.Flags().GetString("license"); len(val) > 0 {
			overrides["inletsProLicense"] = val
		}

		if licenseFile, _ := command.Flags().GetString("license-file"); len(licenseFile) > 0 {
			licenseKey, err := ioutil.ReadFile(licenseFile)
			if err != nil {
				return err
			}

			overrides["inletsProLicense"] = strings.TrimSpace(string(licenseKey))
		}

		if val, _ := command.Flags().GetString("pro-client-image"); len(val) > 0 {
			overrides["proClientImage"] = val
		}

		err = helm.Helm3Upgrade("inlets/inlets-operator",
			namespace, "values.yaml", "", overrides, wait)
		if err != nil {
			return err
		}

		fmt.Println(inletsOperatorPostInstallMsg)

		return nil
	}

	return inletsOperator
}

func getInletsOperatorOverrides(command *cobra.Command) (map[string]string, error) {
	overrides := map[string]string{}
	provider, _ := command.Flags().GetString("provider")
	overrides["provider"] = strings.ToLower(provider)

	secretKeyFile, _ := command.Flags().GetString("secret-key-file")

	if len(secretKeyFile) > 0 {
		overrides["secretKeyFile"] = "/var/secrets/inlets/secret/inlets-secret-key"
	}

	providers := []string{
		"digitalocean", "equinix-metal", "ec2", "scaleway", "civo", "gce", "linode", "azure",
	}

	found := false
	for _, p := range providers {
		if p == provider {
			found = true
		}
	}
	if !found {
		return overrides, fmt.Errorf("provider: %s not supported at this time", provider)
	}

	if provider == "gce" {
		gceProjectID, err := command.Flags().GetString("project-id")
		if err != nil {
			return overrides, err
		}
		overrides["projectID"] = gceProjectID

		zone, err := command.Flags().GetString("zone")
		if err != nil {
			return overrides, err
		}
		overrides["zone"] = strings.ToLower(zone)

		if len(zone) == 0 {
			return overrides, fmt.Errorf("zone is required for provider %s", provider)
		}

		if len(gceProjectID) == 0 {
			return overrides, fmt.Errorf("project-id is required for provider %s", provider)
		}
	} else if provider == "equinix-metal" {
		equinixMetalProjectID, err := command.Flags().GetString("project-id")
		if err != nil {
			return overrides, err
		}
		overrides["projectID"] = equinixMetalProjectID

		if len(equinixMetalProjectID) == 0 {
			return overrides, fmt.Errorf("project-id is required for provider %s", provider)
		}

	} else if provider == "scaleway" {
		orgID, err := command.Flags().GetString("organization-id")
		if err != nil {
			return overrides, err
		}
		overrides["organizationID"] = orgID

		if len(secretKeyFile) == 0 {
			return overrides, fmt.Errorf("secret-key-file is required for provider %s", provider)
		}

		if len(orgID) == 0 {
			return overrides, fmt.Errorf("organization-id is required for provider %s", provider)
		}
	} else if provider == "azure" {
		subscriptionID, err := command.Flags().GetString("subscription-id")
		if err != nil {
			return overrides, err
		}
		if len(subscriptionID) == 0 {
			return overrides, fmt.Errorf("subscription-id is required for provider %s", provider)
		}
		overrides["subscriptionID"] = subscriptionID
	} else if provider == "ec2" {
		if len(secretKeyFile) == 0 {
			return overrides, fmt.Errorf("secret-key-file is required for provider %s", provider)
		}
	}

	return overrides, nil
}

const InletsOperatorInfoMsg = `# The default configuration is for DigitalOcean and your secret is
# stored as "inlets-access-key" in the "default" namespace or the namespace 
# you gave if installing with helm3

# To get your first Public IP run the following:

# K8s 1.17
kubectl run nginx-1 --image=nginx --port=80 --restart=Always

# K8s 1.18 and higher:

kubectl apply -f \
 https://raw.githubusercontent.com/inlets/inlets-operator/master/contrib/nginx-sample-deployment.yaml

# Then expose the Deployment as a LoadBalancer:

kubectl expose deployment nginx-1 --port=80 --type=LoadBalancer

# Find your IP in the "EXTERNAL-IP" field, watch for "<pending>" to 
# change to an IP

kubectl get svc -w

# When you're done, remove the tunnel by deleting the service
kubectl delete svc/nginx-1

# Check the logs
kubectl logs deploy/inlets-operator -f

# Find out more at:
# https://github.com/inlets/inlets-operator`

const inletsOperatorPostInstallMsg = `=======================================================================
= inlets-operator has been installed.                                  =
=======================================================================` +
	"\n\n" + InletsOperatorInfoMsg + "\n\n" + pkg.ThanksForUsing

type Secret struct {
	Namespace string
	Name      string
	Stdin     io.Reader
	Literals  []SecretLiteral
}

type SecretLiteral struct {
	Name      string
	FromFile  string
	FromValue string
}

type K8sVer struct {
	ClientVersion struct {
		Major string `json:"major"`
		Minor string `json:"minor"`
	} `json:"clientVersion"`
}

func getKubectlVersion() (int, int, error) {
	res, err := k8s.KubectlTask("version", "--client", "-o=json")
	if err != nil {
		return 0, 0, err
	}

	if err != nil {
		return 0, 0, err
	}

	v := K8sVer{}
	if err := json.Unmarshal([]byte(res.Stdout), &v); err != nil {
		return 0, 0, err
	}

	major, _ := strconv.Atoi(v.ClientVersion.Major)
	minor, _ := strconv.Atoi(v.ClientVersion.Minor)
	return major, minor, nil

}
func applySecret(s Secret) error {

	_, minor, err := getKubectlVersion()
	if err != nil {
		return err
	}
	dryRunSuffix := ""
	if minor >= 19 {
		dryRunSuffix = "=client"
	}

	parts := []string{"create", "secret", "generic", s.Name, "--dry-run" + dryRunSuffix, "-o=yaml"}

	for _, l := range s.Literals {
		if len(l.FromFile) > 0 {
			parts = append(parts, "--from-file", s.Name+"="+l.FromFile)
		} else {
			parts = append(parts, "--from-literal", s.Name+"="+l.FromValue)
		}
	}

	res, err := k8s.KubectlTask(parts...)
	if err != nil {
		return err
	} else if len(res.Stderr) > 0 && strings.Contains(res.Stderr, "Warning") == false {
		return fmt.Errorf("error from kubectl\n%q", res.Stderr)
	}

	manifest := bytes.NewReader([]byte(res.Stdout))
	res, err = k8s.KubectlTaskStdin(manifest, "apply", "-f", "-")

	if err != nil {
		return err
	} else if len(res.Stderr) > 0 && strings.Contains(res.Stderr, "Warning") == false {
		return fmt.Errorf("error from kubectl\n%q", res.Stderr)
	}
	return nil
}
