package template_test

import (
	"maps"

	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {
	Describe("Open", func() {
		opts := map[string]any{
			"where":    func(args ...any) string { return "" },
			"arg":      func(args ...any) string { return "" },
			"by_one":   func(args ...any) string { return "" },
			"by_many":  func(args ...any) string { return "" },
			"by_index": func(args ...any) string { return "" },
			"cursor":   func(args ...any) string { return "" },
			"order":    func(args ...any) string { return "" },
		}

		It("opens and parses template successfully", func() {
			file, err := template.Open("template.sql.tmpl", opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})

		It("supports built-in functions", func() {
			file, err := template.Open("template.sql.tmpl", opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})

		It("supports custom functions", func() {
			customFunc := map[string]any{
				"custom": func(s string) string {
					return "custom_" + s
				},
			}

			maps.Copy(customFunc, opts)

			file, err := template.Open("template.sql.tmpl", customFunc)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
		})

		When("the template does not exist", func() {
			It("returns an error", func() {
				_, err := template.Open("nonexistent.tmpl")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
