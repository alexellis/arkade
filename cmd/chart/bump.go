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
		Short: "Bump the version of the Helm chart.",
		Long: `Bump the version present in the Chart.yaml of a Helm chart.
To bump the version only if the adjacent values file has changes then specify
the --check-for-value-updates flag. If the values file has no changes the command
returns early with an exit code zero.
`,
		Example: `arkade chart bump --dir ./chart
  arkade chart bump --dir ./charts --check-for-value-updates values.yaml`,
		SilenceUsage: true,
	}

	command.Flags().StringP("dir", "d", "", "Path to the Helm chart directory or a directory containing Helm charts")
	command.Flags().BoolP("verbose", "v", false, "Verbose output")
	command.Flags().BoolP("write", "w", false, "Write the updated values back to the file, or stdout when set to false")
	command.Flags().String("check-for-value-updates", "", "Name of the values file to check if the chart's values have been modified before bumping version")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		chartDir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return fmt.Errorf("invalid value for --dir")
		}
		if chartDir == "" {
			return fmt.Errorf("flag --dir is required")
		}
		verbose, _ := cmd.Flags().GetBool("verbose")
		write, err := cmd.Flags().GetBool("write")
		if err != nil {
			return fmt.Errorf("invalid value for --write")
		}
		valuesFile, err := cmd.Flags().GetString("check-for-value-updates")
		if err != nil {
			return fmt.Errorf("invalid value for --check-for-value-updates")
		}

		// Map with key as the path to Chart.yaml and the value as the parsed contents of Chart.yaml
		chartYamlPath := filepath.Join(chartDir, ChartYamlFileName)
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
			if valuesFile != "" {
				absPath, err := filepath.Abs(chartDir)
				if err != nil {
					return err
				}
				absValuesFile := filepath.Join(absPath, valuesFile)
				_, err = os.Stat(absValuesFile)
				if err != nil {
					return fmt.Errorf("unable to find values file: %s", absValuesFile)
				}

				// Run `git diff --exit-code <file>` to check if the values file has any changes.
				// An exit code of 0 indicates that there are no changes, thus we skip bumping the
				// version of the chart.
				cmd := execute.ExecTask{
					Command: "git",
					Args:    []string{"diff", "--exit-code", valuesFile},
					Cwd:     absPath,
				}
				res, err := cmd.Execute(context.Background())
				if err != nil {
					return fmt.Errorf("could not check updates to chart values: %s", err)
				}

				if res.ExitCode == 0 {
					fmt.Printf("no changes detected in %s; skipping version bump\n", filepath.Join(chartDir, valuesFile))
					os.Exit(0)
				}
			}

			ver, err := semver.NewVersion(version)
			if err != nil {
				return fmt.Errorf("%s", err)
			}
			newVer := ver.IncMinor()
			fmt.Printf("%s %s => %s\n", chartYamlPath, ver.String(), newVer.String())
			if write {
				if verbose {
					log.Printf("Bumping version")
				}
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
			}
		}
		return nil
	}
	return command
}
