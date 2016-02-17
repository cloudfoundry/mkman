package commands_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	outputsDir string
)

var _ = BeforeSuite(func() {
	currentDir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	By("Locating output manifest path")
	outputsDir = filepath.Join(currentDir, "outputs")
})

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}