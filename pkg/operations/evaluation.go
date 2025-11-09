package operations

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

// Evaluate executes JavaScript in the page context and returns the result
func Evaluate(ctx context.Context, script string) (interface{}, error) {
	var result interface{}
	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &result)); err != nil {
		return nil, fmt.Errorf("failed to evaluate JavaScript: %w", err)
	}
	return result, nil
}
