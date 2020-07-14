// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallNfsProvisioner() *cobra.Command {
	var nfsProvisionerApp = &cobra.Command{
		Use:          "nfs-client-provisioner",
		Short:        "Install nfs client provisioner",
		Long:         "Install nfs client provisioner to create dynamic persistent volumes",
		Example:      "arkade install nfs-client-provisioner --nfs-server=x.x.x.x --nfs-path=/exported/path",
		SilenceUsage: true,
	}

	nfsProvisionerApp.Flags().StringP("namespace", "n", "default", "The namespace to install nfs-client (default: default")
	nfsProvisionerApp.Flags().String("nfs-server", "", "IP or hostname of the NFS server ")
	nfsProvisionerApp.Flags().String("nfs-path", "", "Basepath of the mount point to be used")
	nfsProvisionerApp.Flags().Bool("update-repo", true, "Update the helm repo")
	nfsProvisionerApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set =true)")

	nfsProvisionerApp.RunE = func(command *cobra.Command, args []string) error {
		helm3 := true

		namespace, _ := nfsProvisionerApp.Flags().GetString("namespace")
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		log.Printf("Client: %s, %s\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		if err := os.Setenv("HELM_HOME", path.Join(userPath, ".helm")); err != nil {
			return err
		}

		nfsServer, _ := command.Flags().GetString("nfs-server")
		nfsPath, _ := command.Flags().GetString("nfs-path")

		if len(nfsServer) == 0 {
			return fmt.Errorf("--nfs-server required")
		}

		if len(nfsPath) == 0 {
			return fmt.Errorf("--nfs-path required")
		}

		overrides := map[string]string{}

		overrides["nfs.server"] = nfsServer
		overrides["nfs.path"] = nfsPath

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if suffix := getValuesSuffix(arch); suffix == "-armhf" || suffix == "-arm64" {
			overrides["image.repository"] = "quay.io/external_storage/nfs-client-provisioner-arm:latest"
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		nfsProvisionerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("stable/nfs-client-provisioner").
			WithHelmURL("https://kubernetes-charts.storage.googleapis.com").
			WithOverrides(overrides)

		if command.Flags().Changed("kubeconfig") {
			kubeconfigPath, _ := command.Flags().GetString("kubeconfig")
			nfsProvisionerOptions.WithKubeconfigPath(kubeconfigPath)
		}

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(nfsProvisionerOptions)
		if err != nil {
			return err
		}

		println(nfsClientInstallMsg)
		return nil
	}

	return nfsProvisionerApp
}

const NfsClientProvisioneriInfoMsg = `# Test your NFS provisioner:

cat <<EOF | kubectl apply -f -
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: test-claim
  annotations:
    volume.beta.kubernetes.io/storage-class: "nfs-client"
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
EOF

# Create pod:

cat <<EOF | kubectl apply -f -
kind: Pod
apiVersion: v1
metadata:
  name: test-pod
spec:
  containers:
  - name: test-pod
    image: gcr.io/google_containers/busybox:1.24
    command:
      - "/bin/sh"
    args:
      - "-c"
      - "touch /mnt/SUCCESS && exit 0 || exit 1"
    volumeMounts:
      - name: nfs-pvc
        mountPath: "/mnt"
  restartPolicy: "Never"
  volumes:
    - name: nfs-pvc
      persistentVolumeClaim:
        claimName: test-claim
EOF

# Now check your NFS Server for the file SUCCESS.

kubectl delete -f deploy/test-pod.yaml -f deploy/test-claim.yaml

# Now check the folder has been deleted.



`

const nfsClientInstallMsg = `=======================================================================
= nfs-client-provisioner has been installed.                                   =
=======================================================================` +
	"\n\n" + NfsClientProvisioneriInfoMsg + "\n\n" + pkg.ThanksForUsing
