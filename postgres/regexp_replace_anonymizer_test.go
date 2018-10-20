package postgres

import "testing"

func TestRegexReplaceAnonymizerBuild(t *testing.T) {
	tableName := "foo"
	columnName := "bar"

	cases := []struct {
		title       string
		pattern     string
		replacement string

		expectedString string
	}{
		{
			title:       "replace column following a pattern with a string",
			pattern:     "1234",
			replacement: "a",

			expectedString: `REGEXP_REPLACE(foo.bar, '1234', 'a')`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			anonymizer := NewRegexReplaceAnonymizer(c.pattern, c.replacement)

			compareStrings(
				anonymizer.Build(tableName, columnName),
				c.expectedString,
				t,
			)

		})
	}
}
