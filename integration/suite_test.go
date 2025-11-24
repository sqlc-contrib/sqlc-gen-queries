package integration_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var command *exec.Cmd

var _ = BeforeSuite(func() {
	By("Building the copay-coupon-api", func() {
		path, err := gexec.Build("github.com/sqlc-contrib/sqlc-gen-queries/cmd/sqlc-gen-queries")
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(gexec.CleanupBuildArtifacts)
		Expect(path).NotTo(BeEmpty())
		// prepare the command
		command = exec.Command(path)
	})
})
