import { describe, expect, test } from 'vitest'

import { isPalindrome } from '@/problems/valid-palindrome/valid-palindrome'

describe('Valid Palindrome', () => {
  test('should return true if it is a palindrome, or false otherwise', () => {
    expect(isPalindrome('A man, a plan,_ a canal: Panama')).toBeTruthy()
    expect(isPalindrome('race a car')).toBeFalsy()
    expect(isPalindrome(' ')).toBeTruthy()
  })
})
