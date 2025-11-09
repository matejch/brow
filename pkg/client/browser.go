package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/matejch/brow/pkg/config"
)

const (
	defaultHost = "localhost"
)

// TabInfo contains metadata about a browser tab
type TabInfo struct {
	Index    int
	TargetID string
	Title    string
	URL      string
}

// tabContext holds the context and metadata for a single tab
type tabContext struct {
	targetID target.ID
	title    string
	url      string
	ctx      context.Context
	cancel   context.CancelFunc
}

// Browser represents a connection to a Chrome browser instance
type Browser struct {
	config      *config.Config
	allocCtx    context.Context
	allocCancel context.CancelFunc
	tabs        []*tabContext
	mu          sync.RWMutex // Protects tabs slice
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

	// Discover ALL page targets (not just first)
	tabs := make([]*tabContext, 0)
	for _, t := range targets {
		if t.Type == "page" {
			// Create context for this tab
			tabCtx, tabCancel := chromedp.NewContext(allocCtx, chromedp.WithTargetID(t.TargetID))

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

			tabs = append(tabs, &tabContext{
				targetID: t.TargetID,
				title:    t.Title,
				url:      t.URL,
				ctx:      tabCtx,
				cancel:   tabCancel,
			})
		}
	}

	if len(tabs) == 0 {
		allocCancel()
		return nil, fmt.Errorf("no page tabs available")
	}

	return &Browser{
		config:      cfg,
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		tabs:        tabs,
	}, nil
}

// Page returns a Page instance for interacting with the first tab (backward compatible)
func (b *Browser) Page() *Page {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.tabs) == 0 {
		return nil
	}

	return &Page{
		ctx:    b.tabs[0].ctx,
		config: b.config,
	}
}

// TabCount returns the number of open tabs
func (b *Browser) TabCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.tabs)
}

// Tabs returns metadata about all open tabs
func (b *Browser) Tabs() ([]*TabInfo, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	tabs := make([]*TabInfo, len(b.tabs))
	for i, tab := range b.tabs {
		tabs[i] = &TabInfo{
			Index:    i,
			TargetID: string(tab.targetID),
			Title:    tab.title,
			URL:      tab.url,
		}
	}
	return tabs, nil
}

// TabByIndex returns a Page instance for the tab at the specified index
func (b *Browser) TabByIndex(index int) (*Page, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if index < 0 || index >= len(b.tabs) {
		return nil, fmt.Errorf("tab index %d out of range (have %d tabs)", index, len(b.tabs))
	}

	return &Page{
		ctx:    b.tabs[index].ctx,
		config: b.config,
	}, nil
}

// NewTab creates a new tab and navigates to the specified URL (empty string for blank tab)
func (b *Browser) NewTab(url string) (*Page, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Create a new tab context (chromedp will create the tab on first use)
	tabCtx, tabCancel := chromedp.NewContext(b.allocCtx)

	// Add timeout if specified
	if b.config.Timeout > 0 {
		var timeoutCancel context.CancelFunc
		tabCtx, timeoutCancel = context.WithTimeout(tabCtx, b.config.Timeout)
		// Wrap the cancel function to call both
		oldTabCancel := tabCancel
		tabCancel = func() {
			timeoutCancel()
			oldTabCancel()
		}
	}

	// Create tab context
	newTab := &tabContext{
		targetID: "", // Will be populated when tab is actually created
		title:    "",
		url:      url,
		ctx:      tabCtx,
		cancel:   tabCancel,
	}

	b.tabs = append(b.tabs, newTab)

	page := &Page{
		ctx:    tabCtx,
		config: b.config,
	}

	// If URL provided, navigate to it (this creates the tab)
	if url != "" {
		if _, err := page.Navigate(url, true); err != nil {
			// Remove the tab if navigation failed
			b.tabs = b.tabs[:len(b.tabs)-1]
			tabCancel()
			return nil, fmt.Errorf("failed to navigate new tab: %w", err)
		}
	} else {
		// Navigate to blank page to create the tab
		if _, err := page.Navigate("about:blank", false); err != nil {
			b.tabs = b.tabs[:len(b.tabs)-1]
			tabCancel()
			return nil, fmt.Errorf("failed to create new tab: %w", err)
		}
	}

	return page, nil
}

// CloseTab closes the tab at the specified index
func (b *Browser) CloseTab(index int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if index < 0 || index >= len(b.tabs) {
		return fmt.Errorf("tab index %d out of range (have %d tabs)", index, len(b.tabs))
	}

	// Cancel the tab's context (this closes the tab)
	if b.tabs[index].cancel != nil {
		b.tabs[index].cancel()
	}

	// Remove from slice
	b.tabs = append(b.tabs[:index], b.tabs[index+1:]...)

	return nil
}

// Context returns the underlying context for the first tab (backward compatible)
func (b *Browser) Context() context.Context {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.tabs) == 0 {
		return nil
	}
	return b.tabs[0].ctx
}

// SetTimeout updates the timeout for all tab operations
func (b *Browser) SetTimeout(timeout time.Duration) {
	if timeout <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Apply timeout to all tabs
	for _, tab := range b.tabs {
		ctx, cancel := context.WithTimeout(tab.ctx, timeout)
		// Store old cancel to call both
		oldCancel := tab.cancel
		tab.cancel = func() {
			cancel()
			if oldCancel != nil {
				oldCancel()
			}
		}
		tab.ctx = ctx
	}
}

// Close cleanly shuts down the browser connection
// Note: This does NOT close the Chrome browser itself, only disconnects from it
func (b *Browser) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Note: We intentionally do NOT call tab cancel functions here
	// because that would close the tabs in Chrome
	// We only disconnect from the remote debugging session
	// This preserves the tab state for subsequent brow commands

	if b.allocCancel != nil {
		b.allocCancel()
	}

	return nil
}
