//go:build unix

package oci

import (
	"bytes"
	"syscall"
	"testing"
)

func TestTarDirSkipsSpecialFile(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(dir+"/ok.txt", "x"); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(dir+"/pipe", 0600); err != nil {
		t.Skipf("mkfifo unavailable: %s", err)
	}
	var buf bytes.Buffer
	if err := tarDir(dir, &buf); err != nil {
		t.Fatalf("tarDir should not error or hang on a FIFO: %s", err)
	}
	if bytesContains(buf.Bytes(), "pipe") {
		t.Fatal("want FIFO to be skipped, not written to the tar")
	}
}
