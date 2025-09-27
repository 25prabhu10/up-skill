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
