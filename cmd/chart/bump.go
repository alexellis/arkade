package chart

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/alexellis/arkade/pkg/helm"
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
		Short: "Bump the version of the Helm chart(s)",
		Long: `Bump the version present in the Chart.yaml of a Helm chart.
If the provided directory contains multiple charts, then the --recursive flag
can be used to bump the version in all charts.`,
		Example: `arkade bump --dir ./chart
  arkade --dir ./charts --recursive`,
		SilenceUsage: true,
	}

	command.Flags().StringP("dir", "d", "", "Path to the Helm chart directory or a directory containing Helm charts")
	command.Flags().BoolP("recursive", "r", false, "Recursively iterate through directory while bumping chart versions")
	command.Flags().BoolP("verbose", "v", false, "Verbose output")
	command.Flags().BoolP("write", "w", false, "Write the updated values back to the file, or stdout when set to false")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		chartDir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return fmt.Errorf("invalid value for --dir")
		}
		if chartDir == "" {
			return fmt.Errorf("flag --dir is required")
		}
		verbose, _ := cmd.Flags().GetBool("verbose")
		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return fmt.Errorf("invalid value for --recursive")
		}
		write, err := cmd.Flags().GetBool("write")
		if err != nil {
			return fmt.Errorf("invalid value for --write")
		}

		// Map with key as the path to Chart.yaml and the value as the parsed contents of Chart.yaml
		chartYamls := make(map[string]helm.ValuesMap, 0)
		if !recursive {
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
			chartYamls[chartYamlPath] = values
		} else {
			filepath.WalkDir(chartDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.Name() == ChartYamlFileName || d.Name() == ChartYmlFileName {
					values, err := helm.Load(path)
					if err != nil {
						return err
					}
					chartYamls[path] = values
				}
				return nil
			})
			if len(chartYamls) > 0 {
				fmt.Printf("Found %d chart(s)\n", len(chartYamls))
			}
		}

		for file, contents := range chartYamls {
			// If the yaml does not contain a `version` key then skip it.
			if val, ok := contents[versionKey]; !ok {
				continue
			} else {
				version, ok := val.(string)
				if !ok {
					log.Printf("unable to find a valid version in %s", file)
					continue
				}
				ver, err := semver.NewVersion(version)
				if err != nil {
					continue
				}
				newVer := ver.IncMinor()
				fmt.Printf("%s %s => %s\n", file, ver.String(), newVer.String())
				if write {
					if verbose {
						log.Printf("Bumping version")
					}
					update := map[string]string{
						fmt.Sprintf("%s: %s", versionKey, ver.String()): fmt.Sprintf("%s: %s", versionKey, newVer.String()),
					}
					rawChartYaml, err := helm.ReplaceValuesInHelmValuesFile(update, file)
					if err != nil {
						return fmt.Errorf("unable to bump chart version in %s", file)
					}
					if err = os.WriteFile(file, []byte(rawChartYaml), 0600); err != nil {
						return fmt.Errorf("unable to write updated yaml to %s", file)
					}
				}
			}
		}
		return nil
	}
	return command
}
