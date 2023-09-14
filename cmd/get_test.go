package cmd

import (
	"strings"
	"testing"
)

func Test_GetCommandWithInvalidArch(t *testing.T) {
	tests := []struct {
		args      []string
		wantError string
	}{
		{
			args:      []string{"kind", "--arch", "invalid"},
			wantError: "cpu architecture \"invalid\" is not supported",
		},
		{
			args:      []string{"kind", "--os", "invalid"},
			wantError: "operating system \"invalid\" is not supported",
		},
	}

	for _, tc := range tests {
		cmd := MakeGet()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		if !strings.Contains(err.Error(), tc.wantError) {
			t.Fatalf("for args %q\n want: %q\n but got: %q", tc.args, tc.wantError, err)
		}
	}
}
