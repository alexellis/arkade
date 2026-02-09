// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	units "github.com/docker/go-units"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
)

// ── tool state tracking ────────────────────────────────────────────

const (
	stQueued      = "queued"
	stResolving   = "resolving"
	stDownloading = "downloading"
	stDone        = "done"
	stFailed      = "failed"
)

// toolProgress holds the live state for one tool in the download group.
type toolProgress struct {
	name       string
	version    string // resolved version, set after resolving completes
	status     string
	bytesRead  int64 // updated atomically from callback
	totalBytes int64 // updated atomically from callback
	started    time.Time
	elapsed    time.Duration
	err        error
	path       string
}

// ── MakeGet ────────────────────────────────────────────────────────

// MakeGet creates the Get command to download software
func MakeGet() *cobra.Command {
	tools := get.MakeTools()
	sort.Sort(tools)
	var validToolOptions []string = make([]string, len(tools))

	for _, t := range tools {
		validToolOptions = append(validToolOptions, t.Name)
	}

	var command = &cobra.Command{
		Use:   "get",
		Short: `The get command downloads a tool`,
		Long: `The get command downloads a CLI or application from the specific tool's
releases or downloads page. The tool is usually downloaded in binary format
and provides a fast and easy alternative to a package manager.`,
		Example: `  arkade get helm

  # Download multiple tools in parallel
  arkade get kubectl helm k9s stern

  # Increase parallel downloads
  arkade get kubectl helm k9s stern --parallel 8

  # Override the version
  arkade get kubectl@v1.19.3
  arkade get terraform --version=1.7.4

  # Override the OS
  arkade get helm --os darwin --arch aarch64
  arkade get helm --os linux --arch armv7l

  # Get a complete list of CLIs to download:
  arkade get`,
		SilenceUsage: true,
		Aliases:      []string{"g", "d", "download"},
		ValidArgs:    validToolOptions,
	}

	clientArch, clientOS := env.GetClientArch()

	command.Flags().Bool("progress", true, "Display a progress bar")
	command.Flags().StringP("format", "o", "", "Format format of the list of tools (table/markdown/list)")
	command.Flags().String("path", "", "Leave empty to store in HOME/.arkade/bin/, otherwise give a path for the resulting binaries")
	command.Flags().StringP("version", "v", "", "Download a specific version")
	command.Flags().String("arch", clientArch, "CPU architecture for the tool")
	command.Flags().String("os", clientOS, "Operating system for the tool")
	command.Flags().Bool("quiet", false, "Suppress most additional format")
	command.Flags().Bool("verify", true, "Verify the checksum of the downloaded file where a download has a verify strategy defined")
	command.Flags().IntP("parallel", "p", 4, "Maximum number of parallel downloads")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		verify, _ := command.Flags().GetBool("verify")

		if len(args) == 0 {
			format, _ := command.Flags().GetString("format")

			if len(format) > 0 {
				if get.TableFormat(format) == get.MarkdownStyle {
					get.CreateToolsTable(tools, get.MarkdownStyle)
				} else if get.TableFormat(format) == get.ListStyle {
					for _, r := range tools {
						fmt.Printf("%s\n", r.Name)
					}

				} else {
					get.CreateToolsTable(tools, get.TableStyle)
				}
			} else {
				get.CreateToolsTable(tools, get.TableStyle)
			}
			return nil
		}

		version := ""
		if command.Flags().Changed("version") {
			version, _ = command.Flags().GetString("version")
		}

		downloadURLs, err := get.GetDownloadURLs(tools, args, version)
		if err != nil {
			return err
		}

		movePath, _ := command.Flags().GetString("path")
		quiet, _ := command.Flags().GetBool("quiet")
		showProgress, _ := command.Flags().GetBool("progress")
		parallel, _ := command.Flags().GetInt("parallel")
		if parallel < 1 {
			return fmt.Errorf("--parallel must be at least 1")
		}

		movePath = os.ExpandEnv(movePath)

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		arch, _ := command.Flags().GetString("arch")
		if err := get.ValidateArch(arch); err != nil {
			return err
		}

		operatingSystem, _ := command.Flags().GetString("os")
		if err := get.ValidateOS(operatingSystem); err != nil {
			return err
		}

		if parallel > len(downloadURLs) {
			parallel = len(downloadURLs)
		}

		// ── per-tool progress state ──────────────────────────
		progress := make([]toolProgress, len(downloadURLs))
		for i, tool := range downloadURLs {
			progress[i] = toolProgress{name: tool.Name, status: stQueued}
		}

		// ── completion events ────────────────────────────────
		type downloadEvent struct {
			toolIndex int
			resolving bool   // tool entered resolving phase
			started   bool   // tool entered downloading phase
			version   string // resolved version (set with started=true)
			path      string
			err       error
		}

		events := make(chan downloadEvent, len(downloadURLs)*3)
		indexes := make(chan int, len(downloadURLs))

		worker := func() {
			for idx := range indexes {
				tool := downloadURLs[idx]

				// Phase 1: resolve version.
				events <- downloadEvent{toolIndex: idx, resolving: true}

				resolved, err := get.ResolveVersion(&tool, version)
				if err != nil {
					events <- downloadEvent{toolIndex: idx, err: err}
					continue
				}

				// Phase 2: download.
				events <- downloadEvent{toolIndex: idx, started: true, version: resolved}

				cb := func(bytesRead, totalBytes int64) {
					atomic.StoreInt64(&progress[idx].bytesRead, bytesRead)
					atomic.StoreInt64(&progress[idx].totalBytes, totalBytes)
				}

				toolPath, _, dlErr := get.DownloadWithProgress(&tool,
					arch,
					operatingSystem,
					resolved,
					movePath,
					true, // quiet — the renderer owns the display
					verify,
					cb)
				events <- downloadEvent{toolIndex: idx, path: toolPath, err: dlErr}
			}
		}

		var wg sync.WaitGroup
		for i := 0; i < parallel; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				worker()
			}()
		}

		for i := range downloadURLs {
			indexes <- i
		}
		close(indexes)

		tty := !quiet && isTTY()

		// --progress=false suppresses the live download renderer.
		// Summary and post-install messages are still shown unless --quiet.
		renderProgress := showProgress && !quiet

		// In TTY mode output goes to stderr so that stdout stays clean
		// for piping. In plain/CI mode output goes to stdout so that
		// `> out.txt` captures everything.
		out := os.Stdout
		if tty {
			out = os.Stderr
		}

		// Hide cursor during live rendering.
		if tty && renderProgress {
			fmt.Fprint(out, "\033[?25l")
			defer fmt.Fprint(out, "\033[?25h\n")
		}

		// Restore cursor on signal.
		go func() {
			<-signalChan
			if tty && renderProgress {
				fmt.Fprint(out, "\033[?25h\n")
			}
			os.Exit(2)
		}()

		finished := 0
		var firstErr error
		var firstErrTool *get.Tool
		groupStart := time.Now()

		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for finished < len(downloadURLs) {
			select {
			case ev := <-events:
				if ev.resolving {
					progress[ev.toolIndex].status = stResolving
					progress[ev.toolIndex].started = time.Now()
					continue
				}
				if ev.started {
					progress[ev.toolIndex].status = stDownloading
					progress[ev.toolIndex].version = ev.version
					continue
				}

				finished++
				if ev.err != nil {
					progress[ev.toolIndex].status = stFailed
					progress[ev.toolIndex].err = ev.err
					progress[ev.toolIndex].elapsed = time.Since(progress[ev.toolIndex].started)
					if firstErr == nil {
						firstErr = ev.err
						t := downloadURLs[ev.toolIndex]
						firstErrTool = &t
					}
				} else {
					progress[ev.toolIndex].status = stDone
					progress[ev.toolIndex].path = ev.path
					progress[ev.toolIndex].elapsed = time.Since(progress[ev.toolIndex].started)
				}

			case <-ticker.C:
			}

			// Render after every event or tick.
			if renderProgress {
				if tty {
					renderTTY(out, progress, parallel)
				} else {
					renderPlain(out, progress)
				}
			}
		}

		wg.Wait()
		close(events)

		// Final render.
		if renderProgress {
			if tty {
				renderTTY(out, progress, parallel)
			} else {
				renderPlain(out, progress)
			}
		}

		if firstErr != nil {
			if firstErrTool != nil && errors.Is(firstErr, &get.ErrNotFound{}) {
				printGetNotFoundError(*firstErrTool, operatingSystem, arch)
			}
		}

		// Collect successful downloads.
		var localToolsStore []get.ToolLocal
		for _, p := range progress {
			if p.status == stDone && len(p.path) > 0 {
				localToolsStore = append(localToolsStore, get.ToolLocal{
					Name: p.name,
					Path: p.path,
				})
			}
		}

		if !quiet {
			fmt.Fprintln(out)

			// ── Group summary ────────────────────────────
			printGroupSummary(out, progress, time.Since(groupStart), tty)

			if len(localToolsStore) > 0 {
				arkadeBinInPath := movePath == "" && get.ArkadeInPath()

				// When .arkade/bin is already in PATH, the tools are
				// directly usable — skip the installation section to
				// keep output minimal. Conflict warnings still show.
				if !arkadeBinInPath {
					// ── Installation instructions ────────────────
					if tty {
						fmt.Fprintf(out, "\033[1m── Installation ──\033[0m\n\n")
					} else {
						fmt.Fprintf(out, "-- Installation --\n\n")
					}

					// Post-installation message.
					msg, err := get.PostInstallationMsg(movePath, localToolsStore)
					if err != nil {
						return err
					}
					fmt.Fprintf(out, "%s\n", msg)
				}

				// Warn about conflicting binaries found elsewhere in $PATH.
				if movePath == "" {
					homeDir, _ := os.UserHomeDir()
					arkadeBin := filepath.Clean(filepath.Join(homeDir, ".arkade", "bin"))

					pathDirs := filepath.SplitList(os.Getenv("PATH"))

					for _, tl := range localToolsStore {
						var hits []string
						seen := map[string]bool{}
						for _, dir := range pathDirs {
							clean := filepath.Clean(dir)
							// Expand literal ~ so comparison with arkadeBin works.
							if homeDir != "" && strings.HasPrefix(clean, "~") {
								clean = filepath.Clean(homeDir + clean[1:])
							}
							// Skip the arkade bin dir itself — that's where we just installed.
							if clean == arkadeBin {
								continue
							}
							// Deduplicate (PATH often has repeated entries).
							if seen[clean] {
								continue
							}
							candidate := filepath.Join(clean, tl.Name)
							if _, err := os.Stat(candidate); err == nil {
								seen[clean] = true
								// Shorten $HOME prefix to ~ for readability.
								display := clean
								if homeDir != "" && strings.HasPrefix(clean, homeDir) {
									display = "~" + clean[len(homeDir):]
								}
								hits = append(hits, display)
							}
						}
						if len(hits) > 0 {
							fmt.Fprintf(out, " ⚠  %s conflicts: %s\n", tl.Name, strings.Join(hits, ", "))
						}
					}
				}
			}
		}

		return firstErr
	}
	return command
}

