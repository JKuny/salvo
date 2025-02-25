// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Confirms your Kubernetes configuration is setup",
	Long: `Cobra is a CLI library for Go that empowers applications.
           This application is a tool to generate the needed files
           to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting Kubernetes config")

		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			os.Exit(1)
		}

		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		fmt.Printf("Using Kubernetes config file: %s\n", kubeConfigPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
