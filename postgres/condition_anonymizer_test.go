package postgres

import (
	"testing"

	"github.com/Barzahlen/gotidus"
)

func TestConditionAnonymizerBuild(t *testing.T) {
	tableName := "foo"
	columnName := "bar"

	cases := []struct {
		title             string
		dataType          string
		defaultAnonymizer gotidus.Anonymizer
		conditions        []AnonymizationCondition

		expectedString string
	}{
		{
			title:             "no condition",
			dataType:          "integer",
			defaultAnonymizer: gotidus.NewStaticAnonymizer("static", "TEXT"),

			expectedString: `'static'::TEXT`,
		},
		{
			title:             "one condition",
			dataType:          "TEXT",
			defaultAnonymizer: gotidus.NewStaticAnonymizer("static", "TEXT"),
			conditions: []AnonymizationCondition{
				NewAnonymizationCondition("id", "=", "5", "integer", NewNullAnonymizer()),
			},

			expectedString: `(CASE ` +
				`WHEN ((foo.id)::integer) = '5'::integer THEN (NULL::unknown) ` +
				`ELSE 'static'::TEXT END)::TEXT`,
		},
		{
			title:             "two conditions",
			dataType:          "TEXT",
			defaultAnonymizer: gotidus.NewStaticAnonymizer("static", "TEXT"),
			conditions: []AnonymizationCondition{
				NewAnonymizationCondition("id", "=", "5", "integer", NewNullAnonymizer()),
				NewAnonymizationCondition(
					"id",
					">",
					"10",
					"integer",
					gotidus.NewStaticAnonymizer("static", "TEXT"),
				),
			},

			expectedString: `(CASE ` +
				`WHEN ((foo.id)::integer) = '5'::integer THEN (NULL::unknown) ` +
				`WHEN ((foo.id)::integer) > '10'::integer THEN ('static'::TEXT) ` +
				`ELSE 'static'::TEXT END)::TEXT`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			anonymizer := NewConditionAnonymizer(
				c.dataType,
				c.defaultAnonymizer,
				c.conditions...,
			)

			compareStrings(
				anonymizer.Build(tableName, columnName),
				c.expectedString,
				t,
			)

		})
	}
}
