// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "salvo",
	Short: "A CLI for getting Kubernetes logs fast",
	Long: `Salvo is a CLI that uses your local machine's Kubernetes configuration
in order to write your pod logs to a directory for local inspection.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
