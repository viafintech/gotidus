package gotidus

import (
	"testing"

	"github.com/Barzahlen/gotidus/testutils"
)

func TestNewTable(t *testing.T) {
	expectedTable := &Table{
		columns: make(map[string]Anonymizer),
	}

	testutils.CompareStructs(NewTable(), expectedTable, t)
}

func TestTableAddAndGetAnonymizer(t *testing.T) {
	table := NewTable()

	anonymizer := NewStaticAnonymizer("bar", "TEXT")

	// Set anonymizer
	table.AddAnonymizer("foo", anonymizer)

	testutils.CompareStructs(anonymizer, table.columns["foo"], t)

	// Check loading
	anon := table.GetAnonymizer("foo")

	testutils.CompareStructs(anonymizer, anon, t)

	// Check anonymizer not set
	testutils.CompareStructs(table.columns["baz"], nil, t)

	// Check default anonymizer
	defaultAnon := table.GetAnonymizer("baz")

	testutils.CompareStructs(defaultAnon, NewNoopAnonymizer(), t)
}
