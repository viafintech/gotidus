package postgres

import (
	"fmt"

	"github.com/viafintech/gotidus"
)

// RegexReplaceAnonymizer is a gotidus.Anonymizer interface implementation which allows
// defining a pattern which is used to pick a part of a string which is
// then replaced with the defined replacement.
type RegexReplaceAnonymizer struct {
	pattern     string
	replacement string
}

// NewRegexReplaceAnonymizer intializes a new RegexReplaceAnonymizer object.
func NewRegexReplaceAnonymizer(pattern, replacement string) *RegexReplaceAnonymizer {
	return &RegexReplaceAnonymizer{
		pattern:     pattern,
		replacement: replacement,
	}
}

// Build returns the partial query containing the PostgreSQL REGEXP_REPLACE function.
// It uses the configured pattern and replacement and applies it to the given column.
func (a *RegexReplaceAnonymizer) Build(tableName, columnName string) string {
	return fmt.Sprintf(
		`REGEXP_REPLACE(%s, '%s', '%s')`,
		gotidus.FullColumnName(tableName, columnName),
		a.pattern,
		a.replacement,
	)
}
