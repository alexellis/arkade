package get

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func packageArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, contents := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0644, Size: int64(len(contents))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(contents)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestPostInstallationMsgPackage(t *testing.T) {
	for _, movePath := range []string{"", "/tmp/tools"} {
		msg, err := PostInstallationMsg(movePath, []ToolLocal{
			{Name: "codex", Path: "/tmp/tools/codex", Package: true},
			{Name: "jq", Path: "/tmp/tools/jq"},
		})
		if err != nil {
			t.Fatal(err)
		}
		text := string(msg)
		if strings.Contains(text, "sudo mv /tmp/tools/codex") || strings.Contains(text, "sudo install -m 755 /tmp/tools/codex") {
			t.Fatalf("instructions would detach the package entry point: %s", text)
		}
		if !strings.Contains(text, "sudo arkade get codex --path /usr/local/bin") {
			t.Fatalf("missing package installation instructions: %s", text)
		}
		ordinary := "sudo mv /tmp/tools/jq /usr/local/bin/"
		if movePath != "" {
			ordinary = "sudo install -m 755 /tmp/tools/jq /usr/local/bin/jq"
		}
		if !strings.Contains(text, ordinary) || strings.Contains(text, "sudo arkade get jq") {
			t.Fatalf("incorrect ordinary tool installation instructions: %s", text)
		}
	}
}

func TestDownloadPackage(t *testing.T) {
	for _, targetOS := range []string{"linux", "darwin", "mingw"} {
		t.Run(targetOS, func(t *testing.T) {
			ext := ""
			if targetOS == "mingw" {
				ext = ".exe"
			}
			files := map[string]string{
				"bin/tool" + ext: "main v1", "bin/helper" + ext: "helper v1",
				"resources/nested/data": "resource", "manifest.json": "manifest",
			}
			body := packageArchive(t, files)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(body)
			}))
			defer server.Close()
			tool := Tool{Name: "tool", URLTemplate: server.URL + "/package.tar.gz", BinaryTemplate: "tool", PackageBinaries: []string{"tool", "helper"}}
			dir := t.TempDir()
			// Also exercise migration from the original single-binary install.
			if err := os.WriteFile(filepath.Join(dir, "tool"+ext), []byte("old standalone"), 0755); err != nil {
				t.Fatal(err)
			}
			for _, version := range []string{"v1", "v2"} {
				files["bin/tool"+ext] = "main " + version
				body = packageArchive(t, files)
				got, _, err := Download(&tool, "amd64", targetOS, version, dir, false, true, false)
				if err != nil {
					t.Fatal(err)
				}
				if got != filepath.Join(dir, "tool"+ext) {
					t.Fatalf("unexpected path: %s", got)
				}
				for _, name := range []string{"tool", "helper"} {
					p := filepath.Join(dir, name+ext)
					info, err := os.Stat(p)
					if err != nil || info.Mode().Perm() != 0755 {
						t.Fatalf("executable %s: info=%v err=%v", p, info, err)
					}
					content, err := os.ReadFile(p)
					if err != nil || string(content) != files["bin/"+name+ext] {
						t.Fatalf("binary %s: %q, %v", name, content, err)
					}
				}
			}
			for _, name := range []string{"resources/nested/data", "manifest.json"} {
				content, err := os.ReadFile(filepath.Join(dir, ".arkade-tool", name))
				if err != nil || string(content) != files[name] {
					t.Fatalf("resource %s: %q, %v", name, content, err)
				}
			}
			// A missing companion must fail before replacing a working install.
			delete(files, "bin/helper"+ext)
			body = packageArchive(t, files)
			if _, _, err := Download(&tool, "amd64", targetOS, "v3", dir, false, true, false); err == nil {
				t.Fatal("expected missing companion error")
			}
			content, err := os.ReadFile(filepath.Join(dir, "tool"+ext))
			if err != nil || string(content) != "main v2" {
				t.Fatalf("previous install damaged: %q, %v", content, err)
			}
			staging, err := filepath.Glob(filepath.Join(dir, ".arkade-package-*"))
			if err != nil || len(staging) != 0 {
				t.Fatalf("staging directories left behind: %v, %v", staging, err)
			}
		})
	}
}
