export function twoSumBrute(nums: number[], target: number): number[] {
  for (let i = 0; i < nums.length; i += 1) {
    for (let j = i + 1; j < nums.length; j += 1) {
      if (nums[i] + nums[j] === target) {
        return [i, j]
      }
    }
  }

  return []
}

export function twoSum(nums: number[], target: number): number[] {
  const diffs = new Map<number, number>()

  for (let i = 0; i < nums.length; i += 1) {
    if (diffs.has(nums[i])) {
      return [diffs.get(nums[i]) ?? 0, i]
    }

    diffs.set(target - nums[i], i)
  }

  return []
}

export function twoSumFastest(nums: number[], target: number): number[] {
  const map = new Map<number, number>()
  let ans: number[] = []
  nums.some((number, index) => {
    const requiredValue = target - number
    if (map.has(requiredValue)) {
      // @ts-expect-error -- IGNORE --
      ans = [map.get(requiredValue), index]
      return true
    }

    map.set(number, index)
    return false
  })

  return ans
}
