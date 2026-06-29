package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UntarNested reads the gzip-compressed tar file from r and writes it into dir.
// When allowSymlinks is false, any symlink entry in the archive causes an
// error; when true, symlinks are extracted subject to containment checks.
// When flatExtract is true, all files are extracted directly into dir using
// only their basename, ignoring the archive's directory structure (e.g.
// usr/local/bin/foo -> dir/foo). This is analogous to tar's --strip-components,
// but strips all levels in one go.
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func UntarNested(r io.Reader, dir string, gzipped, quiet, allowSymlinks, flatExtract bool) error {
	return untarNested(r, dir, gzipped, quiet, allowSymlinks, flatExtract)
}

func untarNested(r io.Reader, dir string, gzipped, quiet, allowSymlinks, flatExtract bool) (err error) {
	t0 := time.Now()
	nFiles := 0
	madeDir := map[string]bool{}
	defer func() {
		td := time.Since(t0)
		if err == nil {
			if !quiet {
				log.Printf("extracted tarball into %s: %d files, %d dirs (%v)", dir, nFiles, len(madeDir), td)
			}
		} else {
			log.Printf("error extracting tarball into %s after %d files, %d dirs, %v: %v", dir, nFiles, len(madeDir), td, err)
		}
	}()

	if gzipped {
		zr, err := gzip.NewReader(r)
		if err != nil {
			return fmt.Errorf("requires gzip-compressed body: %v", err)
		}
		r = zr
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Resolve dir to its real path so containment checks are not confused by a
	// symlinked install directory (e.g. /usr/local/bin on some systems).
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}
	dir = resolvedDir
	cleanDir := filepath.Clean(dir)

	tr := tar.NewReader(r)

	loggedChtimesError := false
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("tar reading error: %v", err)
			return fmt.Errorf("tar error: %v", err)
		}
		if !validRelPath(f.Name) {
			return fmt.Errorf("tar contained invalid name error %q", f.Name)
		}
		name := f.Name
		if flatExtract {
			name = filepath.Base(name)
		}
		rel := filepath.FromSlash(name)
		abs := filepath.Join(dir, rel)

		fi := f.FileInfo()
		mode := fi.Mode()
		if !quiet {
			fmt.Printf("Extracting: %s\n", abs)
		}
		switch {
		case mode.IsRegular():
			parent := filepath.Dir(abs)
			if !madeDir[parent] {
				// Guard before MkdirAll: it follows a pre-existing symlink and
				// would otherwise create directories outside root.
				if err := assertExistingPrefixWithinRoot(cleanDir, parent); err != nil {
					return err
				}
				if err := os.MkdirAll(parent, 0755); err != nil {
					return err
				}
				madeDir[parent] = true
			}
			// Resolve the physical parent (containment already guaranteed above)
			// to locate the write, allowing write-through of internal symlinks.
			resolvedParent, err := filepath.EvalSymlinks(parent)
			if err != nil {
				return fmt.Errorf("cannot resolve parent of %s: %v", abs, err)
			}
			abs = filepath.Join(resolvedParent, filepath.Base(abs))
			// Don't write through a pre-existing symlink at the leaf; O_CREATE
			// would follow it outside root.
			if fi, err := os.Lstat(abs); err == nil && fi.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("refusing to write through symlink %q", abs)
			}
			wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}
			n, err := io.Copy(wf, tr)
			if closeErr := wf.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("error writing to %s: %v", abs, err)
			}
			if n != f.Size {
				return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, f.Size)
			}
			modTime := f.ModTime
			if modTime.After(t0) {
				// Clamp modtimes at system time. See
				// golang.org/issue/19062 when clock on
				// buildlet was behind the gitmirror server
				// doing the git-archive.
				modTime = t0
			}
			if !modTime.IsZero() {
				if err := os.Chtimes(abs, modTime, modTime); err != nil && !loggedChtimesError {
					// benign error. Gerrit doesn't even set the
					// modtime in these, and we don't end up relying
					// on it anywhere (the gomote push command relies
					// on digests only), so this is a little pointless
					// for now.
					log.Printf("error changing modtime: %v (further Chtimes errors suppressed)", err)
					loggedChtimesError = true // once is enough
				}
			}
			nFiles++
		case mode.IsDir():
			if flatExtract {
				continue
			}
			// Guard before MkdirAll, as with regular files.
			if err := assertExistingPrefixWithinRoot(cleanDir, abs); err != nil {
				return err
			}
			if err := os.MkdirAll(abs, 0755); err != nil {
				return err
			}
			madeDir[abs] = true
		case mode.Type() == os.ModeSymlink:
			if flatExtract {
				log.Printf("skipping symlink %q during flat extraction (target may not resolve correctly)", f.Name)
				continue
			}
			if !allowSymlinks {
				return fmt.Errorf("tar file entry %s is a symlink, but symlink extraction is disabled", f.Name)
			}
			parent := filepath.Dir(abs)
			if !madeDir[parent] {
				if err := assertExistingPrefixWithinRoot(cleanDir, parent); err != nil {
					return err
				}
				if err := os.MkdirAll(parent, 0755); err != nil {
					return err
				}
				madeDir[parent] = true
			}
			// Resolve the physical parent (containment already guaranteed above)
			// to locate where the symlink is created.
			resolvedParent, err := filepath.EvalSymlinks(parent)
			if err != nil {
				return fmt.Errorf("cannot resolve parent of %s: %v", abs, err)
			}
			abs = filepath.Join(resolvedParent, filepath.Base(abs))
			// Validate the link target stays within root. resolvedParent is
			// symlink-free, so this lexical check matches the physical location.
			target := f.Linkname
			if !filepath.IsAbs(target) {
				target = filepath.Join(resolvedParent, target)
			}
			if !inDir(filepath.Clean(target), cleanDir) {
				return fmt.Errorf("refusing symlink %q -> %q (escapes %q)", abs, f.Linkname, dir)
			}
			// ...and physically (pre-existing symlink in the target path).
			if err := assertExistingPrefixWithinRoot(cleanDir, target); err != nil {
				return err
			}
			if err := os.Symlink(f.Linkname, abs); err != nil {
				return err
			}
		default:
			return fmt.Errorf("tar file entry %s contained unsupported file type %v", f.Name, mode)
		}
	}
	return nil
}

// assertExistingPrefixWithinRoot resolves the longest existing ancestor of p and
// returns an error if it does not stay within root.
func assertExistingPrefixWithinRoot(root, p string) error {
	cur := filepath.Clean(p)
	for {
		if _, err := os.Lstat(cur); err == nil {
			break
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			// Reached the filesystem root without an existing component.
			return nil
		}
		cur = parent
	}
	resolved, err := filepath.EvalSymlinks(cur)
	if err != nil {
		return err
	}
	if !inDir(filepath.Clean(resolved), root) {
		return fmt.Errorf("refusing to create %q: existing path %q resolves to %q outside %q", p, cur, resolved, root)
	}
	return nil
}

// inDir reports whether path is equal to root or a direct descendant.
// Both arguments must be clean paths (output of filepath.Clean or filepath.EvalSymlinks).
func inDir(path, root string) bool {
	return path == root || strings.HasPrefix(path, root+string(os.PathSeparator))
}
