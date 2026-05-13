// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	units "github.com/docker/go-units"
	"github.com/mattn/go-isatty"
)

// ── status constants ───────────────────────────────────────────────

const (
	stResolving   = "resolving"
	stDownloading = "downloading"
	stExtracting  = "extracting"
	stDone        = "done"
	stFailed      = "failed"
)

// imageProgress holds the live state for a single OCI image pull.
type imageProgress struct {
	imageName  string
	platform   string // e.g. "linux/amd64"
	status     string
	bytesRead  int64 // updated atomically by countingTransport
	totalBytes int64 // sum of compressed layer sizes from the manifest
	started    time.Time
	elapsed    time.Duration
	err        error
}

// ── counting transport ────────────────────────────────────────────
//
// Wraps an http.RoundTripper and counts bytes streamed from every
// response body, giving us live download progress for crane operations
// without needing a callback in the upstream API.

type countingTransport struct {
	base http.RoundTripper
	n    *int64
}

func (c *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.base.RoundTrip(req)
	if err != nil || resp == nil {
		return resp, err
	}
	resp.Body = &countingReadCloser{rc: resp.Body, n: c.n}
	return resp, nil
}

type countingReadCloser struct {
	rc io.ReadCloser
	n  *int64
}

func (c *countingReadCloser) Read(p []byte) (int, error) {
	n, err := c.rc.Read(p)
	if n > 0 {
		atomic.AddInt64(c.n, int64(n))
	}
	return n, err
}

func (c *countingReadCloser) Close() error { return c.rc.Close() }

// ── TTY live renderer ─────────────────────────────────────────────

var spinFrames = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// renderTTY paints the alternate-screen frame.
func renderTTY(out io.Writer, p *imageProgress) {
	var b strings.Builder
	b.WriteString("\033[H\033[2J") // home + clear
	b.WriteString("\033[1m Arkade by Alex Ellis - https://github.com/sponsors/alexellis\033[0m\n\n")

	tick := int(time.Now().UnixMilli() / 80)
	frame := spinFrames[tick%len(spinFrames)]
	barW := 30

	displayName := p.imageName
	if p.platform != "" {
		displayName = fmt.Sprintf("%s (%s)", p.imageName, p.platform)
	}

	read := atomic.LoadInt64(&p.bytesRead)
	total := atomic.LoadInt64(&p.totalBytes)

	switch p.status {
	case stResolving:
		b.WriteString(fmt.Sprintf(
			" \033[33m%s\033[0m  %s  \033[2mresolving manifest…\033[0m\n",
			frame, displayName))

	case stDownloading:
		if total > 0 {
			pct := int64(math.Round(float64(read) / float64(total) * 100))
			if pct > 100 {
				pct = 100
			}
			bar := renderBar(pct, barW)
			elapsed := time.Since(p.started)
			speed := float64(0)
			if elapsed > 0 {
				speed = float64(read) / elapsed.Seconds()
			}
			eta := fmtETA(total, read, speed)
			b.WriteString(fmt.Sprintf(
				" \033[36m%s\033[0m  %s  %3d%%  %s  %s / %s  %9s/s  ETA %s\n",
				frame, displayName, pct, bar,
				units.HumanSize(float64(read)),
				units.HumanSize(float64(total)),
				units.HumanSize(speed), eta))
		} else if read > 0 {
			elapsed := time.Since(p.started)
			speed := float64(0)
			if elapsed > 0 {
				speed = float64(read) / elapsed.Seconds()
			}
			bar := renderBar(-1, barW)
			b.WriteString(fmt.Sprintf(
				" \033[36m%s\033[0m  %s   --%%  %s  %8s  %9s/s\n",
				frame, displayName, bar,
				units.HumanSize(float64(read)),
				units.HumanSize(speed)))
		} else {
			b.WriteString(fmt.Sprintf(
				" \033[36m%s\033[0m  %s  \033[2mconnecting…\033[0m\n",
				frame, displayName))
		}

	case stExtracting:
		b.WriteString(fmt.Sprintf(
			" \033[36m%s\033[0m  %s  \033[2mextracting…\033[0m\n",
			frame, displayName))
	}

	fmt.Fprint(out, b.String())
}

