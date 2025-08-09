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
		fail("Can't find home directory", err, true)
		return
	}

	dir := filepath.Join(home, ".nocommit")
	hooksDir := filepath.Join(dir, "hooks")
	phrases := filepath.Join(dir, "phrases.txt")
	config := filepath.Join(dir, "config.yaml")

	firstRun := false

	if err := os.MkdirAll(dir, 0755); err != nil {
		fail("Create ~/.nocommit directory", err, true)
		return
	}

	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		fail("Create hooks directory", err, true)
		return
	}

	if _, err := os.Stat(phrases); os.IsNotExist(err) {
		firstRun = true
		if err := os.WriteFile(phrases, []byte("NOCOMMIT\n"), 0644); err != nil {
			fail("Write phrases.txt", err, false)
		}
	}

	if _, err := os.Stat(config); os.IsNotExist(err) {
		firstRun = true
		if err := os.WriteFile(config, []byte("version: 1\ncaseSensitive: false\n"), 0600); err != nil {
			fail("Write config.yaml", err, false)
		}
	}

	installHook(hooksDir)

	if firstRun {
		fmt.Println("✅ nocommit set up at ~/.nocommit")
	}
}

func installHook(hooksDir string) {
	src := filepath.Join("assets", "pre-commit")
	dst := filepath.Join(hooksDir, "pre-commit")

	data, err := os.ReadFile(src)
	if err != nil {
		fail("Read hook template", err, true)
	}

	if err := os.WriteFile(dst, data, 0755); err != nil {
		fail("Write hook", err, true)
	}

}

func fail(ctx string, err error, isFatal bool) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "❌ %s: %v\n", ctx, err)
	if isFatal {
		fmt.Fprintln(os.Stderr, "This was a fatal failure, please open an issue: https://github.com/waleedbaroudi/nocommit/issues")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "If this keeps happening, please open an issue: https://github.com/waleedbaroudi/nocommit/issues")
}

func main() {
	cmd.Execute()
}
