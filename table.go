package gotidus

// NewTable initializes a Table object wich blank columns.
func NewTable() *Table {
	table := &Table{
		columns: make(map[string]Anonymizer),
	}

	return table
}

// Table is the type holding the column configuration
type Table struct {
	columns map[string]Anonymizer
}

// AddAnonymizer allows setting a specific Anonymizer for a column of the given name.
// If an Anonymizer was previously configured for a column name, it will be overwritten.
func (t *Table) AddAnonymizer(columnName string, anonymizer Anonymizer) *Table {
	t.columns[columnName] = anonymizer

	return t
}

// GetAnonymizer retrieves an Anonymizer from the Table configuration.
// If an Anonymizer was configured for the given name, that Anonymizer will be returned.
// If no Anonymizer was configured for the given name, the NoopAnonymizer will be returned.
func (t *Table) GetAnonymizer(columnName string) Anonymizer {
	anonymizer, ok := t.columns[columnName]
	if ok {
		return anonymizer
	}

	return NewNoopAnonymizer()
}
