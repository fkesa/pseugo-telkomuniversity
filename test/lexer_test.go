package test

import (
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestLexerRecognizesDivModAndAssignment(t *testing.T) {
	tokens, err := compiler.NewLexer("hasil <- 100 DIV 8 MOD 3").AllTokens()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"hasil", "<-", "100", "DIV", "8", "MOD", "3"}
	if len(tokens) != len(want)+1 {
		t.Fatalf("got %d tokens, want %d", len(tokens), len(want)+1)
	}
	for i, w := range want {
		if tokens[i].Text != w {
			t.Fatalf("token %d = %q, want %q", i, tokens[i].Text, w)
		}
	}
}

func TestLexerReportsInvalidCharacter(t *testing.T) {
	_, err := compiler.NewLexer("nilai <- 3 @ 4").AllTokens()
	if err == nil {
		t.Fatal("expected lexer error")
	}
}

func TestLexerReportsUnterminatedString(t *testing.T) {
	_, err := compiler.NewLexer("output(\"belum selesai)").AllTokens()
	if err == nil {
		t.Fatal("expected unterminated string error")
	}
}

func TestLexerReadsQuotedChar(t *testing.T) {
	tokens, err := compiler.NewLexer("'A'").AllTokens()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Kind != compiler.TokenChar {
		t.Fatalf("got %s, want char literal", tokens[0].Kind)
	}
}

func TestLexerRecognizesNotEqual(t *testing.T) {
	tokens, err := compiler.NewLexer("a != b").AllTokens()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[1].Text != "!=" {
		t.Fatalf("got %q, want !=", tokens[1].Text)
	}
}
