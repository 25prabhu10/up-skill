from __future__ import annotations

import importlib.util
from pathlib import Path
from typing import TYPE_CHECKING, cast

import pytest

if TYPE_CHECKING:
    from collections.abc import Callable


def _load_selection_sort() -> Callable[[list[int]], list[int]]:
    module_path = (
        Path(__file__).resolve().parents[3]
        / "src"
        / "algorithms"
        / "sort"
        / "selection-sort"
        / "selection-sort.py"
    )
    spec = importlib.util.spec_from_file_location("selection_sort_module", module_path)
    if spec is None or spec.loader is None:
        msg = f"Unable to load selection_sort module from {module_path}"
        raise RuntimeError(msg)

    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    if not hasattr(module, "selection_sort"):
        msg = "selection_sort module does not define a selection_sort function"
        raise AttributeError(msg)

    return cast("Callable[[list[int]], list[int]]", module.selection_sort)


selection_sort = _load_selection_sort()


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
