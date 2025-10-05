import { describe, expect, test } from 'vitest'

import { lengthOfLongestSubstring } from '@/problems/longest-substring-without-repeating-characters/longest-substring-without-repeating-characters'

describe('Longest Substring Without Repeating Characters', () => {
  test('should find the lenght of the longest substring without duplicate charaters', () => {
    expect(lengthOfLongestSubstring('abcabcbb')).toEqual(3)
    expect(lengthOfLongestSubstring('bbbbb')).toEqual(1)
    expect(lengthOfLongestSubstring('pwwkew')).toEqual(3)
    expect(lengthOfLongestSubstring(' ')).toEqual(1)
    expect(lengthOfLongestSubstring('aabaab!bb')).toEqual(3)
  })
})
