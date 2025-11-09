package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
)

// StorageType represents the type of browser storage
type StorageType string

const (
	// LocalStorage represents window.localStorage
	LocalStorage StorageType = "localStorage"
	// SessionStorage represents window.sessionStorage
	SessionStorage StorageType = "sessionStorage"
)

// GetAllStorage retrieves all items from the specified storage type
func GetAllStorage(ctx context.Context, storageType StorageType) (map[string]interface{}, error) {
	storageName := string(storageType)

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
	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &result)); err != nil {
		return nil, fmt.Errorf("failed to get storage items: %w", err)
	}

	// Convert to map[string]interface{}
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return make(map[string]interface{}), nil
}

// GetStorageItem retrieves a specific item from storage
func GetStorageItem(ctx context.Context, storageType StorageType, key string) (interface{}, error) {
	// Safely escape the key using JSON encoding
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("failed to escape key: %w", err)
	}

	script := fmt.Sprintf("%s.getItem(%s)", string(storageType), string(keyJSON))

	var result interface{}
	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &result)); err != nil {
		return nil, fmt.Errorf("failed to get value: %w", err)
	}

	return result, nil
}

// SetStorageItem sets a value in storage
func SetStorageItem(ctx context.Context, storageType StorageType, key, value string) error {
	// Safely escape both key and value using JSON encoding
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("failed to escape key: %w", err)
	}
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to escape value: %w", err)
	}

	script := fmt.Sprintf("%s.setItem(%s, %s)", string(storageType), string(keyJSON), string(valueJSON))

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, nil)); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// RemoveStorageItem removes an item from storage
func RemoveStorageItem(ctx context.Context, storageType StorageType, key string) error {
	// Safely escape the key using JSON encoding
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("failed to escape key: %w", err)
	}

	script := fmt.Sprintf("%s.removeItem(%s)", string(storageType), string(keyJSON))

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, nil)); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// ClearStorage clears all items from the specified storage
func ClearStorage(ctx context.Context, storageType StorageType) error {
	script := fmt.Sprintf("%s.clear()", string(storageType))

	if err := chromedp.Run(ctx, chromedp.Evaluate(script, nil)); err != nil {
		return fmt.Errorf("failed to clear storage: %w", err)
	}

	return nil
}
