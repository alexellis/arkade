# AGENTS.md - Guide for AI Agents Contributing to Arkade

This document provides guidance for AI agents working on arkade, specifically for reviewing and adding new CLI tools to the `arkade get` command.

## Types of Arkade Commands

Arkade provides several types of installers:

- **`arkade get`** - CLI tools usually to be placed at `/usr/local/bin/` or `$HOME/.arkade/bin/`. These are standalone binaries that can be downloaded and executed directly.

- **`arkade system install`** - Linux-only system-level tools like Node.js, Go, Prometheus. These require additional installation steps or system configuration.

- **`arkade oci install`** - Fetches binaries out of OCI images. Ideal for projects that use private repositories like slicer/actuated/k3sup-pro.

- **`arkade install`** - Kubernetes Helm charts or manifests for add-ons like OpenFaaS CE, Istio, PostgreSQL. These deploy software to Kubernetes clusters.

**This guide focuses on `arkade get`** - adding CLI tools that provide static binaries for download.

## 1. How to Add a New CLI (Tool) to Arkade

### What Can Be Added

**Only tools with static binaries** can be added to arkade. The tool must provide pre-compiled binaries for download.

**Cannot be added:**
- Python-based tools (e.g., `aws-cli`, `azure-cli`) - require Python runtime
- Node.js-based tools without static binaries - require Node.js runtime
- Tools that require installation scripts or package managers
- Tools that need runtime dependencies beyond the binary itself

### Prerequisites

1. Fork and create a branch: `git checkout -b add-TOOL_NAME`

### Step 1: Check GitHub Releases

**CRITICAL**: Before writing code, check the latest stable release on GitHub to see what OS/architecture combinations are available.

1. Run a `curl -i -X HEAD https://github.com/OWNER/REPO/releases/latest` (adds `/latest` to go directly to the latest release) - change OWNER and REPO accordingly. The `location` header in the response will show the actual latest version tag without using up an API request. To obtain the `location` header, you must not use the `-L` (follow redirects) flag.
2. Examine ALL download URLs in the "Assets" section, you can obtain this via HTML, again to avoid consuming API requests: `https://github.com/OWNER/REPO/releases/expanded_assets/VERSION` - replace VERSION with the actual version tag from 1. and the OWNER/REPO accordingly. This returns HTML, you can grep it efficiently for anchor tags.
3. Note available combinations:
   - Linux amd64 (x86_64)
   - Linux arm64 (aarch64)
   - Darwin amd64 (x86_64)
   - Darwin arm64
   - Windows amd64 (x86_64)

**Important**: Match the exact naming used by the upstream project (`amd64` vs `x86_64`, `arm64` vs `aarch64`).

### Step 2: Add Tool Definition

Edit `pkg/get/tools.go` and add a new `Tool` entry. **Reference existing examples** like `faas-cli` (lines 27-50) for the structure.

