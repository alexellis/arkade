package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"github.com/ulikunitz/xz"
)

func untarXZ(r io.Reader, dir string, quiet bool) (err error) {
	t0 := time.Now()
	nFiles := 0
	madeDir := map[string]bool{}
	defer func() {
		td := time.Since(t0)
		if err == nil {
			if !quiet {
				log.Printf("extracted tar.xz into %s: %d files, %d dirs (%v)", dir, nFiles, len(madeDir), td)
			}
		} else {
			log.Printf("error extracting tar.xz into %s after %d files, %d dirs, %v: %v", dir, nFiles, len(madeDir), td, err)
		}
	}()

	xzr, err := xz.NewReader(r)
	if err != nil {
		return fmt.Errorf("requires xz-compressed body: %v", err)
	}

	tr := tar.NewReader(xzr)
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
	baseFile := filepath.Base(f.Name)
	abs := filepath.Join(dir, baseFile)
		if !quiet {
			fmt.Printf("Extracting: %s to\t%s\n", f.Name, abs)
		}

		fi := f.FileInfo()
		mode := fi.Mode()
		switch {
		case mode.IsDir():
			break
		case mode.IsRegular():
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
				modTime = t0
			}
			if !modTime.IsZero() {
				if err := os.Chtimes(abs, modTime, modTime); err != nil && !loggedChtimesError {
					log.Printf("error changing modtime: %v (further Chtimes errors suppressed)", err)
					loggedChtimesError = true
				}
			}
			nFiles++
		default:
		}
	}
	return nil
}
