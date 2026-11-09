package compiler

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	index  int
}

func NewParser(tokens []Token) *Parser { return &Parser{tokens: tokens} }
func (p *Parser) current() Token {
	if p.index >= len(p.tokens) {
		return Token{Kind: TokenEOF, Pos: Position{Line: 1, Column: 1}}
	}
	return p.tokens[p.index]
}
// func (p *Parser) peek(n int) Token {
// 	i := p.index + n
// 	if i >= len(p.tokens) {
// 		return Token{Kind: TokenEOF, Pos: p.current().Pos}
// 	}
// 	return p.tokens[i]
// }
func (p *Parser) advance() Token {
	t := p.current()
	if p.index < len(p.tokens) {
		p.index++
	}
	return t
}
func (p *Parser) isWord(word string) bool {
	return p.current().Kind == TokenIdent && strings.EqualFold(p.current().Text, word)
}
func (p *Parser) isOp(op string) bool {
	return p.current().Kind == TokenOp && strings.EqualFold(p.current().Text, op)
}
func (p *Parser) isSymbol(s string) bool {
	if p.current().Text != s {
		return false
	}
	if p.current().Kind == TokenSymbol {
		return true
	}
	return p.current().Kind == TokenOp && (s == "<" || s == ">")
}
func (p *Parser) expectWord(word string) (Token, error) {
	if !p.isWord(word) {
		return Token{}, p.errorf(p.current(), "expected %q", word)
	}
	return p.advance(), nil
}
func (p *Parser) expectSymbol(s string) (Token, error) {
	if !p.isSymbol(s) {
		return Token{}, p.errorf(p.current(), "expected %q", s)
	}
	return p.advance(), nil
}
func (p *Parser) expectIdent() (Token, error) {
	if p.current().Kind != TokenIdent {
		return Token{}, p.errorf(p.current(), "expected identifier")
	}
	return p.advance(), nil
}
func (p *Parser) errorf(t Token, format string, args ...any) error {
	return Diagnostic{Pos: t.Pos, Message: fmt.Sprintf(format, args...)}
}
func (p *Parser) skipSeparators() {
	for p.current().Kind == TokenNewline || p.isSymbol(";") {
		p.advance()
	}
}
func (p *Parser) expectLineEnd() error {
	if p.current().Kind == TokenNewline || p.isSymbol(";") || p.current().Kind == TokenEOF {
		p.skipSeparators()
		return nil
	}
	return p.errorf(p.current(), "expected end of line")
}

func (p *Parser) ParsePrograms() ([]*Program, error) {
	var programs []*Program
	p.skipSeparators()

	for p.current().Kind != TokenEOF {
		prog, err := p.parseSingleProgram()
		if err != nil {
			return nil, err
		}
		programs = append(programs, prog)
		p.skipSeparators()
	}

	if len(programs) == 0 {
		return nil, p.errorf(p.current(), "expected at least one program")
	}

	return programs, nil
}

func (p *Parser) ParseProgram() (*Program, error) {
	programs, err := p.ParsePrograms()
	if err != nil {
		return nil, err
	}
	if len(programs) != 1 {
		return nil, p.errorf(
			p.current(),
			"expected exactly one program, found %d",
			len(programs),
		)
	}
	return programs[0], nil
}

func (p *Parser) parseSingleProgram() (*Program, error) {
	p.skipSeparators()
	start, err := p.expectWord("program")
	if err != nil {
		return nil, err
	}
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	if _, err = p.expectWord("kamus"); err != nil {
		return nil, err
	}
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	prog := &Program{Name: name.Text, Pos: start.Pos}
	for !p.isWord("algoritma") {
		p.skipSeparators()
		if p.isWord("algoritma") {
			break
		}
		if p.current().Kind == TokenEOF {
			return nil, p.errorf(p.current(), "expected 'algoritma' section")
		}
		d, err := p.parseDecl()
		if err != nil {
			return nil, err
		}
		prog.Decls = append(prog.Decls, d)
	}
	p.advance() // algoritma
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	body, err := p.parseStatements("endprogram")
	if err != nil {
		return nil, err
	}
	prog.Body = body
	if _, err = p.expectWord("endprogram"); err != nil {
		return nil, err
	}

	return prog, nil
}

func (p *Parser) parseDecl() (Decl, error) {
	p.skipSeparators()
	if p.isWord("type") {
		return p.parseTypeDecl()
	}
	if p.isWord("constant") {
		kw := p.advance()
		name, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		if _, err = p.expectSymbol(":"); err != nil {
			return nil, err
		}
		typ, err := p.parseTypeName()
		if err != nil {
			return nil, err
		}
		if !p.isOp("<-") {
			return nil, p.errorf(p.current(), "expected '<-' in constant declaration")
		}
		p.advance()
		value, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		if err = p.expectLineEnd(); err != nil {
			return nil, err
		}
		return VarDecl{Names: []string{name.Text}, TypeName: typ, Constant: true, Value: value, Pos: kw.Pos}, nil
	}
	pos := p.current().Pos
	names, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	if _, err = p.expectSymbol(":"); err != nil {
		return nil, err
	}
	typ, err := p.parseTypeName()
	if err != nil {
		return nil, err
	}
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	return VarDecl{Names: names, TypeName: typ, Pos: pos}, nil
}

