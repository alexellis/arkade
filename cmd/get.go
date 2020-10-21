// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

// MakeGet creates the Get command to download software
func MakeGet() *cobra.Command {
	tools := get.MakeTools()

	var command = &cobra.Command{
		Use:   "get",
		Short: `The get command downloads a tool`,
		Long: `The get command downloads a CLI or application from the specific tool's 
releases or downloads page. The tool is usually downloaded in binary format 
and provides a fast and easy alternative to a package manager.`,
		Example: `  arkade get helm
  arkade get linkerd2 --stash=false
  arkade get terraform --version=0.12.0
  arkade get kubectl --progress=false

  # Get a complete list of CLIs to download:
  arkade get --help`,
		SilenceUsage: true,
	}

	command.Flags().Bool("progress", true, "Display a progress bar")
	command.Flags().Bool("stash", true, "When set to true, stash binary in HOME/.arkade/bin/, otherwise store in /tmp/")
	command.Flags().StringP("version", "v", "", "Download a specific version")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			const arkadeGet = `Use "arkade get TOOL" to download a tool or application:`

			buf := ""
			for _, t := range tools {
				buf = buf + t.Name + "\n"
			}
			fmt.Println(arkadeGet + "\n" + buf)
			return nil
		}

		var tool *get.Tool

		if len(args) == 1 {
			for _, t := range tools {
				if t.Name == args[0] {
					tool = &t
					break
				}
			}
		}
		if tool == nil {
			return fmt.Errorf("cannot get tool: %s", args[0])
		}

		fmt.Printf("Downloading %s\n", tool.Name)

		arch, operatingSystem := env.GetClientArch()
		version := ""

		if command.Flags().Changed("version") {
			version, _ = command.Flags().GetString("version")
		}

		stash, _ := command.Flags().GetBool("stash")
		progress, _ := command.Flags().GetBool("progress")
		if p, ok := os.LookupEnv("ARKADE_PROGRESS"); ok {
			b, err := strconv.ParseBool(p)
			if err != nil {
				return fmt.Errorf("ARKADE_PROGRESS is not a valid boolean")
			}

			progress = b
		}

		dlMode := get.DownloadTempDir
		if stash {
			dlMode = get.DownloadArkadeDir
		}

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			select {
			case <-signalChan:
				os.Exit(2)
			}
		}()

		outFilePath, finalName, err := get.Download(tool, arch, operatingSystem, version, dlMode, progress)

		if err != nil {
			return err
		}

		fmt.Printf("Tool written to: %s\n\n", outFilePath)

		if dlMode == get.DownloadTempDir {
			fmt.Printf(`Run the following to copy to install the tool:

chmod +x %s
sudo install -m 755 %s /usr/local/bin/%s
`, outFilePath, outFilePath, finalName)
		} else {
			fmt.Printf(`Run the following to add the (%s) binary to your PATH variable

export PATH=$PATH:$HOME/.arkade/bin/

%s

`, finalName, outFilePath)

		}
		return err
	}

	return command
}
