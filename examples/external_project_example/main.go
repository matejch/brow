package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
)

func main() {
	fmt.Println("Example: Using brow as a library")
	fmt.Println("==================================")

	// Connect to Chrome (make sure it's running with: brow start)
	browser, err := client.New(&config.Config{
		Port:    9222,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to browser: %v\n", err)
	}
	defer browser.Close()

	page := browser.Page()

	// Navigate to a website
	fmt.Println("\n1. Navigating to example.com...")
	result, err := page.Navigate("https://example.com", true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Page title: %s\n", result.Title)

	// Execute JavaScript
	fmt.Println("\n2. Counting links on the page...")
	linkCount, err := page.Eval("document.querySelectorAll('a').length")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Found %v links\n", linkCount)

	// Capture screenshot
	fmt.Println("\n3. Taking screenshot...")
	screenshot, err := page.Screenshot(operations.ScreenshotOptions{
		FullPage: true,
		Quality:  90,
	})
	if err != nil {
		log.Fatal(err)
	}

	filename := "example_screenshot.png"
	if err := os.WriteFile(filename, screenshot, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Saved to %s (%d bytes)\n", filename, len(screenshot))

	// Manage cookies
	fmt.Println("\n4. Setting a cookie...")
	err = page.SetCookie("my_cookie=hello_world; path=/")
	if err != nil {
		log.Fatal(err)
	}

	cookies, err := page.GetCookies("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Total cookies: %d\n", len(cookies))

	// Work with localStorage
	fmt.Println("\n5. Using localStorage...")
	err = page.SetStorageItem(operations.LocalStorage, "user_id", "12345")
	if err != nil {
		log.Fatal(err)
	}

	value, err := page.GetStorageItem(operations.LocalStorage, "user_id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Retrieved value: %v\n", value)

	// Generate PDF
	fmt.Println("\n6. Generating PDF...")
	pdf, err := page.PDF(operations.PDFOptions{
		Landscape:       false,
		PrintBackground: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	pdfFilename := "example_page.pdf"
	if err := os.WriteFile(pdfFilename, pdf, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Saved to %s (%d bytes)\n", pdfFilename, len(pdf))

	fmt.Println("\n✅ All operations completed successfully!")
}
