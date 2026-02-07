package system

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

const (
	githubDownloadTemplate = "https://github.com/%s/%s/releases/download/%s/%s"
	githubLatest           = "latest"

	readWriteExecuteEveryone = 0755
)

// spinWhile runs fn in a goroutine and displays a braille spinner on
// stderr while waiting for it to finish. In non-TTY environments it
// prints a single line with the message and "done" when complete.
func spinWhile(msg string, fn func() error) error {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan error, 1)
	go func() { done <- fn() }()

	tty := isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
	if !tty {
		fmt.Fprintf(os.Stderr, "%s...", msg)
		err := <-done
		if err == nil {
			fmt.Fprintln(os.Stderr, " done.")
		} else {
			fmt.Fprintln(os.Stderr, " failed.")
		}
		return err
	}

	tick := time.NewTicker(80 * time.Millisecond)
	defer tick.Stop()
	i := 0
	for {
		select {
		case err := <-done:
			// Clear the spinner line.
			fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", len(msg)+4))
			return err
		case <-tick.C:
			fmt.Fprintf(os.Stderr, "\r%s %s", frames[i%len(frames)], msg)
			i++
		}
	}
}
