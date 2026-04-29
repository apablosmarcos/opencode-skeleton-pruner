# AGENTS.md — opencode-context-pruner

## Project Overview

`opencode-context-pruner` is an OpenCode plugin that overrides the built-in `read` tool. When the LLM reads a source file larger than 200 lines without specifying an offset, it receives a compact skeleton (function/class signatures with line numbers) instead of the full file. This reduces context token usage by 70–95% on large files.

The project has two components: a Go binary (`opencode-skeleton`) that does the actual parsing, and a TypeScript npm plugin (`opencode-context-pruner`) that hooks into OpenCode and calls the binary.

---

## Architecture Decisions

| Decision | Detail |
|----------|--------|
| Override `read` | The plugin shadows OpenCode's built-in `read` tool entirely |
| Threshold | 200 lines — configurable via `PRUNER_THRESHOLD` env var |
| Go binary optional | The binary provides best results; the plugin includes a TS regex fallback |
| TS fallback | Pure regex, no external deps — always works even without the binary |
| Pass-through | If `offset` or `limit` is given, skip pruning entirely (LLM is expanding a section) |
| Exit code 2 | Binary exits with code 2 for unsupported extensions — NOT an error, plugin falls through to normal read |

---

## Repository Structure

```
opencode-context-pruner/
  skeleton/               ← Go binary source
    main.go
    internal/
      parser/             ← Language parsers (Parser interface + per-language files)
      detect.go           ← ForFile() dispatches by extension
  plugin/                 ← TypeScript OpenCode plugin
    src/
      index.ts            ← Plugin entry point
      smart-read.ts       ← Core override logic
      ts-fallback.ts      ← Pure regex fallback parser
    package.json
```

---

## Building Both Components

### Go binary

```bash
cd skeleton
make build       # output: ./bin/opencode-skeleton
make install     # installs to $GOPATH/bin
make test
```

### TypeScript plugin

```bash
cd plugin
npm install
npm run build
npm test
```

---

## Adding a New Language Parser

1. Create `skeleton/internal/parser/<language>.go` implementing the `Parser` interface from `parser.go`
2. Register the new extension(s) in `skeleton/internal/detect.go` → `ForFile()` function
3. Add the extension to the supported list in `plugin/src/smart-read.ts` so the plugin knows when to attempt the binary call

---

## Important Notes for Agents

- **Exit code 2** from the binary means "unsupported file extension". The plugin MUST treat this as a pass-through (call the original `read`), not as an error.
- **Exit code 1** = real error (parse failure, file not found, etc.)
- **Exit code 0** = success, stdout contains the skeleton
- The plugin uses `@opencode-ai/plugin@1.14.29` with **Zod** for argument schemas — NOT JSON Schema
- Always pass-through when `offset` or `limit` args are present — the LLM is doing a targeted read
- `PRUNER_VERBOSE=1` logs to stderr; never write pruning output to stdout (it would corrupt the tool response)
