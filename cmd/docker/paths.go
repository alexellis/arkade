package docker

import (
	"os"
	"path/filepath"
)

// resolveDockerfilePath accepts either:
// - a directory path, in which case it implies a "Dockerfile" within it
// - a file path, returned as-is
//
// If the input does not exist, it is returned as-is and the caller's subsequent
// file operations will surface the error.
func resolveDockerfilePath(input string) (string, error) {
	st, err := os.Stat(input)
	if err != nil {
		return input, nil
	}

	if st.IsDir() {
		return filepath.Join(input, "Dockerfile"), nil
	}

	return input, nil
}
