package cmd

import (
	"testing"

	"github.com/alexellis/arkade/pkg/get"
)

func Test_GetTools(t *testing.T) {
	tools := get.MakeTools()
	for _, tool := range tools {
		t.Run(tool.Name, func(t *testing.T) {
			cmd := MakeGet()
			cmd.SetArgs([]string{"--stash=false", tool.Name})
			err := cmd.Execute()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
