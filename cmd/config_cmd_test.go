package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigGetSetList(t *testing.T) {
	home := setTempHome(t)

	// seed a minimal config.yaml
	cfgPath := filepath.Join(home, ".nocommit", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, []byte("version: 1\ncaseSensitive: false\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var outBuf, errBuf bytes.Buffer
	RootCmd.SetOut(&outBuf)
	RootCmd.SetErr(&errBuf)
	t.Cleanup(func() {
		RootCmd.SetOut(os.Stdout)
		RootCmd.SetErr(os.Stderr)
	})

	// get existing key
	outBuf.Reset()
	RootCmd.SetArgs([]string{"config", "version"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := outBuf.String(); got == "" || got[:8] != "version:" && got[:8] != "version" {
		t.Fatalf("expected version output, got %q", got)
	}

	// set new key
	outBuf.Reset()
	RootCmd.SetArgs([]string{"config", "caseSensitive", "true"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// get updated key
	outBuf.Reset()
	RootCmd.SetArgs([]string{"config", "caseSensitive"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := outBuf.String(); got == "" || !bytes.Contains([]byte(got), []byte("true")) {
		t.Fatalf("expected true, got %q", got)
	}

	// list prints raw yaml
	outBuf.Reset()
	RootCmd.SetArgs([]string{"config", "list"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	got := outBuf.String()
	if !strings.Contains(got, `caseSensitive: true`) &&
		!strings.Contains(got, `caseSensitive: "true"`) {
		t.Fatalf("list missing updated key; got:\n%s", got)
	}
}
