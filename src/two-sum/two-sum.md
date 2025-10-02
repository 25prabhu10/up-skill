---
topics: HashMap
---

# 1. Two Sum (Easy)

Problem link: [Two Sum](https://leetcode.com/problems/two-sum)

Given an array of integers `nums` and an `integer` target, return indices of the two numbers such that they add up to `target`.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

You can return the answer in any order.

- Example 1:

  Input: `nums = [2,7,11,15], target = 9`
  Output: `[0,1]`
  Explanation: `Because nums[0] + nums[1] == 9, we return [0, 1].`

- Example 2:

  Input: `nums = [3,2,4], target = 6`
  Output: `[1,2]`

- Example 3:

  Input: `nums = [3,3], target = 6`
  Output: `[0,1]`

Constraints:

- `2 <= nums.length <= 104`
- `-109 <= nums[i] <= 109`
- `-109 <= target <= 109`
- Only one valid answer exists.

## Hint

We can use brute force to solve, but that would be `O(n^2)` time complexity

We can use a **hash map to store the difference between the target and the current number**, and check if the current number exists in the hash map. This would be `O(n)` time complexity and `O(n)` space complexity
