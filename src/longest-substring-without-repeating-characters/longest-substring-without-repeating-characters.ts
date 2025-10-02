export function lengthOfLongestSubstring(s: string): number {
  let substring: string[] = []
  let longestSubstring = 0

  for (let i = 0; i < s.length; i++) {
    if (substring.includes(s[i])) {
      substring = substring.slice(substring.indexOf(s[i]) + 1)
    }
    substring.push(s[i])
    longestSubstring = longestSubstring > substring.length ? longestSubstring : substring.length
  }

  return longestSubstring
}

export function lengthOfLongestSubstringFastest(s: string): number {
  let left = 0
  let maxLength = 0
  let charSet = new Set()

  for (let right = 0; right < s.length; right++) {
    while (charSet.has(s[right])) {
      charSet.delete(s[left])
      left++
    }

    charSet.add(s[right])
    maxLength = Math.max(maxLength, right - left + 1)
  }

  return maxLength
}
