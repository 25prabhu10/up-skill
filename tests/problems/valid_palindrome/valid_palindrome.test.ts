import { isPalindrome } from "@/problems/valid_palindrome/valid_palindrome"
import { describe, expect, it } from "vitest"

describe("valid Palindrome", () => {
  it("should return true if it is a palindrome, or false otherwise", () => {
    expect(isPalindrome("A man, a plan,_ a canal: Panama")).toBeTruthy()
    expect(isPalindrome("race a car")).toBeFalsy()
    expect(isPalindrome(" ")).toBeTruthy()
  })
})
