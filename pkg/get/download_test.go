package get

import (
	"os"
	"testing"
)

func TestCopyFileP(t *testing.T) {
	// Create a temporary source file
	srcFile, err := os.CreateTemp("", "src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	// Write some data to the source file
	data := []byte("Hello, World!")
	if _, err := srcFile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := srcFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a temporary destination file
	dstFile, err := os.CreateTemp("", "dst")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dstFile.Name())

	// Use the function to copy the file
	if _, err := CopyFileP(srcFile.Name(), dstFile.Name(), 0644); err != nil {
		t.Fatal(err)
	}

	// Read the destination file
	gotData, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check if the data matches
	if string(gotData) != string(data) {
		t.Fatalf("got %q, want %q", string(gotData), string(data))
	}
}
