package postgres

// NullAnonymizer is a gotidus.Anonymizer interface implementation.
// It overwrites any given value with NULL.
type NullAnonymizer struct{}

// NewNullAnonymizer initializes a new NullAnonymizer object.
func NewNullAnonymizer() *NullAnonymizer {
	return &NullAnonymizer{}
}

// Build returns 'NULL::unknown' as partial query.
func (a *NullAnonymizer) Build(tableName, columnName string) string {
	return "NULL::unknown"
}
