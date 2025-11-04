package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/matej/brow/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	storageType string
	key         string
	value       string
	deleteKey   bool
	clearStorage bool
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Interact with localStorage and sessionStorage",
	Long: `Get, set, or clear browser storage (localStorage or sessionStorage).
By default, retrieves all items from localStorage as JSON.
Use --type to specify localStorage (default) or sessionStorage.`,
	RunE: runStorage,
}

func init() {
	rootCmd.AddCommand(storageCmd)
	storageCmd.Flags().StringVarP(&storageType, "type", "t", "local", "Storage type: local or session")
	storageCmd.Flags().StringVarP(&key, "key", "k", "", "Storage key to get or set")
	storageCmd.Flags().StringVarP(&value, "value", "v", "", "Value to set (requires --key)")
	storageCmd.Flags().BoolVarP(&deleteKey, "delete", "d", false, "Delete the specified key")
	storageCmd.Flags().BoolVarP(&clearStorage, "clear", "c", false, "Clear all storage")
}

func runStorage(cmd *cobra.Command, args []string) error {
	// Determine storage type
	var storageName string
	if storageType == "session" {
		storageName = "sessionStorage"
	} else {
		storageName = "localStorage"
	}

	// Clear storage
	if clearStorage {
		script := fmt.Sprintf("%s.clear()", storageName)
		if err := browser.Run(chromedp.Evaluate(script, nil)); err != nil {
			return fmt.Errorf("failed to clear storage: %w", err)
		}
		fmt.Printf("%s cleared\n", storageName)
		return nil
	}

	// Delete key
	if deleteKey && key != "" {
		script := fmt.Sprintf("%s.removeItem(%q)", storageName, key)
		if err := browser.Run(chromedp.Evaluate(script, nil)); err != nil {
			return fmt.Errorf("failed to delete key: %w", err)
		}
		fmt.Printf("Deleted key: %s\n", key)
		return nil
	}

	// Set value
	if key != "" && value != "" {
		script := fmt.Sprintf("%s.setItem(%q, %q)", storageName, key, value)
		if err := browser.Run(chromedp.Evaluate(script, nil)); err != nil {
			return fmt.Errorf("failed to set value: %w", err)
		}
		fmt.Printf("Set %s[%s] = %s\n", storageName, key, value)
		return nil
	}

	// Get specific key
	if key != "" {
		script := fmt.Sprintf("%s.getItem(%q)", storageName, key)
		var result interface{}
		if err := browser.Run(chromedp.Evaluate(script, &result)); err != nil {
			return fmt.Errorf("failed to get value: %w", err)
		}
		fmt.Printf("%v\n", result)
		return nil
	}

	// Get all items (default)
	script := fmt.Sprintf(`
		(() => {
			let items = {};
			for (let i = 0; i < %s.length; i++) {
				let key = %s.key(i);
				items[key] = %s.getItem(key);
			}
			return items;
		})()
	`, storageName, storageName, storageName)

	var result interface{}
	if err := browser.Run(chromedp.Evaluate(script, &result)); err != nil {
		return fmt.Errorf("failed to get storage items: %w", err)
	}

	// Format as JSON
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format storage: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
