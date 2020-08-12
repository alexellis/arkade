package apps

import "github.com/spf13/cobra"

// App stores details about each app to be used by the CLI
type App struct {
	Name        string
	MakeInstall *cobra.Command
	Info        string
}

// Apps lists all available apps
var Apps = []*App{}

func registerApp(a *App) {
	Apps = append(Apps, a)
}
