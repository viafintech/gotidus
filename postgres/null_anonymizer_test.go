package postgres

import "testing"

func TestNullAnonymizerBuild(t *testing.T) {
	anonymizer := NewNullAnonymizer()

	tableName := "foo"
	columnName := "bar"

	compareStrings(
		anonymizer.Build(tableName, columnName),
		"NULL::unknown",
		t,
	)
}
