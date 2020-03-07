package env

import (
	"os"
	"path"
	"runtime"
)

// GetClientArch returns a pair of arch and os
func GetClientArch() (string, string) {
	return runtime.GOARCH, runtime.GOOS
}

func LocalBinary(name, subdir string) string {
	home := os.Getenv("HOME")
	val := path.Join(home, ".arkade/bin/")
	if len(subdir) > 0 {
		val = path.Join(val, subdir)
	}

	return path.Join(val, name)
}
