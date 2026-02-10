// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

func MakeRelease() *cobra.Command {
	var command = &cobra.Command{
		Use:   "release [version] [title]",
		Short: "Create a GitHub release",
		Long: `Create a GitHub release by querying the latest release tag,
incrementing the version, and using the latest commit message as the
release title.

The latest release tag is obtained via "gh release list", so private
repositories are supported when "gh" is authenticated.

Version bump strategy (mutually exclusive, default is --patch):
  --patch   v0.1.2 -> v0.1.3
  --minor   v0.1.2 -> v0.2.0
  --major   v0.1.2 -> v1.0.0

Semver handling:
  - Tags such as "v0.1.0" (prefixed) and "0.1.0" (bare) are both recognised.
  - The "v" prefix is preserved if the previous tag used one.

Repository visibility changes the defaults:
  - Public repos:  --prerelease=true  --latest=false
    A separate bot promotes pre-releases to latest after verification.
  - Private repos: --prerelease=false --latest=true
    Releases are immediately marked as latest and stable.

Visibility is auto-detected via "gh repo view". Use --prerelease or --latest
to override.

Both version and title can be overridden with positional arguments. If only
one argument is given it is treated as a version when it looks like semver,
otherwise it is used as the title. Supply both to set both.
When a version is given explicitly, --major/--minor/--patch are ignored.`,
		Example: `  # Auto-detect everything (version, title, visibility)
  # Bumps patch by default: v0.1.2 -> v0.1.3
  arkade release

  # Bump minor version: v0.1.2 -> v0.2.0
  arkade release --minor

  # Bump major version: v0.1.2 -> v1.0.0
  arkade release --major

  # Public repo: creates a pre-release, not marked latest
  arkade release

  # Private repo: creates a full release, marked latest
  arkade rel

  # Override the version (--major/--minor/--patch ignored)
  arkade rel v0.5.0

  # Override just the title
  arkade r "Fix CI pipeline"

  # Override both version and title
  arkade release 0.3.0 "Add ARM64 support"

  # Force a full (non-prerelease) release on a public repo
  arkade release --prerelease=false --latest=true

  # Dry run to see what would happen
  arkade release --dry-run`,
		Hidden:       true,
		Aliases:      []string{"rel", "r"},
		SilenceUsage: true,
	}

	command.Flags().Bool("dry-run", false, "Print what would be created without making the release")
	command.Flags().Bool("major", false, "Bump the major version (v1.2.3 -> v2.0.0)")
	command.Flags().Bool("minor", false, "Bump the minor version (v1.2.3 -> v1.3.0)")
	command.Flags().Bool("patch", false, "Bump the patch version (v1.2.3 -> v1.2.4) (default when none specified)")
	command.Flags().String("prerelease", "", "Mark as pre-release (default: true for public repos, false for private)")
	command.Flags().String("latest", "", "Mark as latest release (default: false for public repos, true for private)")
	command.Flags().String("token", "", "GitHub token, if not using gh's default auth")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		var versionOverride, titleOverride string

		if len(args) == 1 {
			if looksLikeSemver(args[0]) {
				versionOverride = args[0]
			} else {
				titleOverride = args[0]
			}
		} else if len(args) == 2 {
			versionOverride = args[0]
			titleOverride = args[1]
		} else if len(args) > 2 {
			return fmt.Errorf("expected at most 2 arguments: [version] [title], got %d", len(args))
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		major, _ := cmd.Flags().GetBool("major")
		minor, _ := cmd.Flags().GetBool("minor")
		patch, _ := cmd.Flags().GetBool("patch")
		prereleaseStr, _ := cmd.Flags().GetString("prerelease")
		latestStr, _ := cmd.Flags().GetString("latest")
		token, _ := cmd.Flags().GetString("token")

		bump, err := parseBumpStrategy(major, minor, patch)
		if err != nil {
			return err
		}

		ctx := context.Background()

		// Detect repo visibility to set defaults
		private, err := repoIsPrivate(ctx)
		if err != nil {
			fmt.Printf("Warning: could not detect repo visibility: %s\nDefaulting to public repo behaviour (pre-release).\n", err)
			private = false
		}

		prerelease := !private // public -> true, private -> false
		if prereleaseStr != "" {
			prerelease = prereleaseStr == "true"
		}

		latest := private // public -> false, private -> true
		if latestStr != "" {
			latest = latestStr == "true"
		}

		if private {
			fmt.Println("Private repo detected.")
		} else {
			fmt.Println("Public repo detected.")
		}

		// Resolve the new version tag
		newTag, err := resolveNewTag(ctx, versionOverride, bump)
		if err != nil {
			return err
		}

		// Resolve the release title
		title, err := resolveTitle(ctx, titleOverride)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("Would create release:\n  Tag:        %s\n  Title:      %s\n  Bump:       %s\n  Notes:      (none)\n  Prerelease: %v\n  Latest:     %v\n", newTag, title, bump, prerelease, latest)
			return nil
		}

		ghArgs := []string{
			"release", "create", newTag,
			"--title", title,
			"--notes", "",
		}

		if prerelease {
			ghArgs = append(ghArgs, "--prerelease")
		}

		if latest {
			ghArgs = append(ghArgs, "--latest")
		} else {
			ghArgs = append(ghArgs, "--latest=false")
		}

		task := execute.ExecTask{
			Command:      "gh",
			Args:         ghArgs,
			StreamStdio:  true,
			PrintCommand: true,
		}

		if token != "" {
			task.Env = []string{fmt.Sprintf("GH_TOKEN=%s", token)}
		}

		res, err := task.Execute(ctx)
		if err != nil {
			return fmt.Errorf("failed to run gh: %w", err)
		}
		if res.ExitCode != 0 {
			return fmt.Errorf("gh exited with code %d: %s", res.ExitCode, res.Stderr)
		}

		kind := "release"
		if prerelease {
			kind = "pre-release"
		}
		fmt.Printf("Created %s %s.\n", kind, newTag)
		return nil
	}

	return command
}

