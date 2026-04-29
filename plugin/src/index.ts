import type { Plugin } from '@opencode-ai/plugin'
import { smartReadTool } from './tools/smart-read.js'

/**
 * opencode-context-pruner plugin
 *
 * Overrides the built-in `read` tool to return compact skeletons for large
 * source files, significantly reducing token usage when the AI needs to
 * understand structure without reading every line.
 */
const plugin: Plugin = async (_input) => {
  return {
    tool: {
      read: smartReadTool,
    },
  }
}

export default plugin
