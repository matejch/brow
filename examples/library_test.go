package examples_test

import (
	"testing"
	"time"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
)

// TestLibraryUsageBasic demonstrates basic usage of the brow library
func TestLibraryUsageBasic(t *testing.T) {
	// Note: This test requires Chrome to be running with remote debugging enabled
	// Start Chrome first with: ./brow start

	// Create a browser instance
	browser, err := client.New(&config.Config{
		Port:    9222,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Skip("Skipping test: Chrome not running. Start with './brow start'")
	}
	defer browser.Close()

	// Get the page
	page := browser.Page()

	// Navigate to a test website
	result, err := page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatalf("navigation failed: %v", err)
	}

	// Verify the page title
	if result.Title != "Example Domain" {
		t.Errorf("expected title 'Example Domain', got %q", result.Title)
	}

	t.Logf("Successfully navigated to %s with title: %s", result.URL, result.Title)
}

// TestLibraryUsageJavaScript demonstrates JavaScript evaluation
func TestLibraryUsageJavaScript(t *testing.T) {
	browser, err := client.New(nil) // Use default config
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate first
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Evaluate JavaScript
	title, err := page.Eval("document.title")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Page title from JavaScript: %v", title)

	// Count links on the page
	linkCount, err := page.Eval("document.querySelectorAll('a').length")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Number of links on page: %v", linkCount)
}

