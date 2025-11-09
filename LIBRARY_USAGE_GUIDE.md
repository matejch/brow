# Library Usage Guide

Complete guide for using brow as a Go library in your projects.

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Complete API Reference](#complete-api-reference)
4. [Common Use Cases](#common-use-cases)
5. [Best Practices](#best-practices)
6. [Troubleshooting](#troubleshooting)

---

## Installation

### Prerequisites

- Go 1.21 or later
- Google Chrome or Chromium installed
- Chrome running with remote debugging enabled

### Step 1: Add brow to your project

```bash
# In your project directory
go get github.com/matejch/brow@latest
```

This adds brow to your `go.mod`:
```go
module myproject

go 1.21

require github.com/matejch/brow v0.1.0
```

### Step 2: Start Chrome with remote debugging

```bash
# Option 1: Use brow CLI
brow start --headless

# Option 2: Start Chrome manually
google-chrome --remote-debugging-port=9222 --headless

# Option 3: With a persistent profile
brow start --profile --headless
```

---

## Quick Start

### Minimal Example

```go
package main

import (
    "fmt"
    "github.com/matejch/brow/pkg/client"
)

func main() {
    // Connect to Chrome
    browser, err := client.New(nil) // Uses default config
    if err != nil {
        panic(err)
    }
    defer browser.Close()

    // Navigate and get title
    page := browser.Page()
    result, _ := page.Navigate("https://example.com", true)
    fmt.Println(result.Title) // "Example Domain"
}
```

### With Configuration

```go
import (
    "time"
    "github.com/matejch/brow/pkg/client"
    "github.com/matejch/brow/pkg/config"
)

browser, err := client.New(&config.Config{
    Port:    9222,              // Chrome debugging port
    Timeout: 30 * time.Second,  // Operation timeout
})
if err != nil {
    panic(err)
}
defer browser.Close()
```

---

## Complete API Reference

### Configuration

```go
type Config struct {
    Port    int           // Chrome DevTools port (default: 9222)
    Timeout time.Duration // Operation timeout (default: 30s)
}

// Create with defaults
config.Default()

// Resolve port (flag > env var > default)
config.ResolvePort(flagPort int) int
```

### Browser

```go
// Create new browser instance
browser, err := client.New(cfg *config.Config) (*Browser, error)

// Get page interface
page := browser.Page() *Page

// Get underlying context (advanced)
ctx := browser.Context() context.Context

// Set timeout dynamically
browser.SetTimeout(timeout time.Duration)

// Close and cleanup
browser.Close() error
```

### Page - Navigation

```go
// Navigate to URL
result, err := page.Navigate(url string, waitReady bool) (*NavigationResult, error)

// NavigationResult contains:
type NavigationResult struct {
    URL   string
    Title string
}
```

### Page - JavaScript

```go
// Execute JavaScript and get result
result, err := page.Eval(script string) (interface{}, error)

// Examples:
title, _ := page.Eval("document.title")
linkCount, _ := page.Eval("document.querySelectorAll('a').length")
data, _ := page.Eval(`({title: document.title, url: location.href})`)
```

### Page - Screenshots

```go
// Capture screenshot
screenshot, err := page.Screenshot(opts operations.ScreenshotOptions) ([]byte, error)

// Options:
type ScreenshotOptions struct {
    FullPage bool // Capture entire page vs viewport
    Quality  int  // JPEG quality (0-100, default 100)
}

// Examples:
viewport, _ := page.Screenshot(operations.ScreenshotOptions{})
fullPage, _ := page.Screenshot(operations.ScreenshotOptions{
    FullPage: true,
    Quality:  90,
})
```

### Page - PDF

```go
// Generate PDF
pdf, err := page.PDF(opts operations.PDFOptions) ([]byte, error)

// Options:
type PDFOptions struct {
    Landscape       bool // Landscape vs portrait
    PrintBackground bool // Include background graphics
}

// Example:
pdf, _ := page.PDF(operations.PDFOptions{
    Landscape:       false,
    PrintBackground: true,
})
```

### Page - Cookies

```go
// Get all cookies (or filtered by domain)
cookies, err := page.GetCookies(domain string) ([]*network.Cookie, error)

// Set a cookie
err := page.SetCookie(cookie string) error

// Clear all cookies
err := page.ClearCookies() error

// Examples:
allCookies, _ := page.GetCookies("")
exampleCookies, _ := page.GetCookies("example.com")
page.SetCookie("session=abc123; path=/; domain=.example.com")
page.ClearCookies()
```

### Page - Storage (localStorage/sessionStorage)

```go
// Storage types
const (
    LocalStorage   StorageType = "localStorage"
    SessionStorage StorageType = "sessionStorage"
)

// Get all items
items, err := page.GetAllStorage(storageType) (map[string]interface{}, error)

// Get specific item
value, err := page.GetStorageItem(storageType, key string) (interface{}, error)

// Set item
err := page.SetStorageItem(storageType, key, value string) error

// Remove item
err := page.RemoveStorageItem(storageType, key string) error

// Clear all
err := page.ClearStorage(storageType) error

// Examples:
page.SetStorageItem(operations.LocalStorage, "user_id", "12345")
userId, _ := page.GetStorageItem(operations.LocalStorage, "user_id")
all, _ := page.GetAllStorage(operations.LocalStorage)
page.RemoveStorageItem(operations.LocalStorage, "user_id")
page.ClearStorage(operations.SessionStorage)
```

### Page - Element Picker

```go
// Inject interactive element picker
err := page.InjectPicker(useXPath bool) error

// Get picked selector
selector, err := page.GetPickedSelector() (string, error)

// Example:
page.InjectPicker(false) // Use CSS selectors
// User clicks element in browser...
selector, _ := page.GetPickedSelector()
fmt.Println(selector) // "#main > div.content"
```

### Page - Advanced

```go
// Get underlying context for custom chromedp operations
ctx := page.Context() context.Context
```

### Multi-Tab Operations

```go
// Get tab count
count := browser.TabCount() int

// List all tabs with metadata
tabs, err := browser.Tabs() ([]*TabInfo, error)
// TabInfo contains: Index, TargetID, Title, URL

// Access specific tab by index
page, err := browser.TabByIndex(index int) (*Page, error)

// Create new tab
newPage, err := browser.NewTab(url string) (*Page, error)

// Close specific tab
err := browser.CloseTab(index int) error

// Examples:
browser, _ := client.New(nil)

// List all tabs
tabs, _ := browser.Tabs()
for _, tab := range tabs {
    fmt.Printf("Tab %d: %s (%s)\n", tab.Index, tab.Title, tab.URL)
}

// Access second tab
page2, _ := browser.TabByIndex(1)
page2.Navigate("https://example.com", true)

// Create new tab with URL
newTab, _ := browser.NewTab("https://example.org")

// Create blank tab
blankTab, _ := browser.NewTab("")
blankTab.Navigate("https://example.net", true)

// Close tab by index
browser.CloseTab(2)
```

---

## Common Use Cases

### 1. End-to-End Testing

```go
func TestLoginFlow(t *testing.T) {
    browser, err := client.New(nil)
    if err != nil {
        t.Skip("Chrome not running")
    }
    defer browser.Close()

    page := browser.Page()

    // Navigate to login page
    _, err = page.Navigate("https://myapp.com/login", true)
    if err != nil {
        t.Fatal(err)
    }

    // Fill out form
    _, err = page.Eval(`
        document.querySelector('#username').value = 'testuser';
        document.querySelector('#password').value = 'testpass';
        document.querySelector('#login-form').submit();
    `)
    if err != nil {
        t.Fatal(err)
    }

    // Wait for redirect
    time.Sleep(2 * time.Second)

    // Verify logged in
    url, _ := page.Eval("window.location.href")
    if !strings.Contains(url.(string), "/dashboard") {
        t.Error("Login failed: not redirected to dashboard")
    }

    // Check for auth cookie
    cookies, _ := page.GetCookies("")
    hasAuthCookie := false
    for _, cookie := range cookies {
        if cookie.Name == "auth_token" {
            hasAuthCookie = true
            break
        }
    }
    if !hasAuthCookie {
        t.Error("No auth cookie found")
    }
}
```

### 2. Web Scraping

```go
func scrapeQuotes() ([]Quote, error) {
    browser, err := client.New(nil)
    if err != nil {
        return nil, err
    }
    defer browser.Close()

    page := browser.Page()

    // Navigate to quotes site
    _, err = page.Navigate("https://quotes.toscrape.com", true)
    if err != nil {
        return nil, err
    }

    // Extract quotes using JavaScript
    result, err := page.Eval(`
        Array.from(document.querySelectorAll('.quote')).map(quote => ({
            text: quote.querySelector('.text').textContent.trim(),
            author: quote.querySelector('.author').textContent,
            tags: Array.from(quote.querySelectorAll('.tag')).map(t => t.textContent)
        }))
    `)
    if err != nil {
        return nil, err
    }

    // Convert to Go structs
    // ... (parse result into []Quote)

    return quotes, nil
}
```

### 3. Screenshot Service

```go
func generateScreenshot(url string) ([]byte, error) {
    browser, err := client.New(&config.Config{
        Timeout: 60 * time.Second,
    })
    if err != nil {
        return nil, err
    }
    defer browser.Close()

    page := browser.Page()

    // Navigate
    _, err = page.Navigate(url, true)
    if err != nil {
        return nil, err
    }

    // Capture full-page screenshot
    return page.Screenshot(operations.ScreenshotOptions{
        FullPage: true,
        Quality:  90,
    })
}

// HTTP handler
func screenshotHandler(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    screenshot, err := generateScreenshot(url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "image/png")
    w.Write(screenshot)
}
```

### 4. PDF Report Generator

```go
func generateReport(reportURL string) ([]byte, error) {
    browser, err := client.New(nil)
    if err != nil {
        return nil, err
    }
    defer browser.Close()

    page := browser.Page()

    // Navigate to report page
    _, err = page.Navigate(reportURL, true)
    if err != nil {
        return nil, err
    }

    // Wait for dynamic content to load
    time.Sleep(2 * time.Second)

    // Generate PDF
    return page.PDF(operations.PDFOptions{
        Landscape:       true,
        PrintBackground: true,
    })
}
```

### 5. Automated Form Filling

```go
func fillAndSubmitForm(formData map[string]string) error {
    browser, err := client.New(nil)
    if err != nil {
        return err
    }
    defer browser.Close()

    page := browser.Page()

    // Navigate
    _, err = page.Navigate("https://forms.example.com", true)
    if err != nil {
        return err
    }

    // Fill each field
    for selector, value := range formData {
        script := fmt.Sprintf(
            "document.querySelector(%s).value = %s",
            jsonEncode(selector),
            jsonEncode(value),
        )
        _, err = page.Eval(script)
        if err != nil {
            return err
        }
    }

    // Submit
    _, err = page.Eval("document.querySelector('form').submit()")
    return err
}
```

### 6. Performance Monitoring

```go
func measurePageLoad(url string) (time.Duration, error) {
    browser, err := client.New(nil)
    if err != nil {
        return 0, err
    }
    defer browser.Close()

    page := browser.Page()

    start := time.Now()
    _, err = page.Navigate(url, true)
    if err != nil {
        return 0, err
    }

    loadTime := time.Since(start)

    // Get performance metrics
    perfData, _ := page.Eval(`
        ({
            loadTime: performance.timing.loadEventEnd - performance.timing.navigationStart,
            domReady: performance.timing.domContentLoadedEventEnd - performance.timing.navigationStart,
            firstPaint: performance.getEntriesByType('paint')[0]?.startTime || 0
        })
    `)

    fmt.Printf("Load time: %v\n", loadTime)
    fmt.Printf("Performance data: %+v\n", perfData)

    return loadTime, nil
}
```

### 7. Concurrent Multi-Tab Operations

```go
func scrapeMultipleSites(urls []string) ([]SiteData, error) {
    browser, err := client.New(nil)
    if err != nil {
        return nil, err
    }
    defer browser.Close()

    // Create tabs for each URL
    tabs := make([]*client.Page, len(urls))
    for i, url := range urls {
        tab, err := browser.NewTab(url)
        if err != nil {
            return nil, fmt.Errorf("failed to create tab %d: %w", i, err)
        }
        tabs[i] = tab
    }

    // Scrape concurrently
    type result struct {
        index int
        data  SiteData
        err   error
    }

    results := make(chan result, len(urls))

    for i, tab := range tabs {
        go func(index int, page *client.Page) {
            // Extract data from this tab
            title, err := page.Eval("document.title")
            if err != nil {
                results <- result{index, SiteData{}, err}
                return
            }

            headings, err := page.Eval(`
                Array.from(document.querySelectorAll('h1, h2')).map(h => h.textContent)
            `)
            if err != nil {
                results <- result{index, SiteData{}, err}
                return
            }

            results <- result{
                index: index,
                data: SiteData{
                    URL:      urls[index],
                    Title:    title.(string),
                    Headings: headings,
                },
            }
        }(i, tab)
    }

    // Collect results
    data := make([]SiteData, len(urls))
    for i := 0; i < len(urls); i++ {
        res := <-results
        if res.err != nil {
            return nil, fmt.Errorf("tab %d error: %w", res.index, res.err)
        }
        data[res.index] = res.data
    }

    return data, nil
}
```

---

## Best Practices

### 1. Always defer Close()

```go
browser, err := client.New(nil)
if err != nil {
    return err
}
defer browser.Close() // IMPORTANT: Always cleanup
```

### 2. Handle errors properly

```go
result, err := page.Navigate(url, true)
if err != nil {
    return fmt.Errorf("navigation failed: %w", err)
}
```

### 3. Use timeouts

```go
browser, _ := client.New(&config.Config{
    Timeout: 30 * time.Second, // Prevent hanging
})
```

### 4. Sanitize user input

The library automatically sanitizes JavaScript for cookies and storage, but be careful with eval:

```go
// UNSAFE: Direct interpolation
userInput := "alert('xss')"
page.Eval(userInput) // BAD!

// SAFE: Use proper encoding
import "encoding/json"
safe, _ := json.Marshal(userInput)
page.Eval(fmt.Sprintf("console.log(%s)", safe)) // GOOD
```

### 5. Wait for page readiness

```go
// Wait for page to be ready
page.Navigate(url, true) // true = wait for body

// Or wait for specific elements with JavaScript
page.Eval(`
    new Promise(resolve => {
        const interval = setInterval(() => {
            if (document.querySelector('.loaded')) {
                clearInterval(interval);
                resolve(true);
            }
        }, 100);
    })
`)
```

### 6. Reuse browser instances

```go
// INEFFICIENT: Creating new browser for each operation
for _, url := range urls {
    browser, _ := client.New(nil)
    browser.Page().Navigate(url, true)
    browser.Close()
}

// EFFICIENT: Reuse browser
browser, _ := client.New(nil)
defer browser.Close()
page := browser.Page()
for _, url := range urls {
    page.Navigate(url, true)
    // ... do work
}
```

---

## Troubleshooting

### "failed to connect to browser"

**Problem:** Chrome is not running or not on the expected port.

**Solution:**
```bash
# Start Chrome
brow start --headless

# Or specify custom port
brow --port 9223 start
browser, _ := client.New(&config.Config{Port: 9223})
```

### "no tabs available"

**Problem:** Chrome has no open tabs.

**Solution:** Chrome automatically opens a tab. If you see this, restart Chrome:
```bash
pkill chrome
brow start
```

### Operations timing out

**Problem:** Pages take too long to load.

**Solution:** Increase timeout:
```go
browser, _ := client.New(&config.Config{
    Timeout: 120 * time.Second, // Longer timeout
})
```

### JavaScript eval returns nil/unexpected results

**Problem:** JavaScript hasn't finished executing or DOM isn't ready.

**Solution:** Add waits:
```go
// Wait for element to exist
page.Eval(`
    new Promise(resolve => {
        const check = () => {
            const el = document.querySelector('.dynamic-content');
            if (el) resolve(el.textContent);
            else setTimeout(check, 100);
        };
        check();
    })
`)
```

### Resource leaks / too many connections

**Problem:** Not closing browsers.

**Solution:** Always use defer:
```go
browser, err := client.New(nil)
if err != nil {
    return err
}
defer browser.Close() // ← Essential!
```

### Can't run multiple instances

**Problem:** Trying to use same port.

**Solution:** Use different ports:
```bash
# Terminal 1
brow start --port 9222

# Terminal 2
brow start --port 9223
```

```go
browser1, _ := client.New(&config.Config{Port: 9222})
browser2, _ := client.New(&config.Config{Port: 9223})
```

---

## Examples

See complete working examples in:
- `examples/library_test.go` - Comprehensive test suite
- `examples/external_project_example/` - Standalone project example
- `examples/books.sh` - CLI scraping example
- `examples/quotes.sh` - CLI scraping example

---

## Support

- GitHub Issues: https://github.com/matejch/brow/issues
- Documentation: README.md
- Examples: examples/
