export function twoSumBrute(nums: number[], target: number): number[] {
  for (let i = 0; i < nums.length; i++) {
    for (let j = i + 1; j < nums.length; j++) {
      if (nums[i] + nums[j] === target) {
        return [i, j]
      }
    }
  }

  return []
}

export function twoSum(nums: number[], target: number): number[] {
  const diffs = new Map()

  for (let i = 0; i < nums.length; i++) {
    if (diffs.get(nums[i]) !== undefined) {
      return [diffs.get(nums[i]), i]
    }

    diffs.set(target - nums[i], i)
  }

  return []
}
