package sqlc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSqlc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQLC Suite")
}
