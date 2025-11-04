package browser

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

const (
	// DefaultPort is the default Chrome DevTools Protocol port
	DefaultPort = 9222
	// DefaultHost is the default host for Chrome remote debugging
	DefaultHost = "localhost"
)

// GetRemoteAllocator creates a new allocator that connects to an existing Chrome instance
func GetRemoteAllocator() (context.Context, context.CancelFunc, error) {
	debugURL := fmt.Sprintf("http://%s:%d", DefaultHost, DefaultPort)

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), debugURL)

	return allocCtx, cancel, nil
}

// NewContext creates a new Chrome context connected to the remote debugging instance
func NewContext() (context.Context, context.CancelFunc, error) {
	allocCtx, allocCancel, err := GetRemoteAllocator()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote allocator: %w", err)
	}

	ctx, cancel := chromedp.NewContext(allocCtx)

	// Create a combined cancel function
	combinedCancel := func() {
		cancel()
		allocCancel()
	}

	return ctx, combinedCancel, nil
}

// Run executes chromedp actions with proper context management
func Run(actions ...chromedp.Action) error {
	ctx, cancel, err := NewContext()
	if err != nil {
		return err
	}
	defer cancel()

	return chromedp.Run(ctx, actions...)
}
