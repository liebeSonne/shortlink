package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	type on struct {
		testdataDir string
	}
	type fileData struct {
		file               string
		structs            []string
		patterns           []string
		packages           []string
		unexpectedStructs  []string
		unexpectedPatterns []string
	}
	type expected struct {
		files []fileData
		err   error
	}

	testCases := []struct {
		name     string
		on       on
		expected expected
	}{
		{
			"basic structures",
			on{"testdata/basic"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"User", "Profile", "Product", "Child"},
						patterns: []string{
							"package testdata",
							"r.ID = 0",
							"r.Name = \"\"",
							"r.Email = \"\"",
							"r.Tags = r.Tags[:0]",
							"clear(r.Metadata)",
							"if r.Profile != nil {",
							"if resetter, ok := r.Child.(interface{ Reset() }); ok {",
						},
						packages:           []string{"testdata"},
						unexpectedStructs:  []string{"NoResetStruct"},
						unexpectedPatterns: []string{"func (r *NoResetStruct) Reset())"},
					},
				},
			},
		},
		{
			"subpackages",
			on{"testdata/subpackages"},
			expected{
				files: []fileData{
					{
						file:    fmt.Sprintf("pkg1/%s", resetFile),
						structs: []string{"Struct1"},
						patterns: []string{
							"package pkg1",
							"func (r *Struct1) Reset()",
						},
						packages: []string{"pkg1"},
					},
					{
						file:    fmt.Sprintf("pkg2/%s", resetFile),
						structs: []string{"Struct2", "Struct3"},
						patterns: []string{
							"package pkg2",
							"func (r *Struct2) Reset()",
							"func (r *Struct3) Reset()",
						},
						packages: []string{"pkg2"},
					},
					{
						file:    fmt.Sprintf("pkg3/%s", resetFile),
						structs: []string{"Struct4"},
						patterns: []string{
							"package pkg3",
							"func (r *Struct4) Reset()",
						},
						packages: []string{"pkg3"},
					},
				},
			},
		},
		{
			"complex types",
			on{"testdata/complex_types"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"ComplexStruct"},
						patterns: []string{
							"package complex_types",
							"r.BoolField = false",
							"r.IntField = 0",
							"r.StringField = \"\"",
							"r.FloatField = 0",
							"r.ByteField = 0",
							"r.RuneField = 0",
							"r.ComplexField = 0 + 0i",
							"r.SliceField = r.SliceField[:0]",
							"clear(r.MapField)",
							"r.ArrayField = [3]int{}",
							"r.ChanField = nil",
							"r.FuncField = nil",
							"if r.StructPtr != nil {",
							"*r.StructPtr = SubStruct{}",
						},
						packages: []string{"complex"},
					},
				},
			},
		},
		{
			"nested structures",
			on{"testdata/nested"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"Parent", "Child", "GrandChild"},
						patterns: []string{
							"package nested",
							"func (r *Parent) Reset()",
							"func (r *Child) Reset()",
							"func (r *GrandChild) Reset()",
							"if resetter, ok := r.Child.(interface{ Reset() }); ok {",
							"resetter.Reset()",
							"if resetter, ok := r.GrandChild.(interface{ Reset() }); ok {",
							"r.Child = Child{}",
							"if r.ChildPtr != nil {",
						},
						packages: []string{"nested"},
					},
				},
			},
		},
		{
			"embedded structures",
			on{"testdata/embedded"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"Embedded", "WithEmbedded"},
						patterns: []string{
							"package embedded",
							"func (r *Embedded) Reset()",
							"func (r *WithEmbedded) Reset()",
							"r.Value = 0",
							"r.Name = \"\"",
						},
						packages: []string{"embedded"},
						unexpectedPatterns: []string{
							"r.Embedded",
						},
					},
				},
			},
		},
		{
			"ignore patterns",
			on{"testdata/ignore"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"RootStruct"},
						patterns: []string{
							"package ignore",
							"func (r *RootStruct) Reset()",
							"r.ID = 0",
						},
						packages:          []string{"ignore"},
						unexpectedStructs: []string{"VendorStruct", "GitStruct"},
						unexpectedPatterns: []string{
							"func (r *VendorStruct) Reset()",
							"func (r *GitStruct) Reset()",
						},
					},
				},
			},
		},
		{
			"multiple structs with different fields",
			on{"testdata/multiple"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"User", "Admin", "Guest"},
						patterns: []string{
							"package multiple",
							"func (r *User) Reset()",
							"func (r *Admin) Reset()",
							"func (r *Guest) Reset()",
							"r.Age = 0",
							"r.Role = \"\"",
							"r.Permissions = r.Permissions[:0]",
							"clear(r.Settings)",
						},
						packages: []string{"multiple"},
					},
				},
			},
		},
		{
			"pointers and references",
			on{"testdata/pointers"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"PointerStruct"},
						patterns: []string{
							"package pointers",
							"func (r *PointerStruct) Reset()",
							"if r.IntPtr != nil {",
							"*r.IntPtr = 0",
							"if r.StringPtr != nil {",
							"*r.StringPtr = \"\"",
							"if r.BoolPtr != nil {",
							"*r.BoolPtr = false",
							"if r.StructPtr != nil {",
							"if resetter, ok := r.StructPtr.(interface{ Reset() }); ok {",
							"resetter.Reset()",
						},
						packages: []string{"pointers"},
					},
				},
			},
		},
		{
			"maps and slices with complex types",
			on{"testdata/collections"},
			expected{
				files: []fileData{
					{
						file:    resetFile,
						structs: []string{"CollectionStruct"},
						patterns: []string{
							"package collections",
							"func (r *CollectionStruct) Reset()",
							"r.IntSlice = r.IntSlice[:0]",
							"r.StringSlice = r.StringSlice[:0]",
							"clear(r.StringMap)",
							"clear(r.IntMap)",
							"r.StructSlice = r.StructSlice[:0]",
							"clear(r.StructMap)",
							"r.InterfaceSlice = r.InterfaceSlice[:0]",
							"clear(r.InterfaceMap)",
						},
						packages: []string{"collections"},
					},
				},
			},
		},
		{
			"empty directory",
			on{"testdata/empty"},
			expected{},
		},
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("faile to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := os.Stat(tc.on.testdataDir); os.IsNotExist(err) {
				t.Skipf("testdata directory %s does not exist, skipping test", tc.on.testdataDir)
			}

			t.Cleanup(func() {
				cleanGeneratedFiles(t, tc.on.testdataDir)
			})

			err := os.Chdir(tc.on.testdataDir)
			if err != nil {
				t.Fatalf("faile to chdir to test data directory: %v", err)
			}

			err = run()
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
				return
			} else {
				require.NoError(t, err)
			}

			for _, file := range tc.expected.files {
				if _, err := os.Stat(file.file); os.IsNotExist(err) {
					t.Errorf("expected file %s was not created:", file.file)
				}
			}

			for _, file := range tc.expected.files {
				content, err := os.ReadFile(file.file)
				if err != nil {
					t.Errorf("faile to read file %s: %v", file.file, err)
				}
				contentStr := string(content)

				for _, structName := range file.structs {
					if !strings.Contains(contentStr, "func (r *"+structName+") Reset()") {
						t.Errorf("faile to find Reset method by struct %s in file %s", structName, file.file)
					}
				}

				for _, structName := range file.unexpectedStructs {
					if strings.Contains(contentStr, "func (r *"+structName+") Reset()") {
						t.Errorf("Reset method by struct %s should not exist in file %s", structName, file.file)
					}
				}

				for _, pattern := range file.patterns {
					if !strings.Contains(contentStr, pattern) {
						t.Errorf("pattern %s not found in file %s", pattern, file.file)
					}
				}

				for _, pattern := range file.unexpectedPatterns {
					if strings.Contains(contentStr, pattern) {
						t.Errorf("pattern %s nshould not exist in file %s", pattern, file.file)
					}
				}

				for _, pkg := range file.packages {
					if !strings.Contains(contentStr, "package "+pkg) {
						t.Errorf("package %s not found in file %s", pkg, file.file)
					}
				}
			}

			err = os.Chdir(originalDir)
			if err != nil {
				t.Fatalf("failed to change back to original directory: %v", err)
			}
		})
	}
}

