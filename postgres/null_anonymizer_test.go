package postgres

import (
	"testing"

	"github.com/Barzahlen/gotidus/testutils"
)

func TestNullAnonymizerBuild(t *testing.T) {
	anonymizer := NewNullAnonymizer()

	tableName := "foo"
	columnName := "bar"

	testutils.CompareStrings(
		anonymizer.Build(tableName, columnName),
		"NULL::unknown",
		t,
	)
}
