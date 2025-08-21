package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/waleedbaroudi/nocommit/cmd"
)

func withTempHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)        // Unix / Git Bash
	t.Setenv("USERPROFILE", tmp) // Windows
	return tmp
}

func TestPhrasesAddListRemove(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outBuf)
	cmd.RootCmd.SetErr(&errBuf)
	t.Cleanup(func() {
		cmd.RootCmd.SetOut(os.Stdout)
		cmd.RootCmd.SetErr(os.Stderr)
	})

	home := withTempHome(t) // must set HOME and USERPROFILE
	_ = home                // no need to mkdir ~/.nocommit; created on demand

	// add
	outBuf.Reset()
	cmd.RootCmd.SetArgs([]string{"phrases", "add", "NOCOMMIT"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("phrases add failed: %v", err)
	}

	// list (should contain)
	outBuf.Reset()
	cmd.RootCmd.SetArgs([]string{"phrases", "list"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("phrases list failed: %v", err)
	}
	if !strings.Contains(outBuf.String(), "NOCOMMIT") {
		t.Errorf("expected 'NOCOMMIT' in list, got: %s", outBuf.String())
	}

	// remove
	outBuf.Reset()
	cmd.RootCmd.SetArgs([]string{"phrases", "remove", "NOCOMMIT"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("phrases remove failed: %v", err)
	}

	// list again (should be gone)
	outBuf.Reset()
	cmd.RootCmd.SetArgs([]string{"phrases", "list"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("phrases list after remove failed: %v", err)
	}
	if strings.Contains(outBuf.String(), "NOCOMMIT") {
		t.Errorf("expected phrase removed, got: %s", outBuf.String())
	}
}
