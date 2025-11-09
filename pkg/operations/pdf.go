package operations

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// PDFOptions configures PDF generation
type PDFOptions struct {
	// Landscape orientation (default false)
	Landscape bool
	// PrintBackground includes background graphics (default true)
	PrintBackground bool
}

// GeneratePDF generates a PDF from the current page
func GeneratePDF(ctx context.Context, opts PDFOptions) ([]byte, error) {
	var buf []byte

	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		buf, _, err = page.PrintToPDF().
			WithPrintBackground(opts.PrintBackground).
			WithLandscape(opts.Landscape).
			Do(ctx)
		return err
	})); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf, nil
}
