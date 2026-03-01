package dockerfile

import (
	"testing"
)

func TestFindImages(t *testing.T) {
	content := `ARG PYTHON_VERSION=3.12

FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:0.11.3 AS watchdog
FROM --platform=${TARGETPLATFORM:-linux/amd64} python:${PYTHON_VERSION}-alpine AS build

COPY --from=watchdog /fwatchdog /usr/bin/fwatchdog

FROM build AS test
FROM build AS ship`

	images := FindImages(content)

	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d: %v", len(images), images)
	}

	if images[0].Image != "ghcr.io/openfaas/of-watchdog" {
		t.Fatalf("expected image ghcr.io/openfaas/of-watchdog, got %s", images[0].Image)
	}

	if images[0].Tag != "0.11.3" {
		t.Fatalf("expected tag 0.11.3, got %s", images[0].Tag)
	}
}

func TestFindImagesMultiple(t *testing.T) {
	content := `FROM alpine:3.23.0
FROM ghcr.io/openfaas/of-watchdog:0.11.3 AS watchdog`

	images := FindImages(content)

	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}

	want := []ImageRef{
		{Image: "alpine", Tag: "3.23.0"},
		{Image: "ghcr.io/openfaas/of-watchdog", Tag: "0.11.3"},
	}

	for i, w := range want {
		if images[i].Image != w.Image || images[i].Tag != w.Tag {
			t.Fatalf("image[%d]: want %s:%s, got %s:%s", i, w.Image, w.Tag, images[i].Image, images[i].Tag)
		}
	}
}

func TestFindImagesSkipsVariables(t *testing.T) {
	content := `FROM python:${PYTHON_VERSION}-alpine
FROM golang:${GO_VERSION}`

	images := FindImages(content)

	if len(images) != 0 {
		t.Fatalf("expected 0 images (variables should be skipped), got %d", len(images))
	}
}

func TestFindImagesSkipsNoTag(t *testing.T) {
	content := `FROM scratch
FROM build`

	images := FindImages(content)

	if len(images) != 0 {
		t.Fatalf("expected 0 images, got %d", len(images))
	}
}

func TestFindImagesDeduplicates(t *testing.T) {
	content := `FROM alpine:3.23.0 AS builder
FROM alpine:3.23.0 AS runtime`

	images := FindImages(content)

	if len(images) != 1 {
		t.Fatalf("expected 1 image (deduplicated), got %d", len(images))
	}
}

func TestFindImagesRegistryWithPort(t *testing.T) {
	content := `FROM myregistry.com:5000/myimage:1.2.3`

	images := FindImages(content)

	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}

	if images[0].Image != "myregistry.com:5000/myimage" {
		t.Fatalf("expected image myregistry.com:5000/myimage, got %s", images[0].Image)
	}

	if images[0].Tag != "1.2.3" {
		t.Fatalf("expected tag 1.2.3, got %s", images[0].Tag)
	}
}

func TestReplaceImage(t *testing.T) {
	content := `FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:0.11.3 AS watchdog
FROM --platform=${TARGETPLATFORM:-linux/amd64} python:${PYTHON_VERSION}-alpine AS build`

	result := ReplaceImage(content,
		"ghcr.io/openfaas/of-watchdog:0.11.3",
		"ghcr.io/openfaas/of-watchdog:0.11.4")

	expected := `FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:0.11.4 AS watchdog
FROM --platform=${TARGETPLATFORM:-linux/amd64} python:${PYTHON_VERSION}-alpine AS build`

	if result != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestReplaceImagePreservesNonFromLines(t *testing.T) {
	content := `# Comment about alpine:3.23.0
FROM alpine:3.23.0 AS builder
RUN echo "alpine:3.23.0"`

	result := ReplaceImage(content, "alpine:3.23.0", "alpine:3.23.3")

	expected := `# Comment about alpine:3.23.0
FROM alpine:3.23.3 AS builder
RUN echo "alpine:3.23.0"`

	if result != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestReplaceImageDuplicateFromLines(t *testing.T) {
	content := `FROM alpine:3.23.0 AS builder
FROM alpine:3.23.0 AS runtime`

	result := ReplaceImage(content, "alpine:3.23.0", "alpine:3.23.3")

	expected := `FROM alpine:3.23.3 AS builder
FROM alpine:3.23.3 AS runtime`

	if result != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestImageRefRef(t *testing.T) {
	ref := ImageRef{Image: "alpine", Tag: "3.23.0"}
	if ref.Ref() != "alpine:3.23.0" {
		t.Fatalf("expected alpine:3.23.0, got %s", ref.Ref())
	}
}

func TestSplitRefWithPort(t *testing.T) {
	image, tag, ok := splitRef("myregistry.com:5000/myimage:1.2.3")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if image != "myregistry.com:5000/myimage" {
		t.Fatalf("expected myregistry.com:5000/myimage, got %s", image)
	}
	if tag != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %s", tag)
	}
}

func TestSplitRefSimple(t *testing.T) {
	image, tag, ok := splitRef("alpine:3.23.0")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if image != "alpine" {
		t.Fatalf("expected alpine, got %s", image)
	}
	if tag != "3.23.0" {
		t.Fatalf("expected 3.23.0, got %s", tag)
	}
}

func TestSplitRefNoTag(t *testing.T) {
	_, _, ok := splitRef("scratch")
	if ok {
		t.Fatal("expected ok=false for image without tag")
	}
}
