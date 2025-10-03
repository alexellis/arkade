package helm

import (
	"testing"
)

func Test_isURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://example.com/values.yaml", true},
		{"http://example.com/values.yaml", true},
		{"ftp://example.com/values.yaml", true},
		{"/local/path/values.yaml", false},
		{"relative/path/values.yaml", false},
		{"values.yaml", false},
		{"", false},
		{"not-a-url", false},
	}

	for _, test := range tests {
		result := isURL(test.input)
		if result != test.expected {
			t.Errorf("isURL(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func Test_processValuesFiles(t *testing.T) {
	input := []string{
		"https://example.com/values1.yaml",
		"/local/values2.yaml",
		"relative/values3.yaml",
	}

	result, err := processValuesFiles(input)
	if err != nil {
		t.Fatalf("processValuesFiles failed: %v", err)
	}

	// Should return the same files unchanged since Helm handles URLs natively
	expected := input
	if len(result) != len(expected) {
		t.Fatalf("Expected %d files, got %d", len(expected), len(result))
	}

	for i, file := range result {
		if file != expected[i] {
			t.Errorf("File %d: expected %q, got %q", i, expected[i], file)
		}
	}
}
