package archive

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type tarEntry struct {
	hdr  tar.Header
	body []byte
}

func buildTar(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	var b bytes.Buffer
	tw := tar.NewWriter(&b)
	for _, e := range entries {
		h := e.hdr
		if len(e.body) > 0 {
			h.Size = int64(len(e.body))
		}
		if err := tw.WriteHeader(&h); err != nil {
			t.Fatalf("write header %q: %v", h.Name, err)
		}
		if len(e.body) > 0 {
			if _, err := tw.Write(e.body); err != nil {
				t.Fatalf("write body %q: %v", h.Name, err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	return b.Bytes()
}

// A symlink whose target is an absolute path outside the install dir must be rejected,
// even when a subsequent entry attempts to write through it.
func Test_UntarNested_RejectsAbsoluteSymlinkWriteThrough(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}
	outsideDir := filepath.Join(baseDir, "outside")
	if err := os.MkdirAll(outsideDir, 0755); err != nil {
		t.Fatal(err)
	}

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "escape-link", Typeflag: tar.TypeSymlink, Linkname: outsideDir, Mode: 0777}},
		{hdr: tar.Header{Name: "escape-link/escape.txt", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("escaped\n")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatal("want error, got nil")
	}

	if content, err := os.ReadFile(filepath.Join(outsideDir, "escape.txt")); err == nil {
		t.Fatalf("file written outside install dir: content=%q", string(content))
	}
}

// A symlink whose target is a relative path that escapes the install dir must be rejected,
// even when a subsequent entry attempts to write through it.
func Test_UntarNested_RejectsRelativeSymlinkWriteThrough(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}
	outsideDir := filepath.Join(baseDir, "outside")
	if err := os.MkdirAll(outsideDir, 0755); err != nil {
		t.Fatal(err)
	}

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "escape-link", Typeflag: tar.TypeSymlink, Linkname: "../outside", Mode: 0777}},
		{hdr: tar.Header{Name: "escape-link/escape.txt", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("escaped\n")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatal("want error, got nil")
	}

	if content, err := os.ReadFile(filepath.Join(outsideDir, "escape.txt")); err == nil {
		t.Fatalf("file written outside install dir: content=%q", string(content))
	}
}

// A chain of symlinks that appears valid lexically but escapes the install dir
// at runtime must be rejected; hop1 -> "." (resolves to install), hop1/hop2 -> ".." (escapes to base).
func Test_UntarNested_RejectsChainedSymlinkEscape(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "hop1", Typeflag: tar.TypeSymlink, Linkname: ".", Mode: 0777}},
		{hdr: tar.Header{Name: "hop1/hop2", Typeflag: tar.TypeSymlink, Linkname: "..", Mode: 0777}},
		{hdr: tar.Header{Name: "hop1/hop2/outside/escape.txt", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("escaped\n")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatal("want error, got nil")
	}

	if content, err := os.ReadFile(filepath.Join(baseDir, "outside", "escape.txt")); err == nil {
		t.Fatalf("file written outside install dir: content=%q", string(content))
	}
}

// A symlink that resolves outside the install dir must not be left on disk,
// even without a subsequent write; hop1 -> "." (resolves to install), hop1/hop2 -> ".." (escapes to base).
func Test_UntarNested_RejectsPlantedEscapingSymlink(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "hop1", Typeflag: tar.TypeSymlink, Linkname: ".", Mode: 0777}},
		{hdr: tar.Header{Name: "hop1/hop2", Typeflag: tar.TypeSymlink, Linkname: "..", Mode: 0777}},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatal("want error, got nil")
	}

	planted := filepath.Join(installDir, "hop2")
	if fi, err := os.Lstat(planted); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		target, _ := os.Readlink(planted)
		t.Fatalf("escaping symlink left on disk: %s -> %q", planted, target)
	}
}

// A symlink whose target traverses a pre-existing symlink (inside the install dir
// but pointing outside) must be rejected and not left on disk. The lexical target
// check alone passes here, so this exercises the physical-prefix resolution.
func Test_UntarNested_RejectsSymlinkTargetViaPreExistingSymlink(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}
	outsideDir := filepath.Join(baseDir, "outside")
	if err := os.MkdirAll(outsideDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Pre-existing symlink inside the install dir that points outside.
	if err := os.Symlink(outsideDir, filepath.Join(installDir, "safe")); err != nil {
		t.Fatal(err)
	}

	// "safe/file" is lexically under installDir, but "safe" resolves outside.
	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "planted", Typeflag: tar.TypeSymlink, Linkname: "safe/file", Mode: 0777}},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatalf("expected extraction to be rejected, got nil error")
	}

	planted := filepath.Join(installDir, "planted")
	if fi, err := os.Lstat(planted); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		target, _ := os.Readlink(planted)
		t.Fatalf("escaping symlink left on disk: %s -> %q", planted, target)
	}
}

