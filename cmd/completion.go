// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const completionCmd = `Use "arkade completion SHELL" to generate SHELL completion for:
    - bash
    - zsh
`

func MakeShellCompletion() *cobra.Command {

	completion := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Output shell completion for the given shell (bash or zsh)",
		Long: `
Outputs shell completion for the given shell (bash or zsh)
This depends on the bash-completion binary.  Example installation instructions:
OS X:
	$ brew install bash-completion
	$ source $(brew --prefix)/etc/bash_completion
	$ arkade completion bash > ~/.arkade-completion  # for bash users
	$ arkade completion zsh > ~/.arkade-completion   # for zsh users
	$ source ~/.arkade-completion
Ubuntu:
	$ apt-get install bash-completion
	$ source /etc/bash-completion
	$ source <(arkade completion bash) # for bash users
	$ source <(arkade completion zsh)  # for zsh users
Additionally, you may want to output the completion to a file and source in your .bashrc
`,
		Example:      completionCmd,
		SilenceUsage: true,
		ValidArgs:    []string{"bash", "zsh"},
	}

	completion.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(completionCmd)
			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf(completionCmd)
		}

		shellName := args[0]

		switch shellName {
		case "bash":
			rootCmd(cmd).GenBashCompletion(os.Stdout)
		case "zsh":
			runCompletionZsh(cmd, os.Stdout)
		default:
			return fmt.Errorf("shell completion not supported for shell: %s", shellName)
		}

		return nil
	}

	return completion
}

func runCompletionZsh(cmd *cobra.Command, out io.Writer) {
	var zshCompdef = "\ncompdef _arkade arkade\n"

	rootCmd(cmd).GenZshCompletion(out)
	io.WriteString(out, zshCompdef)
}

func rootCmd(cmd *cobra.Command) *cobra.Command {
	parent := cmd
	for parent.HasParent() {
		parent = parent.Parent()
	}
	return parent
}