// ── TTY renderer ───────────────────────────────────────────────────
//
// Produces SCP / sftp-style output. Each tool gets one line that is
// updated in-place with a progress bar, percentage, speed and ETA.
// Done tools show a tick (✔), failed tools show a cross (✘).
//
// Example:
//
//   ✔  kubectl    100%  [████████████████████]  48.3 MB   12.1 MB/s  00:04
//   ⠸  helm        63%  [████████████░░░░░░░░]  18.2 MB   9.4 MB/s  ETA 02s
//   ⠼  k9s         21%  [████░░░░░░░░░░░░░░░░]   6.1 MB   5.8 MB/s  ETA 05s
//      stern             queued
//

var spinFrames = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func renderTTY(out io.Writer, progress []toolProgress, parallel int) {
	// Move cursor to top-left and clear screen.
	var b strings.Builder

	b.WriteString("\033[H\033[2J")

	// Title.
	b.WriteString("\033[1m Arkade by Alex Ellis - https://github.com/sponsors/alexellis\033[0m\n\n")

	// Status bar.
	total, done, active, queued, failed := summariseCounts(progress)
	b.WriteString(fmt.Sprintf(
		" Downloading %d tool(s)  \033[32m✔ %d\033[0m  ↓ %d  · %d  \033[31m✘ %d\033[0m  parallel:%d\n\n",
		total, done, active, queued, failed, parallel))

	// Determine max name width (including version suffix when known).
	nameW := 10
	for i := range progress {
		l := len(progress[i].name)
		if progress[i].version != "" {
			l += len(progress[i].version) + 3 // " (v1.2.3)"
		}
		if l > nameW {
			nameW = l
		}
	}
	if nameW > 40 {
		nameW = 40
	}

	tick := int(time.Now().UnixMilli() / 80)
	barW := 20

	for i := range progress {
		p := &progress[i]
		read := atomic.LoadInt64(&p.bytesRead)
		total := atomic.LoadInt64(&p.totalBytes)

		switch p.status {
		case stDone:
			pct := 100
			if total > 0 {
				pct = int(math.Round(float64(read) / float64(total) * 100))
				if pct > 100 {
					pct = 100
				}
			}
			bar := renderBar(int64(pct), barW)
			sizeStr := units.HumanSize(float64(read))
			speed := float64(0)
			if p.elapsed > 0 {
				speed = float64(read) / p.elapsed.Seconds()
			}
			elapsed := fmtDuration(p.elapsed)
			displayName := p.name
			if p.version != "" {
				displayName = fmt.Sprintf("%s (%s)", p.name, p.version)
			}
			b.WriteString(fmt.Sprintf(
				" \033[32m✔\033[0m  %-*s  %3d%%  %s  %8s  %9s/s  %s\n",
				nameW, displayName, pct, bar, sizeStr, units.HumanSize(speed), elapsed))

		case stFailed:
			errMsg := ""
			if p.err != nil {
				errMsg = p.err.Error()
				if len(errMsg) > 40 {
					errMsg = errMsg[:40] + "…"
				}
			}
			b.WriteString(fmt.Sprintf(
				" \033[31m✘\033[0m  %-*s  \033[31mfailed: %s\033[0m\n",
				nameW, p.name, errMsg))

		case stDownloading:
			frame := spinFrames[tick%len(spinFrames)]
			displayName := p.name
			if p.version != "" {
				displayName = fmt.Sprintf("%s (%s)", p.name, p.version)
			}

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
					" \033[36m%s\033[0m  %-*s  %3d%%  %s  %8s  %9s/s  ETA %s\n",
					frame, nameW, displayName, pct, bar, units.HumanSize(float64(read)), units.HumanSize(speed), eta))
			} else if read > 0 {
				// Unknown total — show transferred bytes.
				elapsed := time.Since(p.started)
				speed := float64(0)
				if elapsed > 0 {
					speed = float64(read) / elapsed.Seconds()
				}
				bar := renderBar(-1, barW) // indeterminate
				b.WriteString(fmt.Sprintf(
					" \033[36m%s\033[0m  %-*s   --%%  %s  %8s  %9s/s\n",
					frame, nameW, displayName, bar, units.HumanSize(float64(read)), units.HumanSize(speed)))
			} else {
				b.WriteString(fmt.Sprintf(
					" \033[36m%s\033[0m  %-*s  \033[2mconnecting…\033[0m\n",
					frame, nameW, displayName))
			}

		case stResolving:
			frame := spinFrames[tick%len(spinFrames)]
			b.WriteString(fmt.Sprintf(
				" \033[33m%s\033[0m  %-*s  \033[2mresolving…\033[0m\n",
				frame, nameW, p.name))

		case stQueued:
			b.WriteString(fmt.Sprintf(
				" \033[2m·\033[0m  %-*s  \033[2mqueued\033[0m\n",
				nameW, p.name))
		}
	}

	fmt.Fprint(out, b.String())
}

