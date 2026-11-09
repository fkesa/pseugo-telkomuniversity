package test

import (
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestVariableAliasSharesStorage(t *testing.T) {
	source := `program Alias
kamus
  original : integer
  alias : original
algoritma
  alias <- 7
  output(original)
endprogram
`
	result, err := compiler.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result.GoSource, "var alias") {
		t.Fatalf("alias should not allocate a second variable:\n%s", result.GoSource)
	}
	if !strings.Contains(result.GoSource, "original = 7") {
		t.Fatalf("assignment through alias did not target original variable:\n%s", result.GoSource)
	}
}
