package chart

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

const (
	versionKey        = "version"
	ChartYamlFileName = "Chart.yaml"
	ChartYmlFileName  = "Chart.yml"
)

func MakeBump() *cobra.Command {
	var command = &cobra.Command{
		Use:   "bump",
		Short: "Bump the patch version of the Helm chart.",
		Long: `Bump the version present in the Chart.yaml of a Helm chart
if there are changes staged in git.

Add --force to bump the version even if there are not changes staged in git.
`,
		Example: `  # Bump if there are changes
  arkade chart bump -f ./chart/values.yaml

  # Force a bump, even if there are no changes 
  arkade chart bump -f ./charts/values.yaml --check-for-updates`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", "", "Path to values.yaml file")
	command.Flags().BoolP("verbose", "v", false, "Verbose output")
	command.Flags().BoolP("write", "w", false, "Write the updated values back to the file, or stdout when set to false")
	command.Flags().Bool("force", false, "Update the version even if there are no changes staged in the adjacent folder to Chart.yaml")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		valuesFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("invalid value for --file")
		}
		if valuesFile == "" {
			return fmt.Errorf("flag --file is required")
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}

		write, err := cmd.Flags().GetBool("write")
		if err != nil {
			return err
		}
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		chartDir := filepath.Dir(valuesFile)
		chartYamlPath := filepath.Join(chartDir, ChartYamlFileName)

		// Map with key as the path to Chart.yaml and the value as the parsed contents of Chart.yaml
		var values helm.ValuesMap
		// Try to read a Chart.yaml, but if thats unsuccessful then fall back to Chart.yml
		if values, err = helm.Load(chartYamlPath); err != nil {
			if verbose {
				log.Printf("unable to read %s, falling back to Chart.yml\n", chartYamlPath)
			}
			chartYamlPath = filepath.Join(chartDir, ChartYmlFileName)
			if values, err = helm.Load(chartYamlPath); err != nil {
				return fmt.Errorf("unable to read Chart.yaml or Chart.yml in directory %s", chartDir)
			}
		}

		// If the yaml does not contain a `version` key then error out.
		if val, ok := values[versionKey]; !ok {
			return fmt.Errorf("unable to find a version in %s", chartYamlPath)
		} else {
			version, ok := val.(string)
			if !ok {
				log.Printf("unable to find a valid version in %s", chartYamlPath)
			}
			if !force {
				absPath, err := filepath.Abs(chartDir)
				if err != nil {
					return err
				}

				// Run `git diff --exit-code <file>` to check if any files in the chart dir changed.
				// An exit code of 0 indicates that there are no changes, thus we skip bumping the
				// version of the chart.
				cmd := execute.ExecTask{
					Command: "git",
					Args:    []string{"diff", "--exit-code", "."},
					Cwd:     absPath,
				}
				res, err := cmd.Execute(context.Background())
				if err != nil {
					return fmt.Errorf("could not check updates to chart values: %s", err)
				}

				if res.ExitCode == 0 {
					if verbose {
						fmt.Printf("No changes in: %s, ignoring.\n", chartDir)
					}
					os.Exit(0)
				}
			}

			ver, err := semver.NewVersion(version)
			if err != nil {
				return fmt.Errorf("%s", err)
			}

			newVer := ver.IncPatch()

			fmt.Printf("%s %s => %s\n", chartYamlPath, ver.String(), newVer.String())

			if write {

				update := map[string]string{
					fmt.Sprintf("%s: %s", versionKey, ver.String()): fmt.Sprintf("%s: %s", versionKey, newVer.String()),
				}
				rawChartYaml, err := helm.ReplaceValuesInHelmValuesFile(update, chartYamlPath)
				if err != nil {
					return fmt.Errorf("unable to bump chart version in %s", chartYamlPath)
				}
				if err = os.WriteFile(chartYamlPath, []byte(rawChartYaml), 0600); err != nil {
					return fmt.Errorf("unable to write updated yaml to %s", chartYamlPath)
				}
				fmt.Printf("Wrote to: %s. OK.\n", chartYamlPath)
			} else {
				fmt.Printf("Not updating: %s (--write=false)\n", chartYamlPath)
			}

		}
		return nil
	}
	return command
}
