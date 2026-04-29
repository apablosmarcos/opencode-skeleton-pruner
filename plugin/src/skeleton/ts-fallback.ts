export interface SkeletonSymbol {
  kind: string
  name: string
  receiver?: string
  params?: string
  returns?: string
  line: number
  comment?: string
}

/**
 * Lightweight regex-based skeleton parser for TypeScript/JavaScript.
 * Used when the Go binary is not available.
 */
export function parseTsskeleton(content: string, _filePath: string): SkeletonSymbol[] {
  const lines = content.split('\n')
  const symbols: SkeletonSymbol[] = []

  // Patterns
  const classRe = /^(?:export\s+)?(?:abstract\s+)?class\s+(\w+)/
  const interfaceRe = /^(?:export\s+)?interface\s+(\w+)/
  const typeRe = /^(?:export\s+)?type\s+(\w+)\s*=/
  const enumRe = /^(?:export\s+)?enum\s+(\w+)/
  const funcRe = /^(?:export\s+)?(?:async\s+)?function\s+(\w+)\s*(?:<[^>]*>)?\s*\(([^)]*)\)/
  const arrowRe = /^(?:export\s+)?(?:const|let)\s+(\w+)\s*=\s*(?:async\s+)?\(([^)]*)\)\s*(?::\s*([^=]+))?\s*=>/
  const methodRe = /^  (?:(?:public|private|protected|static|async)\s+)*(\w+)\s*\(([^)]*)\)\s*(?::\s*(.+))?/

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const lineNum = i + 1

    let m: RegExpMatchArray | null

    if ((m = line.match(classRe))) {
      symbols.push({ kind: 'class', name: m[1], line: lineNum })
    } else if ((m = line.match(interfaceRe))) {
      symbols.push({ kind: 'interface', name: m[1], line: lineNum })
    } else if ((m = line.match(typeRe))) {
      symbols.push({ kind: 'type', name: m[1], line: lineNum })
    } else if ((m = line.match(enumRe))) {
      symbols.push({ kind: 'enum', name: m[1], line: lineNum })
    } else if ((m = line.match(funcRe))) {
      const returnMatch = lines[i].match(/\)\s*:\s*([^{]+)/)
      symbols.push({
        kind: 'func',
        name: m[1],
        params: m[2].length > 40 ? '...' : m[2],
        returns: returnMatch ? returnMatch[1].trim() : undefined,
        line: lineNum,
      })
    } else if ((m = line.match(arrowRe))) {
      symbols.push({
        kind: 'func',
        name: m[1],
        params: m[2].length > 40 ? '...' : m[2],
        returns: m[3]?.trim(),
        line: lineNum,
      })
    } else if (line.startsWith('  ') && !line.startsWith('   ') && (m = line.match(methodRe))) {
      // Methods at 2-space indent inside class
      if (!['if', 'for', 'while', 'switch', 'return', 'const', 'let', 'var'].includes(m[1])) {
        symbols.push({
          kind: 'method',
          name: m[1],
          params: m[2].length > 40 ? '...' : m[2],
          returns: m[3]?.trim(),
          line: lineNum,
        })
      }
    }
  }

  return symbols
}

export function formatTsSkeleton(symbols: SkeletonSymbol[], filePath: string, originalLines: number): string {
  const reduction = symbols.length > 0
    ? Math.max(0, 100 - Math.round((symbols.length + 3) * 100 / originalLines))
    : 0

  const lines: string[] = [
    `# Skeleton: ${filePath}`,
    `# Original: ${originalLines} lines → Skeleton: ~${symbols.length} symbols (${reduction}% reduction)`,
    `# Use read(filePath, offset=N) to expand any symbol`,
    '',
  ]

  for (const s of symbols) {
    let sig: string
    switch (s.kind) {
      case 'class':
      case 'interface':
      case 'enum':
        sig = `${s.kind} ${s.name} [line ${s.line}]`
        break
      case 'type':
        sig = `type ${s.name} = ... [line ${s.line}]`
        break
      case 'func':
        sig = s.returns
          ? `function ${s.name}(${s.params ?? ''}): ${s.returns} [line ${s.line}]`
          : `function ${s.name}(${s.params ?? ''}) [line ${s.line}]`
        break
      case 'method':
        sig = s.returns
          ? `  ${s.name}(${s.params ?? ''}): ${s.returns} [line ${s.line}]`
          : `  ${s.name}(${s.params ?? ''}) [line ${s.line}]`
        break
      default:
        sig = `${s.kind} ${s.name} [line ${s.line}]`
    }
    if (s.comment) lines.push(`// ${s.comment}`)
    lines.push(sig)
  }

  return lines.join('\n') + '\n'
}
