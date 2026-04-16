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

		It("generates valid SQL content with default PK queries", func() {
			err := generator.Generate()
			Expect(err).NotTo(HaveOccurred())

			for _, config := range generator.Config.SQL {
				path := filepath.Join(config.Queries, "users.sql")
				Expect(path).To(BeAnExistingFile())
				content, err := os.ReadFile(path)
				Expect(err).NotTo(HaveOccurred())

				// Default PK CRUD queries are always present
				Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
				Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))
				Expect(string(content)).To(ContainSubstring("name: UpdateUser :one"))
				Expect(string(content)).To(ContainSubstring("name: DeleteUser :one"))

				// Default List queries are present
				Expect(string(content)).To(ContainSubstring("name: ListUsers :many"))

				// Opt-in queries are not present without config
				Expect(string(content)).NotTo(ContainSubstring("name: CopyUsers :copyfrom"))
			}

			// Verify posts.sql FK-based list queries
			for _, config := range generator.Config.SQL {
				path := filepath.Join(config.Queries, "posts.sql")
				content, err := os.ReadFile(path)
				Expect(err).NotTo(HaveOccurred())

				// FK-matching index list query is present by default
				Expect(string(content)).To(ContainSubstring("name: ListPostsByUserId :many"))
				// Non-FK index list query is not present without opt-in
				Expect(string(content)).NotTo(ContainSubstring("name: ListPostsByTitle :many"))
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

		Context("with queries configuration", func() {
			BeforeEach(func() {
				dir, err := os.MkdirTemp("", "sqlc-gen-test-queries-*")
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
								Codegen: []sqlc.Codegen{
									{
										Plugin: "gen-queries",
										Out:    dir,
										Options: sqlc.CodegenOptions{
											Queries: []string{
												"CopyUsers",
												"ListPostsByTitle",
											},
										},
									},
								},
							},
						},
					},
				}
			})

			It("generates opt-in queries when listed", func() {
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Default PK CRUD queries are always present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))
					Expect(string(content)).To(ContainSubstring("name: UpdateUser :one"))
					Expect(string(content)).To(ContainSubstring("name: DeleteUser :one"))

					// Default List queries are present
					Expect(string(content)).To(ContainSubstring("name: ListUsers :many"))

					// Opt-in queries that were listed are present
					Expect(string(content)).To(ContainSubstring("name: CopyUsers :copyfrom"))
				}

				// Verify posts.sql opt-in list queries
				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "posts.sql")
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// FK-matching index list query is present by default
					Expect(string(content)).To(ContainSubstring("name: ListPostsByUserId :many"))
					// Non-FK index list query is present because it was opted in
					Expect(string(content)).To(ContainSubstring("name: ListPostsByTitle :many"))
				}
			})

			It("only generates default queries when queries is empty", func() {
				generator.Config.SQL[0].Codegen[0].Options.Queries = []string{}
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Default PK CRUD queries are present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))

					// Default List queries are present
					Expect(string(content)).To(ContainSubstring("name: ListUsers :many"))

					// Opt-in queries are not present
					Expect(string(content)).NotTo(ContainSubstring("name: CopyUsers :copyfrom"))
				}

				// Verify posts.sql FK-based list queries
				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "posts.sql")
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(ContainSubstring("name: ListPostsByUserId :many"))
					Expect(string(content)).NotTo(ContainSubstring("name: ListPostsByTitle :many"))
				}
			})

			It("only generates default queries when codegen is nil", func() {
				generator.Config.SQL[0].Codegen = nil
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Default PK CRUD queries are present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: DeleteUser :one"))

					// Default List queries are present
					Expect(string(content)).To(ContainSubstring("name: ListUsers :many"))
				}

				// Verify posts.sql FK-based list queries
				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "posts.sql")
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(ContainSubstring("name: ListPostsByUserId :many"))
					Expect(string(content)).NotTo(ContainSubstring("name: ListPostsByTitle :many"))
				}
			})

			It("handles non-existent query names gracefully", func() {
				generator.Config.SQL[0].Codegen[0].Options.Queries = []string{"NonExistentQuery", "AnotherFakeQuery"}
				Expect(generator.Generate()).NotTo(HaveOccurred())

				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "users.sql")
					Expect(path).To(BeAnExistingFile())
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					// Default PK CRUD queries are still present
					Expect(string(content)).To(ContainSubstring("name: GetUser :one"))
					Expect(string(content)).To(ContainSubstring("name: InsertUser :one"))

					// Default List queries are still present
					Expect(string(content)).To(ContainSubstring("name: ListUsers :many"))
				}

				// Verify posts.sql FK-based list queries
				for _, config := range generator.Config.SQL {
					path := filepath.Join(config.Queries, "posts.sql")
					content, err := os.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(ContainSubstring("name: ListPostsByUserId :many"))
					Expect(string(content)).NotTo(ContainSubstring("name: ListPostsByTitle :many"))
				}
			})
		})
	})
})
