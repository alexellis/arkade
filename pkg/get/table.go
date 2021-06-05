package get

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// CreateToolTable creates table to show the avaiable CLI tools
func CreateToolsTable(tools Tools) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetColWidth(60)
	table.SetHeader([]string{"Tool", "Description"})
	table.SetCaption(true, "Use 'arkade get TOOL' to download a tool or application.")

	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}, tablewriter.Colors{})

	for _, t := range tools {
		table.Append([]string{t.Name, t.Description})
	}

	table.Render()
}
