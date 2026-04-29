package parser

import (
	"regexp"
)

// PHPParser parses PHP source files using regex
type PHPParser struct{}

var (
	phpNamespaceRe = regexp.MustCompile(`(?m)^namespace\s+([\w\\]+);`)
	phpClassRe     = regexp.MustCompile(`(?m)^(?:abstract\s+|final\s+)?class\s+(\w+)`)
	phpInterfaceRe = regexp.MustCompile(`(?m)^interface\s+(\w+)`)
	phpTraitRe     = regexp.MustCompile(`(?m)^trait\s+(\w+)`)
	phpFuncRe      = regexp.MustCompile(`(?m)^function\s+(\w+)\s*\(([^)]*)\)`)
	phpMethodRe    = regexp.MustCompile(`(?m)^\s+(?:public|protected|private)(?:\s+static)?\s+function\s+(\w+)\s*\(([^)]*)\)`)
)

func (p *PHPParser) Parse(content string) ([]Symbol, error) {
	var symbols []Symbol

	type entry struct {
		offset int
		sym    Symbol
	}
	var entries []entry

	add := func(re *regexp.Regexp, handler func(m []int) Symbol) {
		for _, loc := range re.FindAllStringSubmatchIndex(content, -1) {
			s := handler(loc)
			if s.Kind != "" {
				entries = append(entries, entry{offset: loc[0], sym: s})
			}
		}
	}

	add(phpNamespaceRe, func(m []int) Symbol {
		name := content[m[2]:m[3]]
		return Symbol{Kind: "namespace", Name: name, Line: lineNumberAt(content, m[0])}
	})
	add(phpClassRe, func(m []int) Symbol {
		name := content[m[2]:m[3]]
		return Symbol{Kind: "class", Name: name, Line: lineNumberAt(content, m[0])}
	})
	add(phpInterfaceRe, func(m []int) Symbol {
		name := content[m[2]:m[3]]
		return Symbol{Kind: "interface", Name: name, Line: lineNumberAt(content, m[0])}
	})
	add(phpTraitRe, func(m []int) Symbol {
		name := content[m[2]:m[3]]
		return Symbol{Kind: "trait", Name: name, Line: lineNumberAt(content, m[0])}
	})
	add(phpFuncRe, func(m []int) Symbol {
		name := content[m[2]:m[3]]
		params := content[m[4]:m[5]]
		return Symbol{Kind: "func", Name: name, Params: params, Line: lineNumberAt(content, m[0])}
	})
	add(phpMethodRe, func(m []int) Symbol {
		name := content[m[2]:m[3]]
		params := content[m[4]:m[5]]
		return Symbol{Kind: "method", Name: name, Params: params, Line: lineNumberAt(content, m[0])}
	})

	// sort by offset
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].offset < entries[j-1].offset; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
	for _, e := range entries {
		symbols = append(symbols, e.sym)
	}
	return symbols, nil
}
