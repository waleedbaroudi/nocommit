package cmd

import (
	"fmt"
	"io"
	"os"
)

var errWriter io.Writer = os.Stderr

func SetErrWriter(w io.Writer) {
	if w != nil {
		errWriter = w
	}
}

func fail(ctx string, err error, isFatal bool) {
	if err == nil {
		return
	}
	fmt.Fprintf(errWriter, "❌ %s: %v\n", ctx, err)
	if isFatal {
		fmt.Fprintln(errWriter, "This was a fatal failure, please open an issue: https://github.com/waleedbaroudi/nocommit/issues")
		os.Exit(1)
	}
	fmt.Fprintln(errWriter, "If this keeps happening, please open an issue: https://github.com/waleedbaroudi/nocommit/issues")
}
