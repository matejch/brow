package cmd

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/matej/brow/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	xpath   bool
	closest bool
)

var pickCmd = &cobra.Command{
	Use:   "pick",
	Short: "Interactive element picker to get CSS selectors",
	Long: `Injects an interactive overlay into the page that allows you to click on elements
and get their CSS selector or XPath.

The overlay highlights elements on hover and copies the selector on click.
Press ESC to exit the picker mode.`,
	RunE: runPick,
}

func init() {
	rootCmd.AddCommand(pickCmd)
	pickCmd.Flags().BoolVarP(&xpath, "xpath", "x", false, "Return XPath instead of CSS selector")
}

func runPick(cmd *cobra.Command, args []string) error {
	// JavaScript to inject the element picker
	pickerScript := `
(function() {
	if (window.__browPicker) {
		console.log('Picker already active');
		return;
	}

	window.__browPicker = true;
	let overlay = null;
	let selectedElement = null;

	// Get CSS selector for an element
	function getCSSSelector(el) {
		if (el.id) return '#' + el.id;

		let path = [];
		while (el.parentElement) {
			let selector = el.tagName.toLowerCase();
			if (el.className) {
				selector += '.' + Array.from(el.classList).join('.');
			}

			let siblings = Array.from(el.parentElement.children).filter(
				e => e.tagName === el.tagName
			);
			if (siblings.length > 1) {
				let index = siblings.indexOf(el) + 1;
				selector += ':nth-of-type(' + index + ')';
			}

			path.unshift(selector);
			el = el.parentElement;
		}
		return path.join(' > ');
	}

	// Get XPath for an element
	function getXPath(el) {
		if (el.id) return '//*[@id="' + el.id + '"]';

		let path = [];
		while (el.parentElement) {
			let siblings = Array.from(el.parentElement.children).filter(
				e => e.tagName === el.tagName
			);
			let index = siblings.indexOf(el) + 1;
			path.unshift(el.tagName.toLowerCase() + '[' + index + ']');
			el = el.parentElement;
		}
		return '/' + path.join('/');
	}

	// Create overlay
	overlay = document.createElement('div');
	overlay.style.cssText = 'position: absolute; border: 2px solid red; pointer-events: none; z-index: 999999; background: rgba(255, 0, 0, 0.1);';
	document.body.appendChild(overlay);

	// Info box
	let infoBox = document.createElement('div');
	infoBox.style.cssText = 'position: fixed; top: 10px; right: 10px; background: black; color: white; padding: 10px; z-index: 1000000; font-family: monospace; font-size: 12px;';
	infoBox.textContent = 'Hover to highlight, Click to select, ESC to exit';
	document.body.appendChild(infoBox);

	// Mouse move handler
	function handleMouseMove(e) {
		if (e.target === overlay || e.target === infoBox) return;

		let rect = e.target.getBoundingClientRect();
		overlay.style.left = (rect.left + window.scrollX) + 'px';
		overlay.style.top = (rect.top + window.scrollY) + 'px';
		overlay.style.width = rect.width + 'px';
		overlay.style.height = rect.height + 'px';
	}

	// Click handler
	function handleClick(e) {
		e.preventDefault();
		e.stopPropagation();

		selectedElement = e.target;
		let selector = ${xpath ? 'getXPath(selectedElement)' : 'getCSSSelector(selectedElement)'};

		// Store in window for retrieval
		window.__browPickedSelector = selector;

		cleanup();
	}

	// ESC handler
	function handleKeyDown(e) {
		if (e.key === 'Escape') {
			cleanup();
		}
	}

	// Cleanup function
	function cleanup() {
		document.removeEventListener('mousemove', handleMouseMove);
		document.removeEventListener('click', handleClick, true);
		document.removeEventListener('keydown', handleKeyDown);
		if (overlay) overlay.remove();
		if (infoBox) infoBox.remove();
		window.__browPicker = false;
	}

	// Attach event listeners
	document.addEventListener('mousemove', handleMouseMove);
	document.addEventListener('click', handleClick, true);
	document.addEventListener('keydown', handleKeyDown);
})();
`

	// Inject the picker script
	if err := browser.Run(chromedp.Evaluate(pickerScript, nil)); err != nil {
		return fmt.Errorf("failed to inject picker: %w", err)
	}

	fmt.Println("Element picker activated!")
	fmt.Println("Hover over elements to highlight, click to select, press ESC to exit.")
	fmt.Println("")
	fmt.Println("After selecting an element, run:")
	fmt.Println("  brow eval 'window.__browPickedSelector'")
	fmt.Println("")
	fmt.Println("To get the selected element's selector.")

	return nil
}
