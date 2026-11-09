package test

import (
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestCompileDependOn(t *testing.T) {
	source := `program DemoSwitch
kamus
    indeks : char
    batasA, batasB : integer
algoritma
    input(indeks)
    depend on indeks
        'A':
            batasA <- 100
            batasB <- 75
        'B':
            batasA <- 75
            batasB <- 65
        case 'C':
            batasA <- 65
            batasB <- 50
        default:
            batasA <- 50
            batasB <- 0
    enddepend
    output(batasA, batasB)
endprogram`

	res, err := compiler.Compile(source)
	if err != nil {
		t.Fatalf("expected depend on to compile successfully, got error: %v", err)
	}

	if !strings.Contains(res.GoSource, "switch indeks {") {
		t.Errorf("expected switch statement in generated Go, got:\n%s", res.GoSource)
	}
	if !strings.Contains(res.GoSource, "case 'A':") || !strings.Contains(res.GoSource, "default:") {
		t.Errorf("expected cases in generated Go, got:\n%s", res.GoSource)
	}
}

func TestDecompileGoSwitchToDependOn(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	var nama string
	switch n {
	case 1:
		nama = "Satu"
	case 2, 3:
		nama = "Dua atau Tiga"
	default:
		nama = "Lainnya"
	}
	fmt.Println(nama)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed on switch: %v", err)
	}

	if !strings.Contains(pseudo, "depend on n") {
		t.Errorf("expected depend on n, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "1:") || !strings.Contains(pseudo, "2, 3:") {
		t.Errorf("expected case labels in pseudocode, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "default:") || !strings.Contains(pseudo, "enddepend") {
		t.Errorf("expected default and enddepend, got:\n%s", pseudo)
	}

	// Verify round-trip compilation
	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("re-compiling decompiled switch failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileTaglessSwitchToIfElse(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var skor int
	fmt.Scan(&skor)
	var indeks string
	switch {
	case skor >= 80:
		indeks = "A"
	case skor >= 70:
		indeks = "B"
	default:
		indeks = "C"
	}
	fmt.Println(indeks)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed on tagless switch: %v", err)
	}

	if !strings.Contains(pseudo, "if skor >= 80 then") {
		t.Errorf("expected if skor >= 80 then, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "elseif skor >= 70 then") {
		t.Errorf("expected elseif skor >= 70 then, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "else") || !strings.Contains(pseudo, "endif") {
		t.Errorf("expected else and endif, got:\n%s", pseudo)
	}

	// Verify round-trip compilation
	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("re-compiling decompiled tagless switch failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestCheckSourceRejectsMismatchedCaseType(t *testing.T) {
	source := `program BadSwitch
kamus
    x : integer
algoritma
    depend on x
        "teks":
            output("Error")
    enddepend
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("expected case type mismatch error, got: %v", err)
	}
}
