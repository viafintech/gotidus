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
  mailAnonymizedPartLength int
  mailAnonymizeDomainPart  bool
}

// EmailAnonymizerOption is the function type for passing options
// to the EmailAnonymizer during initialization.
type EmailAnonymizerOption func(*EmailAnonymizer)

// EmailAnonymizedPartLengthOption allows setting the target length for the local part
// of the anonymization rule
func EmailAnonymizedPartLengthOption(length int) EmailAnonymizerOption {
  return func(anonymizer *EmailAnonymizer) {
    anonymizer.mailAnonymizedPartLength = length
  }
}

// EmailAnonymizeDomainPartOption allows configuring whether the domain part should
// also explicitly be anonymized
func EmailAnonymizeDomainPartOption(anonymizeDomainPart bool) EmailAnonymizerOption {
  return func(anonymizer *EmailAnonymizer) {
    anonymizer.mailAnonymizeDomainPart = anonymizeDomainPart
  }
}

// NewEmailAnonymizer initializes a new EmailAnonymizer object
func NewEmailAnonymizer(options ...EmailAnonymizerOption) *EmailAnonymizer {
  anonymizer := &EmailAnonymizer{
    mailAnonymizedPartLength: 15,
  }

  for _, option := range options {
    option(anonymizer)
  }

  return anonymizer
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
        %[3]s
      )
    )::CHARACTER VARYING
    ELSE %[1]s
    END`,
    gotidus.FullColumnName(tableName, columnName),
    a.mailAnonymizedPartLength,
    a.domainPart(tableName, columnName),
  )
}

func (a *EmailAnonymizer) domainPart(tableName, columnName string) string {
  if a.mailAnonymizeDomainPart {
    return fmt.Sprintf(
      `("left"(md5(split_part((%[1]s)::text, '@'::text, 2)::text), %[2]d) || '.com')`,
      gotidus.FullColumnName(tableName, columnName),
      a.mailAnonymizedPartLength,
    )
  }

  return fmt.Sprintf(
    `split_part((%s)::text, '@'::text, 2)`,
    gotidus.FullColumnName(tableName, columnName),
  )
}
