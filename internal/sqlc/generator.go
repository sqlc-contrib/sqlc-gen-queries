package sqlc

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-openapi/inflect"
	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc/template"
)

// blank matches two or more consecutive blank lines.
var blank = regexp.MustCompile(`\n{3,}`)

func init() {
	inflect.AddSingular("quota", "quota")
	inflect.AddPlural("quota", "quotas")
}

// Generator generates an SQL queries
type Generator struct {
	Config  *Config
	Catalog *Catalog
}

// Generate generates the queries based on the configuration.
func (x *Generator) Generate() error {
	// Context holds data for template execution
	type Context struct {
		Engine              string
		Schema              string
		Table               *Table
		QueryInclude        map[string]bool
		QueryExclude        map[string]bool
		InsertColumnExclude map[string]bool
		UpdateColumnExclude map[string]bool
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
		"query_argument": func(column Column) string {
			argument := Argument{
				Column: &column,
			}
			return argument.String()
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
		// Pagination functions
		"query_order": func(table Table) string {
			if table.PrimaryKey == nil {
				return ""
			}
			cols := make([]string, 0, len(table.PrimaryKey.Parts))
			for _, p := range table.PrimaryKey.Parts {
				cols = append(cols, p.Column)
			}
			return strings.Join(cols, ", ")
		},
		// Foreign key index check
		"is_fk_index": func(table Table, index *Index) bool {
			return table.IsForeignKeyIndex(index)
		},
		// Query selection: a query renders when it belongs to the default set
		// or is explicitly included, and never when excluded (exclude wins).
		"should_generate": func(ctx Context, queryName string, isDefault bool) bool {
			if ctx.QueryExclude[queryName] {
				return false
			}
			return isDefault || ctx.QueryInclude[queryName]
		},
		"insert_columns": func(ctx Context) []Column {
			columns := make([]Column, 0, len(ctx.Table.Columns))
			for _, column := range ctx.Table.Columns {
				if columnSelected(ctx.InsertColumnExclude, ctx.Schema, ctx.Table.Name, column.Name) {
					columns = append(columns, column)
				}
			}
			return columns
		},
		"update_columns": func(ctx Context) []Column {
			tableColumns := ctx.Table.GetNonPrimaryKeyColumns()
			columns := make([]Column, 0, len(tableColumns))
			for _, column := range tableColumns {
				if columnSelected(ctx.UpdateColumnExclude, ctx.Schema, ctx.Table.Name, column.Name) {
					columns = append(columns, column)
				}
			}
			return columns
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

		queryInclude := config.GetQueryIncludeSet()
		queryExclude := config.GetQueryExcludeSet()
		insertColumnExclude := config.GetInsertColumnExcludeSet()
		updateColumnExclude := config.GetUpdateColumnExcludeSet()
		include := config.GetIncludeSet()
		exclude := config.GetExcludeSet()

		for _, schema := range x.Catalog.Schemas {
			for _, table := range schema.Tables {
				if !tableSelected(include, exclude, schema.Name, table.Name) {
					continue
				}

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
					Engine:              config.Engine,
					Schema:              schema.Name,
					Table:               &table,
					QueryInclude:        queryInclude,
					QueryExclude:        queryExclude,
					InsertColumnExclude: insertColumnExclude,
					UpdateColumnExclude: updateColumnExclude,
				}
				// Execute template into buffer, then squeeze blank lines
				var buffer bytes.Buffer
				if err := template.Execute(&buffer, ctx); err != nil {
					return err
				}

				data := blank.ReplaceAll(buffer.Bytes(), []byte("\n\n"))
				if _, err := file.Write(data); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
