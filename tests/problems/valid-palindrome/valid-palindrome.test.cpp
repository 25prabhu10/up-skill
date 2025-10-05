#include "gtest/gtest.h"
#include <string>

extern "C" {
#include "../../../src/problems/valid-palindrome/valid-palindrome.h"
}

struct ValidPalindromeCase {
  std::string input;
  bool expected;
};

class ValidPalindromeTest
    : public ::testing::TestWithParam<ValidPalindromeCase> {};

TEST_P(ValidPalindromeTest, ChecksIfInputIsPalindrome) {
  ValidPalindromeCase param = GetParam();
  bool result = isPalindrome(param.input.c_str());
  ASSERT_EQ(result, param.expected);
}

INSTANTIATE_TEST_SUITE_P(
    ValidPalindromeCases, ValidPalindromeTest,
    ::testing::Values(ValidPalindromeCase{"A man, a plan, a canal: Panama",
                                          true},
                      ValidPalindromeCase{"race a car", false},
                      ValidPalindromeCase{" ", true}));
