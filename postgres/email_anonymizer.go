package postgres

import (
  "fmt"

  "github.com/viafintech/gotidus"
)

// EmailAnonymizer is a gotidus.Anonymizer interface implementation,
// which allows anonymizing the local part of a string matching the general pattern of a mail address.
// For example
//
//  "test@example.com"
//
// would be anonymized by a 15 character string based on the md5 hash of the local part.
type EmailAnonymizer struct {
  mailLocalPartLength int
}

// NewEmailAnonymizer initializes a new EmailAnonymizer object
func NewEmailAnonymizer() *EmailAnonymizer {
  return &EmailAnonymizer{
    mailLocalPartLength: 15,
  }
}

// Build returns the partial query holding the logic to overwrite
// the local part of an email address.
func (a *EmailAnonymizer) Build(tableName, columnName string) string {
  return fmt.Sprintf(
    `CASE WHEN ((%[1]s)::TEXT ~~ '%%@%%'::TEXT)
    THEN (
      (
        (
          "left"(
            md5((%[1]s)::text),
            %[2]d
          ) ||
          '@'::text
        ) ||
        split_part(
          (%[1]s)::text,
          '@'::text,
          2
        )
      )
    )::CHARACTER VARYING
    ELSE %[1]s
    END`,
    gotidus.FullColumnName(tableName, columnName),
    a.mailLocalPartLength,
  )
}
