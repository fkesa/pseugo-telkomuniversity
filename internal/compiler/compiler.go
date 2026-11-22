package compiler

import (
    "bytes"
    "errors"
    "go/ast"
    "go/format"
    "go/parser"
    "go/token"
	"fmt"
)

type Result struct {
	Tokens   []Token
	Program  *Program
	GoSource string
}

// CheckSource performs lexical analysis, syntax parsing, and semantic/type checking
// on pseudocode source code without running code generation.
func CheckSource(source string) ([]*Program, error) {
	tokens, err := NewLexer(source).AllTokens()
	if err != nil {
		return nil, err
	}

	programs, err := NewParser(tokens).ParsePrograms()
	if err != nil {
		return nil, err
	}

	programNames := make(map[string]bool, len(programs))
	for _, prog := range programs {
		programNames[prog.Name] = true
	}

	for _, prog := range programs {
		checker := NewChecker()
		checker.SetPrograms(programNames)
		diagnostics := checker.Check(prog)

		if len(diagnostics) > 0 {
			return nil, fmt.Errorf(
				"semantic analysis failed in program %q:\n%w",
				prog.Name,
				errors.Join(diagnostics...),
			)
		}
	}

	return programs, nil
}

func Compile(source string) (*Result, error) {
	tokens, err := NewLexer(source).AllTokens()
	if err != nil {
		return nil, err
	}

	programs, err := CheckSource(source)
	if err != nil {
		return nil, err
	}

	goSource, err := generateGoPrograms(programs)
	if err != nil {
		return nil, err
	}

	return &Result{
		Tokens:   tokens,
		Program:  programs[0],
		GoSource: goSource,
	}, nil
}

func generateGoPrograms(programs []*Program) (string, error) {
    fset := token.NewFileSet()
    out := &ast.File{
        Name: ast.NewIdent("main"),
    }

    calledPrograms := make(map[string]bool)
    for _, prog := range programs {
        collectCalledPrograms(prog.Body, calledPrograms)
    }

    var declarations []ast.Decl
    var allCalls []ast.Stmt
    var topLevelCalls []ast.Stmt

    for i, prog := range programs {
        src, err := GenerateGo(prog)
        if err != nil {
            return "", fmt.Errorf(
                "program %q: %w",
                prog.Name,
                err,
            )
        }

        parsed, err := parser.ParseFile(
            fset,
            "",
            src,
            parser.ParseComments,
        )
        if err != nil {
            return "", err
        }

        var mainFn *ast.FuncDecl
        var localDecls []ast.Decl
        var helpers []ast.Decl

        for _, d := range parsed.Decls {
            switch v := d.(type) {
            case *ast.GenDecl:
                if v.Tok == token.IMPORT {
                    if i == 0 {
                        out.Decls = append(out.Decls, v)
                    }
                    continue
                }

                localDecls = append(localDecls, v)

            case *ast.FuncDecl:
                if v.Name.Name == "main" {
                    mainFn = v
                } else if i == 0 {
                    helpers = append(helpers, v)
                }
            }
        }

        if mainFn == nil {
            return "", fmt.Errorf(
                "generated function missing for %q",
                prog.Name,
            )
        }

        // Berikan nama unik pada setiap program.
        functionName := goIdent(prog.Name)

        mainFn.Name = ast.NewIdent(functionName)

        // Pindahkan deklarasi variabel dan tipe
        // ke dalam fungsi program masing-masing.
        body := make([]ast.Stmt, 0, len(localDecls)+len(mainFn.Body.List))

        for _, d := range localDecls {
            body = append(body, &ast.DeclStmt{
                Decl: d,
            })
        }

        body = append(body, mainFn.Body.List...)
        mainFn.Body.List = body

        declarations = append(declarations, mainFn)

        // Tambahkan pemanggilan program secara berurutan.
        callStmt := &ast.ExprStmt{
            X: &ast.CallExpr{
                Fun: ast.NewIdent(functionName),
            },
        }
        allCalls = append(allCalls, callStmt)
        if !calledPrograms[prog.Name] {
            topLevelCalls = append(topLevelCalls, callStmt)
        }

        if i == 0 {
            declarations = append(declarations, helpers...)
        }
    }

    calls := topLevelCalls
    if len(calls) == 0 {
        calls = allCalls
    }

    // Buat main Go yang menjalankan seluruh program.
    mainFn := &ast.FuncDecl{
        Name: ast.NewIdent("main"),
        Type: &ast.FuncType{
            Params: &ast.FieldList{},
        },
        Body: &ast.BlockStmt{
            List: calls,
        },
    }

    out.Decls = append(out.Decls, declarations...)
    out.Decls = append(out.Decls, mainFn)

    var buffer bytes.Buffer

    if err := format.Node(&buffer, fset, out); err != nil {
        return "", err
    }

    return buffer.String(), nil
}

func collectCalledPrograms(stmts []Stmt, called map[string]bool) {
    for _, s := range stmts {
        switch v := s.(type) {
        case CallStmt:
            called[v.Name] = true
        case ForStmt:
            collectCalledPrograms(v.Body, called)
        case WhileStmt:
            collectCalledPrograms(v.Body, called)
        case IfStmt:
            collectCalledPrograms(v.Then, called)
            collectCalledPrograms(v.Else, called)
        case RepeatStmt:
            collectCalledPrograms(v.Body, called)
        case SwitchStmt:
            for _, cl := range v.Cases {
                collectCalledPrograms(cl.Body, called)
            }
            collectCalledPrograms(v.Default, called)
        }
    }
}