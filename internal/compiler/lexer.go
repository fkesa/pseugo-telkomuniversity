package compiler

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	source []rune
	index  int
	line   int
	column int
}

func NewLexer(source string) *Lexer {
	return &Lexer{source: []rune(source), line: 1, column: 1}
}

func (l *Lexer) current() rune {
	if l.index >= len(l.source) {
		return 0
	}
	return l.source[l.index]
}

func (l *Lexer) peek(n int) rune {
	p := l.index + n
	if p < 0 || p >= len(l.source) {
		return 0
	}
	return l.source[p]
}

func (l *Lexer) position() Position {
	return Position{Offset: l.index, Line: l.line, Column: l.column}
}

func (l *Lexer) advance() rune {
	ch := l.current()
	if ch == 0 {
		return 0
	}
	l.index++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) skipHorizontalSpace() {
	for l.current() == ' ' || l.current() == '\t' || l.current() == '\r' {
		l.advance()
	}
}

func isIdentifierStart(ch rune) bool { return ch == '_' || unicode.IsLetter(ch) }
func isIdentifierPart(ch rune) bool  { return ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) }

func (l *Lexer) readIdentifier() Token {
	pos := l.position()
	start := l.index
	for isIdentifierPart(l.current()) {
		l.advance()
	}
	word := string(l.source[start:l.index])
	upper := strings.ToUpper(word)
	if upper == "DIV" || upper == "MOD" || upper == "AND" || upper == "OR" || upper == "XOR" || upper == "NOT" {
		return Token{Kind: TokenOp, Text: upper, Pos: pos}
	}
	return Token{Kind: TokenIdent, Text: word, Pos: pos}
}

func (l *Lexer) readNumber() (Token, error) {
	pos := l.position()
	start := l.index
	for unicode.IsDigit(l.current()) {
		l.advance()
	}
	kind := TokenInt
	if l.current() == '.' && unicode.IsDigit(l.peek(1)) {
		kind = TokenReal
		l.advance()
		for unicode.IsDigit(l.current()) {
			l.advance()
		}
	}
	text := string(l.source[start:l.index])
	if kind == TokenInt {
		if _, err := strconv.ParseInt(text, 10, 64); err != nil {
			return Token{}, Diagnostic{Pos: pos, Message: fmt.Sprintf("integer literal %q is out of supported range", text)}
		}
	} else if _, err := strconv.ParseFloat(text, 64); err != nil {
		return Token{}, Diagnostic{Pos: pos, Message: fmt.Sprintf("invalid real literal %q", text)}
	}
	return Token{Kind: kind, Text: text, Pos: pos}, nil
}

func (l *Lexer) readQuoted(quote rune, kind TokenKind) (Token, error) {
	pos := l.position()
	start := l.index
	l.advance() // opening quote
	escaped := false
	for l.current() != 0 {
		ch := l.current()
		if ch == '\n' {
			return Token{}, Diagnostic{Pos: pos, Message: "newline inside quoted literal"}
		}
		if escaped {
			escaped = false
			l.advance()
			continue
		}
		if ch == '\\' {
			escaped = true
			l.advance()
			continue
		}
		if ch == quote {
			l.advance()
			text := string(l.source[start:l.index])
			if kind == TokenString {
				if _, err := strconv.Unquote(text); err != nil {
					return Token{}, Diagnostic{Pos: pos, Message: "invalid string escape: " + err.Error()}
				}
			} else {
				r, err := strconv.Unquote(text)
				if err != nil {
					return Token{}, Diagnostic{Pos: pos, Message: "invalid character literal: " + err.Error()}
				}
				if len([]rune(r)) != 1 {
					return Token{}, Diagnostic{Pos: pos, Message: "character literal must contain exactly one character"}
				}
			}
			return Token{Kind: kind, Text: text, Pos: pos}, nil
		}
		l.advance()
	}
	return Token{}, Diagnostic{Pos: pos, Message: "unterminated quoted literal"}
}

func (l *Lexer) NextToken() (Token, error) {
	for {
		l.skipHorizontalSpace()
		pos := l.position()
		ch := l.current()
		switch {
		case ch == 0:
			return Token{Kind: TokenEOF, Pos: pos}, nil
		case ch == '{':
			l.advance()
			for l.current() != 0 && l.current() != '}' {
				l.advance()
			}
			if l.current() == '}' {
				l.advance()
			}
			continue
		case ch == '/' && l.peek(1) == '/':
			l.advance()
			l.advance()
			for l.current() != 0 && l.current() != '\n' {
				l.advance()
			}
			continue
		case ch == '/' && l.peek(1) == '*':
			l.advance()
			l.advance()
			for l.current() != 0 && !(l.current() == '*' && l.peek(1) == '/') {
				l.advance()
			}
			if l.current() == '*' {
				l.advance()
				l.advance()
			}
			continue
		case ch == '\n':
			l.advance()
			return Token{Kind: TokenNewline, Text: "\\n", Pos: pos}, nil
		case isIdentifierStart(ch):
			return l.readIdentifier(), nil
		case unicode.IsDigit(ch):
			return l.readNumber()
		case ch == '"':
			return l.readQuoted('"', TokenString)
		case ch == '\'':
			return l.readQuoted('\'', TokenChar)
		case ch == '<' && l.peek(1) == '-':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "<-", Pos: pos}, nil
		case ch == '+' && l.peek(1) == '+':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "++", Pos: pos}, nil
		case ch == '-' && l.peek(1) == '-':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "--", Pos: pos}, nil
		case (ch == '<' || ch == '>' || ch == '=') && l.peek(1) == '=':
			text := string([]rune{ch, '='})
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: text, Pos: pos}, nil
		case ch == '/' && l.peek(1) == '=':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "/=", Pos: pos}, nil
		case ch == '!' && l.peek(1) == '=':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "!=", Pos: pos}, nil
		case ch == '&' && l.peek(1) == '&':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "AND", Pos: pos}, nil
		case ch == '|' && l.peek(1) == '|':
			l.advance()
			l.advance()
			return Token{Kind: TokenOp, Text: "OR", Pos: pos}, nil
		case strings.ContainsRune("+-*/%=<>!", ch):
			l.advance()
			return Token{Kind: TokenOp, Text: string(ch), Pos: pos}, nil
		case strings.ContainsRune("(),:;<>.[] ", ch):
			l.advance()
			return Token{Kind: TokenSymbol, Text: string(ch), Pos: pos}, nil
		default:
			return Token{}, Diagnostic{Pos: pos, Message: fmt.Sprintf("unexpected character %q", ch)}
		}
	}
}

func (l *Lexer) AllTokens() ([]Token, error) {
	var tokens []Token
	for {
		t, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
		if t.Kind == TokenEOF {
			return tokens, nil
		}
	}
}
