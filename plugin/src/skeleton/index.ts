import { runGoBinary } from './go-binary.js'
import { parseTsskeleton, formatTsSkeleton } from './ts-fallback.js'

export interface SkeletonResult {
  skeleton: string
  method: 'go-binary' | 'ts-fallback'
  originalLines: number
  symbolCount: number
}

export async function generateSkeleton(filePath: string, content: string): Promise<SkeletonResult> {
  const originalLines = content.split('\n').length

  // Try Go binary first
  const goBinaryResult = runGoBinary(filePath)
  if (goBinaryResult) {
    const symbolCount = (goBinaryResult.match(/\[line \d+\]/g) ?? []).length
    return { skeleton: goBinaryResult, method: 'go-binary', originalLines, symbolCount }
  }

  // Fallback: pure TS regex parser
  const symbols = parseTsskeleton(content, filePath)
  const skeleton = formatTsSkeleton(symbols, filePath, originalLines)
  return { skeleton, method: 'ts-fallback', originalLines, symbolCount: symbols.length }
}
