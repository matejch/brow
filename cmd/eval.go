package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/matej/brow/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	rawOutput  bool
	jsonOutput bool
)

var evalCmd = &cobra.Command{
	Use:   "eval <javascript>",
	Short: "Execute JavaScript in the current page",
	Long: `Executes JavaScript code in the current page context and returns the result.
Results are automatically formatted as JSON unless --raw is specified.`,
	Args: cobra.ExactArgs(1),
	RunE: runEval,
}

func init() {
	rootCmd.AddCommand(evalCmd)
	evalCmd.Flags().BoolVarP(&rawOutput, "raw", "r", false, "Output raw result without formatting")
	evalCmd.Flags().BoolVarP(&jsonOutput, "json", "j", true, "Format output as JSON (default true)")
}

func runEval(_ *cobra.Command, args []string) error {
	script := args[0]

	// Attach to existing tab
	ctx, cancel, err := browser.GetExistingTabContext()
	if err != nil {
		return err
	}
	defer cancel()

	var result interface{}
	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &result)); err != nil {
		return fmt.Errorf("failed to evaluate JavaScript: %w", err)
	}

	// Handle output formatting
	if rawOutput {
		fmt.Printf("%v\n", result)
	} else if jsonOutput {
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format result as JSON: %w", err)
		}
		fmt.Println(string(output))
	} else {
		fmt.Printf("%v\n", result)
	}

	return nil
}
