package get

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	units "github.com/docker/go-units"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/mattn/go-isatty"
)

// ProgressCallback is called by the download layer to report bytes
// transferred and total size so the caller can render progress.
type ProgressCallback func(bytesRead, totalBytes int64)

const (
	DownloadTempDir   = iota
	DownloadArkadeDir = iota
)

type ErrNotFound struct {
}

func (e *ErrNotFound) Error() string {
	return "server returned status: 404"
}

// callbackReader wraps an io.ReadCloser and calls a ProgressCallback
// on every Read so the caller can track bytes transferred.
type callbackReader struct {
	r         io.ReadCloser
	total     int64
	read      int64
	callback  ProgressCallback
	lastFlush time.Time
}

func newCallbackReader(r io.ReadCloser, total int64, cb ProgressCallback) io.ReadCloser {
	// Fire an initial event so the caller knows the total size immediately.
	if cb != nil {
		cb(0, total)
	}
	return &callbackReader{r: r, total: total, callback: cb}
}

func (c *callbackReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	if n > 0 {
		c.read += int64(n)
		now := time.Now()
		// Throttle callbacks to at most every 50ms to avoid overwhelming
		// the renderer, but always fire on the first and last read.
		if c.callback != nil && (now.Sub(c.lastFlush) >= 50*time.Millisecond || err == io.EOF) {
			c.callback(c.read, c.total)
			c.lastFlush = now
		}
	}
	return n, err
}

func (c *callbackReader) Close() error {
	// Final flush so the caller sees 100%.
	if c.callback != nil {
		c.callback(c.read, c.total)
	}
	return c.r.Close()
}

// DownloadWithProgress is like Download but accepts an optional
// ProgressCallback for the HTTP transfer phase. When cb is non-nil
// the built-in progress bar is suppressed and the callback receives
// byte counts instead, letting the caller render progress centrally.
func DownloadWithProgress(tool *Tool, arch, operatingSystem, version string, movePath string, quiet, verify bool, cb ProgressCallback) (string, string, error) {
	return downloadTool(tool, arch, operatingSystem, version, movePath, false, quiet, verify, cb)
}

func Download(tool *Tool, arch, operatingSystem, version string, movePath string, displayProgress, quiet, verify bool) (string, string, error) {
	return downloadTool(tool, arch, operatingSystem, version, movePath, displayProgress, quiet, verify, nil)
}

