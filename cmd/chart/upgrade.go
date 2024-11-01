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

type tagAttributes struct {
	hasSuffix bool
	hasMajor  bool
	hasMinor  bool
	hasPatch  bool
	original  string
}

func (c *tagAttributes) attributesMatch(n tagAttributes) bool {
	return c.hasMajor == n.hasMajor &&
		c.hasMinor == n.hasMinor &&
		c.hasPatch == n.hasPatch &&
		c.hasSuffix == n.hasSuffix
}

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

	candidateTag, hasSemVerTag := getCandidateTag(ref, tag)

	if !hasSemVerTag {
		return false, iName, fmt.Errorf("no valid semver tags of current format found for %s", imageName)
	}

	laterVersionB := false

	// AE: Don't upgrade to an RC tag, even if it's newer.
	if tagIsUpgradeable(tag, candidateTag) {

		laterVersionB = true

		iName = fmt.Sprintf("%s:%s", imageName, candidateTag)
		if v {
			log.Printf("[%s] %s => %s", imageName, tag, candidateTag)
		}
	}

	return laterVersionB, iName, nil
}

func tagIsUpgradeable(current, candidate string) bool {

	if strings.EqualFold(current, "latest") {
		return false
	}

	currentSemVer, _ := semver.NewVersion(current)
	candidateSemVer, _ := semver.NewVersion(candidate)

	return candidateSemVer.Compare(currentSemVer) == 1 && candidateSemVer.Prerelease() == currentSemVer.Prerelease()

}

func getCandidateTag(discoveredTags []string, currentTag string) (string, bool) {

	var candidateTags []*semver.Version
	for _, tag := range discoveredTags {
		v, err := semver.NewVersion(tag)
		if err == nil {
			candidateTags = append(candidateTags, v)
		}
	}

	if len(candidateTags) > 0 {

		currentTagAttr := getTagAttributes(currentTag)
		sort.Sort(sort.Reverse(semver.Collection(candidateTags)))

		for _, candidate := range candidateTags {
			candidateTagAttr := getTagAttributes(candidate.Original())
			if currentTagAttr.attributesMatch(candidateTagAttr) {
				return candidate.Original(), true
			}
		}
	}

	return "", false

}

func getTagAttributes(t string) tagAttributes {

	tagParts := strings.Split(t, "-")
	tagLevels := strings.Split(tagParts[0], ".")

	return tagAttributes{
		hasSuffix: len(tagParts) > 1,
		hasMajor:  len(tagLevels) >= 1 && tagLevels[0] != "",
		hasMinor:  len(tagLevels) >= 2,
		hasPatch:  len(tagLevels) == 3,
		original:  t,
	}
}
