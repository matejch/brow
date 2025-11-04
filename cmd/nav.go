package cmd

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/matej/brow/pkg/browser"
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

func runNav(cmd *cobra.Command, args []string) error {
	url := args[0]

	var title string
	actions := []chromedp.Action{
		chromedp.Navigate(url),
	}

	if waitReady {
		actions = append(actions, chromedp.WaitReady("body"))
	}

	actions = append(actions, chromedp.Title(&title))

	if err := browser.Run(actions...); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}

	fmt.Printf("Navigated to: %s\n", url)
	fmt.Printf("Page title: %s\n", title)

	return nil
}
