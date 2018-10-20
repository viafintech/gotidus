package gotidus

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Barzahlen/gotidus/testutils"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewGenerator(t *testing.T) {
	queryBuilder := &mockQueryBuilder{}

	expectedGenerator := &Generator{
		queryBuilder: queryBuilder,
		tables:       make(map[string]*Table),
		viewPostfix:  "bazbaz",
	}

	testutils.CompareStructs(
		NewGenerator(queryBuilder, WithViewPostfix("bazbaz")),
		expectedGenerator,
		t,
	)
}

func TestGeneratorAddAndGetTable(t *testing.T) {
	generator := NewGenerator(&mockQueryBuilder{})

	fooTable := NewTable()
	fooTable.AddAnonymizer("bar", NewNoopAnonymizer())

	// Set table
	generator.AddTable("foo", fooTable)

	testutils.CompareStructs(fooTable, generator.tables["foo"], t)

	// Check loading
	tbl := generator.GetTable("foo")

	testutils.CompareStructs(fooTable, tbl, t)

	var nullTable *Table

	// Check table not set
	testutils.CompareStructs(generator.tables["baz"], nullTable, t)

	// Check default table
	defaultTable := generator.GetTable("baz")
	expectedTable := NewTable()

	testutils.CompareStructs(defaultTable, expectedTable, t)
}

func TestGeneratorClearViews(t *testing.T) {
	queryBuilder := &mockQueryBuilder{}

	cases := []struct {
		title         string
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			title: "view selection fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"viewname"})
				rows.AddRow("foo_anonymized")
				rows.AddRow("foo2_anonymized")

				mock.
					ExpectQuery(queryBuilder.ListViewsQuery()).
					WithArgs("anonymized").
					WillReturnError(errors.New("simulated failure"))
			},
			expectedError: errors.New("Failed to select views: simulated failure"),
		},
		{
			title: "second view removal fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"viewname"})
				rows.AddRow("foo_anonymized")
				rows.AddRow("foo2_anonymized")

				mock.
					ExpectQuery(queryBuilder.ListViewsQuery()).
					WithArgs("anonymized").
					WillReturnRows(rows)

				mock.
					ExpectExec(queryBuilder.DropViewQuery("foo_anonymized")).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.
					ExpectExec(queryBuilder.DropViewQuery("foo2_anonymized")).
					WillReturnError(errors.New("simulated failure"))
			},
			expectedError: errors.New("Failed to drop view 'foo2_anonymized': simulated failure"),
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to initialize DB mock")
			}

			c.setupMock(dbMock)

			generator := NewGenerator(queryBuilder)

			testutils.CompareStructs(
				generator.ClearViews(db),
				c.expectedError,
				t,
			)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Did not execute expected queries")
			}
		})
	}
}

