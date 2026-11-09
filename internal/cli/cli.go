package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"pseudocode-compiler/internal/compiler"
)

// Run executes the CLI with the provided arguments and standard I/O streams.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	cmd := args[0]

	// Handle legacy flags if first arg starts with "-"
	if strings.HasPrefix(cmd, "-") {
		return runLegacy(args, stdout, stderr)
	}

	switch strings.ToLower(cmd) {
	case "convert":
		return runConvert(args[1:], stdout, stderr)
	case "check":
		return runCheck(args[1:], stdout, stderr)
	case "tokens":
		return runTokens(args[1:], stdout, stderr)
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		// Check if user passed a file directly as first argument: pseugo <file>
		if _, err := os.Stat(cmd); err == nil {
			return runAuto(cmd, args[1:], stdout, stderr)
		}
		fmt.Fprintf(stderr, "error: unknown command %q\n\n", cmd)
		printUsage(stderr)
		return 2
	}
}

func runConvert(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var toGolang bool
	var toPseu bool
	var input string
	var output string

	fs.BoolVar(&toGolang, "to-golang", false, "convert Pseudocode to Go")
	fs.BoolVar(&toGolang, "to-go", false, "alias for -to-golang")
	fs.BoolVar(&toPseu, "to-pseu", false, "convert Go to Pseudocode")
	fs.BoolVar(&toPseu, "to-pseudo", false, "alias for -to-pseu")

	fs.StringVar(&input, "i", "", "input file path")
	fs.StringVar(&input, "in", "", "alias for -i")

	fs.StringVar(&output, "o", "", "output file path (default: stdout)")
	fs.StringVar(&output, "out", "", "alias for -o")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Positional arguments fallback
	if input == "" && len(fs.Args()) > 0 {
		input = fs.Args()[0]
	}
	if output == "" && len(fs.Args()) > 1 {
		output = fs.Args()[1]
	}

	if input == "" {
		fmt.Fprintln(stderr, "error: input file is required (use -i <file>)")
		return 2
	}

	if toGolang && toPseu {
		fmt.Fprintln(stderr, "error: cannot specify both -to-golang and -to-pseu")
		return 2
	}

	// Auto-detect direction if neither flag is set
	if !toGolang && !toPseu {
		ext := strings.ToLower(filepath.Ext(input))
		if ext == ".go" {
			toPseu = true
		} else if ext == ".pseudo" {
			toGolang = true
		} else {
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
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var input string
	fs.StringVar(&input, "i", "", "input file path")
	fs.StringVar(&input, "in", "", "alias for -i")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if input == "" && len(fs.Args()) > 0 {
		input = fs.Args()[0]
	}

	if input == "" {
		fmt.Fprintln(stderr, "error: file to check is required (e.g. pseugo check <file.pseudo>)")
		return 2
	}

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
	fs := flag.NewFlagSet("tokens", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var input string
	fs.StringVar(&input, "i", "", "input file path")
	fs.StringVar(&input, "in", "", "alias for -i")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if input == "" && len(fs.Args()) > 0 {
		input = fs.Args()[0]
	}

	if input == "" {
		fmt.Fprintln(stderr, "error: input file is required (e.g. pseugo tokens <file.pseudo>)")
		return 2
	}

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
		return runConvert(append([]string{"-to-pseu", "-i", file}, extraArgs...), stdout, stderr)
	}
	return runConvert(append([]string{"-to-golang", "-i", file}, extraArgs...), stdout, stderr)
}

func runLegacy(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("legacy", flag.ContinueOnError)
	fs.SetOutput(stderr)

	input := fs.String("in", "", "path to input source file (.pseudo or .go)")
	output := fs.String("out", "", "path for output file (default: stdout)")
	checkOnly := fs.Bool("check", false, "run typechecker on pseudocode")
	showTokens := fs.Bool("tokens", false, "print lexer tokens")
	reverse := fs.Bool("reverse", false, "convert Go source code to pseudocode")
	toPseudo := fs.Bool("to-pseudo", false, "alias for -reverse")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *input == "" {
		printUsage(stderr)
		return 2
	}

	if *showTokens {
		return runTokens([]string{"-i", *input}, stdout, stderr)
	}
	if *checkOnly {
		return runCheck([]string{"-i", *input}, stdout, stderr)
	}

	convertArgs := []string{"-i", *input}
	if *output != "" {
		convertArgs = append(convertArgs, "-o", *output)
	}
	if *reverse || *toPseudo {
		convertArgs = append(convertArgs, "-to-pseu")
	}
	return runConvert(convertArgs, stdout, stderr)
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "pseugo - Pseudocode Compiler & Go Transpiler (Alpro Telkom University)")
	fmt.Fprintln(w, "\nUsage:")
	fmt.Fprintln(w, "  pseugo convert -to-golang -i <input.pseudo> [-o <output.go>]")
	fmt.Fprintln(w, "  pseugo convert -to-pseu -i <input.go> [-o <output.pseudo>]")
	fmt.Fprintln(w, "  pseugo check <file.pseudo>")
	fmt.Fprintln(w, "  pseugo tokens <file.pseudo>")
	fmt.Fprintln(w, "\nCommands:")
	fmt.Fprintln(w, "  convert    Convert between Pseudocode and Golang")
	fmt.Fprintln(w, "  check      Run static typechecker on pseudocode without code generation")
	fmt.Fprintln(w, "  tokens     Print lexer tokens of pseudocode")
	fmt.Fprintln(w, "  help       Show this help message")
	fmt.Fprintln(w, "\nConvert Options:")
	fmt.Fprintln(w, "  -to-golang, -to-go     Convert pseudocode to Go")
	fmt.Fprintln(w, "  -to-pseu, -to-pseudo   Convert Go to pseudocode")
	fmt.Fprintln(w, "  -i, -in string         Input file path")
	fmt.Fprintln(w, "  -o, -out string        Output file path (default: stdout)")
	fmt.Fprintln(w, "\nExamples:")
	fmt.Fprintln(w, "  pseugo convert -to-golang -i examples/div_mod.pseudo -o build/div_mod.go")
	fmt.Fprintln(w, "  pseugo convert -to-pseu -i build/div_mod.go -o div_mod.pseudo")
	fmt.Fprintln(w, "  pseugo check examples/div_mod.pseudo")
	fmt.Fprintln(w, "  pseugo tokens examples/div_mod.pseudo")
}
