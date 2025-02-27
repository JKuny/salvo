// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"salvo/cmd"
)

// TestVersionCommand tests the "version" command
func TestVersionCommand(t *testing.T) {
	// Create a buffer to capture output
	outputBuffer := &bytes.Buffer{}

	// Set the buffer as the output for the root command
	cmd.RootCmd.SetOut(outputBuffer)

	// Simulate passing "version" as an argument to the root command
	cmd.RootCmd.SetArgs([]string{"version"})

	// Execute the root command
	err := cmd.RootCmd.Execute()
	if err != nil {
		t.Fatalf("Error executing root command with version subcommand: %v", err)
	}

	// Capture the output
	actualOutput := strings.TrimSpace(outputBuffer.String())
	expectedOutput := "0.0.1"

	// Assert the captured output matches the expected version
	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output: got %q, want %q", actualOutput, expectedOutput)
	}
}
