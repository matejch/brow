package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/matej/brow/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	landscape    bool
	printBg      bool
	pdfOutput    string
)

var pdfCmd = &cobra.Command{
	Use:   "pdf [output-file]",
	Short: "Export the current page as PDF",
	Long: `Generates a PDF from the current page.
If no output file is specified, saves to 'output.pdf'.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPDF,
}

func init() {
	rootCmd.AddCommand(pdfCmd)
	pdfCmd.Flags().BoolVarP(&landscape, "landscape", "l", false, "Use landscape orientation")
	pdfCmd.Flags().BoolVarP(&printBg, "background", "b", true, "Print background graphics (default true)")
}

func runPDF(cmd *cobra.Command, args []string) error {
	// Determine output file
	if len(args) > 0 {
		pdfOutput = args[0]
	} else {
		pdfOutput = "output.pdf"
	}

	var buf []byte

	if err := browser.Run(chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		buf, _, err = page.PrintToPDF().
			WithPrintBackground(printBg).
			WithLandscape(landscape).
			Do(ctx)
		return err
	})); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Write to file
	if err := os.WriteFile(pdfOutput, buf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF to file: %w", err)
	}

	fmt.Printf("PDF saved to: %s\n", pdfOutput)
	return nil
}
