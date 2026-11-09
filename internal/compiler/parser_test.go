package compiler

import "testing"

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
	result, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Program.Body) != 2 {
		t.Fatalf("got %d top-level statements, want 2", len(result.Program.Body))
	}
	outer, ok := result.Program.Body[1].(ForStmt)
	if !ok {
		t.Fatalf("second statement has type %T, want ForStmt", result.Program.Body[1])
	}
	if len(outer.Body) != 1 {
		t.Fatalf("outer loop has %d body statements, want 1", len(outer.Body))
	}
	if _, ok := outer.Body[0].(ForStmt); !ok {
		t.Fatalf("outer loop child has type %T, want nested ForStmt", outer.Body[0])
	}
}

func TestParserRejectsUnsupportedIf(t *testing.T) {
	source := `program Unsupported
kamus
  x : integer
algoritma
  if x = 1 then
    output(x)
  endif
endprogram
`
	_, err := Compile(source)
	if err == nil {
		t.Fatal("expected parser error because IF is not part of the implemented grammar")
	}
}
