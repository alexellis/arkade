package docker

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/dockerfile"
	"github.com/spf13/cobra"
)

func MakeGen() *cobra.Command {
	var command = &cobra.Command{
		Use:   "gen",
		Short: "Generate an arkade.yaml from detected images in a Dockerfile",
		Long: `Generate an arkade.yaml from detected images in a Dockerfile.

This command scans a Dockerfile and generates an arkade.yaml file
with the detected images. You can then edit this file and use it
with 'arkade docker upgrade' to automatically upgrade images.

Only images with explicit tags are detected. Images using variable
substitution (e.g., ${VERSION}) are skipped.
`,
		Example: `  # Generate arkade.yaml from current directory's Dockerfile
  arkade docker gen

  # Generate from a specific Dockerfile
  arkade docker gen -f ./Dockerfile.prod

  # Output to stdout only (don't create file)
  arkade docker gen --stdout
  arkade docker gen -f Dockerfile | tee arkade.yaml
`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", "Dockerfile", "Path to Dockerfile")
	command.Flags().BoolP("stdout", "s", false, "Output to stdout instead of creating file")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		toStdout, _ := cmd.Flags().GetBool("stdout")

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		found := dockerfile.FindImages(string(content))
		if len(found) == 0 {
			return fmt.Errorf("no images found in %s", file)
		}

		// Sort images for consistent output
		images := make([]string, len(found))
		for i, img := range found {
			images[i] = img.Image
		}

		var yamlBuilder strings.Builder
		yamlBuilder.WriteString("images:\n")
		for _, img := range images {
			yamlBuilder.WriteString(fmt.Sprintf("- %s\n", img))
		}

		output := yamlBuilder.String()

		if toStdout {
			fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		}

		basePath := path.Dir(file)
		outputPath := path.Join(basePath, "arkade.yaml")

		if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outputPath, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Generated %s with %d images:\n", outputPath, len(images))
		for _, img := range images {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", img)
		}

		return nil
	}

	return command
}
