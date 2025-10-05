#include "gtest/gtest.h"

extern "C" {
#include "../../../src/problems/two-sum/two-sum.h"
}

TEST(test_twoSum, case_1) {
  int nums[] = {2, 7, 11, 15};
  int target = 9;
  int returnSize = 0;
  int *result = twoSum(nums, 4, target, &returnSize);
  ASSERT_EQ(returnSize, 2);
  EXPECT_EQ(result[0], 0);
  EXPECT_EQ(result[1], 1);
  free(result);
}

TEST(test_twoSum, case_2) {
  int nums[] = {3, 2, 4};
  int target = 6;
  int returnSize = 0;
  int *result = twoSum(nums, 4, target, &returnSize);
  ASSERT_EQ(returnSize, 2);
  EXPECT_EQ(result[0], 1);
  EXPECT_EQ(result[1], 2);
  free(result);
}

TEST(test_twoSum, case_3) {
  int nums[] = {3, 3};
  int target = 6;
  int returnSize = 0;
  int *result = twoSum(nums, 4, target, &returnSize);
  ASSERT_EQ(returnSize, 2);
  EXPECT_EQ(result[0], 0);
  EXPECT_EQ(result[1], 1);
  free(result);
}
