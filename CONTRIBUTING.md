## Contributing

### Guidelines

Here are a few guidelines for contributing:

* If you would like to contribute to the codebase then please raise an issue to propose the change or feature.
* Do not work on an issue / PR until it gets a `design/approved` label from a maintainer.
* Do not mix feature changes or fixes with refactoring - it makes the code harder to review and means there is more for
  the maintainers (with limited time) to test.

* If you have found a bug please raise an issue and fill out the whole template.
* Don't raise PRs for typos, these aren't necessary - just raise an Issue.
* If the documentation can be improved / translated etc. please raise an issue to discuss.

* Please always provide a summary of what you changed, how you did it and how it can be tested.
* Most of the time we like to keep one commit per PR or if you have more you should have a perfect reason for it.
  (That's why we like to amend to a commit and commonly require rebase a multiple commit PR, before we approve/merge
  it.)
* If the Issue involves adding multiple features, like both a Helm chart and a CLI tool, create a PR per tool/chart, to
  cleanly separate them.

All commits must have a `Signed-off-by:` line in accordance with the Developer Certificate of Origin, which you can read
about at the end of this document.

To add the sign-off, simply run:

```bash
git config --global user.name "Full Name"
git config --global user.email "you@example.com"

git commit -s / --signoff
```

We do not expect to contact you in relation to your contributions, however a real, genuine email address must be used which is contactable. A no-reply email address will result in the DCO being marked as invalid by Derek, the DCO bot we use.

This is not cryptography (GPG commit signing), does not require any keys and does not take any longer than typing in the above three
commands.

### Recommended settings for apps for `arkade install`

`arkade install` is primarily meant for use during development and experimentation, however most apps should also have
options available to make them suitable for production.

By default, turn off:

* persistence and volumes
* multiple replicas or HA
* clustering
* additional sidecars and components

But make them configurable, so `--persistence` would be `false` by default for an app like Postgresql, but easy to
enable.

### Recommended way for adding binaries to `arkade get`

For installable CLIs, add the following architectures when available:

**For Darwin (MacOS):**

* AMD64 / `x86_64` (Intel)
* ARM64 / `arm64` (M1)

**For Linux:**

* AMD64 / `x86_64` (Intel/AMD)
* ARM64 / `aarch64` (64-bit Raspberry Pi OS, AWS Graviton and commercial 64-bit ARM servers)
* ARMv7 / `armhf` (32-bit Raspberry Pi OS)

**For Windows:**

* AMD64 (Intel/AMD)

Do not add ARMv6 support, or 32-bit architecture.

FreeBSD and other operating systems are not supported at this time.

Binaries that are packaged as `.zip`, `.tgz` or `.tar.gz` are also supported, look at previous apps to find an example
to copy.

Always write a unit tests for your tool in the `pkg/get/get_test.go`.

In the unit tests always pin an exact version and write a test case for that version.

You can try different architectures and OSes with `arkade get` via `--arch` and `--os` flags, to try out arm64 binaries
for example.

Please always test your tools manually (`hack/test-tool.sh`) and with the end-to-end test tool (`make e2e`). Add the output of these commands to the PR description.

**Available variables in the `URLTemplate`:**

* `OS` (current OS or value of the `--os` flag)
* `Arch` (current architecture or value of the `--arch` flag)
* `Name` (name of the tool)
* `Version` (version of the tool)
* `VersionNumber` (version of the tool with the `v` removed)

Example:

```json
{
  "OS": "darwin",
  "Arch": "arm64",
  "Name": "gh",
  "Version": "v2.14.7",
  "VersionNumber": "2.14.7"
}
```

**Available variables in the `BinaryTemplate`:**

* `OS` (current OS or value of the `--os` flag)
* `Arch` (current architecture or value of the `--arch` flag)
* `Name` (name of the tool)
* `Version` (version of the tool)
* `VersionNumber` (version of the tool with the `v` removed)
* `Repo` (name of the GitHub repository)
* `Owner` (owner of the GitHub repository)

Example:

```json
{
  "OS": "darwin",
  "Arch": "arm64",
  "Name": "gh",
  "Version": "v2.14.7",
  "VersionNumber": "2.14.7",
  "Repo": "cli",
  "Owner": "cli"
}
```

### Workflow for a first-time contributing

Here's the basics for contributing:

#### Cloning

```bash
mkdir $GOPATH/go/src/github.com/alexellis
git clone https://github.com/alexellis/arkade
cd arkade

go build
```

#### Adding your fork:

```bash
git remote add fork https://github.com/NAME/arkade
```

#### To verify changes:

```bash
gofmt -w -s ./pkg
gofmt -w -s ./cmd

# Run all the unit tests
make test

# Use e2e tests ot check that URLs can be downloaded for all tools
make e2e
```

#### Checkout a branch to start work

```bash
git checkout -b fork/add-NAME-of-APP
```

#### Push up your changes for a PR

```bash
git config user.name "Full name"
git config user.email "real@email.com"

git commit -s

git push fork add-NAME-of-APP
```

#### Test other people's PRs

You can also check out other people's PRs and test them:

```bash
arkade get gh
gh auth

$ gh pr list

Showing 10 of 10 open pull requests in alexellis/arkade

#477  Add comma escaping for --set flag        yankeexe:fix-openfaas-helm-set
#438  Added support to install kube-ps1        andreppires:master

gh checkout 477

go build

# Try the new version of arkade
./arkade
```

### Compliance

All commits need to be signed-off in accordance with the Developer Certificate of Origin (DCO) as per below.

The [Derek bot](https://github.com/alexellis/derek) will enforce this policy, if you need help please get in touch.

#### License

This project is licensed under the MIT License.

### Reporting a suspected vulnerability / security issue

If you would like to report a suspected vulnerability / security issue, please email alex@openfaas.com. Bear in mind
that this is a community project, and it may take a few days to get back to you. If you have a working code sample in a
private GitHub repo, please feel free to give access to that also.

#### Sign-off your work

The sign-off is a simple line at the end of the explanation for a patch. Your
signature certifies that you wrote the patch or otherwise have the right to pass
it on as an open-source patch. The rules are pretty simple: if you can certify
the below (from [developercertificate.org](http://developercertificate.org/)):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

    Signed-off-by: Joe Smith <joe.smith@email.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

* Please sign your commits with `git commit -s` so that commits are traceable.
