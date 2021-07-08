package get

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableFormat string

const (
	TableStyle    TableFormat = "table"
	MarkdownStyle TableFormat = "markdown"
)

// CreateToolTable creates table to show the avaiable CLI tools
func CreateToolsTable(tools Tools, format TableFormat) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tool", "Description"})
	table.SetCaption(true, "Use 'arkade get TOOL' to download a tool or application.")
	if format == MarkdownStyle {
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
