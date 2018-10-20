package gotidus

import (
	"database/sql"
	"fmt"
)

// QueryBuilder is the interface used to implement support for different databases.
type QueryBuilder interface {
	ListViewsQuery() string
	DropViewQuery(viewName string) string

	ListTablesQuery() string

	ListColumnsQuery() string

	CreateViewQuery(viewName string, tableName string, columns []string) string
}

// DefaultViewPostfix defines the postfix given to views to distinguish them from the table names.
const DefaultViewPostfix = "anonymized"

// NewGenerator initializes a new Generator object.
// It requires a QueryBuilder object and can be enhanced with GeneratorOption functions.
func NewGenerator(queryBuilder QueryBuilder, options ...GeneratorOption) *Generator {
	generator := &Generator{
		queryBuilder: queryBuilder,
		tables:       make(map[string]*Table),
		viewPostfix:  DefaultViewPostfix,
	}

	for _, option := range options {
		option(generator)
	}

	return generator
}

// Generator is the type orchestrating the view clearing and creation,
// based on the table config.
type Generator struct {
	queryBuilder QueryBuilder
	tables       map[string]*Table
	viewPostfix  string
}

// AddTable adds a Table configuration to the generator with the given name.
// If this function is called again with the same name, it will overwrite the existing table.
func (g *Generator) AddTable(name string, table *Table) *Generator {
	g.tables[name] = table

	return g
}

// GetTable retrieves a Table from the config.
// If a Table was configured for the given name, that Table object will be returned.
// If no Table was configured for the given name, a blank table configuration is returned.
func (g *Generator) GetTable(name string) *Table {
	table, ok := g.tables[name]
	if ok {
		return table
	}

	return NewTable()
}

func (g *Generator) loopExistingViews(db *sql.DB, viewFunc func(viewName string) error) error {
	rows, err := db.Query(
		g.queryBuilder.ListViewsQuery(),
		g.viewPostfix,
	)
	if err != nil {
		return fmt.Errorf("Failed to select views: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var viewName string

		if err := rows.Scan(&viewName); err != nil {
			return fmt.Errorf("Failed to scan viewname: %+v", err)
		}

		if err := viewFunc(viewName); err != nil {
			return err
		}
	}

	return nil
}

// ClearViews removes any potentially existing views that exist with the configured postfix.
func (g *Generator) ClearViews(db *sql.DB) error {
	return g.loopExistingViews(
		db,
		func(viewName string) error {
			if _, err := db.Exec(g.queryBuilder.DropViewQuery(viewName)); err != nil {
				return fmt.Errorf("Failed to drop view '%s': %+v", viewName, err)
			}

			return nil
		},
	)
}

func (g *Generator) loopTables(db *sql.DB, tableFunc func(tableName string) error) error {
	tableRows, err := db.Query(g.queryBuilder.ListTablesQuery())
	if err != nil {
		return fmt.Errorf("Failed to select tables: %+v", err)
	}
	defer tableRows.Close()

	for tableRows.Next() {
		var tableName string

		if err := tableRows.Scan(&tableName); err != nil {
			return fmt.Errorf("Failed to scan table name: %+v", err)
		}

		if err := tableFunc(tableName); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) loopColumns(
	db *sql.DB,
	tableName string,
	columnFunc func(columnName string) error,
) error {
	columnRows, err := db.Query(g.queryBuilder.ListColumnsQuery(), tableName)
	if err != nil {
		return fmt.Errorf("Failed to select columns: %+v", err)
	}
	defer columnRows.Close()

	for columnRows.Next() {
		var columnName string

		if err := columnRows.Scan(&columnName); err != nil {
			return err
		}

		if err := columnFunc(columnName); err != nil {
			return err
		}
	}

	return nil
}

// CreateViews creates views named <table_name>_<postfix> for each table that could be found.
// It uses the configuration set before CreateViews was called.
func (g *Generator) CreateViews(db *sql.DB) error {
	return g.loopTables(
		db,
		func(tableName string) error {
			columns := make([]string, 0)

			table := g.GetTable(tableName)

			if err := g.loopColumns(
				db,
				tableName,
				func(columnName string) error {
					anonymizer := table.GetAnonymizer(columnName)

					columns = append(
						columns,
						fmt.Sprintf("%s AS %s", anonymizer.Build(tableName, columnName), columnName),
					)

					return nil
				},
			); err != nil {
				return err
			}

			viewName := g.ViewName(tableName)

			if _, err := db.Exec(
				g.queryBuilder.CreateViewQuery(viewName, tableName, columns),
			); err != nil {
				return fmt.Errorf("Failed to create view '%s': %+v", viewName, err)
			}

			return nil
		},
	)
}

// ViewName builds the view name from the table name and the postfix to <table_name>_<postfix>.
func (g *Generator) ViewName(tableName string) string {
	return fmt.Sprintf("%s_%s", tableName, g.viewPostfix)
}

// GeneratorOption is a function type following the option function pattern.
// It can be used to define methods of configuring the Generator object.
type GeneratorOption func(*Generator)

// WithViewPostfix is a GeneratorOption builder, which allows configuring the view postfix.
func WithViewPostfix(viewPostfix string) GeneratorOption {
	return func(g *Generator) {
		g.viewPostfix = viewPostfix
	}
}
