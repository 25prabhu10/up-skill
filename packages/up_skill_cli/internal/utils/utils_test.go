package utils

import (
	"testing"
)

func TestNormalizeName(t *testing.T) {
	passTestCases := map[string]string{
		"Inverse Binary Tree": "inverse_binary_tree",
		"ABC-asdjk":           "abc_asdjk",
		"98 _ asd:A":          "98___asd_a",
		"":                    "",
		" ":                   "_",
		" _-":                 "___",
	}

	for input, expected := range passTestCases {
		result := NormalizeName(input)

		if result != expected {
			t.Errorf("\ntesting:%s expected:%s got:%s", input, expected, result)
		}

	}
}
