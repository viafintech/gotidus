package gotidus

import (
	"fmt"
)

// FullColumnName is a helper function that allows building the name
// based on the table and column names.
func FullColumnName(tableName, columnName string) string {
	return fmt.Sprintf("%s.%s", tableName, columnName)
}

// Anonymizer is the interface for functions that build the query snippet to anonymize a specific column.
type Anonymizer interface {
	Build(tableName, columnName string) string
}

// NoopAnonymizer is an Anonymizer interface implementation which returns the column value is as.
// It is also the default anonymizer for every column unless otherwise defined.
type NoopAnonymizer struct{}

// NewNoopAnonymizer initializes a new NoopAnonymizer object
func NewNoopAnonymizer() *NoopAnonymizer {
	return &NoopAnonymizer{}
}

// Build returns the column name build from the table and column name
func (a *NoopAnonymizer) Build(tableName, columnName string) string {
	return FullColumnName(tableName, columnName)
}

// StaticAnonymizer is an Anonymizer interfface implementation that ensures that every row returns the same static value.
type StaticAnonymizer struct {
	staticValue string
	dataType    string
}

// NewStaticAnonymizer initializes a new StaticAnonymizer object
func NewStaticAnonymizer(staticValue, dataType string) *StaticAnonymizer {
	return &StaticAnonymizer{
		staticValue: staticValue,
		dataType:    dataType,
	}
}

// Build returns a partial query from the static value and data type given on object initialization.
// table and column name are ignored here.
func (a *StaticAnonymizer) Build(tableName, columnName string) string {
	return fmt.Sprintf("'%s'::%s", a.staticValue, a.dataType)
}
