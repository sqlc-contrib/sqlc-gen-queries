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
			Expect(postsTable.ForeignKeys).To(HaveLen(1))
			fk := postsTable.ForeignKeys[0]
			Expect(fk.Name).To(Equal("fk_posts_user_id"))
			Expect(fk.Columns).To(Equal([]string{"user_id"}))
			Expect(fk.References.Table).To(Equal("users"))
			Expect(fk.References.Columns).To(Equal([]string{"id"}))

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
})
