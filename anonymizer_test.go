package gotidus

import (
	"testing"

	"github.com/viafintech/gotidus/testutils"
)

func TestFullColumnName(t *testing.T) {
	testutils.CompareStrings(
		FullColumnName("table", "column"),
		"table.column",
		t,
	)
}

func TestNoopAnonymizerBuild(t *testing.T) {
	anonymizer := NewNoopAnonymizer()

	tableName := "some_table"
	columnName := "some_column"

	testutils.CompareStrings(
		anonymizer.Build(tableName, columnName),
		FullColumnName(tableName, columnName),
		t,
	)
}

func TestStaticAnonymizerBuild(t *testing.T) {
	anonymizer := NewStaticAnonymizer("23", "integer")

	tableName := "some_table"
	columnName := "some_column"

	testutils.CompareStrings(
		anonymizer.Build(tableName, columnName),
		"'23'::integer",
		t,
	)
}
