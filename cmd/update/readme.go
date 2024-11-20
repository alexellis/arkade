// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package update

import (
	"fmt"
	"os"
	"sort"
	"strings"

	install "github.com/alexellis/arkade/cmd"
	system "github.com/alexellis/arkade/cmd/system"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeReadme() *cobra.Command {
	var command = &cobra.Command{
		Use:   "readme",
		Short: "Update the system / tools / apps tables in the readme",
		Long: `
A command for contributors to use when adding or updating any of the developer tools
to quickly perform an in-place update of the various tables within the readme .`,
		Example:       `  arkade update readme`,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	command.RunE = func(cmd *cobra.Command, args []string) error {

		var readmeTables = make(map[string]string)

		//update system installs <!-- start system content -->
		systemList := system.MakeInstall().Commands()
		readmeTables["system"] = system.CreateSystemTable(systemList)

		//update apps  <!-- start apps content -->
		appList := install.GetApps()
		readmeTables["apps"] = install.CreateAppsTable(appList)

		//update tools  <!-- start tools content -->
		tools := get.MakeTools()
		sort.Sort(tools)
		readmeTables["tools"] = get.CreateToolsTable(tools, get.MarkdownStyle)

		return writeTableToReadme(readmeTables)

	}
	return command
}

func writeTableToReadme(tables map[string]string) error {

	filePath := "README.md"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}
	content := string(fileContent)

	for k, v := range tables {

		startMarker := fmt.Sprintf("<!-- start %s content -->", k)
		endMarker := fmt.Sprintf("<!-- end %s content -->", k)

		startIdx := strings.Index(content, startMarker)
		endIdx := strings.Index(content, endMarker)
		if startIdx == -1 || endIdx == -1 || startIdx > endIdx {
			return fmt.Errorf("%s readme markers not found or are in incorrect order", k)
		}

		content = content[:startIdx+len(startMarker)] + "\n" +
			v + "\n" +
			content[endIdx:]

		fmt.Printf("Updated %s table.\n", k)
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("README tables updated successfully")
	return nil
}