func downloadTool(tool *Tool, arch, operatingSystem, version string, movePath string, displayProgress, quiet, verify bool, cb ProgressCallback) (string, string, error) {

	downloadURL, resolvedVersion, err := GetDownloadURL(tool,
		strings.ToLower(operatingSystem),
		strings.ToLower(arch),
		version, quiet)
	if err != nil {
		return "", "", err
	}

	if !quiet {
		log.Printf("Downloading: %s", downloadURL)
	}

	start := time.Now()

	// When a ProgressCallback is provided the caller owns the display,
	// so we suppress the built-in per-file progress bar.
	var outFilePath string
	if cb != nil {
		outFilePath, err = downloadFileWithCallback(downloadURL, cb)
	} else {
		outFilePath, err = downloadFile(downloadURL, displayProgress)
	}
	if err != nil {
		return "", "", err
	}

	if !quiet {
		filename := path.Base(downloadURL)
		stat, err := os.Stat(outFilePath)
		size := ""
		if err == nil {
			size = "(" + units.HumanSize(float64(stat.Size())) + ")"
		}
		log.Printf("Downloaded %s %s in %s.", filename, size, time.Since(start).Round(time.Millisecond))
	}

	if verify {
		if tool.VerifyStrategy == ClaudeShasumStrategy {
			st := time.Now()
			tmpl := template.New(tool.Name + "sha")
			tmpl = tmpl.Funcs(templateFuncs)
			t, err := tmpl.Parse(tool.VerifyTemplate)
			if err != nil {
				return "", "", err
			}

			var buf bytes.Buffer
			inputs := map[string]string{
				"Name":          tool.Name,
				"Owner":         tool.Owner,
				"Repo":          tool.Repo,
				"Version":       resolvedVersion,
				"VersionNumber": strings.TrimPrefix(resolvedVersion, "v"),
				"Arch":          arch,
				"OS":            operatingSystem,
			}

			if err = t.Execute(&buf, inputs); err != nil {
				return "", "", err
			}

			verifyURL := strings.TrimSpace(buf.String())
			log.Printf("Downloading SHA sum from: %s", verifyURL)
			shaSumManifest, err := fetchText(verifyURL)
			if err != nil {
				return "", "", err
			}

			var manifest struct {
				Version   string `json:"version"`
				BuildDate string `json:"buildDate"`
				Platforms map[string]struct {
					Checksum string `json:"checksum"`
					Size     int64  `json:"size"`
				} `json:"platforms"`
			}
			if err := json.Unmarshal([]byte(shaSumManifest), &manifest); err != nil {
				return "", "", err
			}

			var archMappingForClaude = map[string]string{
				"amd64":   "amd64",
				"x86_64":  "x64",
				"arm64":   "arm64",
				"aarch64": "arm64",
			}

			platformKey := fmt.Sprintf("%s-%s", strings.ToLower(operatingSystem), archMappingForClaude[arch])

			platformInfo, found := manifest.Platforms[platformKey]
			if !found {
				return "", "", fmt.Errorf("no checksum info found for platform: %s", platformKey)
			}

			if err := verifySHA(platformInfo.Checksum, outFilePath); err != nil {
				return "", "", err
			} else {
				log.Printf("SHA sum verified in %s.", time.Since(st).Round(time.Millisecond))
			}

		} else if tool.VerifyStrategy == HashicorpShasumStrategy {
			st := time.Now()
			tmpl := template.New(tool.Name + "sha")
			tmpl = tmpl.Funcs(templateFuncs)
			t, err := tmpl.Parse(tool.VerifyTemplate)
			if err != nil {
				return "", "", err
			}

			var buf bytes.Buffer
			inputs := map[string]string{
				"Name":          tool.Name,
				"Owner":         tool.Owner,
				"Repo":          tool.Repo,
				"Version":       resolvedVersion,
				"VersionNumber": strings.TrimPrefix(resolvedVersion, "v"),
				"Arch":          arch,
				"OS":            operatingSystem,
			}

			if err = t.Execute(&buf, inputs); err != nil {
				return "", "", err
			}

			verifyURL := strings.TrimSpace(buf.String())
			log.Printf("Downloading SHA sum from: %s", verifyURL)
			shaSum, err := fetchText(verifyURL)
			if err != nil {
				return "", "", err
			}
			if err := verifySHA(shaSum, outFilePath); err != nil {
				return "", "", err
			} else {
				log.Printf("SHA sum verified in %s.", time.Since(st).Round(time.Millisecond))
			}
		} else if tool.VerifyStrategy == AmpShasumStrategy {
			st := time.Now()
			tmpl := template.New(tool.Name + "sha")
			tmpl = tmpl.Funcs(templateFuncs)
			t, err := tmpl.Parse(tool.VerifyTemplate)
			if err != nil {
				return "", "", err
			}

			var buf bytes.Buffer
			inputs := map[string]string{
				"Name":          tool.Name,
				"Owner":         tool.Owner,
				"Repo":          tool.Repo,
				"Version":       resolvedVersion,
				"VersionNumber": strings.TrimPrefix(resolvedVersion, "v"),
				"Arch":          arch,
				"OS":            operatingSystem,
			}

			if err = t.Execute(&buf, inputs); err != nil {
				return "", "", err
			}

			verifyURL := strings.TrimSpace(buf.String())
			log.Printf("Downloading SHA sum from: %s", verifyURL)
			shaSum, err := fetchText(verifyURL)
			if err != nil {
				return "", "", err
			}
			if err := verifySHA(shaSum, outFilePath); err != nil {
				return "", "", err
			} else {
				log.Printf("SHA sum verified in %s.", time.Since(st).Round(time.Millisecond))
			}
		}
	}

	if isArchiveStr(downloadURL) {

		outPath, err := decompress(tool, downloadURL, outFilePath, operatingSystem, arch, version, quiet)
		if err != nil {
			return "", "", err
		}

		outFilePath = outPath
		if v, ok := os.LookupEnv("ARK_DEBUG"); ok && v == "1" {
			log.Printf("Extracted: %s", outFilePath)
		}
	}

	finalName := tool.Name
	if strings.Contains(strings.ToLower(operatingSystem), "mingw") && !tool.NoExtension {
		finalName = finalName + ".exe"
	}

	var localPath string

	if movePath == "" {
		_, err := config.InitUserDir()
		if err != nil {
			return "", "", err
		}

		localPath = env.LocalBinary(finalName, "")
	} else {
		localPath = filepath.Join(movePath, finalName)
	}

	if v, ok := os.LookupEnv("ARK_DEBUG"); ok && v == "1" {
		log.Printf("Copying %s to %s\n", outFilePath, localPath)
	}

	if _, err = CopyFile(outFilePath, localPath); err != nil {
		return "", "", err
	}

	// Remove parent folder of the binary
	tempPath := filepath.Dir(outFilePath)
	if err := os.RemoveAll(tempPath); err != nil {
		log.Printf("Error removing temporary directory: %s", err)
	}

	outFilePath = localPath

	return outFilePath, finalName, nil
}

