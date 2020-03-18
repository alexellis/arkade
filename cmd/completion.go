// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var shell string

func MakeCompletion() *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:   "completion SHELL",
		Short: "Generates shell completion",
		Long: `Generates shell auto completion for Bash or ZSH.

You can enable completion using the following:
source <(arkade completion)

To configure your bash/zsh shell to load completions for each session, add to your bashrc/zshrc

# ~/.bashrc or ~/.zshrc
source <(arkade completion)
`,
		Example: `  arkade completion --shell bash
  arkade completion --shell zsh`,
		RunE: runCompletion,
	}

	completionCmd.Flags().StringVar(&shell, "shell", "", "Outputs shell completion, must be bash or zsh")

	return completionCmd
}

func runCompletion(cmd *cobra.Command, args []string) (err error) {
	if shell == "" {
		re := regexp.MustCompile(`.*/`)
		shell = re.ReplaceAllString(os.Getenv("SHELL"), "")
	}

	switch shell {
	case "bash":
		err = bashCompletion(cmd)
		if err != nil {
			return err
		}
		return nil

	case "zsh":
		err = generateZshCompletion(cmd)
		if err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("%q shell not supported, must be bash or zsh", shell)
	}
}

func bashCompletion(cmd *cobra.Command) error {
	err := cmd.Parent().GenBashCompletion(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

func generateZshCompletion(cmd *cobra.Command) error {
	zshHead := "#compdef arkade\n"

	out := os.Stdout

	_, err := out.Write([]byte(zshHead))
	if err != nil {
		return err
	}

	zshInitialization := `
__arkade_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__arkade_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__arkade_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__arkade_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__arkade_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__arkade_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__arkade_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__arkade_filedir() {
	local RET OLD_IFS w qw
	__arkade_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __arkade_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__arkade_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__arkade_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
	printf %q "$1"
	fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__arkade_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__arkade_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__arkade_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__arkade_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__arkade_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__arkade_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__arkade_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	_, err = out.Write([]byte(zshInitialization))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	cmd.Parent().GenBashCompletion(buf)

	_, err = out.Write(buf.Bytes())
	if err != nil {
		return err
	}

	zshTail := `
BASH_COMPLETION_EOF
}
__arkade_bash_source <(__arkade_convert_bash_to_zsh)
_complete arkade 2>/dev/null
`
	_, err = out.Write([]byte(zshTail))
	if err != nil {
		return err
	}

	return nil
}