// bumpStrategy represents which semver component to increment.
type bumpStrategy int

const (
	bumpPatchStrategy bumpStrategy = iota
	bumpMinorStrategy
	bumpMajorStrategy
)

func (b bumpStrategy) String() string {
	switch b {
	case bumpMinorStrategy:
		return "minor"
	case bumpMajorStrategy:
		return "major"
	default:
		return "patch"
	}
}

// parseBumpStrategy validates the --major/--minor/--patch flags are mutually
// exclusive and returns the selected strategy. Defaults to patch.
func parseBumpStrategy(major, minor, patch bool) (bumpStrategy, error) {
	set := 0
	if major {
		set++
	}
	if minor {
		set++
	}
	if patch {
		set++
	}
	if set > 1 {
		return 0, fmt.Errorf("--major, --minor, and --patch are mutually exclusive")
	}

	if major {
		return bumpMajorStrategy, nil
	}
	if minor {
		return bumpMinorStrategy, nil
	}
	return bumpPatchStrategy, nil
}

// repoIsPrivate uses gh to check whether the current repo is private.
func repoIsPrivate(ctx context.Context) (bool, error) {
	task := execute.ExecTask{
		Command: "gh",
		Args:    []string{"api", "repos/{owner}/{repo}", "--jq", ".visibility"},
	}

	res, err := task.Execute(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to run gh: %w", err)
	}
	if res.ExitCode != 0 {
		return false, fmt.Errorf("gh exited with code %d: %s", res.ExitCode, res.Stderr)
	}

	visibility := strings.TrimSpace(strings.ToUpper(res.Stdout))
	return visibility == "PRIVATE" || visibility == "INTERNAL", nil
}

// resolveNewTag returns the tag for the new release. If versionOverride is
// non-empty it is normalised and returned directly. Otherwise the latest
// release tag is fetched via gh and bumped according to the strategy.
func resolveNewTag(ctx context.Context, versionOverride string, bump bumpStrategy) (string, error) {
	if versionOverride != "" {
		return normaliseSemver(versionOverride)
	}

	latestTag, err := latestReleaseTag(ctx)
	if err != nil {
		return "", fmt.Errorf("could not determine latest release: %w\nUse a positional argument to set the version explicitly", err)
	}

	return bumpVersion(latestTag, bump)
}

// resolveTitle returns a release title. If titleOverride is non-empty it is
// used directly. Otherwise the first line of the latest git commit message is
// returned.
func resolveTitle(ctx context.Context, titleOverride string) (string, error) {
	if titleOverride != "" {
		return titleOverride, nil
	}

	task := execute.ExecTask{
		Command: "git",
		Args:    []string{"log", "-1", "--format=%s"},
	}

	res, err := task.Execute(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get latest commit message: %w", err)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git log exited with code %d: %s", res.ExitCode, res.Stderr)
	}

	msg := strings.TrimSpace(res.Stdout)
	if msg == "" {
		return "", fmt.Errorf("latest commit message is empty")
	}

	return msg, nil
}

// latestReleaseTag queries gh for the most recent release tag.
func latestReleaseTag(ctx context.Context) (string, error) {
	task := execute.ExecTask{
		Command: "gh",
		Args:    []string{"api", "repos/{owner}/{repo}/releases", "--jq", ".[0].tag_name"},
	}

	res, err := task.Execute(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run gh: %w", err)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("gh exited with code %d: %s", res.ExitCode, res.Stderr)
	}

	tag := strings.TrimSpace(res.Stdout)
	if tag == "" {
		return "", fmt.Errorf("no releases found in this repository")
	}

	return tag, nil
}

// bumpVersion increments the specified semver component of a tag.
// Both "v"-prefixed and bare versions are supported; the prefix style is
// preserved in the output.
func bumpVersion(tag string, bump bumpStrategy) (string, error) {
	hasV := strings.HasPrefix(tag, "v")

	ver, err := semver.NewVersion(tag)
	if err != nil {
		return "", fmt.Errorf("cannot parse %q as semver: %w", tag, err)
	}

	var bumped semver.Version
	switch bump {
	case bumpMajorStrategy:
		bumped = ver.IncMajor()
	case bumpMinorStrategy:
		bumped = ver.IncMinor()
	default:
		bumped = ver.IncPatch()
	}

	if hasV {
		return "v" + bumped.String(), nil
	}
	return bumped.String(), nil
}

// normaliseSemver validates that s looks like semver and returns it unchanged.
func normaliseSemver(s string) (string, error) {
	raw := s
	if strings.HasPrefix(s, "v") {
		s = s[1:]
	}
	if _, err := semver.NewVersion(s); err != nil {
		return "", fmt.Errorf("%q is not a valid semver version", raw)
	}
	return raw, nil
}

// looksLikeSemver returns true when s resembles a semver string
// (with or without a "v" prefix). This is intentionally lenient so that
// users can pass e.g. "1.0" and have it treated as a version.
func looksLikeSemver(s string) bool {
	raw := s
	if strings.HasPrefix(raw, "v") || strings.HasPrefix(raw, "V") {
		raw = raw[1:]
	}

	// Quick heuristic: starts with a digit and contains a dot.
	if len(raw) == 0 {
		return false
	}
	if raw[0] < '0' || raw[0] > '9' {
		return false
	}
	return strings.Contains(raw, ".")
}
