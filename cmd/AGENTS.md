# AGENTS.md - Guide for AI Agents working on `arkade release`

This document covers the hidden `release` command (`cmd/release.go`), which creates GitHub releases from any repo using `gh`.

## Overview

The command is registered in `main.go` as a hidden top-level command with aliases `rel` and `r`. It is not shown in `arkade --help` output.

**Source file:** `cmd/release.go`
**Registration:** `main.go` via `rootCmd.AddCommand(cmd.MakeRelease())`

## What the command does

1. Detects repo visibility (public/private) via `gh api repos/{owner}/{repo}`
2. Queries the latest release tag via `gh api repos/{owner}/{repo}/releases`
3. Parses the tag as semver and increments the patch version
4. Gets the latest commit subject line via `git log` for the release title
5. Creates the release via `gh release create`

No release notes are generated. A separate webhook bot handles release note generation.

## Public vs Private repo rules

This is the most important design decision in the command.

| | Public repo | Private / Internal repo |
|---|---|---|
| `--prerelease` | `true` (default) | `false` (default) |
| `--latest` | `false` (default) | `true` (default) |

**Why:** Public repos use a bot that promotes pre-releases to latest after verification. Private repos have no such bot, so releases are immediately marked as latest and stable.

Visibility is auto-detected via `gh api repos/{owner}/{repo}`. The result is compared against `PRIVATE` and `INTERNAL` (both treated as private). Everything else (including `PUBLIC`) uses public-repo defaults.

Both flags can be overridden explicitly:
```bash
arkade release --prerelease=false --latest=true
```

**When modifying this logic, never change the defaults without understanding the bot promotion workflow for public repos.**

## Semver handling

- Both `v`-prefixed (`v0.1.0`) and bare (`0.1.0`) tags are supported
- The prefix style from the previous release is preserved in the new tag
- The bump strategy is controlled by `--major`, `--minor`, `--patch` flags (mutually exclusive)
- Default is `--patch` when none are specified
- When a version is given as a positional argument, the bump flags are ignored
- The `github.com/Masterminds/semver` library is used (already vendored, same as `cmd/chart/bump.go`)

### Bump strategies

| Flag | Example |
|---|---|
| `--patch` (default) | `v1.2.3` -> `v1.2.4` |
| `--minor` | `v1.2.3` -> `v1.3.0` |
| `--major` | `v1.2.3` -> `v2.0.0` |

The `bumpVersion()` function handles all three via the `bumpStrategy` type. The `parseBumpStrategy()` function validates mutual exclusion of the flags.

## Positional arguments

The command accepts 0, 1, or 2 positional arguments:

| Arguments | Behaviour |
|---|---|
| None | Version auto-detected, title from latest commit |
| One, looks like semver | Used as version override, title from latest commit |
| One, not semver | Used as title override, version auto-detected |
| Two | First is version, second is title |

The `looksLikeSemver()` heuristic checks: optional `v`/`V` prefix, starts with a digit, contains a dot.

## External dependencies

The command shells out to two tools:

- **`gh`** (GitHub CLI) - for release listing, repo visibility, and release creation. Must be authenticated (`gh auth login`). Supports private repos.
- **`git`** - for reading the latest commit message.

Both are called via `github.com/alexellis/go-execute/v2`, not `os/exec` directly. This is consistent with the rest of the arkade codebase.

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--dry-run` | bool | `false` | Print what would be created, don't call `gh release create` |
| `--major` | bool | `false` | Bump the major version |
| `--minor` | bool | `false` | Bump the minor version |
| `--patch` | bool | `false` | Bump the patch version (this is the default when none specified) |
| `--prerelease` | string | `""` (auto) | `"true"` or `"false"`. Auto-detected from repo visibility when empty |
| `--latest` | string | `""` (auto) | `"true"` or `"false"`. Auto-detected from repo visibility when empty |
| `--token` | string | `""` | GitHub token passed via `GH_TOKEN` env var to `gh` |

Note: `--prerelease` and `--latest` are string flags (not bool) so that the empty string means "not set by user, use auto-detected default". This is intentional.

Note: `--major`, `--minor`, `--patch` are mutually exclusive. If more than one is set, the command returns an error.

## Testing

There are no unit tests for this command. To verify behaviour:

```bash
# Build
go build -o arkade .

# Dry run in a public repo (default = patch bump)
./arkade release --dry-run
# Expected: Prerelease: true, Latest: false, Bump: patch

# Minor bump
./arkade release --dry-run --minor

# Major bump
./arkade release --dry-run --major

# Mutual exclusion error
./arkade release --dry-run --major --minor
# Expected: error

# Dry run with overrides
./arkade release --dry-run --prerelease=false --latest=true

# Test argument parsing
./arkade rel --dry-run v1.0.0
./arkade r --dry-run "Some title"
./arkade release --dry-run 1.0.0 "Some title"

# Verify aliases
./arkade rel --help
./arkade r --help
```

## Modifying the command

### Adding new flags
Add them in `MakeRelease()` before the `command.RunE` assignment. Use the same pattern as existing flags.

### Changing visibility detection
The `repoIsPrivate()` function is self-contained. It returns `true` for `PRIVATE` and `INTERNAL` visibility values from `gh`.

### Changing the version bump strategy
The `bumpVersion()` function handles version incrementing using the `bumpStrategy` type. It supports `bumpPatchStrategy`, `bumpMinorStrategy`, and `bumpMajorStrategy`. The user selects the strategy via `--major`, `--minor`, `--patch` flags, validated by `parseBumpStrategy()`.

### Error handling for repos with no releases
When `gh release list` returns no tags, `latestReleaseTag()` returns an error. The error message tells the user to pass a version explicitly as a positional argument.
