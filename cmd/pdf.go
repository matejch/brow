package cmd

import (
	"fmt"
	"os"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
	"github.com/spf13/cobra"
)

var (
	landscape bool
	printBg   bool
	pdfOutput string
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

func runPDF(_ *cobra.Command, args []string) error {
	// Determine output file
	if len(args) > 0 {
		pdfOutput = args[0]
	} else {
		pdfOutput = "output.pdf"
	}

	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	buf, err := browser.Page().PDF(operations.PDFOptions{
		Landscape:       landscape,
		PrintBackground: printBg,
	})
	if err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile(pdfOutput, buf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF to file: %w", err)
	}

	fmt.Printf("PDF saved to: %s\n", pdfOutput)
	return nil
}
