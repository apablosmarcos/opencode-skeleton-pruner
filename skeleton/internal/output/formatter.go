package output

import (
	"fmt"
	"strings"

	"github.com/apablosmarcos/opencode-skeleton-pruner/skeleton/internal/parser"
)

type Options struct {
	FilePath      string
	OriginalLines int
	JSON          bool
}

func Format(symbols []parser.Symbol, opts Options) string {
	// Header
	skeletonLines := len(symbols) + 3 // approximate
	reduction := 0
	if opts.OriginalLines > 0 {
		reduction = 100 - (skeletonLines * 100 / opts.OriginalLines)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Skeleton: %s\n", opts.FilePath))
	sb.WriteString(fmt.Sprintf("# Original: %d lines → Skeleton: ~%d symbols (%d%% reduction)\n",
		opts.OriginalLines, len(symbols), reduction))
	sb.WriteString("# Use read(filePath, offset=N) to expand any symbol\n\n")

	for _, s := range symbols {
		line := formatSymbol(s)
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatSymbol(s parser.Symbol) string {
	prefix := ""
	if s.Comment != "" {
		prefix = fmt.Sprintf("// %s\n", s.Comment)
	}

	var sig string
	switch s.Kind {
	case "import_block":
		sig = fmt.Sprintf("import (...) [line %d]", s.Line)
	case "class", "interface", "trait":
		sig = fmt.Sprintf("%s %s [line %d]", s.Kind, s.Name, s.Line)
	case "type":
		sig = fmt.Sprintf("type %s [line %d]", s.Name, s.Line)
	case "func":
		if s.Returns != "" {
			sig = fmt.Sprintf("func %s(%s) %s [line %d]", s.Name, s.Params, s.Returns, s.Line)
		} else {
			sig = fmt.Sprintf("func %s(%s) [line %d]", s.Name, s.Params, s.Line)
		}
	case "method":
		recv := ""
		if s.Receiver != "" {
			recv = fmt.Sprintf("(%s) ", s.Receiver)
		}
		if s.Returns != "" {
			sig = fmt.Sprintf("func %s%s(%s) %s [line %d]", recv, s.Name, s.Params, s.Returns, s.Line)
		} else {
			sig = fmt.Sprintf("func %s%s(%s) [line %d]", recv, s.Name, s.Params, s.Line)
		}
	default:
		sig = fmt.Sprintf("%s %s [line %d]", s.Kind, s.Name, s.Line)
	}
	return prefix + sig
}
