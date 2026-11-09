package compiler

type Program struct {
	Name  string
	Decls []Decl
	Body  []Stmt
	Pos   Position
}

type Decl interface{ isDecl() }

type VarDecl struct {
	Names    []string
	TypeName string
	Constant bool
	Value    Expr
	Pos      Position
}

func (VarDecl) isDecl() {}

type TypeDecl struct {
	Name   string
	Fields []FieldDecl
	Pos    Position
}

func (TypeDecl) isDecl() {}

type FieldDecl struct {
	Names    []string
	TypeName string
	Pos      Position
}

type Stmt interface{ isStmt() }

type AssignStmt struct {
	Target Expr
	Name   string
	Value  Expr
	Pos    Position
}

func (AssignStmt) isStmt() {}

type InputStmt struct {
	Targets []Expr
	Names   []string
	Pos     Position
}

func (InputStmt) isStmt() {}

type OutputStmt struct {
	Values []Expr
	Pos    Position
}

func (OutputStmt) isStmt() {}

type ForStmt struct {
	Counter string
	Start   Expr
	End     Expr
	Down    bool
	Body    []Stmt
	Pos     Position
}

func (ForStmt) isStmt() {}

type WhileStmt struct {
	Condition Expr
	Body      []Stmt
	Pos       Position
}

func (WhileStmt) isStmt() {}

type Expr interface{ isExpr() }

type LiteralExpr struct {
	Kind  TokenKind
	Value string
	Pos   Position
}

func (LiteralExpr) isExpr() {}

type IdentExpr struct {
	Name string
	Pos  Position
}

func (IdentExpr) isExpr() {}

type FieldExpr struct {
	Base  Expr
	Field string
	Pos   Position
}

func (FieldExpr) isExpr() {}

type UnaryExpr struct {
	Op    string
	Right Expr
	Pos   Position
}

func (UnaryExpr) isExpr() {}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
	Pos   Position
}

func (BinaryExpr) isExpr() {}

type IfStmt struct {
	Condition Expr
	Then      []Stmt
	Else      []Stmt
	Pos       Position
}

func (IfStmt) isStmt() {}

type RepeatStmt struct {
	Body      []Stmt
	Condition Expr
	Pos       Position
}

func (RepeatStmt) isStmt() {}

type CaseClause struct {
	Values []Expr
	Body   []Stmt
	Pos    Position
}

type SwitchStmt struct {
	Expr    Expr
	Cases   []CaseClause
	Default []Stmt
	Pos     Position
}

func (SwitchStmt) isStmt() {}

type CastExpr struct {
	TargetType string
	Expr       Expr
	Pos        Position
}

func (CastExpr) isExpr() {}