// ── Non-TTY / plain renderer ───────────────────────────────────────
//
// Prints each state change on its own line, no ANSI, no overwrites.
// Suitable for CI, piped output, TERM=dumb, etc.

// plainLastStatus tracks what was last printed per tool so we avoid
// spamming duplicate lines.
var plainLastStatus = map[int]string{}

func renderPlain(out io.Writer, progress []toolProgress) {
	for i := range progress {
		p := &progress[i]
		key := fmt.Sprintf("%s:%s", p.name, p.status)
		if plainLastStatus[i] == key {
			continue
		}
		plainLastStatus[i] = key

		switch p.status {
		case stResolving:
			fmt.Fprintf(out, "[resolving]   %s\n", p.name)
		case stDownloading:
			fmt.Fprintf(out, "[downloading] %s\n", p.name)
		case stDone:
			read := atomic.LoadInt64(&p.bytesRead)
			sizeStr := units.HumanSize(float64(read))
			elapsed := fmtDuration(p.elapsed)
			displayName := p.name
			if p.version != "" {
				displayName = fmt.Sprintf("%s (%s)", p.name, p.version)
			}
			fmt.Fprintf(out, "[done]        %s  %s  %s\n", displayName, sizeStr, elapsed)
		case stFailed:
			errMsg := ""
			if p.err != nil {
				errMsg = p.err.Error()
			}
			fmt.Fprintf(out, "[failed]      %s  %s\n", p.name, errMsg)
		}
	}
}