// Ordinary nested directories and files must extract correctly.
func Test_UntarNested_AllowsValidNestedFiles(t *testing.T) {
	installDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(installDir)

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "bin", Typeflag: tar.TypeDir, Mode: 0755}},
		{hdr: tar.Header{Name: "bin/tool", Typeflag: tar.TypeReg, Mode: 0755}, body: []byte("#!/bin/sh\n")},
		{hdr: tar.Header{Name: "README.md", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("hello\n")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err != nil {
		t.Fatalf("expected clean extraction, got: %v", err)
	}
	for _, rel := range []string{"bin/tool", "README.md"} {
		if _, err := os.Stat(filepath.Join(installDir, rel)); err != nil {
			t.Fatalf("expected %q to exist: %v", rel, err)
		}
	}
}

// A file written through an internal symlink (one whose target stays within the install dir)
// must land at the symlink's target location.
func Test_UntarNested_AllowsWriteThroughInternalSymlink(t *testing.T) {
	installDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(installDir)

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "subdir", Typeflag: tar.TypeDir, Mode: 0755}},
		{hdr: tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "subdir", Mode: 0777}},
		{hdr: tar.Header{Name: "link/file.txt", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("hello\n")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err != nil {
		t.Fatalf("expected clean extraction with write through internal symlink, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(installDir, "subdir", "file.txt")); err != nil {
		t.Fatalf("expected file to exist at symlink target: %v", err)
	}
}

// A directory created through an internal symlink must land at the symlink's target location.
func Test_UntarNested_AllowsWriteThroughInternalSymlinkDir(t *testing.T) {
	installDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(installDir)

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "subdir", Typeflag: tar.TypeDir, Mode: 0755}},
		{hdr: tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "subdir", Mode: 0777}},
		{hdr: tar.Header{Name: "link/newdir", Typeflag: tar.TypeDir, Mode: 0755}},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err != nil {
		t.Fatalf("expected clean extraction with dir write-through internal symlink, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(installDir, "subdir", "newdir")); err != nil {
		t.Fatalf("expected directory to exist at symlink target: %v", err)
	}
}

// A pre-existing symlink inside the extraction root that points outside must not cause
// MkdirAll to create directories outside the root, even though no file is written there.
func Test_UntarNested_PreExistingSymlinkDoesNotCreateDirOutsideRoot(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}
	outsideDir := filepath.Join(baseDir, "outside")
	if err := os.MkdirAll(outsideDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Plant a symlink inside the extraction root pointing outside — pre-existing, not from tar.
	if err := os.Symlink(outsideDir, filepath.Join(installDir, "link")); err != nil {
		t.Fatal(err)
	}

	// Tar only contains a regular file whose parent traverses the pre-existing symlink.
	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "link/subdir/escape.txt", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("escaped\n")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatal("want error, got nil")
	}

	// MkdirAll must not follow the symlink and create outsideDir/subdir before EvalSymlinks catches the escape.
	if _, err := os.Stat(filepath.Join(outsideDir, "subdir")); err == nil {
		t.Fatal("directory created outside extraction root via pre-existing symlink")
	}
}

// A pre-existing symlink at the final path component must not be written through;
// a regular-file entry of the same name must not redirect the write outside root.
func Test_UntarNested_RejectsLeafSymlinkWriteThrough(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	installDir := filepath.Join(baseDir, "install")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatal(err)
	}
	outsideTarget := filepath.Join(baseDir, "target.txt")
	if err := os.WriteFile(outsideTarget, []byte("ORIGINAL"), 0644); err != nil {
		t.Fatal(err)
	}
	// Plant a leaf symlink inside the root pointing at a file outside the root.
	if err := os.Symlink(outsideTarget, filepath.Join(installDir, "evil")); err != nil {
		t.Fatal(err)
	}

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "evil", Typeflag: tar.TypeReg, Mode: 0644}, body: []byte("HACKED")},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err == nil {
		t.Fatal("want error, got nil")
	}

	if b, _ := os.ReadFile(outsideTarget); string(b) != "ORIGINAL" {
		t.Fatalf("file outside root overwritten through leaf symlink: now %q", string(b))
	}
}

// A symlink whose target stays within the install dir must be created.
func Test_UntarNested_AllowsValidInternalSymlink(t *testing.T) {
	installDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(installDir)

	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "tool-v1", Typeflag: tar.TypeReg, Mode: 0755}, body: []byte("bin\n")},
		{hdr: tar.Header{Name: "tool", Typeflag: tar.TypeSymlink, Linkname: "tool-v1", Mode: 0777}},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err != nil {
		t.Fatalf("expected clean extraction with internal symlink, got: %v", err)
	}
	linkPath := filepath.Join(installDir, "tool")
	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("expected symlink %q to exist: %v", linkPath, err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected %q to be a symlink", linkPath)
	}
}

// A symlink entry whose parent directory is not listed as its own entry in the tar
// must still be created; the parent directory is made on demand.
func Test_UntarNested_CreatesParentDirForSymlink(t *testing.T) {
	installDir, err := os.MkdirTemp("", "arkade-untar-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(installDir)

	// No "nested" directory entry precedes the symlink; the target stays within root.
	data := buildTar(t, []tarEntry{
		{hdr: tar.Header{Name: "nested/link", Typeflag: tar.TypeSymlink, Linkname: "tool-v1", Mode: 0777}},
	})

	if err := UntarNested(bytes.NewReader(data), installDir, false, true); err != nil {
		t.Fatalf("expected clean extraction with on-demand parent dir for symlink, got: %v", err)
	}
	linkPath := filepath.Join(installDir, "nested", "link")
	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("expected symlink %q to exist: %v", linkPath, err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected %q to be a symlink", linkPath)
	}
}
