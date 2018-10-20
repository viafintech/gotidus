package gotidus

import (
	"encoding/json"
	"reflect"
	"testing"
)

func compareStrings(got, expected string, t *testing.T) {
	if got != expected {
		t.Errorf(
			"\nGot unexpected string: \n%#v\nExpected string:\n%#v\n",
			got,
			expected,
		)
	}
}

func compareStructs(got, expected interface{}, t *testing.T) {
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
