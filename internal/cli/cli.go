package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"pseudocode-compiler/internal/compiler"
)

// Run executes the CLI with the provided arguments and standard I/O streams.
func Run(args []string, stdout, stderr io.Writer) int {
	return RunWithIO(args, os.Stdin, stdout, stderr)
}

// RunWithIO executes the CLI with explicit stdin, stdout, and stderr streams.
func RunWithIO(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	cmd := args[0]

	switch strings.ToLower(cmd) {
	case "run":
		return runPseudo(args[1:], stdin, stdout, stderr)
	case "convert":
		return runConvert(args[1:], stdout, stderr)
	case "check":
		return runCheck(args[1:], stdout, stderr)
	case "tokens":
		return runTokens(args[1:], stdout, stderr)
	case "help", "--help", "-help", "-h":
		printUsage(stdout)
		return 0
	default:
		// Check if user passed a file directly as first argument: pseugo <input> [output]
		if !strings.HasPrefix(cmd, "-") {
			if _, err := os.Stat(cmd); err == nil {
				return runAuto(cmd, args[1:], stdout, stderr)
			}
		}
		fmt.Fprintf(stderr, "error: unknown command %q\n\n", cmd)
		printUsage(stderr)
		return 2
	}
}

func runPseudo(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch strings.ToLower(arg) {
			case "-h", "-help", "--help":
				printUsage(stdout)
				return 0
			default:
				fmt.Fprintf(stderr, "error: unknown flag %q\n", arg)
				return 2
			}
		}
		positional = append(positional, arg)
	}

	if len(positional) == 0 {
		fmt.Fprintln(stderr, "error: input file is required (usage: pseugo run <file.pseudo>)")
		return 2
	}
	if len(positional) > 1 {
		fmt.Fprintln(stderr, "error: too many arguments (usage: pseugo run <file.pseudo>)")
		return 2
	}

	input := positional[0]
	sourceBytes, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", input, err)
		return 1
	}

	result, err := compiler.Compile(string(sourceBytes))
	if err != nil {
		fmt.Fprintf(stderr, "compile error: %v\n", err)
		return 1
	}

	tmpDir, err := os.MkdirTemp("", "pseugo-run-*")
	if err != nil {
		fmt.Fprintf(stderr, "error creating temp directory: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte(result.GoSource), 0644); err != nil {
		fmt.Fprintf(stderr, "error writing temporary Go file: %v\n", err)
		return 1
	}

	cmd := exec.Command("go", "run", goFile)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "error running generated Go program: %v\n", err)
		return 1
	}
	return 0
}

func runConvert(args []string, stdout, stderr io.Writer) int {
	var toGolang bool
	var toPseu bool
	var positional []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch strings.ToLower(arg) {
			case "-to-golang", "--to-golang", "-to-go", "--to-go":
				toGolang = true
			case "-to-pseu", "--to-pseu", "-to-pseudo", "--to-pseudo":
				toPseu = true
			case "-h", "-help", "--help":
				printUsage(stdout)
				return 0
			default:
				fmt.Fprintf(stderr, "error: unknown flag %q\n", arg)
				return 2
			}
			continue
		}
		positional = append(positional, arg)
	}

	if len(positional) == 0 {
		fmt.Fprintln(stderr, "error: input file is required (usage: pseugo convert [-to-golang|-to-pseu] <input> [output])")
		return 2
	}
	if len(positional) > 2 {
		fmt.Fprintln(stderr, "error: too many arguments (usage: pseugo convert [-to-golang|-to-pseu] <input> [output])")
		return 2
	}

	input := positional[0]
	output := ""
	if len(positional) == 2 {
		output = positional[1]
	}

	if toGolang && toPseu {
		fmt.Fprintln(stderr, "error: cannot specify both -to-golang and -to-pseu")
		return 2
	}

	// Auto-detect direction if neither flag is set
	if !toGolang && !toPseu {
		ext := strings.ToLower(filepath.Ext(input))
		switch ext {
		case ".go":
			toPseu = true
		case ".pseudo", ".pseudocode":
			toGolang = true
		default:
			fmt.Fprintln(stderr, "error: please specify conversion mode (-to-golang or -to-pseu)")
			return 2
		}
	}

	sourceBytes, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", input, err)
		return 1
	}
	source := string(sourceBytes)

	if toPseu {
		// Go -> Pseudocode
		pseudoSource, err := compiler.DecompileGo(source)
		if err != nil {
			fmt.Fprintf(stderr, "decompile error: %v\n", err)
			return 1
		}

		if output == "" {
			fmt.Fprintln(stdout, pseudoSource)
			return 0
		}

		if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			fmt.Fprintf(stderr, "error creating directory: %v\n", err)
			return 1
		}
		if err = os.WriteFile(output, []byte(pseudoSource+"\n"), 0644); err != nil {
			fmt.Fprintf(stderr, "error writing output: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "Generated pseudocode: %s\n", output)
		return 0
	}

	// Pseudocode -> Go (runs typechecker before codegen)
	result, err := compiler.Compile(source)
	if err != nil {
		fmt.Fprintf(stderr, "compile error: %v\n", err)
		return 1
	}

	if output == "" {
		fmt.Fprint(stdout, result.GoSource)
		return 0
	}

	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		fmt.Fprintf(stderr, "error creating directory: %v\n", err)
		return 1
	}
	if err = os.WriteFile(output, []byte(result.GoSource), 0644); err != nil {
		fmt.Fprintf(stderr, "error writing output: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Generated Go: %s\n", output)
	return 0
}

