// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallRegistryCredsOperator() *cobra.Command {
	var command = &cobra.Command{
		Use:   "registry-creds",
		Short: "Install registry-creds",
		Long: `Install the registry-creds operator, to take a single registry secret and 
to propagate it to all available namespaces. Works on regular Intel, ARM 
and ARM64 clusters.`,
		Example:      `  arkade install registry-creds`,
		SilenceUsage: true,
	}

	command.Flags().String("username", "", "Username for your registry or the Docker Hub")
	command.Flags().String("password", "", "Password for your registry or the Docker Hub")
	command.Flags().String("email", "", "Email address for your registry or the Docker Hub (optional)")
	command.Flags().String("server", "", "Server for your registry or the Docker Hub, default: is blank, for the Docker Hub")
	command.Flags().Bool("from-env", false, "Read flags from the environment instead of flags, prefixed with DOCKER_, i.e. DOCKER_EMAIL")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		var (
			username string
			password string
			email    string
			server   string
			fromEnv  bool
		)

		fmt.Printf("Applying controller's manifests.\n")
		_, err := k8s.KubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/alexellis/registry-creds/master/mainfest.yaml")
		if err != nil {
			return err
		}

		if command.Flags().Changed("username") {
			var err error
			username, err = command.Flags().GetString("username")
			if err != nil {
				return err
			}
			password, err = command.Flags().GetString("password")
			if err != nil {
				return err
			}
			email, err = command.Flags().GetString("email")
			if err != nil {
				return err
			}
			server, err = command.Flags().GetString("server")
			if err != nil {
				return err
			}
		}
		if fromEnv {
			username = os.Getenv("DOCKER_USERNAME")
			password = os.Getenv("DOCKER_PASSWORD")
			email = os.Getenv("DOCKER_EMAIL")
			server = os.Getenv("DOCKER_SERVER")
		}

		if len(username) > 0 && len(password) == 0 {
			return fmt.Errorf("both a username, and password are required when a username is given")
		}

		if len(username) > 0 {
			fmt.Printf("Attempting to create secret for user: %s\n", username)
			serverStr := ""
			if len(server) > 0 {
				serverStr = "--docker-server=" + server
			}
			res, err := k8s.KubectlTask("create", "secret", "docker-registry", "registry-seed-secret",
				"--namespace=kube-system",
				"--docker-username="+username, "--docker-password="+password,
				"--docker-email="+email, serverStr)
			if err != nil {
				return err
			}

			if res.ExitCode != 0 {
				fmt.Printf("Warning: %s\n", res.Stderr)
			} else {
				fmt.Printf("%s", res.Stdout)
			}

			dir := os.TempDir()
			cr := path.Join(dir, "clusterpullsecret.yaml")

			err = ioutil.WriteFile(cr, []byte(`apiVersion: ops.alexellis.io/v1
kind: ClusterPullSecret
metadata:
  name: primary
spec:
  secretRef:
    name: registry-seed-secret
    namespace: kube-system
`), 0644)

			if err != nil {
				return err
			}
			fmt.Printf("Wrote temporary file: %s\n", cr)

			fmt.Printf("Creating ClusterPullSecret\n")
			crRes, err := k8s.KubectlTask("apply", "-f", cr)
			if err != nil {
				return err
			}
			if crRes.ExitCode != 0 {
				fmt.Printf("Warning: %s\n", crRes.Stderr)
			}

			fmt.Printf(`
# To view your ClusterPullSecret
kubectl get ClusterPullSecret/primary

# View your Seed Secret
kubectl get secret -n kube-system registry-seed-secret

# View your ServiceAccount's imagePullSecrets list
kubectl get serviceaccount default -o yaml

`)
		}

		fmt.Println(`=======================================================================
= registry-creds has been installed.                                  =
=======================================================================` +
			"\n\n" + RegistryCredsOperatorInfoMsg + "\n\n" + pkg.ThanksForUsing)

		return nil
	}

	return command
}

const RegistryCredsOperatorInfoMsg = `This operator is used to propagate a single ImagePullSecret to all 
namespaces within your cluster, so that images can be pulled with 
authentication.

Why is that required? For private images, and to cope with the 
Docker Hub's recent rate limiting for anonymous pulls of layers.

# Find out more at
https://github.com/alexellis/registry-creds`
