package sqlc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/inflect"
	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc/template"
)

// Generator generates an SQL queries
type Generator struct {
	Config  *Config
	Catalog *Catalog
}

// Generate generates the queries based on the configuration.
func (x *Generator) Generate() error {
	// Context holds data for template execution
	type Context struct {
		Engine string
		Schema string
		Table  *Table
	}

	opts := map[string]any{
		// Table Functions
		"table_ref": func(fk ForeignKey) string {
			column := fk.Columns[0]
			// If the column name starts with the referenced table name, use that as the role
			if strings.HasPrefix(column, fk.References.Table+"_") {
				return fk.References.Table
			}

			suffixes := []string{"_id", "_fk", "_ref", "_key"}
			// Remove common FK suffixes
			for _, suffix := range suffixes {
				column = strings.TrimSuffix(column, suffix)
			}

			return column
		},
		"table_name": func(table string, kind string) string {
			switch kind {
			case "many":
				return inflect.Camelize(inflect.Pluralize(table))
			case "one":
				return inflect.Camelize(inflect.Singularize(table))
			}
			return ""
		},
		"table_join": func(table Table, fk ForeignKey) string {
			// We assume the foreign key references another table in the same catalog
			// Otherwise it will return nil and likely panic
			tableRef := x.Catalog.GetTable(fk.References.Table)

			jtype := "INNER JOIN"
			// Determine join type: LEFT JOIN if any FK column is nullable
			for _, name := range fk.Columns {
				if column := table.GetColumn(name); column != nil && column.Null {
					jtype = "LEFT JOIN"
					break
				}
			}

			condition := &CompositeCondition{Operator: "AND"}
			// Build the join condition
			for i := range fk.Columns {
				column := fk.Columns[i]
				columnRef := fk.References.Columns[i]
				// We assume the columns exist; otherwise it will likely panic
				condition.AddColumnRef(
					&ColumnRef{
						Table: &table,
						Name:  column,
					},
					&ColumnRef{
						Table: tableRef,
						Name:  columnRef,
					},
				)
			}

			return fmt.Sprintf("%s %s ON %s", jtype, tableRef.Name, condition.String())
		},
		"table_embed": func(table string) string {
			return fmt.Sprintf("sqlc.embed(%s)", table)
		},
		// Query Functions
		"query_condition": func(table Table, index *Index) string {
			condition := &CompositeCondition{Operator: "AND"}
			// Build the condition clause
			for _, part := range index.Parts {
				if column := table.GetColumn(part.Column); column != nil {
					condition.AddColumn(column)
				}
			}
			return condition.String()
		},
		"query_index": func(index *Index) string {
			// Don't add suffix for primary key lookups
			if index.Name == "primary key" {
				return ""
			}

			var items []string
			// Build the suffix based on index parts
			for _, part := range index.Parts {
				items = append(items, inflect.Camelize(part.Column))
			}
			return "By" + strings.Join(items, "And")
		},
		"query_argument": func(column Column) string {
			argument := Argument{
				Column: &column,
			}
			return argument.String()
		},
		// Pagination functions
		"page_start": func(table Table) string {
			if table.PrimaryKey != nil && len(table.PrimaryKey.Parts) > 0 {
				argument := &Argument{
					Column: &Column{
						Name: "page_start",
						Type: "text",
						Null: true,
					},
				}
				column := table.PrimaryKey.Parts[0].Column
				return fmt.Sprintf("(%v::text IS NULL OR %s::text > %v::text)", argument, column, argument)
			}
			return ""
		},
		"page_order": func(table Table) string {
			if table.PrimaryKey != nil && len(table.PrimaryKey.Parts) > 0 {
				return table.PrimaryKey.Parts[0].Column
			}
			return ""
		},
	}

	// Open the template file
	template, err := template.Open("template.sql.tmpl", opts)
	if err != nil {
		return err
	}

	for _, config := range x.Config.SQL {
		if err := os.MkdirAll(config.Queries, os.ModePerm); err != nil {
			return err
		}

		for _, schema := range x.Catalog.Schemas {
			for _, table := range schema.Tables {
				file, err := os.Create(
					filepath.Join(config.Queries,
						fmt.Sprintf("%s.sql", table.Name)),
				)
				if err != nil {
					return err
				}
				//nolint:all
				defer file.Close()

				ctx := Context{
					Engine: config.Engine,
					Schema: schema.Name,
					Table:  &table,
				}
				// Execute template
				if err := template.Execute(file, ctx); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
