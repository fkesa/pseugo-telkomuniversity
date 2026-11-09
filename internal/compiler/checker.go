package compiler

import (
	"fmt"
	"strings"
)

type Symbol struct {
	Type     string
	Constant bool
	AliasOf  string
}
type TypeInfo struct{ Fields []FieldDecl }

type Checker struct {
	globals map[string]Symbol
	types   map[string]TypeInfo
	errors  []error
}

func NewChecker() *Checker {
	return &Checker{globals: map[string]Symbol{}, types: map[string]TypeInfo{}}
}
func (c *Checker) report(pos Position, format string, args ...any) {
	c.errors = append(c.errors, Diagnostic{Pos: pos, Message: fmt.Sprintf(format, args...)})
}
func (c *Checker) knownType(name string) bool {
	if isPrimitiveType(name) {
		return true
	}
	_, ok := c.types[name]
	return ok
}
func (c *Checker) Check(prog *Program) []error {
	// Register record names first so declarations may refer to a type declared later.
	for _, d := range prog.Decls {
		if td, ok := d.(TypeDecl); ok {
			key := td.Name
			if _, exists := c.types[key]; exists || isPrimitiveType(td.Name) {
				c.report(td.Pos, "type %q is already defined", td.Name)
				continue
			}
			c.types[key] = TypeInfo{Fields: td.Fields}
		}
	}
	// Validate record fields after all record names have been registered.
	for _, d := range prog.Decls {
		if td, ok := d.(TypeDecl); ok {
			seen := map[string]bool{}
			for _, f := range td.Fields {
				if !c.knownType(f.TypeName) {
					c.report(f.Pos, "unknown field type %q", f.TypeName)
				}
				for _, n := range f.Names {
					k := strings.ToLower(n)
					if seen[k] {
						c.report(f.Pos, "duplicate field %q in type %q", n, td.Name)
					}
					seen[k] = true
				}
			}
		}
	}
	// Register variables, constants, and aliases. An alias declaration has the
	// form: new_name : previously_declared_variable.
	for _, d := range prog.Decls {
		if vd, ok := d.(VarDecl); ok {
			target, isAlias := c.globals[vd.TypeName]
			if !vd.Constant && !c.knownType(vd.TypeName) && !isAlias {
				c.report(vd.Pos, "unknown type or alias target %q", vd.TypeName)
			}
			if isAlias && (vd.Constant || len(vd.Names) != 1) {
				c.report(vd.Pos, "an alias declaration must define exactly one non-constant name")
			}
			for _, name := range vd.Names {
				if isReservedWord(name) {
					c.report(vd.Pos, "%q is reserved and cannot be a variable name", name)
					continue
				}
				if _, exists := c.globals[name]; exists {
					c.report(vd.Pos, "%q is already declared", name)
					continue
				}
				if isAlias {
					target.AliasOf = vd.TypeName
					c.globals[name] = target
				} else {
					c.globals[name] = Symbol{Type: normalizeType(vd.TypeName), Constant: vd.Constant}
				}
			}
		}
	}
	// Check constant initializers after symbol collection so names and types are
	// known, but only earlier constants may be referenced in an initializer.
	visibleConstants := map[string]Symbol{}
	for _, d := range prog.Decls {
		if vd, ok := d.(VarDecl); ok && vd.Constant {
			t := c.infer(vd.Value, c.globals)
			if !sameType(t, normalizeType(vd.TypeName)) {
				c.report(vd.Pos, "constant %q has type %s, but its value has type %s", vd.Names[0], vd.TypeName, t)
			}
			if !isPrimitiveType(vd.TypeName) {
				c.report(vd.Pos, "constant %q must use a primitive type", vd.Names[0])
			}
			if !c.isConstantExpr(vd.Value, visibleConstants) {
				c.report(vd.Pos, "initializer of constant %q must use literals or previously declared constants", vd.Names[0])
			}
			if len(vd.Names) == 1 {
				visibleConstants[vd.Names[0]] = Symbol{Type: normalizeType(vd.TypeName), Constant: true}
			}
		}
	}
	c.checkStatements(prog.Body, copySymbols(c.globals))
	return c.errors
}
func isReservedWord(name string) bool {
	words := []string{
    "program", "kamus", "algoritma", "endprogram", "constant", "type",
    "input", "output", "for", "while", "endwhile", "to", "downto", "do", "endfor",
    "repeat", "until", // <-- Tambahkan ini
    "if", "then", "elseif", "else", "endif",
    "div", "mod", "and", "or", "xor", "not", "integer", "int", "real", "float",
    "boolean", "bool", "string", "char", "true", "false", "benar", "salah",
}
	for _, w := range words {
		if strings.EqualFold(name, w) {
			return true
		}
	}
	return false
}
func copySymbols(src map[string]Symbol) map[string]Symbol {
	dst := map[string]Symbol{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
func sameType(a, b string) bool { return normalizeType(a) == normalizeType(b) }
func (c *Checker) checkStatements(stmts []Stmt, env map[string]Symbol) {
	for _, s := range stmts {
		switch v := s.(type) {
		case AssignStmt:
			target := v.Target
			if target == nil && v.Name != "" {
				target = IdentExpr{Name: v.Name, Pos: v.Pos}
			}
			sym, ok := c.resolveTarget(target, env)
			if !ok {
				c.report(v.Pos, "assignment target is not declared")
				break
			}
			if sym.Constant {
				c.report(v.Pos, "cannot assign to constant target")
			}
			got := c.infer(v.Value, env)
			if !sameType(sym.Type, got) && !(sym.Type == TypeReal && got == TypeInteger) {
				c.report(v.Pos, "cannot assign %s to %s target", got, sym.Type)
			}
		case InputStmt:
			targets := v.Targets
			if len(targets) == 0 {
				for _, n := range v.Names {
					targets = append(targets, IdentExpr{Name: n, Pos: v.Pos})
				}
			}
			for _, target := range targets {
				sym, ok := c.resolveTarget(target, env)
				if !ok {
					c.report(v.Pos, "input target is not declared")
				} else if sym.Constant {
					c.report(v.Pos, "cannot input into constant target")
				}
			}
		case OutputStmt:
			for _, e := range v.Values {
				c.infer(e, env)
			}
		case WhileStmt:
			cond := c.infer(v.Condition, env)
			if cond != TypeBoolean {
				c.report(v.Pos, "while condition must be boolean, got %s", cond)
			}
			c.checkStatements(v.Body, copySymbols(env))
		case ForStmt:
			start := c.infer(v.Start, env)
			end := c.infer(v.End, env)
			if start != TypeInteger || end != TypeInteger {
				c.report(v.Pos, "for loop bounds must both be integer, got %s and %s", start, end)
			}
			inner := copySymbols(env)
			inner[v.Counter] = Symbol{Type: TypeInteger}
			c.checkStatements(v.Body, inner)
		case IfStmt:
			cond := c.infer(v.Condition, env)
			if cond != TypeBoolean {
				c.report(v.Pos, "if condition must be boolean, got %s", cond)
			}
			c.checkStatements(v.Then, copySymbols(env))
			if len(v.Else) > 0 {
				c.checkStatements(v.Else, copySymbols(env))
			}
		case RepeatStmt:
			c.checkStatements(v.Body, copySymbols(env))
			cond := c.infer(v.Condition, env)
			if cond != TypeBoolean {
				c.report(v.Pos, "until condition must be boolean, got %s", cond)
			}
		}
	}
}

func (c *Checker) resolveTarget(e Expr, env map[string]Symbol) (Symbol, bool) {
	switch v := e.(type) {
	case IdentExpr:
		sym, ok := env[v.Name]
		return sym, ok
	case FieldExpr:
		typ := c.fieldType(v, env)
		return Symbol{Type: typ}, typ != "<error>"
	}
	return Symbol{}, false
}
func (c *Checker) fieldType(e FieldExpr, env map[string]Symbol) string {
	base, ok := e.Base.(IdentExpr)
	if !ok {
		c.report(e.Pos, "nested record fields are not supported")
		return "<error>"
	}
	sym, ok := env[base.Name]
	if !ok {
		c.report(e.Pos, "variable %q is not declared", base.Name)
		return "<error>"
	}
	info, ok := c.types[sym.Type]
	if !ok {
		c.report(e.Pos, "%q is not a record type", sym.Type)
		return "<error>"
	}
	for _, f := range info.Fields {
		for _, n := range f.Names {
			if strings.EqualFold(n, e.Field) {
				return normalizeType(f.TypeName)
			}
		}
	}
	c.report(e.Pos, "record type %q has no field %q", sym.Type, e.Field)
	return "<error>"
}

func (c *Checker) infer(e Expr, env map[string]Symbol) string {
	switch v := e.(type) {
	case LiteralExpr:
		switch v.Kind {
		case TokenInt:
			return TypeInteger
		case TokenReal:
			return TypeReal
		case TokenString:
			return TypeString
		case TokenChar:
			return TypeChar
		default:
			if v.Value == "true" || v.Value == "benar" || v.Value == "false" || v.Value == "salah" {
				return TypeBoolean
			}
		}
	case FieldExpr:
		return c.fieldType(v, env)
	case IdentExpr:
		sym, ok := env[v.Name]
		if !ok {
			c.report(v.Pos, "variable %q is not declared", v.Name)
			return "<error>"
		}
		return sym.Type
	case UnaryExpr:
		r := c.infer(v.Right, env)
		switch v.Op {
		case "NOT":
			if r != TypeBoolean {
				c.report(v.Pos, "NOT requires boolean operand, got %s", r)
			}
			return TypeBoolean
		case "+", "-":
			if !isNumericType(r) {
				c.report(v.Pos, "unary %s requires numeric operand, got %s", v.Op, r)
			}
			return r
		}
	case BinaryExpr:
		left := c.infer(v.Left, env)
		right := c.infer(v.Right, env)
		switch v.Op {
		case "+", "-", "*":
			if !isNumericType(left) || !isNumericType(right) {
				c.report(v.Pos, "operator %s requires numeric operands, got %s and %s", v.Op, left, right)
				return "<error>"
			}
			if left == TypeReal || right == TypeReal {
				return TypeReal
			}
			return TypeInteger
		case "DIV", "MOD":
			if left != TypeInteger || right != TypeInteger {
				c.report(v.Pos, "%s is defined only for integer operands, got %s and %s", v.Op, left, right)
				return "<error>"
			}
			return TypeInteger
		case "/":
			if !isNumericType(left) || !isNumericType(right) {
				c.report(v.Pos, "/ requires numeric operands")
				return "<error>"
			}
			return TypeReal
		case "%":
			c.report(v.Pos, "use MOD instead of %% in pseudocode")
			return "<error>"
		case "AND", "OR", "XOR":
			if left != TypeBoolean || right != TypeBoolean {
				c.report(v.Pos, "operator %s requires boolean operands, got %s and %s", v.Op, left, right)
			}
			return TypeBoolean
		case "=", "==", "/=", "!=":
			if !sameType(left, right) && !(isNumericType(left) && isNumericType(right)) {
				c.report(v.Pos, "comparison operands must have compatible types, got %s and %s", left, right)
			}
			return TypeBoolean
		case "<", "<=", ">", ">=":
			if !isNumericType(left) || !isNumericType(right) {
				c.report(v.Pos, "ordered comparison requires numeric operands, got %s and %s", left, right)
			}
			return TypeBoolean
		}
	}
	return "<error>"
}

func (c *Checker) isConstantExpr(e Expr, env map[string]Symbol) bool {
	switch v := e.(type) {
	case LiteralExpr:
		return true
	case IdentExpr:
		sym, ok := env[v.Name]
		return ok && sym.Constant
	case UnaryExpr:
		return c.isConstantExpr(v.Right, env)
	case BinaryExpr:
		return c.isConstantExpr(v.Left, env) && c.isConstantExpr(v.Right, env)
	default:
		return false
	}
}
