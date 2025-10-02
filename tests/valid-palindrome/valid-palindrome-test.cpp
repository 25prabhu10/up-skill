#include "gtest/gtest.h"
#include <string>

extern "C" {
#include "../../src/valid-palindrome/valid-palindrome.h"
}

TEST(test_isPalindrome, case_1) {
  std::string str = "A man, a plan,_ a canal: Panama";
  bool result = isPalindrome(str.c_str());
  ASSERT_EQ(result, true);
}

TEST(test_isPalindrome, case_2) {
  std::string str = "race a car";
  bool result = isPalindrome(str.c_str());
  ASSERT_EQ(result, false);
}

TEST(test_isPalindrome, case_3) {
  std::string str = " ";
  bool result = isPalindrome(str.c_str());
  ASSERT_EQ(result, true);
}
