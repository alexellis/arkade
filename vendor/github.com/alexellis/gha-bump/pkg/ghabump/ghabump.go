// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package ghabump

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v3"
)

// RunOptions contains the configuration for running gha-bump.
type RunOptions struct {
	Target  string
	Verbose bool
	Write   bool
}

// Workflow represents a parsed GitHub Actions workflow.
type Workflow map[string]interface{}

// Run upgrades GitHub Actions workflow files to the latest major version.
func Run(opts RunOptions) error {
	files, err := FindWorkflows(opts.Target)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no workflow files found in %s", opts.Target)
	}

	if opts.Verbose {
		fmt.Printf("Found %d workflow file(s)\n\n", len(files))
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	totalUpdates := 0

	for _, file := range files {
		if opts.Verbose {
			fmt.Printf("Processing: %s\n", file)
		}

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		replacements, err := ProcessWorkflow(data, client, opts.Verbose)
		if err != nil {
			return err
		}

		if opts.Verbose && len(replacements) > 0 {
			fmt.Println("Detected following replacements:")
			for old, newVer := range replacements {
				fmt.Printf("  %s -> %s\n", old, newVer)
			}
		}

		if len(replacements) > 0 {
			updated := ApplyReplacements(data, replacements)
			totalUpdates += len(replacements)

			if opts.Write {
				if err := os.WriteFile(file, []byte(updated), 0644); err != nil {
					return err
				}
			} else {
				fmt.Print(updated)
			}
		}

		fmt.Println()
	}

	return nil
}

// FindWorkflows finds all workflow YAML files in the target path.
func FindWorkflows(target string) ([]string, error) {
	st, err := os.Stat(target)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		return []string{target}, nil
	}

	p := filepath.Join(target, ".github", "workflows")
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(p, "*.yaml"))
	if err != nil {
		return nil, err
	}

	yml, err := filepath.Glob(filepath.Join(p, "*.yml"))
	if err != nil {
		return nil, err
	}
	files = append(files, yml...)

	return files, nil
}

// ProcessWorkflow parses a workflow and suggests major version upgrades.
func ProcessWorkflow(data []byte, client *http.Client, verbose bool) (map[string]string, error) {
	workflow, err := parseWorkflow(data)
	if err != nil {
		return nil, err
	}

	jobs, ok := workflow["jobs"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("jobs not found in workflow")
	}

	replacements := map[string]string{}

	for jobName, job := range jobs {
		jobMap, ok := job.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("job %s is not a map", jobName)
		}

		steps, ok := jobMap["steps"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("steps not found or not a list in job %s", jobName)
		}

		for _, step := range steps {
			stepMap, ok := step.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("step is not a map in job %s", jobName)
			}

			var name string
			var uses string

			if v, ok := stepMap["name"].(string); ok {
				name = v
			}
			if v, ok := stepMap["uses"].(string); ok {
				uses = v
			}

			st := name
			if name == "" {
				st = uses
			}

			if verbose && len(st) > 0 {
				if len(uses) > 0 {
					fmt.Printf("  %s: %s\n", st, uses)
				} else {
					fmt.Printf("  %s\n", st)
				}
			}

			if len(uses) > 0 {
				newVer, err := suggestMajorUpgrade(client, uses)
				if err != nil {
					return nil, err
				}
				if newVer != "" {
					replacements[uses] = newVer
				}
			}
		}
	}

	return replacements, nil
}

// ApplyReplacements applies version upgrades to workflow content.
func ApplyReplacements(data []byte, replacements map[string]string) string {
	content := string(data)
	for old, newVer := range replacements {
		workflowPath, _, _ := strings.Cut(old, "@")
		newFull := fmt.Sprintf("%s@%s", workflowPath, newVer)
		content = strings.ReplaceAll(content, old, newFull)
	}
	return content
}

// suggestMajorUpgrade suggests the latest major version for an action.
func suggestMajorUpgrade(client *http.Client, uses string) (string, error) {
	ownerRepo, currentVer, ok := strings.Cut(uses, "@")
	if !ok || currentVer == "master" {
		return "", nil
	}
	owner, repo, ok := strings.Cut(ownerRepo, "/")
	if !ok {
		return "", nil
	}

	if !strings.HasPrefix(currentVer, "v") {
		return "", nil
	}

	version, err := getLatestVersion(client, owner, repo)
	if err != nil {
		return "", err
	}

	oldSemver, err := semver.NewVersion(currentVer)
	if err != nil {
		return "", err
	}
	newSemver, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	if newSemver.Major() > oldSemver.Major() {
		return fmt.Sprintf("v%d", newSemver.Major()), nil
	}

	return "", nil
}

// getLatestVersion fetches the latest release version from GitHub.
func getLatestVersion(client *http.Client, owner, repo string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://github.com/%s/%s/releases/latest", owner, repo), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		body, _ = io.ReadAll(res.Body)
	}

	if res.StatusCode != http.StatusFound {
		return "", fmt.Errorf("failed to get latest version for %s/%s: %s, body: %s", owner, repo, res.Status, string(body))
	}

	location := res.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no location header found for %s/%s", owner, repo)
	}

	parts := strings.Split(location, "/")
	if len(parts) < 7 {
		return "", fmt.Errorf("invalid location header: %s", location)
	}

	return parts[len(parts)-1], nil
}

// parseWorkflow parses YAML workflow content.
func parseWorkflow(data []byte) (map[string]interface{}, error) {
	var wf map[string]interface{}
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, err
	}
	return wf, nil
}
