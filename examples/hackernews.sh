#!/bin/bash
# Example: Scrape Hacker News front page stories
# This demonstrates the composable, low-overhead philosophy

echo "Starting Chrome..."
./brow start --profile &
sleep 3  # Give Chrome time to start

echo "Navigating to Hacker News..."
./brow nav https://news.ycombinator.com

echo "Getting story count..."
./brow eval 'document.querySelectorAll(".titleline > a").length'

echo "Extracting all stories to JSON..."
./brow eval 'Array.from(document.querySelectorAll(".athing")).map(story => ({
  rank: story.querySelector(".rank")?.textContent,
  title: story.querySelector(".titleline > a")?.textContent,
  url: story.querySelector(".titleline > a")?.href,
  score: story.nextElementSibling?.querySelector(".score")?.textContent,
  user: story.nextElementSibling?.querySelector(".hnuser")?.textContent
}))' > stories.json

echo "Taking screenshot..."
./brow screenshot hackernews.png

echo "Exporting to PDF..."
./brow pdf hackernews.pdf

echo "Done! Check stories.json, hackernews.png, and hackernews.pdf"