func runCheck(args []string, stdout, stderr io.Writer) int {
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch strings.ToLower(arg) {
			case "-h", "-help", "--help":
				printUsage(stdout)
				return 0
			default:
				fmt.Fprintf(stderr, "error: unknown flag %q\n", arg)
				return 2
			}
		}
		positional = append(positional, arg)
	}

	if len(positional) == 0 {
		fmt.Fprintln(stderr, "error: file to check is required (usage: pseugo check <file.pseudo>)")
		return 2
	}
	if len(positional) > 1 {
		fmt.Fprintln(stderr, "error: too many arguments (usage: pseugo check <file.pseudo>)")
		return 2
	}

	input := positional[0]
	sourceBytes, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", input, err)
		return 1
	}

	programs, err := compiler.CheckSource(string(sourceBytes))
	if err != nil {
		fmt.Fprintf(stderr, "typecheck error:\n%v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "[OK] Type check passed: %s is valid (%d program(s) verified, no errors)\n", input, len(programs))
	return 0
}

func runTokens(args []string, stdout, stderr io.Writer) int {
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch strings.ToLower(arg) {
			case "-h", "-help", "--help":
				printUsage(stdout)
				return 0
			default:
				fmt.Fprintf(stderr, "error: unknown flag %q\n", arg)
				return 2
			}
		}
		positional = append(positional, arg)
	}

	if len(positional) == 0 {
		fmt.Fprintln(stderr, "error: input file is required (usage: pseugo tokens <file.pseudo>)")
		return 2
	}
	if len(positional) > 1 {
		fmt.Fprintln(stderr, "error: too many arguments (usage: pseugo tokens <file.pseudo>)")
		return 2
	}

	input := positional[0]
	sourceBytes, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", input, err)
		return 1
	}

	tokens, err := compiler.NewLexer(string(sourceBytes)).AllTokens()
	if err != nil {
		fmt.Fprintf(stderr, "lexer error: %v\n", err)
		return 1
	}

	for _, t := range tokens {
		fmt.Fprintln(stdout, t.String())
	}
	return 0
}

func runAuto(file string, extraArgs []string, stdout, stderr io.Writer) int {
	ext := strings.ToLower(filepath.Ext(file))
	if ext == ".go" {
		return runConvert(append([]string{"-to-pseu", file}, extraArgs...), stdout, stderr)
	}
	return runConvert(append([]string{"-to-golang", file}, extraArgs...), stdout, stderr)
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "pseugo - Pseudocode Compiler & Go Transpiler (Alpro Telkom University)")
	fmt.Fprintln(w, "\nUsage:")
	fmt.Fprintln(w, "  pseugo run <file.pseudo>")
	fmt.Fprintln(w, "  pseugo convert [-to-golang|-to-pseu] <input> [output]")
	fmt.Fprintln(w, "  pseugo <input> [output]")
	fmt.Fprintln(w, "  pseugo check <file.pseudo>")
	fmt.Fprintln(w, "  pseugo tokens <file.pseudo>")
	fmt.Fprintln(w, "\nCommands:")
	fmt.Fprintln(w, "  run        Compile and execute a pseudocode file directly")
	fmt.Fprintln(w, "  convert    Convert between Pseudocode and Golang (direction auto-detected from extension)")
	fmt.Fprintln(w, "  check      Run static typechecker on pseudocode without code generation")
	fmt.Fprintln(w, "  tokens     Print lexer tokens of pseudocode")
	fmt.Fprintln(w, "  help       Show this help message")
	fmt.Fprintln(w, "\nConvert Options:")
	fmt.Fprintln(w, "  -to-golang, -to-go     Convert pseudocode to Go")
	fmt.Fprintln(w, "  -to-pseu, -to-pseudo   Convert Go to pseudocode")
	fmt.Fprintln(w, "\nExamples:")
	fmt.Fprintln(w, "  pseugo run examples/div_mod.pseudo")
	fmt.Fprintln(w, "  pseugo convert examples/div_mod.pseudo build/div_mod.go")
	fmt.Fprintln(w, "  pseugo convert -to-golang examples/div_mod.pseudo build/div_mod.go")
	fmt.Fprintln(w, "  pseugo convert -to-pseu build/div_mod.go div_mod.pseudo")
	fmt.Fprintln(w, "  pseugo build/div_mod.go")
	fmt.Fprintln(w, "  pseugo check examples/div_mod.pseudo")
	fmt.Fprintln(w, "  pseugo tokens examples/div_mod.pseudo")
}