// DownloadFile downloads a file to a temporary directory
// and returns the path to the file and any error.
func DownloadFileP(downloadURL string, displayProgress bool) (string, error) {
	return downloadFile(downloadURL, displayProgress)
}

func downloadFile(downloadURL string, displayProgress bool) (string, error) {
	return retryWithBackoff(func() (string, error) {
		req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("User-Agent", pkg.UserAgent())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode == http.StatusNotFound {
			return "", &ErrNotFound{}
		}

		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("server returned status: %d", res.StatusCode)
		}

		_, fileName := path.Split(downloadURL)
		tmp := os.TempDir()

		customTmp, err := os.MkdirTemp(tmp, "arkade-*")
		if err != nil {
			return "", err
		}

		outFilePath := path.Join(customTmp, fileName)
		wrappedReader := withProgressBar(res.Body, int(res.ContentLength), displayProgress)

		// Owner/Group read/write/execute
		// World - execute
		out, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
		if err != nil {
			return "", err
		}

		defer out.Close()
		defer wrappedReader.Close()

		if _, err := io.Copy(out, wrappedReader); err != nil {
			return "", err
		}

		return outFilePath, nil
	}, 10, 100*time.Millisecond)
}

// downloadFileWithCallback is like downloadFile but wraps the response
// body with a callbackReader so the caller gets byte-level progress.
func downloadFileWithCallback(downloadURL string, cb ProgressCallback) (string, error) {
	return retryWithBackoff(func() (string, error) {
		req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("User-Agent", pkg.UserAgent())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode == http.StatusNotFound {
			return "", &ErrNotFound{}
		}

		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("server returned status: %d", res.StatusCode)
		}

		_, fileName := path.Split(downloadURL)
		tmp := os.TempDir()

		customTmp, err := os.MkdirTemp(tmp, "arkade-*")
		if err != nil {
			return "", err
		}

		outFilePath := path.Join(customTmp, fileName)
		wrappedReader := newCallbackReader(res.Body, res.ContentLength, cb)

		out, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
		if err != nil {
			return "", err
		}

		defer out.Close()
		defer wrappedReader.Close()

		if _, err := io.Copy(out, wrappedReader); err != nil {
			return "", err
		}

		return outFilePath, nil
	}, 10, 100*time.Millisecond)
}

func CopyFile(src, dst string) (int64, error) {
	return CopyFileP(src, dst, 0700)
}

func CopyFileP(src, dst string, permMode int) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(permMode))
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func withProgressBar(r io.ReadCloser, length int, displayProgress bool) io.ReadCloser {
	if !displayProgress {
		return r
	}

	if usePlainProgressOutput() {
		return newLineProgressReader(r, int64(length), os.Stderr)
	}

	return newTTYProgressReader(r, int64(length), os.Stderr)
}

func usePlainProgressOutput() bool {
	if envValue, ok := os.LookupEnv("ARKADE_PROGRESS_TTY"); ok {
		v, err := strconv.ParseBool(envValue)
		if err == nil {
			return v
		}
	}

	isTTY := isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
	if !isTTY {
		return true
	}

	term := strings.TrimSpace(strings.ToLower(os.Getenv("TERM")))
	return term == "" || term == "dumb"
}

