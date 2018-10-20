package postgres

import "testing"

func TestEmailAnonymizerBuild(t *testing.T) {
	tableName := "foo"
	columnName := "bar"

	cases := []struct {
		title string

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
        split_part(
          (foo.bar)::text,
          '@'::text,
          2
        )
      )
    )::CHARACTER VARYING
    ELSE foo.bar
    END`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			anonymizer := NewEmailAnonymizer()

			compareStrings(
				anonymizer.Build(tableName, columnName),
				c.expectedString,
				t,
			)

		})
	}
}
