package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_MergeFlags(t *testing.T) {

	tests := []struct {
		title     string
		flags     map[string]string
		overrides []string
		want      map[string]string
		wantErr   error
	}{
		// positive cases:
		{"Single key with numeric value and no flags",
			map[string]string{},
			[]string{"b=1"},
			map[string]string{"b": "1"},
			nil,
		},
		{"Empty set and no flags, should not fail",
			map[string]string{},
			[]string{},
			map[string]string{},
			nil,
		},
		{"No set key and an existing flag",
			map[string]string{"a": "1"},
			[]string{},
			map[string]string{"a": "1"},
			nil,
		},
		{"Single key with numeric value, override the flag with numeric value",
			map[string]string{"a": "1"},
			[]string{"a=2"},
			map[string]string{"a": "2"},
			nil,
		},
		{"Single key with numeric value and single flag key with numeric value",
			map[string]string{"a": "1"},
			[]string{"b=1"},
			map[string]string{"a": "1", "b": "1"},
			nil,
		},
		{"Single key with numeric value, update existing key and a new key",
			map[string]string{"a": "1"},
			[]string{"a=2", "b=1"},
			map[string]string{"a": "2", "b": "1"},
			nil,
		},
		{"Update all existing flags in the map",
			map[string]string{"a": "1", "b": "2"},
			[]string{"a=2", "b=3"},
			map[string]string{"a": "2", "b": "3"},
			nil,
		},
		{"Multiple = in value",
			map[string]string{},
			[]string{"a=b=3=c=1=y=5"},
			map[string]string{"a": "b=3=c=1=y=5"},
			nil,
		},
		{"Quote the value string using '",
			map[string]string{},
			[]string{"a='b=3 c=1 y=5'"},
			map[string]string{"a": "b=3 c=1 y=5"},
			nil,
		},

		// check errors
		{"Incorrect flag format, providing : as a delimiter",
			map[string]string{"a": "1"},
			[]string{"a:2"},
			nil,
			fmt.Errorf("incorrect format for custom flag `a:2`"),
		},
		{"Incorrect flag format, providing space as a delimiter",
			map[string]string{"a": "1"}, []string{"a 2"},
			nil,
			fmt.Errorf("incorrect format for custom flag `a 2`"),
		},
		{"Incorrect flag format, providing - as a delimiter",
			map[string]string{"a": "1"},
			[]string{"a-2"},
			nil,
			fmt.Errorf("incorrect format for custom flag `a-2`"),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if err := MergeFlags(test.flags, test.overrides); err != nil {
				if test.wantErr == nil {
					t.Fatalf("failed to merge err: %v, existing flags: %v, set overrides: %v", err, test.flags, test.overrides)
				} else if !strings.EqualFold(err.Error(), test.wantErr.Error()) {
					t.Fatalf("inconsistent error return want: %v, got: %v", test.wantErr, err)
				}
				return
			}
			if !reflect.DeepEqual(test.want, test.flags) {
				t.Fatalf("error merging flags, want: %v, got: %v", test.want, test.flags)
			}
		})
	}
}
