package gotidus

import (
	"testing"
)

func TestNewTable(t *testing.T) {
	expectedTable := &Table{
		columns: make(map[string]Anonymizer),
	}

	compareStructs(NewTable(), expectedTable, t)
}

func TestTableAddAndGetAnonymizer(t *testing.T) {
	table := NewTable()

	anonymizer := NewStaticAnonymizer("bar", "TEXT")

	// Set anonymizer
	table.AddAnonymizer("foo", anonymizer)

	compareStructs(anonymizer, table.columns["foo"], t)

	// Check loading
	anon := table.GetAnonymizer("foo")

	compareStructs(anonymizer, anon, t)

	// Check anonymizer not set
	compareStructs(table.columns["baz"], nil, t)

	// Check default anonymizer
	defaultAnon := table.GetAnonymizer("baz")

	compareStructs(defaultAnon, NewNoopAnonymizer(), t)
}
