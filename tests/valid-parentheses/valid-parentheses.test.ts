import { describe, expect, test } from 'vitest'

import { isValid } from '@/valid-parentheses/valid-parentheses'

describe('Valid Parentheses', () => {
  test('should determine if the input string is valid', () => {
    expect(isValid('()')).toBeTruthy()
    expect(isValid('()[]{}')).toBeTruthy()
    expect(isValid('(]')).toBeFalsy()
    expect(isValid('([])')).toBeTruthy()
    expect(isValid('([])]')).toBeFalsy()
    expect(isValid('[')).toBeFalsy()
  })
})
