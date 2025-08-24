package cmd

import (
	"os"
	"path/filepath"
)

var (
	hookStart string
	hookEnd   string
	HookLine  string
	hookBlock string
)

func init() {
	hookStart = "# nocommit hook start"
	hookEnd = "# nocommit hook end"
	home, _ := os.UserHomeDir()
	hookPath := filepath.Join(home, ".nocommit", "hooks", "pre-commit")
	correctedHookPath := "\"" + filepath.ToSlash(hookPath) + "\""
	HookLine = "sh " + correctedHookPath
	hookBlock = hookStart + "\n" + HookLine + "\n" + hookEnd + "\n"
}
