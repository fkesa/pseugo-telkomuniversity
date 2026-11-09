package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"pseudocode-compiler/internal/compiler"
)

func main() {
	input := flag.String("in", "", "path to pseudocode source file")
	output := flag.String("out", "", "path for generated .go file (default: stdout)")
	showTokens := flag.Bool("tokens", false, "print lexer tokens instead of compiling")
	flag.Parse()
	if *input == "" {
		fmt.Fprintln(os.Stderr, "usage: pseudoc -in program.pseudo [-out program.go] [-tokens]")
		os.Exit(2)
	}
	source, err := os.ReadFile(*input)
	if err != nil {
		fail(err)
	}
	if *showTokens {
		tokens, err := compiler.NewLexer(string(source)).AllTokens()
		if err != nil {
			fail(err)
		}
		for _, t := range tokens {
			fmt.Println(t.String())
		}
		return
	}
	result, err := compiler.Compile(string(source))
	if err != nil {
		fail(err)
	}
	if *output == "" {
		fmt.Print(result.GoSource)
		return
	}
	if err = os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		fail(err)
	}
	if err = os.WriteFile(*output, []byte(result.GoSource), 0644); err != nil {
		fail(err)
	}
	fmt.Printf("Generated %s\n", *output)
}
func fail(err error) { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
