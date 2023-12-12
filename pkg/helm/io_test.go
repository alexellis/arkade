package helm

import (
	"testing"
)

func Test_LoadFrom(t *testing.T) {
	yamlText := `faas-netes:
  image: openfaas/faas-netes:0.1.0
`

	imageWant := `openfaas/faas-netes:0.1.0`

	values, err := LoadFrom(yamlText)
	if err != nil {
		t.Fatal(err)
	}

	imageGot := values["faas-netes"].(ValuesMap)["image"]
	if imageGot != imageWant {
		t.Fatalf("got %q, want %q", imageGot, imageWant)
	}
}
