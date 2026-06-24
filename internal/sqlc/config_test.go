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
			Expect(opts.Queries.Include).To(HaveLen(2))
			Expect(opts.Queries.Include).To(ContainElements("CopyUsers", "GetUserByEmail"))
			Expect(opts.Tables.Exclude).To(ContainElement("posts"))
			Expect(opts.InsertColumns.Exclude).To(ContainElements("id", "users.created_at", "public.users.updated_at"))
			Expect(opts.UpdateColumns.Exclude).To(ContainElements("created_at", "users.updated_at", "public.users.deleted_at"))
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
			Expect(opts.Queries.Include).To(BeNil())
			Expect(opts.Tables.Exclude).To(BeNil())
		})

		It("returns empty options when no matching plugin", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{Plugin: "other-plugin", Out: "out"},
				},
			}
			opts := sql.GetOptions()
			Expect(opts.Queries.Include).To(BeNil())
			Expect(opts.Tables.Exclude).To(BeNil())
		})

		It("returns options for the gen-queries plugin", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "ent/query",
						Options: sqlc.CodegenOptions{
							Queries: sqlc.QueryOptions{Include: []string{"ListUsers", "CopyUsers"}},
							Tables:  sqlc.TableOptions{Exclude: []string{"audit_logs"}},
							InsertColumns: sqlc.ColumnOptions{
								Exclude: []string{"id", "created_at", "updated_at"},
							},
							UpdateColumns: sqlc.ColumnOptions{
								Exclude: []string{"created_at", "updated_at"},
							},
						},
					},
				},
			}
			opts := sql.GetOptions()
			Expect(opts.Queries.Include).To(HaveLen(2))
			Expect(opts.Queries.Include).To(ContainElements("ListUsers", "CopyUsers"))
			Expect(opts.Tables.Exclude).To(ContainElement("audit_logs"))
			Expect(opts.InsertColumns.Exclude).To(ContainElements("id", "created_at", "updated_at"))
			Expect(opts.UpdateColumns.Exclude).To(ContainElements("created_at", "updated_at"))
		})
	})

	Describe("SQL.GetQueryIncludeSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			querySet := sql.GetQueryIncludeSet()
			Expect(querySet).NotTo(BeNil())
			Expect(querySet).To(BeEmpty())
		})

		It("returns an empty map when include is empty", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin:  "gen-queries",
						Out:     "ent/query",
						Options: sqlc.CodegenOptions{Queries: sqlc.QueryOptions{Include: []string{}}},
					},
				},
			}
			querySet := sql.GetQueryIncludeSet()
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
							Queries: sqlc.QueryOptions{Include: []string{"ListUsers", "CopyUsers", "GetUserByEmail"}},
						},
					},
				},
			}
			querySet := sql.GetQueryIncludeSet()
			Expect(querySet).To(HaveLen(3))
			Expect(querySet["ListUsers"]).To(BeTrue())
			Expect(querySet["CopyUsers"]).To(BeTrue())
			Expect(querySet["GetUserByEmail"]).To(BeTrue())
			Expect(querySet["OtherQuery"]).To(BeFalse())
		})
	})

	Describe("SQL.GetQueryExcludeSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			querySet := sql.GetQueryExcludeSet()
			Expect(querySet).NotTo(BeNil())
			Expect(querySet).To(BeEmpty())
		})

		It("returns a map with excluded query names", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "ent/query",
						Options: sqlc.CodegenOptions{
							Queries: sqlc.QueryOptions{Exclude: []string{"DeleteUser", "BatchDeleteUsers"}},
						},
					},
				},
			}
			querySet := sql.GetQueryExcludeSet()
			Expect(querySet).To(HaveLen(2))
			Expect(querySet["DeleteUser"]).To(BeTrue())
			Expect(querySet["BatchDeleteUsers"]).To(BeTrue())
			Expect(querySet["GetUser"]).To(BeFalse())
		})
	})

	Describe("SQL.GetExcludeSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			excludeSet := sql.GetExcludeSet()
			Expect(excludeSet).NotTo(BeNil())
			Expect(excludeSet).To(BeEmpty())
		})

		It("returns an empty map when exclude is empty", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin:  "gen-queries",
						Out:     "out",
						Options: sqlc.CodegenOptions{Tables: sqlc.TableOptions{Exclude: []string{}}},
					},
				},
			}
			excludeSet := sql.GetExcludeSet()
			Expect(excludeSet).NotTo(BeNil())
			Expect(excludeSet).To(BeEmpty())
		})

		It("returns a map with excluded table names", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "out",
						Options: sqlc.CodegenOptions{
							Tables: sqlc.TableOptions{Exclude: []string{"users", "analytics.events"}},
						},
					},
				},
			}
			excludeSet := sql.GetExcludeSet()
			Expect(excludeSet).To(HaveLen(2))
			Expect(excludeSet["users"]).To(BeTrue())
			Expect(excludeSet["analytics.events"]).To(BeTrue())
			Expect(excludeSet["posts"]).To(BeFalse())
		})
	})

	Describe("SQL.GetIncludeSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			includeSet := sql.GetIncludeSet()
			Expect(includeSet).NotTo(BeNil())
			Expect(includeSet).To(BeEmpty())
		})

		It("returns an empty map when include is empty", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin:  "gen-queries",
						Out:     "out",
						Options: sqlc.CodegenOptions{Tables: sqlc.TableOptions{Include: []string{}}},
					},
				},
			}
			includeSet := sql.GetIncludeSet()
			Expect(includeSet).NotTo(BeNil())
			Expect(includeSet).To(BeEmpty())
		})

		It("returns a map with included table names", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "out",
						Options: sqlc.CodegenOptions{
							Tables: sqlc.TableOptions{Include: []string{"users", "analytics.events"}},
						},
					},
				},
			}
			includeSet := sql.GetIncludeSet()
			Expect(includeSet).To(HaveLen(2))
			Expect(includeSet["users"]).To(BeTrue())
			Expect(includeSet["analytics.events"]).To(BeTrue())
			Expect(includeSet["posts"]).To(BeFalse())
		})
	})

	Describe("SQL.GetInsertColumnExcludeSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			excludeSet := sql.GetInsertColumnExcludeSet()
			Expect(excludeSet).NotTo(BeNil())
			Expect(excludeSet).To(BeEmpty())
		})

		It("returns a map with excluded insert column names", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "out",
						Options: sqlc.CodegenOptions{
							InsertColumns: sqlc.ColumnOptions{
								Exclude: []string{"id", "todos.created_at", "public.todos.updated_at"},
							},
						},
					},
				},
			}
			excludeSet := sql.GetInsertColumnExcludeSet()
			Expect(excludeSet).To(HaveLen(3))
			Expect(excludeSet["id"]).To(BeTrue())
			Expect(excludeSet["todos.created_at"]).To(BeTrue())
			Expect(excludeSet["public.todos.updated_at"]).To(BeTrue())
			Expect(excludeSet["title"]).To(BeFalse())
		})
	})

	Describe("SQL.GetUpdateColumnExcludeSet", func() {
		It("returns an empty map when codegen is nil", func() {
			sql := sqlc.SQL{}
			excludeSet := sql.GetUpdateColumnExcludeSet()
			Expect(excludeSet).NotTo(BeNil())
			Expect(excludeSet).To(BeEmpty())
		})

		It("returns a map with excluded update column names", func() {
			sql := sqlc.SQL{
				Codegen: []sqlc.Codegen{
					{
						Plugin: "gen-queries",
						Out:    "out",
						Options: sqlc.CodegenOptions{
							UpdateColumns: sqlc.ColumnOptions{
								Exclude: []string{"created_at", "todos.updated_at", "public.todos.user_id"},
							},
						},
					},
				},
			}
			excludeSet := sql.GetUpdateColumnExcludeSet()
			Expect(excludeSet).To(HaveLen(3))
			Expect(excludeSet["created_at"]).To(BeTrue())
			Expect(excludeSet["todos.updated_at"]).To(BeTrue())
			Expect(excludeSet["public.todos.user_id"]).To(BeTrue())
			Expect(excludeSet["title"]).To(BeFalse())
		})
	})
})
