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
})
