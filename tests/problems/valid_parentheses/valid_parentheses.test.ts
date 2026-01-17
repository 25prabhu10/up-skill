import { isValid } from "@/problems/valid_parentheses/valid_parentheses"
import { describe, expect, it } from "vitest"

describe("valid parentheses", () => {
  it("should determine if the input string is valid", () => {
    expect(isValid("()")).toBeTruthy()
    expect(isValid("()[]{}")).toBeTruthy()
    expect(isValid("(]")).toBeFalsy()
    expect(isValid("([])")).toBeTruthy()
    expect(isValid("([])]")).toBeFalsy()
    expect(isValid("[")).toBeFalsy()
  })
})
