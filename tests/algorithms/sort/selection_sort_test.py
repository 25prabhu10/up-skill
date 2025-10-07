from __future__ import annotations

import pytest

from src.algorithms.sort.selection_sort import selection_sort


@pytest.mark.parametrize(
    ("original", "expected"),
    [
        ([], []),
        ([1], [1]),
        ([5, 3, 6, 2, 10], [2, 3, 5, 6, 10]),
        ([3, 1, 2, 3, 1], [1, 1, 2, 3, 3]),
        ([-2, -5, 0, 3], [-5, -2, 0, 3]),
    ],
)
def test_selection_sort_sorts_values(original: list[int], expected: list[int]) -> None:
    working = original.copy()

    result = selection_sort(working)

    assert result == expected
    assert working == []
    assert result is not working
