package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

const generateRest = "generate:reset"
const resetFile = "reset.gen.go"

type FieldInfo struct {
	Name       string
	Type       string
	IsPtr      bool
	IsSlice    bool
	IsMap      bool
	IsStruct   bool
	StructName string
	IsBasic    bool
}

type StructInfo struct {
	File    string
	Package string
	Name    string
	Fields  []FieldInfo
}

func run() error {
	packageStructsMap, err := collectPackageStructsMap()
	if err != nil {
		return fmt.Errorf("error on collect package structs map: %w", err)
	}

	for pkgPath, structs := range packageStructsMap {
		err = generateResetFile(pkgPath, structs)
		if err != nil {
			return fmt.Errorf("error on generate reset file: %w", err)
		}
	}

	fmt.Printf("Generated reset files for %d packages\n", len(packageStructsMap))

	return nil
}

func collectPackageStructsMap() (map[string][]StructInfo, error) {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles | packages.NeedName,
		Dir:  ".",
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("error on load packages: %w", err)
	}

	fmt.Printf("Loaded %d packages\n", len(pkgs))

	packageStructsMap := make(map[string][]StructInfo)

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			fmt.Printf("skip package '%s' with errors: %v\n", pkg.PkgPath, pkg.Errors)
			continue
		}

		for _, file := range pkg.Syntax {
			structs := parseFileNode(file, pkg.Name)
			if len(structs) > 0 && len(pkg.GoFiles) > 0 {
				pkgDIr := filepath.Dir(pkg.GoFiles[0])
				packageStructsMap[pkgDIr] = append(packageStructsMap[pkgDIr], structs...)
			}
		}
	}

	return packageStructsMap, nil
}

func parseFileContent(content, pkgName string) ([]StructInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return parseFileNode(node, pkgName), nil
}

func parseFileNode(file *ast.File, pkgName string) []StructInfo {
	var results []StructInfo
	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, isStruct := typeSpec.Type.(*ast.StructType)
			if !isStruct {
				continue
			}

			if hasRestCommentComment(genDecl, typeSpec) {
				fields := parseFields(structType)
				results = append(results, StructInfo{
					File:    file.Name.Name,
					Package: pkgName,
					Name:    typeSpec.Name.Name,
					Fields:  fields,
				})
			}
		}
		return true
	})

	return results
}

func hasRestCommentComment(
	genDecl *ast.GenDecl,
	typeSpec *ast.TypeSpec,
) bool {
	if typeSpec.Doc != nil {
		for _, comment := range typeSpec.Doc.List {
			text := strings.TrimSpace(comment.Text)
			isComment := strings.Contains(text, generateRest)
			if isComment {
				return true
			}
		}
	}

	if genDecl.Doc != nil {
		for _, comment := range genDecl.Doc.List {
			text := strings.TrimSpace(comment.Text)
			isComment := strings.Contains(text, generateRest)
			if isComment {
				return true
			}
		}
	}

	return false
}

func parseFields(structType *ast.StructType) []FieldInfo {
	var fields []FieldInfo

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		for _, name := range field.Names {
			fieldInfo := parseField(field.Type, name.Name)
			fields = append(fields, fieldInfo)
		}
	}

	return fields
}

func parseField(expr ast.Expr, fieldName string) FieldInfo {
	switch x := expr.(type) {
	case *ast.Ident:
		isBasic := isBasicType(x.Name)
		structName := x.Name
		if isBasic {
			structName = ""
		}
		return FieldInfo{
			Name:       fieldName,
			Type:       x.Name,
			IsStruct:   isStructType(x.Name),
			StructName: structName,
			IsBasic:    isBasic,
		}
	case *ast.StarExpr:
		field := parseField(x.X, "")
		return FieldInfo{
			Name:       fieldName,
			Type:       field.Type,
			IsPtr:      true,
			IsStruct:   field.IsStruct,
			StructName: field.StructName,
			IsBasic:    isBasicType(field.Type),
		}
	case *ast.ArrayType:
		field := parseField(x.Elt, "")
		var lenStr string
		if x.Len == nil {
			lenStr = ""
		} else {
			if bl, ok := x.Len.(*ast.BasicLit); ok {
				lenStr = bl.Value
			} else {
				lenStr = "N"
			}
		}
		fieldType := "[" + lenStr + "]" + field.Type
		return FieldInfo{
			Name:    fieldName,
			Type:    fieldType,
			IsSlice: true,
		}
	case *ast.MapType:
		keyField := parseField(x.Key, "")
		valueField := parseField(x.Value, "")
		fieldType := "map[" + keyField.Type + "]" + valueField.Type
		return FieldInfo{
			Name:  fieldName,
			Type:  fieldType,
			IsMap: true,
		}
	case *ast.SelectorExpr:
		var fieldType string
		if ident, ok := x.X.(*ast.Ident); ok {
			fieldType = ident.Name + "." + x.Sel.Name
		} else {
			fieldType = x.Sel.Name
		}
		return FieldInfo{
			Name: fieldName,
			Type: fieldType,
		}
	case *ast.InterfaceType:
		fieldType := "interface{}"
		return FieldInfo{
			Name: fieldName,
			Type: fieldType,
		}
	case *ast.ChanType:
		dir := ""
		switch x.Dir {
		case ast.SEND:
			dir = "chan<- "
		case ast.RECV:
			dir = "<-chan"
		default:
			dir = "chan "

		}
		field := parseField(x.Value, "")
		fieldType := dir + field.Type
		return FieldInfo{
			Name: fieldName,
			Type: fieldType,
		}
	case *ast.FuncType:
		fieldType := "func()"
		return FieldInfo{
			Name: fieldName,
			Type: fieldType,
		}
	default:
		return FieldInfo{
			Name: fieldName,
		}
	}
}

func isStructType(typeName string) bool {
	return !isBasicType(typeName)
}

func isBasicType(typeName string) bool {
	basicTypes := map[string]bool{
		"bool":       true,
		"byte":       true,
		"complex64":  true,
		"complex128": true,
		"error":      true,
		"float32":    true,
		"float64":    true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"rune":       true,
		"string":     true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uintptr":    true,
	}
	return basicTypes[typeName]
}

func generateResetFile(pkgPath string, structs []StructInfo) error {
	if len(structs) == 0 {
		return nil
	}

	outputPath := filepath.Join(pkgPath, resetFile)

	var buf bytes.Buffer

	buf.WriteString("// Code generated by reset generator; DO NOT EDIT.\n\n")
	buf.WriteString("package " + structs[0].Package + "\n\n")

	for _, s := range structs {
		err := generateResetMethod(&buf, s)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func generateResetMethod(buf *bytes.Buffer, s StructInfo) error {
	tpl := resetMethodTemplate

	funcMap := template.FuncMap{
		zeroValue: prepareZeroValue,
	}

	t := template.Must(template.New("reset").Funcs(funcMap).Parse(tpl))
	return t.Execute(buf, s)
}

func prepareZeroValue(typeName string) string {
	cleanType := typeName
	if idx := strings.LastIndex(typeName, "."); idx != -1 {
		cleanType = typeName[idx+1:]
	}

	switch cleanType {
	case "bool":
		return "false"
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64",
		"byte", "rune":
		return "0"
	case "complex64", "complex128":
		return "0 + 0i"
	default:
		if strings.Contains(typeName, ".") {
			return typeName + "{}"
		}
		return typeName + "{}"
	}
}
