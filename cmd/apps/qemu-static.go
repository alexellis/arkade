// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"context"
	"fmt"
	"os"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/env"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

func MakeInstallQemuStatic() *cobra.Command {
	var qemuStatic = &cobra.Command{
		Use:   "qemu-static",
		Short: "Install qemu-user-static",
		Long: `Runs the qemu-user-static container in Docker to enable 
support for multi-arch builds.

Learn more:

https://github.com/multiarch/qemu-user-static`,
		Example:      `  arkade install qemu-static`,
		Aliases:      []string{"qemu-user-static"},
		SilenceUsage: true,
	}

	qemuStatic.RunE = func(command *cobra.Command, args []string) error {

		arch, _ := env.GetClientArch()

		if arch != "x86_64" {
			return fmt.Errorf(`qemu-user-static is only supported on the AMD64 architecture, found: %s`, arch)
		}

		fmt.Printf("Running \"docker run --rm --privileged multiarch/qemu-user-static --reset -p yes\"\n\n")

		if err := runQemuStaticContainer(); err != nil {
			return err
		}

		fmt.Printf("\n\n%s\n\n", qemuStaticPostInstallMsg)

		return nil
	}

	return qemuStatic
}

const QemuStaticInfoMsg = `# Find out more at:
# https://github.com/multiarch/qemu-user-static`

const qemuStaticPostInstallMsg = `=======================================================================
= qemu-user-static has been installed.                                        =
=======================================================================` +
	"\n\n" + QemuStaticInfoMsg + "\n\n" + pkg.SupportMessageShort

func runQemuStaticContainer() error {
	task := execute.ExecTask{
		Command: "docker",
		Args: []string{"run", "--rm", "--privileged",
			"multiarch/qemu-user-static", "--reset", "-p", "yes"},
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute(context.Background())

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	return nil
}
