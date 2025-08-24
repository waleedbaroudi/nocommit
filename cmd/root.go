package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:           "nocommit",
	Short:         "Block commits containing forbidden phrases",
	Long:          "nocommit installs a Git pre-commit hook that scans staged files for forbidden phrases (e.g., NOCOMMIT) and blocks the commit if found.",
	SilenceErrors: true, 
	// Version is set in version.go and injected at build time via -ldflags
}

// Execute is called by main.main().
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Route our fail(...) helper to Cobra's error writer for every command
	RootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		SetErrWriter(cmd.ErrOrStderr())
	}
}
