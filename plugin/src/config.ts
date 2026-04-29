export interface ContextPrunerConfig {
  /** Lines threshold above which skeleton is returned instead of full content */
  threshold: number
  /** File extensions to prune (others pass through unchanged) */
  extensions: string[]
  /** Glob patterns to always pass through (never prune) */
  passthrough: string[]
  /** Whether to log pruning activity to stderr */
  verbose: boolean
}

export const DEFAULT_CONFIG: ContextPrunerConfig = {
  threshold: 200,
  extensions: ['.ts', '.tsx', '.js', '.jsx', '.mjs', '.cjs', '.go', '.py', '.php'],
  passthrough: ['**/node_modules/**', '**/*.test.*', '**/*.spec.*', '**/*.d.ts'],
  verbose: false,
}

export function loadConfig(): ContextPrunerConfig {
  // Allow overriding via env vars
  const threshold = process.env.PRUNER_THRESHOLD
    ? parseInt(process.env.PRUNER_THRESHOLD, 10)
    : DEFAULT_CONFIG.threshold
  const verbose = process.env.PRUNER_VERBOSE === '1'

  return { ...DEFAULT_CONFIG, threshold, verbose }
}
