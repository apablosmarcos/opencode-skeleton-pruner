package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/apablosmarcos/opencode-skeleton-pruner/skeleton/internal/output"
	"github.com/apablosmarcos/opencode-skeleton-pruner/skeleton/internal/parser"
)

func main() {
	args := os.Args[1:]
	jsonMode := false

	// parse flags
	filtered := args[:0]
	for _, a := range args {
		if a == "--json" {
			jsonMode = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: opencode-skeleton [--json] <file>")
		os.Exit(1)
	}

	filePath := args[0]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	lineCount := countLines(string(content))

	p, err := parser.ForFile(filePath)
	if err != nil {
		// Unsupported extension — print a message to stderr, exit 2 (signal to caller: pass-through)
		fmt.Fprintf(os.Stderr, "unsupported: %v\n", err)
		os.Exit(2)
	}

	symbols, err := p.Parse(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	opts := output.Options{
		FilePath:      filePath,
		OriginalLines: lineCount,
		JSON:          jsonMode,
	}

	fmt.Print(output.Format(symbols, opts))
}

func countLines(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	n := 0
	for scanner.Scan() {
		n++
	}
	return n
}
