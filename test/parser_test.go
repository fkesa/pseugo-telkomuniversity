package test

import (
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestParserBuildsNestedForLoops(t *testing.T) {
	source := `program Nested
kamus
  total : integer
algoritma
  total <- 0
  for i <- 1 to 2 do
    for j <- 1 to 2 do
      total <- total + i * j
    endfor
  endfor
endprogram
`
	result, err := compiler.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Program.Body) != 2 {
		t.Fatalf("got %d top-level statements, want 2", len(result.Program.Body))
	}
	outer, ok := result.Program.Body[1].(compiler.ForStmt)
	if !ok {
		t.Fatalf("second statement has type %T, want ForStmt", result.Program.Body[1])
	}
	if len(outer.Body) != 1 {
		t.Fatalf("outer loop has %d body statements, want 1", len(outer.Body))
	}
	if _, ok := outer.Body[0].(compiler.ForStmt); !ok {
		t.Fatalf("outer loop child has type %T, want nested ForStmt", outer.Body[0])
	}
}

func TestParserSupportsIf(t *testing.T) {
	source := `program Supported
kamus
  x : integer
algoritma
  if x = 1 then
    output(x)
  endif
endprogram
`
	_, err := compiler.Compile(source)
	if err != nil {
		t.Fatalf("expected compile to succeed with IF, got: %v", err)
	}
}
