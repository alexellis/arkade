package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeInfo() *cobra.Command {

	command := &cobra.Command{
		Use:          "info",
		Short:        "Info for an app",
		Long:         `Info for an app`,
		Example:      `  arkade venafi info [APP]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("give an app as an argument")
		}
		fmt.Printf("Info for your app: %s\n", args[0])

		return nil
	}

	return command
}
