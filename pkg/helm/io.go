package helm

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ValuesMap is an alias for map[string]interface{}
type ValuesMap map[string]interface{}

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

// FilterImages takes a ValuesMap and returns a map of images that
// were found at the top level, or one level down with keys of
// "image: "
func FilterImages(values ValuesMap) map[string]string {
	images := map[string]string{}

	for k, v := range values {

		// Match anything at the top level called "image: ..."
		if k == "image" {
			images[k] = v.(string)
		}

		// Match anything at one level down i.e. "gateway.image: ..."
		if c, ok := v.(map[interface{}]interface{}); ok && c != nil {
			for kk, vv := range c {
				if kk == "image" {
					images[k] = vv.(string)
				}
			}
		}
	}

	return images
}
