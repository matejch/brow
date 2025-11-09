package client

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/matejch/brow/pkg/config"
	"github.com/matejch/brow/pkg/operations"
)

// Page represents a browser page/tab and provides methods for automation
type Page struct {
	ctx    context.Context
	config *config.Config
}

// Navigate navigates to the specified URL
func (p *Page) Navigate(url string, waitReady bool) (*operations.NavigationResult, error) {
	return operations.Navigate(p.ctx, url, waitReady)
}

// Eval executes JavaScript in the page context and returns the result
func (p *Page) Eval(script string) (interface{}, error) {
	return operations.Evaluate(p.ctx, script)
}

// Screenshot captures a screenshot of the current page
func (p *Page) Screenshot(opts operations.ScreenshotOptions) ([]byte, error) {
	return operations.CaptureScreenshot(p.ctx, opts)
}

// PDF generates a PDF from the current page
func (p *Page) PDF(opts operations.PDFOptions) ([]byte, error) {
	return operations.GeneratePDF(p.ctx, opts)
}

// GetCookies retrieves all cookies, optionally filtered by domain
func (p *Page) GetCookies(domain string) ([]*network.Cookie, error) {
	return operations.GetCookies(p.ctx, domain)
}

// SetCookie sets a cookie
func (p *Page) SetCookie(cookie string) error {
	return operations.SetCookie(p.ctx, cookie)
}

// ClearCookies clears all browser cookies
func (p *Page) ClearCookies() error {
	return operations.ClearCookies(p.ctx)
}

// GetAllStorage retrieves all items from the specified storage type
func (p *Page) GetAllStorage(storageType operations.StorageType) (map[string]interface{}, error) {
	return operations.GetAllStorage(p.ctx, storageType)
}

// GetStorageItem retrieves a specific item from storage
func (p *Page) GetStorageItem(storageType operations.StorageType, key string) (interface{}, error) {
	return operations.GetStorageItem(p.ctx, storageType, key)
}

// SetStorageItem sets a value in storage
func (p *Page) SetStorageItem(storageType operations.StorageType, key, value string) error {
	return operations.SetStorageItem(p.ctx, storageType, key, value)
}

// RemoveStorageItem removes an item from storage
func (p *Page) RemoveStorageItem(storageType operations.StorageType, key string) error {
	return operations.RemoveStorageItem(p.ctx, storageType, key)
}

// ClearStorage clears all items from the specified storage
func (p *Page) ClearStorage(storageType operations.StorageType) error {
	return operations.ClearStorage(p.ctx, storageType)
}

// InjectPicker injects an interactive element picker into the page
func (p *Page) InjectPicker(useXPath bool) error {
	return operations.InjectPicker(p.ctx, useXPath)
}

// GetPickedSelector retrieves the selector picked by the user
func (p *Page) GetPickedSelector() (string, error) {
	return operations.GetPickedSelector(p.ctx)
}

// Context returns the underlying context for advanced usage
func (p *Page) Context() context.Context {
	return p.ctx
}
