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
		"cursor": func(table Table) string {
			if table.PrimaryKey != nil && len(table.PrimaryKey.Parts) > 0 {
				argument := &Argument{
					Column: &Column{
						Name: "page_start",
						Type: "text",
					},
				}
				column := table.PrimaryKey.Parts[0].Column

				return fmt.Sprintf("(%v IS NULL OR %s::text > %v)", argument, column, argument)
			}
			return ""
		},
		"by_index": func(table Table, index *Index) string {
			if table.PrimaryKey.Name == index.Name {
				return ""
			}

			var items []string
			// Prepare the clause
			for _, part := range index.Parts {
				items = append(items, inflect.Camelize(part.Column))
			}

			return "By" + strings.Join(items, "And")
		},
		"by_many": func(table Table) string {
			return inflect.Camelize(inflect.Pluralize(table.Name))
		},
		"by_one": func(table Table) string {
			return inflect.Camelize(inflect.Singularize(table.Name))
		},
		"where": func(table Table, index *Index) string {
			where := &CompositeCondition{Operator: "AND"}
			// Build the where clause
			for _, part := range index.Parts {
				if column := table.GetColumn(part.Column); column != nil {
					where.AddColumn(column)
				}
			}

			return where.String()
		},
		"order": func(table Table) string {
			if table.PrimaryKey != nil && len(table.PrimaryKey.Parts) > 0 {
				return table.PrimaryKey.Parts[0].Column
			}
			return ""
		},
		"arg": func(column Column) string {
			arggument := Argument{
				Column: &column,
			}
			return arggument.String()
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
