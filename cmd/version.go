// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the application",
	Long:  `Print the version of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
