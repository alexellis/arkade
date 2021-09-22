// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import (
	"bytes"
	"html/template"
	"os"
	"path"
)

func writeFile(name string, data []byte) (string, error) {
	d := os.TempDir()
	p := path.Join(d, "issuer.yaml")

	err := os.WriteFile(p, data, os.ModePerm)

	return p, err
}

func templateManifest(templateSt string, data interface{}) ([]byte, error) {
	tmpl, err := template.New("yaml").Parse(templateSt)

	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, data)

	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), err
}
