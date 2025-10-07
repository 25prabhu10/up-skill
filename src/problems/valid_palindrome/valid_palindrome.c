#include "valid_palindrome.h"
#include <ctype.h>
#include <stdbool.h>
#include <string.h>

bool isPalindrome(const char *s) {
  int left = 0;
  int right = (int)strlen(s) - 1;

  while (left < right) {
    if (!isalnum(s[left])) {
      left++;
      continue;
    } else if (!isalnum(s[right])) {
      right--;
      continue;
    } else if (tolower(s[left]) == tolower(s[right])) {
      left++;
      right--;
      continue;
    } else {
      return false;
    }
  }

  return true;
}

bool isPalindromeFastest(char *s) {
  int i = 0;
  int j = strlen(s) - 1;
  while (i < j) {
    if (!isalnum(s[i])) {
      i++;
      continue;
    }

    if (!isalnum(s[j])) {
      j--;
      continue;
    }

    if (tolower(s[i]) != tolower(s[j]))
      return false;
    i++;
    j--;
  }
  return true;
}
