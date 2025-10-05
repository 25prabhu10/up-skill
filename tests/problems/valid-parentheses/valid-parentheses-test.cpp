#include "gtest/gtest.h"
#include <string>

extern "C" {
#include "../../../src/problems/valid-parentheses/valid-parentheses.h"
}

TEST(test_isValid, case_1) {
  std::string str = "()";
  bool result = isValid(str.c_str());
  ASSERT_EQ(result, true);
}

TEST(test_isValid, case_2) {
  std::string str = "()[]{}";
  bool result = isValid(str.c_str());
  ASSERT_EQ(result, true);
}

TEST(test_isValid, case_3) {
  std::string str = "(]";
  bool result = isValid(str.c_str());
  ASSERT_EQ(result, false);
}

TEST(test_isValid, case_4) {
  std::string str = "([])";
  bool result = isValid(str.c_str());
  ASSERT_EQ(result, true);
}

TEST(test_isValid, case_5) {
  std::string str = "([])]";
  bool result = isValid(str.c_str());
  ASSERT_EQ(result, false);
}

TEST(test_isValid, case_6) {
  std::string str = "[";
  bool result = isValid(str.c_str());
  ASSERT_EQ(result, false);
}
