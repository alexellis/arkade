package cmd

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

func MakeNewApp() *cobra.Command {

	newApp := &cobra.Command{
		Use:          "new-app",
		Short:        "Command to initialise a new app",
		Long:         "This command can be used to generate a new app, it wont be registered to arkade but the basic template is created",
		Example:      `  arkade new-app --name MyNewApp --command newApp`,
		SilenceUsage: true,
		Hidden:       true,
	}

	newApp.Flags().String("name", "", "Set --name to the name for your app, users will see this")
	newApp.Flags().String("command", "", "Set --command to the command value for your app, used for the cli command")
	newApp.Flags().String("chartURL", "", "optionally set the chart URL in the command, otherwise it will be up to you to set")
	newApp.Flags().String("helmRepoName", "", "optionally set the helm repo name, if ommited it will need to be filled later")
	newApp.RunE = func(cmd *cobra.Command, args []string) error {

		name, _ := newApp.Flags().GetString("name")
		command, _ := newApp.Flags().GetString("command")
		chartURL, _ := newApp.Flags().GetString("chartURL")
		helmRepoName, _ := newApp.Flags().GetString("helmRepoName")

		if len(name) == 0 || len(command) == 0 {
			return fmt.Errorf("name and command must be present, got name: [%s] and command: [%s]", name, command)
		}

		app := App{
			Name:         name,
			Command:      command,
			ChartURL:     chartURL,
			HelmRepoName: helmRepoName,
			Year:         time.Now().Year(),
		}
		tmpl, err := template.New("test").Parse(NewAppTemplate)
		if err != nil {
			return err
		}

		f, err := os.Create(fmt.Sprintf("./cmd/apps/%s_app.go", command))
		if err != nil {
			return err
		}

		err = tmpl.Execute(f, app)

		if err != nil {
			return err
		}
		return nil
	}

	return newApp

}

type App struct {
	Name         string
	Command      string
	ChartURL     string
	HelmRepoName string
	Year         int
}

var NewAppTemplate = `// Copyright (c) arkade author(s) {{.Year}}. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstall{{.Name}}() *cobra.Command {
	var {{.Command}}App = &cobra.Command{
		Use:          "{{.Command}}",
		Short:        "",
		Long:         "",
		Example:      "",
		SilenceUsage: true,
	}

	{{.Command}}App.Flags().StringP("namespace", "n", "default", "The namespace to install chartmuseum (default: default")
	{{.Command}}App.Flags().Bool("update-repo", true, "Update the helm repo")

	{{.Command}}App.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := {{.Command}}App.Flags().GetString("namespace")
		updateRepo, _ := {{.Command}}App.Flags().GetBool("update-repo")

		overrides := map[string]string{}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		{{.Command}}Options := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("{{.HelmRepoName}}").
			WithHelmURL("{{.ChartURL}}").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath).
			WithHelmUpdateRepo(updateRepo)

		_, err := apps.MakeInstallChart({{.Command}}Options)
		if err != nil {
			return err
		}

		println({{.Name}}InfoMsg)
		return nil
	}

	return {{.Command}}App
}


const {{.Name}}InfoMsg = ` + "`# Get started with {{.Name}} here:\n`" + `

const {{.Name}}InstallMsg = ` + "`" + `=======================================================================
= {{.Name}} has been installed.                                   =
=======================================================================` + "`" + `+
	"\n\n" + {{.Name}}InfoMsg + "\n\n" + pkg.ThanksForUsing`
