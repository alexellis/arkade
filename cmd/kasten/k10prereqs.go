// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// kasten contains a suite of Sponsored Apps for arkade
package kasten

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alexellis/arkade/pkg/config"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
)

const k10preflightInfoMsg = `Run pre-flight checks for Kasten's Backup and Migration solution
Read more at:

https://docs.kasten.io/latest/install/requirements.html#pre-flight-checks`

func MakeInstallK10Preflight() *cobra.Command {
	var k10cmd = &cobra.Command{
		Use:          "preflight",
		Short:        "Install preflight",
		Long:         k10preflightInfoMsg,
		SilenceUsage: true,
	}

	k10cmd.Flags().Bool("dry-run", false, "Print the commands that would be run by the preflight script.")

	k10cmd.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}
		dryRun, err := command.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}

		primerURL := "https://docs.kasten.io/tools/k10_primer.sh"
		req, err := http.NewRequest(http.MethodGet, primerURL, nil)
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error downloading k10 preflight tool: %w", err)
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected server status %d downloading k10 preflight tool", res.StatusCode)
		}

		body, _ := ioutil.ReadAll(res.Body)

		outFile := os.TempDir() + "/k10_primer.sh"

		if err := ioutil.WriteFile(outFile, body, 0755); err != nil {
			return fmt.Errorf("error writing k10 preflight tool to: %s %w", outFile, err)
		}

		fmt.Printf("Downloaded %s to %s\n", primerURL, outFile)

		if dryRun {
			fmt.Printf("Preflight script contents:\n%s\n", string(body))
			fmt.Printf("Run this command with --dry-run=false to execute the script.\n")
			return nil
		}

		scriptRes, err := runScript(outFile)
		if err != nil {
			return err
		}
		if scriptRes.ExitCode != 0 {
			return fmt.Errorf("exit code %d, stderr: %s", scriptRes.ExitCode, scriptRes.Stderr)
		}

		return nil
	}

	return k10cmd
}

func runScript(file string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s", "/bin/bash"),
		Args:        []string{file},
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute()
	if err != nil {
		return res, err
	}

	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	return res, nil
}
