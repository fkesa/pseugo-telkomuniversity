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

	// 1. Pseudocode -> Go
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"convert", "-to-golang", "-i", "../examples/div_mod.pseudo", "-o", goOut}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}

	// 2. Go -> Pseudocode
	stdout.Reset()
	stderr.Reset()
	code = cli.Run([]string{"convert", "-to-pseu", "-i", goOut, "-o", pseudoOut}, &stdout, &stderr)
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

func TestCLILegacy(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"-check", "-in", "../examples/div_mod.pseudo"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr: %s", code, stderr.String())
	}
}
