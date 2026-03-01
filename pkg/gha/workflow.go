package gha

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

func ApplyReplacements(data []byte, replacements map[string]string) string {
	content := string(data)
	for old, newVer := range replacements {
		workflowPath, _, _ := strings.Cut(old, "@")
		newFull := fmt.Sprintf("%s@%s", workflowPath, newVer)
		content = strings.ReplaceAll(content, old, newFull)
	}
	return content
}

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

func parseWorkflow(data []byte) (map[string]interface{}, error) {
	var wf map[string]interface{}
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, err
	}
	return wf, nil
}
