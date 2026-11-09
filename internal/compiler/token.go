package compiler

import "fmt"

type TokenKind string

const (
	TokenEOF     TokenKind = "EOF"
	TokenNewline TokenKind = "NEWLINE"
	TokenIdent   TokenKind = "IDENT"
	TokenInt     TokenKind = "INT_LITERAL"
	TokenReal    TokenKind = "REAL_LITERAL"
	TokenString  TokenKind = "STRING_LITERAL"
	TokenChar    TokenKind = "CHAR_LITERAL"
	TokenSymbol  TokenKind = "SYMBOL"
	TokenOp      TokenKind = "OPERATOR"
)

type Position struct {
	Offset int
	Line   int
	Column int
}

type Token struct {
	Kind TokenKind
	Text string
	Pos  Position
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%d:%d", t.Kind, t.Text, t.Pos.Line, t.Pos.Column)
}

type Diagnostic struct {
	Pos     Position
	Message string
}

func (d Diagnostic) Error() string {
	return fmt.Sprintf("%d:%d: %s", d.Pos.Line, d.Pos.Column, d.Message)
}
