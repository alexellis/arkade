module github.com/alexellis/arkade

go 1.25.0

require (
	github.com/Masterminds/semver v1.5.0
	github.com/alexellis/go-execute/v2 v2.2.1
	github.com/docker/go-units v0.5.0
	github.com/google/go-containerregistry v0.20.7
	github.com/mattn/go-isatty v0.0.20
	github.com/morikuni/aec v1.0.0
	github.com/olekukonko/tablewriter v1.1.1
	github.com/otiai10/copy v1.14.1
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.3.1
	github.com/spf13/cobra v1.10.2
	golang.org/x/crypto v0.47.0
	golang.org/x/mod v0.32.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/alexellis/gha-bump v0.0.0

require (
	github.com/alexellis/fstail v0.0.0-20250917111842-2ab578ec2afb
	github.com/clipperhouse/displaywidth v0.3.1 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.18.1 // indirect
	github.com/docker/cli v29.0.3+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.4 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.1.0 // indirect
	github.com/olekukonko/ll v0.1.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/otiai10/mint v1.6.3 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/vbatts/tar-split v0.12.2 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
)

replace github.com/alexellis/fstail => ../fstail

replace github.com/alexellis/gha-bump => ../gha-bump
