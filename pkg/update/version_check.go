package update

import (
	"context"
	"os"
	"strings"

	"github.com/alexellis/go-execute/v2"
)

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
