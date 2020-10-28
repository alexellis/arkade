package cmd_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alexellis/arkade/cmd"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/k3s"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestInstall(t *testing.T) {
	c, err := gnomock.Start(
		k3s.Preset(k3s.WithVersion("v1.19.3")),
		gnomock.WithContainerName("gnomock-k3s"),
	)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, gnomock.Stop(c))
	}()

	cfgBytes, err := k3s.ConfigBytes(c)
	require.NoError(t, err)

	f, err := ioutil.TempFile("", "gnomock-kubeconfig-")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, f.Close())
		require.NoError(t, os.Remove(f.Name()))
	}()

	_, err = f.Write(cfgBytes)
	require.NoError(t, err)

	require.NoError(t, os.Setenv("KUBECONFIG", f.Name()))
	ctx := context.Background()

	cfg, err := k3s.Config(c)
	require.NoError(t, err)

	client, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	t.Run("install cli tools", func(t *testing.T) {
		command := cmd.MakeGet()
		command.SetArgs([]string{"kubectl"})
		require.NoError(t, command.Execute())

		home, err := os.UserHomeDir()
		require.NoError(t, err)
		path := os.Getenv("PATH")
		os.Setenv("PATH", fmt.Sprintf("%s:%s/.arkade/bin", path, home))
	})

	t.Run("openfaas", func(t *testing.T) {
		command := cmd.MakeInstall()
		command.SetArgs([]string{"openfaas"})
		require.NoError(t, command.Execute())

		deploys, err := client.AppsV1().Deployments("openfaas").List(ctx, v1.ListOptions{})
		require.NoError(t, err)
		require.Len(t, deploys.Items, 7)

		actualDeploys := make([]string, 0, 7)
		for _, deploy := range deploys.Items {
			actualDeploys = append(actualDeploys, deploy.Name)
		}
		expectedDeploys := []string{
			"alertmanager", "nats", "queue-worker", "basic-auth-plugin",
			"prometheus", "gateway", "faas-idler",
		}
		require.ElementsMatch(t, expectedDeploys, actualDeploys)
	})
}
