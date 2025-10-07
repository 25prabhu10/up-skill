#include "gtest/gtest.h"
#include <string>

extern "C" {
#include "../../../src/problems/valid_parentheses/valid_parentheses.h"
}

struct ValidParenthesesCase {
  std::string input;
  bool expected;
};

class ValidParenthesesTest
    : public ::testing::TestWithParam<ValidParenthesesCase> {};

TEST_P(ValidParenthesesTest, ChecksIfInputIsValid) {
  ValidParenthesesCase param = GetParam();
  bool result = isValid(param.input.c_str());
  ASSERT_EQ(result, param.expected);
}

INSTANTIATE_TEST_SUITE_P(ValidParenthesesCases, ValidParenthesesTest,
                         ::testing::Values(ValidParenthesesCase{"()", true},
                                           ValidParenthesesCase{"()[]{}", true},
                                           ValidParenthesesCase{"(]", false},
                                           ValidParenthesesCase{"([])", true},
                                           ValidParenthesesCase{"([])]", false},
                                           ValidParenthesesCase{"[", false}));