func (p *Parser) parseTypeDecl() (Decl, error) {
	kw := p.advance()
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if _, err = p.expectSymbol("<"); err != nil {
		return nil, err
	}
	fields := []FieldDecl{}
	p.skipSeparators()
	for !p.isSymbol(">") {
		if p.current().Kind == TokenEOF {
			return nil, p.errorf(p.current(), "unterminated type declaration; expected '>'")
		}
		pos := p.current().Pos
		names, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		if _, err = p.expectSymbol(":"); err != nil {
			return nil, err
		}
		typ, err := p.parseTypeName()
		if err != nil {
			return nil, err
		}
		fields = append(fields, FieldDecl{Names: names, TypeName: typ, Pos: pos})
		if err = p.expectLineEnd(); err != nil {
			return nil, err
		}
		p.skipSeparators()
	}
	p.advance() // >
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	return TypeDecl{Name: name.Text, Fields: fields, Pos: kw.Pos}, nil
}

func (p *Parser) parseNameList() ([]string, error) {
	first, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	names := []string{first.Text}
	for p.isSymbol(",") {
		p.advance()
		next, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		names = append(names, next.Text)
	}
	return names, nil
}

func (p *Parser) parseTypeName() (string, error) {
	t, err := p.expectIdent()
	if err != nil {
		return "", err
	}
	return t.Text, nil
}

func (p *Parser) parseRepeat() (Stmt, error) {
	kw := p.advance() // lewati token 'repeat'
	if err := p.expectLineEnd(); err != nil {
		return nil, err
	}

	body, err := p.parseStatements("until")
	if err != nil {
		return nil, err
	}

	if _, err = p.expectWord("until"); err != nil {
		return nil, err
	}

	cond, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}

	return RepeatStmt{Body: body, Condition: cond, Pos: kw.Pos}, nil
}

func (p *Parser) parseStatements(endWords ...string) ([]Stmt, error) {
	var stmts []Stmt
	p.skipSeparators()
	for {
		if p.current().Kind == TokenEOF {
			return nil, p.errorf(p.current(), "expected %v before end of file", endWords)
		}
		matched := false
		for _, w := range endWords {
			if p.isWord(w) {
				matched = true
				break
			}
		}
		if matched {
			break
		}
		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
		if err = p.expectLineEnd(); err != nil {
			return nil, err
		}
		p.skipSeparators()
	}
	return stmts, nil
}

func (p *Parser) parseStatement() (Stmt, error) {
	t := p.current()
	if p.isWord("input") {
		return p.parseInput()
	}
	if p.isWord("output") {
		return p.parseOutput()
	}
	if p.isWord("for") {
		return p.parseFor()
	}
	if p.isWord("while") {
		return p.parseWhile()
	}
	if p.isWord("if") {
		return p.parseIf()
	}
	if p.isWord("repeat") {
		return p.parseRepeat()
	}
	if t.Kind == TokenIdent {
		target, err := p.parsePostfixIdent()
		if err != nil {
			return nil, err
		}
		if !p.isOp("<-") {
			return nil, p.errorf(p.current(), "expected '<-' after assignment target")
		}
		p.advance()
		value, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		return AssignStmt{Target: target, Value: value, Pos: t.Pos}, nil
	}
	return nil, p.errorf(t, "expected assignment, input, output, if, for, or while statement")
}

func (p *Parser) parseIf() (Stmt, error) {
	kw := p.advance() // if
	cond, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if _, err = p.expectWord("then"); err != nil {
		return nil, err
	}
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}

	thenBody, err := p.parseStatements("elseif", "else", "endif")
	if err != nil {
		return nil, err
	}

	var elseBody []Stmt
	if p.isWord("elseif") {
		nestedIf, err := p.parseIf()
		if err != nil {
			return nil, err
		}
		elseBody = []Stmt{nestedIf}
		return IfStmt{Condition: cond, Then: thenBody, Else: elseBody, Pos: kw.Pos}, nil
	}

	if p.isWord("else") {
		p.advance() // else
		if err = p.expectLineEnd(); err != nil {
			return nil, err
		}
		elseBody, err = p.parseStatements("endif")
		if err != nil {
			return nil, err
		}
	}

	if _, err = p.expectWord("endif"); err != nil {
		return nil, err
	}

	return IfStmt{Condition: cond, Then: thenBody, Else: elseBody, Pos: kw.Pos}, nil
}

func (p *Parser) parseInput() (Stmt, error) {
	kw := p.advance()
	if _, err := p.expectSymbol("("); err != nil {
		return nil, err
	}
	var targets []Expr
	for {
		e, eerr := p.parsePostfixIdent()
		if eerr != nil {
			return nil, eerr
		}
		targets = append(targets, e)
		if !p.isSymbol(",") {
			break
		}
		p.advance()
	}
	if _, err := p.expectSymbol(")"); err != nil {
		return nil, err
	}
	return InputStmt{Targets: targets, Pos: kw.Pos}, nil
}

