package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const resetFile = "reset.gen.go"

//go:embed templates/*.tmpl
var templateFS embed.FS

type TemplateManager struct {
	templates *template.Template
}

type TemplateData struct {
	Package string
	Structs []StructInfo
}

func NewTemplateManager() (*TemplateManager, error) {
	funcMap := template.FuncMap{
		"zeroValue": prepareZeroValue,
		"isPtr":     func(f FieldInfo) bool { return f.IsPtr },
		"isSlice":   func(f FieldInfo) bool { return f.IsSlice },
		"isMap":     func(f FieldInfo) bool { return f.IsMap },
		"isStruct":  func(f FieldInfo) bool { return f.IsStruct },
		"isBasic":   func(f FieldInfo) bool { return f.IsBasic },
	}

	templates, err := template.New("reset").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded templates: %w", err)
	}

	return &TemplateManager{
		templates: templates,
	}, nil
}

func (tm *TemplateManager) Generate(data TemplateData) (string, error) {
	var buf bytes.Buffer
	if err := tm.templates.ExecuteTemplate(&buf, "file.tmpl", data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

func (tm *TemplateManager) GenerateFile(pkgPath string, structs []StructInfo) error {
	if len(structs) == 0 {
		return nil
	}

	data := TemplateData{
		Package: structs[0].Package,
		Structs: structs,
	}

	content, err := tm.Generate(data)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(pkgPath, resetFile)
	return os.WriteFile(outputPath, []byte(content), 0644)
}
