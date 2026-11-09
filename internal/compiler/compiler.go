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

func Compile(source string) (*Result, error) {
    tokens, err := NewLexer(source).AllTokens()
    if err != nil {
        return nil, err
    }

    programs, err := NewParser(tokens).ParsePrograms()
    if err != nil {
        return nil, err
    }

    for _, prog := range programs {
        checker := NewChecker()
        diagnostics := checker.Check(prog)

        if len(diagnostics) > 0 {
            return nil, fmt.Errorf(
                "semantic analysis failed in program %q:\n%w",
                prog.Name,
                errors.Join(diagnostics...),
            )
        }
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

    var declarations []ast.Decl
    var calls []ast.Stmt

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
        calls = append(calls, &ast.ExprStmt{
            X: &ast.CallExpr{
                Fun: ast.NewIdent(functionName),
            },
        })

        if i == 0 {
            declarations = append(declarations, helpers...)
        }
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