type lineProgressReader struct {
	r            io.ReadCloser
	total        int64
	read         int64
	nextPercent  int64
	lastPercent  int64
	lastReported time.Time
	started      time.Time
	out          io.Writer
}

type ttyProgressReader struct {
	r          io.ReadCloser
	total      int64
	read       int64
	out        io.Writer
	started    time.Time
	lastDraw   time.Time
	frameIndex int
	printedHdr bool
}

func newTTYProgressReader(r io.ReadCloser, total int64, out io.Writer) io.ReadCloser {
	return &ttyProgressReader{
		r:       r,
		total:   total,
		out:     out,
		started: time.Now(),
	}
}

func (t *ttyProgressReader) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	if n > 0 {
		t.read += int64(n)
		t.draw(false)
	}
	return n, err
}

func (t *ttyProgressReader) Close() error {
	t.draw(true)
	fmt.Fprint(t.out, "\n")
	return t.r.Close()
}

func (t *ttyProgressReader) draw(final bool) {
	now := time.Now()
	if !final && !t.lastDraw.IsZero() && now.Sub(t.lastDraw) < 80*time.Millisecond {
		return
	}

	elapsed := now.Sub(t.started)
	if elapsed <= 0 {
		elapsed = time.Millisecond
	}

	speed := float64(t.read) / elapsed.Seconds()
	percent := int64(0)
	if t.total > 0 {
		percent = int64(math.Round((float64(t.read) / float64(t.total)) * 100))
		if percent > 100 {
			percent = 100
		}
	}

	spinFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := spinFrames[t.frameIndex%len(spinFrames)]
	t.frameIndex++
	if final {
		frame = "✓"
	}

	bar := renderASCIIBar(percent, 12)
	readableRead := units.HumanSize(float64(t.read))
	readableSpeed := units.HumanSize(speed) + "/s"
	prefix := "\033[1;36mDownloading\033[0m"
	if !t.printedHdr {
		fmt.Fprintf(t.out, "%s...\n", prefix)
		t.printedHdr = true
	}

	line := ""
	if t.total > 0 {
		readableTotal := units.HumanSize(float64(t.total))
		eta := formatETA(t.total, t.read, speed)
		line = fmt.Sprintf("\r  %s %s %s Downloading... %3d%%  %8s/%-8s  %8s  ETA %s\033[K",
			prefix, frame, bar, percent, readableRead, readableTotal, readableSpeed, eta)
	} else {
		line = fmt.Sprintf("\r  %s %s %s Downloading...      %8s downloaded  %8s\033[K", prefix, frame, bar, readableRead, readableSpeed)
	}

	fmt.Fprint(t.out, line)
	t.lastDraw = now
}

func renderASCIIBar(percent int64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int((percent * int64(width)) / 100)
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}

func formatETA(total, read int64, speed float64) string {
	if total <= 0 || speed <= 0 || read >= total {
		return "00:00"
	}

	remaining := float64(total - read)
	seconds := int64(math.Ceil(remaining / speed))
	if seconds < 0 {
		seconds = 0
	}

	min := seconds / 60
	sec := seconds % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

func newLineProgressReader(r io.ReadCloser, total int64, out io.Writer) io.ReadCloser {
	return &lineProgressReader{
		r:            r,
		total:        total,
		nextPercent:  10,
		lastReported: time.Now(),
		started:      time.Now(),
		out:          out,
	}
}

func (l *lineProgressReader) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)

	if n > 0 {
		l.read += int64(n)
		l.report()
	}

	return n, err
}

func (l *lineProgressReader) Close() error {
	if l.total > 0 && l.read > 0 {
		if l.lastPercent < 100 {
			fmt.Fprintf(l.out, "Download progress: %d%% (%s/%s)\n",
				100,
				units.HumanSize(float64(l.read)),
				units.HumanSize(float64(l.total)))
		}
	} else if l.read > 0 {
		fmt.Fprintf(l.out, "Download progress: %s downloaded\n", units.HumanSize(float64(l.read)))
	}

	fmt.Fprintf(l.out, "Download complete in %s\n", time.Since(l.started).Round(time.Millisecond))
	return l.r.Close()
}

