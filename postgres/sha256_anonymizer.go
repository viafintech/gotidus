package postgres

import "fmt"

// SHA256Anonymizer is a gotidus.Anonymizer interface implementation.
// It overwrites any given value with the SHA256 value limited to the given length.
type SHA256Anonymizer struct {
  length int
}

// NewSHA256Anonymizer initializes a new SHA256Anonymizer object.
func NewSHA256Anonymizer(length int) *SHA256Anonymizer {
  return &SHA256Anonymizer{
    length: length,
  }
}

// Build returns 'NULL::unknown' as partial query.
func (a *SHA256Anonymizer) Build(tableName, columnName string) string {
  return fmt.Sprintf(
    "SUBSTRING(ENCODE(DIGEST('%s.%s', 'sha256'), 'HEX'), 0, %d)",
    tableName,
    columnName,
    // +1 as substring is excluding the last character
    // and passing 10 would only result in 9 characters
    a.length+1,
  )
}
