package cmd

import (
	"fmt"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/spf13/cobra"
)

var (
	xpath bool
)

var pickCmd = &cobra.Command{
	Use:   "pick",
	Short: "Interactive element picker to get CSS selectors",
	Long: `Injects an interactive overlay into the page that allows you to click on elements
and get their CSS selector or XPath.

The overlay highlights elements on hover and copies the selector on click.
Press ESC to exit the picker mode.`,
	RunE: runPick,
}

func init() {
	rootCmd.AddCommand(pickCmd)
	pickCmd.Flags().BoolVarP(&xpath, "xpath", "x", false, "Return XPath instead of CSS selector")
}

func runPick(_ *cobra.Command, _ []string) error {
	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	if err := browser.Page().InjectPicker(xpath); err != nil {
		return err
	}

	fmt.Println("Element picker activated!")
	fmt.Println("Hover over elements to highlight, click to select, press ESC to exit.")
	fmt.Println("")
	fmt.Println("After selecting an element, run:")
	fmt.Println("  brow eval 'window.__browPickedSelector'")
	fmt.Println("")
	fmt.Println("To get the selected element's selector.")

	return nil
}
