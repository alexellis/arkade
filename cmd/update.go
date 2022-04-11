// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

func MakeUpdate() *cobra.Command {
	var command = &cobra.Command{
		Use:          "update",
		Short:        "Print update instructions",
		Example:      `  arkade update`,
		Aliases:      []string{"u"},
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println(arkadeUpdate)

		fmt.Println("\n", aec.Bold.Apply(pkg.SupportMessageShort))
	}
	return command
}

const arkadeUpdate = `You can update arkade with the following:

# Remove cached versions of tools
rm -rf $HOME/.arkade

# For Linux/MacOS:
curl -SLfs https://get.arkade.dev | sudo sh

# For Windows (using Git Bash)
curl -SLfs https://get.arkade.dev | sh

# Or download from GitHub: https://github.com/alexellis/arkade/releases

Thanks for using arkade!`
