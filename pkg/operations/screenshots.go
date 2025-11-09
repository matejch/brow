package operations

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

// ScreenshotOptions configures screenshot capture
type ScreenshotOptions struct {
	// FullPage captures the entire page instead of just the viewport
	FullPage bool
	// Quality for full-page screenshots (0-100, default 100)
	Quality int
}

// CaptureScreenshot captures a screenshot of the current page
func CaptureScreenshot(ctx context.Context, opts ScreenshotOptions) ([]byte, error) {
	var buf []byte
	var action chromedp.Action

	if opts.FullPage {
		quality := opts.Quality
		if quality == 0 {
			quality = 100
		}
		action = chromedp.FullScreenshot(&buf, quality)
	} else {
		action = chromedp.CaptureScreenshot(&buf)
	}

	if err := chromedp.Run(ctx, action); err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return buf, nil
}
