package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		command.Dir, _ = os.MkdirTemp("", "sqcl-gen-test-*")
		command.Env = append(command.Env, GetFileEnv("SQLC_CONFIG_FILE", "../internal/sqlc/config_test.yaml"))
		command.Env = append(command.Env, GetFileEnv("SQLC_CATALOG_FILE", "../internal/sqlc/catalog_test.json"))
	})
})

// GetFileEnv is a helper to set environment variables from file contents in tests
func GetFileEnv(name, key string) string {
	dir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	return fmt.Sprintf("%v=%v", name, filepath.Join(dir, key))
}
