package test

import (
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestElseIfWithSpaceSupport(t *testing.T) {
	src := `program TestElseIf
kamus
    a, b : integer
algoritma
    input(a, b)
    if a < b then
        output("kurang")
    else if a > b then
        output("lebih")
    else
        output("sama")
    endif
endprogram`

	res, err := compiler.Compile(src)
	if err != nil {
		t.Fatalf("expected 'else if' with space to compile, got: %v", err)
	}
	if !strings.Contains(res.GoSource, "else if") {
		t.Errorf("expected generated Go to contain 'else if', got:\n%s", res.GoSource)
	}
}

func TestAmpersandAndPipeOperatorsInPseudocode(t *testing.T) {
	src := `program TestLogicalSymbols
kamus
    n : integer
    valid : boolean
algoritma
    input(n)
    if n >= 1 && n <= 100 then
        valid <- true
    elseif n < 0 || n > 200 then
        valid <- false
    endif
    output(valid)
endprogram`

	res, err := compiler.Compile(src)
	if err != nil {
		t.Fatalf("expected '&&' and '||' in pseudocode to compile, got: %v", err)
	}
	if !strings.Contains(res.GoSource, "&&") || !strings.Contains(res.GoSource, "||") {
		t.Errorf("expected logical operators in generated Go, got:\n%s", res.GoSource)
	}
}

func TestProgramNameMatchingOperatorKeyword(t *testing.T) {
	src := `program AND
kamus
    n : integer
    hasil : boolean
algoritma
    input(n)
    hasil <- true
    output(hasil)
endprogram`

	res, err := compiler.Compile(src)
	if err != nil {
		t.Fatalf("expected program named AND to compile, got: %v", err)
	}
	if !strings.Contains(res.GoSource, "func AND()") {
		t.Errorf("expected func AND() in generated Go, got:\n%s", res.GoSource)
	}
}

func TestMixedTypeArithmeticAutoPromotion(t *testing.T) {
	src := `program TestMixedArithmetic
kamus
    harga : real
    qty : integer
    total : real
algoritma
    harga <- 12.5
    qty <- 4
    total <- harga * qty
    output(total)
endprogram`

	res, err := compiler.Compile(src)
	if err != nil {
		t.Fatalf("expected mixed arithmetic to compile, got: %v", err)
	}
	if !strings.Contains(res.GoSource, "float64(qty)") {
		t.Errorf("expected float64 promotion for int operand, got:\n%s", res.GoSource)
	}
}

func TestDecompileStructFieldRealDivision(t *testing.T) {
	goSrc := `package main

import "fmt"

type d_bmi struct {
	nama      string
	berat     float64
	tinggi    float64
	nilai_bmi float64
}

func bmi() {
	var b d_bmi
	fmt.Scan(&b.nama, &b.berat, &b.tinggi)
	b.nilai_bmi = b.berat / (b.tinggi * b.tinggi)
	fmt.Println(b.nilai_bmi)
}

func main() {
	bmi()
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("expected DecompileGo to succeed on struct real division, got: %v", err)
	}

	if strings.Contains(pseudo, "DIV") {
		t.Errorf("expected real division '/' instead of 'DIV', got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "b.nilai_bmi <- b.berat / (b.tinggi * b.tinggi)") {
		t.Errorf("expected b.nilai_bmi <- b.berat / ..., got:\n%s", pseudo)
	}

	// Verify round trip
	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("recompilation of decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}
