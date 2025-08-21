package cmd

import "testing"

func TestContainsHookDetectsHookLine(t *testing.T) {
	content := "#!/bin/sh\n" + HookLine + "\necho ok\n"
	if !containsHook(content) {
		t.Fatalf("expected containsHook to detect HookLine")
	}
}

func TestRemoveHookBlockRemovesMarkersOnly(t *testing.T) {
	before := "x\ny\n" + HookLine + "\nz\n"
	after := removeHookBlock(before)
	if containsHook(after) {
		t.Fatalf("expected hook markers removed")
	}
	if after != "x\ny\nz\n" {
		t.Fatalf("unexpected content after removal: %q", after)
	}
}
