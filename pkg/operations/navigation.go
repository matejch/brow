package operations

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

// NavigationResult holds the results of a navigation operation
type NavigationResult struct {
	URL   string
	Title string
}

// Navigate navigates to the specified URL and optionally waits for the page to be ready
func Navigate(ctx context.Context, url string, waitReady bool) (*NavigationResult, error) {
	var title string
	actions := []chromedp.Action{
		chromedp.Navigate(url),
	}

	if waitReady {
		actions = append(actions, chromedp.WaitReady("body"))
	}

	actions = append(actions, chromedp.Title(&title))

	if err := chromedp.Run(ctx, actions...); err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	return &NavigationResult{
		URL:   url,
		Title: title,
	}, nil
}
