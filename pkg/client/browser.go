package client

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/matejch/brow/pkg/config"
)

const (
	defaultHost = "localhost"
)

// Browser represents a connection to a Chrome browser instance
type Browser struct {
	config      *config.Config
	allocCtx    context.Context
	allocCancel context.CancelFunc
	tabCtx      context.Context
	tabCancel   context.CancelFunc
}

// New creates a new Browser instance that connects to an existing Chrome instance
func New(cfg *config.Config) (*Browser, error) {
	if cfg == nil {
		cfg = config.Default()
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	debugURL := fmt.Sprintf("http://%s:%d", defaultHost, cfg.Port)

	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), debugURL)

	// Create a temporary context to query targets
	tempCtx, tempCancel := chromedp.NewContext(allocCtx)

	// Get existing targets
	targets, err := chromedp.Targets(tempCtx)
	tempCancel() // Clean up temp context immediately

	if err != nil {
		allocCancel()
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	if len(targets) == 0 {
		allocCancel()
		return nil, fmt.Errorf("no tabs available - please start Chrome first with 'brow start'")
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
		return nil, fmt.Errorf("no page tabs available")
	}

	// Attach to the existing tab
	tabCtx, tabCancel := chromedp.NewContext(allocCtx, chromedp.WithTargetID(targetID))

	// Add timeout if specified
	if cfg.Timeout > 0 {
		var timeoutCancel context.CancelFunc
		tabCtx, timeoutCancel = context.WithTimeout(tabCtx, cfg.Timeout)
		// Wrap the cancel function to call both
		oldTabCancel := tabCancel
		tabCancel = func() {
			timeoutCancel()
			oldTabCancel()
		}
	}

	return &Browser{
		config:      cfg,
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		tabCtx:      tabCtx,
		tabCancel:   tabCancel,
	}, nil
}

// Page returns a Page instance for interacting with the current tab
func (b *Browser) Page() *Page {
	return &Page{
		ctx:    b.tabCtx,
		config: b.config,
	}
}

// Context returns the underlying context for advanced usage
func (b *Browser) Context() context.Context {
	return b.tabCtx
}

// SetTimeout updates the timeout for browser operations
// Returns a new context with the timeout applied
func (b *Browser) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(b.tabCtx, timeout)
		// Store old cancel to call both
		oldCancel := b.tabCancel
		b.tabCancel = func() {
			cancel()
			if oldCancel != nil {
				oldCancel()
			}
		}
		b.tabCtx = ctx
	}
}

// Close cleanly shuts down the browser connection
// Note: This does NOT close the Chrome browser itself, only disconnects from it
func (b *Browser) Close() error {
	// Note: We intentionally do NOT call tabCancel() here because that would close the tab
	// We only disconnect from the remote debugging session by calling allocCancel()
	// This preserves the tab state for subsequent brow commands
	if b.allocCancel != nil {
		b.allocCancel()
	}
	return nil
}
