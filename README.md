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

## Example Agent Workflow

```bash
# Start browser
brow start --profile

# Navigate and scrape
brow nav https://news.ycombinator.com
brow eval 'document.querySelectorAll(".titleline > a").length'

# Get all story titles
brow eval 'Array.from(document.querySelectorAll(".titleline > a")).map(a => ({title: a.textContent, url: a.href}))' > stories.json

# Screenshot
brow screenshot hn.png

# Export PDF
brow pdf hn.pdf
```

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
git clone https://github.com/matej/brow
cd brow
go build -o brow
```

Or install directly:
```bash
go install github.com/matej/brow@latest
```
