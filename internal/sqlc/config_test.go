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

		It("loads a config file with skip_queries", func() {
			config, err := sqlc.LoadConfig("./config_test_exclude.yaml")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.SQL).To(HaveLen(1))
			Expect(config.SQL[0].SkipQueries).To(HaveLen(3))
			Expect(config.SQL[0].SkipQueries).To(ContainElements("DeleteUser", "UpdateAuditLog", "GetUserByEmail"))
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

	Describe("SQL.GetSkipQueriesSet", func() {
		It("returns an empty map when skip_queries is nil", func() {
			sql := sqlc.SQL{}
			skipSet := sql.GetSkipQueriesSet()
			Expect(skipSet).NotTo(BeNil())
			Expect(skipSet).To(BeEmpty())
		})

		It("returns an empty map when skip_queries is empty", func() {
			sql := sqlc.SQL{
				SkipQueries: []string{},
			}
			skipSet := sql.GetSkipQueriesSet()
			Expect(skipSet).NotTo(BeNil())
			Expect(skipSet).To(BeEmpty())
		})

		It("returns a map with correct query names", func() {
			sql := sqlc.SQL{
				SkipQueries: []string{"DeleteUser", "UpdateAuditLog", "GetUserByEmail"},
			}
			skipSet := sql.GetSkipQueriesSet()
			Expect(skipSet).To(HaveLen(3))
			Expect(skipSet["DeleteUser"]).To(BeTrue())
			Expect(skipSet["UpdateAuditLog"]).To(BeTrue())
			Expect(skipSet["GetUserByEmail"]).To(BeTrue())
			Expect(skipSet["OtherQuery"]).To(BeFalse())
		})

		It("provides O(1) lookup", func() {
			sql := sqlc.SQL{
				SkipQueries: []string{"Query1", "Query2", "Query3"},
			}
			skipSet := sql.GetSkipQueriesSet()
			_, exists := skipSet["Query2"]
			Expect(exists).To(BeTrue())
			_, exists = skipSet["NonExistent"]
			Expect(exists).To(BeFalse())
		})
	})
})
