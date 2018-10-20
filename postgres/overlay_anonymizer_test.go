package postgres

import "testing"

func TestOverlayAnonymizerBuild(t *testing.T) {
	tableName := "foo"
	columnName := "bar"

	cases := []struct {
		title       string
		overlayBase string
		start       int
		count       int

		expectedString string
	}{
		{
			title:       "overlay with XXXXX and start at 11",
			overlayBase: "X",
			start:       11,
			count:       5,

			expectedString: `"overlay"((foo.bar)::text, 'XXXXX'::text, 11)`,
		},
		{
			title:       "overlay with YWYWYW and start at 6",
			overlayBase: "YW",
			start:       6,
			count:       3,

			expectedString: `"overlay"((foo.bar)::text, 'YWYWYW'::text, 6)`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			anonymizer := NewOverlayAnonymizer(c.overlayBase, c.start, c.count)

			compareStrings(
				anonymizer.Build(tableName, columnName),
				c.expectedString,
				t,
			)

		})
	}
}
