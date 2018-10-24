package testutils

import (
	"encoding/json"
	"reflect"
	"testing"
)

// CompareStrings is a function to reduce boilerplace for string comparison for go tests
func CompareStrings(got, expected string, t *testing.T) {
	if got != expected {
		t.Errorf(
			"\nGot unexpected string: \n%#v\nExpected string:\n%#v\n",
			got,
			expected,
		)
	}
}

// CompareStructs is a function to reduce boilerplace for struct comparison for go tests
func CompareStructs(got, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(got, expected) {
		// Also show JSON for easier visual comparison
		o1, _ := json.MarshalIndent(got, "", "  ")
		o2, _ := json.MarshalIndent(expected, "", "  ")

		t.Errorf(
			"\nGot unexpected output: \n%#v\nExpected output:\n%#v\n"+
				"Got JSON:\n%s\nExpected JSON:\n%s\n",
			got, expected, o1, o2)
	}
}
