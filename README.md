# opencode-skeleton-pruner

> OpenCode plugin that replaces large file reads with compact skeletons — reducing context token usage by up to 98%.

[![npm](https://img.shields.io/npm/v/opencode-skeleton-pruner)](https://www.npmjs.com/package/opencode-skeleton-pruner)
[![license](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

---

## The Problem

When an LLM reads a 2,000-line source file just to find one function, it consumes thousands of tokens on code it doesn't need. Multiply that by a full session and the cost — in tokens and latency — adds up fast.

## The Solution

`opencode-skeleton-pruner` intercepts every `read` call. When a file is larger than the configured threshold (default: 200 lines), instead of returning the full content it returns a **skeleton** — a compact list of class/function signatures with line numbers.

The LLM then calls `read(filePath, offset=N)` to expand only the section it actually needs.

```
Full file:  2,393 lines  →  Skeleton: 60 symbols  →  98% token reduction
```

---

## How It Works

1. Plugin overrides the built-in `read` tool
2. File below threshold → full content (unchanged)
3. File above threshold + no `offset`/`limit` → skeleton
4. LLM reads symbol map, identifies the relevant section
5. LLM calls `read(filePath, offset=N)` to expand just that section

---

## Skeleton Output Example

Instead of 342 lines of source code, the LLM sees:

```
# Skeleton: src/UserService.ts
# Original: 342 lines → Skeleton: ~12 symbols (97% reduction)
# Use read(filePath, offset=N) to expand any symbol

class UserService [line 8]
  constructor(...) [line 15]
  findById(id: string): Promise<User> [line 22]
  create(dto: CreateUserDto): Promise<User> [line 38]
  update(id: string, dto: UpdateUserDto): Promise<User> [line 57]
  delete(id: string): Promise<void> [line 74]

interface UserRepository [line 91]
type CreateUserDto = ... [line 112]
type UpdateUserDto = ... [line 118]
```

---

## Installation

### Step 1 — Install the Go binary (recommended)

The Go binary uses `go/ast` for accurate Go parsing and provides better results for TypeScript/JS/Python/PHP. Without it, the plugin falls back to a pure TypeScript regex parser automatically.

```bash
git clone git@github.com:apablosmarcos/opencode-skeleton-pruner.git
cd opencode-skeleton-pruner/skeleton
make install
```

This installs `opencode-skeleton` to `~/.local/bin`. Make sure `~/.local/bin` is in your `$PATH`.

### Step 2 — Add the plugin to OpenCode

**Global install (applies to all projects):**
```bash
opencode plugin opencode-skeleton-pruner --global
```

**Or place the standalone file directly:**
```bash
cp ~/.config/opencode/plugins/context-pruner.ts ~/.config/opencode/plugins/
```

---

## Configuration

Configure via environment variables (e.g., in your shell profile):

| Variable           | Default | Description                                   |
|--------------------|---------|-----------------------------------------------|
| `PRUNER_THRESHOLD` | `200`   | Minimum line count to trigger skeleton mode   |
| `PRUNER_VERBOSE`   | `0`     | Set to `1` to log pruning activity to stderr  |

---

## Supported Languages

| Language              | Extensions                        | Parser                  |
|-----------------------|-----------------------------------|-------------------------|
| TypeScript / TSX      | `.ts`, `.tsx`                     | Go binary / TS fallback |
| JavaScript / JSX      | `.js`, `.jsx`, `.mjs`, `.cjs`     | Go binary / TS fallback |
| Go                    | `.go`                             | Go binary (`go/ast`)    |
| Python                | `.py`                             | Go binary / TS fallback |
| PHP                   | `.php`                            | Go binary / TS fallback |

Files with unsupported extensions are always passed through unchanged.

---

## Complementary Plugin

**[DCP — `@tarquinen/opencode-dcp`](https://github.com/Opencode-DCP/opencode-dynamic-context-pruning)** works at the conversation history layer (compressing old messages). This plugin works at the file read layer (preventing large reads from entering context in the first place). They solve different problems and work well together:

```bash
opencode plugin @tarquinen/opencode-dcp --global
opencode plugin opencode-skeleton-pruner --global
```

---

## Development

### Build the Go binary

```bash
cd skeleton
make build    # builds ./build/opencode-skeleton
make install  # installs to ~/.local/bin
make test     # runs Go tests
```

### Build the TypeScript plugin

```bash
cd plugin
bun install
bun run build   # compiles to dist/
bun test        # runs tests
```

### Test locally end-to-end

```bash
PRUNER_VERBOSE=1 opencode
```

---

## License

MIT — see [LICENSE](./LICENSE)
