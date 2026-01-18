package sqlc_test

import (
	"os"
	"path/filepath"

	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	var generator *sqlc.Generator

	BeforeEach(func() {
		dir, err := os.MkdirTemp("", "sqlc-gen-test-*")
		Expect(err).NotTo(HaveOccurred())

		catalog, err := sqlc.LoadCatalog("./catalog_test.json")
		Expect(err).NotTo(HaveOccurred())

		generator = &sqlc.Generator{
			Catalog: catalog,
			Config: &sqlc.Config{
				Version: "2",
				SQL: []sqlc.SQL{
					{
						Schema:  "schema.sql",
						Engine:  "postgresql",
						Queries: dir,
					},
				},
			},
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(generator.Config.SQL[0].Queries)).To(Succeed())
	})

	Describe("Generate", func() {
		It("generates SQL query files for all tables", func() {
			Expect(generator.Generate()).NotTo(HaveOccurred())

			for _, config := range generator.Config.SQL {
				Expect(filepath.Join(config.Queries, "posts.sql")).To(BeAnExistingFile())
				Expect(filepath.Join(config.Queries, "users.sql")).To(BeAnExistingFile())
			}
		})

		It("generates valid SQL content", func() {
			err := generator.Generate()
			Expect(err).NotTo(HaveOccurred())

			for _, config := range generator.Config.SQL {
				path := filepath.Join(config.Queries, "users.sql")
				Expect(path).To(BeAnExistingFile())
				content, err := os.ReadFile(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).NotTo(BeEmpty())
			}
		})

		When("the queries directory does not exist", func() {
			It("returns an error", func() {
				for index := range generator.Config.SQL {
					generator.Config.SQL[index].Queries = "/nonexistent/path"
					Expect(generator.Generate()).To(HaveOccurred())
				}
			})
		})

		Context("with skip_queries configuration", func() {
			BeforeEach(func() {
				dir, err := os.MkdirTemp("", "sqlc-gen-test-exclude-*")
				Expect(err).NotTo(HaveOccurred())

				catalog, err := sqlc.LoadCatalog("./catalog_test.json")
				Expect(err).NotTo(HaveOccurred())

				generator = &sqlc.Generator{
					Catalog: catalog,
					Config: &sqlc.Config{
						Version: "2",
						SQL: []sqlc.SQL{
							{
								Schema:  "schema.sql",
								Engine:  "postgresql",
								Queries: dir,
								SkipQueries: []string{
									"GetUser",
									"DeleteUser",
									"BatchGetUsers",
								},
							},
						},
					},
				}
			})

			It("does not generate skipped queries", func() {
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Verify excluded queries are not present
					Expect(string(content)).NotTo(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).NotTo(ContainSubstring("name: DeleteUser :one"))
					Expect(string(content)).NotTo(ContainSubstring("name: BatchGetUsers :batchone"))

					// Verify other queries are still present
					Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))
					Expect(string(content)).To(ContainSubstring("name: UpdateUser :one"))
				}
			})

			It("generates all queries when skip_queries is empty", func() {
				generator.Config.SQL[0].SkipQueries = []string{}
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Verify all queries are present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: DeleteUser :one"))
					Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))
					Expect(string(content)).To(ContainSubstring("name: UpdateUser :one"))
				}
			})

			It("generates all queries when skip_queries is nil", func() {
				generator.Config.SQL[0].SkipQueries = nil
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Verify all queries are present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: DeleteUser :one"))
				}
			})

			It("handles non-existent query names gracefully", func() {
				generator.Config.SQL[0].SkipQueries = []string{"NonExistentQuery", "AnotherFakeQuery"}
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Verify all real queries are still present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))
				}
			})
		})
	})
})
