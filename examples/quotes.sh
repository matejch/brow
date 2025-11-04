#!/bin/bash
# Example: Scrape quotes from quotes.toscrape.com
# This site is specifically designed for web scraping practice
# No rate limits, no ethical concerns - perfect for testing!

echo "Starting Chrome..."
../brow start --profile &
sleep 3  # Give Chrome time to start

echo "Navigating to Quotes to Scrape..."
../brow nav https://quotes.toscrape.com

echo ""
echo "Getting quote count..."
../brow eval 'document.querySelectorAll(".quote").length'

echo ""
echo "Extracting all quotes with authors and tags to JSON..."
../brow eval '
Array.from(document.querySelectorAll(".quote")).map(quote => ({
  text: quote.querySelector(".text").textContent.trim(),
  author: quote.querySelector(".author").textContent,
  tags: Array.from(quote.querySelectorAll(".tag")).map(tag => tag.textContent)
}))
' > quotes.json

echo ""
echo "Taking screenshot..."
../brow screenshot quotes.png

echo ""
echo "Exporting to PDF..."
../brow pdf quotes.pdf

echo ""
echo "Done! Check the following files:"
echo "  - quotes.json (structured quote data)"
echo "  - quotes.png (screenshot)"
echo "  - quotes.pdf (PDF export)"
echo ""
echo "Sample data from quotes.json:"
head -20 quotes.json
