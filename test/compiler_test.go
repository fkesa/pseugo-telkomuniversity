package test

import (
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestCompileDivMod(t *testing.T) {
	source := `program Test
kamus
  a, b, q, r : integer
algoritma
  a <- 100
  b <- 8
  q <- a DIV b
  r <- a MOD b
  output(q, r)
endprogram
`
	result, err := compiler.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.GoSource, "(a / b)") || !strings.Contains(result.GoSource, "(a % b)") {
		t.Fatalf("generated Go did not map DIV/MOD correctly:\n%s", result.GoSource)
	}
}

func TestCompileRejectsTypeMismatch(t *testing.T) {
	source := `program Bad
kamus
  age : integer
algoritma
  age <- "teks"
endprogram
`
	_, err := compiler.Compile(source)
	if err == nil || !strings.Contains(err.Error(), "cannot assign") {
		t.Fatalf("expected type mismatch, got %v", err)
	}
}

func TestCompileRejectsNonIntegerMOD(t *testing.T) {
	source := `program Bad
kamus
  a : real
  b : real
  c : real
algoritma
  c <- a MOD b
endprogram
`
	_, err := compiler.Compile(source)
	if err == nil || !strings.Contains(err.Error(), "only for integer") {
		t.Fatalf("expected MOD type error, got %v", err)
	}
}

func TestCompileRejectsAssignmentToConstant(t *testing.T) {
	source := `program Bad
kamus
  constant limit : integer <- 9
algoritma
  limit <- 10
endprogram
`
	_, err := compiler.Compile(source)
	if err == nil || !strings.Contains(err.Error(), "cannot assign to constant") {
		t.Fatalf("expected constant assignment error, got %v", err)
	}
}

func TestCompileSupportsRealDivision(t *testing.T) {
	source := `program Bad
kamus
  a, b : integer
  c : real
algoritma
  c <- a / b
endprogram
`
	if _, err := compiler.Compile(source); err != nil {
		t.Fatalf("expected real division to compile: %v", err)
	}
}

func TestConstantCanUseEarlierConstant(t *testing.T) {
	source := `program Constants
kamus
  constant first : integer <- 10
  constant second : integer <- first + 5
algoritma
  output(second)
endprogram
`
	if _, err := compiler.Compile(source); err != nil {
		t.Fatalf("expected earlier-constant reference to compile: %v", err)
	}
}

func TestConstantCannotUseLaterConstant(t *testing.T) {
	source := `program Constants
kamus
  constant first : integer <- second + 1
  constant second : integer <- 5
algoritma
  output(first)
endprogram
`
	_, err := compiler.Compile(source)
	if err == nil || !strings.Contains(err.Error(), "previously declared constants") {
		t.Fatalf("expected forward constant reference error, got %v", err)
	}
}
