package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "brow",
	Short: "Simple CLI tools for browser automation via Chrome DevTools Protocol",
	Long: `brow provides a suite of composable CLI tools for browser automation.
It connects to Chrome running with remote debugging enabled on port 9222.

Philosophy: Simple, composable tools that output text for easy agent consumption.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
