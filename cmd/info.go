// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeInfo() *cobra.Command {

	info := &cobra.Command{
		Use:   "info",
		Short: "Find info about a Kubernetes app",
		Long:  "Find info about how to use the installed Kubernetes app",
		Example: `  arkade info [APP]
arkade info openfaas
arkade info inlets-operator
arkade info mongodb
arkade info
arkade info --help`,
		SilenceUsage: true,
	}

	info.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("Run arkade info APP_NAME for more")
			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf("you can only get info about exactly one installed app")
		}

		appList := GetApps()
		appName := args[0]
		if _, ok := appList[appName]; !ok {
			return fmt.Errorf("no info available for app: %s", appName)
		}
		fmt.Printf("Info for app: %s\n", appName)
		fmt.Println(appList[appName].InfoMessage)
		return nil

	}

	return info
}
