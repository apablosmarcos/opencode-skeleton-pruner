# AGENTS.md — plugin/

## Purpose

This directory contains the OpenCode plugin (`opencode-context-pruner`). It overrides the built-in `read` tool with a context-pruning version that returns file skeletons instead of full content for large files.

---

## Plugin Architecture

- **Entry point**: `src/index.ts` — exports a function `(input) => Promise<Hooks>`
- **Hooks shape**: `Hooks.tool` (singular, not `tools`) is where tools are registered
- **SDK**: `@opencode-ai/plugin` — import from this package; use `tool.schema` for Zod arg definitions
- **Core logic**: `src/smart-read.ts` — the overridden `read` tool implementation
- **Fallback parser**: `src/ts-fallback.ts` — pure regex skeleton generator, no external dependencies

---

## Key Rules for This Package

### smart-read.ts

- **Always pass-through** when `offset` or `limit` args are given — the LLM is expanding a specific section, do not intercept
- Call the Go binary first; if it exits with code 2 (unsupported extension), fall through to normal read
- If the binary is not found in `$PATH`, use `ts-fallback.ts` automatically
- Never swallow real errors (exit code 1 from binary = surface the error)

### ts-fallback.ts

- Pure regex, zero external dependencies — keep it that way
- Must remain self-contained; do not import from other project files except types
- Used as a last resort when the binary is unavailable

### Zod Schemas

```ts
import { tool } from "@opencode-ai/plugin";

// Use tool.schema for Zod arg definitions — NOT JSON Schema objects
const args = tool.schema({
  path: z.string(),
  offset: z.number().optional(),
  limit: z.number().optional(),
});
```

---

## Environment Variables

| Variable           | Default | Effect                                         |
|--------------------|---------|------------------------------------------------|
| `PRUNER_THRESHOLD` | `200`   | Line count above which skeleton mode activates |
| `PRUNER_VERBOSE`   | `0`     | `1` = log pruning decisions to stderr          |

---

## Tests

Place tests in `src/__tests__/` or alongside source files as `*.test.ts`. Run with:

```bash
npm test
```

Focus test coverage on:
- Pass-through when offset/limit given
- Skeleton returned when file exceeds threshold
- Fallback to TS parser when binary missing
- Pass-through for unsupported extensions (binary exit code 2)
