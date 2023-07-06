package postgres

import (
  "testing"

  "github.com/viafintech/gotidus/testutils"
)

func TestSHA256AnonymizerBuild(t *testing.T) {
  anonymizer := NewSHA256Anonymizer(7)

  tableName := "foo"
  columnName := "bar"

  testutils.CompareStrings(
    anonymizer.Build(tableName, columnName),
    "SUBSTRING(ENCODE(DIGEST('foo.bar', 'sha256'), 'HEX'), 0, 8)",
    t,
  )
}
