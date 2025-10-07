#include "two_sum.h"
#include <stdlib.h>

int *twoSumBrute(int *nums, int numsSize, int target, int *returnSize) {
  for (int i = 0; i < numsSize; i++) {
    for (int j = i + 1; j < numsSize; j++) {
      if (nums[i] + nums[j] == target) {
        int *indexs = (int *)malloc(2 * sizeof(int));
        indexs[0] = i;
        indexs[1] = j;
        *returnSize = 2;
        return indexs;
      }
    }
  }

  // Return an empty array if no solution is found
  *returnSize = 0;
  return malloc(sizeof(int) * *returnSize);
}

int *twoSum(int *nums, int numsSize, int target, int *returnSize) {
  return twoSumBrute(nums, numsSize, target, returnSize);
}
