package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable nocommit in this repo",
	Run: func(cmd *cobra.Command, args []string) {
		gitDir, err := exec.Command("git", "rev-parse", "--git-dir").Output()
		if err != nil {
			fmt.Println("❌ Not a git repository. Please call this command within a git repo.")
			return
		}
		hooksDir := filepath.Join(strings.TrimSpace(string(gitDir)), "hooks")
		if err := os.MkdirAll(hooksDir, 0755); err != nil {
			fail("Create hooks dir", err, true)
			return
		}

		hook := filepath.Join(hooksDir, "pre-commit")

		content := ""
		if b, err := os.ReadFile(hook); err == nil {
			content = string(b)
			if strings.Contains(content, HookLine) {
				fmt.Println("✅ nocommit already enabled for this repo.")
				return
			}
		}

		// create or append
		if content == "" {
			content = "#!/bin/sh\n\n" + hookBlock
		} else {
			content += "\n\n" + hookBlock
		}

		if err := os.WriteFile(hook, []byte(content), 0755); err != nil {
			fail("Write pre-commit", err, true)
			return
		}
		fmt.Println("✅ nocommit enabled for this repo.")
	},
}

func init() {
	RootCmd.AddCommand(enableCmd)
}