// renderTTYFinal prints the completed state as static text into the
// main scrollback after the alt-screen has been exited.
func renderTTYFinal(out io.Writer, p *imageProgress) {
	var b strings.Builder

	displayName := p.imageName
	if p.platform != "" {
		displayName = fmt.Sprintf("%s (%s)", p.imageName, p.platform)
	}

	read := atomic.LoadInt64(&p.bytesRead)

	switch p.status {
	case stDone:
		size := units.HumanSize(float64(read))
		speed := float64(0)
		if p.elapsed > 0 {
			speed = float64(read) / p.elapsed.Seconds()
		}
		b.WriteString(fmt.Sprintf(
			" \033[32m✔\033[0m  %s  %8s  %9s/s  %s\n",
			displayName, size, units.HumanSize(speed), fmtDuration(p.elapsed)))

	case stFailed:
		msg := ""
		if p.err != nil {
			msg = p.err.Error()
			if len(msg) > 120 {
				msg = msg[:120] + "…"
			}
		}
		b.WriteString(fmt.Sprintf(
			" \033[31m✘\033[0m  %s  \033[31mfailed: %s\033[0m\n",
			displayName, msg))
	}

	fmt.Fprint(out, b.String())
}

// ── Non-TTY / plain renderer ──────────────────────────────────────
//
// Prints each state change on its own line. During the download
// phase we also emit a line each time progress crosses a 10% bucket
// so piped output (CI, `>`, etc.) is still useful. The caller threads
// a key through to suppress duplicate lines.

func renderPlain(out io.Writer, p *imageProgress, prev string) string {
	displayName := p.imageName
	if p.platform != "" {
		displayName = fmt.Sprintf("%s (%s)", p.imageName, p.platform)
	}

	read := atomic.LoadInt64(&p.bytesRead)
	total := atomic.LoadInt64(&p.totalBytes)

	key := p.status
	if p.status == stDownloading && total > 0 {
		pct := int64(read * 100 / total)
		if pct > 100 {
			pct = 100
		}
		bucket := pct / 10 * 10
		key = fmt.Sprintf("downloading-%d", bucket)
	}
	if key == prev {
		return prev
	}

	switch p.status {
	case stResolving:
		fmt.Fprintf(out, "[resolving]   %s\n", displayName)
	case stDownloading:
		if total > 0 {
			pct := int64(read * 100 / total)
			if pct > 100 {
				pct = 100
			}
			fmt.Fprintf(out, "[downloading] %s  %3d%%  %s / %s\n",
				displayName, pct,
				units.HumanSize(float64(read)),
				units.HumanSize(float64(total)))
		} else {
			fmt.Fprintf(out, "[downloading] %s\n", displayName)
		}
	case stExtracting:
		fmt.Fprintf(out, "[extracting]  %s\n", displayName)
	case stDone:
		fmt.Fprintf(out, "[done]        %s  %s  %s\n",
			displayName, units.HumanSize(float64(read)), fmtDuration(p.elapsed))
	case stFailed:
		msg := ""
		if p.err != nil {
			msg = p.err.Error()
		}
		fmt.Fprintf(out, "[failed]      %s  %s\n", displayName, msg)
	}
	return key
}

// ── small helpers (mirror cmd/get.go) ─────────────────────────────

func renderBar(pct int64, width int) string {
	if pct < 0 {
		t := int(time.Now().UnixMilli()/120) % (width * 2)
		pos := t
		if pos >= width {
			pos = width*2 - pos - 1
		}
		buf := make([]rune, width)
		for j := range buf {
			buf[j] = '░'
		}
		if pos >= 0 && pos < width {
			buf[pos] = '█'
		}
		return "\033[2m[\033[0m" + string(buf) + "\033[2m]\033[0m"
	}
	if pct > 100 {
		pct = 100
	}
	filled := int((pct * int64(width)) / 100)
	if filled > width {
		filled = width
	}
	return "\033[2m[\033[0m" +
		"\033[32m" + strings.Repeat("█", filled) + "\033[0m" +
		"\033[2m" + strings.Repeat("░", width-filled) + "\033[0m" +
		"\033[2m]\033[0m"
}

func fmtETA(total, read int64, speed float64) string {
	if total <= 0 || speed <= 0 || read >= total {
		return "00s"
	}
	secs := int(math.Ceil(float64(total-read) / speed))
	if secs < 0 {
		secs = 0
	}
	if secs >= 60 {
		return fmt.Sprintf("%dm%02ds", secs/60, secs%60)
	}
	return fmt.Sprintf("%02ds", secs)
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	s := int(d.Seconds())
	ms := int(d.Milliseconds()) % 1000
	if s >= 60 {
		return fmt.Sprintf("%dm%02d.%01ds", s/60, s%60, ms/100)
	}
	if s > 0 {
		return fmt.Sprintf("%d.%01ds", s, ms/100)
	}
	return fmt.Sprintf("0.%01ds", ms/100)
}

func isTTY() bool {
	if ci := strings.ToLower(os.Getenv("CI")); ci == "1" || ci == "true" {
		return false
	}
	if !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		return false
	}
	term := strings.TrimSpace(strings.ToLower(os.Getenv("TERM")))
	return term != "" && term != "dumb"
}
