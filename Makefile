Version := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
IsDirty := $(if $(shell git status --porcelain),-dirty,)
GitCommit := $(shell git rev-parse HEAD)
BuildTimestamp := $(shell date +%s)
LDFLAGS := "-s -w \
	-X github.com/alexellis/arkade/pkg.Version=$(Version)$(IsDirty) \
	-X github.com/alexellis/arkade/pkg.GitCommit=$(GitCommit) \
	-X github.com/alexellis/arkade/pkg.BuildTimestamp=$(BuildTimestamp)"
PLATFORM := $(shell ./hack/platform-tag.sh)
SOURCE_DIRS = cmd pkg main.go
export GO111MODULE=on

.PHONY: all
all: gofmt test build dist hash

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS)

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ |grep -v vendor/| tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make gofmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: e2e
e2e:
	CGO_ENABLED=0 go test github.com/alexellis/arkade/pkg/get -cover --tags e2e -v

.PHONY: dist-local
dist-local:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o bin/arkade

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o bin/arkade
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -o bin/arkade-darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -o bin/arkade-darwin-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -o bin/arkade-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o bin/arkade-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o bin/arkade.exe

.PHONY: hash
hash:
	rm -rf bin/*.sha256 && ./hack/hashgen.sh
