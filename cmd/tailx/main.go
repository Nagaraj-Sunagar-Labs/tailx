// Command tailx is a command line utility similar to the Unix tail command.
//
// main.go is intentionally minimal: it only handles version and help.
// All other logic lives in internal/.
package main

import (
	"fmt"
	"os"

	"github.com/Nagaraj-Sunagar-Labs/tailx/internal/cli"
)

func main() {
	args := os.Args[1:]

	// With no arguments, show help.
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, cli.Usage())
		os.Exit(cli.ExitUsage)
	}

	switch args[0] {
	case "--version":
		fmt.Printf("tailx version %s\n", cli.Version)
		os.Exit(0)
	case "--help":
		fmt.Print(cli.Usage())
		os.Exit(0)
	}

	// Everything else (-n, -f) is handled in internal/cli.
	os.Exit(cli.Execute(args, os.Stdout, os.Stderr))
}
