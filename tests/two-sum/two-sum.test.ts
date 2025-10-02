import { describe, expect, test } from 'vitest'

import { twoSum } from '@/two-sum/two-sum'

describe('Two Sum', () => {
  test('should return indices of the two numbers such that they add up', () => {
    expect(twoSum([2, 7, 11, 15], 9)).toStrictEqual([0, 1])
    expect(twoSum([3, 2, 4], 6)).toStrictEqual([1, 2])
    expect(twoSum([3, 3], 6)).toStrictEqual([0, 1])
  })
})
