package get

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"testing"
)

type nopReadCloser struct {
	io.Reader
}

func (n *nopReadCloser) Close() error { return nil }

func TestLineProgressReader_WithTotal(t *testing.T) {
	data := bytes.Repeat([]byte("a"), 100)
	out := bytes.Buffer{}

	r := newLineProgressReader(&nopReadCloser{Reader: bytes.NewReader(data)}, int64(len(data)), &out)

	b := make([]byte, 10)
	for {
		_, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
	}

	if err := r.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Download progress: 10%") {
		t.Fatalf("expected 10%% progress output, got: %s", output)
	}
	if !strings.Contains(output, "Download progress: 100%") {
		t.Fatalf("expected 100%% progress output, got: %s", output)
	}
	if !strings.Contains(output, "Download complete in") {
		t.Fatalf("expected completion output, got: %s", output)
	}
}

func TestLineProgressReader_UnknownTotal(t *testing.T) {
	data := bytes.Repeat([]byte("a"), 32)
	out := bytes.Buffer{}

	r := newLineProgressReader(&nopReadCloser{Reader: bytes.NewReader(data)}, 0, &out)

	b := make([]byte, 8)
	for {
		_, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
	}

	if err := r.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "downloaded") {
		t.Fatalf("expected downloaded output, got: %s", output)
	}
	if !strings.Contains(output, "Download complete in") {
		t.Fatalf("expected completion output, got: %s", output)
	}
}

func TestRenderASCIIBar(t *testing.T) {
	got := renderASCIIBar(50, 10)
	want := "[█████░░░░░]"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestFormatETA(t *testing.T) {
	got := formatETA(100, 50, 10)
	want := "00:05"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestUsePlainProgressOutput_Override(t *testing.T) {
	t.Setenv("ARKADE_PROGRESS_TTY", "true")
	if !usePlainProgressOutput() {
		t.Fatalf("expected plain output when ARKADE_PROGRESS_TTY=true")
	}

	t.Setenv("ARKADE_PROGRESS_TTY", "false")
	if usePlainProgressOutput() {
		t.Fatalf("expected tty output when ARKADE_PROGRESS_TTY=false")
	}
}

func TestUsePlainProgressOutput_DumbTerm(t *testing.T) {
	t.Setenv("ARKADE_PROGRESS_TTY", "")
	_ = os.Unsetenv("ARKADE_PROGRESS_TTY")
	t.Setenv("TERM", "dumb")
	if !usePlainProgressOutput() {
		t.Fatalf("expected plain output for TERM=dumb")
	}
}

func TestCallbackReader_ReportsProgress(t *testing.T) {
	data := bytes.Repeat([]byte("x"), 256)
	total := int64(len(data))

	var lastRead, lastTotal int64
	var callCount int64

	cb := func(bytesRead, totalBytes int64) {
		atomic.StoreInt64(&lastRead, bytesRead)
		atomic.StoreInt64(&lastTotal, totalBytes)
		atomic.AddInt64(&callCount, 1)
	}

	r := newCallbackReader(&nopReadCloser{Reader: bytes.NewReader(data)}, total, cb)

	buf := make([]byte, 32)
	for {
		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
	}

	if err := r.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	if atomic.LoadInt64(&lastRead) != total {
		t.Fatalf("expected final read=%d, got %d", total, atomic.LoadInt64(&lastRead))
	}
	if atomic.LoadInt64(&lastTotal) != total {
		t.Fatalf("expected total=%d reported, got %d", total, atomic.LoadInt64(&lastTotal))
	}
	if atomic.LoadInt64(&callCount) < 2 {
		t.Fatalf("expected at least 2 callbacks (initial + close), got %d", atomic.LoadInt64(&callCount))
	}
}

func TestCallbackReader_NilCallback(t *testing.T) {
	data := []byte("hello world")
	r := newCallbackReader(&nopReadCloser{Reader: bytes.NewReader(data)}, int64(len(data)), nil)

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(out) != string(data) {
		t.Fatalf("data mismatch: got %q", string(out))
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}

func TestArkadeInPath(t *testing.T) {
	t.Setenv("PATH", "/usr/bin:/home/user/.arkade/bin:/usr/local/bin")
	if !ArkadeInPath() {
		t.Fatal("expected ArkadeInPath()=true when .arkade/bin is in PATH")
	}

	t.Setenv("PATH", "/usr/bin:/usr/local/bin")
	if ArkadeInPath() {
		t.Fatal("expected ArkadeInPath()=false when .arkade/bin is not in PATH")
	}
}
