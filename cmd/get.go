// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

// MakeGet creates the Get command to download software
func MakeGet() *cobra.Command {
	tools := get.MakeTools()
	sort.Sort(tools)
	var validToolOptions []string = make([]string, len(tools))
	for _, t := range tools {
		validToolOptions = append(validToolOptions, t.Name)
	}

	var command = &cobra.Command{
		Use:   "get",
		Short: `The get command downloads a tool`,
		Long: `The get command downloads a CLI or application from the specific tool's
releases or downloads page. The tool is usually downloaded in binary format
and provides a fast and easy alternative to a package manager.`,
		Example: `  arkade get helm

  # Options for the download
  arkade get linkerd2 --stash=false
  arkade get kubectl --progress=false

  # Override the version
  arkade get terraform --version=0.12.0
  arkade get kubectl@v1.19.3

  # Override the OS
  arkade get helm --os darwin --arch aarch64 
  arkade get helm --os linux --arch armv7l 

  # Get a complete list of CLIs to download:
  arkade get`,
		SilenceUsage: true,
		Aliases:      []string{"g", "d", "download"},
		ValidArgs:    validToolOptions,
	}

	clientArch, clientOS := env.GetClientArch()
	command.Flags().Bool("progress", true, "Display a progress bar")
	command.Flags().StringP("output", "o", "", "Output format of the list of tools (table/markdown/list)")
	command.Flags().Bool("stash", true, "When set to true, stash binary in HOME/.arkade/bin/, otherwise store in /tmp/")
	command.Flags().StringP("version", "v", "", "Download a specific version")
	command.Flags().String("arch", clientArch, "CPU architecture for the tool")
	command.Flags().String("os", clientOS, "Operating system for the tool")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			output, _ := command.Flags().GetString("output")

			if len(output) > 0 {
				if get.TableFormat(output) == get.MarkdownStyle {
					get.CreateToolsTable(tools, get.MarkdownStyle)
				} else if get.TableFormat(output) == get.ListStyle {
					for _, r := range tools {
						fmt.Printf("%s\n", r.Name)
					}

				} else {
					get.CreateToolsTable(tools, get.TableStyle)
				}
			} else {
				get.CreateToolsTable(tools, get.TableStyle)
			}
			return nil
		}

		version := ""
		if command.Flags().Changed("version") {
			version, _ = command.Flags().GetString("version")
		}

		downloadURLs, err := get.GetDownloadURLs(tools, args, version)
		if err != nil {
			return err
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
			<-signalChan
			os.Exit(2)
		}()

		var outFilePath string
		var localToolsStore []get.ToolLocal

		arch, _ := command.Flags().GetString("arch")
		if err := get.ValidateArch(arch); err != nil {
			return err
		}
		operatingSystem, _ := command.Flags().GetString("os")
		if err := get.ValidateOS(operatingSystem); err != nil {
			return err
		}

		for _, tool := range downloadURLs {
			fmt.Printf("Downloading: %s\n", tool.Name)
			outFilePath, _, err = get.Download(&tool,
				arch,
				operatingSystem,
				version,
				dlMode,
				progress)
			if err != nil {
				return err
			}

			localToolsStore = append(localToolsStore, get.ToolLocal{Name: tool.Name, Path: outFilePath})
			fmt.Printf("\nTool written to: %s\n\n", outFilePath)
		}

		msg, err := get.PostInstallationMsg(dlMode, localToolsStore)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", msg)

		return err
	}
	return command
}
