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
	HookLine = filepath.Join(home, ".nocommit", "hooks", "pre-commit")
	hookBlock = hookStart + "\n" + HookLine + "\n" + hookEnd + "\n"
}
