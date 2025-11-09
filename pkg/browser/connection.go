package browser

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

const (
	// DefaultPort is the default Chrome DevTools Protocol port
	DefaultPort = 9222
	// DefaultHost is the default host for Chrome remote debugging
	DefaultHost = "localhost"
)

// ResolvePort determines which port to use based on command-line flag, environment variable, and default
// Priority: flagPort (if != 0) > BROW_DEBUG_PORT env var > DefaultPort
func ResolvePort(flagPort int) int {
	// If flag was explicitly set, use it
	if flagPort != 0 {
		return flagPort
	}

	// Check environment variable
	if envPort := os.Getenv("BROW_DEBUG_PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil && port > 0 {
			return port
		}
	}

	// Fall back to default
	return DefaultPort
}

// GetRemoteAllocator creates a new allocator that connects to an existing Chrome instance
func GetRemoteAllocator(port int) (context.Context, context.CancelFunc, error) {
	debugURL := fmt.Sprintf("http://%s:%d", DefaultHost, port)

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), debugURL)

	return allocCtx, cancel, nil
}

// GetExistingTabContext attaches to an existing browser tab without creating a new one
// Returns a context that should NOT be cancelled if you want to keep the tab open
// Only cancels the allocator context, not the tab context itself
func GetExistingTabContext(port int) (context.Context, context.CancelFunc, error) {
	allocCtx, allocCancel, err := GetRemoteAllocator(port)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote allocator: %w", err)
	}

	// Create a temporary context to query targets
	tempCtx, tempCancel := chromedp.NewContext(allocCtx)

	// Get existing targets
	targets, err := chromedp.Targets(tempCtx)
	tempCancel() // Clean up temp context immediately

	if err != nil {
		allocCancel()
		return nil, nil, fmt.Errorf("failed to get targets: %w", err)
	}

	if len(targets) == 0 {
		allocCancel()
		return nil, nil, fmt.Errorf("no tabs available - please start Chrome first with 'brow start'")
	}

	// Find the first page target (ignore background pages, extensions, etc.)
	var targetID target.ID
	for _, t := range targets {
		if t.Type == "page" {
			targetID = t.TargetID
			break
		}
	}

	if targetID == "" {
		allocCancel()
		return nil, nil, fmt.Errorf("no page tabs available")
	}

	// Attach to the existing tab - this does NOT create a new tab
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithTargetID(targetID))

	// CRITICAL: Only cancel the allocator to disconnect from the debugging session
	// Do NOT call cancel() on the tab context - that would close the tab!
	// The tab must remain open for subsequent commands to work
	cancelFunc := func() {
		// cancel()  // Intentionally NOT calling this - it would close the tab
		allocCancel() // Only disconnect from remote debugging
	}

	// Prevent unused variable warning
	_ = cancel

	return ctx, cancelFunc, nil
}
