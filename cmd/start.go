package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/matejch/brow/pkg/browser"
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
	Long: `Starts Chrome with remote debugging enabled.
Port can be configured with --port flag or BROW_DEBUG_PORT env var (default: 9222).
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

func runStart(_ *cobra.Command, _ []string) error {
	chromePath, err := findChrome()
	if err != nil {
		return fmt.Errorf("Chrome not found: %w", err)
	}

	// Resolve the port to use (flag > env > default)
	debugPort := browser.ResolvePort(Port)

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
		fmt.Sprintf("--remote-debugging-port=%d", debugPort),
		fmt.Sprintf("--user-data-dir=%s", userDataDir),
		"--no-first-run",
		"--no-default-browser-check",
		"about:blank", // Force Chrome to create an initial tab
	}

	if headless {
		chromeArgs = append(chromeArgs, "--headless=new")
	}

	// Start Chrome
	chromeCmd := exec.Command(chromePath, chromeArgs...)

	// Detach Chrome from the parent process so it survives after brow exits
	// Use Setsid (not Setpgid) to create a new session and fully detach from terminal
	chromeCmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Create new session - makes Chrome a true daemon
	}

	// Disconnect stdio streams to prevent hanging
	// Chrome will run independently in the background
	chromeCmd.Stdout = nil
	chromeCmd.Stderr = nil
	chromeCmd.Stdin = nil

	if err := chromeCmd.Start(); err != nil {
		return fmt.Errorf("failed to start Chrome: %w", err)
	}

	// Save PID before calling Release() (Release invalidates the handle)
	pid := chromeCmd.Process.Pid

	// Release the process so it continues after brow exits
	if err := chromeCmd.Process.Release(); err != nil {
		return fmt.Errorf("failed to release Chrome process: %w", err)
	}

	fmt.Printf("Chrome started (PID: %d)\n", pid)
	fmt.Printf("Remote debugging: http://localhost:%d\n", debugPort)
	fmt.Printf("Profile: %s\n", userDataDir)
	fmt.Println("Chrome is running in the background. Close Chrome manually when done.")

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