func TestGeneratorCreateViews(t *testing.T) {
	queryBuilder := &mockQueryBuilder{}

	defaultGeneratorFunc := func() *Generator {
		fooTable := NewTable()
		fooTable.AddAnonymizer("bar", NewStaticAnonymizer("var", "TEXT"))
		fooTable.AddAnonymizer("unknown", NewStaticAnonymizer("value", "TEXT"))
		generator := NewGenerator(queryBuilder)
		generator.AddTable("foo", fooTable)

		return generator
	}

	cases := []struct {
		title          string
		buildGenerator func() *Generator
		setupMock      func(sqlmock.Sqlmock)
		expectedError  error
	}{
		{
			title:          "table selection fails",
			buildGenerator: defaultGeneratorFunc,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"tablename"})
				rows.AddRow("foo")
				rows.AddRow("foo2")

				mock.
					ExpectQuery(queryBuilder.ListTablesQuery()).
					WillReturnError(errors.New("simulated failure"))
			},
			expectedError: errors.New("Failed to select tables: simulated failure"),
		},
		{
			title:          "column selection fails",
			buildGenerator: defaultGeneratorFunc,
			setupMock: func(mock sqlmock.Sqlmock) {
				tableRows := sqlmock.NewRows([]string{"tablename"})
				tableRows.AddRow("foo")
				tableRows.AddRow("foo2")

				mock.
					ExpectQuery(queryBuilder.ListTablesQuery()).
					WillReturnRows(tableRows)

				mock.
					ExpectQuery(queryBuilder.ListColumnsQuery()).
					WithArgs("foo").
					WillReturnError(errors.New("simulated failure"))
			},
			expectedError: errors.New("Failed to select columns: simulated failure"),
		},
		{
			title:          "view creation fails",
			buildGenerator: defaultGeneratorFunc,
			setupMock: func(mock sqlmock.Sqlmock) {
				tableRows := sqlmock.NewRows([]string{"tablename"})
				tableRows.AddRow("foo")
				tableRows.AddRow("foo2")

				mock.
					ExpectQuery(queryBuilder.ListTablesQuery()).
					WillReturnRows(tableRows)

				fooColumnRows := sqlmock.NewRows([]string{"columnname"})
				fooColumnRows.AddRow("id")
				fooColumnRows.AddRow("bar")

				mock.
					ExpectQuery(queryBuilder.ListColumnsQuery()).
					WithArgs("foo").
					WillReturnRows(fooColumnRows)

				mock.
					ExpectExec(
						queryBuilder.CreateViewQuery(
							"foo_anonymized",
							"foo",
							[]string{
								"foo.id AS id",
								"'var'::TEXT AS bar",
							},
						),
					).
					WillReturnError(errors.New("simulated failure"))
			},
			expectedError: errors.New("Failed to create view 'foo_anonymized': simulated failure"),
		},
		{
			title:          "view creation succeeds",
			buildGenerator: defaultGeneratorFunc,
			setupMock: func(mock sqlmock.Sqlmock) {
				tableRows := sqlmock.NewRows([]string{"tablename"})
				tableRows.AddRow("foo")
				tableRows.AddRow("foo2")

				mock.
					ExpectQuery(queryBuilder.ListTablesQuery()).
					WillReturnRows(tableRows)

				fooColumnRows := sqlmock.NewRows([]string{"columnname"})
				fooColumnRows.AddRow("id")
				fooColumnRows.AddRow("bar")

				mock.
					ExpectQuery(queryBuilder.ListColumnsQuery()).
					WithArgs("foo").
					WillReturnRows(fooColumnRows)

				mock.
					ExpectExec(
						queryBuilder.CreateViewQuery(
							"foo_anonymized",
							"foo",
							[]string{
								"foo.id AS id",
								"'var'::TEXT AS bar",
							},
						),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				foo2ColumnRows := sqlmock.NewRows([]string{"columnname"})
				foo2ColumnRows.AddRow("id")
				foo2ColumnRows.AddRow("amount")

				mock.
					ExpectQuery(queryBuilder.ListColumnsQuery()).
					WithArgs("foo2").
					WillReturnRows(foo2ColumnRows)

				mock.
					ExpectExec(
						queryBuilder.CreateViewQuery(
							"foo2_anonymized",
							"foo2",
							[]string{
								"foo2.id AS id",
								"foo2.amount AS amount",
							},
						),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to initialize DB mock")
			}

			c.setupMock(dbMock)

			generator := c.buildGenerator()

			testutils.CompareStructs(
				generator.CreateViews(db),
				c.expectedError,
				t,
			)

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("Did not execute expected queries")
			}
		})
	}
}

func TestGeneratorViewName(t *testing.T) {
	cases := []struct {
		title     string
		tableName string
		options   []GeneratorOption

		expectedViewName string
	}{
		{
			title:     "with default post",
			tableName: "foo",

			expectedViewName: "foo_anonymized",
		},
		{
			title:     "with default post",
			tableName: "bar",
			options: []GeneratorOption{
				WithViewPostfix("bazbaz"),
			},

			expectedViewName: "bar_bazbaz",
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			g := NewGenerator(&mockQueryBuilder{}, c.options...)

			viewName := g.ViewName(c.tableName)

			testutils.CompareStrings(viewName, c.expectedViewName, t)
		})
	}
}

type mockQueryBuilder struct{}

func (mqb *mockQueryBuilder) ListViewsQuery() string {
	return "list_view_query"
}
func (mqb *mockQueryBuilder) DropViewQuery(viewName string) string {
	return fmt.Sprintf("drop_view_query:%s", viewName)
}

func (mqb *mockQueryBuilder) ListTablesQuery() string {
	return "list_tables_query"
}

func (mqb *mockQueryBuilder) ListColumnsQuery() string {
	return "list_columns_query"
}

func (mqb *mockQueryBuilder) CreateViewQuery(
	viewName string,
	tableName string,
	columns []string,
) string {
	columnsString := strings.Join(columns, "|")

	return fmt.Sprintf("create_view_query:%s;%s;%s", viewName, tableName, columnsString)
}
