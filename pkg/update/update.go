package update

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/go-execute/v2"
)

type Updater struct {
	resolver     Resolver
	verify       bool
	verifier     Verifier
	force        bool
	versionCheck VersionCheck
}

type Resolver interface {
	GetRelease() (string, error)
	GetDownloadURL(release string) (string, error)
}

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

func (u Updater) WithVerifier(verifier Verifier) Updater {
	u.verifier = verifier
	return u
}

func (u Updater) WithResolver(resolver Resolver) Updater {
	u.resolver = resolver
	return u
}

func NewUpdater() Updater {
	return Updater{
		verify:       true,
		force:        false,
		versionCheck: DefaultVersionCheck{},
	}
}

func (u Updater) WithVersionCheck(check VersionCheck) Updater {
	u.versionCheck = check
	return u
}

func (u Updater) WithVerify(verify bool) Updater {
	u.verify = verify
	return u
}

func (u Updater) WithForce(force bool) Updater {
	u.force = force
	return u
}

func (u Updater) Do() error {

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	execName := filepath.Base(executable)

	targetVersion, err := u.resolver.GetRelease()
	if err != nil {
		return err
	}

	updateNeeded, err := u.versionCheck.UpdateRequired(targetVersion)
	if err != nil {
		return err
	}

	if !updateNeeded && !u.force {
		fmt.Printf("You are already using %s@%s\n", execName, targetVersion)
		return nil
	}

	downloadUrl, err := u.resolver.GetDownloadURL(targetVersion)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading: %s\n", downloadUrl)
	newBinary, err := get.DownloadFileP(downloadUrl, true)
	if err != nil {
		return err
	}

	if u.verify {
		if u.verifier == nil {
			return fmt.Errorf("verifier is nil")
		}

		if err := u.verifier.Verify(downloadUrl, newBinary); err != nil {
			return err
		}
	}

	if err := replaceExec(executable, newBinary); err != nil {
		return err
	}

	fmt.Printf("Replaced: %s..OK.\n", executable)

	return nil
}

// Copy the new binary to the same directory as the current binary before calling os.Rename to prevent an
// 'invalid cross-device link' error because the source and destination are not on the same file system.
func replaceExec(currentExec, newBinary string) error {
	targetDir := filepath.Dir(currentExec)
	filename := filepath.Base(currentExec)
	newExec := filepath.Join(targetDir, fmt.Sprintf(".%s.new", filename))

	// Copy the contents of newbinary to a new executable file
	sf, err := os.Open(newBinary)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(newExec, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}

	// Replace the current executable file with the new executable file
	if err := os.Rename(newExec, currentExec); err != nil {
		return err
	}

	return nil
}

type VersionCheck interface {
	UpdateRequired(target string) (bool, error)
}

type DefaultVersionCheck struct {
	Command  string
	Argument string
}

func (d DefaultVersionCheck) UpdateRequired(target string) (bool, error) {
	executable, err := os.Executable()
	if err != nil {
		return false, err
	}

	task := execute.ExecTask{
		Command: executable,
		Args:    []string{"version"},
	}

	res, err := task.Execute(context.TODO())
	if err != nil {
		return false, err
	}

	if !strings.Contains(res.Stdout, target) {
		return true, nil
	}

	return false, nil
}
