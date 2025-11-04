package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/matej/brow/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	domain       string
	setCookie    string
	clearCookies bool
)

var cookiesCmd = &cobra.Command{
	Use:   "cookies",
	Short: "Get, set, or clear browser cookies",
	Long: `Manages browser cookies.
By default, retrieves all cookies as JSON.
Use --set to set a cookie (format: "name=value; domain=.example.com; path=/")
Use --clear to clear all cookies.`,
	RunE: runCookies,
}

func init() {
	rootCmd.AddCommand(cookiesCmd)
	cookiesCmd.Flags().StringVarP(&domain, "domain", "d", "", "Filter cookies by domain")
	cookiesCmd.Flags().StringVarP(&setCookie, "set", "s", "", "Set a cookie (format: name=value)")
	cookiesCmd.Flags().BoolVarP(&clearCookies, "clear", "c", false, "Clear all cookies")
}

func runCookies(_ *cobra.Command, _ []string) error {
	// Clear cookies
	if clearCookies {
		return clearAllCookies()
	}

	// Set cookie
	if setCookie != "" {
		return setACookie()
	}

	// Get cookies (default)
	return getCookies()
}

func getCookies() error {
	ctx, cancel, err := browser.GetExistingTabContext()
	if err != nil {
		return err
	}
	defer cancel()

	var cookies []*network.Cookie

	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		allCookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}

		// Filter by domain if specified
		if domain != "" {
			for _, cookie := range allCookies {
				if cookie.Domain == domain || "."+cookie.Domain == domain {
					cookies = append(cookies, cookie)
				}
			}
		} else {
			cookies = allCookies
		}

		return nil
	})); err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	// Format as JSON
	output, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format cookies: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func setACookie() error {
	// Simple cookie parsing (name=value format)
	// For more complex cookies, users can use JavaScript via eval
	fmt.Println("Setting cookie via JavaScript...")

	ctx, cancel, err := browser.GetExistingTabContext()
	if err != nil {
		return err
	}
	defer cancel()

	script := fmt.Sprintf("document.cookie = %q", setCookie)

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, nil)); err != nil {
		return fmt.Errorf("failed to set cookie: %w", err)
	}

	fmt.Println("Cookie set successfully")
	return nil
}

func clearAllCookies() error {
	ctx, cancel, err := browser.GetExistingTabContext()
	if err != nil {
		return err
	}
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return network.ClearBrowserCookies().Do(ctx)
	})); err != nil {
		return fmt.Errorf("failed to clear cookies: %w", err)
	}

	fmt.Println("All cookies cleared")
	return nil
}
