package postgres

import (
	"testing"

	"github.com/viafintech/gotidus/testutils"
)

func TestEmailAnonymizerBuild(t *testing.T) {
	tableName := "foo"
	columnName := "bar"

	cases := []struct {
		title   string
		options []EmailAnonymizerOption

		expectedString string
	}{
		{
			title: "email overlay",

			expectedString: `CASE WHEN ((foo.bar)::TEXT ~~ '%@%'::TEXT)
    THEN (
      (
        (
          "left"(
            md5((foo.bar)::text),
            15
          ) ||
          '@'::text
        ) ||
        split_part((foo.bar)::text, '@'::text, 2)
      )
    )::CHARACTER VARYING
    ELSE foo.bar
    END`,
		},
		{
			title: "email overlay with different length",
			options: []EmailAnonymizerOption{
				EmailAnonymizedPartLengthOption(10),
			},

			expectedString: `CASE WHEN ((foo.bar)::TEXT ~~ '%@%'::TEXT)
    THEN (
      (
        (
          "left"(
            md5((foo.bar)::text),
            10
          ) ||
          '@'::text
        ) ||
        split_part((foo.bar)::text, '@'::text, 2)
      )
    )::CHARACTER VARYING
    ELSE foo.bar
    END`,
		},
		{
			title: "email overlay with anonymized domain part length",
			options: []EmailAnonymizerOption{
				EmailAnonymizedPartLengthOption(10),
				EmailAnonymizeDomainPartOption(true),
			},

			expectedString: `CASE WHEN ((foo.bar)::TEXT ~~ '%@%'::TEXT)
    THEN (
      (
        (
          "left"(
            md5((foo.bar)::text),
            10
          ) ||
          '@'::text
        ) ||
        "left"(md5(split_part((foo.bar)::text, '@'::text, 2)::text) 10)
      )
    )::CHARACTER VARYING
    ELSE foo.bar
    END`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			anonymizer := NewEmailAnonymizer(c.options...)

			testutils.CompareStrings(
				anonymizer.Build(tableName, columnName),
				c.expectedString,
				t,
			)

		})
	}
}
