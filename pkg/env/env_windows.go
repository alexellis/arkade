// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

//go:build windows
// +build windows

package env

import (
	"os"
	"path"
	"runtime"
)

// GetClientArch returns a pair of arch and os
func GetClientArch() (arch string, os string) {
	arch = runtime.GOARCH
	return arch, "ming"
}

func LocalBinary(name, subdir string) string {
	home := os.Getenv("HOME")
	val := path.Join(home, ".arkade/bin/")
	if len(subdir) > 0 {
		val = path.Join(val, subdir)
	}

	return path.Join(val, name)
}
