package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var phrasesCmd = &cobra.Command{
	Use:   "phrases",
	Short: "Manage block phrases",
}

var phrasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List phrases",
	Run: func(cmd *cobra.Command, args []string) {
		ps, err := readPhrases()
		if err != nil {
			fail("Read phrases", err, false)
			return
		}
		if len(ps) == 0 {
			fmt.Println("(no phrases)")
			return
		}
		for _, p := range ps {
			fmt.Println(p)
		}
	},
}

var phrasesAddCmd = &cobra.Command{
	Use:   "add <phrase>",
	Short: "Add a phrase",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		phrase := strings.TrimSpace(args[0])
		if phrase == "" || strings.HasPrefix(phrase, "#") {
			fmt.Println("❌ Invalid phrase")
			_ = cmd.Usage()
			return
		}
		ps, err := readPhrases()
		if err != nil {
			fail("Read phrases", err, false)
			return
		}
		for _, p := range ps {
			if p == phrase {
				fmt.Println("✅ Already exists")
				return
			}
		}
		ps = append(ps, phrase)
		if err := writePhrases(ps); err != nil {
			fail("Write phrases", err, false)
			return
		}
		fmt.Printf("✅ Added '%s' to the list of block phrases", phrase)
	},
}

var phrasesRemoveCmd = &cobra.Command{
	Use:   "remove <phrase>",
	Short: "Remove a phrase",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := strings.TrimSpace(args[0])
		ps, err := readPhrases()
		if err != nil {
			fail("Read phrases", err, false)
			return
		}
		out := make([]string, 0, len(ps))
		removed := false
		for _, p := range ps {
			if p == target {
				removed = true
				continue
			}
			out = append(out, p)
		}
		if !removed {
			fmt.Println("⚠️ Not found")
			return
		}
		if err := writePhrases(out); err != nil {
			fail("Write phrases", err, false)
			return
		}
		fmt.Printf("✅ Removed '%s' from the list of block phrases", target)
	},
}

func init() {
	phrasesCmd.AddCommand(phrasesListCmd, phrasesAddCmd, phrasesRemoveCmd)
	rootCmd.AddCommand(phrasesCmd)
}

// helpers

func phrasesPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".nocommit", "phrases.txt")
}

func readPhrases() ([]string, error) {
	path := phrasesPath()
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		// create empty file if missing
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			return nil, err
		}
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var res []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		res = append(res, line)
	}
	return res, sc.Err()
}

func writePhrases(list []string) error {
	path := phrasesPath()
	var b strings.Builder
	for _, p := range list {
		if p == "" {
			continue
		}
		b.WriteString(p)
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0644)
}