// TestLibraryUsageScreenshot demonstrates screenshot capture
func TestLibraryUsageScreenshot(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate first
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Capture viewport screenshot
	screenshot, err := page.Screenshot(operations.ScreenshotOptions{
		FullPage: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(screenshot) == 0 {
		t.Error("screenshot data is empty")
	}

	t.Logf("Captured screenshot: %d bytes", len(screenshot))

	// Capture full-page screenshot
	fullScreenshot, err := page.Screenshot(operations.ScreenshotOptions{
		FullPage: true,
		Quality:  90,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Captured full-page screenshot: %d bytes", len(fullScreenshot))
}

// TestLibraryUsageCookies demonstrates cookie management
func TestLibraryUsageCookies(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate first
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Set a cookie
	err = page.SetCookie("test_cookie=hello; path=/")
	if err != nil {
		t.Fatal(err)
	}

	// Get all cookies
	cookies, err := page.GetCookies("")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Found %d cookies", len(cookies))

	// Verify our cookie was set
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "test_cookie" && cookie.Value == "hello" {
			found = true
			break
		}
	}

	if !found {
		t.Error("test_cookie was not found after setting")
	}
}

// TestLibraryUsageStorage demonstrates localStorage/sessionStorage
func TestLibraryUsageStorage(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate first
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Set localStorage items
	err = page.SetStorageItem(operations.LocalStorage, "test_key", "test_value")
	if err != nil {
		t.Fatal(err)
	}

	// Get the item back
	value, err := page.GetStorageItem(operations.LocalStorage, "test_key")
	if err != nil {
		t.Fatal(err)
	}

	if value != "test_value" {
		t.Errorf("expected 'test_value', got %v", value)
	}

	// Get all localStorage items
	allItems, err := page.GetAllStorage(operations.LocalStorage)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("localStorage contains %d items", len(allItems))

	// Clean up
	err = page.RemoveStorageItem(operations.LocalStorage, "test_key")
	if err != nil {
		t.Fatal(err)
	}
}

// TestLibraryUsagePDF demonstrates PDF generation
func TestLibraryUsagePDF(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate first
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Generate PDF
	pdf, err := page.PDF(operations.PDFOptions{
		Landscape:       false,
		PrintBackground: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(pdf) == 0 {
		t.Error("PDF data is empty")
	}

	t.Logf("Generated PDF: %d bytes", len(pdf))
}

// TestLibraryUsageMultipleTabs demonstrates working with the same browser across multiple operations
func TestLibraryUsageMultipleTabs(t *testing.T) {
	browser, err := client.New(&config.Config{
		Port:    9222,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate to first site
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Get title
	title1, err := page.Eval("document.title")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("First page title: %v", title1)

	// Navigate to second site (same tab)
	_, err = page.Navigate("https://example.org", true)
	if err != nil {
		t.Fatal(err)
	}

	// Get title
	title2, err := page.Eval("document.title")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Second page title: %v", title2)

	// Verify we're on a different page
	if title1 == title2 {
		t.Error("titles should be different after navigation")
	}
}

// TestMultiTabSupport demonstrates working with multiple tabs
func TestMultiTabSupport(t *testing.T) {
	browser, err := client.New(&config.Config{
		Port:    9222,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	// Check initial tab count
	initialCount := browser.TabCount()
	t.Logf("Initial tab count: %d", initialCount)

	if initialCount == 0 {
		t.Fatal("Expected at least one tab")
	}

	// List all tabs
	tabs, err := browser.Tabs()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Found %d tabs", len(tabs))
	for _, tab := range tabs {
		t.Logf("  Tab %d: %s (%s)", tab.Index, tab.Title, tab.URL)
	}

	// Access different tabs
	if len(tabs) > 1 {
		page1 := browser.Page() // First tab
		page2, err := browser.TabByIndex(1)
		if err != nil {
			t.Fatal(err)
		}

		// Navigate both tabs
		_, err = page1.Navigate("https://example.com", true)
		if err != nil {
			t.Fatal(err)
		}

		_, err = page2.Navigate("https://example.org", true)
		if err != nil {
			t.Fatal(err)
		}

		// Verify they're different
		title1, _ := page1.Eval("document.title")
		title2, _ := page2.Eval("document.title")

		if title1 == title2 {
			t.Error("Tabs should have different titles")
		}

		t.Logf("Tab 1 title: %v", title1)
		t.Logf("Tab 2 title: %v", title2)
	}
}

// TestNewTab demonstrates creating new tabs
func TestNewTab(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	initialCount := browser.TabCount()

	// Create a new tab with URL
	newPage, err := browser.NewTab("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	// Verify tab was created
	newCount := browser.TabCount()
	if newCount != initialCount+1 {
		t.Errorf("Expected %d tabs, got %d", initialCount+1, newCount)
	}

	// Verify the new tab works
	title, err := newPage.Eval("document.title")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("New tab title: %v", title)

	// Create blank tab
	blankPage, err := browser.NewTab("")
	if err != nil {
		t.Fatal(err)
	}

	// Verify blank tab count
	if browser.TabCount() != initialCount+2 {
		t.Errorf("Expected %d tabs after creating blank tab", initialCount+2)
	}

	// Navigate blank tab
	_, err = blankPage.Navigate("https://example.org", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Final tab count: %d", browser.TabCount())
}

// TestConcurrentTabOperations demonstrates concurrent operations on multiple tabs
func TestConcurrentTabOperations(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	// Create multiple tabs
	tab1, err := browser.NewTab("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	tab2, err := browser.NewTab("https://example.org")
	if err != nil {
		t.Fatal(err)
	}

	tab3, err := browser.NewTab("https://example.net")
	if err != nil {
		t.Fatal(err)
	}

	// Perform concurrent operations
	type result struct {
		tab   int
		title interface{}
		err   error
	}

	results := make(chan result, 3)

	// Execute JavaScript concurrently on all tabs
	go func() {
		title, err := tab1.Eval("document.title")
		results <- result{1, title, err}
	}()

	go func() {
		title, err := tab2.Eval("document.title")
		results <- result{2, title, err}
	}()

	go func() {
		title, err := tab3.Eval("document.title")
		results <- result{3, title, err}
	}()

	// Collect results
	for i := 0; i < 3; i++ {
		res := <-results
		if res.err != nil {
			t.Errorf("Tab %d error: %v", res.tab, res.err)
		} else {
			t.Logf("Tab %d title: %v", res.tab, res.title)
		}
	}

	t.Logf("Successfully performed concurrent operations on %d tabs", 3)
}

// TestTabByIndex demonstrates accessing tabs by index
func TestTabByIndex(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	// Ensure we have at least 2 tabs
	if browser.TabCount() < 2 {
		_, err := browser.NewTab("https://example.com")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test valid index
	page, err := browser.TabByIndex(0)
	if err != nil {
		t.Fatal(err)
	}

	if page == nil {
		t.Error("Expected page, got nil")
	}

	// Test invalid index
	_, err = browser.TabByIndex(999)
	if err == nil {
		t.Error("Expected error for invalid index, got nil")
	}

	// Test negative index
	_, err = browser.TabByIndex(-1)
	if err == nil {
		t.Error("Expected error for negative index, got nil")
	}

	t.Logf("TabByIndex validation working correctly")
}

// TestCloseTab demonstrates closing specific tabs
func TestCloseTab(t *testing.T) {
	browser, err := client.New(nil)
	if err != nil {
		t.Skip("Skipping test: Chrome not running")
	}
	defer browser.Close()

	initialCount := browser.TabCount()

	// Create a new tab
	_, err = browser.NewTab("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	if browser.TabCount() != initialCount+1 {
		t.Error("Tab was not created")
	}

	// Close the last tab
	err = browser.CloseTab(browser.TabCount() - 1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify tab was closed
	if browser.TabCount() != initialCount {
		t.Errorf("Expected %d tabs after closing, got %d", initialCount, browser.TabCount())
	}

	// Test closing invalid index
	err = browser.CloseTab(999)
	if err == nil {
		t.Error("Expected error when closing invalid index")
	}

	t.Logf("CloseTab working correctly")
}
