package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/matejch/brow/pkg/client"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
	"github.com/spf13/cobra"
)

var (
	storageType  string
	key          string
	value        string
	deleteKey    bool
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

func runStorage(_ *cobra.Command, _ []string) error {
	browser, err := client.New(&config.Config{
		Port: config.ResolvePort(Port),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	page := browser.Page()

	// Determine storage type
	var st operations.StorageType
	var storageName string
	if storageType == "session" {
		st = operations.SessionStorage
		storageName = "sessionStorage"
	} else {
		st = operations.LocalStorage
		storageName = "localStorage"
	}

	// Clear storage
	if clearStorage {
		if err := page.ClearStorage(st); err != nil {
			return err
		}
		fmt.Printf("%s cleared\n", storageName)
		return nil
	}

	// Delete key
	if deleteKey && key != "" {
		if err := page.RemoveStorageItem(st, key); err != nil {
			return err
		}
		fmt.Printf("Deleted key: %s\n", key)
		return nil
	}

	// Set value
	if key != "" && value != "" {
		if err := page.SetStorageItem(st, key, value); err != nil {
			return err
		}
		fmt.Printf("Set %s[%s] = %s\n", storageName, key, value)
		return nil
	}

	// Get specific key
	if key != "" {
		result, err := page.GetStorageItem(st, key)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", result)
		return nil
	}

	// Get all items (default)
	result, err := page.GetAllStorage(st)
	if err != nil {
		return err
	}

	// Format as JSON
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format storage: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
