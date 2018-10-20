package postgres

import "testing"

func compareStrings(got, expected string, t *testing.T) {
	if got != expected {
		t.Errorf(
			"\nGot unexpected string: \n%#v\nExpected string:\n%#v\n",
			got,
			expected,
		)
	}
}
