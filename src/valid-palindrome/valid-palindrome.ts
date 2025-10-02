export function isPalindrome(s: string): boolean {
  const str = s.replace(/[\W_]/gi, '').toLowerCase()

  for (let i = 0; i < str.length / 2; i++) {
    if (str[i] !== str[str.length - i - 1]) {
      return false
    }
  }

  return true
}

export function isPalindromeFastest(s: string): boolean {
  // Note that \W is the equivalent of [^0-9a-zA-Z_]
  let filtered = s.replace(/[^0-9a-z]/gi, '').toLowerCase()

  if (filtered.length === 0) {
    return true
  }

  let l = 0,
    r = filtered.length - 1

  while (l < r) {
    if (filtered[l] !== filtered[r]) {
      return false
    }

    l++
    r--
  }

  return true
}

export function isPalindromeFastestJS(s: string): boolean {
  const newStr = s.toLowerCase().replace(/[^a-z0-9]/g, '')
  // oxlint-disable-next-line prefer-spread
  return newStr === newStr.split('').reverse().join('')
}
