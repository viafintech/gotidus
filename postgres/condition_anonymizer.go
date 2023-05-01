package postgres

import (
	"fmt"
	"strings"

	"github.com/viafintech/gotidus"
)

// ConditionAnonymizer is a gotidus.Anonymizer interface implementation
// which allows elaborate case statements for anonymizing columns
// based on values in the same or other columns.
type ConditionAnonymizer struct {
	dataType          string
	defaultAnonymizer gotidus.Anonymizer
	conditions        []AnonymizationCondition
}

// NewConditionAnonymizer initializes a new ConditionAnonymizer object
func NewConditionAnonymizer(
	dataType string,
	defaultAnonymizer gotidus.Anonymizer,
	conditions ...AnonymizationCondition,
) *ConditionAnonymizer {
	return &ConditionAnonymizer{
		dataType:          dataType,
		defaultAnonymizer: defaultAnonymizer,
		conditions:        conditions,
	}
}

// Build returns the partial query from the configured data.
// If no conditions were given, it behaves as if the defaultAnonymizer was called directly.
// If one or more conditions are configured, they are defined as cases
// in a CASE statement and the defaultAnonymizer covers the ELSE branch.
func (a *ConditionAnonymizer) Build(tableName, columnName string) string {
	defaultCase := a.defaultAnonymizer.Build(tableName, columnName)

	if len(a.conditions) < 1 {
		return defaultCase
	}

	whenStrings := []string{}

	for _, condition := range a.conditions {
		whenStrings = append(whenStrings, condition.BuildCase(tableName, columnName))
	}

	return fmt.Sprintf(
		"(CASE %s ELSE %s END)::%s",
		strings.Join(whenStrings, " "),
		defaultCase,
		a.dataType,
	)
}

// AnonymizationCondition is the implementation which holds the conditions for the ConditionAnonymizer.
type AnonymizationCondition struct {
	column     string
	comparator string
	value      string
	dataType   string
	anonymizer gotidus.Anonymizer
}

// NewAnonymizationCondition initializes a new AnonymizationCondition object
func NewAnonymizationCondition(
	column string,
	comparator string,
	value string,
	dataType string,
	anonymizer gotidus.Anonymizer,
) AnonymizationCondition {
	return AnonymizationCondition{
		column:     column,
		comparator: comparator,
		value:      value,
		dataType:   dataType,
		anonymizer: anonymizer,
	}
}

// BuildCase builds the WHEN/THEN line of the CASE statement for the ConditionAnonymizer
func (ac AnonymizationCondition) BuildCase(
	tableName string,
	columnName string,
) string {
	return fmt.Sprintf(
		"WHEN (%s) %s '%s'::%s THEN (%s)",
		applyType(gotidus.FullColumnName(tableName, ac.column), ac.dataType),
		ac.comparator,
		ac.value,
		ac.dataType,
		ac.anonymizer.Build(tableName, columnName),
	)
}

func applyType(partial string, dataType string) string {
	return fmt.Sprintf("(%s)::%s", partial, dataType)
}