func (p *Parser) parseOutput() (Stmt, error) {
	kw := p.advance()
	if _, err := p.expectSymbol("("); err != nil {
		return nil, err
	}
	var values []Expr
	for {
		v, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
		if !p.isSymbol(",") {
			break
		}
		p.advance()
	}
	if _, err := p.expectSymbol(")"); err != nil {
		return nil, err
	}
	return OutputStmt{Values: values, Pos: kw.Pos}, nil
}

func (p *Parser) parseFor() (Stmt, error) {
	kw := p.advance()
	counter, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if !p.isOp("<-") {
		return nil, p.errorf(p.current(), "expected '<-' after for counter")
	}
	p.advance()
	start, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	down := false
	if p.isWord("to") {
		p.advance()
	} else if p.isWord("downto") {
		p.advance()
		down = true
	} else {
		return nil, p.errorf(p.current(), "expected 'to' or 'downto' in for loop")
	}
	end, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if _, err = p.expectWord("do"); err != nil {
		return nil, err
	}
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	body, err := p.parseStatements("endfor")
	if err != nil {
		return nil, err
	}
	if _, err = p.expectWord("endfor"); err != nil {
		return nil, err
	}
	return ForStmt{Counter: counter.Text, Start: start, End: end, Down: down, Body: body, Pos: kw.Pos}, nil
}

func (p *Parser) parsePostfixIdent() (Expr, error) {
	t, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	var e Expr = IdentExpr{Name: t.Text, Pos: t.Pos}
	for p.isSymbol(".") {
		p.advance()
		f, er := p.expectIdent()
		if er != nil {
			return nil, er
		}
		e = FieldExpr{Base: e, Field: f.Text, Pos: f.Pos}
	}
	return e, nil
}

func (p *Parser) parseWhile() (Stmt, error) {
	kw := p.advance()
	cond, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if _, err = p.expectWord("do"); err != nil {
		return nil, err
	}
	if err = p.expectLineEnd(); err != nil {
		return nil, err
	}
	body, err := p.parseStatements("endwhile")
	if err != nil {
		return nil, err
	}
	if _, err = p.expectWord("endwhile"); err != nil {
		return nil, err
	}
	return WhileStmt{Condition: cond, Body: body, Pos: kw.Pos}, nil
}

// Pratt parser: precedence 1 is the weakest binding.
func precedence(op string) int {
	switch strings.ToUpper(op) {
	case "OR":
		return 2
	case "XOR":
		return 3
	case "AND":
		return 4
	case "=", "==", "/=", "!=", "<", "<=", ">", ">=":
		return 5
	case "+", "-":
		return 6
	case "*", "DIV", "MOD", "/", "%":
		return 7
	default:
		return 0
	}
}

func (p *Parser) parseExpression(minPrec int) (Expr, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}
	for p.current().Kind == TokenOp {
		op := p.current().Text
		prec := precedence(op)
		if prec < minPrec || prec == 0 {
			break
		}
		tok := p.advance()
		right, err := p.parseExpression(prec + 1)
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{Left: left, Op: strings.ToUpper(tok.Text), Right: right, Pos: tok.Pos}
	}
	return left, nil
}

func (p *Parser) parsePrefix() (Expr, error) {
	t := p.current()
	if p.isOp("NOT") || p.isOp("-") || p.isOp("+") {
		p.advance()
		right, err := p.parseExpression(8)
		if err != nil {
			return nil, err
		}
		return UnaryExpr{Op: strings.ToUpper(t.Text), Right: right, Pos: t.Pos}, nil
	}
	switch t.Kind {
	case TokenInt, TokenReal, TokenString, TokenChar:
		p.advance()
		return LiteralExpr{Kind: t.Kind, Value: t.Text, Pos: t.Pos}, nil
	case TokenIdent:
		if strings.EqualFold(t.Text, "true") || strings.EqualFold(t.Text, "false") || strings.EqualFold(t.Text, "benar") || strings.EqualFold(t.Text, "salah") {
			p.advance()
			return LiteralExpr{Kind: TokenIdent, Value: strings.ToLower(t.Text), Pos: t.Pos}, nil
		}
		p.advance()
		var e Expr = IdentExpr{Name: t.Text, Pos: t.Pos}
		for p.isSymbol(".") {
			p.advance()
			f, err := p.expectIdent()
			if err != nil {
				return nil, err
			}
			e = FieldExpr{Base: e, Field: f.Text, Pos: f.Pos}
		}
		return e, nil
	case TokenSymbol:
		if t.Text == "(" {
			p.advance()
			e, err := p.parseExpression(1)
			if err != nil {
				return nil, err
			}
			if _, err := p.expectSymbol(")"); err != nil {
				return nil, err
			}
			return e, nil
		}
	}
	return nil, p.errorf(t, "expected expression")
}