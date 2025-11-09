# brow

Simple CLI tools for browser automation via Chrome DevTools Protocol. Inspired by the philosophy of composable, low-overhead tools for AI agents.

**Now available as a Go library!** Use brow programmatically in your applications.

## Installation

### As a Library

```bash
go get github.com/matejch/brow@latest
```

Then import in your code:
```go
import "github.com/matejch/brow/pkg/client"
```

### As a CLI Tool

```bash
# Install from source
git clone https://github.com/matejch/brow
cd brow
go build -o brow

# Or install directly
go install github.com/matejch/brow@latest
```

## Quick Start

### CLI Usage

```bash
# Start Chrome with remote debugging
./brow start --profile

# Navigate and interact
./brow nav https://example.com
./brow eval 'document.title'
./brow screenshot page.png
```

### Library Usage

```go
package main

import (
    "fmt"
    "github.com/matejch/brow/pkg/client"
)

func main() {
    browser, _ := client.New(nil)
    defer browser.Close()

    page := browser.Page()
    result, _ := page.Navigate("https://example.com", true)
    fmt.Println(result.Title) // "Example Domain"
}
```

## Commands

### start
Launch Chrome with remote debugging (default port 9222).
```bash
brow start              # Fresh session
brow start --profile    # Persistent profile (keeps cookies/logins)
brow start --headless   # Run headless
brow start --port 9223  # Use custom port
```

### nav
Navigate to a URL.
```bash
brow nav https://example.com
```

### eval
Execute JavaScript in the current page.
```bash
brow eval 'document.querySelectorAll("a").length'
brow eval 'document.body.innerText' --raw
```

### screenshot
Capture a screenshot.
```bash
brow screenshot output.png
brow screenshot --full-page  # Capture entire page
brow screenshot --base64     # Output base64 data
```

### pick
Interactive element picker to get CSS selectors.
```bash
brow pick              # Returns CSS selector
brow pick --xpath      # Returns XPath
# After clicking an element:
brow eval 'window.__browPickedSelector'
```

### cookies
Manage cookies.
```bash
brow cookies                        # Get all cookies as JSON
brow cookies --domain example.com   # Filter by domain
brow cookies --set "name=value"     # Set cookie
brow cookies --clear                # Clear all cookies
```

### storage
Interact with localStorage/sessionStorage.
```bash
brow storage                              # Get all localStorage
brow storage --type session               # Get sessionStorage
brow storage --key name --value data      # Set item
brow storage --key name                   # Get specific item
brow storage --key name --delete          # Delete item
brow storage --clear                      # Clear all
```

### pdf
Export page as PDF.
```bash
brow pdf output.pdf
brow pdf --landscape
brow pdf --no-background
```

## Port Configuration

By default, brow connects to Chrome on port 9222. You can customize the port in three ways:

1. **Command-line flag** (highest precedence):
```bash
brow --port 9223 start
brow --port 9223 nav https://example.com
```

2. **Environment variable**:
```bash
export BROW_DEBUG_PORT=9223
brow start
brow nav https://example.com
```

3. **Default fallback**: Port 9222 (if neither flag nor env var is set)

**Use cases for custom ports:**
- Run multiple Chrome instances simultaneously
- Avoid port conflicts with other applications
- Use in CI/CD environments with dynamic port allocation

**Example with multiple instances:**
```bash
# Terminal 1: Chrome on port 9222
brow start --profile

# Terminal 2: Chrome on port 9223
brow --port 9223 start --profile

# Now you can control both independently
brow nav https://example.com                    # Controls port 9222
brow --port 9223 nav https://different-site.com # Controls port 9223
```

## Library Usage

### Installation

Add brow to your project:

```bash
go get github.com/matejch/brow@v1.2.0
```

### Complete Example

brow provides a clean, type-safe API for browser automation:

```go
package main

import (
    "fmt"
    "time"

    "github.com/matejch/brow/pkg/client"
    "github.com/matejch/brow/pkg/config"
    "github.com/matejch/brow/pkg/operations"
)

func main() {
    // Create browser instance (connects to existing Chrome)
    browser, err := client.New(&config.Config{
        Port:    9222,
        Timeout: 30 * time.Second,
    })
    if err != nil {
        panic(err)
    }
    defer browser.Close()

    page := browser.Page()

    // Navigate to a page
    result, _ := page.Navigate("https://example.com", true)
    fmt.Printf("Title: %s\n", result.Title)

    // Execute JavaScript
    linkCount, _ := page.Eval("document.querySelectorAll('a').length")
    fmt.Printf("Links: %v\n", linkCount)

    // Capture screenshot
    screenshot, _ := page.Screenshot(operations.ScreenshotOptions{
        FullPage: true,
    })
    // ... save screenshot to file

    // Manage cookies
    cookies, _ := page.GetCookies("")
    fmt.Printf("Found %d cookies\n", len(cookies))

    // Work with localStorage
    page.SetStorageItem(operations.LocalStorage, "key", "value")
    value, _ := page.GetStorageItem(operations.LocalStorage, "key")

    // Generate PDF
    pdf, _ := page.PDF(operations.PDFOptions{
        Landscape: false,
        PrintBackground: true,
    })
    // ... save PDF to file
}
```

