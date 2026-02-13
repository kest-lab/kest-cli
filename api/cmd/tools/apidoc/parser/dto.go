package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ParseDTOs parses all DTO definitions from the app directory
func ParseDTOs(appDir string) (map[string]*DTO, error) {
	dtos := make(map[string]*DTO)

	err := filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only parse dto.go files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// Parse the file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files that can't be parsed
		}

		// Get module name from directory
		module := filepath.Base(filepath.Dir(path))

		// Find all struct definitions
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				dto := &DTO{
					Name:   typeSpec.Name.Name,
					Module: module,
					Type:   classifyDTO(typeSpec.Name.Name),
					Fields: parseFields(structType),
				}

				// Get doc comment
				if genDecl.Doc != nil {
					dto.Comment = strings.TrimSpace(genDecl.Doc.Text())
				}

				// Use full key: module.DTOName
				key := module + "." + dto.Name
				dtos[key] = dto
			}
		}

		return nil
	})

	return dtos, err
}

// classifyDTO determines if a DTO is a request, response, or model
func classifyDTO(name string) string {
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "request") {
		return "request"
	}
	if strings.Contains(nameLower, "response") {
		return "response"
	}
	return "model"
}

// parseFields parses struct fields into DTOField slice
func parseFields(structType *ast.StructType) []DTOField {
	var fields []DTOField

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue // Skip embedded fields
		}

		dtoField := DTOField{
			Name: field.Names[0].Name,
			Type: typeToString(field.Type),
		}

		// Parse struct tags
		if field.Tag != nil {
			tag := field.Tag.Value
			dtoField.JSONName = extractTag(tag, "json")
			dtoField.Binding = extractTag(tag, "binding")
			dtoField.Required = strings.Contains(dtoField.Binding, "required")
			dtoField.Validation = formatValidation(dtoField.Binding)
		}

		// Get field comment
		if field.Comment != nil {
			dtoField.Comment = strings.TrimSpace(field.Comment.Text())
		} else if field.Doc != nil {
			dtoField.Comment = strings.TrimSpace(field.Doc.Text())
		}

		fields = append(fields, dtoField)
	}

	return fields
}

// typeToString converts an ast.Expr type to string
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "any"
	default:
		return "unknown"
	}
}

// extractTag extracts a specific tag value from struct tag
func extractTag(tag, key string) string {
	// Remove backticks
	tag = strings.Trim(tag, "`")

	// Find the key
	pattern := regexp.MustCompile(key + `:"([^"]*)"`)
	matches := pattern.FindStringSubmatch(tag)
	if len(matches) > 1 {
		// Return first part before comma for json tag
		if key == "json" {
			parts := strings.Split(matches[1], ",")
			return parts[0]
		}
		return matches[1]
	}
	return ""
}

// formatValidation formats binding rules into human-readable validation
func formatValidation(binding string) string {
	if binding == "" {
		return ""
	}

	var rules []string
	parts := strings.Split(binding, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case part == "required":
			rules = append(rules, "Required")
		case strings.HasPrefix(part, "min="):
			rules = append(rules, "Min: "+strings.TrimPrefix(part, "min="))
		case strings.HasPrefix(part, "max="):
			rules = append(rules, "Max: "+strings.TrimPrefix(part, "max="))
		case part == "email":
			rules = append(rules, "Email format")
		case part == "omitempty":
			rules = append(rules, "Optional")
		case strings.HasPrefix(part, "oneof="):
			rules = append(rules, "One of: "+strings.TrimPrefix(part, "oneof="))
		}
	}

	return strings.Join(rules, ", ")
}
