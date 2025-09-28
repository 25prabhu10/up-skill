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
    if (diffs.has(nums[i])) {
      return [diffs.get(nums[i]), i]
    }

    diffs.set(target - nums[i], i)
  }

  return []
}

export function twoSumFastest(nums, target) {
  const map = new Map()
  let ans
  nums.some((number, index) => {
    const requiredValue = target - number
    if (map.has(requiredValue)) {
      ans = [map.get(requiredValue), index]
      return true
    }

    map.set(number, index)
    return false
  })

  return ans
}
