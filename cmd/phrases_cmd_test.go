package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestPhrasesCommands_AddListRemove(t *testing.T) {
	_ = setTempHome(t)

	var outBuf, errBuf bytes.Buffer
	RootCmd.SetOut(&outBuf)
	RootCmd.SetErr(&errBuf)
	t.Cleanup(func() {
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	// add
	outBuf.Reset()
	RootCmd.SetArgs([]string{"phrases", "add", "NOCOMMIT"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// list (contains)
	outBuf.Reset()
	RootCmd.SetArgs([]string{"phrases", "list"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(outBuf.String(), "NOCOMMIT") {
		t.Fatalf("expected NOCOMMIT in list; got %q", outBuf.String())
	}

	// remove
	outBuf.Reset()
	RootCmd.SetArgs([]string{"phrases", "remove", "NOCOMMIT"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	// list (gone)
	outBuf.Reset()
	RootCmd.SetArgs([]string{"phrases", "list"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("list2 failed: %v", err)
	}
	if strings.Contains(outBuf.String(), "NOCOMMIT") {
		t.Fatalf("expected removal; got %q", outBuf.String())
	}
}

func TestPhrasesCommands_AddRejectsInvalid(t *testing.T) {
	_ = setTempHome(t)

	var outBuf, errBuf bytes.Buffer
	RootCmd.SetOut(&outBuf)
	RootCmd.SetErr(&errBuf)
	t.Cleanup(func() {
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	outBuf.Reset()
	RootCmd.SetArgs([]string{"phrases", "add", "#comment"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("add comment failed: %v", err)
	}
	if !strings.Contains(outBuf.String(), "Invalid phrase") {
		t.Fatalf("expected invalid phrase message; got %q", outBuf.String())
	}
}
