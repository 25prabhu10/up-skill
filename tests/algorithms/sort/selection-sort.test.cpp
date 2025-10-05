#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include <vector>

extern "C" {
#include "../../../src/algorithms/sort/selection-sort/selection-sort.h"
}

struct SelectionSortCase {
  std::vector<int> input;
  std::vector<int> expected;
};

class SelectionSortTest : public ::testing::TestWithParam<SelectionSortCase> {};

TEST_P(SelectionSortTest, SortsInPlaceAccordingToExpectation) {
  SelectionSortCase param = GetParam();
  std::vector<int> data = param.input;

  selection_sort(data.data(), static_cast<int>(data.size()));

  EXPECT_THAT(data, testing::ElementsAreArray(param.expected));
}

INSTANTIATE_TEST_SUITE_P(
    SelectionSortCases, SelectionSortTest,
    ::testing::Values(SelectionSortCase{{5, 2, 9, 1, 5}, {1, 2, 5, 5, 9}},
                      SelectionSortCase{{4, 1, 3, 9, 7}, {1, 3, 4, 7, 9}},
                      SelectionSortCase{{3, -1, 3, 2, -1}, {-1, -1, 2, 3, 3}},
                      SelectionSortCase{{1, 2, 3, 4, 5}, {1, 2, 3, 4, 5}},
                      SelectionSortCase{{42}, {42}},
                      SelectionSortCase{{}, {}}));