### Library Features

- **Simple API**: Browser and Page objects with methods for all operations
- **No global state**: Create multiple browser instances with different configurations
- **Type-safe**: Strongly-typed options for screenshots, PDFs, storage, etc.
- **Resource management**: Proper cleanup with defer browser.Close()
- **Timeout support**: Configure operation timeouts
- **Input sanitization**: Built-in protection against JavaScript injection
- **Well-tested**: Comprehensive test suite and working examples
- **Zero breaking changes**: CLI remains fully compatible

### Documentation

- **[Library Usage Guide](LIBRARY_USAGE_GUIDE.md)** - Complete API reference and examples
- **[Refactor Summary](REFACTOR_SUMMARY.md)** - Architecture and implementation details
- **[Example Tests](examples/library_test.go)** - Comprehensive test suite
- **[Example Project](examples/external_project_example/)** - Standalone working example

## Philosophy

- **Composable**: Each command is independent and outputs text/files
- **Low overhead**: ~200 token documentation vs 13K+ for MCP servers
- **Simple**: Just CLI tools, no complex protocols
- **Stateful**: Chrome instance maintains state between commands
- **Extensible**: Add new commands by creating new Go files in cmd/
- **Library-first**: Core functionality available as clean Go API

## Example Workflows

### Scrape Quotes (using quotes.toscrape.com)
```bash
# Start browser
brow start --profile

# Navigate to practice scraping site
brow nav https://quotes.toscrape.com

# Count quotes on page
brow eval 'document.querySelectorAll(".quote").length'

# Extract all quotes with authors and tags
brow eval 'Array.from(document.querySelectorAll(".quote")).map(quote => ({
  text: quote.querySelector(".text").textContent.trim(),
  author: quote.querySelector(".author").textContent,
  tags: Array.from(quote.querySelectorAll(".tag")).map(tag => tag.textContent)
}))' > quotes.json

# Screenshot and PDF
brow screenshot quotes.png
brow pdf quotes.pdf
```

### Scrape Book Catalog (using books.toscrape.com)
```bash
# Navigate to book catalog
brow nav https://books.toscrape.com

# Extract book data
brow eval 'Array.from(document.querySelectorAll(".product_pod")).map(book => ({
  title: book.querySelector("h3 a").getAttribute("title"),
  price: book.querySelector(".price_color").textContent,
  availability: book.querySelector(".availability").textContent.trim()
}))' > books.json

# Capture catalog
brow screenshot books.png
```

**Note**: These examples use sites specifically designed for web scraping practice. See `examples/` directory for complete scripts.

## Running Examples

Complete example scripts are available in the `examples/` directory.

### Run from examples directory:
```bash
cd examples
./quotes.sh    # Scrape quotes.toscrape.com
./books.sh     # Scrape books.toscrape.com
```

### What the examples produce:
- **quotes.sh**: `quotes.json`, `quotes.png`, `quotes.pdf`
- **books.sh**: `books.json`, `book-titles.json`, `books.png`, `books.pdf`

**Note:** The example scripts must be run from the `examples/` directory because they reference `../brow` to access the binary in the project root.

## Architecture

### CLI Architecture
- Single binary with subcommands
- Connects to Chrome via remote debugging (default port 9222)
- Uses chromedp for Chrome DevTools Protocol communication
- State persists in Chrome instance, not in brow CLI
- Each command runs independently and exits

### Library Architecture
```
┌─────────────────────────────────┐
│   Your Application              │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│   pkg/client (Browser, Page)    │  ← Public API
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│   pkg/operations                │  ← Business logic
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│   chromedp (CDP)                │  ← Chrome DevTools Protocol
└─────────────────────────────────┘
```

## Requirements

- **Go 1.21+** (for building or using as library)
- **Google Chrome or Chromium** installed
- **Chrome remote debugging** enabled (via `brow start` or manual launch)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

For major changes, please open an issue first to discuss what you would like to change.

## License

See LICENSE file for details.

## Links

- **GitHub**: https://github.com/matejch/brow
- **Issues**: https://github.com/matejch/brow/issues
- **Latest Release**: https://github.com/matejch/brow/releases/latest
