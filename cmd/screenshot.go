package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
	"github.com/spf13/cobra"
)

var (
	fullPage   bool
	outputFile string
	base64Out  bool
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot [output-file]",
	Short: "Capture a screenshot of the current page",
	Long: `Captures a screenshot of the current page.
If no output file is specified, outputs base64-encoded image data.
Use --full-page to capture the entire page instead of just the viewport.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScreenshot,
}

func init() {
	rootCmd.AddCommand(screenshotCmd)
	screenshotCmd.Flags().BoolVarP(&fullPage, "full-page", "f", false, "Capture full page (not just viewport)")
	screenshotCmd.Flags().BoolVarP(&base64Out, "base64", "b", false, "Output base64-encoded image data")
}

func runScreenshot(_ *cobra.Command, args []string) error {
	// Determine output file
	if len(args) > 0 {
		outputFile = args[0]
	}

	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	buf, err := browser.Page().Screenshot(operations.ScreenshotOptions{
		FullPage: fullPage,
		Quality:  100,
	})
	if err != nil {
		return err
	}

	// Handle output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, buf, 0644); err != nil {
			return fmt.Errorf("failed to write screenshot to file: %w", err)
		}
		fmt.Printf("Screenshot saved to: %s\n", outputFile)
	} else if base64Out {
		fmt.Println(base64.StdEncoding.EncodeToString(buf))
	} else {
		// Default: write to screenshot.png
		defaultFile := "screenshot.png"
		if err := os.WriteFile(defaultFile, buf, 0644); err != nil {
			return fmt.Errorf("failed to write screenshot to file: %w", err)
		}
		fmt.Printf("Screenshot saved to: %s\n", defaultFile)
	}

	return nil
}
