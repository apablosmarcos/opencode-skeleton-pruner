package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ForFile(path string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return &GoParser{}, nil
	case ".ts", ".tsx":
		return &TypeScriptParser{}, nil
	case ".js", ".jsx", ".mjs", ".cjs":
		return &JavaScriptParser{}, nil
	case ".py":
		return &PythonParser{}, nil
	case ".php":
		return &PHPParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported extension: %s", ext)
	}
}
