package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable nocommit in this repo",
	Run: func(cmd *cobra.Command, args []string) {
		gitDir, err := exec.Command("git", "rev-parse", "--git-dir").Output()
		if err != nil {
			fmt.Println("❌ Not a git repository. Please call this command within a git repo.")
			return
		}

		hook := filepath.Join(strings.TrimSpace(string(gitDir)), "hooks", "pre-commit")

		// Read hook file
		data, err := os.ReadFile(hook)
		if err != nil {
			fmt.Println("⚠️ No pre-commit hook found.")
			return
		}
		content := strings.ReplaceAll(string(data), "\r\n", "\n") // handle CRLF to be safe on Windows

		// Check if block exists
		if !containsHook(content) {
			fmt.Println("⚠️ nocommit is not enabled for this repo.")
			return
		}

		// Remove the block
		newContent := removeHookBlock(content)

		// Write back updated hook
		if err := os.WriteFile(hook, []byte(newContent), 0755); err != nil {
			fail("Write pre-commit", err, false)
			return
		}
		fmt.Println("✅ nocommit disabled for this repo.")
	},
}

// Returns true if the given text contains any part of the hook block
func containsHook(text string) bool {
	return strings.Contains(text, hookStart) || strings.Contains(text, hookLine) || strings.Contains(text, hookEnd)
}

func removeHookBlock(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		if containsHook(line) {
			continue
		}
		result = append(result, line)
	}

	// Trim trailing newlines
	return strings.TrimRight(strings.Join(result, "\n"), "\n") + "\n"
}

func init() {
	rootCmd.AddCommand(disableCmd)
}
