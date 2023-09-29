package update

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alexellis/arkade/pkg/get"
)

type Resolver interface {
	GetRelease() (string, error)
	GetDownloadURL(release string) (string, error)
}

type Updater struct {
	resolver     Resolver
	verify       bool
	verifier     Verifier
	force        bool
	versionCheck VersionCheck
}

func NewUpdater() Updater {
	return Updater{
		verify:       true,
		force:        false,
		versionCheck: DefaultVersionCheck{},
	}
}

func (u Updater) WithVerifier(verifier Verifier) Updater {
	u.verifier = verifier
	return u
}

func (u Updater) WithResolver(resolver Resolver) Updater {
	u.resolver = resolver
	return u
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
