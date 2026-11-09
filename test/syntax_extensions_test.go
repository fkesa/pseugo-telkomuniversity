package test

import (
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestCommentsSupport(t *testing.T) {
	source := `program TestComments
{ ini adalah komentar gaya Pascal/Algoritmik }
kamus
    // komentar satu baris
    a : integer
    /* komentar
       blok bergaya C */
    b : real
algoritma
    a <- 10 { inisialisasi a }
    b <- 2.5 // inisialisasi b
    output(a, b)
endprogram`

	res, err := compiler.Compile(source)
	if err != nil {
		t.Fatalf("expected program with comments to compile, got: %v", err)
	}

	if !strings.Contains(res.GoSource, "a = 10") || !strings.Contains(res.GoSource, "b = 2.5") {
		t.Errorf("expected assignments in generated Go, got:\n%s", res.GoSource)
	}
}

func TestEqualSignAssignment(t *testing.T) {
	source := `program TestEqualAssign
kamus
    constant batas : integer = 100
    x, y, i : integer
algoritma
    x = 10
    y = x + batas
    for i = 1 to 5 do
        y = y + i
    endfor
    output(y)
endprogram`

	res, err := compiler.Compile(source)
	if err != nil {
		t.Fatalf("expected assignment with '=' to compile, got: %v", err)
	}

	if !strings.Contains(res.GoSource, "x = 10") || !strings.Contains(res.GoSource, "batas int = 100") {
		t.Errorf("expected assignments in generated Go, got:\n%s", res.GoSource)
	}
}

func TestIncDecOperators(t *testing.T) {
	source := `program TestIncDec
kamus
    c : integer
algoritma
    c <- 0
    c++
    c++
    c--
    output(c)
endprogram`

	res, err := compiler.Compile(source)
	if err != nil {
		t.Fatalf("expected ++ and -- to compile, got: %v", err)
	}

	if !strings.Contains(res.GoSource, "c = (c + 1)") || !strings.Contains(res.GoSource, "c = (c - 1)") {
		t.Errorf("expected inc/dec in generated Go, got:\n%s", res.GoSource)
	}
}

func TestTypeCasting(t *testing.T) {
	source := `program TestCast
kamus
    pi : real
    nilai : integer
algoritma
    pi <- 3.14
    nilai <- int(pi)
    output(nilai)
endprogram`

	res, err := compiler.Compile(source)
	if err != nil {
		t.Fatalf("expected type casting to compile, got: %v", err)
	}

	if !strings.Contains(res.GoSource, "int(pi)") {
		t.Errorf("expected int(pi) in generated Go, got:\n%s", res.GoSource)
	}
}
