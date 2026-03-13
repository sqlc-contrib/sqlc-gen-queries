package sqlc_test

import (
	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("LoadConfig", func() {
		It("loads a valid config file", func() {
			config, err := sqlc.LoadConfig("./config_test.yaml")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
		})

		It("loads a config file with codegen options", func() {
			config, err := sqlc.LoadConfig("./config_test_exclude.yaml")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.SQL).To(HaveLen(1))
			Expect(config.SQL[0].Codegen).To(HaveLen(1))
			Expect(config.SQL[0].Codegen[0].Plugin).To(Equal("gen-queries"))
			Expect(config.SQL[0].Codegen[0].Out).To(Equal("ent/query"))
			opts := config.SQL[0].GetOptions()
			Expect(opts.Queries).To(HaveLen(3))
			Expect(opts.Queries).To(ContainElements("ListUsers", "CopyUsers", "GetUserByEmail"))
		})

		When("the file does not exist", func() {
			It("returns an error", func() {
				config, err := sqlc.LoadConfig("./config_test.json")
				Expect(err).To(MatchError(ContainSubstring("file does not exist")))
				Expect(config).To(BeNil())
			})
		})

		When("the file is invalid YAML", func() {
			It("returns an error", func() {
				config, err := sqlc.LoadConfig("./config.go")
				Expect(err).To(HaveOccurred())
				Expect(config).To(BeNil())
			})
		})
	})

	Describe("SQL.GetOptions", func() {
		It("returns empty options when codegen is nil", func() {
			sql := sqlc.SQL{}
			opts := sql.GetOptions()
			Expect(opts.Queries).To(BeNil())
		})

		It("returns empty options when no matching plugin", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{Plugin: "other-plugin", Out: "out"},
				},
			}
			opts := sql.GetOptions()
			Expect(opts.Queries).To(BeNil())
		})

		It("returns options for the gen-queries plugin", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "ent/query",
						Options: sqlc.CodegenOptions{
							Queries: []string{"ListUsers", "CopyUsers"},
						},
					},
				},
			}
			opts := sql.GetOptions()
			Expect(opts.Queries).To(HaveLen(2))
			Expect(opts.Queries).To(ContainElements("ListUsers", "CopyUsers"))
		})
	})

	Describe("SQL.GetQueriesSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			querySet := sql.GetQueriesSet()
			Expect(querySet).NotTo(BeNil())
			Expect(querySet).To(BeEmpty())
		})

		It("returns an empty map when queries is empty", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin:  "gen-queries",
						Out:     "ent/query",
						Options: sqlc.CodegenOptions{Queries: []string{}},
					},
				},
			}
			querySet := sql.GetQueriesSet()
			Expect(querySet).NotTo(BeNil())
			Expect(querySet).To(BeEmpty())
		})

		It("returns a map with correct query names", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "ent/query",
						Options: sqlc.CodegenOptions{
							Queries: []string{"ListUsers", "CopyUsers", "GetUserByEmail"},
						},
					},
				},
			}
			querySet := sql.GetQueriesSet()
			Expect(querySet).To(HaveLen(3))
			Expect(querySet["ListUsers"]).To(BeTrue())
			Expect(querySet["CopyUsers"]).To(BeTrue())
			Expect(querySet["GetUserByEmail"]).To(BeTrue())
			Expect(querySet["OtherQuery"]).To(BeFalse())
		})

		It("provides O(1) lookup", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "ent/query",
						Options: sqlc.CodegenOptions{
							Queries: []string{"Query1", "Query2", "Query3"},
						},
					},
				},
			}
			querySet := sql.GetQueriesSet()
			_, exists := querySet["Query2"]
			Expect(exists).To(BeTrue())
			_, exists = querySet["NonExistent"]
			Expect(exists).To(BeFalse())
		})
	})
})
