package sqlc_test

import (
	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Catalog", func() {
	Describe("LoadCatalog", func() {
		It("loads and parses a complete catalog file correctly", func() {
			catalog, err := sqlc.LoadCatalog("./catalog_test.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(catalog).NotTo(BeNil())

			// Verify schemas
			Expect(catalog.Schemas).To(HaveLen(2))
			Expect(catalog.Schemas[0].Name).To(Equal("public"))
			Expect(catalog.Schemas[1].Name).To(Equal("analytics"))

			// Verify schema attributes
			publicSchema := catalog.Schemas[0]
			Expect(publicSchema.Comment).To(Equal("Public schema"))

			// Verify tables
			Expect(publicSchema.Tables).To(HaveLen(2))
			Expect(publicSchema.Tables[0].Name).To(Equal("users"))
			Expect(publicSchema.Tables[1].Name).To(Equal("posts"))

			// Verify users table
			usersTable := publicSchema.Tables[0]
			Expect(usersTable.Comment).To(Equal("User accounts"))

			// Verify columns
			Expect(usersTable.Columns).To(HaveLen(3))
			Expect(usersTable.Columns[0].Name).To(Equal("id"))
			Expect(usersTable.Columns[0].Type).To(Equal("integer"))
			Expect(usersTable.Columns[0].Null).To(BeFalse())
			Expect(usersTable.Columns[0].Comment).To(Equal("Primary key"))

			Expect(usersTable.Columns[1].Name).To(Equal("email"))
			Expect(usersTable.Columns[1].Type).To(Equal("varchar(255)"))
			Expect(usersTable.Columns[1].Null).To(BeFalse())

			Expect(usersTable.Columns[2].Name).To(Equal("name"))
			Expect(usersTable.Columns[2].Type).To(Equal("text"))
			Expect(usersTable.Columns[2].Null).To(BeTrue())

			// Verify primary key
			Expect(usersTable.PrimaryKey).NotTo(BeNil())
			Expect(usersTable.PrimaryKey.Name).To(Equal("users_pkey"))
			Expect(usersTable.PrimaryKey.Unique).To(BeTrue())
			Expect(usersTable.PrimaryKey.Parts).To(HaveLen(1))
			Expect(usersTable.PrimaryKey.Parts[0].Column).To(Equal("id"))
			Expect(usersTable.PrimaryKey.Parts[0].Desc).To(BeFalse())

			// Verify indexes
			Expect(usersTable.Indexes).To(HaveLen(1))
			emailIndex := usersTable.Indexes[0]
			Expect(emailIndex.Name).To(Equal("idx_users_email"))
			Expect(emailIndex.Unique).To(BeTrue())
			Expect(emailIndex.Parts).To(HaveLen(1))
			Expect(emailIndex.Parts[0].Column).To(Equal("email"))

			// Verify posts table
			postsTable := publicSchema.Tables[1]

			// Verify foreign keys
			Expect(postsTable.ForeignKeys).To(HaveLen(2))

			// First FK: NOT NULL user_id (for INNER JOIN)
			fk := postsTable.ForeignKeys[0]
			Expect(fk.Name).To(Equal("fk_posts_user_id"))
			Expect(fk.Columns).To(Equal([]string{"user_id"}))
			Expect(fk.References.Table).To(Equal("users"))
			Expect(fk.References.Columns).To(Equal([]string{"id"}))

			// Second FK: NULLABLE author_id (for LEFT JOIN)
			fkAuthor := postsTable.ForeignKeys[1]
			Expect(fkAuthor.Name).To(Equal("fk_posts_author_id"))
			Expect(fkAuthor.Columns).To(Equal([]string{"author_id"}))
			Expect(fkAuthor.References.Table).To(Equal("users"))
			Expect(fkAuthor.References.Columns).To(Equal([]string{"id"}))

			// Verify expression-based indexes
			Expect(postsTable.Indexes).To(HaveLen(2))
			exprIndex := postsTable.Indexes[1]
			Expect(exprIndex.Name).To(Equal("idx_posts_title_expr"))
			Expect(exprIndex.Parts).To(HaveLen(1))
			Expect(exprIndex.Parts[0].Expr).To(Equal("lower(title)"))
			Expect(exprIndex.Parts[0].Column).To(BeEmpty())
		})

		When("the file does not exist", func() {
			It("returns an error", func() {
				catalog, err := sqlc.LoadCatalog("./catalog_test.yaml")
				Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
				Expect(catalog).To(BeNil())
			})
		})

		When("the file is invalid JSON", func() {
			It("returns an error", func() {
				catalog, err := sqlc.LoadCatalog("./catalog.go")
				Expect(err).To(HaveOccurred())
				Expect(catalog).To(BeNil())
			})
		})
	})

	Describe("Table methods", func() {
		var catalog *sqlc.Catalog
		var usersTable *sqlc.Table

		BeforeEach(func() {
			var err error
			catalog, err = sqlc.LoadCatalog("./catalog_test.json")
			Expect(err).NotTo(HaveOccurred())
			usersTable = &catalog.Schemas[0].Tables[0]
		})

		Describe("GetColumn", func() {
			It("returns the column when it exists", func() {
				column := usersTable.GetColumn("email")
				Expect(column).NotTo(BeNil())
				Expect(column.Name).To(Equal("email"))
				Expect(column.Type).To(Equal("varchar(255)"))
			})

			It("returns nil when the column does not exist", func() {
				column := usersTable.GetColumn("nonexistent")
				Expect(column).To(BeNil())
			})
		})

		Describe("GetNonUniqueIndexes", func() {
			var postsTable *sqlc.Table

			BeforeEach(func() {
				postsTable = &catalog.Schemas[0].Tables[1]
			})

			It("returns non-unique indexes", func() {
				keys := postsTable.GetNonUniqueIndexes()
				Expect(keys).To(HaveLen(1))
				Expect(keys[0].Name).To(Equal("idx_posts_user_id"))
				Expect(keys[0].Unique).To(BeFalse())
			})

			It("returns empty slice when there are no non-unique indexes", func() {
				keys := usersTable.GetNonUniqueIndexes()
				Expect(keys).To(BeEmpty())
			})
		})

		Describe("GetUniqueKeys", func() {
			It("returns primary key and unique indexes", func() {
				keys := usersTable.GetUniqueKeys()
				Expect(keys).To(HaveLen(2))

				// First should be primary key (normalized to "primary key")
				Expect(keys[0].Name).To(Equal("primary key"))
				Expect(keys[0].Unique).To(BeTrue())

				// Second should be unique index
				Expect(keys[1].Name).To(Equal("idx_users_email"))
				Expect(keys[1].Unique).To(BeTrue())
			})

			It("returns only indexes when there is no primary key", func() {
				// Create a table without primary key
				table := &sqlc.Table{
					Name: "test_table",
					Indexes: []sqlc.Index{
						{Name: "idx_unique", Unique: true},
						{Name: "idx_normal", Unique: false},
					},
				}

				keys := table.GetUniqueKeys()
				Expect(keys).To(HaveLen(1))
				Expect(keys[0].Name).To(Equal("idx_unique"))
			})
		})
	})

	Describe("Condition", func() {
		It("generates correct SQL condition string", func() {
			column := &sqlc.Column{
				Name: "user_id",
				Type: "integer",
				Null: false,
			}

			arg := &sqlc.Argument{
				Column: column,
			}

			condition := &sqlc.ArgumentCondition{
				Column:   column,
				Argument: arg,
			}

			Expect(condition.String()).To(Equal("user_id = sqlc.arg(user_id)"))
		})
	})

	Describe("Argument", func() {
		Context("when column is nullable", func() {
			It("uses sqlc.narg", func() {
				column := &sqlc.Column{
					Name: "description",
					Type: "text",
					Null: true,
				}

				arg := &sqlc.Argument{
					Column: column,
				}

				Expect(arg.String()).To(Equal("sqlc.narg(description)"))
			})
		})

		Context("when column is not nullable", func() {
			It("uses sqlc.arg", func() {
				column := &sqlc.Column{
					Name: "id",
					Type: "integer",
					Null: false,
				}

				arg := &sqlc.Argument{
					Column: column,
				}

				Expect(arg.String()).To(Equal("sqlc.arg(id)"))
			})
		})
	})

	Describe("CompositeCondition", func() {
		Describe("AddColumn", func() {
			It("adds a condition for the column", func() {
				composite := &sqlc.CompositeCondition{
					Operator: "AND",
				}

				column := &sqlc.Column{
					Name: "id",
					Type: "integer",
					Null: false,
				}

				composite.AddColumn(column)

				Expect(composite.Conditions).To(HaveLen(1))
			})
		})

		Describe("String", func() {
			It("joins conditions with the operator", func() {
				column1 := &sqlc.Column{
					Name: "id",
					Type: "integer",
					Null: false,
				}

				column2 := &sqlc.Column{
					Name: "email",
					Type: "text",
					Null: false,
				}

				composite := &sqlc.CompositeCondition{
					Operator: "AND",
				}

				composite.AddColumn(column1)
				composite.AddColumn(column2)

				result := composite.String()
				Expect(result).To(Equal("id = sqlc.arg(id) AND email = sqlc.arg(email)"))
			})

			It("returns empty string when there are no conditions", func() {
				composite := &sqlc.CompositeCondition{
					Operator: "AND",
				}

				Expect(composite.String()).To(BeEmpty())
			})

			It("handles OR operator", func() {
				column1 := &sqlc.Column{
					Name: "id",
					Type: "integer",
					Null: false,
				}

				column2 := &sqlc.Column{
					Name: "email",
					Type: "text",
					Null: false,
				}

				composite := &sqlc.CompositeCondition{
					Operator: "OR",
				}

				composite.AddColumn(column1)
				composite.AddColumn(column2)

				result := composite.String()
				Expect(result).To(Equal("id = sqlc.arg(id) OR email = sqlc.arg(email)"))
			})
		})

		Describe("GetNonPrimaryKeyColumns", func() {
			It("excludes primary key columns", func() {
				catalog, err := sqlc.LoadCatalog("./catalog_test.json")
				Expect(err).NotTo(HaveOccurred())
				Expect(catalog).NotTo(BeNil())
				usersTable := &catalog.Schemas[0].Tables[0]
				nonPKColumns := usersTable.GetNonPrimaryKeyColumns()

				// Should have 2 columns (email, name) but not id
				Expect(nonPKColumns).To(HaveLen(2))

				columnNames := make([]string, len(nonPKColumns))
				for i, col := range nonPKColumns {
					columnNames[i] = col.Name
				}

				Expect(columnNames).To(ConsistOf("email", "name"))
				Expect(columnNames).NotTo(ContainElement("id"))
			})

			It("returns all columns when table has no primary key", func() {
				tableWithoutPK := &sqlc.Table{
					Name: "test_table",
					Columns: []sqlc.Column{
						{Name: "id", Type: "integer"},
						{Name: "name", Type: "text"},
						{Name: "email", Type: "varchar(255)"},
					},
				}

				nonPKColumns := tableWithoutPK.GetNonPrimaryKeyColumns()
				Expect(nonPKColumns).To(HaveLen(3))

				columnNames := make([]string, len(nonPKColumns))
				for i, col := range nonPKColumns {
					columnNames[i] = col.Name
				}

				Expect(columnNames).To(ConsistOf("id", "name", "email"))
			})

			It("handles composite primary keys", func() {
				tableWithCompositePK := &sqlc.Table{
					Name: "junction_table",
					Columns: []sqlc.Column{
						{Name: "user_id", Type: "integer"},
						{Name: "role_id", Type: "integer"},
						{Name: "created_at", Type: "timestamp"},
						{Name: "updated_at", Type: "timestamp"},
					},
					PrimaryKey: &sqlc.Index{
						Name: "junction_pkey",
						Parts: []sqlc.IndexPart{
							{Column: "user_id"},
							{Column: "role_id"},
						},
					},
				}

				nonPKColumns := tableWithCompositePK.GetNonPrimaryKeyColumns()
				Expect(nonPKColumns).To(HaveLen(2))

				columnNames := make([]string, len(nonPKColumns))
				for i, col := range nonPKColumns {
					columnNames[i] = col.Name
				}

				Expect(columnNames).To(ConsistOf("created_at", "updated_at"))
				Expect(columnNames).NotTo(ContainElement("user_id"))
				Expect(columnNames).NotTo(ContainElement("role_id"))
			})

			It("returns empty slice for table with only primary key", func() {
				tableWithOnlyPK := &sqlc.Table{
					Name: "simple_table",
					Columns: []sqlc.Column{
						{Name: "id", Type: "integer"},
					},
					PrimaryKey: &sqlc.Index{
						Name: "simple_pkey",
						Parts: []sqlc.IndexPart{
							{Column: "id"},
						},
					},
				}

				nonPKColumns := tableWithOnlyPK.GetNonPrimaryKeyColumns()
				Expect(nonPKColumns).To(BeEmpty())
			})

			It("preserves column order", func() {
				tableWithOrderedColumns := &sqlc.Table{
					Name: "ordered_table",
					Columns: []sqlc.Column{
						{Name: "id", Type: "integer"},
						{Name: "name", Type: "text"},
						{Name: "email", Type: "varchar(255)"},
						{Name: "created_at", Type: "timestamp"},
						{Name: "updated_at", Type: "timestamp"},
					},
					PrimaryKey: &sqlc.Index{
						Name: "ordered_pkey",
						Parts: []sqlc.IndexPart{
							{Column: "id"},
						},
					},
				}

				nonPKColumns := tableWithOrderedColumns.GetNonPrimaryKeyColumns()
				Expect(nonPKColumns).To(HaveLen(4))

				// Verify order is preserved (excluding id)
				Expect(nonPKColumns[0].Name).To(Equal("name"))
				Expect(nonPKColumns[1].Name).To(Equal("email"))
				Expect(nonPKColumns[2].Name).To(Equal("created_at"))
				Expect(nonPKColumns[3].Name).To(Equal("updated_at"))
			})

			It("handles table with no columns", func() {
				emptyTable := &sqlc.Table{
					Name:    "empty_table",
					Columns: []sqlc.Column{},
				}

				nonPKColumns := emptyTable.GetNonPrimaryKeyColumns()
				Expect(nonPKColumns).To(BeEmpty())
			})

			It("handles table with nullable and non-nullable columns", func() {
				tableWithMixedColumns := &sqlc.Table{
					Name: "mixed_table",
					Columns: []sqlc.Column{
						{Name: "id", Type: "integer", Null: false},
						{Name: "name", Type: "text", Null: true},
						{Name: "email", Type: "varchar(255)", Null: false},
						{Name: "description", Type: "text", Null: true},
					},
					PrimaryKey: &sqlc.Index{
						Name: "mixed_pkey",
						Parts: []sqlc.IndexPart{
							{Column: "id"},
						},
					},
				}

				nonPKColumns := tableWithMixedColumns.GetNonPrimaryKeyColumns()
				Expect(nonPKColumns).To(HaveLen(3))

				// Verify column properties are preserved
				Expect(nonPKColumns[0].Name).To(Equal("name"))
				Expect(nonPKColumns[0].Null).To(BeTrue())

				Expect(nonPKColumns[1].Name).To(Equal("email"))
				Expect(nonPKColumns[1].Null).To(BeFalse())

				Expect(nonPKColumns[2].Name).To(Equal("description"))
				Expect(nonPKColumns[2].Null).To(BeTrue())
			})
		})
	})
})
