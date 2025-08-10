package cmd

import (
	"os"
	"path/filepath"
)

var (
	hookStart string
	hookEnd   string
	hookLine  string
	hookBlock string
)

func init() {
	hookStart = "# nocommit hook start"
	hookEnd = "# nocommit hook end"
	home, _ := os.UserHomeDir()
	hookLine = filepath.Join(home, ".nocommit", "hooks", "pre-commit")
	hookBlock = hookStart + "\n" + hookLine + "\n" + hookEnd + "\n"
}
