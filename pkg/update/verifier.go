package update

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg"
)

type Verifier interface {
	Verify(digestUrl, newBinary string) error
}

type DefaultVerifier struct {
}

func (d DefaultVerifier) Verify(downloadUrl, newBinary string) error {
	digest, err := downloadDigest(downloadUrl + ".sha256")
	if err != nil {
		return err
	}

	if err := compareSHA(digest, newBinary); err != nil {
		return fmt.Errorf("checksum failed for %s, error: %w", newBinary, err)
	}

	fmt.Printf("Checksum verified..OK.\n")
	return nil
}

func downloadDigest(uri string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		body, _ = io.ReadAll(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d, body: %s", res.StatusCode, string(body))
	}

	return string(body), nil
}

// compareSHA returns a nil error if the local digest matches the remote digest
func compareSHA(remoteDigest, localFile string) error {

	// GitHub format may sometimes include the binary name and a space, i.e.
	// "9dcfd1611440aa15333980b860220bcd55ca1d6875692facc458caf7eb1cd042  bin/arkade-darwin-arm64"
	if strings.Contains(remoteDigest, " ") {
		t, _, _ := strings.Cut(remoteDigest, " ")
		remoteDigest = t
	}

	localDigest, err := getSHA256Checksum(localFile)
	if err != nil {
		return err
	}

	if remoteDigest != localDigest {
		return fmt.Errorf("checksum mismatch, want: %s, but got: %s", remoteDigest, localDigest)
	}

	return nil
}

func getSHA256Checksum(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(f)), nil
}