// ── Group summary ──────────────────────────────────────────────────

func printGroupSummary(out io.Writer, progress []toolProgress, wall time.Duration, tty bool) {
	_, done, _, _, failed := summariseCounts(progress)

	if tty {
		fmt.Fprintf(out, "\033[1m── Summary ──\033[0m\n\n")
	} else {
		fmt.Fprintf(out, "-- Summary --\n\n")
	}

	for i := range progress {
		p := &progress[i]
		read := atomic.LoadInt64(&p.bytesRead)

		if p.status == stDone {
			stat, err := os.Stat(p.path)
			size := units.HumanSize(float64(read))
			if err == nil {
				size = units.HumanSize(float64(stat.Size()))
			}
			speed := float64(0)
			if p.elapsed > 0 {
				speed = float64(read) / p.elapsed.Seconds()
			}
			displayName := p.name
			if p.version != "" {
				displayName = fmt.Sprintf("%s (%s)", p.name, p.version)
			}
			if tty {
				fmt.Fprintf(out, " \033[32m✔\033[0m  %-20s %8s  %9s/s  %s  %s\n",
					displayName, size, units.HumanSize(speed), fmtDuration(p.elapsed), p.path)
			} else {
				fmt.Fprintf(out, " OK  %-20s %8s  %9s/s  %s  %s\n",
					displayName, size, units.HumanSize(speed), fmtDuration(p.elapsed), p.path)
			}
		} else if p.status == stFailed {
			errMsg := ""
			if p.err != nil {
				errMsg = p.err.Error()
			}
			if tty {
				fmt.Fprintf(out, " \033[31m✘\033[0m  %-20s %s\n", p.name, errMsg)
			} else {
				fmt.Fprintf(out, " ERR %-20s %s\n", p.name, errMsg)
			}
		}
	}

	fmt.Fprintf(out, "\n %d succeeded, %d failed in %s\n\n",
		done, failed, fmtDuration(wall))
}

