package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// GetCookies retrieves all cookies, optionally filtered by domain
func GetCookies(ctx context.Context, domain string) ([]*network.Cookie, error) {
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
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}

	return cookies, nil
}

// SetCookie sets a cookie using the document.cookie API
// The cookie parameter should be in the format: "name=value; domain=.example.com; path=/"
func SetCookie(ctx context.Context, cookie string) error {
	// Safely escape the cookie value using JSON encoding
	cookieJSON, err := json.Marshal(cookie)
	if err != nil {
		return fmt.Errorf("failed to escape cookie value: %w", err)
	}

	script := fmt.Sprintf("document.cookie = %s", string(cookieJSON))

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, nil)); err != nil {
		return fmt.Errorf("failed to set cookie: %w", err)
	}

	return nil
}

// ClearCookies clears all browser cookies
func ClearCookies(ctx context.Context) error {
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return network.ClearBrowserCookies().Do(ctx)
	})); err != nil {
		return fmt.Errorf("failed to clear cookies: %w", err)
	}

	return nil
}
