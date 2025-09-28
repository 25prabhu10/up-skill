export function isValid(s: string): boolean {
  const validPairs: Record<string, string> = {
    '(': ')',
    '[': ']',
    '{': '}',
  }
  let isBracketValid = false
  const expectedBrackets: string[] = []
  for (let i = 0; i < s.length; i++) {
    const expectedBracket = validPairs[s[i]]
    console.log('Open', s[i], 'Close', expectedBracket)

    if (expectedBracket === undefined) {
      const lastBracket = expectedBrackets.pop()
      console.log('lastBracket: -->', lastBracket)
      isBracketValid = expectedBracket === lastBracket
    } else {
      expectedBrackets.push(expectedBracket)
    }
    // if (validPairs[s[i]] === s[i + 1]) {
    //   isBracketValid = true
    // }
  }

  console.log(expectedBrackets)
  return isBracketValid
}
