package get

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
)

// installPackage preserves resources beside bin/ and exposes only the declared
// executables. Relative links keep custom --path installations relocatable.
func installPackage(tool *Tool, archivePath, downloadURL, operatingSystem, movePath string) (installedPath string, err error) {
	if !strings.HasSuffix(downloadURL, ".tar.gz") {
		return "", fmt.Errorf("%s: package installation requires a tar.gz archive", tool.Name)
	}
	if movePath == "" {
		if _, err := config.InitUserDir(); err != nil {
			return "", err
		}
		movePath = filepath.Dir(env.LocalBinary(tool.Name, ""))
	}

	stage, err := os.MkdirTemp(movePath, ".arkade-package-*")
	if err != nil {
		return "", err
	}
	keepStage := false
	defer func() {
		if !keepStage {
			os.RemoveAll(stage)
		}
	}()
	packageDir := filepath.Join(stage, "package")
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	err = archive.UntarNested(f, packageDir, true, true, true, false)
	f.Close()
	if err != nil {
		return "", err
	}

	packageName := ".arkade-" + tool.Name
	packagePath := filepath.Join(movePath, packageName)
	names := make([]string, 0, len(tool.PackageBinaries))
	seen := map[string]bool{}
	for _, name := range tool.PackageBinaries {
		if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) || seen[name] {
			return "", fmt.Errorf("%s: invalid or duplicate package binary %q", tool.Name, name)
		}
		seen[name] = true
		if strings.HasPrefix(strings.ToLower(operatingSystem), "ming") && !tool.NoExtension {
			name += ".exe"
		}
		source := filepath.Join(packageDir, "bin", name)
		info, err := os.Lstat(source)
		if err != nil {
			return "", fmt.Errorf("package binary %s: %w", name, err)
		}
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("package binary %s is not a regular file", name)
		}
		if err := os.Chmod(source, 0755); err != nil {
			return "", err
		}
		if info, err := os.Lstat(filepath.Join(movePath, name)); err == nil && info.IsDir() {
			return "", fmt.Errorf("package binary destination %s is a directory", name)
		} else if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		// Create links before changing the installation, including checking the
		// host's ability to create symlinks (Developer Mode/elevation on Windows).
		if err := os.Symlink(filepath.Join(packageName, "bin", name), filepath.Join(stage, name)); err != nil {
			return "", fmt.Errorf("create package link %s: %w", name, err)
		}
		names = append(names, name)
	}
	if !seen[tool.Name] {
		return "", fmt.Errorf("%s: PackageBinaries must include the main binary", tool.Name)
	}

	// Keep the previous package and entry points until every link is in place.
	// If publication fails, restore them before removing the staging directory.
	previousPackage := filepath.Join(stage, "previous-package")
	previous := map[string]string{}
	published := []string{}
	hadPackage, movedPackage := false, false
	defer func() {
		if err == nil {
			return
		}
		for _, name := range published {
			os.Remove(filepath.Join(movePath, name))
		}
		for name, backup := range previous {
			if restoreErr := os.Rename(backup, filepath.Join(movePath, name)); restoreErr != nil {
				keepStage = true
				log.Printf("Restore %s: %s; backup retained at %s", name, restoreErr, backup)
			}
		}
		if movedPackage {
			os.RemoveAll(packagePath)
		}
		if hadPackage {
			if restoreErr := os.Rename(previousPackage, packagePath); restoreErr != nil {
				keepStage = true
				log.Printf("Restore package %s: %s; backup retained at %s", packagePath, restoreErr, previousPackage)
			}
		}
	}()
	if err = os.Rename(packagePath, previousPackage); err == nil {
		hadPackage = true
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err = os.Rename(packageDir, packagePath); err != nil {
		return "", err
	}
	movedPackage = true
	for _, name := range names {
		destination := filepath.Join(movePath, name)
		backup := filepath.Join(stage, "previous-"+name)
		if err = os.Rename(destination, backup); err == nil {
			previous[name] = backup
		} else if !os.IsNotExist(err) {
			return "", err
		}
		if err = os.Rename(filepath.Join(stage, name), destination); err != nil {
			return "", err
		}
		published = append(published, name)
	}
	mainName := tool.Name
	if strings.HasPrefix(strings.ToLower(operatingSystem), "ming") && !tool.NoExtension {
		mainName += ".exe"
	}
	return filepath.Join(movePath, mainName), nil
}
