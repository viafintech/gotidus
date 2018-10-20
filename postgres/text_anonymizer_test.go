package postgres

import (
	"strings"
	"testing"
)

func TestTextAnonymizerBuild(t *testing.T) {
	tableName := "foo"
	columnName := "bar"

	cases := []struct {
		title          string
		mappingBuilder func(string) string

		expectedString string
	}{
		{
			title: "uses the built pattern",
			mappingBuilder: func(string) string {
				return "9876543210ÖÄÜöäüßZYXWVUTSRQPONMLKJIHGFEDCBAzyxwvutsrqponmlkjihgfedcba"
			},
			expectedString: `translate(` +
				`foo.bar::TEXT, ` +
				`'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZßüäöÜÄÖ0123456789'::TEXT, ` +
				`'9876543210ÖÄÜöäüßZYXWVUTSRQPONMLKJIHGFEDCBAzyxwvutsrqponmlkjihgfedcba'::TEXT` +
				`)`,
		},
		{
			title: "uses the build pattern (2)",
			mappingBuilder: func(string) string {
				return "JIHGFEDCBAzyxwvutsrqponmlkjihgfedcba9876543210ÖÄÜöäüßZYXWVUTSRQPONMLK"
			},
			expectedString: `translate(` +
				`foo.bar::TEXT, ` +
				`'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZßüäöÜÄÖ0123456789'::TEXT, ` +
				`'JIHGFEDCBAzyxwvutsrqponmlkjihgfedcba9876543210ÖÄÜöäüßZYXWVUTSRQPONMLK'::TEXT` +
				`)`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {

			anonymizer := NewTextAnonymizer()

			anonymizer.buildMapping = c.mappingBuilder

			compareStrings(
				anonymizer.Build(tableName, columnName),
				c.expectedString,
				t,
			)

		})
	}
}

func TestBuildMapping(t *testing.T) {
	cases := []struct {
		title string
		base  string
	}{
		{
			title: "randomized string from base string",
			base:  "ABCDEFGHIJ",
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			str1 := buildMapping(c.base)
			str2 := buildMapping(c.base)

			if str1 == str2 {
				t.Error("Generated two strings identical")
			}

			for _, c := range strings.Split(c.base, "") {
				if !strings.Contains(str1, c) {
					t.Errorf("Missing expected character '%s' from %s", c, str1)
				}

				if !strings.Contains(str2, c) {
					t.Errorf("Missing expected character '%s' from %s", c, str2)
				}
			}
		})
	}
}
