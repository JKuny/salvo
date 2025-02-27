// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the application",
	Long:  `Print the version of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s\n", VERSION)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
