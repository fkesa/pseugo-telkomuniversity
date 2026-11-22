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
    end
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

	if !strings.Contains(pseudo, "    depend on n\n    1:\n        nama <- \"Satu\"") {
		t.Errorf("expected case labels aligned with depend on and body indented once, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "\n    2, 3:\n") {
		t.Errorf("expected multi-value case label aligned with depend on, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "\n    default:\n        nama <- \"Lainnya\"\n    end\n") || strings.Contains(pseudo, "enddepend") {
		t.Errorf("expected default and end aligned with depend on (not enddepend), got:\n%s", pseudo)
	}

	// Verify round-trip compilation
	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("re-compiling decompiled switch failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileTaglessSwitchToDependOn(t *testing.T) {
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

	if strings.Contains(pseudo, "if ") || strings.Contains(pseudo, "elseif ") {
		t.Errorf("expected tagless switch to decompile to depend on instead of if/elseif, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "depend on\n") {
		t.Errorf("expected tagless depend on, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "skor >= 80:") || !strings.Contains(pseudo, "skor >= 70:") {
		t.Errorf("expected boolean condition cases, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "default:") || !strings.Contains(pseudo, "\n    end\n") || strings.Contains(pseudo, "enddepend") {
		t.Errorf("expected default and end (not enddepend), got:\n%s", pseudo)
	}

	// Verify round-trip compilation back to Go tagless switch
	res, err := compiler.Compile(pseudo)
	if err != nil {
		t.Fatalf("re-compiling decompiled tagless switch failed: %v\nPseudocode:\n%s", err, pseudo)
	}
	if !strings.Contains(res.GoSource, "switch {") {
		t.Errorf("expected round-tripped Go to contain tagless switch, got:\n%s", res.GoSource)
	}
}

func TestDecompileTaglessEqualitySwitchExtractsSubject(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var jenis, kendaraan string
	fmt.Scan(&jenis)
	switch {
	case jenis == "A":
		kendaraan = "Mobil"
	case jenis == "B", jenis == "b":
		kendaraan = "Motor"
	default:
		kendaraan = "Tidak tersedia"
	}
	fmt.Println(kendaraan)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed on equality tagless switch: %v", err)
	}

	if !strings.Contains(pseudo, "depend on jenis") {
		t.Errorf("expected depend on jenis, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, `"A":`) || !strings.Contains(pseudo, `"B", "b":`) {
		t.Errorf("expected extracted case values, got:\n%s", pseudo)
	}
	if strings.Contains(pseudo, "if ") {
		t.Errorf("did not expect if/else for switch, got:\n%s", pseudo)
	}

	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("re-compiling extracted depend on failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestCheckSourceRejectsNonBooleanTaglessDependOnCase(t *testing.T) {
	source := `program BadTaglessSwitch
kamus
    x : integer
algoritma
    depend on
        x + 1:
            output("Error")
    end
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "case condition must be boolean") {
		t.Fatalf("expected non-boolean case condition error, got: %v", err)
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
    end
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("expected case type mismatch error, got: %v", err)
	}
}

func TestParserRejectsEnddepend(t *testing.T) {
	source := `program RejectEnddepend
kamus
    x : integer
algoritma
    depend on x
        1:
            output(x)
    enddepend
endprogram`

	_, err := compiler.Compile(source)
	if err == nil {
		t.Fatalf("expected enddepend to be rejected, but compilation succeeded")
	}
}
