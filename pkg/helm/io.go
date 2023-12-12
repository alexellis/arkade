package helm

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// ValuesMap is an alias for map[string]interface{}
type ValuesMap map[string]interface{}

// Load a values.yaml file and return a ValuesMap with the keys
// and values from the YAML file as a map[string]interface{}
func Load(yamlPath string) (ValuesMap, error) {
	body, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s, error: %s", yamlPath, err)
	}

	values := ValuesMap{}

	if err = yaml.Unmarshal(body, &values); err != nil {
		return nil, fmt.Errorf("unable to parse %s, error: %s", yamlPath, err)
	}

	return values, nil
}

// LoadFrom loads a values.yaml snippet from memory and
// returns a ValuesMap with the keys and values from the YAML
func LoadFrom(yamlText string) (ValuesMap, error) {

	values := ValuesMap{}

	if err := yaml.Unmarshal([]byte(yamlText), &values); err != nil {
		return nil, fmt.Errorf("unable to parse %s, error: %s", yamlText, err)
	}

	return values, nil
}

// ReplaceValuesInHelmValuesFile takes a values.yaml file and replaces values in it with the values provided in the map
// and returns the updated values.yaml file as a string
func ReplaceValuesInHelmValuesFile(values map[string]string, yamlPath string) (string, error) {
	readFile, err := os.ReadFile(yamlPath)
	if err != nil {
		return "", err
	}

	fileContent := string(readFile)
	for k, v := range values {
		fileContent = strings.ReplaceAll(fileContent, k, v)
	}
	return fileContent, nil
}

// FilterImagesUptoDepth takes a ValuesMap and returns a map of images that
// were found upto max level
func FilterImagesUptoDepth(values ValuesMap, depth int) map[string]string {
	images := map[string]string{}

	for k, v := range values {

		if k == "image" && reflect.TypeOf(v).Kind() == reflect.String {
			imageUrl := v.(string)
			images[imageUrl] = imageUrl
		}

		if c, ok := v.(ValuesMap); ok && depth > 0 {
			images = mergeMaps(images, FilterImagesUptoDepth(c, depth-1))
		}
	}
	return images
}

func mergeMaps(original, latest map[string]string) map[string]string {
	for k, v := range latest {
		original[k] = v
	}
	return original
}
