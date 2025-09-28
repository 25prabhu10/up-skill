import { describe, expect, test } from 'vitest'

import { lengthOfLongestSubstring } from '../../src/longest-substring-without-repeating-characters/longest-substring-without-repeating-characters'

describe('Longest Substring Without Repeating Characters', () => {
  test('should find the lenght of the longest substring without duplicate charaters', () => {
    expect(lengthOfLongestSubstring('abcabcbb')).equals(3)
    expect(lengthOfLongestSubstring('bbbbb')).equals(1)
    expect(lengthOfLongestSubstring('pwwkew')).equals(3)
    expect(lengthOfLongestSubstring(' ')).equals(1)
    expect(lengthOfLongestSubstring('aabaab!bb')).equals(3)
  })
})
