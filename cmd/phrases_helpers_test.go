package cmd

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func setTempHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // Windows
	return home
}

func TestReadPhrasesCreatesFileAndSkipsComments(t *testing.T) {
	_ = setTempHome(t)

	// file doesn't exist yet → readPhrases creates it and returns empty
	ps, err := readPhrases()
	if err != nil {
		t.Fatalf("readPhrases: %v", err)
	}
	if len(ps) != 0 {
		t.Fatalf("expected empty list, got %v", ps)
	}

	// write with blanks and comments
	path := phrasesPath()
	if err := os.WriteFile(path, []byte("\n# comment\nNOCOMMIT\n  \n#x\nTODO\n"), 0o644); err != nil {
		t.Fatalf("write phrases: %v", err)
	}
	ps, err = readPhrases()
	if err != nil {
		t.Fatalf("readPhrases: %v", err)
	}
	got := strings.Join(ps, ",")
	if got != "NOCOMMIT,TODO" {
		t.Fatalf("want NOCOMMIT,TODO; got %q", got)
	}
}

func TestWritePhrasesRoundTrip(t *testing.T) {
	_ = setTempHome(t)
	_, _ = readPhrases()
	in := []string{"ONE", "", "TWO"}
	if err := writePhrases(in); err != nil {
		t.Fatalf("writePhrases: %v", err)
	}
	ps, err := readPhrases()
	if err != nil {
		t.Fatalf("readPhrases: %v", err)
	}
	if len(ps) != 2 || ps[0] != "ONE" || ps[1] != "TWO" {
		t.Fatalf("round-trip mismatch: %v", ps)
	}

	// ensure file path under HOME
	p := phrasesPath()
	if runtime.GOOS == "windows" && !strings.Contains(p, "\\.nocommit\\phrases.txt") && !strings.Contains(p, "/.nocommit/phrases.txt") {
		t.Fatalf("unexpected phrases path: %s", p)
	}
}
