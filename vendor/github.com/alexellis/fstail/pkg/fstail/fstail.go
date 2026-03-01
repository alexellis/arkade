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
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

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
}

// NewStreamer creates a new Streamer for the given file.
func NewStreamer(f *os.File, k8sPrefix bool, disablePrefix bool) *Streamer {
	return &Streamer{f: f, k8sPrefix: k8sPrefix, disablePrefix: disablePrefix}
}

// Stream reads and outputs log lines from the file.
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
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			break
		}

		fmt.Printf("%s%s", prefix, string(line))
	}
}

// Close closes the file handled by the Streamer.
func (s *Streamer) Close() {
	if s.f != nil {
		s.f.Close()
	}
}

// Run starts watching the directory for file changes and tails them.
func Run(opts RunOptions) error {
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
				if len(opts.Match) > 0 && !strings.Contains(event.Name, opts.Match) {
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					if _, ok := printers[event.Name]; !ok {
						if f, err := os.Open(event.Name); err == nil {
							s := NewStreamer(f, opts.PrefixStyle == PrefixStyleK8s, opts.DisableLogPrefix)
							go s.Stream()
							printers[event.Name] = s
						} else {
							log.Println(err)
						}
					}
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					if _, ok := printers[event.Name]; !ok {
						if f, err := os.Open(event.Name); err == nil {
							s := NewStreamer(f, opts.PrefixStyle == PrefixStyleK8s, opts.DisableLogPrefix)
							go s.Stream()
							printers[event.Name] = s
						} else {
							log.Println(err)
						}
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if _, ok := printers[event.Name]; ok {
						printers[event.Name].Close()
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
