package get

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type TableFormat string

const (
	TableStyle    TableFormat = "table"
	MarkdownStyle TableFormat = "markdown"
	ListStyle     TableFormat = "list"
	// Direct ANSI Code Embedding
	Bold    = "\033[1m"
	FgGreen = "\033[32m"
	Reset   = "\033[0m"
)

// CreateToolTable creates table to show the avaiable CLI tools
func CreateToolsTable(tools Tools, format TableFormat) {

	var border tw.Border
	var settings tw.Settings
	var colWidth tw.CellWidth

	symbols := tw.NewSymbolCustom("Lines").WithRow("-").WithColumn("|")

	caption := tw.Caption{
		Text: fmt.Sprintf("There are %d tools, use `arkade get NAME` to download one.", len(tools)),
		Spot: tw.SpotBottomLeft,
	}

	switch format {
	case MarkdownStyle:
		symbols.WithCenter("|").
			WithMidLeft("|").
			WithMidRight("|")
		border.Left, border.Right = tw.On, tw.On
		border.Top, border.Bottom = tw.Off, tw.Off

	default:
		symbols.WithTopLeft("+").
			WithTopMid("+").
			WithTopRight("+").
			WithCenter("+").
			WithMidLeft("+").
			WithMidRight("+").
			WithBottomLeft("+").
			WithBottomMid("+").
			WithBottomRight("+")

		settings.Separators.BetweenRows = tw.On
		colWidth = tw.CellWidth{Global: 60}
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{
				Borders:  border,
				Symbols:  symbols,
				Settings: settings,
			})),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNone},
				ColMaxWidths: colWidth,
			},
		}),
	)
	table.Header([]string{"Tool", "Description"})
	table.Caption(caption)

	for _, t := range tools {

		name := Bold + FgGreen + t.Name + Reset
		url := fmt.Sprintf("https://github.com/%s/%s", t.Owner, t.Repo)

		if format == MarkdownStyle {
			name = fmt.Sprintf("[%s](%s)", t.Name, url)
		}
		table.Append([]string{name, t.Description})

	}

	table.Render()
}
