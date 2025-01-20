package get

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/env"
)

const GitHubVersionStrategy = "github"
const GitLabVersionStrategy = "gitlab"
const k8sVersionStrategy = "k8s"

var supportedOS = [...]string{"linux", "darwin", "ming"}
var supportedArchitectures = [...]string{"x86_64", "arm", "amd64", "armv6l", "armv7l", "arm64", "aarch64"}

// Tool describes how to download a CLI tool from a binary
// release - whether a single binary, or an archive.
type Tool struct {
	// The name of the tool for download
	Name string

	// Repo is a GitHub repo, when no repo exists, use the same
	// as the name.
	Repo string

	// Owner is the name of the GitHub account, when no account
	// exists, use the vendor name lowercase.
	Owner string

	// Version pinned or left empty to pull the latest release
	// if any only if only BinaryTemplate is specified.
	Version string

	// Bespoke approach for finding version when none is set.
	VersionStrategy string

	// Description of what the tool is used for.
	Description string

	// URLTemplate specifies a Go template for the download URL
	// override the OS, architecture and extension
	// All whitespace will be trimmed/
	URLTemplate string

	// The binary template can be used when downloading GitHub
	// It assumes that the only part of the URL needing to be
	// templated is the binary name on a standard GitHub download
	// URL.
	BinaryTemplate string

	// NoExtension is required for tooling such as kubectx
	// which at time of writing is a bash script.
	NoExtension bool
}

type ReleaseLocation struct {
	Url     string
	Timeout time.Duration
	Method  string
}

var releaseLocations = map[string]ReleaseLocation{
	GitHubVersionStrategy: {
		Url:     "https://github.com/%s/%s/releases/latest",
		Timeout: time.Second * 10,
		Method:  http.MethodHead,
	},
	GitLabVersionStrategy: {
		Url:     "https://gitlab.com/%s/%s/-/releases/permalink/latest",
		Timeout: time.Second * 5,
		Method:  http.MethodHead,
	},
	k8sVersionStrategy: {
		Url:     "https://cdn.dl.k8s.io/release/stable.txt",
		Timeout: time.Second * 5,
		Method:  http.MethodGet,
	},
}

type ToolLocal struct {
	Name string
	Path string
}

var templateFuncs = map[string]interface{}{
	"HasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
}

func (tool Tool) IsArchive(quiet bool) (bool, error) {
	arch, operatingSystem := env.GetClientArch()
	version := ""

	downloadURL, err := GetDownloadURL(&tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version, quiet)
	if err != nil {
		return false, err
	}

	return strings.HasSuffix(downloadURL, "tar.gz") ||
		strings.HasSuffix(downloadURL, "zip") ||
		strings.HasSuffix(downloadURL, "tgz"), nil
}

func isArchiveStr(downloadURL string) bool {

	return strings.HasSuffix(downloadURL, "tar.gz") ||
		strings.HasSuffix(downloadURL, "zip") ||
		strings.HasSuffix(downloadURL, "tgz")
}

// GetDownloadURL fetches the download URL for a release of a tool
// for a given os, architecture and version
func GetDownloadURL(tool *Tool, os, arch, version string, quiet bool) (string, error) {
	ver := GetToolVersion(tool, version)

	dlURL, err := tool.GetURL(os, arch, ver, quiet)
	if err != nil {
		return "", err
	}

	return dlURL, nil
}

func GetToolVersion(tool *Tool, version string) string {

	if len(version) > 0 {
		return version
	}

	return tool.Version
}

func (tool Tool) Head(uri string) (int, string, http.Header, error) {
	req, err := http.NewRequest(http.MethodHead, uri, nil)
	if err != nil {
		return http.StatusBadRequest, "", nil, err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusBadRequest, "", nil, err
	}

	var body string
	if res.Body != nil {
		b, _ := io.ReadAll(res.Body)
		body = string(b)
	}

	return res.StatusCode, body, res.Header, nil
}

