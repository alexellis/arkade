// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package fstail

import (
	"os"

	"github.com/alexellis/fstail/pkg/fstail"
	"github.com/spf13/cobra"
)

func MakeFstail() *cobra.Command {
	var prefixStyle string
	var disablePrefix bool

	command := &cobra.Command{
		Use:   "fstail [FLAGS] [PATH] [MATCH]",
		Short: "Tail files with optional string matching",
		Long:  `Tail files in a directory with optional string matching and prefix styles.`,
		Example: `  arkade fstail
  arkade fstail /var/log/containers
  arkade fstail /var/log/containers "server error"
  arkade fstail --prefix k8s /var/log/containers
  arkade fstail --prefix none /var/log/containers`,
		Aliases: []string{"ft"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts fstail.RunOptions

			if len(args) == 1 {
				opts.WorkDir = args[0]
			} else if len(args) == 2 {
				opts.WorkDir = args[0]
				opts.Match = args[1]
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}
				opts.WorkDir = cwd
			}

			opts.DisableLogPrefix = disablePrefix
			if disablePrefix {
				opts.PrefixStyle = fstail.PrefixStyleNone
			} else if prefixStyle == "k8s" {
				opts.PrefixStyle = fstail.PrefixStyleK8s
			} else {
				opts.PrefixStyle = fstail.PrefixStyleFilename
			}

			return fstail.Run(opts)
		},
	}

	command.Flags().StringVar(&prefixStyle, "prefix", "filename", "Prefix style: filename, k8s, or none")
	command.Flags().BoolVar(&disablePrefix, "no-prefix", false, "Disable prefix printing (equivalent to --prefix none)")

	return command
}
