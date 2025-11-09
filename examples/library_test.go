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
