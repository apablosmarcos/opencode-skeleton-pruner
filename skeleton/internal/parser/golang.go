package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// GoParser parses Go source files using go/ast
type GoParser struct{}

func (g *GoParser) Parse(content string) ([]Symbol, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("go parse error: %w", err)
	}

	lines := strings.Split(content, "\n")
	_ = lines

	var symbols []Symbol

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			pos := fset.Position(d.TokPos)
			switch d.Tok.String() {
			case "import":
				symbols = append(symbols, Symbol{
					Kind: "import_block",
					Line: pos.Line,
				})
			case "type":
				for _, spec := range d.Specs {
					ts := spec.(*ast.TypeSpec)
					kind := "type"
					switch ts.Type.(type) {
					case *ast.InterfaceType:
						kind = "interface"
					case *ast.StructType:
						kind = "class"
					}
					comment := ""
					if d.Doc != nil && len(d.Doc.List) > 0 {
						comment = strings.TrimPrefix(d.Doc.List[0].Text, "// ")
					}
					symbols = append(symbols, Symbol{
						Kind:    kind,
						Name:    ts.Name.Name,
						Line:    fset.Position(ts.Name.NamePos).Line,
						Comment: comment,
					})
				}
			case "const", "var":
				n := len(d.Specs)
				kind := d.Tok.String()
				name := fmt.Sprintf("%s block (%d items)", kind, n)
				if n == 1 {
					if vs, ok := d.Specs[0].(*ast.ValueSpec); ok && len(vs.Names) == 1 {
						name = vs.Names[0].Name
					}
				}
				symbols = append(symbols, Symbol{
					Kind: kind,
					Name: name,
					Line: pos.Line,
				})
			}

		case *ast.FuncDecl:
			pos := fset.Position(d.Name.NamePos)
			comment := ""
			if d.Doc != nil && len(d.Doc.List) > 0 {
				comment = strings.TrimPrefix(d.Doc.List[0].Text, "// ")
			}

			params := formatFieldList(d.Type.Params, content, fset)
			returns := formatFieldList(d.Type.Results, content, fset)

			if d.Recv != nil && len(d.Recv.List) > 0 {
				recv := formatReceiver(d.Recv.List[0], content, fset)
				symbols = append(symbols, Symbol{
					Kind:     "method",
					Name:     d.Name.Name,
					Receiver: recv,
					Params:   params,
					Returns:  returns,
					Line:     pos.Line,
					Comment:  comment,
				})
			} else {
				symbols = append(symbols, Symbol{
					Kind:    "func",
					Name:    d.Name.Name,
					Params:  params,
					Returns: returns,
					Line:    pos.Line,
					Comment: comment,
				})
			}
		}
	}

	return symbols, nil
}

func formatReceiver(field *ast.Field, content string, fset *token.FileSet) string {
	start := fset.Position(field.Type.Pos()).Offset
	end := fset.Position(field.Type.End()).Offset
	if start >= 0 && end <= len(content) && start < end {
		return content[start:end]
	}
	return ""
}

func formatFieldList(fl *ast.FieldList, content string, fset *token.FileSet) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}
	start := fset.Position(fl.Opening).Offset + 1 // skip '('
	end := fset.Position(fl.Closing).Offset
	if fl.Opening == token.NoPos {
		// results without parens: single return type
		start = fset.Position(fl.List[0].Type.Pos()).Offset
		end = fset.Position(fl.List[len(fl.List)-1].Type.End()).Offset
	}
	if start >= 0 && end <= len(content) && start < end {
		return strings.TrimSpace(content[start:end])
	}
	return ""
}
