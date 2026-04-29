package parser

// JavaScriptParser wraps TypeScriptParser for JS/JSX files
type JavaScriptParser struct {
	ts TypeScriptParser
}

func (j *JavaScriptParser) Parse(content string) ([]Symbol, error) {
	return j.ts.Parse(content)
}
