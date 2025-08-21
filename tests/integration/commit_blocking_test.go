//go:build !windows
// +build !windows

package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCommitBlocking(t *testing.T) {
	// Build a hermetic nocommit binary from the repo root
	repoRoot := moduleRoot(t)
	tmp := t.TempDir()
	bin := filepath.Join(tmp, "nocommit"+exe())
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = repoRoot
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}

	// Sandbox HOME/USERPROFILE so ~/.nocommit is isolated
	tmpHome := filepath.Join(tmp, "home")
	_ = os.MkdirAll(tmpHome, 0o755)

	homePOSIX := strings.ReplaceAll(tmpHome, "\\", "/")
	env := append(os.Environ(),
		"HOME="+homePOSIX,      // shell-friendly for the hook
		"USERPROFILE="+tmpHome, // Windows-friendly for Go/os.UserHomeDir
	)

	// Create a temp git repo
	repo := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	mustRun(t, env, repo, "git", "init")
	mustRun(t, env, repo, "git", "config", "user.email", "test@example.com")
	mustRun(t, env, repo, "git", "config", "user.name", "Test User")

	// Enable nocommit in this repo
	mustRunBin(t, env, repo, bin, "enable")

	// Ensure phrase exists (your first run likely writes NOCOMMIT; this keeps the test explicit)
	mustRunBin(t, env, repo, bin, "phrases", "add", "NOCOMMIT")

	// Create a file containing the forbidden phrase
	bad := filepath.Join(repo, "bad.txt")
	if err := os.WriteFile(bad, []byte("this line has NOCOMMIT tag"), 0o644); err != nil {
		t.Fatalf("write bad.txt: %v", err)
	}
	mustRun(t, env, repo, "git", "add", "bad.txt")

	// Commit should be blocked by the pre-commit hook
	cmd := exec.Command("git", "commit", "-m", "should fail")
	cmd.Dir = repo
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected commit to fail, but it succeeded. Output:\n%s", out)
	}
	// Optional: assert the output mentions nocommit (best-effort)
	if !strings.Contains(string(out), "nocommit") {
		t.Logf("commit failed as expected, output did not contain 'nocommit':\n%s", out)
	}

	// Now make a clean commit that should pass
	good := filepath.Join(repo, "good.txt")
	if err := os.WriteFile(good, []byte("all clear"), 0o644); err != nil {
		t.Fatalf("write good.txt: %v", err)
	}
	mustRun(t, env, repo, "git", "add", "good.txt")
	cmd = exec.Command("git", "commit", "-m", "should pass")
	cmd.Dir = repo
	cmd.Env = env
	if out, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("expected clean commit to pass, but it failed: %v\n%s", err, out)
	}
}

func moduleRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func exe() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func mustRunBin(t *testing.T, env []string, dir, bin string, args ...string) {
	t.Helper()
	c := exec.Command(bin, args...)
	c.Dir = dir
	c.Env = env
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("%s %v failed: %v\n%s", bin, args, err, out)
	}
}
