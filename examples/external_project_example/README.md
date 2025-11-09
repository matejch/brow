# External Project Example

This demonstrates how to use brow as a library in your own Go project.

## Setup

### 1. Add brow to your project

```bash
# In your project directory
go get github.com/matejch/brow@latest
```

### 2. Start Chrome with remote debugging

```bash
# Option 1: Use brow CLI
brow start --headless

# Option 2: Start Chrome manually
google-chrome --remote-debugging-port=9222 --headless
```

### 3. Run the example

```bash
# Run the main program
go run main.go

# Or run tests
go test -v
```

## Usage in Your Project

### Basic Usage

```go
import (
    "github.com/matejch/brow/pkg/client"
    "github.com/matejch/brow/pkg/config"
)

// Connect to browser
browser, err := client.New(&config.Config{Port: 9222})
if err != nil {
    panic(err)
}
defer browser.Close()

// Use it
page := browser.Page()
page.Navigate("https://example.com", true)
title, _ := page.Eval("document.title")
```

### In Tests

```go
func TestMyWebApp(t *testing.T) {
    browser, err := client.New(nil)
    if err != nil {
        t.Skip("Chrome not running")
    }
    defer browser.Close()

    page := browser.Page()

    // Your test logic
    result, _ := page.Navigate("https://myapp.com", true)
    if result.Title != "Expected Title" {
        t.Error("Wrong title")
    }
}
```

## Real-World Use Cases

### 1. E2E Testing
```go
// Test login flow
page.Navigate("https://app.com/login", true)
page.Eval(`document.querySelector('#username').value = 'user'`)
page.Eval(`document.querySelector('#password').value = 'pass'`)
page.Eval(`document.querySelector('form').submit()`)
// Verify logged in...
```

### 2. Web Scraping
```go
page.Navigate("https://quotes.toscrape.com", true)
quotes, _ := page.Eval(`
    Array.from(document.querySelectorAll('.quote')).map(q => ({
        text: q.querySelector('.text').textContent,
        author: q.querySelector('.author').textContent
    }))
`)
```

### 3. Screenshot Service
```go
page.Navigate(userUrl, true)
screenshot, _ := page.Screenshot(operations.ScreenshotOptions{
    FullPage: true,
})
// Serve screenshot to user...
```

### 4. PDF Generation
```go
page.Navigate("https://report.com", true)
pdf, _ := page.PDF(operations.PDFOptions{
    PrintBackground: true,
})
// Save or send PDF...
```

## Configuration Options

```go
browser, _ := client.New(&config.Config{
    Port:    9223,              // Custom port
    Timeout: 60 * time.Second,  // Operation timeout
})
```

## Multiple Instances

You can run multiple browser instances concurrently:

```go
// Browser 1 on port 9222
browser1, _ := client.New(&config.Config{Port: 9222})
defer browser1.Close()

// Browser 2 on port 9223
browser2, _ := client.New(&config.Config{Port: 9223})
defer browser2.Close()

// Use them independently
browser1.Page().Navigate("https://site1.com", true)
browser2.Page().Navigate("https://site2.com", true)
```
