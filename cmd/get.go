// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeGet() *cobra.Command {
	tools := get.MakeTools()

	var command = &cobra.Command{
		Use:   "get",
		Short: "Get a release of a tool or application and install it on your local computer.",
		Example: `  arkade get kubectl
  arkade get openfaas`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(arkadeGet)
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

		downloadURL, err := get.GetDownloadURL(tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
		if err != nil {
			return err
		}

		fmt.Println(downloadURL)

		res, err := http.DefaultClient.Get(downloadURL)
		if err != nil {
			return err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("incorrect status for downloading tool: %d", res.StatusCode)
		}

		_, fileName := path.Split(downloadURL)
		tmp := os.TempDir()

		outFilePath := path.Join(tmp, fileName)

		out, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err = io.Copy(out, res.Body); err != nil {
			return err
		}

		finalName := tool.Name
		if strings.Contains(strings.ToLower(operatingSystem), "mingw") {
			finalName = finalName + ".exe"
		}

		fmt.Printf(`Tool written to: %s

Run the following to copy to install the tool:

chmod +x %s
sudo install -m 755 %s /usr/local/bin/%s
`, outFilePath, outFilePath, outFilePath, finalName)

		return err
	}

	return command
}

const arkadeGet = `Use "arkade get TOOL" to download a tool or application:

  - kubectl
  - faas-cli`
