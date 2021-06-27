package get

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/alexellis/arkade/pkg/env"
)

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

var templateFuncs = map[string]interface{}{
	"HasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
}

func (tool Tool) IsArchive() bool {
	arch, operatingSystem := env.GetClientArch()
	version := ""

	downloadURL, _ := GetDownloadURL(&tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
	return strings.HasSuffix(downloadURL, "tar.gz") || strings.HasSuffix(downloadURL, "zip") || strings.HasSuffix(downloadURL, "tgz")
}

func GetBinaryName(tool *Tool, os, arch, version string) (string, error) {
	if len(tool.BinaryTemplate) > 0 {
		var err error
		t := template.New(tool.Name + "_binaryname")
		t = t.Funcs(templateFuncs)
		t, err = t.Parse(tool.BinaryTemplate)
		if err != nil {
			return "", err
		}

		ver := getToolVersion(tool, version)

		var buf bytes.Buffer
		err = t.Execute(&buf, map[string]string{
			"OS":            os,
			"Arch":          arch,
			"Name":          tool.Name,
			"Version":       ver,
			"VersionNumber": strings.TrimPrefix(ver, "v"),
		})
		if err != nil {
			return "", err
		}
		res := strings.TrimSpace(buf.String())
		return res, nil
	}

	return "", errors.New("BinaryTemplate is not set")
}

// GetDownloadURL fetches the download URL for a release of a tool
// for a given os,  architecture and version
func GetDownloadURL(tool *Tool, os, arch, version string) (string, error) {
	ver := getToolVersion(tool, version)

	dlURL, err := tool.GetURL(os, arch, ver)
	if err != nil {
		return "", err
	}

	return dlURL, nil
}

func (tool Tool) GetURL(os, arch, version string) (string, error) {
	if len(tool.URLTemplate) == 0 {
		return getURLByGithubTemplate(tool, os, arch, version)
	}

	return getByDownloadTemplate(tool, os, arch, version)
}

func getURLByGithubTemplate(tool Tool, os, arch, version string) (string, error) {
	if len(version) == 0 {
		var err error
		version, err = findGitHubRelease(tool.Owner, tool.Repo)
		if err != nil {
			return "", err
		}
	}

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

func findGitHubRelease(owner, repo string) (string, error) {

	url := fmt.Sprintf("https://github.com/%s/%s/releases/latest", owner, repo)

	timeout := time.Second * 5
	client := makeHTTPClient(&timeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 302 {
		return "", fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	loc := res.Header.Get("Location")
	if len(loc) == 0 {
		return "", fmt.Errorf("unable to determine release of tool")
	}

	version := loc[strings.LastIndex(loc, "/")+1:]
	return version, nil
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

	err = t.Execute(&buf, inputs)

	if err != nil {
		return "", err
	}
	res := strings.TrimSpace(buf.String())
	return res, nil
}

// https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/kubeseal-linux-amd64
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

func getToolVersion(tool *Tool, version string) string {
	ver := tool.Version
	if len(version) > 0 {
		ver = version
	}
	return ver
}
