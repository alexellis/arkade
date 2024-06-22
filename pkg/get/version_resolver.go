package get

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alexellis/arkade/pkg"
)

// Tools can implement version resolver to resolve versions differently
type VersionResolver interface {
	GetVersion() (string, error)
	Inputs() map[string]string
}

var _ VersionResolver = (*GithubVersionResolver)(nil)

type GithubVersionResolver struct {
	Owner string
	Repo  string
}

// Inputs implements VersionResolver.
func (r *GithubVersionResolver) Inputs() map[string]string {
	return map[string]string{}
}

func (r *GithubVersionResolver) GetVersion() (string, error) {
	url := fmt.Sprintf("https://github.com/%s/%s/releases/latest", r.Owner, r.Repo)
	client := makeHTTPClient(&githubTimeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

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

var _ VersionResolver = (*K8VersionResolver)(nil)

type K8VersionResolver struct{}

// Inputs implements VersionResolver.
func (r *K8VersionResolver) Inputs() map[string]string {
	return map[string]string{}
}

func (r *K8VersionResolver) GetVersion() (string, error) {
	url := "https://cdn.dl.k8s.io/release/stable.txt"

	timeout := time.Second * 5
	client := makeHTTPClient(&timeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body == nil {
		return "", fmt.Errorf("unable to determine release of tool")
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	version := string(bodyBytes)
	return version, nil
}

var _ VersionResolver = (*GoVersionResolver)(nil)

type GoVersionResolver struct{}

// Inputs implements VersionResolver.
func (r *GoVersionResolver) Inputs() map[string]string {
	return map[string]string{}
}

func (r *GoVersionResolver) GetVersion() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://go.dev/VERSION?m=text", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body == nil {
		return "", fmt.Errorf("unexpected empty body")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	content := strings.TrimSpace(string(body))
	exp, err := regexp.Compile(`^go(\d+.\d+.\d+)`)
	if err != nil {
		return "", err
	}

	version := exp.FindStringSubmatch(content)
	if len(version) < 2 {
		return "", fmt.Errorf("failed to fetch go latest version number")
	}

	return version[1], nil
}

var _ VersionResolver = (*NodeVersionResolver)(nil)

type NodeVersionResolver struct {
	Channel string
	Version string
}

// Inputs implements VersionResolver.
func (n *NodeVersionResolver) Inputs() map[string]string {
	return map[string]string{
		"Channel": n.Channel,
	}
}

func (n *NodeVersionResolver) GetVersion() (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://nodejs.org/download/%s/%s", n.Channel, n.Version), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not find latest version for %s, (%d), body: %s", n.Version, res.StatusCode, string(body))
	}

	regex := regexp.MustCompile(`(?m)node-v(\d+.\d+.\d+)-linux-.*`)
	result := regex.FindStringSubmatch(string(body))

	if len(result) < 2 {
		if v, ok := os.LookupEnv("ARK_DEBUG"); ok && v == "1" {
			fmt.Printf("Body: %s\n", string(body))
		}
		return "", fmt.Errorf("could not find latest version for %s, (%d), %s", n.Version, res.StatusCode, result)
	}
	return result[1], nil
}