**Key points:**
- Use `BinaryTemplate` for GitHub releases (simpler)
- Use `URLTemplate` for custom URLs or non-GitHub sources
- Supported archive formats: `.tar.gz`, `.zip` (`.tar.xz` is NOT supported)
- Template variables: `.OS`, `.Arch`, `.Name`, `.Version`, `.VersionNumber`, `.Repo`, `.Owner`
- Windows detection: `HasPrefix .OS "ming"`
- **CRITICAL**: If a binary is missing for a specific OS/arch (e.g., Windows amd64), the template must still generate a URL that results in a 404 error, NOT download the wrong binary (e.g., don't download Linux binary when Windows was requested)

#### Archive tools: when the binary name inside the archive differs from the tool name

When a tool is distributed as an archive (`.tar.gz`, `.tgz`, `.zip`) and the **binary inside the archive** has a platform-specific name (e.g., `mytool-darwin-arm64` rather than just `mytool`), you **must** use both `URLTemplate` and `BinaryTemplate` together:

- **`URLTemplate`** — the full download URL including the archive extension (e.g., `https://github.com/.../mytool-darwin-arm64.tgz`)
- **`BinaryTemplate`** — the name of the **binary inside the archive**, without the archive extension (e.g., `mytool-darwin-arm64`)

**Do NOT** put the archive filename (with `.tgz`/`.tar.gz`/`.zip` extension) in `BinaryTemplate` alone. The `decompress()` function in `pkg/get/download.go` uses `BinaryTemplate` to locate the extracted binary. If `BinaryTemplate` contains an archive extension, decompress falls back to `tool.Name` which will be wrong when the inner binary has a platform suffix.

**Reference example**: `inletsctl` in `pkg/get/tools.go` — uses `URLTemplate` for the download URL and `BinaryTemplate` for the inner binary name.

**When `BinaryTemplate` alone is safe**: Only when the tool is a **plain binary** (not an archive). In that case `BinaryTemplate` is the release asset filename, and the downloaded file is used directly without decompression.

### Step 3: Write Unit Tests

Add a test function in `pkg/get/get_test.go`. **Reference `Test_DownloadFaasCli`** (around line 2761) as an example.

**Requirements:**
- Use a pinned version (not "latest")
- Test all available OS/arch combinations
- Verify URLs match actual GitHub release URLs

### Step 4: Download and Verify Every OS/Arch Combination

First test downloading the current OS/arch (no flags needed) i.e. `go run . get TOOL`. And run `file` on the output to verify the type and if it's valid or invalid i.e. a gzip, or a HTML page, or got a non zero exit code.

**MANDATORY**: Download and verify EVERY combination using the `file` command.

```bash
# Build arkade
go build

# Test all combinations (script automates this)
./hack/test-tool.sh TOOL_NAME
```

For each combination, verify the `file` command output:
- Linux amd64: `ELF 64-bit LSB executable, x86-64`
- Linux arm64: `ELF 64-bit LSB executable, ARM aarch64`
- Darwin amd64: `Mach-O 64-bit x86_64 executable`
- Darwin arm64: `Mach-O 64-bit arm64 executable`
- Windows amd64: `PE32+ executable (console) x86-64`

Tools built with Rust often have `unknown` in their filename, that's OK. If deciding between GNU aka libc or musl, pick the musl version.

**Include the full output of `./hack/test-tool.sh TOOL_NAME` in your PR description.**

### Step 5: Update Documentation

The README.md file contains instructions for updating itself. Follow the note at the bottom of the "Catalog of CLIs" section: run `go run . get --format markdown` to generate the updated table, then replace the existing catalog section. Write it to a file in the workspace that you delete after, to avoid needing extra permissions.

There are two tokens in the README.md - replace al text between them with what you've generated.

Start of replaceable block is inside: `<!-- start of tool list -->` and the end is inside: `<!-- end of tool list -->`.


### Step 6: Create Pull Request

**PR Description must include:**
- List of available/unavailable OS/arch combinations from GitHub releases page
- Full output from `./hack/test-tool.sh TOOL_NAME` showing `file` command results
- Output from `make e2e` (if applicable)

**Checklist:**
- [ ] All commits signed off (`git commit -s`)
- [ ] Unit tests pass
- [ ] All OS/arch combinations verified with `file` command
- [ ] README.md updated
- [ ] PR description includes verification output

### Architecture Support Reference

| OS | Architecture | Const name | Notes |
|---|---|---|---|
| macOS (Intel) | x86_64 | `arch64bit` | Intel Macs |
| macOS (Apple Silicon) | arm64 | `archDarwinARM64` | M1/M2/M3 Macs |
| Linux | x86_64 | `arch64bit` | Standard Linux |
| Linux | aarch64/arm64 | `archARM64` | ARM64 Linux |
| Windows | x86_64 | `arch64bit` | Windows (Git Bash) |

**Note**: Do not add ARMv6 or 32-bit x86 support.

### Troubleshooting

- **URLs don't match**: Check actual release URLs on GitHub and adjust template
- **Wrong architecture in binary**: Verify binary names on GitHub releases page
- **Missing combinations**: Document why in PR description if upstream doesn't provide them. The template must still generate a URL that returns 404 (not download the wrong binary)
- **Downloads wrong binary**: If requesting Windows but getting Linux binary, the template is incorrectly falling back. Each OS/arch must have a unique URL that matches the actual release or returns 404
- **"stat ... no such file or directory" after extraction**: The binary name inside the archive doesn't match what `decompress()` expects. This happens when `BinaryTemplate` alone contains an archive extension (`.tgz`, `.tar.gz`, `.zip`) — the code falls back to `tool.Name` instead of the platform-specific binary name. Fix by splitting into `URLTemplate` (download URL) + `BinaryTemplate` (inner binary name without extension). See the "Archive tools" section above.

---

## 2. How to Review a New CLI Being Added as an AI Agent

### Pre-Review Checklist

- [ ] Issue has `design/approved` label
- [ ] All commits signed off
- [ ] PR adds only one tool

### Code Review

#### Tool Definition (`pkg/get/tools.go`)

- [ ] Tool provides static binaries (not Python/Node.js-based)
- [ ] Required fields: `Name`, `Owner`, `Repo`, `Description`
- [ ] Either `BinaryTemplate` or `URLTemplate` provided
- [ ] Supports required OS/arch combinations (Linux amd64/arm64, Darwin amd64/arm64, Windows amd64)
- [ ] Archive format is `.tar.gz` or `.zip` (not `.tar.xz`)
- [ ] Missing OS/arch combinations generate URLs that return 404 (not download wrong binary)

#### Unit Tests (`pkg/get/get_test.go`)

- [ ] Test function exists with pinned version
- [ ] Test cases for all available platforms
- [ ] URLs match actual GitHub release URLs

#### Documentation

- [ ] README.md updated with tool entry

#### PR Description

- [ ] **CRITICAL**: Includes `file` command output for **every** OS/arch combination
- [ ] Documents which OS/arch combinations are available from upstream
- [ ] Includes output from `./hack/test-tool.sh TOOL_NAME`
- [ ] Includes output from `make e2e` (if applicable)

### Critical Review: Binary Verification

**MANDATORY**: Verify `file` command output for every combination shows correct architecture:
- Linux amd64: `ELF 64-bit LSB executable, x86-64`
- Linux arm64: `ELF 64-bit LSB executable, ARM aarch64`
- Darwin amd64: `Mach-O 64-bit x86_64 executable`
- Darwin arm64: `Mach-O 64-bit arm64 executable`
- Windows amd64: `PE32+ executable (console) x86-64`

**If missing, request it before approving.**

### Review Commands

```bash
go build && ./hack/test-tool.sh TOOL_NAME
go test ./pkg/get/... -v
./arkade get --format markdown | grep TOOL_NAME
```

### Common Issues

1. Tool requires runtime (Python/Node.js) - cannot be added
2. Missing `file` command output for all combinations
3. URLs don't match actual GitHub releases
4. Missing architecture support
5. Wrong architecture mapping (`arm64` vs `aarch64`, `amd64` vs `x86_64`)
6. Using unsupported archive format (`.tar.xz`)
7. Template downloads wrong binary when combination is missing (e.g., downloads Linux when Windows requested) - must return 404 instead

---

## Reference Examples

- **Simple BinaryTemplate**: `faas-cli` (lines 27-50 in `pkg/get/tools.go`)
- **Test example**: `Test_DownloadFaasCli` (around line 2761 in `pkg/get/get_test.go`)
- **Recent additions**: `dufs` (commit a120f8c), `logcli` (commit 4f72efe), `ripgrep` (commit a80f284)

## Additional Resources

- [CONTRIBUTING.md](CONTRIBUTING.md) - General contribution guidelines
- [.github/PULL_REQUEST_TEMPLATE.md](.github/PULL_REQUEST_TEMPLATE.md) - PR template
