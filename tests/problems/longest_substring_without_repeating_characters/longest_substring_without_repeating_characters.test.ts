import { lengthOfLongestSubstring } from "@/problems/longest_substring_without_repeating_characters/longest_substring_without_repeating_characters"
import { describe, expect, it } from "vitest"

describe("longest Substring Without Repeating Characters", () => {
  it("should find the lenght of the longest substring without duplicate charaters", () => {
    expect(lengthOfLongestSubstring("abcabcbb")).toBe(3)
    expect(lengthOfLongestSubstring("bbbbb")).toBe(1)
    expect(lengthOfLongestSubstring("pwwkew")).toBe(3)
    expect(lengthOfLongestSubstring(" ")).toBe(1)
    expect(lengthOfLongestSubstring("aabaab!bb")).toBe(3)
  })
})
