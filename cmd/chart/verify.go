package chart

import (
	"fmt"
	"os"
	"path"
	"text/tabwriter"

	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"

	"github.com/google/go-containerregistry/pkg/crane"
)

func MakeVerify() *cobra.Command {
	var command = &cobra.Command{
		Use:   "verify",
		Short: "Verify images from a values.yaml file exist on the remote registry",
		Long: `Verify images in a values.yaml file exist within a remote registry.
Container images must be specified at the top level, or one level down in the 
"image: " or "component.image: " field in a values.yaml file.

Returns exit code zero if all images were found on the remote registry.

Otherwise, it returns a non-zero exit code and a table of images not found:

COMPONENT           IMAGE
dashboard           ghcr.io/openfaasltd/openfaas-dashboard:0.9.8
autoscaler          ghcr.io/openfaasltd/autoscaler:0.2.5

`,
		Example: `  chartctl verify -f ./chart/values.yaml
  chartctl verify --verbose -f ./chart/values.yaml`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", "", "Path to values.yaml file")
	command.Flags().BoolP("verbose", "v", false, "Verbose output")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("invalid value for flag --file")
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		if len(file) == 0 {
			return fmt.Errorf("flag --file is required")
		}

		if ext := path.Ext(file); ext != ".yaml" && ext != ".yml" {
			return fmt.Errorf("--file must be a YAML file")
		}

		if verbose {
			fmt.Printf("Verifying images in: %s\n", file)
		}

		values, err := helm.Load(file)
		if err != nil {
			return err
		}

		filtered := helm.FilterImages(values)
		if len(filtered) == 0 {
			return fmt.Errorf("no images found in %s", file)
		}

		missed := []verifyError{}
		if verbose {
			if len(filtered) > 0 {
				fmt.Printf("Found %d images\n", len(filtered))
			}
		}
		for k, v := range filtered {
			if verbose {
				fmt.Printf("> [%s] %s\n", k, v)
			}

			ref, err := crane.Head(v)
			if err != nil {
				missed = append(missed, verifyError{
					Err:       err,
					Image:     v,
					Component: k,
				})
			} else {
				if verbose {
					fmt.Printf("< [%s] %v\n", v, ref.Digest)
				}
			}
		}

		if len(missed) > 0 {
			fmt.Fprintf(os.Stderr, "%d images are missing in %s\n\n", len(missed), file)

			w := tabwriter.NewWriter(os.Stderr, 20, 3, 1, ' ', 0)

			if verbose {
				fmt.Fprintf(w, "COMPONENT\tIMAGE\tERROR\n")
			} else {
				fmt.Fprintf(w, "COMPONENT\tIMAGE\n")
			}

			for _, err := range missed {
				if verbose {
					fmt.Fprintf(w, "%s\t%s\t%s\n", err.Component, err.Image, err.Err)
				} else {
					fmt.Fprintf(w, "%s\t%s\n", err.Component, err.Image)
				}

				w.Flush()
			}

			fmt.Fprintf(os.Stderr, "\n")

			return fmt.Errorf("verifying failed")
		}

		return nil
	}

	return command
}

type verifyError struct {
	Err       error
	Image     string
	Component string
}
