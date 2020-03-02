package app

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/kubernetes"
)

// HelmApp implements a standard helm app with no special install requirements
type HelmApp struct {
	Namespace       string
	ChartRepository string
	ChartName       string
	ChartVersion    string
	Name            string
	Values          map[string]string
	InfoMessage     string
	CrdURL          string
}

func (app *HelmApp) Install() error {
	// Download helm client
	// TODO: Move it to arkade init or at a more general place?
	userPath, err := config.InitUserDir()
	if err != nil {
		return err
	}

	clientArch, clientOS := env.GetClientArch()
	log.Printf("Client: %s, %s\n", clientArch, clientOS)
	log.Printf("User dir established as: %s\n", userPath)
	os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

	_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, true)
	if err != nil {
		return err
	}

	// Add the needed Repository
	// arkade-name --> so we handle it by our own and is has no user impact
	// Only support helm 3 here
	err = helm.AddHelmRepo(fmt.Sprintf("arkade-%s", app.Name), app.ChartRepository, true)

	if err != nil {
		return err
	}

	// Update the repo
	err = helm.UpdateHelmRepos(true)

	if err != nil {
		return err
	}

	// create namespace
	// TODO: create a funcion in pkg/kubernetes
	namespaceResult, err := kubernetes.KubectlTask("create", "namespace", app.Namespace)
	if err != nil {
		return err
	}

	if namespaceResult.ExitCode != 0 {
		log.Printf("[Warning] unable to create namespace %s, may already exist: %s", app.Namespace, namespaceResult.Stderr)
	}

	// apply crds if needed
	if app.CrdURL != "" {
		log.Println("Applying CRD")
		res, err := kubernetes.KubectlTask("apply", "--validate=false", "-f", app.CrdURL)
		if err != nil {
			return err
		}

		if res.ExitCode > 0 {
			return fmt.Errorf("error applying CRD from: %s, error: %s", app.CrdURL, res.Stderr)
		}
	}

	// install helm chart
	err = helm.HelmUpgrade(app.Name, fmt.Sprintf("arkade-%s/%s", app.Name, app.ChartName), app.Namespace, app.ChartVersion, app.Values)

	if err != nil {
		return err
	}

	return nil
}

func (app *HelmApp) GetInfoMessage() string {
	return app.InfoMessage
}

//TODO: should be implemented
func (app *HelmApp) Verify() bool {
	return true
}

func (app *HelmApp) SetName(name string) *HelmApp {
	app.Name = name
	return app
}

func (app *HelmApp) SetNamespace(name string) *HelmApp {
	app.Namespace = name
	return app
}

func (app *HelmApp) SetChartName(name string) *HelmApp {
	app.ChartName = name
	return app
}

func (app *HelmApp) SetChartRepository(url string) *HelmApp {
	app.ChartRepository = url
	return app
}

func (app *HelmApp) SetChartVersion(version string) *HelmApp {
	app.ChartVersion = version
	return app
}

func (app *HelmApp) SetCrdUrl(url string) *HelmApp {
	app.CrdURL = url
	return app
}

func (app *HelmApp) SetInfoMessage(msg string) *HelmApp {
	app.InfoMessage = msg
	return app
}

func (app *HelmApp) AddValue(key, value string) *HelmApp {
	if app.Values == nil {
		app.Values = map[string]string{}
	}
	app.Values[key] = value
	return app
}
