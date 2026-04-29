package parser

import (
	"regexp"
	"strings"
)

// TypeScriptParser parses TypeScript/TSX files using regex
type TypeScriptParser struct{}

var (
	tsClassRe     = regexp.MustCompile(`(?m)^(export\s+)?(abstract\s+)?class\s+(\w+)`)
	tsInterfaceRe = regexp.MustCompile(`(?m)^(export\s+)?interface\s+(\w+)`)
	tsTypeRe      = regexp.MustCompile(`(?m)^(export\s+)?type\s+(\w+)\s*=`)
	tsEnumRe      = regexp.MustCompile(`(?m)^(export\s+)?enum\s+(\w+)`)
	tsFuncRe      = regexp.MustCompile(`(?m)^(export\s+)?(?:async\s+)?function\s+(\w+)\s*\(([^)]*)\)`)
	tsArrowRe     = regexp.MustCompile(`(?m)^(export\s+)?const\s+(\w+)\s*=\s*(?:async\s*)?\(([^)]*)\)\s*=>`)
	tsMethodRe    = regexp.MustCompile(`(?m)^(\s{2,})(?:public\s+|private\s+|protected\s+|static\s+|async\s+)*(\w+)\s*\(([^)]*)\)`)
)

func lineNumberAt(content string, offset int) int {
	return strings.Count(content[:offset], "\n") + 1
}

func (t *TypeScriptParser) Parse(content string) ([]Symbol, error) {
	var symbols []Symbol

	addMatches := func(re *regexp.Regexp, handler func(match []int) Symbol) {
		for _, loc := range re.FindAllStringSubmatchIndex(content, -1) {
			s := handler(loc)
			if s.Kind != "" {
				symbols = append(symbols, s)
			}
		}
	}

	// classes
	addMatches(tsClassRe, func(m []int) Symbol {
		name := content[m[6]:m[7]]
		return Symbol{Kind: "class", Name: name, Line: lineNumberAt(content, m[0])}
	})
	// interfaces
	addMatches(tsInterfaceRe, func(m []int) Symbol {
		name := content[m[4]:m[5]]
		return Symbol{Kind: "interface", Name: name, Line: lineNumberAt(content, m[0])}
	})
	// types
	addMatches(tsTypeRe, func(m []int) Symbol {
		name := content[m[4]:m[5]]
		return Symbol{Kind: "type", Name: name, Line: lineNumberAt(content, m[0])}
	})
	// enums
	addMatches(tsEnumRe, func(m []int) Symbol {
		name := content[m[4]:m[5]]
		return Symbol{Kind: "type", Name: name, Line: lineNumberAt(content, m[0])}
	})
	// functions
	addMatches(tsFuncRe, func(m []int) Symbol {
		name := content[m[4]:m[5]]
		params := content[m[6]:m[7]]
		return Symbol{Kind: "func", Name: name, Params: params, Line: lineNumberAt(content, m[0])}
	})
	// arrow functions
	addMatches(tsArrowRe, func(m []int) Symbol {
		name := content[m[4]:m[5]]
		params := content[m[6]:m[7]]
		return Symbol{Kind: "func", Name: name, Params: params, Line: lineNumberAt(content, m[0])}
	})
	// methods (indented)
	addMatches(tsMethodRe, func(m []int) Symbol {
		name := content[m[4]:m[5]]
		params := content[m[6]:m[7]]
		// skip common false positives
		if name == "if" || name == "for" || name == "while" || name == "switch" || name == "return" {
			return Symbol{}
		}
		return Symbol{Kind: "method", Name: name, Params: params, Line: lineNumberAt(content, m[0])}
	})

	return symbols, nil
}