func TestParseRestComments(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCount int
	}{
		{
			name: "comment directly above",
			code: `
package test

// generate:reset
type TestStruct struct {
	ID int
}
`,
			expectCount: 1,
		},
		{
			name: "no comment",
			code: `
package test

type TestStruct struct {
	ID int
}
`,
			expectCount: 0,
		},
		{
			name: "different comment",
			code: `
package test

// some other comment
type TestStruct struct {
	ID int
}
`,
			expectCount: 0,
		},
		{
			name: "comment with space",
			code: `
package test

//  generate:reset  
type TestStruct struct {
	ID int
}
`,
			expectCount: 1,
		},
		{
			name: "comment with empty line between",
			code: `
package test

// generate:reset

type TestStruct struct {
	ID int
}
`,
			expectCount: 0,
		},
		{
			name: "comment before const",
			code: `
package test

// generate:reset
const SomeConst = 1

type TestStruct struct {
	ID int
}
`,
			expectCount: 0,
		},
		{
			name: "first from multiple comments",
			code: `
package test

// generate:reset
type TestStruct1 struct {
	ID int
}

type TestStruct2 struct {
	ID int
}
`,
			expectCount: 1,
		},
		{
			name: "second from multiple comments",
			code: `
package test

type TestStruct1 struct {
	ID int
}

// generate:reset
type TestStruct2 struct {
	ID int
}
`,
			expectCount: 1,
		},
		{
			name: "multiple struct",
			code: `
package test

// some comment
type (
  	// some 1 comment 
  	TestStruct1 struct {
		ID int
  	}
	// some 2 comment
  	// generate:reset
  	TestStruct2 struct {
		ID int
  	}
)
`,
			expectCount: 1,
		},
		{
			name: "multiple struct comment",
			code: `
package test

// some comment
// generate:reset
type (
  	// some 1 comment 
  	TestStruct1 struct {
		ID int
  	}
	// some 2 comment
  	// generate:reset
  	TestStruct2 struct {
		ID int
  	}
)
`,
			expectCount: 2,
		},
		{
			name: "block comment",
			code: `
package test

/*
 generate:reset
*/
type TestStruct struct {
	ID int
}
`,
			expectCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			structs, err := parseFileContent(tc.code, "test")
			require.NoError(t, err)

			require.Len(t, structs, tc.expectCount)
		})
	}
}

