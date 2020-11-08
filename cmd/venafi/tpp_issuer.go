package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeTPPIssuer makes the app for the TPP issuer
func MakeTPPIssuer() *cobra.Command {

	command := &cobra.Command{
		Use:   "tpp-issuer",
		Short: "Install the cert-manager issuer for Venafi TPP",
		Long: `Install the cert-manager issuer for Venafi TPP to obtain 
TLS certificates from enterprise-grade CAs from self-hosted Venafi 
instances.`,
		Example: `  arkade venafi install cloud-issuer
  arkade venafi install tpp-issuer --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing the TPP issuer for you now.")
		return nil
	}

	return command
}
