# brow

Simple CLI tools for browser automation via Chrome DevTools Protocol. Inspired by the philosophy of composable, low-overhead tools for AI agents.

## Quick Start

```bash
# Build
go build -o brow

# Start Chrome with remote debugging
./brow start --profile

# Navigate and interact
./brow nav https://example.com
./brow eval 'document.title'
./brow screenshot page.png
```

## Commands

### start
Launch Chrome with remote debugging on port 9222.
```bash
brow start              # Fresh session
brow start --profile    # Persistent profile (keeps cookies/logins)
brow start --headless   # Run headless
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

## Philosophy

- **Composable**: Each command is independent and outputs text/files
- **Low overhead**: ~200 token documentation vs 13K+ for MCP servers
- **Simple**: Just CLI tools, no complex protocols
- **Stateful**: Chrome instance maintains state between commands
- **Extensible**: Add new commands by creating new Go files in cmd/

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

- Single binary with subcommands
- Connects to Chrome via remote debugging (port 9222)
- Uses chromedp for Chrome DevTools Protocol communication
- State persists in Chrome instance, not in brow CLI
- Each command runs independently and exits

## Requirements

- Go 1.21+
- Google Chrome or Chromium installed

## Installation

```bash
git clone https://github.com/matejch/brow
cd brow
go build -o brow
```

Or install directly:
```bash
go install github.com/matejch/brow@latest
```
