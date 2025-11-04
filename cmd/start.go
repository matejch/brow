package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	profileDir string
	useProfile bool
	headless   bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch Chrome with remote debugging enabled",
	Long: `Starts Chrome with remote debugging on port 9222.
By default, uses a temporary profile for clean sessions.
Use --profile to maintain cookies and login state.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVar(&useProfile, "profile", false, "Use persistent profile (maintains cookies/logins)")
	startCmd.Flags().StringVar(&profileDir, "profile-dir", "", "Custom profile directory path")
	startCmd.Flags().BoolVar(&headless, "headless", false, "Run Chrome in headless mode")
}

func runStart(cmd *cobra.Command, args []string) error {
	chromePath, err := findChrome()
	if err != nil {
		return fmt.Errorf("Chrome not found: %w", err)
	}

	// Determine profile directory
	var userDataDir string
	if profileDir != "" {
		userDataDir = profileDir
	} else if useProfile {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		userDataDir = filepath.Join(home, ".brow-profile")
	} else {
		// Use temp directory for clean sessions
		userDataDir, err = os.MkdirTemp("", "brow-*")
		if err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	// Build Chrome arguments
	chromeArgs := []string{
		fmt.Sprintf("--remote-debugging-port=%d", 9222),
		fmt.Sprintf("--user-data-dir=%s", userDataDir),
		"--no-first-run",
		"--no-default-browser-check",
	}

	if headless {
		chromeArgs = append(chromeArgs, "--headless=new")
	}

	// Start Chrome
	chromeCmd := exec.Command(chromePath, chromeArgs...)
	chromeCmd.Stdout = os.Stdout
	chromeCmd.Stderr = os.Stderr

	if err := chromeCmd.Start(); err != nil {
		return fmt.Errorf("failed to start Chrome: %w", err)
	}

	fmt.Printf("Chrome started (PID: %d)\n", chromeCmd.Process.Pid)
	fmt.Printf("Remote debugging: http://localhost:9222\n")
	fmt.Printf("Profile: %s\n", userDataDir)

	return nil
}

// findChrome attempts to locate the Chrome executable on the system
func findChrome() (string, error) {
	var candidates []string

	switch runtime.GOOS {
	case "darwin":
		candidates = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	case "linux":
		candidates = []string{
			"google-chrome",
			"google-chrome-stable",
			"chromium",
			"chromium-browser",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
	case "windows":
		candidates = []string{
			filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(os.Getenv("LocalAppData"), "Google", "Chrome", "Application", "chrome.exe"),
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Try to find Chrome
	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Chrome not found in standard locations")
}
