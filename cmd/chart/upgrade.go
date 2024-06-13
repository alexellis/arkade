package chart

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/google/go-containerregistry/pkg/crane"
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

Otherwise, it returns a non-zero exit code and the updated values.yaml file.`,
		Example: `arkade upgrade -f ./chart/values.yaml
  arkade upgrade --verbose -f ./chart/values.yaml`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", "", "Path to values.yaml file")
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

		if verbose {
			log.Printf("Verifying images in: %s\n", file)
		}

		values, err := helm.Load(file)
		if err != nil {
			return err
		}

		filtered := helm.FilterImagesUptoDepth(values, depth)
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
						updated, imageNameAndTag, err := updateImages(image, verbose)
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

func splitImageName(reposName string) (string, string) {
	nameParts := strings.SplitN(reposName, ":", 2)
	return nameParts[0], nameParts[1]
}

func updateImages(iName string, v bool) (bool, string, error) {

	imageName, tag := splitImageName(iName)
	ref, err := crane.ListTags(imageName)
	if err != nil {
		return false, iName, errors.New("unable to list tags for " + imageName)
	}

	var vs []*semver.Version
	for _, r := range ref {
		v, err := semver.NewVersion(r)
		if err == nil {
			vs = append(vs, v)
		}
	}

	if len(vs) == 0 {
		return false, iName, fmt.Errorf("no valid semver tags found for %s", imageName)
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	latestTag := vs[0].String()
	// Semver is "eating" the "v" prefix, so we need to add it back, if it was there in first place
	if strings.HasPrefix(tag, "v") {
		latestTag = "v" + latestTag
	}

	laterVersionB := false

	// AE: Don't upgrade to an RC tag, even if it's newer.
	if tagIsUpgradeable(tag, latestTag) {

		laterVersionB = true

		iName = fmt.Sprintf("%s:%s", imageName, latestTag)
		if v {
			log.Printf("[%s] %s => %s", imageName, tag, latestTag)
		}
	}

	return laterVersionB, iName, nil
}

func tagIsUpgradeable(currentTag, latestTag string) bool {

	return latestTag != currentTag && !strings.Contains(strings.ToLower(latestTag), "-rc") && !strings.EqualFold(currentTag, "latest")

}
