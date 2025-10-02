export function isValid(s: string): boolean {
  let brackets: string[] = []

  for (let i = 0; i < s.length; i++) {
    if (s[i] === '(' || s[i] === '{' || s[i] === '[') {
      brackets.push(s[i])
    } else if (s[i] === ')' || s[i] === '}' || s[i] === ']') {
      if (brackets.length === 0) {
        return false
      }

      let top = brackets.pop()

      if (
        (s[i] === ')' && top !== '(') ||
        (s[i] === '}' && top !== '{') ||
        (s[i] === ']' && top !== '[')
      ) {
        return false
      }
    }
  }

  return brackets.length === 0
}

export function isValidFastest(s: string): boolean {
  let parentheses: Record<string, string> = {
    ')': '(',
    ']': '[',
    '}': '{',
  }

  let open = Object.values(parentheses)

  let queue: string[] = []

  // oxlint-disable-next-line prefer-spread
  let sentence: string[] = s.split('')

  for (let i = 0; i < sentence.length; i++) {
    if (open.includes(sentence[i])) {
      queue.push(sentence[i])
    } else {
      let last = queue.pop()
      if (last !== parentheses[sentence[i]]) {
        return false
      }
    }
  }

  if (queue.length > 0) {
    return false
  }

  return true
}

export function isValidFastestJS(s: string): boolean {
  let map = new Map([
    ['{', '}'],
    ['[', ']'],
    ['(', ')'],
  ])
  let stack = []
  for (let i = 0; i < s.length; i++) {
    let char = s[i]
    if (map.has(char)) {
      stack.push(map.get(char))
    } else if (stack.pop() !== char) {
      return false
    }
  }
  return stack.length === 0
}
