package get

import (
	"fmt"
	"os"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type TableFormat string

const (
	TableStyle    TableFormat = "table"
	MarkdownStyle TableFormat = "markdown"
	ListStyle     TableFormat = "list"
)

// CreateToolTable creates table to show the avaiable CLI tools
func CreateToolsTable(tools Tools, format TableFormat) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tool", "Description"})
	table.SetCaption(true,
		fmt.Sprintf("There are %d tools, use `arkade get NAME` to download one.", len(tools)))
	if format == MarkdownStyle {
		table.SetCaption(true)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetAutoWrapText(false)
	} else {
		table.SetRowLine(true)
		table.SetColWidth(60)
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{})
		table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}, tablewriter.Colors{})
	}

	for _, t := range tools {
		table.Append([]string{t.Name, t.Description})
	}

	table.Render()
}

type ArchMatrix struct {
	OS   string
	Arch string
}

func CreateCompatibilityToolsTable() {
	table := tablewriter.NewWriter(os.Stdout)

	arch := []ArchMatrix{
		{OS: "darwin", Arch: "x86_64"},
		{OS: "darwin", Arch: "arm64"},
		{OS: "linux", Arch: "x86_64"},
		{OS: "linux", Arch: "aarch64"},
		{OS: "linux", Arch: "armhf"},
		{OS: "ming", Arch: "x86_64"},
	}
	var header = []string{"Tool"}
	for _, v := range arch {
		header = append(header, v.OS+"/"+v.Arch)
	}
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	var compatible = '✅'
	var notCompatible = '❌'
	// var undefined = '🤷'

	tools := MakeTools()

	// check tools in parallel
	c := make(chan []string)
	var wg sync.WaitGroup
	for _, tool := range tools {
		wg.Add(1)
		go func(t Tool) {
			result := []string{t.Name}
			decision := string(compatible)
			for _, pair := range arch {
				quiet := true
				_, err := t.GetURL(pair.OS, pair.Arch, t.Version, quiet)
				if err != nil {
					decision = string(notCompatible)
				}

				// status, _, headers, err := t.Head(url)
				// if err != nil {
				// 	decision = string(undefined)
				// 	fmt.Println("Head failed for ", t.Name)
				// }

				// if status != http.StatusOK {
				// 	decision = string(notCompatible)
				// 	fmt.Println("StatusNotOK for ", t.Name)
				// }
				result = append(result, decision)
			}
			c <- result
			defer wg.Done()
		}(tool)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	// process data
	for row := range c {
		table.Append(row)
	}
	table.Render()
}
