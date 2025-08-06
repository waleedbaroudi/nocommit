/*
Copyright © 2025 Walid Baroudi wa.baroudi9@gmail.com
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/waleedbaroudi/nocommit/cmd"
)

func init() {
	cobra.OnInitialize(setup)
}

func setup() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Cannot determine home directory:", err)
		// TODO: add logs to direct the user
		return
	}

	nocommitDir := filepath.Join(home, ".nocommit")
	phrasesFile := filepath.Join(nocommitDir, "phrases.txt")
	configFile := filepath.Join(nocommitDir, "config.yaml")

	err = os.MkdirAll(nocommitDir, 0755)
	if err != nil {
		fmt.Println("❌ Failed to create .nocommit directory:", err)
		return
	}

	// Create phrases.txt with defaults if not exists
	if _, err := os.Stat(phrasesFile); os.IsNotExist(err) {
		defaultPhrases := []byte("NOCOMMIT\n")
		os.WriteFile(phrasesFile, defaultPhrases, 0644)
	}

	// Create config.json if not exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.WriteFile(configFile, []byte(""), 0644)
	}
}

func main() {
	cmd.Execute()
}
