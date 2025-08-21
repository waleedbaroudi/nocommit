package cmd_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/waleedbaroudi/nocommit/cmd"
)

func TestEnableDisable(t *testing.T) {
	tmp := t.TempDir()

	// Isolate home (Windows + Unix)
	tmpHome := filepath.Join(tmp, "home")
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	// Init repo (use shared mustRun from commit_blocking_test.go)
	env := os.Environ()
	mustRun(t, env, tmp, "git", "init")

	// In-process RootCmd uses CWD → run inside the repo
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Capture CLI output
	var outBuf, errBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outBuf)
	cmd.RootCmd.SetErr(&errBuf)
	t.Cleanup(func() {
		cmd.RootCmd.SetOut(os.Stdout)
		cmd.RootCmd.SetErr(os.Stderr)
	})

	// enable
	cmd.RootCmd.SetArgs([]string{"enable"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("enable failed: %v\nstderr:\n%s", err, errBuf.String())
	}

	// verify hook
	gitDirBytes, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v", err)
	}
	hookPath := filepath.Join(strings.TrimSpace(string(gitDirBytes)), "hooks", "pre-commit")

	data, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("read pre-commit: %v", err)
	}
	if !strings.Contains(string(data), cmd.HookLine) {
		t.Errorf("expected hook line in pre-commit hook")
	}

	// disable
	cmd.RootCmd.SetArgs([]string{"disable"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("disable failed: %v\nstderr:\n%s", err, errBuf.String())
	}

	data, err = os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("read pre-commit after disable: %v", err)
	}
	if strings.Contains(string(data), cmd.HookLine) {
		t.Errorf("expected hook line removed after disable")
	}
}

func mustRun(t *testing.T, env []string, dir string, name string, args ...string) {
	t.Helper()
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Env = env
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}
