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
	})
})
