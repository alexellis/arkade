// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/go-execute/v2"
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
	command.RunE = func(cmd *cobra.Command, args []string) error {

		name := "arkade"
		toolList := get.MakeTools()
		var tool *get.Tool
		for _, t := range toolList {
			if t.Name == name {
				tool = &t
				break
			}
		}

		release, err := get.FindGitHubRelease("alexellis", name)
		if err != nil {
			return err
		}

		executable, err := os.Executable()
		if err != nil {
			return err
		}

		task := execute.ExecTask{
			Command: executable,
			Args:    []string{"version"},
		}

		res, err := task.Execute(context.TODO())
		if err != nil {
			return err
		}

		fmt.Printf("Latest release: %s\n", release)

		if strings.Contains(res.Stdout, release) {
			fmt.Println("You are already using the latest version of arkade.")

			fmt.Println("\n\n", aec.Bold.Apply(pkg.SupportMessageShort))

			return nil
		}

		arch, operatingSystem := env.GetClientArch()
		arch = strings.ToLower(arch)
		operatingSystem = strings.ToLower(operatingSystem)

		if arch == "x86_64" {
			arch = "amd64"
		}

		downloadUrl, err := get.GetDownloadURL(tool, operatingSystem, arch, release, false)
		if err != nil {
			return err
		}

		binary, err := get.DownloadFileP(downloadUrl, true)
		if err != nil {
			return err
		}

		if err := replaceExec(executable, binary); err != nil {
			return err
		}

		fmt.Printf("Replaced: %s.. OK.", executable)

		fmt.Println("\n\n", aec.Bold.Apply(pkg.SupportMessageShort))
		return nil
	}
	return command
}

// Copy the new binary to the same directory as the current binary before calling os.Rename to prevent an
// 'invalid cross-device link' error because the source and destination are not on the same file system.
func replaceExec(currentExec, newBinary string) error {
	targetDir := filepath.Dir(currentExec)
	filename := filepath.Base(currentExec)
	newExec := filepath.Join(targetDir, fmt.Sprintf(".%s.new", filename))

	// Copy the contents of newbinary to a new executable file
	sf, err := os.Open(newBinary)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(newExec, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}

	// Replace the current executable file with the new executable file
	if err := os.Rename(newExec, currentExec); err != nil {
		return err
	}

	return nil
}
