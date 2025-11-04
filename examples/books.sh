#!/bin/bash
# Example: Scrape book catalog from books.toscrape.com
# Another site specifically designed for web scraping practice
# Demonstrates more complex data extraction

echo "Starting Chrome..."
../brow start --profile &
sleep 3  # Give Chrome time to start

echo "Navigating to Books to Scrape..."
../brow nav https://books.toscrape.com

echo ""
echo "Getting book count..."
../brow eval 'document.querySelectorAll(".product_pod").length'

echo ""
echo "Extracting book catalog data to JSON..."
../brow eval '
Array.from(document.querySelectorAll(".product_pod")).map(book => {
  const title = book.querySelector("h3 a").getAttribute("title");
  const price = book.querySelector(".price_color").textContent;
  const availability = book.querySelector(".availability").textContent.trim();

  // Extract star rating from class name (e.g., "star-rating Three")
  const ratingElement = book.querySelector(".star-rating");
  const ratingClass = ratingElement ? ratingElement.className : "";
  const rating = ratingClass.replace("star-rating ", "");

  return {
    title: title,
    price: price,
    rating: rating,
    availability: availability
  };
})
' > books.json

echo ""
echo "Taking screenshot..."
../brow screenshot books.png

echo ""
echo "Exporting to PDF..."
../brow pdf books.pdf

echo ""
echo "Extracting just book titles..."
../brow eval '
Array.from(document.querySelectorAll(".product_pod h3 a"))
  .map(a => a.getAttribute("title"))
' > book-titles.json

echo ""
echo "Done! Check the following files:"
echo "  - books.json (full book data with prices, ratings, availability)"
echo "  - book-titles.json (just the titles)"
echo "  - books.png (screenshot)"
echo "  - books.pdf (PDF export)"
echo ""
echo "Sample data from books.json:"
head -30 books.json
