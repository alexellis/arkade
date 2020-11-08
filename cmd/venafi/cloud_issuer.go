package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeCloudIssuer makes an app for the Venafi Cloud issuer
func MakeCloudIssuer() *cobra.Command {

	command := &cobra.Command{
		Use:   "cloud-issuer",
		Short: "Install the cert-manager issuer for Venafi cloud",
		Long: `Install the cert-manager issuer for Venafi cloud to obtain 
TLS certificates from enterprise-grade CAs.`,
		Example: `  arkade venafi install cloud-issuer
  arkade venafi install cloud-issuer --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing the cloud-issuer for you now.")
		return nil
	}

	return command
}
