package test

import (
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestCheckSourceValid(t *testing.T) {
	source := `program ValidProgram
kamus
    a, b : integer
    x : real
    nama : string
    status : boolean
    constant batas : integer <- 50
    type Titik <
        x : real
        y : real
    >
    pt : Titik
algoritma
    a <- 10
    b <- a + batas
    x <- 3.14
    pt.x <- x
    pt.y <- 2.0
    status <- a < b
    if status then
        output("Valid:", b)
    endif
    for i <- 1 to 5 do
        a <- a + i
    endfor
    while a > 0 do
        a <- a - 1
    endwhile
    repeat
        b <- b + 1
    until b >= 100
endprogram`

	programs, err := compiler.CheckSource(source)
	if err != nil {
		t.Fatalf("expected valid source to pass typecheck, got error: %v", err)
	}
	if len(programs) != 1 {
		t.Fatalf("expected 1 program, got %d", len(programs))
	}
}

func TestCheckSourceRejectsUndeclaredVariable(t *testing.T) {
	source := `program Bad
kamus
    x : integer
algoritma
    x <- y + 10
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "not declared") {
		t.Fatalf("expected undeclared variable error, got %v", err)
	}
}

func TestCheckSourceRejectsNonBooleanIfCondition(t *testing.T) {
	source := `program Bad
kamus
    x : integer
algoritma
    if x + 5 then
        output(x)
    endif
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "if condition must be boolean") {
		t.Fatalf("expected boolean if condition error, got %v", err)
	}
}

func TestCheckSourceRejectsNonBooleanWhileCondition(t *testing.T) {
	source := `program Bad
kamus
    x : integer
algoritma
    while x + 1 do
        output(x)
    endwhile
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "while condition must be boolean") {
		t.Fatalf("expected boolean while condition error, got %v", err)
	}
}

func TestCheckSourceRejectsNonBooleanUntilCondition(t *testing.T) {
	source := `program Bad
kamus
    x : integer
algoritma
    repeat
        x <- x + 1
    until 100
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "until condition must be boolean") {
		t.Fatalf("expected boolean until condition error, got %v", err)
	}
}

func TestCheckSourceRejectsNonIntegerForBounds(t *testing.T) {
	source := `program Bad
kamus
    x : integer
algoritma
    for i <- 1.5 to 10.0 do
        x <- x + 1
    endfor
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "for loop bounds must both be integer") {
		t.Fatalf("expected integer for bounds error, got %v", err)
	}
}

func TestCheckSourceRejectsNonIntegerDIV(t *testing.T) {
	source := `program Bad
kamus
    a : real
    b, c : integer
algoritma
    c <- a DIV b
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "defined only for integer operands") {
		t.Fatalf("expected DIV integer operands error, got %v", err)
	}
}

func TestCheckSourceRejectsUnknownRecordField(t *testing.T) {
	source := `program Bad
kamus
    type Titik <
        x : real
        y : real
    >
    pt : Titik
algoritma
    pt.z <- 10.0
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "has no field") {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func TestCheckSourceRejectsConstantInput(t *testing.T) {
	source := `program Bad
kamus
    constant batas : integer <- 10
algoritma
    input(batas)
endprogram`

	_, err := compiler.CheckSource(source)
	if err == nil || !strings.Contains(err.Error(), "cannot input into constant") {
		t.Fatalf("expected constant input error, got %v", err)
	}
}
