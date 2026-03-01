package gha

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexellis/arkade/pkg/gha"
	"github.com/spf13/cobra"
)

func MakeUpgrade() *cobra.Command {
	var command = &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade actions in GitHub Actions workflow files to the latest major version",
		Aliases: []string{"u"},
		Long: `Upgrade actions in GitHub Actions workflow files to the latest major version.

Processes all workflow YAML files in .github/workflows/ or a single file.
Only bumps major versions (e.g. actions/checkout@v3 to actions/checkout@v4).
`,
		Example: `  # Upgrade all workflows in the current directory
  arkade gha upgrade

  # Upgrade a single workflow file
  arkade gha upgrade -f .github/workflows/build.yaml

  # Dry-run mode, don't write changes
  arkade gha upgrade --write=false`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", ".", "Path to workflow file or directory")
	command.Flags().BoolP("verbose", "v", true, "Verbose output")
	command.Flags().BoolP("write", "w", true, "Write the updated values back to the file")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("file")
		verbose, _ := cmd.Flags().GetBool("verbose")
		writeFile, _ := cmd.Flags().GetBool("write")

		files, err := gha.FindWorkflows(target)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return fmt.Errorf("no workflow files found in %s", target)
		}

		if verbose {
			fmt.Printf("Found %d workflow file(s)\n\n", len(files))
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		totalUpdates := 0

		for _, file := range files {
			if verbose {
				fmt.Printf("Processing: %s\n", file)
			}

			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			replacements, err := gha.ProcessWorkflow(data, client, verbose)
			if err != nil {
				return err
			}

			if verbose && len(replacements) > 0 {
				fmt.Println("Detected following replacements:")
				for old, newVer := range replacements {
					fmt.Printf("  %s -> %s\n", old, newVer)
				}
			}

			if len(replacements) > 0 {
				updated := gha.ApplyReplacements(data, replacements)
				totalUpdates += len(replacements)

				if writeFile {
					if err := os.WriteFile(file, []byte(updated), 0644); err != nil {
						return err
					}
				} else {
					fmt.Print(updated)
				}
			}

			fmt.Println()
		}

		if totalUpdates > 0 && writeFile {
			log.Printf("Wrote %d update(s) across %d file(s)", totalUpdates, len(files))
		}

		return nil
	}

	return command
}