func (tool Tool) GetURL(os, arch, version string, quiet bool) (string, error) {

	if len(version) == 0 {

		if !quiet {
			log.Printf("Looking up version for: %s", tool.Name)
		}

		var releaseType string
		if len(tool.URLTemplate) == 0 ||
			strings.Contains(tool.URLTemplate, "https://github.com/") {

			releaseType = GitHubVersionStrategy

		}

		if len(tool.VersionStrategy) > 0 {
			releaseType = tool.VersionStrategy
		}

		if _, supported := releaseLocations[releaseType]; supported {

			v, err := FindRelease(releaseType, tool.Owner, tool.Repo)
			if err != nil {
				return "", err
			}
			version = v
		}

		if !quiet {
			log.Printf("Found: %s", version)
		}
	}

	if len(tool.URLTemplate) > 0 {
		return getByDownloadTemplate(tool, os, arch, version)
	}

	return getURLByGithubTemplate(tool, os, arch, version)
}

func getURLByGithubTemplate(tool Tool, os, arch, version string) (string, error) {

	var err error
	t := template.New(tool.Name + "binary")
	t = t.Funcs(templateFuncs)
	t, err = t.Parse(tool.BinaryTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	pref := map[string]string{
		"OS":            os,
		"Arch":          arch,
		"Name":          tool.Name,
		"Version":       version,
		"VersionNumber": strings.TrimPrefix(version, "v"),
	}

	err = t.Execute(&buf, pref)
	if err != nil {
		return "", err
	}

	downloadName := strings.TrimSpace(buf.String())

	return getBinaryURL(tool.Owner, tool.Repo, version, downloadName), nil
}

func FindGitHubRelease(owner, repo string) (string, error) {
	return FindRelease(GitHubVersionStrategy, owner, repo)
}

func FindRelease(location, owner, repo string) (string, error) {
	url := formatUrl(releaseLocations[location].Url, owner, repo)

	clientTimeout := releaseLocations[location].Timeout
	client := makeHTTPClient(&clientTimeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(releaseLocations[location].Method, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if releaseLocations[location].Method == http.MethodHead {

		if res.StatusCode != http.StatusMovedPermanently && res.StatusCode != http.StatusFound {
			return "", fmt.Errorf("server returned status: %d", res.StatusCode)
		}

		loc := res.Header.Get("Location")
		if len(loc) == 0 {
			return "", fmt.Errorf("unable to determine release of tool")
		}

		version := loc[strings.LastIndex(loc, "/")+1:]
		return version, nil
	}

	if res.Body == nil {
		return "", fmt.Errorf("unable to determine release of tool")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	version := string(bodyBytes)
	return version, nil
}

func formatUrl(url, owner, repo string) string {
	if strings.Contains(url, "%s") {
		return fmt.Sprintf(url, owner, repo)
	}
	return url
}

func getBinaryURL(owner, repo, version, downloadName string) string {
	if in := strings.Index(downloadName, "/"); in > -1 {
		return fmt.Sprintf(
			"https://github.com/%s/%s/releases/download/%s",
			owner, repo, downloadName)
	}
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		owner, repo, version, downloadName)
}

func getByDownloadTemplate(tool Tool, os, arch, version string) (string, error) {
	var err error
	t := template.New(tool.Name)
	t = t.Funcs(templateFuncs)
	t, err = t.Parse(tool.URLTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	inputs := map[string]string{
		"OS":            os,
		"Arch":          arch,
		"Version":       version,
		"VersionNumber": strings.TrimPrefix(version, "v"),
		"Repo":          tool.Repo,
		"Owner":         tool.Owner,
		"Name":          tool.Name,
	}

	if err := t.Execute(&buf, inputs); err != nil {
		return "", err
	}

	res := strings.TrimSpace(buf.String())
	return res, nil
}

// makeHTTPClient makes a HTTP client with good defaults for timeouts.
func makeHTTPClient(timeout *time.Duration, tlsInsecure bool) http.Client {
	return makeHTTPClientWithDisableKeepAlives(timeout, tlsInsecure, false)
}

// makeHTTPClientWithDisableKeepAlives makes a HTTP client with good defaults for timeouts.
func makeHTTPClientWithDisableKeepAlives(timeout *time.Duration, tlsInsecure bool, disableKeepAlives bool) http.Client {
	client := http.Client{}

	if timeout != nil || tlsInsecure {
		tr := &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: disableKeepAlives,
		}

		if timeout != nil {
			client.Timeout = *timeout
			tr.DialContext = (&net.Dialer{
				Timeout: *timeout,
			}).DialContext

			tr.IdleConnTimeout = 120 * time.Millisecond
			tr.ExpectContinueTimeout = 1500 * time.Millisecond
		}

		if tlsInsecure {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: tlsInsecure}
		}

		tr.DisableKeepAlives = disableKeepAlives

		client.Transport = tr
	}

	return client
}

// GetBinaryName returns the name of a binary for the given tool or an
// error if the tool's template cannot be parsed or executed.
func GetBinaryName(tool *Tool, os, arch, version string) (string, error) {
	if len(tool.BinaryTemplate) > 0 {
		var err error
		t := template.New(tool.Name + "_binaryname")
		t = t.Funcs(templateFuncs)

		t, err = t.Parse(tool.BinaryTemplate)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		ver := GetToolVersion(tool, version)
		if err := t.Execute(&buf, map[string]string{
			"OS":            os,
			"Arch":          arch,
			"Name":          tool.Name,
			"Version":       ver,
			"VersionNumber": strings.TrimPrefix(ver, "v"),
		}); err != nil {
			return "", err
		}

		res := strings.TrimSpace(buf.String())
		return res, nil
	}

	return "", errors.New("BinaryTemplate is not set")
}

// GetDownloadURLs generates a list of URL for each tool, for download
func GetDownloadURLs(tools Tools, toolArgs []string, version string) (Tools, error) {
	arkadeTools := []Tool{}

	for _, arg := range toolArgs {
		name := arg

		// Handle version specified tool name
		if i := strings.LastIndex(arg, "@"); i > -1 {
			name = arg[:i]
			if len(version) > 0 {
				return nil, fmt.Errorf("cannot specify --version flag and @ syntax at the same time for %s", name)
			}
			version = arg[i+1:]
		}

		err := toolExists(&arkadeTools, tools, name, version)
		if err != nil {
			return nil, err
		}

		// unset value for the next iteration of versioned tool
		version = ""
	}

	return arkadeTools, nil
}

// toolExists checks if user provided tool exists on arkade
func toolExists(arkadeTools *[]Tool, tools Tools, name, version string) error {
	for _, tool := range tools {
		if name == tool.Name {
			if len(version) > 0 {
				tool.Version = version
			}
			*arkadeTools = append(*arkadeTools, tool)

			return nil
		}
	}
	return fmt.Errorf("tool %s not found", name)
}

func PostToolNotFoundMsg(url string) string {
	return fmt.Sprintf("Look like this tool isn't available for your OS or/and platform. Check out the link to see if the tool you're after is available: %s.\nIf it is there, don't hesitate to open an issue and contribute to Arkade!", url)
}

// PostInstallationMsg generates installation message after tool has been downloaded
func PostInstallationMsg(movePath string, localToolsStore []ToolLocal) ([]byte, error) {

	t := template.New("Installation Instructions")

	if movePath != "" {
		t.Parse(`Run the following to copy to install the tool:

{{- range . }}
sudo install -m 755 {{.Path}} /usr/local/bin/{{.Name}}
{{- end}}`)

	} else {
		t.Parse(`# Add arkade binary directory to your PATH variable
export PATH=$PATH:$HOME/.arkade/bin/

# Test the binary:
{{- range . }}
{{.Path}}
{{- end }}

# Or install with:
{{- range . }}
sudo mv {{.Path}} /usr/local/bin/

{{- end}}`)
	}

	var tpl bytes.Buffer

	err := t.Execute(&tpl, localToolsStore)
	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), err
}

// ValidateOS returns whether a given operating system is supported
func ValidateOS(name string) error {
	for _, os := range supportedOS {
		if strings.HasPrefix(strings.ToLower(name), os) {
			return nil
		}
	}

	return fmt.Errorf("operating system %q is not supported. Available prefixes: %s",
		name, strings.Join(supportedOS[:], ", "))
}

// ValidateArch returns whether a given cpu architecture is supported
func ValidateArch(name string) error {
	for _, arch := range supportedArchitectures {
		if arch == strings.ToLower(name) {
			return nil
		}
	}
	return fmt.Errorf("cpu architecture %q is not supported. Available: %s",
		name, strings.Join(supportedArchitectures[:], ", "))
}
