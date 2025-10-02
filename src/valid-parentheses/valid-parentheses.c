#include "valid-parentheses.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

bool isValid(const char *s) {
  const size_t arr_size = strlen(s);
  char *stack = (char *)malloc(sizeof(char) * arr_size);
  int j = -1;

  for (size_t i = 0; i < arr_size; i++) {
    if (s[i] == '(' || s[i] == '{' || s[i] == '[') {
      stack[++j] = s[i];
    } else if (s[i] == ')' || s[i] == '}' || s[i] == ']') {
      if (j == -1 || ((s[i] == ')' && stack[j] != '(') ||
                      (s[i] == '}' && stack[j] != '{') ||
                      (s[i] == ']' && stack[j] != '['))) {
        return false;
      }
      j--;
    }
  }

  free(stack);

  return j == -1;
}

bool isValidFastest(char *s) {
  char stack[100000];
  int top = -1;
  for (int i = 0; s[i] != '\0'; i++) {
    char c = s[i];
    if (c == '(' || c == '{' || c == '[') {

      stack[++top] = c;
    } else {
      if (top == -1)
        return false;

      char f = stack[top--];
      if ((c == ')' && f != '(') || (c == '}' && f != '{') ||
          (c == ']' && f != '[')) {
        return false;
      }
    }
  }
  return top == -1;
}