func (l *lineProgressReader) report() {
	now := time.Now()

	if l.total > 0 {
		percent := (l.read * 100) / l.total
		if percent > 100 {
			percent = 100
		}
		if percent >= l.nextPercent {
			fmt.Fprintf(l.out, "Download progress: %d%% (%s/%s)\n",
				percent,
				units.HumanSize(float64(l.read)),
				units.HumanSize(float64(l.total)))
			l.lastPercent = percent
			l.nextPercent += 10
			l.lastReported = now
			return
		}
	}

	if now.Sub(l.lastReported) >= 2*time.Second {
		if l.total > 0 {
			percent := (l.read * 100) / l.total
			if percent > 100 {
				percent = 100
			}
			fmt.Fprintf(l.out, "Download progress: %d%% (%s/%s)\n",
				percent,
				units.HumanSize(float64(l.read)),
				units.HumanSize(float64(l.total)))
			l.lastPercent = percent
		} else {
			fmt.Fprintf(l.out, "Download progress: %s downloaded\n", units.HumanSize(float64(l.read)))
		}
		l.lastReported = now
	}
}

func decompress(tool *Tool, downloadURL, outFilePath, operatingSystem, arch, version string, quiet bool) (string, error) {

	archiveFile, err := os.Open(outFilePath)
	if err != nil {
		return "", err
	}

	outFilePathDir := filepath.Dir(outFilePath)
	if len(tool.BinaryTemplate) == 0 && len(tool.URLTemplate) > 0 {
		outFilePath = path.Join(outFilePathDir, tool.Name)
	} else if len(tool.BinaryTemplate) > 0 && len(tool.URLTemplate) == 0 &&
		(!strings.Contains(tool.BinaryTemplate, "tar.gz") &&
			!strings.Contains(tool.BinaryTemplate, "zip") &&
			!strings.Contains(tool.BinaryTemplate, "tgz")) {
		fileName, err := GetBinaryName(tool,
			strings.ToLower(operatingSystem),
			strings.ToLower(arch),
			version)
		if err != nil {
			return "", err
		}

		outFilePath = path.Join(outFilePathDir, fileName)
	} else if len(tool.BinaryTemplate) > 0 && len(tool.URLTemplate) > 0 {
		fileName, err := GetBinaryName(tool,
			strings.ToLower(operatingSystem),
			strings.ToLower(arch),
			version)
		if err != nil {
			return "", err
		}
		outFilePath = path.Join(outFilePathDir, fileName)

	} else {
		outFilePath = path.Join(outFilePathDir, tool.Name)
	}

	if strings.Contains(strings.ToLower(operatingSystem), "mingw") && tool.NoExtension == false {
		outFilePath += ".exe"
	}

	forceQuiet := true

	if strings.HasSuffix(downloadURL, "tar.gz") || strings.HasSuffix(downloadURL, "tgz") {
		if err := archive.Untar(archiveFile, outFilePathDir, true, forceQuiet); err != nil {
			return "", err
		}
	} else if strings.HasSuffix(downloadURL, "zip") {
		fInfo, err := archiveFile.Stat()
		if err != nil {
			return "", err
		}

		if !quiet {
			log.Printf("Name: %s, size: %d", fInfo.Name(), fInfo.Size())
		}

		if err := archive.Unzip(archiveFile, fInfo.Size(), outFilePathDir, forceQuiet); err != nil {
			return "", err
		}
	}

	return outFilePath, nil
}

func fetchText(url string) (string, error) {
	return retryWithBackoff(func() (string, error) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", pkg.UserAgent())
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		var body []byte
		if res.Body != nil {
			defer res.Body.Close()
			body, _ = io.ReadAll(res.Body)
		}

		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unexpected status code %d, body: %s", res.StatusCode, string(body))
		}

		return string(body), nil
	}, 10, 100*time.Millisecond)
}

func verifySHA(shaSum, outFilePath string) error {

	outFileBaseName := filepath.Base(outFilePath)

	lines := strings.Split(shaSum, "\n")

	for _, line := range lines {
		remoteHash, file, ok := strings.Cut(line, " ")
		if ok {

			if file == outFileBaseName {
				calculated, err := getSHA256Checksum(outFilePath)
				if err != nil {
					return err
				}
				if calculated != remoteHash {
					return fmt.Errorf("checksum mismatch, want: %s, but got: %s", remoteHash, calculated)
				}
			}
		}
	}
	return nil
}

func getSHA256Checksum(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(f)), nil
}
