// +build docker

package e2etests

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/Barzahlen/gotidus"
	"github.com/Barzahlen/gotidus/postgres"
	"github.com/Barzahlen/gotidus/testutils"
)

type queryCheck struct {
	Query           string
	ExpectationFunc func(rows *sql.Row, t *testing.T)
}

func TestPostgres(t *testing.T) {
	db, err := sql.Open("postgres", PGURI)
	if err != nil {
		t.Fatalf("Failed to connect to database: %+v", err)
	}

	cases := []struct {
		title               string
		setupQueries        []string
		anonymizationConfig map[string]*gotidus.Table
		queryChecks         []queryCheck
	}{
		{
			title: "NoopAnonymizer: return value as is",
			setupQueries: []string{
				"CREATE TABLE test_table (test_column TEXT)",
				"INSERT INTO test_table (test_column) VALUES ('some_unchanged_value')",
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().AddAnonymizer("test_column", gotidus.NewNoopAnonymizer()),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "some_unchanged_value"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "StaticAnonymizer: overwrite with static value",
			setupQueries: []string{
				"CREATE TABLE test_table (test_column TEXT)",
				"INSERT INTO test_table (test_column) VALUES ('value_to_be_overwritten')",
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer("test_column", gotidus.NewStaticAnonymizer("static_value", "TEXT")),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "static_value"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "ConditionAnonymizer",

			setupQueries: []string{
				"CREATE TABLE test_table (id INT, reference INT, test_column TEXT)",
				`INSERT INTO test_table (id, reference, test_column) VALUES ` +
					`(1,  5, 'value1'), ` +
					`(2, 10, 'value2'), ` +
					`(3, 15, 'value3')`,
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer(
						"test_column",
						postgres.NewConditionAnonymizer(
							"TEXT",
							gotidus.NewStaticAnonymizer("static_value", "TEXT"),
							postgres.NewAnonymizationCondition(
								"reference",
								"<",
								"10",
								"INT",
								postgres.NewNullAnonymizer(),
							),
							postgres.NewAnonymizationCondition(
								"reference",
								">",
								"10",
								"INT",
								gotidus.NewNoopAnonymizer(),
							),
						),
					),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 1",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedNullStr := sql.NullString{
							Valid:  false,
							String: "",
						}
						var nullStr sql.NullString

						err := row.Scan(&nullStr)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStructs(nullStr, expectedNullStr, t)
					},
				},
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 2",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "static_value"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 3",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "value3"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "EmailAnonymizer: anonymize if it is an email address",
			setupQueries: []string{
				"CREATE TABLE test_table (id INT, email_or_name TEXT)",
				"INSERT INTO test_table (id, email_or_name) VALUES (1, 'test@example.com'), (2, 'noemail')",
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer("email_or_name", postgres.NewEmailAnonymizer()),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT email_or_name FROM test_table_anonymized WHERE id = 1",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "55502f40dc8b7c7@example.com"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
				{
					Query: "SELECT email_or_name FROM test_table_anonymized WHERE id = 2",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "noemail"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "NullAnonymizer: overwrite with null",
			setupQueries: []string{
				"CREATE TABLE test_table (test_column TEXT)",
				"INSERT INTO test_table (test_column) VALUES ('value_to_be_overwritten')",
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer("test_column", postgres.NewNullAnonymizer()),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedNullStr := sql.NullString{
							Valid:  false,
							String: "",
						}
						var nullStr sql.NullString

						err := row.Scan(&nullStr)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStructs(nullStr, expectedNullStr, t)
					},
				},
			},
		},
		{
			title: "OverlayAnonymizer: print string over value",
			setupQueries: []string{
				"CREATE TABLE test_table (id INT, test_column TEXT)",
				"INSERT INTO test_table (id, test_column) " +
					"VALUES (1, 'value_to_be_overwritten'), (2, 'short_value')",
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer("test_column", postgres.NewOverlayAnonymizer("X", 6, 10)),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 1",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "valueXXXXXXXXXXrwritten"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 2",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "shortXXXXXXXXXX"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "RegexReplaceAnonymizer",
			setupQueries: []string{
				"CREATE TABLE test_table (id INT, test_column TEXT)",
				"INSERT INTO test_table (id, test_column) " +
					"VALUES (1, 'value_to_be_change'), (2, 'unchanged_value')",
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer("test_column", postgres.NewRegexReplaceAnonymizer("to_be", "that_was")),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 1",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "value_that_was_change"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
				{
					Query: "SELECT test_column FROM test_table_anonymized WHERE id = 2",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := "unchanged_value"
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "RemoveJSONKeysAnonymizer",
			setupQueries: []string{
				"CREATE TABLE test_table (json_column TEXT)",
				`INSERT INTO test_table (json_column) ` +
					`VALUES ('{"token":"efgh","secret":"abcd","remaining":"value"}')`,
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer(
						"json_column",
						postgres.NewRemoveJSONKeysAnonymizer([]string{"token", "secret"}),
					),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT json_column FROM test_table_anonymized",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						expectedStr := `{"remaining":"value"}`
						var str string

						err := row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						testutils.CompareStrings(str, expectedStr, t)
					},
				},
			},
		},
		{
			title: "TextAnonymizer",
			setupQueries: []string{
				"CREATE TABLE test_table (test_column TEXT)",
				`INSERT INTO test_table (test_column) ` +
					`VALUES ('scramble THIS text so that it cannot be understood 1234')`,
			},
			anonymizationConfig: map[string]*gotidus.Table{
				"test_table": gotidus.NewTable().
					AddAnonymizer("test_column", postgres.NewTextAnonymizer()),
			},
			queryChecks: []queryCheck{
				{
					Query: "SELECT test_column FROM test_table_anonymized",
					ExpectationFunc: func(row *sql.Row, t *testing.T) {
						blockWidths := []int{8, 4, 4, 2, 4, 2, 6, 2, 10, 4}

						patterns := make([]string, 0, len(blockWidths))
						basePattern := "[A-Za-z0-9ÖöÜüÄäß]{%d}"

						for _, blockWidth := range blockWidths {
							patterns = append(patterns, fmt.Sprintf(basePattern, blockWidth))
						}

						r, err := regexp.Compile(strings.Join(patterns, " "))

						var str string

						err = row.Scan(&str)
						if err != nil {
							t.Errorf("Failed to retrieve value from check query: %+v", err)
						}

						if !r.MatchString(str) {
							t.Errorf("Failed to match regex %s: %s", r, str)
						}
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			for _, query := range c.setupQueries {
				_, err := db.Exec(query)
				if err != nil {
					t.Errorf("Failed to execute query '%s': %+v", query, err)
				}
			}

			generator := gotidus.NewGenerator(postgres.NewQueryBuilder())
			if c.anonymizationConfig != nil {
				for tableName, table := range c.anonymizationConfig {
					generator.AddTable(tableName, table)
				}
			}

			err = generator.CreateViews(db)
			if err != nil {
				t.Errorf("Failed to create views: %+v", err)
			}

			if c.queryChecks != nil {
				for _, check := range c.queryChecks {
					row := db.QueryRow(check.Query)

					check.ExpectationFunc(row, t)
				}
			}

			err = generator.ClearViews(db)
			if err != nil {
				t.Errorf("Failed to clear views: %+v", err)
			}

			resetPGDB(db, t)
		})
	}
}

func resetPGDB(db *sql.DB, t *testing.T) {
	queries := []string{
		"DROP SCHEMA public CASCADE;",
		"CREATE SCHEMA public;",
		"GRANT ALL ON SCHEMA public TO postgres;",
		"GRANT ALL ON SCHEMA public TO public;",
		"COMMENT ON SCHEMA public IS 'standard public schema';",
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			t.Errorf("Failed during schema cleanup on '%s': %+v", query, err)
		}
	}
}
