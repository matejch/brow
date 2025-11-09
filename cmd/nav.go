package cmd

import (
	"fmt"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/spf13/cobra"
)

var (
	waitReady bool
)

var navCmd = &cobra.Command{
	Use:   "nav <url>",
	Short: "Navigate to a URL",
	Long:  `Navigates the browser to the specified URL and waits for the page to load.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNav,
}

func init() {
	rootCmd.AddCommand(navCmd)
	navCmd.Flags().BoolVarP(&waitReady, "wait", "w", true, "Wait for page to be ready (default true)")
}

func runNav(_ *cobra.Command, args []string) error {
	url := args[0]

	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	result, err := browser.Page().Navigate(url, waitReady)
	if err != nil {
		return err
	}

	fmt.Printf("Navigated to: %s\n", url)
	fmt.Printf("Page title: %s\n", result.Title)

	return nil
}
