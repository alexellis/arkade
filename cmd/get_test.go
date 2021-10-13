package cmd

import (
	"strings"
	"testing"
)

func Test_GetCommand(t *testing.T) {
	tests := []struct {
		args          []string
		expectedError string
	}{
		{
			args:          []string{"kind", "--arch", "dummy"},
			expectedError: "cpu architecture \"dummy\" is not supported",
		},
		{
			args:          []string{"kind", "--os", "dummy"},
			expectedError: "operating system \"dummy\" is not supported",
		},
	}

	for _, tc := range tests {
		cmd := MakeGet()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()
		if !strings.Contains(err.Error(), tc.expectedError) {
			t.Fatalf("for args %q\n want: %q\n but got: %q", tc.args, tc.expectedError, err)
		}
	}
}
