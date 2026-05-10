package mapset

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// ListInterfaces parses the Go source file at filePath and returns one string
// per interface type declaration, formatted as a complete type...interface{...} block.
func ListInterfaces(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
				continue
			}
			var buf bytes.Buffer
			err := format.Node(&buf, fset, &ast.GenDecl{
				Tok:   token.TYPE,
				Specs: []ast.Spec{typeSpec},
			})
			if err != nil {
				return nil, err
			}
			result = append(result, stripBlankLines(buf.String()))
		}
	}
	return result, nil
}

func stripBlankLines(s string) string {
	lines := strings.Split(s, "\n")
	out := lines[:0]
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
