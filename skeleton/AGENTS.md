# AGENTS.md — skeleton/

## Purpose

This directory contains the `opencode-skeleton` Go binary. It accepts a file path as a CLI argument and prints a compact skeleton (symbol signatures with line numbers) to stdout.

---

## Key Rules

- **Pure stdlib only** — no CGO, no external dependencies, no `go.sum` with third-party modules
- Go files are parsed with `go/ast` for accurate results
- Other languages (TypeScript, JavaScript, Python, PHP) use regex-based parsers

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0`  | Success — skeleton printed to stdout |
| `1`  | Real error (file not found, parse failure, etc.) — error printed to stderr |
| `2`  | Unsupported file extension — NOT an error; caller should fall through to normal behavior |

**Exit code 2 is NOT an error.** The plugin uses it to decide whether to pass-through to the original `read` tool.

---

## Code Structure

```
skeleton/
  main.go                        ← CLI entry point; reads path arg, calls detect.ForFile()
  internal/
    detect.go                    ← ForFile() — dispatches by file extension to the right parser
    parser/
      parser.go                  ← Parser interface definition
      golang.go                  ← Go parser (uses go/ast)
      typescript.go              ← TS/JS parser (regex); contains lineNumberAt() helper
      python.go                  ← Python parser (regex)
      php.go                     ← PHP parser (regex)
```

---

## Parser Interface

Every language parser implements the `Parser` interface defined in `parser.go`:

```go
type Parser interface {
    Parse(src []byte) ([]Symbol, error)
}

type Symbol struct {
    Name    string
    Kind    string // "func", "class", "interface", "type", etc.
    Line    int
    Indent  int    // nesting level
}
```

---

## Adding a New Language

1. Create `internal/parser/<language>.go` implementing `Parser`
2. Register the extension(s) in `internal/detect.go` → `ForFile()` switch/map
3. Exit with code 2 for any extension not handled in `ForFile()`

---

## Shared Utilities

- `lineNumberAt(src []byte, offset int) int` in `typescript.go` is a shared helper used by all regex-based parsers in the package to convert a byte offset to a line number. Import it within the `parser` package — do not duplicate it.

---

## Building and Testing

```bash
make build     # builds ./bin/opencode-skeleton
make install   # installs to $GOPATH/bin
make test      # runs go test ./...
```

Tests live in `internal/parser/*_test.go`. Each parser should have test cases covering:
- Basic symbol extraction
- Nested symbols (classes with methods)
- Edge cases (empty file, comment-only file)
