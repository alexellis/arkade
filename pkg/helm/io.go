package helm

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ValuesMap is an alias for map[string]interface{}
type ValuesMap map[interface{}]interface{}

// Load a values.yaml file and return a ValuesMap with the keys
// and values from the YAML file as a map[string]interface{}
func Load(yamlPath string) (ValuesMap, error) {
	body, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s, error: %s", yamlPath, err)
	}

	values := ValuesMap{}

	err = yaml.Unmarshal(body, &values)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s, error: %s", yamlPath, err)
	}

	return values, nil
}

// FilterImagesUptoDepth takes a ValuesMap and returns a map of images that
// were found upto max level
func FilterImagesUptoDepth(values ValuesMap, depth int) map[string]string {
	images := map[string]string{}

	for k, v := range values {
		if k == "image" {
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
