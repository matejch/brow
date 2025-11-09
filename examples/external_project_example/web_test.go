package main

import (
	"testing"
	"time"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
)

// Example: Testing a web application with brow
func TestWebApplication(t *testing.T) {
	// Connect to Chrome
	browser, err := client.New(&config.Config{
		Port:    9222,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Skip("Chrome not running. Start with: brow start")
	}
	defer browser.Close()

	page := browser.Page()

	t.Run("Homepage loads", func(t *testing.T) {
		result, err := page.Navigate("https://example.com", true)
		if err != nil {
			t.Fatal(err)
		}

		if result.Title != "Example Domain" {
			t.Errorf("Expected 'Example Domain', got '%s'", result.Title)
		}
	})

	t.Run("JavaScript works", func(t *testing.T) {
		// Check that page is interactive
		result, err := page.Eval("typeof document !== 'undefined'")
		if err != nil {
			t.Fatal(err)
		}

		if result != true {
			t.Error("JavaScript not working properly")
		}
	})

	t.Run("Can take screenshots", func(t *testing.T) {
		screenshot, err := page.Screenshot(operations.ScreenshotOptions{
			FullPage: false,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(screenshot) == 0 {
			t.Error("Screenshot is empty")
		}
	})
}

// Example: E2E test for a login flow
func TestLoginFlow(t *testing.T) {
	browser, err := client.New(nil) // Use defaults
	if err != nil {
		t.Skip("Chrome not running")
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate to login page
	_, err = page.Navigate("https://example.com", true)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate filling out a form (if the page had a login form)
	// This is just an example - adapt to your actual app
	_, err = page.Eval(`
		// Your form filling logic here
		console.log('Would fill out login form here');
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Check cookies were set after "login"
	cookies, err := page.GetCookies("")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Found %d cookies after navigation", len(cookies))
}