func TestParseFields(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []FieldInfo
	}{
		{
			name: "basic types",
			code: `
package test

type TestStruct struct {
	ID     int
	Name   string
	Active bool
}
`,
			expected: []FieldInfo{
				{Name: "ID", Type: "int", IsBasic: true},
				{Name: "Name", Type: "string", IsBasic: true},
				{Name: "Active", Type: "bool", IsBasic: true},
			},
		},
		{
			name: "slice and map",
			code: `
		package test
		
		type TestStruct struct {
			Tags     []string
			Metadata map[string]interface{}
		}
		`,
			expected: []FieldInfo{
				{Name: "Tags", Type: "[]string", IsSlice: true},
				{Name: "Metadata", Type: "map[string]interface{}", IsMap: true},
			},
		},
		{
			name: "pointer types",
			code: `
		package test
		
		type TestStruct struct {
			Ptr     *int
			StrPtr  *string
			Struct  *SubStruct
		}
		
		type SubStruct struct {
			Value int
		}
		`,
			expected: []FieldInfo{
				{Name: "Ptr", Type: "int", IsPtr: true, IsBasic: true},
				{Name: "StrPtr", Type: "string", IsPtr: true, IsBasic: true},
				{Name: "Struct", Type: "SubStruct", IsPtr: true, IsStruct: true, StructName: "SubStruct"},
			},
		},
		{
			name: "nested structs",
			code: `
		package test
		
		type TestStruct struct {
			Child   SubStruct
			PtrChild *SubStruct
		}
		
		type SubStruct struct {
			Value int
		}
		`,
			expected: []FieldInfo{
				{Name: "Child", Type: "SubStruct", IsStruct: true, StructName: "SubStruct"},
				{Name: "PtrChild", Type: "SubStruct", IsPtr: true, IsStruct: true, StructName: "SubStruct"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tc.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			var structType *ast.StructType
			ast.Inspect(file, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						if structType == nil {
							structType = st
						}
						return false
					}
				}
				return true
			})

			if structType == nil {
				t.Fatal("No struct found in test code")
			}

			fields := parseFields(structType)

			require.Len(t, fields, len(tc.expected), fmt.Sprintf("go fields: %+v", fields))

			for i, expected := range tc.expected {
				if i >= len(fields) {
					break
				}
				field := fields[i]

				assert.Equal(t, expected.Name, field.Name, fmt.Sprintf("Name in field %+v", field))
				assert.Equal(t, expected.Type, field.Type, fmt.Sprintf("Type in field %+v", field))
				assert.Equal(t, expected.IsPtr, field.IsPtr, fmt.Sprintf("IsPtr in field %+v", field))
				assert.Equal(t, expected.IsSlice, field.IsSlice, fmt.Sprintf("IsSlice in field %+v", field))
				assert.Equal(t, expected.IsMap, field.IsMap, fmt.Sprintf("IsMap in field %+v", field))
				assert.Equal(t, expected.IsStruct, field.IsStruct, fmt.Sprintf("IsStruct in field %+v", field))
				assert.Equal(t, expected.StructName, field.StructName, fmt.Sprintf("StructName in field %+v", field))
				assert.Equal(t, expected.IsBasic, field.IsBasic, fmt.Sprintf("IsBasic in field %+v", field))
			}
		})
	}
}

func TestIsStructType(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"bool", false},
		{"string", false},
		{"int", false},
		{"MyStruct", true},
		{"SubStruct", true},
		{"time.Time", true},
		{"*int", true},
		{"[]string", true},
	}

	for _, tc := range tests {
		t.Run(tc.typeName, func(t *testing.T) {
			result := isStructType(tc.typeName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsBasicType(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"bool", true},
		{"string", true},
		{"int", true},
		{"int8", true},
		{"int16", true},
		{"int32", true},
		{"int64", true},
		{"uint", true},
		{"uint8", true},
		{"uint16", true},
		{"uint32", true},
		{"uint64", true},
		{"float32", true},
		{"float64", true},
		{"byte", true},
		{"rune", true},
		{"complex64", true},
		{"complex128", true},
		{"uintptr", true},
		{"error", true},
		{"MyStruct", false},
		{"time.Time", false},
		{"*int", false},
		{"[]string", false},
		{"map[string]int", false},
	}

	for _, tc := range tests {
		t.Run(tc.typeName, func(t *testing.T) {
			result := isBasicType(tc.typeName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func cleanGeneratedFiles(t *testing.T, dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == resetFile {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				t.Logf("Warning: failed to remove %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		t.Logf("Warning: error during cleanup: %v", err)
	}
}
