package postgres

import "testing"

func TestQueryBuilderQueries(t *testing.T) {
	queryBuilder := NewQueryBuilder()

	cases := []struct {
		title         string
		query         string
		expectedQuery string
	}{
		{
			title: "list views query",
			query: queryBuilder.ListViewsQuery(),
			expectedQuery: `
  SELECT
    viewname
  FROM pg_catalog.pg_views
  WHERE viewname ILIKE '%' || $1
  ORDER BY viewname ASC`,
		},
		{
			title:         "drop view query",
			query:         queryBuilder.DropViewQuery("transactions_anonymized"),
			expectedQuery: "DROP VIEW IF EXISTS transactions_anonymized",
		},
		{
			title: "list tables query",
			query: queryBuilder.ListTablesQuery(),
			expectedQuery: `
  SELECT tablename
  FROM pg_catalog.pg_tables
  WHERE schemaname = CURRENT_SCHEMA
  ORDER BY tablename ASC`,
		},
		{
			title: "list columns query",
			query: queryBuilder.ListColumnsQuery(),
			expectedQuery: `
  SELECT
    column_name
  FROM information_schema.columns
  WHERE table_name = $1
  ORDER BY ordinal_position ASC`,
		},
		{
			title: "create view query",
			query: queryBuilder.CreateViewQuery(
				"transactions_anonymized",
				"transactions",
				[]string{"id AS id", "amount AS amount"},
			),
			expectedQuery: `
  CREATE OR REPLACE VIEW transactions_anonymized AS
    SELECT id AS id, amount AS amount
    FROM transactions`,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			compareStrings(c.query, c.expectedQuery, t)
		})
	}
}
