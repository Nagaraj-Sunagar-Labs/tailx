// Package cli implements the tailx command line interface:
// flag parsing, validation and dispatch to the core tail logic.
package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Nagaraj-Sunagar-Labs/tailx/internal/tail"
)

// Version is the tailx release version.
const Version = "0.1.0"

// Process exit codes.
const (
	ExitOK    = 0
	ExitError = 1 // runtime failure (e.g. unreadable file)
	ExitUsage = 2 // bad invocation
)

// Usage returns the help text for tailx.
func Usage() string {
	return `tailx - output the last part of files

usage: tailx [-n lines] [-f] <file>...

Options:
  -n lines  number of trailing lines to output (default 10)
  -f        output appended data as the file grows (for 60 seconds)

  --help     show this help and exit
  --version  show version and exit
`
}

// Execute parses args (without the program name) and runs the tail logic.
// It returns the process exit code.
func Execute(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("tailx", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, Usage()) }

	lines := fs.Int("n", 10, "number of trailing lines to output")
	follow := fs.Bool("f", false, "output appended data as the file grows")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return ExitOK
		}
		return ExitUsage
	}

	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(stderr, "tailx: no input file provided")
		fs.Usage()
		return ExitUsage
	}

	opts := tail.Options{
		Lines:  *lines,
		Follow: *follow,
		Files:  files,
	}
	if err := tail.Run(stdout, opts); err != nil {
		fmt.Fprintf(stderr, "tailx: %v\n", err)
		return ExitError
	}
	return ExitOK
}
