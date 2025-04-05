package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// First, verify that the slice lengths match.
		if len(actual) != len(c.expected) {
			t.Errorf("For input %q: expected slice length %d, got %d. Expected slice : %v Actual slice: %v",
				c.input, len(c.expected), len(actual), c.expected, actual)
			// Skip further checks for this test case.
			continue
		}

		// Now, compare each element in the slice.
		for i, word := range actual {
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("For input %q at index %d: expected %q, got %q", c.input, i, expectedWord, word)
			}
		}
	}
}
