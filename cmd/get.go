// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	units "github.com/docker/go-units"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
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
  arkade get kubectl@v1.19.3
  arkade get terraform --version=1.7.4

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
	command.Flags().StringP("format", "o", "", "Format format of the list of tools (table/markdown/list)")
	command.Flags().String("path", "", "Leave empty to store in HOME/.arkade/bin/, otherwise give a path for the resulting binaries")
	command.Flags().StringP("version", "v", "", "Download a specific version")
	command.Flags().String("arch", clientArch, "CPU architecture for the tool")
	command.Flags().String("os", clientOS, "Operating system for the tool")
	command.Flags().Bool("quiet", false, "Suppress most additional format")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			format, _ := command.Flags().GetString("format")

			if len(format) > 0 {
				if get.TableFormat(format) == get.MarkdownStyle {
					get.CreateToolsTable(tools, get.MarkdownStyle)
				} else if get.TableFormat(format) == get.ListStyle {
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

		movePath, _ := command.Flags().GetString("path")
		progress, _ := command.Flags().GetBool("progress")
		quiet, _ := command.Flags().GetBool("quiet")

		if quiet && !command.Flags().Changed("progress") {
			progress = false
		}

		if p, ok := os.LookupEnv("ARKADE_PROGRESS"); ok {
			b, err := strconv.ParseBool(p)
			if err != nil {
				return fmt.Errorf("ARKADE_PROGRESS is not a valid boolean")
			}

			progress = b
		}

		movePath = os.ExpandEnv(movePath)

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
			if !quiet {
				fmt.Printf("Downloading: %s\n", tool.Name)
			}
			outFilePath, _, err = get.Download(&tool,
				arch,
				operatingSystem,
				version,
				movePath,
				progress,
				quiet)

			// handle 404 error gracefully
			if errors.Is(err, &get.ErrNotFound{}) {

				extra := ""
				// 1. The tool either has an explicit GitHub URL
				// 2. or has no URL in the URLTemplate meaning it's on GitHub
				// 3. or there is no URLTemplate because a BinaryTemplate was used instead, meaning the tool is on GitHub
				if strings.Contains(tool.URLTemplate, "https://github.com/") ||
					!strings.Contains(tool.URLTemplate, "https://") ||
					len(tool.URLTemplate) == 0 {
					extra = fmt.Sprintf(`
* View the %s releases page: %s`, tool.Name, fmt.Sprintf("https://github.com/%s/%s/releases", tool.Owner, tool.Repo))
				}

				fmt.Fprintf(os.Stderr, `
The requested version of %s is not available or configured in arkade for %s/%s

* Check if a binary is available from the project for your Operating System%s
* Feel free to raise an issue at https://github.com/alexellis/arkade/issues for help

`, tool.Name, operatingSystem, arch, extra)

				return err
			}
			if err != nil {
				return err
			}

			localToolsStore = append(localToolsStore, get.ToolLocal{Name: tool.Name, Path: outFilePath})
			if !quiet {
				size := ""
				stat, err := os.Stat(outFilePath)
				if err == nil {
					size = "(" + units.HumanSize(float64(stat.Size())) + ")"
				}

				fmt.Printf("\nWrote: %s %s\n\n", outFilePath, size)
			}
		}

		nl := ""
		if !quiet {
			nl = "\n"
			msg, err := get.PostInstallationMsg(movePath, localToolsStore)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", msg)
		}

		if !quiet {
			fmt.Printf("%s%s\n", nl, aec.Bold.Apply(pkg.SupportMessageShort))
		}

		return err
	}
	return command
}
