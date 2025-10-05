#include "gtest/gtest.h"
#include <cstdlib>
#include <vector>

extern "C" {
#include "../../../src/problems/two-sum/two-sum.h"
}

struct TwoSumCase {
  std::vector<int> nums;
  int target;
  std::vector<int> expectedIndices;
};

class TwoSumTest : public ::testing::TestWithParam<TwoSumCase> {};

TEST_P(TwoSumTest, FindsExpectedIndices) {
  const TwoSumCase param = GetParam();
  std::vector<int> numbers = param.nums;
  int returnSize = -1;
  int *result = twoSum(numbers.data(), static_cast<int>(numbers.size()),
                       param.target, &returnSize);

  EXPECT_EQ(returnSize, static_cast<int>(param.expectedIndices.size()));

  if (param.expectedIndices.empty()) {
    // When no solution exists, implementation may return nullptr or a unique
    // pointer.
    if (result != nullptr) {
      free(result);
    }
    SUCCEED();
    return;
  }

  ASSERT_NE(result, nullptr);

  ASSERT_EQ(returnSize, 2);
  ASSERT_EQ(param.expectedIndices.size(), 2u);

  EXPECT_EQ(result[0], param.expectedIndices[0]);
  EXPECT_EQ(result[1], param.expectedIndices[1]);

  ASSERT_LT(result[0], static_cast<int>(numbers.size()));
  ASSERT_LT(result[1], static_cast<int>(numbers.size()));
  ASSERT_GE(result[0], 0);
  ASSERT_GE(result[1], 0);
  EXPECT_NE(result[0], result[1]);
  EXPECT_EQ(numbers[result[0]] + numbers[result[1]], param.target);

  free(result);
}

INSTANTIATE_TEST_SUITE_P(
    TwoSumCases, TwoSumTest,
    ::testing::Values(TwoSumCase{{2, 7, 11, 15}, 9, {0, 1}},
                      TwoSumCase{{3, 2, 4}, 6, {1, 2}},
                      TwoSumCase{{3, 3}, 6, {0, 1}},
                      TwoSumCase{{-1, -2, -3, -4, -5}, -8, {2, 4}},
                      TwoSumCase{{5, 75, 25}, 100, {1, 2}},
                      TwoSumCase{{1, 2, 3, 4}, 100, {}}));
