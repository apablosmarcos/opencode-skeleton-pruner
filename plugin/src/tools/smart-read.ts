import { tool } from '@opencode-ai/plugin'
import { readFileSync } from 'node:fs'
import { extname } from 'node:path'
import { loadConfig } from '../config.js'
import { generateSkeleton } from '../skeleton/index.js'

export const smartReadTool = tool({
  description:
    'Read a file. For large source files, returns a compact skeleton with line numbers to expand specific sections.',
  args: {
    filePath: tool.schema.string().describe('Absolute path to the file'),
    offset: tool.schema
      .number()
      .optional()
      .describe('Line number to start reading from (1-indexed)'),
    limit: tool.schema.number().optional().describe('Maximum number of lines to read'),
  },
  async execute({ filePath, offset, limit }) {
    const config = loadConfig()

    // If offset/limit requested → always pass through (user wants specific section)
    if (offset !== undefined || limit !== undefined) {
      return readWithOffsetLimit(filePath, offset, limit)
    }

    // Check extension
    const ext = extname(filePath).toLowerCase()
    if (!config.extensions.includes(ext)) {
      return readFull(filePath)
    }

    // Read and check line count
    let content: string
    try {
      content = readFileSync(filePath, 'utf8')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err)
      throw new Error(`Cannot read file: ${msg}`)
    }

    const lineCount = content.split('\n').length
    if (lineCount <= config.threshold) {
      return content
    }

    // Large file → generate skeleton
    const result = await generateSkeleton(filePath, content)
    if (config.verbose) {
      process.stderr.write(
        `[context-pruner] ${filePath}: ${result.originalLines} lines → skeleton (${result.method})\n`,
      )
    }
    return result.skeleton
  },
})

function readFull(filePath: string): string {
  return readFileSync(filePath, 'utf8')
}

function readWithOffsetLimit(filePath: string, offset?: number, limit?: number): string {
  const content = readFileSync(filePath, 'utf8')
  const lines = content.split('\n')
  const start = offset ? offset - 1 : 0
  const end = limit ? start + limit : lines.length
  return lines.slice(start, end).join('\n')
}
