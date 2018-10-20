package postgres

import (
	"fmt"
	"strings"
)

// QueryBuilder is the specific implementation of the gotidus.QueryBuilder interface for PostgreSQL.
type QueryBuilder struct{}

// NewQueryBuilder initializes a new QueryBuilder object.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

const listViewsQuery string = `
  SELECT
    viewname
  FROM pg_catalog.pg_views
  WHERE viewname ILIKE '%' || $1
  ORDER BY viewname ASC`

// ListViewsQuery returns the query for listing existing views.
// It requires passing the view postfix on query execution.
func (qb *QueryBuilder) ListViewsQuery() string {
	return listViewsQuery
}

const dropViewQueryTemplate string = "DROP VIEW IF EXISTS %s"

// DropViewQuery returns the query for removing the view for which the name is given.
func (qb *QueryBuilder) DropViewQuery(viewName string) string {
	return fmt.Sprintf(dropViewQueryTemplate, viewName)
}

const listTablesQuery string = `
  SELECT tablename
  FROM pg_catalog.pg_tables
  WHERE schemaname = CURRENT_SCHEMA
  ORDER BY tablename ASC`

// ListTablesQuery returns the query for listing existing tables.
func (qb *QueryBuilder) ListTablesQuery() string {
	return listTablesQuery
}

const listColumnsQuery string = `
  SELECT
    column_name
  FROM information_schema.columns
  WHERE table_name = $1
  ORDER BY ordinal_position ASC`

// ListColumnsQuery returns the query for listing existing columns.
// It requires passing the table name on query execution for which the columns should be listed.
func (qb *QueryBuilder) ListColumnsQuery() string {
	return listColumnsQuery
}

const createViewQueryTemplate string = `
  CREATE OR REPLACE VIEW %s AS
    SELECT %s
    FROM %s`

// CreateViewQuery returns the query for creating a view.
// It builds the query using the view name, table name and the data for the selectable columns.
func (qb *QueryBuilder) CreateViewQuery(
	viewName string,
	tableName string,
	columns []string,
) string {
	joinedSelectString := strings.Join(columns, ", ")

	return fmt.Sprintf(
		createViewQueryTemplate,
		viewName,
		joinedSelectString,
		tableName,
	)
}
