import { execFileSync } from 'node:child_process'
import { existsSync } from 'node:fs'

const BINARY_NAME = 'opencode-skeleton'

function findBinary(): string | null {
  // Check PATH locations
  const paths = [
    `${process.env.HOME}/.local/bin/${BINARY_NAME}`,
    `/usr/local/bin/${BINARY_NAME}`,
    `/usr/bin/${BINARY_NAME}`,
  ]
  for (const p of paths) {
    if (existsSync(p)) return p
  }
  // Try which via execFileSync
  try {
    const result = execFileSync('which', [BINARY_NAME], { encoding: 'utf8', timeout: 2000 })
    const p = result.trim()
    if (p && existsSync(p)) return p
  } catch {
    // which not found or binary not in PATH
  }
  return null
}

let _binaryPath: string | null | undefined = undefined

export function getBinaryPath(): string | null {
  if (_binaryPath === undefined) {
    _binaryPath = findBinary()
  }
  return _binaryPath
}

export function runGoBinary(filePath: string): string | null {
  const bin = getBinaryPath()
  if (!bin) return null

  try {
    const output = execFileSync(bin, [filePath], {
      encoding: 'utf8',
      timeout: 10000,
      maxBuffer: 1024 * 1024, // 1MB
    })
    return output
  } catch (err: unknown) {
    // exit code 2 = unsupported extension (expected)
    // other error = binary failed, fall through to TS fallback
    if (err && typeof err === 'object' && 'status' in err && (err as NodeJS.ErrnoException & { status?: number }).status === 2) {
      return null
    }
    return null
  }
}
