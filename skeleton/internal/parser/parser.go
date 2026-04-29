package parser

// Symbol represents a parsed code symbol (function, method, class, etc.)
type Symbol struct {
	Kind     string // "func", "method", "class", "interface", "type", "const", "var", "import_block"
	Name     string
	Receiver string // for methods: "(*MyStruct)"
	Params   string // function params as-is
	Returns  string // return types as-is
	Line     int    // 1-based line number in original file
	Comment  string // leading doc comment (first line only)
}

// Parser parses source code and returns a list of symbols
type Parser interface {
	Parse(content string) ([]Symbol, error)
}
