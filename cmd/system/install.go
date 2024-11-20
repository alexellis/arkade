// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func MakeInstall() *cobra.Command {

	command := &cobra.Command{
		Use:     "install",
		Short:   "Install system apps",
		Long:    `Install system apps for Linux hosts`,
		Aliases: []string{"i"},
		Example: `  arkade system install [APP]
  arkade system install --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeInstallGo())
	command.AddCommand(MakeInstallFirecracker())
	command.AddCommand(MakeInstallPrometheus())
	command.AddCommand(MakeInstallCNI())
	command.AddCommand(MakeInstallContainerd())
	command.AddCommand(MakeInstallActionsRunner())
	command.AddCommand(MakeInstallNode())
	command.AddCommand(MakeInstallTCRedirectTap())
	command.AddCommand(MakeInstallRegistry())
	command.AddCommand(MakeInstallGitLabRunner())
	command.AddCommand((MakeInstallBuildkitd()))
	command.AddCommand(MakeInstallPowershell())
	command.AddCommand(MakeInstallCaddyServer())

	command.AddCommand(MakeInstallNodeExporter())

	return command
}

func CreateSystemTable(sys []*cobra.Command) string {

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"System Install", "Description"})
	table.SetCaption(true,
		fmt.Sprintf("\nThere are %d system installations available.\n", len(sys)))

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	var sysMap = make(map[string]string, len(sys))
	sortedList := make([]string, 0, len(sys))

	for _, s := range sys {
		sysMap[s.Use] = s.Short
		sortedList = append(sortedList, s.Use)
	}
	sort.Strings(sortedList)

	for _, sysInst := range sortedList {
		table.Append([]string{sysInst, sysMap[sysInst]})
	}

	table.Render()

	return tableString.String()
}
