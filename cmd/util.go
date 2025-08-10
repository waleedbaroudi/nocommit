package cmd

import (
	"fmt"
	"os"
)

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