// ── Helpers ────────────────────────────────────────────────────────

func renderBar(pct int64, width int) string {
	if pct < 0 {
		// Indeterminate — pulse animation.
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

func summariseCounts(progress []toolProgress) (total, done, active, queued, failed int) {
	total = len(progress)
	for i := range progress {
		switch progress[i].status {
		case stDone:
			done++
		case stDownloading, stResolving:
			active++
		case stQueued:
			queued++
		case stFailed:
			failed++
		}
	}
	return
}

func isTTY() bool {
	// CI environments are never interactive.
	if ci := strings.ToLower(os.Getenv("CI")); ci == "1" || ci == "true" {
		return false
	}
	if !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		return false
	}
	term := strings.TrimSpace(strings.ToLower(os.Getenv("TERM")))
	return term != "" && term != "dumb"
}

func printGetNotFoundError(tool get.Tool, operatingSystem, arch string) {
	extra := ""
	if strings.Contains(tool.URLTemplate, "https://github.com/") ||
		!strings.Contains(tool.URLTemplate, "https://") ||
		len(tool.URLTemplate) == 0 {
		extra = fmt.Sprintf(`
* View the %s releases page: %s`, tool.Name, fmt.Sprintf("https://github.com/%s/%s/releases", tool.Owner, tool.Repo))
	}

	fmt.Fprintf(os.Stderr, `
The requested version of %s is not available or configured in arkade for %s/%s

* Check if a binary is available from the project for your Operating System%s
* Feel free to raise an issue at https://github.com/alexellis/arkade/issues for help

`, tool.Name, operatingSystem, arch, extra)
}
