package noexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "do not use os.Exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if strings.Contains(pass.Fset.Position(file.Pos()).Filename, "/go-build/") {
			return nil, nil
		}

		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "main" || funcDecl.Recv != nil {
				continue
			}
			for _, stmt := range funcDecl.Body.List {
				ast.Inspect(stmt, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "os" && fun.Sel.Name == "Exit" {
							pass.Reportf(call.Pos(), "do not use os.Exit in main function")
						}
					}
					return true
				})
			}
		}
	}
	return nil, nil
}
