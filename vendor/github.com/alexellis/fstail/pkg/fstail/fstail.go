// Copyright (c) Alex Ellis 2023. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package fstail

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

// fallbackPollInterval is how often a Streamer re-checks its file for
// new data without a wake-up from the watcher.
const fallbackPollInterval = time.Second

// RunOptions contains the configuration for running fstail.
type RunOptions struct {
	WorkDir          string
	Match            string
	PrefixStyle      PrefixStyle
	DisableLogPrefix bool
}

// PrefixStyle defines the prefix format for log output.
type PrefixStyle string

const (
	PrefixStyleFilename PrefixStyle = "filename"
	PrefixStyleK8s      PrefixStyle = "k8s"
	PrefixStyleNone     PrefixStyle = "none"
)

// Streamer handles streaming log output from a single file.
type Streamer struct {
	f *os.File

	k8sPrefix     bool
	disablePrefix bool

	wake chan struct{}
	stop chan struct{}
	once sync.Once
}

// NewStreamer creates a new Streamer for the given file.
func NewStreamer(f *os.File, k8sPrefix bool, disablePrefix bool) *Streamer {
	return &Streamer{
		f:             f,
		k8sPrefix:     k8sPrefix,
		disablePrefix: disablePrefix,
		wake:          make(chan struct{}, 1),
		stop:          make(chan struct{}),
	}
}

// Wake nudges the Streamer to check its file for new data. Calls coalesce,
// so it is safe to invoke for every watch event.
func (s *Streamer) Wake() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

// Stop signals the Streamer to exit and closes its file. It is safe to
// call more than once.
func (s *Streamer) Stop() {
	s.once.Do(func() {
		close(s.stop)
		s.f.Close()
	})
}

// Stream reads and outputs log lines from the file until Stop is called.
func (s *Streamer) Stream() {
	base := path.Base(s.f.Name())

	var prefix string

	if !s.k8sPrefix && !s.disablePrefix {
		prefix = fmt.Sprintf("%s| ", base)
	} else if s.k8sPrefix {
		podSt, _, ok := strings.Cut(base, "_")
		if ok {
			prefix = fmt.Sprintf("%s| ", podSt)
		}
	}

	reader := bufio.NewReader(s.f)
	ticker := time.NewTicker(fallbackPollInterval)
	defer ticker.Stop()

	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			fmt.Printf("%s%s", prefix, line)
			continue
		}

		if err != io.EOF {
			break
		}

		if rewound(s.f, reader) {
			continue
		}

		select {
		case <-s.wake:
		case <-ticker.C:
		case <-s.stop:
			return
		}
	}
}

// rewound seeks the file back to its start if it shrank below the read
// offset, as done by tail -F, returning true when that happened.
func rewound(f *os.File, reader *bufio.Reader) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}

	offset, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return false
	}

	if info.Size() >= offset {
		return false
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return false
	}

	reader.Reset(f)
	return true
}

// Run starts watching the directory for file changes and tails them.
func Run(opts RunOptions) error {
	if opts.PrefixStyle == PrefixStyleNone {
		opts.DisableLogPrefix = true
	}

	fmt.Printf("Watching: %s match: %s, prefix: %s\n", opts.WorkDir, opts.Match, opts.PrefixStyle)

	printers := make(map[string]*Streamer)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	if len(opts.Match) > 0 {
		files, err := os.ReadDir(opts.WorkDir)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if !strings.Contains(file.Name(), opts.Match) {
				continue
			}

			log.Printf("Attaching to: %s", file.Name())

			if f, err := os.Open(path.Join(opts.WorkDir, file.Name())); err == nil {
				s := NewStreamer(f, opts.PrefixStyle == PrefixStyleK8s, opts.DisableLogPrefix)
				go s.Stream()
				printers[path.Join(opts.WorkDir, file.Name())] = s
			} else {
				log.Println(err)
			}
		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if len(opts.Match) > 0 && !strings.Contains(path.Base(event.Name), opts.Match) {
					continue
				}

				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					if info, statErr := os.Stat(event.Name); statErr != nil || info.IsDir() {
						continue
					}

					if s, ok := printers[event.Name]; ok {
						s.Wake()
					} else if f, err := os.Open(event.Name); err == nil {
						s := NewStreamer(f, opts.PrefixStyle == PrefixStyleK8s, opts.DisableLogPrefix)
						go s.Stream()
						printers[event.Name] = s
					} else {
						log.Println(err)
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if s, ok := printers[event.Name]; ok {
						s.Stop()
						delete(printers, event.Name)
					}
				}

			case err := <-watcher.Errors:
				if err != nil {
					log.Fatalln("Error:", err)
				}
			}
		}
	}()

	log.Printf("Adding watch for: %s", opts.WorkDir)
	if err = watcher.Add(opts.WorkDir); err != nil {
		return fmt.Errorf("failed to add watch: %w", err)
	}

	<-done

	return nil
}
