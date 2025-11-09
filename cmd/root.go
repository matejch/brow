package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Port is the Chrome remote debugging port (can be set via --port flag)
	Port int
)

var rootCmd = &cobra.Command{
	Use:   "brow",
	Short: "Simple CLI tools for browser automation via Chrome DevTools Protocol",
	Long: `brow provides a suite of composable CLI tools for browser automation.
It connects to Chrome running with remote debugging enabled (default port 9222).

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

	// Add persistent flag for Chrome debugging port
	rootCmd.PersistentFlags().IntVar(&Port, "port", 0, "Chrome remote debugging port (default 9222, or set BROW_DEBUG_PORT env var)")
}
