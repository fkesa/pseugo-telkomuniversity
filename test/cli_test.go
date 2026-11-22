package test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"pseudocode-compiler/internal/cli"
)

func TestCLICheck(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"check", "../examples/div_mod.pseudo"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[OK] Type check passed") {
		t.Errorf("unexpected output: %s", stdout.String())
	}
}

func TestCLIConvertTwoWay(t *testing.T) {
	tmpDir := t.TempDir()
	goOut := filepath.Join(tmpDir, "cli_test_div.go")
	pseudoOut := filepath.Join(tmpDir, "cli_test_div.pseudo")

	// 1. Pseudocode -> Go (positional <input> <output>)
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"convert", "-to-golang", "../examples/div_mod.pseudo", goOut}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}

	// 2. Go -> Pseudocode (positional <input> <output>)
	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"convert", "-to-pseu", goOut, pseudoOut}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}

	// 3. Check the generated pseudocode
	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"check", pseudoOut}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected generated pseudocode to pass check, got code %d, stderr: %s", code, stderr.String())
	}
}

func TestCLIConvertPositionalAutoDetect(t *testing.T) {
	tmpDir := t.TempDir()
	goOut := filepath.Join(tmpDir, "auto_div.go")
	pseudoOut := filepath.Join(tmpDir, "auto_div.pseudo")

	var stdout, stderr bytes.Buffer
	// Without -to-golang flag (auto-detected from .pseudo)
	code := cli.Run([]string{"convert", "../examples/div_mod.pseudo", goOut}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}

	// Direct shorthand: pseugo <input.go> <output.pseudo>
	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{goOut, pseudoOut}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}

	// Flag after positional args
	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"convert", goOut, "-to-pseu"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "program ") {
		t.Errorf("expected pseudocode on stdout, got:\n%s", stdout.String())
	}
}

func TestCLITokens(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"tokens", "../examples/hello.pseudo"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "IDENT(\"program\")") && !strings.Contains(stdout.String(), "program") {
		t.Errorf("expected token output, got: %s", stdout.String())
	}
}

func TestCLIRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"run", "../examples/div_mod.pseudo"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "12 4" {
		t.Errorf("expected output '12 4', got: %q", stdout.String())
	}

	// Interactive input test with RunWithIO
	stdout.Reset()
	stderr.Reset()
	code = cli.RunWithIO([]string{"run", "../examples/hello.pseudo"}, strings.NewReader("Budi\n20\n"), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Halo Budi") || !strings.Contains(stdout.String(), "Umur: 20") {
		t.Errorf("unexpected interactive output: %q", stdout.String())
	}
}

func TestCLIRejectsLegacyAndIOFlags(t *testing.T) {
	cases := [][]string{
		{"-check", "-in", "../examples/div_mod.pseudo"},
		{"convert", "-to-golang", "-i", "../examples/div_mod.pseudo"},
		{"convert", "-to-golang", "../examples/div_mod.pseudo", "-o", "out.go"},
		{"check", "-i", "../examples/div_mod.pseudo"},
		{"tokens", "-i", "../examples/hello.pseudo"},
		{"run", "-i", "../examples/div_mod.pseudo"},
	}

	for _, args := range cases {
		var stdout, stderr bytes.Buffer
		code := cli.Run(args, &stdout, &stderr)
		if code == 0 {
			t.Errorf("expected non-zero exit code for %v, got 0", args)
		}
	}
}
