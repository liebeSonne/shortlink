package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	panicCallErrText    = "panic call is not allowed"
	logFatalCallErrText = "log.Fatal call in main package not in main function is not allowed"
	osExitCallErrText   = "os.Exit call in main package not in main function is not allowed"
)

var Analyzer = &analysis.Analyzer{
	Name: "linter",
	Doc:  "check panic, log.Fatal, os.Exit (exclude files with `mock` prefix)",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		if isFileWithMockPrefix(filename) {
			continue
		}

		isMainPkg := pass.Pkg.Name() == "main"

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if isPanicCall(x) {
					pass.Reportf(x.Pos(), panicCallErrText)
				}
				if isLogFatalCall(pass, x) {
					if isMainPkg && !isInMainFunction(file, x.Pos()) {
						pass.Reportf(x.Pos(), logFatalCallErrText)
					}
				}
				if isOsExitCall(pass, x) {
					if isMainPkg && !isInMainFunction(file, x.Pos()) {
						pass.Reportf(x.Pos(), osExitCallErrText)
					}
				}
			}
			return true
		})
	}

	return nil, nil
}

func isPanicCall(call *ast.CallExpr) bool {
	if f, ok := call.Fun.(*ast.Ident); ok {
		return f.Name == "panic"
	}
	return false
}

func isLogFatalCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	return isPkgMethodCall(pass, call, "log", "Fatal")
}

func isOsExitCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	return isPkgMethodCall(pass, call, "os", "Exit")
}

func isPkgMethodCall(pass *analysis.Pass, call *ast.CallExpr, pkgName, methodName string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if pass.TypesInfo == nil {
		return false
	}

	obj := pass.TypesInfo.Uses[sel.Sel]
	if obj == nil {
		return false
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}

	sig := fn.Type().(*types.Signature)
	if sig.Recv() != nil {
		return false
	}

	pkg := fn.Pkg()
	if pkg == nil {
		return false
	}

	return pkg.Path() == pkgName && sel.Sel.Name == methodName
}

func isInMainFunction(file *ast.File, pos token.Pos) bool {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Name.Name == "main" && fn.Recv == nil {
			if pos >= fn.Pos() && pos <= fn.End() {
				return true
			}
		}
	}

	return false
}

func isFileWithMockPrefix(filename string) bool {
	base := filepath.Base(filename)
	base = strings.ToLower(base)
	return strings.HasPrefix(base, "mock")
}
