// Package exitcheck проверяет, что в функции main пакета main не вызывается os.Exit.
package exitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "запрещает прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Проверяем, что мы в пакете main.
		if pass.Pkg.Name() != "main" {
			continue
		}

		// Проходим по всем узлам AST.
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем вызовы функций.
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Проверяем, что это os.Exit.
			fn, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if fn.Sel.Name == "Exit" {
				// Проверяем, что os импортирован.
				if ident, ok := fn.X.(*ast.Ident); ok && ident.Name == "os" {
					pass.Reportf(call.Pos(), "запрещен прямой вызов os.Exit в main")
				}
			}
			return true
		})
	}
	return nil, nil
}
