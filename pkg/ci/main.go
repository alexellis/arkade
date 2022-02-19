package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
	"text/template"
	"time"

	"github.com/alexellis/arkade/pkg/get"
)

func URLExists(name, url, version string) error {
	fmt.Println("--------------->>>>>>>>>>>>>>>>>>>>>>>>>", name)
	timeout := time.Second * 10
	client := get.MakeHTTPClient(&timeout, false)

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("tool %s with version %s not found", name, version)
	}
	return nil
}

func main() {
	tools := get.MakeTools()
	sort.Sort(tools)
	var errorTools []string

	for _, tool := range tools {
		url, err := get.GetDownloadURL(&tool, "linux", "x86_64", "")
		if err != nil {
			errorTools = append(errorTools, tool.Name)
			continue
		}
		err = URLExists(tool.Name, url, "")
		if err != nil {
			errorTools = append(errorTools, tool.Name)
		}
	}

	if len(errorTools) > 0 {
		t := template.New("List of tools with errors")
		t.Parse(`===========================================================================
List of tools that encountered errors:
{{- range .}}
{{. -}}
{{- end }}
`)
		var tpl bytes.Buffer
		err := t.Execute(&tpl, errorTools)
		if err != nil {
			panic(err)
		}
		log.Fatalf("%s", tpl.Bytes())
	}
}
