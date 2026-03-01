package chart

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/images"
	"github.com/spf13/cobra"
)

func MakeUpgrade() *cobra.Command {
	var command = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade all images in a values.yaml file to the latest version",
		Long: `Upgrade all images in a values.yaml file to the latest version.
Container images must be specified at the top level, or one level down in the 
"image: " or "component.image: " field in a values.yaml file.

Returns exit code zero if all images were found on the remote registry.

Otherwise, it returns a non-zero exit code and the updated values.yaml file.

An arkade.yaml file can be colocated with the values.yaml file or specified via
--ignore-file flag

The contents should be in YAML given as:

ignore:
- image
- component1.image
- component2.image
`,
		Example: `  # Upgrade and write the changes
  arkade chart upgrade -f ./chart/values.yaml --write

  # Dry-run mode, don't write the changes (default) 
  arkade chart upgrade --verbose -f ./chart/values.yaml

  # Use an arkade config file to load an ignore list
  arkade chart upgrade --ignore-file ./chart/arkade.yaml  -f ./chart/values.yaml`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", "", "Path to values.yaml file")
	command.Flags().StringP("ignore-file", "i", "", "Path to an arkade.yaml config file with a list of image paths to ignore defined")

	command.Flags().BoolP("verbose", "v", false, "Verbose output")
	command.Flags().BoolP("write", "w", false, "Write the updated values back to the file, or stdout when set to false")
	command.Flags().IntP("depth", "d", 3, "how many levels deep into the YAML structure to walk looking for image: tags")
	command.Flags().IntP("workers", "c", 4, "number of workers to use")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetInt("depth")
		if err != nil {
			return fmt.Errorf("error with --depth usage: %s", err)
		}
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("invalid value for flag --file")
		}

		ignoreFile, _ := cmd.Flags().GetString("ignore-file")

		writeFile, _ := cmd.Flags().GetBool("write")

		verbose, _ := cmd.Flags().GetBool("verbose")
		depth, _ := cmd.Flags().GetInt("depth")
		workers, _ := cmd.Flags().GetInt("workers")

		if len(file) == 0 {
			return fmt.Errorf("flag --file is required")
		}

		if ext := path.Ext(file); ext != ".yaml" && ext != ".yml" {
			return fmt.Errorf("--file must be a YAML file")
		}

		if len(ignoreFile) > 0 {
			if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
				return fmt.Errorf("ignore file %s does not exist", ignoreFile)
			}
		} else {
			basePath := path.Dir(file)
			defaultConfig := path.Join(basePath, "arkade.yaml")

			if _, err := os.Stat(defaultConfig); err == nil {
				ignoreFile = defaultConfig
			}
		}

		var arkadeCfg *config.ArkadeConfig

		if len(ignoreFile) > 0 {
			var err error
			arkadeCfg, err = config.Load(ignoreFile)
			if err != nil {
				return err
			}
		} else {
			arkadeCfg = &config.ArkadeConfig{}
		}

		if verbose {
			log.Printf("Verifying images in: %s\n", file)
		}

		values, err := helm.Load(file)
		if err != nil {
			return err
		}

		filtered := helm.FilterImagesUptoDepth(values, depth, "", arkadeCfg)
		if len(filtered) == 0 {
			return fmt.Errorf("no images found in %s", file)
		}

		if verbose {
			if len(filtered) > 0 {
				log.Printf("Found %d images\n", len(filtered))
			}
		}

		wg := sync.WaitGroup{}
		wg.Add(workers)

		workChan := make(chan string, len(filtered))
		errChan := make(chan error, len(filtered))
		updatedImages := make(map[string]string)

		for i := 0; i < workers; i++ {
			go func() {

				defer wg.Done()

				for image := range workChan {
					if len(image) > 0 {
						updated, imageNameAndTag, err := images.UpdateImage(image, verbose)
						if err != nil {
							errChan <- err
							continue
						}
						if updated {
							updatedImages[image] = imageNameAndTag
						}
					}
				}
			}()
		}

		for k := range filtered {
			workChan <- k
		}

		close(workChan)
		wg.Wait()
		close(errChan)

		var joinedErrors error
		for err := range errChan {
			if err != nil {
				joinedErrors = errors.Join(joinedErrors, err)
			}
		}
		if joinedErrors != nil {
			return joinedErrors
		}

		rawValues, err := helm.ReplaceValuesInHelmValuesFile(updatedImages, file)
		if err != nil {
			return err
		}

		if len(updatedImages) > 0 && writeFile {
			if err := os.WriteFile(file, []byte(rawValues), 0600); err != nil {
				return err
			}
			log.Printf("Wrote %d updates to: %s", len(updatedImages), file)
		}

		return nil
	}

	return command
}
