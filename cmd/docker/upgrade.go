package docker

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/dockerfile"
	"github.com/alexellis/arkade/pkg/images"
	"github.com/spf13/cobra"
)

func MakeUpgrade() *cobra.Command {
	var command = &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade images in a Dockerfile to the latest version",
		Aliases: []string{"u"},
		Long: `Upgrade container images in a Dockerfile to the latest version.

Only images specified via the --image flag will be upgraded.
Images using variable substitution in tags (e.g. ${VERSION}) are skipped.

Use --pin-major-minor to constrain an image to patch updates within its
current major.minor version. For example, golang:1.24 would upgrade to
1.24.4 but not to 1.25.

An arkade.yaml file can be placed in the same directory as the Dockerfile:

images:
- ghcr.io/openfaas/of-watchdog
- golang

pin_major_minor:
- golang
`,
		Example: `  # Upgrade specific images and write the changes
  arkade docker upgrade --image ghcr.io/openfaas/of-watchdog --write

  # Pin golang to its current major.minor version
  arkade docker upgrade \
    --image ghcr.io/openfaas/of-watchdog \
    --image golang \
    --pin-major-minor golang \
    --verbose

  # Use a different Dockerfile
  arkade docker upgrade -f ./Dockerfile.template --image alpine`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", "Dockerfile", "Path to Dockerfile")
	command.Flags().StringArrayP("image", "i", nil, "Image name to upgrade (specify multiple times)")
	command.Flags().StringArray("pin-major-minor", nil, "Pin image to current major.minor version, only upgrade patch (specify multiple times)")
	command.Flags().BoolP("verbose", "v", true, "Verbose output")
	command.Flags().BoolP("write", "w", true, "Write the updated values back to the file, or stdout when set to false")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		imageNames, _ := cmd.Flags().GetStringArray("image")
		pinnedNames, _ := cmd.Flags().GetStringArray("pin-major-minor")
		verbose, _ := cmd.Flags().GetBool("verbose")
		writeFile, _ := cmd.Flags().GetBool("write")

		basePath := path.Dir(file)
		defaultConfig := path.Join(basePath, "arkade.yaml")
		if _, err := os.Stat(defaultConfig); err == nil {
			cfg, err := config.Load(defaultConfig)
			if err != nil {
				return err
			}
			imageNames = append(imageNames, cfg.Images...)
			pinnedNames = append(pinnedNames, cfg.PinMajorMinor...)
		}

		if len(imageNames) == 0 {
			return fmt.Errorf("specify images to upgrade via --image flag or images list in arkade.yaml")
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		found := dockerfile.FindImages(string(content))
		if len(found) == 0 {
			return fmt.Errorf("no images found in %s", file)
		}

		allowSet := map[string]bool{}
		for _, img := range imageNames {
			allowSet[img] = true
		}

		pinnedSet := map[string]bool{}
		for _, img := range pinnedNames {
			pinnedSet[img] = true
		}

		var toUpdate []dockerfile.ImageRef
		for _, ref := range found {
			if allowSet[ref.Image] {
				toUpdate = append(toUpdate, ref)
			}
		}

		if verbose {
			log.Printf("Found %d images, %d to upgrade\n", len(found), len(toUpdate))
		}

		updatedContent := string(content)
		updateCount := 0

		for _, ref := range toUpdate {
			var updated bool
			var newRef string
			var err error

			if pinnedSet[ref.Image] {
				updated, newRef, err = images.UpdateImagePinned(ref.Ref(), verbose)
			} else {
				updated, newRef, err = images.UpdateImage(ref.Ref(), verbose)
			}

			if err != nil {
				if verbose {
					log.Printf("Warning: %s\n", err)
				}
				continue
			}

			if updated {
				updatedContent = dockerfile.ReplaceImage(updatedContent, ref.Ref(), newRef)
				updateCount++
			}
		}

		if updateCount > 0 && writeFile {
			if err := os.WriteFile(file, []byte(updatedContent), 0600); err != nil {
				return err
			}
			log.Printf("Wrote %d updates to: %s", updateCount, file)
		} else if !writeFile {
			fmt.Print(updatedContent)
		}

		return nil
	}

	return command
}
