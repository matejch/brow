package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
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
	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	cookies, err := browser.Page().GetCookies(domain)
	if err != nil {
		return err
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
	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	if err := browser.Page().SetCookie(setCookie); err != nil {
		return err
	}

	fmt.Println("Cookie set successfully")
	return nil
}

func clearAllCookies() error {
	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	if err := browser.Page().ClearCookies(); err != nil {
		return err
	}

	fmt.Println("All cookies cleared")
	return nil
}
