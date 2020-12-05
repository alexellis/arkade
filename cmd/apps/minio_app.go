// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

func MakeInstallMinio() *cobra.Command {
	var minio = &cobra.Command{
		Use:          "minio",
		Short:        "Install minio",
		Long:         `Install minio`,
		Example:      `  arkade install minio`,
		SilenceUsage: true,
	}

	minio.Flags().Bool("update-repo", true, "Update the helm repo")
	minio.Flags().String("access-key", "", "Provide an access key to override the pre-generated value")
	minio.Flags().String("secret-key", "", "Provide a secret key to override the pre-generated value")
	minio.Flags().Bool("distributed", false, "Deploy Minio in Distributed Mode")
	minio.Flags().String("namespace", "default", "Kubernetes namespace for the application")
	minio.Flags().Bool("persistence", false, "Enable persistence")
	minio.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	minio.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetBool("wait")
		if err != nil {
			return fmt.Errorf("error with --wait usage: %s", err)
		}

		_, err = command.Flags().GetBool("persistence")
		if err != nil {
			return fmt.Errorf("error with --persistence usage: %s", err)
		}

		_, err = command.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error with --namespace usage: %s", err)
		}

		_, err = command.Flags().GetBool("update-repo")
		if err != nil {
			return fmt.Errorf("error with --update-repo usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		_, err = command.Flags().GetString("access-key")
		if err != nil {
			return fmt.Errorf("error with --access-key usage: %s", err)
		}
		_, err = command.Flags().GetString("secret-key")
		if err != nil {
			return fmt.Errorf("error with --secret-key usage: %s", err)
		}

		_, err = command.Flags().GetBool("distributed")
		if err != nil {
			return fmt.Errorf("error with --distributed usage: %s", err)
		}

		return nil
	}

	minio.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		wait, _ := command.Flags().GetBool("wait")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		ns, _ := command.Flags().GetString("namespace")
		persistence, _ := command.Flags().GetBool("persistence")
		accessKey, _ := command.Flags().GetString("access-key")
		secretKey, _ := command.Flags().GetString("secret-key")
		dist, _ := command.Flags().GetBool("distributed")
		customFlags, _ := command.Flags().GetStringArray("set")

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		overrides := map[string]string{}

		gen, err := password.NewGenerator(&password.GeneratorInput{
			Symbols: "+/",
		})
		if err != nil {
			return err
		}

		if len(accessKey) == 0 {
			fmt.Printf("Access Key not provided, one will be generated for you\n")
			accessKey, err = gen.Generate(20, 10, 0, false, true)
		}
		if len(secretKey) == 0 {
			fmt.Printf("Secret Key not provided, one will be generated for you\n")
			secretKey, err = gen.Generate(40, 10, 5, false, true)
		}

		overrides["accessKey"] = accessKey
		overrides["secretKey"] = secretKey
		overrides["persistence.enabled"] = strings.ToLower(strconv.FormatBool(persistence))
		if dist {
			overrides["mode"] = "distributed"
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		minioAppOptions := types.DefaultInstallOptions().
			WithNamespace(ns).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("minio/minio").
			WithHelmURL("https://helm.min.io/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(minioAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(minioInstallMsg)
		return nil
	}

	return minio
}

var _, clientOS = env.GetClientArch()

var MinioInfoMsg = `# Forward the minio port to your machine
kubectl port-forward -n default svc/minio 9000:9000 &

# Get the access and secret key to gain access to minio
ACCESSKEY=$(kubectl get secret -n default minio -o jsonpath="{.data.accesskey}" | base64 --decode; echo)
SECRETKEY=$(kubectl get secret -n default minio -o jsonpath="{.data.secretkey}" | base64 --decode; echo)

# Get the Minio Client
curl -SLf https://dl.min.io/client/mc/release/` + strings.ToLower(clientOS) + `-amd64/mc \
  && chmod +x mc

# Add a host
mc config host add minio http://127.0.0.1:9000 $ACCESSKEY $SECRETKEY

# List buckets
mc ls minio

# Find out more at: https://min.io`

var minioInstallMsg = `=======================================================================
= Minio has been installed.                                           =
=======================================================================` +
	"\n\n" + MinioInfoMsg + "\n\n" + pkg.ThanksForUsing
