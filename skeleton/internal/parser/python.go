package parser

import (
	"regexp"
	"strings"
)

// PythonParser parses Python source files using regex
type PythonParser struct{}

var (
	pyClassRe  = regexp.MustCompile(`(?m)^class\s+(\w+)`)
	pyFuncRe   = regexp.MustCompile(`(?m)^(async\s+)?def\s+(\w+)\s*\(([^)]*)\)`)
	pyMethRe   = regexp.MustCompile(`(?m)^( {4,}|\t+)(async\s+)?def\s+(\w+)\s*\(([^)]*)\)`)
	pyDecoRe   = regexp.MustCompile(`@(\w+)`)
)

func (p *PythonParser) Parse(content string) ([]Symbol, error) {
	var symbols []Symbol
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Class
		if m := pyClassRe.FindStringSubmatch(trimmed); m != nil && strings.HasPrefix(line, "class ") {
			symbols = append(symbols, Symbol{Kind: "class", Name: m[1], Line: lineNum})
			continue
		}

		// Module-level func/async func
		if strings.HasPrefix(line, "def ") || strings.HasPrefix(line, "async def ") {
			m := pyFuncRe.FindStringSubmatch(line)
			if m != nil {
				kind := "func"
				if m[1] != "" {
					kind = "async_func"
				}
				comment := ""
				if i > 0 {
					prev := strings.TrimSpace(lines[i-1])
					if dm := pyDecoRe.FindStringSubmatch(prev); dm != nil {
						comment = "@" + dm[1]
					}
				}
				symbols = append(symbols, Symbol{Kind: kind, Name: m[2], Params: m[3], Line: lineNum, Comment: comment})
			}
			continue
		}

		// Indented method
		if m := pyMethRe.FindStringSubmatch(line); m != nil {
			kind := "method"
			if m[2] != "" {
				kind = "async_method"
			}
			comment := ""
			if i > 0 {
				prev := strings.TrimSpace(lines[i-1])
				if dm := pyDecoRe.FindStringSubmatch(prev); dm != nil {
					comment = "@" + dm[1]
				}
			}
			symbols = append(symbols, Symbol{Kind: kind, Name: m[3], Params: m[4], Line: lineNum, Comment: comment})
		}
	}

	return symbols, nil
}
