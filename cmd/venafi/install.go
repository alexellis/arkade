package venafi

import "github.com/spf13/cobra"

func MakeInstall() *cobra.Command {

	command := &cobra.Command{
		Use:     "install",
		Short:   "Install Sponsored Apps for Venafi",
		Long:    `Install Sponsored Apps for Venafi`,
		Aliases: []string{"i"},
		Example: `  arkade venafi install [APP]
  arkade venafi install --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}
	return command
}
