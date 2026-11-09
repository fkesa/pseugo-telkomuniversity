package main

import (
	"os"

	"pseudocode-compiler/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
