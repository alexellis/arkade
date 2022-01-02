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
    - fish
    - powershell
`

func MakeShellCompletion() *cobra.Command {

	completion := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Output shell completion for the given shell (bash, zsh, fish or powershell)",
		Long: `
Outputs shell completion for the given shell (bash, zsh, fish or powershell)
For bash, this depends on the bash-completion binary.  Example installation instructions:
macOS:
	$ brew install bash-completion                   # for bash users
	$ source $(brew --prefix)/etc/bash_completion    # for bash users
	$ arkade completion bash > ~/.arkade-completion  # for bash users
	$ arkade completion zsh > ~/.arkade-completion   # for zsh users
	$ source ~/.arkade-completion
Ubuntu:
	$ apt-get install bash-completion  # for bash users
	$ source /etc/bash-completion      # for bash users
	$ source <(arkade completion bash) # for bash users
	$ source <(arkade completion zsh)  # for zsh users
Additionally, you may want to output the completion to a file and source in your .bashrc / .zshrc / config.fish / profile.ps1
`,
		Example:      completionCmd,
		SilenceUsage: true,
		ValidArgs:    []string{"bash", "zsh", "fish", "powershell"},
	}

	completion.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Print(completionCmd)
			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf(completionCmd)
		}

		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			runCompletionZsh(cmd, os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		default:
			return fmt.Errorf("shell completion not supported for shell: %s", args[0])
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
